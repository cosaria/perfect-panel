package user

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type DeleteUserSubscribeInput struct {
	Body types.DeleteUserSubscribeRequest
}

func DeleteUserSubscribeHandler(deps Deps) func(context.Context, *DeleteUserSubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteUserSubscribeInput) (*struct{}, error) {
		l := NewDeleteUserSubscribeLogic(ctx, deps)
		if err := l.DeleteUserSubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeleteUserSubscribeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewDeleteUserSubscribeLogic Delete user subcribe
func NewDeleteUserSubscribeLogic(ctx context.Context, deps Deps) *DeleteUserSubscribeLogic {
	return &DeleteUserSubscribeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *DeleteUserSubscribeLogic) DeleteUserSubscribe(req *types.DeleteUserSubscribeRequest) error {
	// find user subscribe by ID
	userSubscribe, err := l.deps.UserModel.FindOneSubscribe(l.ctx, req.UserSubscribeId)
	if err != nil {
		l.Errorw("failed to find user subscribe", logger.Field("error", err.Error()), logger.Field("userSubscribeId", req.UserSubscribeId))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "failed to find user subscribe: %v", err.Error())
	}

	err = l.deps.UserModel.DeleteSubscribeById(l.ctx, req.UserSubscribeId)
	if err != nil {
		l.Errorw("failed to delete user subscribe", logger.Field("error", err.Error()), logger.Field("userSubscribeId", req.UserSubscribeId))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "failed to delete user subscribe: %v", err.Error())
	}
	// Clear user subscribe cache
	if err = l.deps.UserModel.ClearSubscribeCache(l.ctx, userSubscribe); err != nil {
		l.Errorw("failed to clear user subscribe cache", logger.Field("error", err.Error()), logger.Field("userSubscribeId", req.UserSubscribeId))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "failed to clear user subscribe cache: %v", err.Error())
	}
	// Clear subscribe cache
	if err = l.deps.SubscribeModel.ClearCache(l.ctx, userSubscribe.SubscribeId); err != nil {
		l.Errorw("failed to clear subscribe cache", logger.Field("error", err.Error()), logger.Field("subscribeId", userSubscribe.SubscribeId))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "failed to clear subscribe cache: %v", err.Error())
	}
	return nil
}
