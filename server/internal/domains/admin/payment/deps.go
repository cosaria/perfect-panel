package payment

import (
	serverconfig "github.com/perfect-panel/server/config"
	modelpayment "github.com/perfect-panel/server/internal/platform/persistence/payment"
)

// Deps holds the narrow admin payment dependencies while Phase 6 removes
// direct ServiceContext usage from service packages.
type Deps struct {
	PaymentModel modelpayment.Model
	Config       *serverconfig.Config
}
