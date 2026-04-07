package subscribe

import (
	"context"
	"github.com/perfect-panel/server/models/subscribe"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type UpdateSubscribeGroupInput struct {
	Body types.UpdateSubscribeGroupRequest
}

func UpdateSubscribeGroupHandler(deps Deps) func(context.Context, *UpdateSubscribeGroupInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateSubscribeGroupInput) (*struct{}, error) {
		l := NewUpdateSubscribeGroupLogic(ctx, deps)
		if err := l.UpdateSubscribeGroup(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateSubscribeGroupLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Update subscribe group
func NewUpdateSubscribeGroupLogic(ctx context.Context, deps Deps) *UpdateSubscribeGroupLogic {
	return &UpdateSubscribeGroupLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateSubscribeGroupLogic) UpdateSubscribeGroup(req *types.UpdateSubscribeGroupRequest) error {
	err := l.deps.DB.Model(&subscribe.Group{}).Where("id = ?", req.Id).Save(&subscribe.Group{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
	}).Error
	if err != nil {
		l.Logger.Error("[UpdateSubscribeGroup] update subscribe group failed", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update subscribe group failed: %v", err.Error())
	}
	return nil
}
