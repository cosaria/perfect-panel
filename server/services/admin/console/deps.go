package console

import (
	modelnode "github.com/perfect-panel/server/models/node"
	modelorder "github.com/perfect-panel/server/models/order"
	modelticket "github.com/perfect-panel/server/models/ticket"
	modeluser "github.com/perfect-panel/server/models/user"
	"gorm.io/gorm"
)

// Deps holds the narrow admin console dependencies while Phase 6 removes
// direct ServiceContext usage from service packages.
type Deps struct {
	OrderModel  modelorder.Model
	UserModel   modeluser.Model
	NodeModel   modelnode.Model
	TicketModel modelticket.Model
	DB          *gorm.DB
}
