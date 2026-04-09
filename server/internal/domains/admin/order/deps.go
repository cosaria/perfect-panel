package order

import (
	"github.com/hibiken/asynq"
	modelorder "github.com/perfect-panel/server/models/order"
	modelpayment "github.com/perfect-panel/server/models/payment"
)

// Deps holds the narrow admin order dependencies while Phase 6 removes
// direct ServiceContext usage from service packages.
type Deps struct {
	OrderModel   modelorder.Model
	PaymentModel modelpayment.Model
	Queue        *asynq.Client
}
