package user

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type KickOfflineByUserDeviceInput struct {
	Body types.KickOfflineRequest
}

func KickOfflineByUserDeviceHandler(svcCtx *svc.ServiceContext) func(context.Context, *KickOfflineByUserDeviceInput) (*struct{}, error) {
	return func(ctx context.Context, input *KickOfflineByUserDeviceInput) (*struct{}, error) {
		l := NewKickOfflineByUserDeviceLogic(ctx, svcCtx)
		if err := l.KickOfflineByUserDevice(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type KickOfflineByUserDeviceLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// kick offline user device
func NewKickOfflineByUserDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickOfflineByUserDeviceLogic {
	return &KickOfflineByUserDeviceLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *KickOfflineByUserDeviceLogic) KickOfflineByUserDevice(req *types.KickOfflineRequest) error {
	device, err := l.svcCtx.UserModel.FindOneDevice(l.ctx, req.Id)
	if err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get Device  error: %v", err.Error())
	}
	l.svcCtx.DeviceManager.KickDevice(device.UserId, device.Identifier)
	device.Online = false
	err = l.svcCtx.UserModel.UpdateDevice(l.ctx, device)
	if err != nil {
		l.Logger.Error("[KickOfflineByUserDeviceLogic] Update Device Error:", logger.Field("err", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update Device error: %v", err.Error())
	}

	return nil
}
