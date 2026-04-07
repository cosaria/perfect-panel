package system

import (
	"context"
	"reflect"

	"github.com/perfect-panel/server/config"
	modelsystem "github.com/perfect-panel/server/models/system"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UpdateRegisterConfigInput struct {
	Body types.RegisterConfig
}

func UpdateRegisterConfigHandler(deps Deps) func(context.Context, *UpdateRegisterConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateRegisterConfigInput) (*struct{}, error) {
		l := NewUpdateRegisterConfigLogic(ctx, deps)
		if err := l.UpdateRegisterConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateRegisterConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

func NewUpdateRegisterConfigLogic(ctx context.Context, deps Deps) *UpdateRegisterConfigLogic {
	return &UpdateRegisterConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateRegisterConfigLogic) UpdateRegisterConfig(req *types.RegisterConfig) error {
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
			// Update the site config
			err = db.Model(&modelsystem.System{}).Where("`category` = 'register' and `key` = ?", fieldName).Update("value", fieldValue).Error
			if err != nil {
				break
			}
		}
		if err != nil {
			return err
		}
		return l.deps.Redis.Del(l.ctx, config.RegisterConfigKey, config.GlobalConfigKey).Err()
	})
	if err != nil {
		l.Errorw("[UpdateRegisterConfig] update register config error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update register config error: %v", err.Error())
	}
	// init system config
	if l.deps.ReloadRegister != nil {
		l.deps.ReloadRegister()
	}
	return nil
}
