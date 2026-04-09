package user

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/pkg/errors"
)

type UpdateUserSubscribeNoteInput struct {
	Body types.UpdateUserSubscribeNoteRequest
}

func UpdateUserSubscribeNoteHandler(deps Deps) func(context.Context, *UpdateUserSubscribeNoteInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserSubscribeNoteInput) (*struct{}, error) {
		l := NewUpdateUserSubscribeNoteLogic(ctx, deps)
		if err := l.UpdateUserSubscribeNote(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateUserSubscribeNoteLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewUpdateUserSubscribeNoteLogic Update User Subscribe Note
func NewUpdateUserSubscribeNoteLogic(ctx context.Context, deps Deps) *UpdateUserSubscribeNoteLogic {
	return &UpdateUserSubscribeNoteLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateUserSubscribeNoteLogic) UpdateUserSubscribeNote(req *types.UpdateUserSubscribeNoteRequest) error {
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

	userSub.Note = req.Note
	var newSub user.Subscribe
	tool.DeepCopy(&newSub, userSub)

	err = l.deps.UserModel.UpdateSubscribe(l.ctx, &newSub)
	if err != nil {
		l.Errorw("UpdateSubscribe failed:", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "UpdateSubscribe failed: %v", err.Error())
	}

	// Clear user subscription cache
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
