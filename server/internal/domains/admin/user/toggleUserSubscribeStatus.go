package user

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type ToggleUserSubscribeStatusInput struct {
	Body types.ToggleUserSubscribeStatusRequest
}

func ToggleUserSubscribeStatusHandler(deps Deps) func(context.Context, *ToggleUserSubscribeStatusInput) (*struct{}, error) {
	return func(ctx context.Context, input *ToggleUserSubscribeStatusInput) (*struct{}, error) {
		l := NewToggleUserSubscribeStatusLogic(ctx, deps)
		if err := l.ToggleUserSubscribeStatus(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type ToggleUserSubscribeStatusLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewToggleUserSubscribeStatusLogic Stop user subscribe
func NewToggleUserSubscribeStatusLogic(ctx context.Context, deps Deps) *ToggleUserSubscribeStatusLogic {
	return &ToggleUserSubscribeStatusLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *ToggleUserSubscribeStatusLogic) ToggleUserSubscribeStatus(req *types.ToggleUserSubscribeStatusRequest) error {
	userSub, err := l.deps.UserModel.FindOneSubscribe(l.ctx, req.UserSubscribeId)
	if err != nil {
		l.Errorw("FindOneSubscribe error", logger.Field("error", err.Error()), logger.Field("userSubscribeId", req.UserSubscribeId))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), " FindOneSubscribe error: %v", err.Error())
	}

	switch userSub.Status {
	case 2: // active
		userSub.Status = 5 // set status to stopped
	case 5: // stopped
		userSub.Status = 2 // set status to active
	default:
		l.Errorw("invalid user subscribe status", logger.Field("userSubscribeId", req.UserSubscribeId), logger.Field("status", userSub.Status))
		return errors.Wrapf(xerr.NewErrCodeMsg(xerr.ERROR, "invalid subscribe status"), "invalid user subscribe status: %d", userSub.Status)
	}

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
