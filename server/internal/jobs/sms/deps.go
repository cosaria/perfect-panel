package smslogic

import (
	"github.com/perfect-panel/server/config"
	modellog "github.com/perfect-panel/server/internal/platform/persistence/log"
)

type Deps struct {
	LogModel modellog.Model
	Config   *config.Config
}

func (d Deps) currentConfig() config.Config {
	if d.Config == nil {
		return config.Config{}
	}
	return *d.Config
}
