package marketing

import (
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

// Deps holds the narrow admin marketing dependencies while Phase 6 removes
// direct ServiceContext usage from service packages.
type Deps struct {
	DB    *gorm.DB
	Queue *asynq.Client
}
