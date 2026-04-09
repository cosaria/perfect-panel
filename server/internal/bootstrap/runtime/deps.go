package runtime

import (
	"net"

	"github.com/hibiken/asynq"
	"github.com/oschwald/geoip2-golang"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/persistence/announcement"
	"github.com/perfect-panel/server/internal/platform/persistence/auth"
	"github.com/perfect-panel/server/internal/platform/persistence/client"
	"github.com/perfect-panel/server/internal/platform/persistence/coupon"
	"github.com/perfect-panel/server/internal/platform/persistence/document"
	"github.com/perfect-panel/server/internal/platform/persistence/log"
	"github.com/perfect-panel/server/internal/platform/persistence/node"
	"github.com/perfect-panel/server/internal/platform/persistence/order"
	"github.com/perfect-panel/server/internal/platform/persistence/payment"
	"github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	"github.com/perfect-panel/server/internal/platform/persistence/system"
	"github.com/perfect-panel/server/internal/platform/persistence/ticket"
	"github.com/perfect-panel/server/internal/platform/persistence/traffic"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/perfect-panel/server/internal/platform/support/limit"
	"github.com/perfect-panel/server/internal/platform/support/verify/device"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type GeoIPCityReader interface {
	City(ipAddress net.IP) (*geoip2.City, error)
}

type Deps struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Config *config.Config
	Queue  *asynq.Client
	Live   *LiveState

	AuthModel         auth.Model
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

	AuthLimiter   *limit.PeriodLimit
	DeviceManager *device.DeviceManager
	GeoIPDB       GeoIPCityReader
}
