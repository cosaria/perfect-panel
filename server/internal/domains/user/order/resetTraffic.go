package order

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/config"
	queue "github.com/perfect-panel/server/internal/jobs/spec"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/models/order"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type ResetTrafficInput struct {
	Body types.ResetTrafficOrderRequest
}

type ResetTrafficOutput struct {
	Body *types.ResetTrafficOrderResponse
}

func ResetTrafficHandler(deps Deps) func(context.Context, *ResetTrafficInput) (*ResetTrafficOutput, error) {
	return func(ctx context.Context, input *ResetTrafficInput) (*ResetTrafficOutput, error) {
		l := NewResetTrafficLogic(ctx, deps)
		resp, err := l.ResetTraffic(&input.Body)
		if err != nil {
			return nil, err
		}
		return &ResetTrafficOutput{Body: resp}, nil
	}
}

type ResetTrafficLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Reset traffic
func NewResetTrafficLogic(ctx context.Context, deps Deps) *ResetTrafficLogic {
	return &ResetTrafficLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *ResetTrafficLogic) ResetTraffic(req *types.ResetTrafficOrderRequest) (resp *types.ResetTrafficOrderResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	// find user subscription
	userSubscribe, err := l.deps.UserModel.FindOneUserSubscribe(l.ctx, req.UserSubscribeID)
	if err != nil {
		l.Errorw("[ResetTraffic] Database query error", logger.Field("error", err.Error()), logger.Field("UserSubscribeID", req.UserSubscribeID))
		return nil, errors.Wrapf(err, "find user subscribe error: %v", err.Error())
	}
	if userSubscribe.Subscribe == nil {
		l.Errorw("[ResetTraffic] subscribe not found", logger.Field("UserSubscribeID", req.UserSubscribeID))
		return nil, errors.New("subscribe not found")
	}
	amount := userSubscribe.Subscribe.Replacement
	var deductionAmount int64
	// Check user deduction amount
	if u.GiftAmount > 0 {
		if u.GiftAmount >= amount {
			deductionAmount = amount
			amount = 0
			u.GiftAmount -= amount
		} else {
			deductionAmount = u.GiftAmount
			amount -= u.GiftAmount
			u.GiftAmount = 0
		}
	}
	// find payment method
	payment, err := l.deps.PaymentModel.FindOne(l.ctx, req.Payment)
	if err != nil {
		l.Errorw("[ResetTraffic] Database query error", logger.Field("error", err.Error()), logger.Field("payment", req.Payment))
		return nil, errors.Wrapf(err, "find payment error: %v", err.Error())
	}
	var feeAmount int64
	// Calculate the handling fee
	if amount > 0 {
		feeAmount = calculateFee(amount, payment)
	}
	// create order
	orderInfo := order.Order{
		Id:             0,
		ParentId:       userSubscribe.OrderId,
		UserId:         u.Id,
		OrderNo:        tool.GenerateTradeNo(),
		Type:           3,
		Price:          userSubscribe.Subscribe.Replacement,
		Amount:         amount + feeAmount,
		GiftAmount:     deductionAmount,
		FeeAmount:      feeAmount,
		PaymentId:      payment.Id,
		Method:         payment.Platform,
		Status:         1,
		SubscribeId:    userSubscribe.SubscribeId,
		SubscribeToken: userSubscribe.Token,
	}
	// Database transaction
	err = l.deps.DB.Transaction(func(db *gorm.DB) error {
		// update user deduction && Pre deduction ,Return after canceling the order
		if orderInfo.GiftAmount > 0 {
			// update user deduction && Pre deduction ,Return after canceling the order
			if err := l.deps.UserModel.Update(l.ctx, u, db); err != nil {
				l.Errorw("[ResetTraffic] Database update error", logger.Field("error", err.Error()), logger.Field("user", u))
				return err
			}
			// create deduction record
			giftLog := log.Gift{
				Type:        log.GiftTypeReduce,
				OrderNo:     orderInfo.OrderNo,
				SubscribeId: 0,
				Amount:      orderInfo.GiftAmount,
				Balance:     u.GiftAmount,
				Remark:      "Renewal order deduction",
				Timestamp:   time.Now().UnixMilli(),
			}
			content, _ := giftLog.Marshal()

			if err = db.Model(&log.SystemLog{}).Create(&log.SystemLog{
				Type:     log.TypeGift.Uint8(),
				Date:     time.Now().Format(time.DateOnly),
				ObjectID: u.Id,
				Content:  string(content),
			}).Error; err != nil {
				l.Errorw("[ResetTraffic] Database insert error", logger.Field("error", err.Error()), logger.Field("deductionLog", content))
				return err
			}
		}
		// insert order
		return db.Model(&order.Order{}).Create(&orderInfo).Error
	})
	if err != nil {
		l.Errorw("[ResetTraffic] Database insert error", logger.Field("error", err.Error()), logger.Field("order", orderInfo))
		return nil, errors.Wrapf(err, "insert order error: %v", err.Error())
	}
	// Deferred task
	payload := queue.DeferCloseOrderPayload{
		OrderNo: orderInfo.OrderNo,
	}
	val, err := json.Marshal(payload)
	if err != nil {
		l.Errorw("[ResetTraffic] Marshal payload error", logger.Field("error", err.Error()), logger.Field("payload", payload))
	}
	task := asynq.NewTask(queue.DeferCloseOrder, val, asynq.MaxRetry(3))
	taskInfo, err := l.deps.Queue.Enqueue(task, asynq.ProcessIn(CloseOrderTimeMinutes*time.Minute))
	if err != nil {
		l.Errorw("[ResetTraffic] Enqueue task error", logger.Field("error", err.Error()), logger.Field("task", task))
	} else {
		l.Infow("[ResetTraffic] Enqueue task success", logger.Field("TaskID", taskInfo.ID))
	}
	return &types.ResetTrafficOrderResponse{
		OrderNo: orderInfo.OrderNo,
	}, nil
}
