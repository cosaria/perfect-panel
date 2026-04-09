package user

import (
	"context"

	"time"

	"github.com/google/uuid"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/order"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/uuidx"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type ResetUserSubscribeTokenInput struct {
	Body types.ResetUserSubscribeTokenRequest
}

func ResetUserSubscribeTokenHandler(deps Deps) func(context.Context, *ResetUserSubscribeTokenInput) (*struct{}, error) {
	return func(ctx context.Context, input *ResetUserSubscribeTokenInput) (*struct{}, error) {
		l := NewResetUserSubscribeTokenLogic(ctx, deps)
		if err := l.ResetUserSubscribeToken(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type ResetUserSubscribeTokenLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewResetUserSubscribeTokenLogic Reset User Subscribe Token
func NewResetUserSubscribeTokenLogic(ctx context.Context, deps Deps) *ResetUserSubscribeTokenLogic {
	return &ResetUserSubscribeTokenLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *ResetUserSubscribeTokenLogic) ResetUserSubscribeToken(req *types.ResetUserSubscribeTokenRequest) error {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	userSub, err := l.deps.UserModel.FindOneUserSubscribe(l.ctx, req.UserSubscribeId)
	if err != nil {
		l.Errorw("FindOneUserSubscribe failed:", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindOneUserSubscribe failed: %v", err.Error())
	}
	if userSub.UserId != u.Id {
		l.Errorw("UserSubscribeId does not belong to the current user")
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "UserSubscribeId does not belong to the current user")
	}

	var orderDetails *order.Details
	// find order
	if userSub.OrderId != 0 {
		orderDetails, err = l.deps.OrderModel.FindOneDetails(l.ctx, userSub.OrderId)
		if err != nil {
			l.Errorw("FindOneDetails failed:", logger.Field("error", err.Error()))
			return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindOneDetails failed: %v", err.Error())
		}
	} else {
		// if order id is 0, this a admin create user subscribe
		orderDetails = &order.Details{}
	}

	userSub.Token = uuidx.SubscribeToken(orderDetails.OrderNo + time.Now().Format("20060102150405.000"))
	userSub.UUID = uuid.New().String()
	var newSub user.Subscribe
	tool.DeepCopy(&newSub, userSub)

	err = l.deps.UserModel.UpdateSubscribe(l.ctx, &newSub)
	if err != nil {
		l.Errorw("UpdateSubscribe failed:", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "UpdateSubscribe failed: %v", err.Error())
	}
	//clear user subscription cache
	if err = l.deps.UserModel.ClearSubscribeCache(l.ctx, &newSub); err != nil {
		l.Errorw("ClearSubscribeCache failed", logger.Field("error", err.Error()), logger.Field("userSubscribeId", userSub.Id))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "ClearSubscribeCache failed: %v", err.Error())
	}
	// Clear subscription cache
	if err = l.deps.SubscribeModel.ClearCache(l.ctx, userSub.SubscribeId); err != nil {
		l.Errorw("ClearSubscribeCache failed", logger.Field("error", err.Error()), logger.Field("subscribeId", userSub.SubscribeId))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "ClearSubscribeCache failed: %v", err.Error())
	}

	return nil
}
