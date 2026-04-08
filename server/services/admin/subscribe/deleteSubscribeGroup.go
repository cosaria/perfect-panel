package subscribe

import (
	"context"
	"github.com/perfect-panel/server/models/subscribe"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type DeleteSubscribeGroupInput struct {
	Body types.DeleteSubscribeGroupRequest
}

func DeleteSubscribeGroupHandler(deps Deps) func(context.Context, *DeleteSubscribeGroupInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteSubscribeGroupInput) (*struct{}, error) {
		l := NewDeleteSubscribeGroupLogic(ctx, deps)
		if err := l.DeleteSubscribeGroup(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type DeleteSubscribeGroupLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Delete subscribe group
func NewDeleteSubscribeGroupLogic(ctx context.Context, deps Deps) *DeleteSubscribeGroupLogic {
	return &DeleteSubscribeGroupLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *DeleteSubscribeGroupLogic) DeleteSubscribeGroup(req *types.DeleteSubscribeGroupRequest) error {
	err := l.deps.DB.Model(&subscribe.Group{}).Where("id = ?", req.Id).Delete(&subscribe.Group{}).Error
	if err != nil {
		l.Error("[DeleteSubscribeGroupLogic] delete subscribe group failed: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete subscribe group failed: %v", err.Error())
	}
	return nil
}
