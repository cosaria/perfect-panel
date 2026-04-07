package portal

import (
	"github.com/hibiken/asynq"
	serverconfig "github.com/perfect-panel/server/config"
	modelcoupon "github.com/perfect-panel/server/models/coupon"
	modelorder "github.com/perfect-panel/server/models/order"
	modelpayment "github.com/perfect-panel/server/models/payment"
	modelsubscribe "github.com/perfect-panel/server/models/subscribe"
	modeluser "github.com/perfect-panel/server/models/user"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Deps struct {
	PaymentModel    modelpayment.Model
	SubscribeModel  modelsubscribe.Model
	CouponModel     modelcoupon.Model
	OrderModel      modelorder.Model
	UserModel       modeluser.Model
	DB              *gorm.DB
	Redis           *redis.Client
	Queue           *asynq.Client
	Config          *serverconfig.Config
	GetExchangeRate func() float64
	SetExchangeRate func(float64)
	GetExchangeRateSnapshot func() ExchangeRateSnapshot
	PrepareExchangeRate     func(string, string) uint64
	StoreExchangeRate       func(uint64, string, string, float64) bool
}

type ExchangeRateSnapshot struct {
	Version uint64
	From    string
	To      string
	Rate    float64
}

func (d Deps) CurrentExchangeRate() float64 {
	if d.GetExchangeRateSnapshot != nil {
		return d.GetExchangeRateSnapshot().Rate
	}
	if d.GetExchangeRate == nil {
		return 0
	}
	return d.GetExchangeRate()
}

func (d Deps) CurrentExchangeRateSnapshot() ExchangeRateSnapshot {
	if d.GetExchangeRateSnapshot != nil {
		return d.GetExchangeRateSnapshot()
	}
	return ExchangeRateSnapshot{Rate: d.CurrentExchangeRate()}
}

func (d Deps) CacheExchangeRate(rate float64) {
	if d.SetExchangeRate != nil {
		d.SetExchangeRate(rate)
	}
}

func (d Deps) PrepareExchangeRateCache(from, to string) uint64 {
	if d.PrepareExchangeRate == nil {
		return 0
	}
	return d.PrepareExchangeRate(from, to)
}

func (d Deps) StoreExchangeRateCache(version uint64, from, to string, rate float64) bool {
	if d.StoreExchangeRate != nil {
		return d.StoreExchangeRate(version, from, to, rate)
	}
	d.CacheExchangeRate(rate)
	return true
}
