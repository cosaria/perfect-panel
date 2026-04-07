package orderLogic

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	serverconfig "github.com/perfect-panel/server/config"
	modelcoupon "github.com/perfect-panel/server/models/coupon"
	modellog "github.com/perfect-panel/server/models/log"
	modelorder "github.com/perfect-panel/server/models/order"
	modelpayment "github.com/perfect-panel/server/models/payment"
	modelsubscribe "github.com/perfect-panel/server/models/subscribe"
	modeluser "github.com/perfect-panel/server/models/user"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Deps struct {
	OrderModel     modelorder.Model
	PaymentModel   modelpayment.Model
	SubscribeModel modelsubscribe.Model
	UserModel      modeluser.Model
	CouponModel    modelcoupon.Model
	LogModel       modellog.Model
	DB             *gorm.DB
	Queue          *asynq.Client
	Redis          *redis.Client
	TelegramBot    func() *tgbotapi.BotAPI
	Config         *serverconfig.Config
}

func (d Deps) currentConfig() serverconfig.Config {
	if d.Config == nil {
		return serverconfig.Config{}
	}
	return *d.Config
}

func (d Deps) CurrentTelegramBot() *tgbotapi.BotAPI {
	if d.TelegramBot == nil {
		return nil
	}
	return d.TelegramBot()
}
