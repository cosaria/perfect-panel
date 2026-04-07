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

type UpdateNodeConfigInput struct {
	Body types.NodeConfig
}

func UpdateNodeConfigHandler(deps Deps) func(context.Context, *UpdateNodeConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateNodeConfigInput) (*struct{}, error) {
		l := NewUpdateNodeConfigLogic(ctx, deps)
		if err := l.UpdateNodeConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateNodeConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

func NewUpdateNodeConfigLogic(ctx context.Context, deps Deps) *UpdateNodeConfigLogic {
	return &UpdateNodeConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateNodeConfigLogic) UpdateNodeConfig(req *types.NodeConfig) error {
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
			// Update the server config
			err = db.Model(&modelsystem.System{}).Where("`category` = 'server' and `key` = ?", fieldName).Update("value", fieldValue).Error
			if err != nil {
				break
			}
		}
		return l.deps.Redis.Del(l.ctx, config.NodeConfigKey).Err()
	})
	if err != nil {
		l.Errorw("[UpdateNodeConfig] update node config error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update server config error: %v", err)
	}
	if l.deps.ReloadNode != nil {
		l.deps.ReloadNode()
	}
	return nil
}
