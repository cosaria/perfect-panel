package user

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type KickOfflineByUserDeviceInput struct {
	Body types.KickOfflineRequest
}

func KickOfflineByUserDeviceHandler(deps Deps) func(context.Context, *KickOfflineByUserDeviceInput) (*struct{}, error) {
	return func(ctx context.Context, input *KickOfflineByUserDeviceInput) (*struct{}, error) {
		l := NewKickOfflineByUserDeviceLogic(ctx, deps)
		if err := l.KickOfflineByUserDevice(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type KickOfflineByUserDeviceLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// kick offline user device
func NewKickOfflineByUserDeviceLogic(ctx context.Context, deps Deps) *KickOfflineByUserDeviceLogic {
	return &KickOfflineByUserDeviceLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *KickOfflineByUserDeviceLogic) KickOfflineByUserDevice(req *types.KickOfflineRequest) error {
	device, err := l.deps.UserModel.FindOneDevice(l.ctx, req.Id)
	if err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get Device  error: %v", err.Error())
	}
	l.deps.DeviceManager.KickDevice(device.UserId, device.Identifier)
	device.Online = false
	err = l.deps.UserModel.UpdateDevice(l.ctx, device)
	if err != nil {
		l.Logger.Error("[KickOfflineByUserDeviceLogic] Update Device Error:", logger.Field("err", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update Device error: %v", err.Error())
	}

	return nil
}
