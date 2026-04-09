package authMethod

import (
	"github.com/perfect-panel/server/config"
	modelauth "github.com/perfect-panel/server/internal/platform/persistence/auth"
)

type Deps struct {
	AuthModel    modelauth.Model
	Config       *config.Config
	ReloadEmail  func()
	ReloadMobile func()
	ReloadDevice func()
}

func (d Deps) currentConfig() config.Config {
	if d.Config == nil {
		return config.Config{}
	}
	return *d.Config
}
