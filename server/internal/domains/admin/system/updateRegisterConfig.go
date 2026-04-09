package system

import (
	"context"
	"reflect"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
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
			err = l.deps.UpdateSystemConfigField(l.ctx, db, "register", fieldName, fieldValue)
			if err != nil {
				break
			}
		}
		if err != nil {
			return err
		}
		return l.deps.DeleteConfigCache(l.ctx, config.RegisterConfigKey, config.GlobalConfigKey)
	})
	if err != nil {
		l.Errorw("[UpdateRegisterConfig] update register config error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update register config error: %v", err.Error())
	}
	if err := l.deps.ReloadRegisterConfig(); err != nil {
		l.Errorw("[UpdateRegisterConfig] reload register config error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "reload register config error: %v", err.Error())
	}
	return nil
}
