package runtime

import (
	"net"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/oschwald/geoip2-golang"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/ads"
	"github.com/perfect-panel/server/models/announcement"
	"github.com/perfect-panel/server/models/auth"
	"github.com/perfect-panel/server/models/client"
	"github.com/perfect-panel/server/models/coupon"
	"github.com/perfect-panel/server/models/document"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/models/order"
	"github.com/perfect-panel/server/models/payment"
	"github.com/perfect-panel/server/models/subscribe"
	"github.com/perfect-panel/server/models/system"
	"github.com/perfect-panel/server/models/ticket"
	"github.com/perfect-panel/server/models/traffic"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/limit"
	"github.com/perfect-panel/server/modules/verify/device"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GeoIPCityReader interface {
	City(ipAddress net.IP) (*geoip2.City, error)
}

type Deps struct {
	DB           *gorm.DB
	Redis        *redis.Client
	Config       *config.Config
	Queue        *asynq.Client
	ExchangeRate float64

	AuthModel         auth.Model
	AdsModel          ads.Model
	LogModel          log.Model
	NodeModel         node.Model
	UserModel         user.Model
	OrderModel        order.Model
	ClientModel       client.Model
	TicketModel       ticket.Model
	SystemModel       system.Model
	CouponModel       coupon.Model
	PaymentModel      payment.Model
	DocumentModel     document.Model
	SubscribeModel    subscribe.Model
	TrafficLogModel   traffic.Model
	AnnouncementModel announcement.Model

	Restart               func() error
	TelegramBot           *tgbotapi.BotAPI
	NodeMultiplierManager *node.Manager
	AuthLimiter           *limit.PeriodLimit
	DeviceManager         *device.DeviceManager
	GeoIPDB               GeoIPCityReader
}
