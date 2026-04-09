package system

import (
	"context"
	"reflect"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UpdateSubscribeConfigInput struct {
	Body types.SubscribeConfig
}

func UpdateSubscribeConfigHandler(deps Deps) func(context.Context, *UpdateSubscribeConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateSubscribeConfigInput) (*struct{}, error) {
		l := NewUpdateSubscribeConfigLogic(ctx, deps)
		if err := l.UpdateSubscribeConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateSubscribeConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

func NewUpdateSubscribeConfigLogic(ctx context.Context, deps Deps) *UpdateSubscribeConfigLogic {
	return &UpdateSubscribeConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateSubscribeConfigLogic) UpdateSubscribeConfig(req *types.SubscribeConfig) error {
	v := reflect.ValueOf(*req)
	// Get the reflection type of the structure
	t := v.Type()
	err := l.deps.SystemModel.Transaction(l.ctx, func(db *gorm.DB) error {
		var err error
		for i := 0; i < v.NumField(); i++ {
			// Get the field name
			fieldName := t.Field(i).Name
			// Get the field value to string
			fieldValue := tool.ConvertValueToString(v.Field(i))
			err = l.deps.UpdateSystemConfigField(l.ctx, db, "subscribe", fieldName, fieldValue)
			if err != nil {
				break
			}
		}
		if err != nil {
			return err
		}
		return l.deps.DeleteConfigCache(l.ctx, config.SubscribeConfigKey, config.GlobalConfigKey)
	})

	if err != nil {
		l.Errorw("[UpdateSubscribeConfigLogic] update subscribe config error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update subscribe config error: %v", err)
	}

	if l.deps.currentConfig().Subscribe.SubscribePath != req.SubscribePath {
		go func() {
			if l.deps.Restart == nil {
				return
			}
			err = l.deps.Restart()
			if err != nil {
				l.Errorw("[UpdateSubscribeConfigLogic] restart error: ", logger.Field("error", err.Error()))
			}
		}()
		return nil
	}

	if err := l.deps.ReloadSubscribeConfig(); err != nil {
		l.Errorw("[UpdateSubscribeConfigLogic] reload subscribe config error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "reload subscribe config error: %v", err)
	}
	return nil
}
