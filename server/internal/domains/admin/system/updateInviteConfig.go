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
			err = l.deps.UpdateSystemConfigField(l.ctx, db, "invite", fieldName, fieldValue)
			if err != nil {
				break
			}
		}
		if err != nil {
			return err
		}
		// clear cache
		return l.deps.DeleteConfigCache(l.ctx, config.InviteConfigKey, config.GlobalConfigKey)
	})
	if err != nil {
		l.Errorw("[UpdateInviteConfig] update invite config error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update invite config error: %v", err)
	}
	if err := l.deps.ReloadInviteConfig(); err != nil {
		l.Errorw("[UpdateInviteConfig] reload invite config error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.ERROR), "reload invite config error: %v", err)
	}
	return nil
}
