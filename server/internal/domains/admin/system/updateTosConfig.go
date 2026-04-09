package system

import (
	"context"
	"reflect"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UpdateTosConfigInput struct {
	Body types.TosConfig
}

func UpdateTosConfigHandler(deps Deps) func(context.Context, *UpdateTosConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateTosConfigInput) (*struct{}, error) {
		l := NewUpdateTosConfigLogic(ctx, deps)
		if err := l.UpdateTosConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateTosConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

func NewUpdateTosConfigLogic(ctx context.Context, deps Deps) *UpdateTosConfigLogic {
	return &UpdateTosConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateTosConfigLogic) UpdateTosConfig(req *types.TosConfig) error {
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
			err = l.deps.UpdateSystemConfigField(l.ctx, db, "tos", fieldName, fieldValue)
			if err != nil {
				break
			}
		}
		if err != nil {
			return err
		}
		return l.deps.DeleteConfigCache(l.ctx, config.TosConfigKey)
	})
	if err != nil {
		l.Errorw("[UpdateTosConfigLogic] update tos config error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update tos config error: %v", err)
	}

	return nil
}
