package user

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type DeleteUserDeviceInput struct {
	Body types.DeleteUserDeviceRequest
}

func DeleteUserDeviceHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteUserDeviceInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteUserDeviceInput) (*struct{}, error) {
		l := NewDeleteUserDeviceLogic(ctx, svcCtx)
		if err := l.DeleteUserDevice(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeleteUserDeviceLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete user device
func NewDeleteUserDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserDeviceLogic {
	return &DeleteUserDeviceLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserDeviceLogic) DeleteUserDevice(req *types.DeleteUserDeviceRequest) error {
	err := l.svcCtx.UserModel.DeleteDevice(l.ctx, req.Id)
	if err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete user error: %v", err.Error())
	}
	return nil
}
