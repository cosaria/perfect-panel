package system

import (
	"context"
	"reflect"

	"github.com/perfect-panel/server/config"
	modelsystem "github.com/perfect-panel/server/models/system"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UpdateSiteConfigInput struct {
	Body types.SiteConfig
}

func UpdateSiteConfigHandler(deps Deps) func(context.Context, *UpdateSiteConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateSiteConfigInput) (*struct{}, error) {
		l := NewUpdateSiteConfigLogic(ctx, deps)
		if err := l.UpdateSiteConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateSiteConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

func NewUpdateSiteConfigLogic(ctx context.Context, deps Deps) *UpdateSiteConfigLogic {
	return &UpdateSiteConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateSiteConfigLogic) UpdateSiteConfig(req *types.SiteConfig) error {
	// Get the reflection value of the structure
	v := reflect.ValueOf(*req)
	// Get the reflection type of the structure
	t := v.Type()
	err := l.deps.SystemModel.Transaction(l.ctx, func(db *gorm.DB) error {
		var err error
		for i := 0; i < v.NumField(); i++ {
			// Get the field name
			fieldName := t.Field(i).Name
			// Get the field value
			fieldValue := v.Field(i)
			err = db.Model(&modelsystem.System{}).Where("`category` = 'site' and `key` = ?", fieldName).Update("value", fieldValue.String()).Error
			if err != nil {
				break
			}
		}
		if err != nil {
			return err
		}

		return l.deps.Redis.Del(l.ctx, config.SiteConfigKey, config.GlobalConfigKey).Err()
	})
	if err != nil {
		l.Logger.Error("[UpdateSiteConfig] update site config error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update site config error: %v", err.Error())
	}
	if l.deps.ReloadSite != nil {
		l.deps.ReloadSite()
	}
	return nil
}
