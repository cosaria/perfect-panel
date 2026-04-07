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

type UpdateVerifyConfigInput struct {
	Body types.VerifyConfig
}

func UpdateVerifyConfigHandler(deps Deps) func(context.Context, *UpdateVerifyConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateVerifyConfigInput) (*struct{}, error) {
		l := NewUpdateVerifyConfigLogic(ctx, deps)
		if err := l.UpdateVerifyConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateVerifyConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

func NewUpdateVerifyConfigLogic(ctx context.Context, deps Deps) *UpdateVerifyConfigLogic {
	return &UpdateVerifyConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateVerifyConfigLogic) UpdateVerifyConfig(req *types.VerifyConfig) error {
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
			err = db.Model(&modelsystem.System{}).Where("`category` = 'verify' and `key` = ?", fieldName).Update("value", fieldValue).Error
			if err != nil {
				break
			}
		}
		if err != nil {
			return err
		}
		// clear cache
		return l.deps.Redis.Del(l.ctx, config.VerifyConfigKey, config.GlobalConfigKey).Err()
	})
	if err != nil {
		l.Errorw("[UpdateVerifyConfigLogic] update verify config error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update verify config error: %v", err)
	}
	// Update the config
	if l.deps.Config != nil {
		tool.DeepCopy(&l.deps.Config.Verify, req)
	}
	if l.deps.ReloadVerify != nil {
		l.deps.ReloadVerify()
	}
	return nil
}
