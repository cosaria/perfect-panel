package user

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type ResetUserSubscribeTrafficInput struct {
	Body types.ResetUserSubscribeTrafficRequest
}

func ResetUserSubscribeTrafficHandler(deps Deps) func(context.Context, *ResetUserSubscribeTrafficInput) (*struct{}, error) {
	return func(ctx context.Context, input *ResetUserSubscribeTrafficInput) (*struct{}, error) {
		l := NewResetUserSubscribeTrafficLogic(ctx, deps)
		if err := l.ResetUserSubscribeTraffic(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type ResetUserSubscribeTrafficLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewResetUserSubscribeTrafficLogic Reset user subscribe traffic
func NewResetUserSubscribeTrafficLogic(ctx context.Context, deps Deps) *ResetUserSubscribeTrafficLogic {
	return &ResetUserSubscribeTrafficLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *ResetUserSubscribeTrafficLogic) ResetUserSubscribeTraffic(req *types.ResetUserSubscribeTrafficRequest) error {
	userSub, err := l.deps.UserModel.FindOneSubscribe(l.ctx, req.UserSubscribeId)
	if err != nil {
		l.Errorw("FindOneSubscribe error", logger.Field("error", err.Error()), logger.Field("userSubscribeId", req.UserSubscribeId))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), " FindOneSubscribe error: %v", err.Error())
	}
	userSub.Download = 0
	userSub.Upload = 0

	err = l.deps.UserModel.UpdateSubscribe(l.ctx, userSub)
	if err != nil {
		l.Errorw("UpdateSubscribe error", logger.Field("error", err.Error()), logger.Field("userSubscribeId", req.UserSubscribeId))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), " UpdateSubscribe error: %v", err.Error())
	}
	// Clear user subscribe cache
	if err = l.deps.UserModel.ClearSubscribeCache(l.ctx, userSub); err != nil {
		l.Errorw("ClearSubscribeCache failed:", logger.Field("error", err.Error()), logger.Field("userSubscribeId", userSub.Id))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "ClearSubscribeCache failed: %v", err.Error())
	}
	// Clear subscribe cache
	if err = l.deps.SubscribeModel.ClearCache(l.ctx, userSub.SubscribeId); err != nil {
		l.Errorw("failed to clear subscribe cache", logger.Field("error", err.Error()), logger.Field("subscribeId", userSub.SubscribeId))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "failed to clear subscribe cache: %v", err.Error())
	}
	return nil
}
