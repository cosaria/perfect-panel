package notify

import (
	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/config"
	modelorder "github.com/perfect-panel/server/internal/platform/persistence/order"
	"gorm.io/gorm"
)

// Deps holds the narrow notify dependencies while Phase 6 removes direct
// ServiceContext usage from service packages.
type Deps struct {
	DB         *gorm.DB
	OrderModel modelorder.Model
	Queue      *asynq.Client
	Config     *config.Config
}

func (d Deps) debugEnabled() bool {
	return d.Config != nil && d.Config.Debug
}
