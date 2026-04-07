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

type UpdateVerifyCodeConfigInput struct {
	Body types.VerifyCodeConfig
}

func UpdateVerifyCodeConfigHandler(deps Deps) func(context.Context, *UpdateVerifyCodeConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateVerifyCodeConfigInput) (*struct{}, error) {
		l := NewUpdateVerifyCodeConfigLogic(ctx, deps)
		if err := l.UpdateVerifyCodeConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateVerifyCodeConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Update Verify Code Config
func NewUpdateVerifyCodeConfigLogic(ctx context.Context, deps Deps) *UpdateVerifyCodeConfigLogic {
	return &UpdateVerifyCodeConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateVerifyCodeConfigLogic) UpdateVerifyCodeConfig(req *types.VerifyCodeConfig) error {
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
			err = db.Model(&modelsystem.System{}).Where("`category` = 'verify_code' and `key` = ?", fieldName).Update("value", fieldValue).Error
			if err != nil {
				break
			}
		}
		err = l.deps.Redis.Del(l.ctx, config.VerifyCodeConfigKey).Err()
		return err
	})
	if err != nil {
		l.Errorw("[UpdateRegisterConfig] update verify code config error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update register config error: %v", err.Error())
	}
	return nil
}
