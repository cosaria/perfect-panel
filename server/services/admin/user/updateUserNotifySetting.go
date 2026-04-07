package user

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type UpdateUserNotifySettingInput struct {
	Body types.UpdateUserNotifySettingRequest
}

func UpdateUserNotifySettingHandler(deps Deps) func(context.Context, *UpdateUserNotifySettingInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserNotifySettingInput) (*struct{}, error) {
		l := NewUpdateUserNotifySettingLogic(ctx, deps)
		if err := l.UpdateUserNotifySetting(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateUserNotifySettingLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewUpdateUserNotifySettingLogic Update user notify setting
func NewUpdateUserNotifySettingLogic(ctx context.Context, deps Deps) *UpdateUserNotifySettingLogic {
	return &UpdateUserNotifySettingLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateUserNotifySettingLogic) UpdateUserNotifySetting(req *types.UpdateUserNotifySettingRequest) error {
	userInfo, err := l.deps.UserModel.FindOne(l.ctx, req.UserId)
	if err != nil {
		l.Errorw("[UpdateUserNotifySettingLogic] Find User Error:", logger.Field("err", err.Error()), logger.Field("userId", req.UserId))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Find User Error")
	}
	tool.DeepCopy(userInfo, req)
	err = l.deps.UserModel.Update(l.ctx, userInfo)
	if err != nil {
		l.Errorw("[UpdateUserNotifySettingLogic] Update User Error:", logger.Field("err", err.Error()), logger.Field("userId", req.UserId))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "Update User Error")
	}
	return nil
}
