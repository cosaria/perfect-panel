package subscribe

import (
	"context"
	"github.com/perfect-panel/server/models/subscribe"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type BatchDeleteSubscribeGroupInput struct {
	Body types.BatchDeleteSubscribeGroupRequest
}

func BatchDeleteSubscribeGroupHandler(svcCtx *svc.ServiceContext) func(context.Context, *BatchDeleteSubscribeGroupInput) (*struct{}, error) {
	return func(ctx context.Context, input *BatchDeleteSubscribeGroupInput) (*struct{}, error) {
		l := NewBatchDeleteSubscribeGroupLogic(ctx, svcCtx)
		if err := l.BatchDeleteSubscribeGroup(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type BatchDeleteSubscribeGroupLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Batch delete subscribe group
func NewBatchDeleteSubscribeGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteSubscribeGroupLogic {
	return &BatchDeleteSubscribeGroupLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteSubscribeGroupLogic) BatchDeleteSubscribeGroup(req *types.BatchDeleteSubscribeGroupRequest) error {
	err := l.svcCtx.DB.Model(&subscribe.Group{}).Where("id in ?", req.Ids).Delete(&subscribe.Group{}).Error
	if err != nil {
		l.Logger.Error("[BatchDeleteSubscribeGroup] Delete Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "batch delete subscribe group failed: %v", err.Error())
	}
	return nil
}
