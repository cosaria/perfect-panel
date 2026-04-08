package common

import (
	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/config"
	modelauth "github.com/perfect-panel/server/models/auth"
	modelclient "github.com/perfect-panel/server/models/client"
	modelsystem "github.com/perfect-panel/server/models/system"
	modeluser "github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/limit"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Deps struct {
	AuthModel   modelauth.Model
	ClientModel modelclient.Model
	SystemModel modelsystem.Model
	UserModel   modeluser.Model
	DB          *gorm.DB
	Redis       *redis.Client
	AuthLimiter *limit.PeriodLimit
	Queue       *asynq.Client
	Config      *config.Config
}

func (d Deps) currentConfig() config.Config {
	if d.Config == nil {
		return config.Config{}
	}
	return *d.Config
}
