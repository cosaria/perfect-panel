package log

import (
	"github.com/perfect-panel/server/config"
	modellog "github.com/perfect-panel/server/models/log"
	modelsystem "github.com/perfect-panel/server/models/system"
	"gorm.io/gorm"
)

// Deps holds the narrow admin log dependencies while Phase 6 removes direct
// ServiceContext usage from service packages.
type Deps struct {
	LogModel    modellog.Model
	SystemModel modelsystem.Model
	DB          *gorm.DB
	Config      *config.Config
}

func (d Deps) currentConfig() config.Config {
	if d.Config == nil {
		return config.Config{}
	}
	return *d.Config
}
