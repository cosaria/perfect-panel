package user

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/perfect-panel/server/config"
	modelauth "github.com/perfect-panel/server/internal/platform/persistence/auth"
	modellog "github.com/perfect-panel/server/internal/platform/persistence/log"
	modelorder "github.com/perfect-panel/server/internal/platform/persistence/order"
	modelsubscribe "github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	modeluser "github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Deps holds the narrow user/user read-path dependencies while Phase 6 removes
// direct ServiceContext usage from service packages.
type Deps struct {
	UserModel      modeluser.Model
	LogModel       modellog.Model
	AuthModel      modelauth.Model
	OrderModel     modelorder.Model
	SubscribeModel modelsubscribe.Model
	Redis          *redis.Client
	Config         *config.Config
	TelegramBot    func() *tgbotapi.BotAPI
	DB             *gorm.DB
}

func (d Deps) currentConfig() config.Config {
	if d.Config == nil {
		return config.Config{}
	}
	return *d.Config
}

func (d Deps) CurrentTelegramBot() *tgbotapi.BotAPI {
	if d.TelegramBot == nil {
		return nil
	}
	return d.TelegramBot()
}
