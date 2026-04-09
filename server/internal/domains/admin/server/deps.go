package server

import (
	modelnode "github.com/perfect-panel/server/internal/platform/persistence/node"
	modeluser "github.com/perfect-panel/server/internal/platform/persistence/user"
	"gorm.io/gorm"
)

// Deps holds the narrow admin server dependencies while Phase 6 removes
// direct ServiceContext usage from service packages.
type Deps struct {
	NodeModel modelnode.Model
	UserModel modeluser.Model
	DB        *gorm.DB
}
