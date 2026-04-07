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
	PaymentModel   modelpayment.Model
	SubscribeModel modelsubscribe.Model
	CouponModel    modelcoupon.Model
	OrderModel     modelorder.Model
	UserModel      modeluser.Model
	DB             *gorm.DB
	Redis          *redis.Client
	Queue          *asynq.Client
	Config         *serverconfig.Config
	ExchangeRate   float64
}
