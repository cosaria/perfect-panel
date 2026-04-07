package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/perfect-panel/server/config"
	modelauth "github.com/perfect-panel/server/models/auth"
	modelsystem "github.com/perfect-panel/server/models/system"
	modeluser "github.com/perfect-panel/server/models/user"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Deps holds the narrow telegram dependencies while Phase 6 removes direct
// ServiceContext usage from service packages.
type Deps struct {
	AuthModel   modelauth.Model
	SystemModel modelsystem.Model
	UserModel   modeluser.Model
	Redis       *redis.Client
	DB          *gorm.DB
	TelegramBot *tgbotapi.BotAPI
	Config      *config.Config
}

func (d Deps) currentConfig() config.Config {
	if d.Config == nil {
		return config.Config{}
	}
	return *d.Config
}
