package system

import (
	"context"
	"errors"

	"github.com/perfect-panel/server/config"
	modelnode "github.com/perfect-panel/server/models/node"
	modelsystem "github.com/perfect-panel/server/models/system"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Deps struct {
	SystemModel           modelsystem.Model
	Redis                 *redis.Client
	Config                *config.Config
	UpdateSiteField       func(context.Context, *gorm.DB, string, string) error
	UpdateConfigField     func(context.Context, *gorm.DB, string, string, string) error
	DeleteCacheKeys       func(context.Context, ...string) error
	NodeMultiplierManager func() *modelnode.Manager
	Restart               func() error
	ReloadVerify          func()
	RunReloadVerify       func() error
	ReloadNode            func()
	RunReloadNode         func() error
	ReloadCurrency        func()
	RunReloadCurrency     func() error
	ReloadInvite          func()
	RunReloadInvite       func() error
	ReloadRegister        func()
	RunReloadRegister     func() error
	ReloadSite            func()
	RunReloadSite         func() error
	ReloadSubscribe       func()
	RunReloadSubscribe    func() error
	ReloadTelegram        func()
}

func (d Deps) currentConfig() config.Config {
	if d.Config == nil {
		return config.Config{}
	}
	return *d.Config
}

func (d Deps) CurrentNodeMultiplierManager() *modelnode.Manager {
	if d.NodeMultiplierManager == nil {
		return nil
	}
	return d.NodeMultiplierManager()
}

func (d Deps) UpdateSiteConfigField(ctx context.Context, db *gorm.DB, fieldName, fieldValue string) error {
	return d.UpdateSystemConfigField(ctx, db, "site", fieldName, fieldValue)
}

func (d Deps) UpdateSystemConfigField(ctx context.Context, db *gorm.DB, category, fieldName, fieldValue string) error {
	if d.UpdateConfigField != nil {
		return d.UpdateConfigField(ctx, db, category, fieldName, fieldValue)
	}
	if d.UpdateSiteField != nil {
		if category == "site" {
			return d.UpdateSiteField(ctx, db, fieldName, fieldValue)
		}
	}
	if db == nil {
		return errors.New("site config transaction db is nil")
	}
	return db.Model(&modelsystem.System{}).
		Where("`category` = ? and `key` = ?", category, fieldName).
		Update("value", fieldValue).Error
}

func (d Deps) DeleteConfigCache(ctx context.Context, keys ...string) error {
	if d.DeleteCacheKeys != nil {
		return d.DeleteCacheKeys(ctx, keys...)
	}
	if d.Redis == nil {
		return errors.New("redis client is nil")
	}
	return d.Redis.Del(ctx, keys...).Err()
}

func (d Deps) ReloadSiteConfig() error {
	if d.RunReloadSite != nil {
		return d.RunReloadSite()
	}
	if d.ReloadSite != nil {
		d.ReloadSite()
	}
	return nil
}

func (d Deps) ReloadCurrencyConfig() error {
	if d.RunReloadCurrency != nil {
		return d.RunReloadCurrency()
	}
	if d.ReloadCurrency != nil {
		d.ReloadCurrency()
	}
	return nil
}

func (d Deps) ReloadInviteConfig() error {
	if d.RunReloadInvite != nil {
		return d.RunReloadInvite()
	}
	if d.ReloadInvite != nil {
		d.ReloadInvite()
	}
	return nil
}

func (d Deps) ReloadRegisterConfig() error {
	if d.RunReloadRegister != nil {
		return d.RunReloadRegister()
	}
	if d.ReloadRegister != nil {
		d.ReloadRegister()
	}
	return nil
}

func (d Deps) ReloadVerifyConfig() error {
	if d.RunReloadVerify != nil {
		return d.RunReloadVerify()
	}
	if d.ReloadVerify != nil {
		d.ReloadVerify()
	}
	return nil
}

func (d Deps) ReloadSubscribeConfig() error {
	if d.RunReloadSubscribe != nil {
		return d.RunReloadSubscribe()
	}
	if d.ReloadSubscribe != nil {
		d.ReloadSubscribe()
	}
	return nil
}

func (d Deps) ReloadNodeConfig() error {
	if d.RunReloadNode != nil {
		return d.RunReloadNode()
	}
	if d.ReloadNode != nil {
		d.ReloadNode()
	}
	return nil
}
