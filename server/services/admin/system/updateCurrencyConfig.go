package system

import (
	"context"
	"reflect"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UpdateCurrencyConfigInput struct {
	Body types.CurrencyConfig
}

func UpdateCurrencyConfigHandler(deps Deps) func(context.Context, *UpdateCurrencyConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateCurrencyConfigInput) (*struct{}, error) {
		l := NewUpdateCurrencyConfigLogic(ctx, deps)
		if err := l.UpdateCurrencyConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateCurrencyConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Update Currency Config
func NewUpdateCurrencyConfigLogic(ctx context.Context, deps Deps) *UpdateCurrencyConfigLogic {
	return &UpdateCurrencyConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateCurrencyConfigLogic) UpdateCurrencyConfig(req *types.CurrencyConfig) error {
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
			err = l.deps.UpdateSystemConfigField(l.ctx, db, "currency", fieldName, fieldValue)
			if err != nil {
				break
			}
		}
		if err != nil {
			return err
		}
		// clear cache
		return l.deps.DeleteConfigCache(l.ctx, config.CurrencyConfigKey, config.GlobalConfigKey)
	})
	if err != nil {
		l.Errorw("[UpdateCurrencyConfig] update currency config error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update invite config error: %v", err)
	}
	if err := l.deps.ReloadCurrencyConfig(); err != nil {
		l.Errorw("[UpdateCurrencyConfig] reload currency config error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "reload currency config error: %v", err)
	}
	return nil
}
