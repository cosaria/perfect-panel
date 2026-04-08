package subscribe

import (
	"context"
	"github.com/perfect-panel/server/models/subscribe"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type CreateSubscribeGroupInput struct {
	Body types.CreateSubscribeGroupRequest
}

func CreateSubscribeGroupHandler(deps Deps) func(context.Context, *CreateSubscribeGroupInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateSubscribeGroupInput) (*struct{}, error) {
		l := NewCreateSubscribeGroupLogic(ctx, deps)
		if err := l.CreateSubscribeGroup(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type CreateSubscribeGroupLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Create subscribe group
func NewCreateSubscribeGroupLogic(ctx context.Context, deps Deps) *CreateSubscribeGroupLogic {
	return &CreateSubscribeGroupLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *CreateSubscribeGroupLogic) CreateSubscribeGroup(req *types.CreateSubscribeGroupRequest) error {
	err := l.deps.DB.Model(&subscribe.Group{}).Create(&subscribe.Group{
		Name:        req.Name,
		Description: req.Description,
	}).Error
	if err != nil {
		l.Error("[CreateSubscribeGroupLogic] create subscribe group failed: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "create subscribe group failed: %v", err.Error())
	}
	return nil
}
