package order

import (
	"github.com/hibiken/asynq"
	serverconfig "github.com/perfect-panel/server/config"
	modelcoupon "github.com/perfect-panel/server/models/coupon"
	modelorder "github.com/perfect-panel/server/models/order"
	modelpayment "github.com/perfect-panel/server/models/payment"
	modelsubscribe "github.com/perfect-panel/server/models/subscribe"
	modeluser "github.com/perfect-panel/server/models/user"
	"gorm.io/gorm"
)

// Deps holds the narrow public order dependencies while Phase 6 removes
// direct ServiceContext usage from service packages.
type Deps struct {
	OrderModel     modelorder.Model
	PaymentModel   modelpayment.Model
	SubscribeModel modelsubscribe.Model
	UserModel      modeluser.Model
	CouponModel    modelcoupon.Model
	DB             *gorm.DB
	Queue          *asynq.Client
	Config         *serverconfig.Config
}
