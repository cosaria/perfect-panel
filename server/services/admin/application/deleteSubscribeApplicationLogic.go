package application

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type DeleteSubscribeApplicationLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteSubscribeApplicationLogic Delete subscribe application
func NewDeleteSubscribeApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSubscribeApplicationLogic {
	return &DeleteSubscribeApplicationLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSubscribeApplicationLogic) DeleteSubscribeApplication(req *types.DeleteSubscribeApplicationRequest) error {
	err := l.svcCtx.ClientModel.Delete(l.ctx, req.Id)
	if err != nil {
		l.Errorf("Failed to delete subscribe application with ID %d: %v", req.Id, err)
		return errors.Wrap(xerr.NewErrCode(xerr.DatabaseDeletedError), err.Error())
	}
	return nil
}
