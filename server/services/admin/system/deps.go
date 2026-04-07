package system

import (
	"github.com/perfect-panel/server/config"
	modelnode "github.com/perfect-panel/server/models/node"
	modelsystem "github.com/perfect-panel/server/models/system"
	"github.com/redis/go-redis/v9"
)

type Deps struct {
	SystemModel           modelsystem.Model
	Redis                 *redis.Client
	Config                *config.Config
	NodeMultiplierManager func() *modelnode.Manager
	Restart               func() error
	ReloadVerify          func()
	ReloadNode            func()
	ReloadCurrency        func()
	ReloadInvite          func()
	ReloadRegister        func()
	ReloadSite            func()
	ReloadSubscribe       func()
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
