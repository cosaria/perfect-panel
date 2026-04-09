package subscribe

import (
	"github.com/perfect-panel/server/config"
	modelnode "github.com/perfect-panel/server/internal/platform/persistence/node"
	modelsubscribe "github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	modeluser "github.com/perfect-panel/server/internal/platform/persistence/user"
	"gorm.io/gorm"
)

// Deps holds the narrow user/subscribe dependencies while Phase 6 removes
// direct ServiceContext usage from service packages.
type Deps struct {
	SubscribeModel modelsubscribe.Model
	UserModel      modeluser.Model
	NodeModel      modelnode.Model
	DB             *gorm.DB
	Config         *config.Config
}

func (d Deps) currentConfig() config.Config {
	if d.Config == nil {
		return config.Config{}
	}
	return *d.Config
}
