package traffic

import (
	"github.com/hibiken/asynq"
	serverconfig "github.com/perfect-panel/server/config"
	modellog "github.com/perfect-panel/server/models/log"
	modelnode "github.com/perfect-panel/server/models/node"
	modelsubscribe "github.com/perfect-panel/server/models/subscribe"
	modeltraffic "github.com/perfect-panel/server/models/traffic"
	modeluser "github.com/perfect-panel/server/models/user"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Deps struct {
	DB                    *gorm.DB
	Redis                 *redis.Client
	Queue                 *asynq.Client
	NodeModel             modelnode.Model
	UserModel             modeluser.Model
	SubscribeModel        modelsubscribe.Model
	TrafficLogModel       modeltraffic.Model
	NodeMultiplierManager *modelnode.Manager
	Config                *serverconfig.Config
	LogModel              modellog.Model
}

func (d Deps) currentConfig() serverconfig.Config {
	if d.Config == nil {
		return serverconfig.Config{}
	}
	return *d.Config
}
