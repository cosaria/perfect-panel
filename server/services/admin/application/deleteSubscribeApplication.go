package application

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type DeleteSubscribeApplicationInput struct {
	Body types.DeleteSubscribeApplicationRequest
}

func DeleteSubscribeApplicationHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteSubscribeApplicationInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteSubscribeApplicationInput) (*struct{}, error) {
		l := NewDeleteSubscribeApplicationLogic(ctx, svcCtx)
		if err := l.DeleteSubscribeApplication(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

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
