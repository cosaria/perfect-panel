package initialize

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/perfect-panel/server/config"
	modelauth "github.com/perfect-panel/server/models/auth"
	modelnode "github.com/perfect-panel/server/models/node"
	modelsystem "github.com/perfect-panel/server/models/system"
	modeluser "github.com/perfect-panel/server/models/user"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Deps struct {
	DB                       *gorm.DB
	Redis                    *redis.Client
	Config                   *config.Config
	AuthModel                modelauth.Model
	SystemModel              modelsystem.Model
	UserModel                modeluser.Model
	SetExchangeRate          func(float64)
	SetNodeMultiplierManager func(*modelnode.Manager)
	SetTelegramBot           func(*tgbotapi.BotAPI)
}

func (d Deps) currentConfig() config.Config {
	if d.Config == nil {
		return config.Config{}
	}
	return *d.Config
}
