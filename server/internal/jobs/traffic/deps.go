package traffic

import (
	"context"

	"github.com/hibiken/asynq"
	serverconfig "github.com/perfect-panel/server/config"
	modellog "github.com/perfect-panel/server/internal/platform/persistence/log"
	modelnode "github.com/perfect-panel/server/internal/platform/persistence/node"
	modelsubscribe "github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	modeltraffic "github.com/perfect-panel/server/internal/platform/persistence/traffic"
	modeluser "github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Deps struct {
	DB                        *gorm.DB
	Redis                     *redis.Client
	Queue                     *asynq.Client
	NodeModel                 modelnode.Model
	UserModel                 modeluser.Model
	SubscribeModel            modelsubscribe.Model
	TrafficLogModel           modeltraffic.Model
	NodeMultiplierManager     func() *modelnode.Manager
	LoadNodeMultiplierManager func(context.Context) (*modelnode.Manager, error)
	Config                    *serverconfig.Config
	LogModel                  modellog.Model
}

func (d Deps) currentConfig() serverconfig.Config {
	if d.Config == nil {
		return serverconfig.Config{}
	}
	return *d.Config
}

func (d Deps) CurrentNodeMultiplierManager() *modelnode.Manager {
	if d.NodeMultiplierManager == nil {
		return nil
	}
	return d.NodeMultiplierManager()
}

func (d Deps) ResolveNodeMultiplierManager(ctx context.Context) (*modelnode.Manager, error) {
	if manager := d.CurrentNodeMultiplierManager(); manager != nil {
		return manager, nil
	}
	if d.LoadNodeMultiplierManager == nil {
		return nil, nil
	}
	return d.LoadNodeMultiplierManager(ctx)
}
