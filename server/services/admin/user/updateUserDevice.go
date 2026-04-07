package user

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type UpdateUserDeviceInput struct {
	Body types.UserDevice
}

func UpdateUserDeviceHandler(deps Deps) func(context.Context, *UpdateUserDeviceInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserDeviceInput) (*struct{}, error) {
		l := NewUpdateUserDeviceLogic(ctx, deps)
		if err := l.UpdateUserDevice(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateUserDeviceLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// User device
func NewUpdateUserDeviceLogic(ctx context.Context, deps Deps) *UpdateUserDeviceLogic {
	return &UpdateUserDeviceLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateUserDeviceLogic) UpdateUserDevice(req *types.UserDevice) error {
	device, err := l.deps.UserModel.FindOneDevice(l.ctx, req.Id)
	if err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get Device  error: %v", err.Error())
	}
	device.Enabled = req.Enabled
	err = l.deps.UserModel.UpdateDevice(l.ctx, device)
	if err != nil {
		l.Logger.Error("[UpdateUserDeviceLogic] Update Device Error:", logger.Field("err", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update Device error: %v", err.Error())
	}
	return nil
}
