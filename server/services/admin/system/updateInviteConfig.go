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

type UpdateInviteConfigInput struct {
	Body types.InviteConfig
}

func UpdateInviteConfigHandler(deps Deps) func(context.Context, *UpdateInviteConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateInviteConfigInput) (*struct{}, error) {
		l := NewUpdateInviteConfigLogic(ctx, deps)
		if err := l.UpdateInviteConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateInviteConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

func NewUpdateInviteConfigLogic(ctx context.Context, deps Deps) *UpdateInviteConfigLogic {
	return &UpdateInviteConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateInviteConfigLogic) UpdateInviteConfig(req *types.InviteConfig) error {
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
			// Update the invite config
			err = db.Model(&modelsystem.System{}).Where("`category` = 'invite' and `key` = ?", fieldName).Update("value", fieldValue).Error
			if err != nil {
				break
			}
		}
		if err != nil {
			return err
		}
		// clear cache
		return l.deps.Redis.Del(l.ctx, config.InviteConfigKey, config.GlobalConfigKey).Err()
	})
	if err != nil {
		l.Errorw("[UpdateInviteConfig] update invite config error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update invite config error: %v", err)
	}
	if l.deps.ReloadInvite != nil {
		l.deps.ReloadInvite()
	}
	return nil
}
