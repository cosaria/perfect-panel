package server

import (
	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/config"
	modelnode "github.com/perfect-panel/server/internal/platform/persistence/node"
	modelsubscribe "github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	modeluser "github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/redis/go-redis/v9"
)

// Deps holds the narrow node-service dependencies while Phase 6 removes direct
// ServiceContext usage from service packages.
type Deps struct {
	NodeModel      modelnode.Model
	SubscribeModel modelsubscribe.Model
	UserModel      modeluser.Model
	Redis          *redis.Client
	Queue          *asynq.Client
	Config         *config.Config
}

func (d Deps) currentConfig() config.Config {
	if d.Config == nil {
		return config.Config{}
	}
	return *d.Config
}
