package svc

import (
	"context"
	"errors"

	"github.com/perfect-panel/server/models/client"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/modules/verify/device"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/ads"
	"github.com/perfect-panel/server/models/announcement"
	"github.com/perfect-panel/server/models/auth"
	"github.com/perfect-panel/server/models/coupon"
	"github.com/perfect-panel/server/models/document"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/models/order"
	"github.com/perfect-panel/server/models/payment"
	"github.com/perfect-panel/server/models/subscribe"
	"github.com/perfect-panel/server/models/system"
	"github.com/perfect-panel/server/models/ticket"
	"github.com/perfect-panel/server/models/traffic"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/limit"
	"github.com/perfect-panel/server/modules/infra/orm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ServiceContext is a temporary composition-root shell for Phase 6 migration.
// Downstream packages should translate it into package-local deps instead of
// treating it as a universal service-layer container.
type ServiceContext struct {
	DB           *gorm.DB
	Redis        *redis.Client
	Config       config.Config
	Queue        *asynq.Client
	ExchangeRate float64
	GeoIP        *IPLocation

	//NodeCache   *cache.NodeCacheClient
	AuthModel   auth.Model
	AdsModel    ads.Model
	LogModel    log.Model
	NodeModel   node.Model
	UserModel   user.Model
	OrderModel  order.Model
	ClientModel client.Model
	TicketModel ticket.Model
	//ServerModel        server.Model
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
}

func NewServiceContext(c config.Config) *ServiceContext {
	// gorm initialize
	db, err := orm.ConnectMysql(orm.Mysql{
		Config: c.MySQL,
	})

	if err != nil {
		panic(err.Error())
	}

	// IP location initialize
	geoIP, err := NewIPLocation("./cache/GeoLite2-City.mmdb")
	if err != nil {
		panic(err.Error())
	}

	rds := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host,
		Password: c.Redis.Pass,
		DB:       c.Redis.DB,
	})
	err = rds.Ping(context.Background()).Err()
	if err != nil {
		panic(err.Error())
	} else {
		if err := clearSendCountLimitKeys(context.Background(), rds); err != nil {
			panic(err.Error())
		}
	}
	authLimiter := limit.NewPeriodLimit(86400, 15, rds, config.SendCountLimitKeyPrefix, limit.Align())
	srv := &ServiceContext{
		DB:           db,
		Redis:        rds,
		Config:       c,
		Queue:        NewAsynqClient(c),
		ExchangeRate: 0,
		GeoIP:        geoIP,
		//NodeCache:   cache.NewNodeCacheClient(rds),
		AuthLimiter: authLimiter,
		AdsModel:    ads.NewModel(db, rds),
		LogModel:    log.NewModel(db),
		NodeModel:   node.NewModel(db, rds),
		AuthModel:   auth.NewModel(db, rds),
		UserModel:   user.NewModel(db, rds),
		OrderModel:  order.NewModel(db, rds),
		ClientModel: client.NewSubscribeApplicationModel(db),
		TicketModel: ticket.NewModel(db, rds),
		//ServerModel:       server.NewModel(db, rds),
		SystemModel:       system.NewModel(db, rds),
		CouponModel:       coupon.NewModel(db, rds),
		PaymentModel:      payment.NewModel(db, rds),
		DocumentModel:     document.NewModel(db, rds),
		SubscribeModel:    subscribe.NewModel(db, rds),
		TrafficLogModel:   traffic.NewModel(db),
		AnnouncementModel: announcement.NewModel(db, rds),
	}
	srv.DeviceManager = NewDeviceManager(srv)
	return srv

}

func clearSendCountLimitKeys(ctx context.Context, rds *redis.Client) error {
	var cursor uint64

	for {
		keys, nextCursor, err := rds.Scan(ctx, cursor, config.SendCountLimitKeyPrefix+"*", 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := rds.Del(ctx, keys...).Err(); err != nil && !errors.Is(err, redis.Nil) {
				return err
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			return nil
		}
	}
}
