package auth

import (
	"github.com/perfect-panel/server/config"
	modellog "github.com/perfect-panel/server/internal/platform/persistence/log"
	modelsubscribe "github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	modeluser "github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/redis/go-redis/v9"
)

// Deps holds the narrow auth dependencies while Phase 6 removes direct
// ServiceContext usage from service packages.
type Deps struct {
	UserModel      modeluser.Model
	LogModel       modellog.Model
	SubscribeModel modelsubscribe.Model
	Redis          *redis.Client
	Config         *config.Config
}

func (d Deps) currentConfig() config.Config {
	if d.Config == nil {
		return config.Config{}
	}
	return *d.Config
}
