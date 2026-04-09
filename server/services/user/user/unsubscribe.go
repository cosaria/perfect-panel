package user

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type UnsubscribeInput struct {
	Body types.UnsubscribeRequest
}

func UnsubscribeHandler(deps Deps) func(context.Context, *UnsubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *UnsubscribeInput) (*struct{}, error) {
		l := NewUnsubscribeLogic(ctx, deps)
		if err := l.Unsubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UnsubscribeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewUnsubscribeLogic creates a new instance of UnsubscribeLogic for handling subscription cancellation
func NewUnsubscribeLogic(ctx context.Context, deps Deps) *UnsubscribeLogic {
	return &UnsubscribeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

// Unsubscribe handles the subscription cancellation process with proper refund distribution
// It prioritizes refunding to gift amount for balance-paid orders, then to regular balance
func (l *UnsubscribeLogic) Unsubscribe(req *types.UnsubscribeRequest) error {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}

	// find user subscription by ID
	userSub, err := l.deps.UserModel.FindOneSubscribe(l.ctx, req.Id)
	if err != nil {
		l.Errorw("FindOneSubscribe failed", logger.Field("error", err.Error()), logger.Field("reqId", req.Id))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindOneSubscribe failed: %v", err.Error())
	}

	activate := []uint8{0, 1, 2}

	if !tool.Contains(activate, userSub.Status) {
		// Only active (2) or paused (5) subscriptions can be cancelled
		l.Errorw("Subscription status invalid for cancellation", logger.Field("userSubscribeId", userSub.Id), logger.Field("status", userSub.Status))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Subscription status invalid for cancellation")
	}

	// Calculate the remaining amount to refund based on unused subscription time/traffic
	remainingAmount, err := CalculateRemainingAmount(l.ctx, l.deps, req.Id)
	if err != nil {
		return err
	}

	// Process unsubscription in a database transaction to ensure data consistency
	err = l.deps.UserModel.Transaction(l.ctx, func(db *gorm.DB) error {
		// Find and update subscription status to cancelled (status = 4)
		userSub.Status = 4 // Set status to cancelled
		if err = l.deps.UserModel.UpdateSubscribe(l.ctx, userSub); err != nil {
			return err
		}

		// Query the original order information to determine refund strategy
		orderInfo, err := l.deps.OrderModel.FindOne(l.ctx, userSub.OrderId)
		if err != nil {
			return err
		}
		// Calculate refund distribution based on payment method and gift amount priority
		var balance, gift int64
		if orderInfo.Method == "balance" {
			// For balance-paid orders, prioritize refunding to gift amount first
			if orderInfo.GiftAmount >= remainingAmount {
				// Gift amount covers the entire refund - refund all to gift balance
				gift = remainingAmount
				balance = u.Balance // Regular balance remains unchanged
			} else {
				// Gift amount insufficient - refund to gift first, remainder to regular balance
				gift = orderInfo.GiftAmount
				balance = u.Balance + (remainingAmount - orderInfo.GiftAmount)
			}
		} else {
			// For non-balance payment orders, refund entirely to regular balance
			balance = remainingAmount + u.Balance
			gift = 0
		}

		// Create balance log entry only if there's an actual regular balance refund
		balanceRefundAmount := balance - u.Balance
		if balanceRefundAmount > 0 {
			balanceLog := log.Balance{
				OrderNo:   orderInfo.OrderNo,
				Amount:    balanceRefundAmount,
				Type:      log.BalanceTypeRefund, // Type 4 represents refund transaction
				Balance:   balance,
				Timestamp: time.Now().UnixMilli(),
			}
			content, _ := balanceLog.Marshal()

			if err := db.Model(&log.SystemLog{}).Create(&log.SystemLog{
				Type:     log.TypeBalance.Uint8(),
				Date:     time.Now().Format(time.DateOnly),
				ObjectID: u.Id,
				Content:  string(content),
			}).Error; err != nil {
				return err
			}
		}

		// Create gift amount log entry if there's a gift balance refund
		if gift > 0 {

			giftLog := log.Gift{
				SubscribeId: userSub.Id,
				OrderNo:     orderInfo.OrderNo,
				Type:        log.GiftTypeIncrease, // Type 1 represents gift amount increase
				Amount:      gift,
				Balance:     u.GiftAmount + gift,
				Remark:      "Unsubscribe refund",
			}
			content, _ := giftLog.Marshal()

			if err := db.Model(&log.SystemLog{}).Create(&log.SystemLog{
				Type:     log.TypeGift.Uint8(),
				Date:     time.Now().Format(time.DateOnly),
				ObjectID: u.Id,
				Content:  string(content),
			}).Error; err != nil {
				return err
			}
			// Update user's gift amount
			u.GiftAmount += gift
		}

		// Update user's regular balance and save changes to database
		u.Balance = balance
		return l.deps.UserModel.Update(l.ctx, u)
	})

	if err != nil {
		l.Errorw("Unsubscribe transaction failed", logger.Field("error", err.Error()), logger.Field("userId", u.Id), logger.Field("reqId", req.Id))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "Unsubscribe transaction failed: %v", err.Error())
	}

	//clear user subscription cache
	if err = l.deps.UserModel.ClearSubscribeCache(l.ctx, userSub); err != nil {
		l.Errorw("ClearSubscribeCache failed", logger.Field("error", err.Error()), logger.Field("userSubscribeId", userSub.Id))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "ClearSubscribeCache failed: %v", err.Error())
	}
	// Clear subscription cache
	if err = l.deps.SubscribeModel.ClearCache(l.ctx, userSub.SubscribeId); err != nil {
		l.Errorw("ClearSubscribeCache failed", logger.Field("error", err.Error()), logger.Field("subscribeId", userSub.SubscribeId))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "ClearSubscribeCache failed: %v", err.Error())
	}

	return err
}
