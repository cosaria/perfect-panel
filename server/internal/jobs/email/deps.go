package emailLogic

import (
	"github.com/perfect-panel/server/config"
	modellog "github.com/perfect-panel/server/models/log"
	"gorm.io/gorm"
)

type Deps struct {
	DB       *gorm.DB
	LogModel modellog.Model
	Config   *config.Config
}

func (d Deps) currentConfig() config.Config {
	if d.Config == nil {
		return config.Config{}
	}
	return *d.Config
}
