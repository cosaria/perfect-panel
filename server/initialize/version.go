package initialize

import (
	"errors"
	"time"

	"github.com/perfect-panel/server/models/user"
	"gorm.io/gorm"

	"github.com/perfect-panel/server/models/migrate"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/orm"
)

func Migrate(deps Deps) {
	cfg := deps.currentConfig()
	mc := orm.Mysql{
		Config: cfg.MySQL,
	}
	now := time.Now()
	if err := migrate.Migrate(mc.Dsn()).Up(); err != nil {
		if errors.Is(err, migrate.NoChange) {
			logger.Info("[Migrate] database not change")
			return
		}
		logger.Errorf("[Migrate] Up error: %v", err.Error())
		panic(err)
	} else {
		logger.Info("[Migrate] Database change, took " + time.Since(now).String())
	}
	// if not found admin user
	err := deps.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&user.User{}).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := migrate.CreateAdminUser(cfg.Administrator.Email, cfg.Administrator.Password, tx); err != nil {
				logger.Errorf("[Migrate] CreateAdminUser error: %v", err.Error())
				return err
			}
			logger.Info("[Migrate] Create admin user success")
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
