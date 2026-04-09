package ticket

import modelticket "github.com/perfect-panel/server/internal/platform/persistence/ticket"

// Deps holds the narrow user/ticket dependencies while Phase 6 removes direct
// ServiceContext usage from service packages.
type Deps struct {
	TicketModel modelticket.Model
}
