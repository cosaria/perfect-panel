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

type DeleteSubscribeGroupInput struct {
	Body types.DeleteSubscribeGroupRequest
}

func DeleteSubscribeGroupHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteSubscribeGroupInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteSubscribeGroupInput) (*struct{}, error) {
		l := NewDeleteSubscribeGroupLogic(ctx, svcCtx)
		if err := l.DeleteSubscribeGroup(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeleteSubscribeGroupLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete subscribe group
func NewDeleteSubscribeGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSubscribeGroupLogic {
	return &DeleteSubscribeGroupLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSubscribeGroupLogic) DeleteSubscribeGroup(req *types.DeleteSubscribeGroupRequest) error {
	err := l.svcCtx.DB.Model(&subscribe.Group{}).Where("id = ?", req.Id).Delete(&subscribe.Group{}).Error
	if err != nil {
		l.Logger.Error("[DeleteSubscribeGroupLogic] delete subscribe group failed: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete subscribe group failed: %v", err.Error())
	}
	return nil
}
