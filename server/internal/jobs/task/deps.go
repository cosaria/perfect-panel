package task

import (
	serverconfig "github.com/perfect-panel/server/config"
	modelsubscribe "github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	modelsystem "github.com/perfect-panel/server/internal/platform/persistence/system"
	modeluser "github.com/perfect-panel/server/internal/platform/persistence/user"
	"gorm.io/gorm"
)

type Deps struct {
	DB                  *gorm.DB
	SystemModel         modelsystem.Model
	SubscribeModel      modelsubscribe.Model
	UserModel           modeluser.Model
	SetExchangeRate     func(float64)
	PrepareExchangeRate func(string, string) uint64
	StoreExchangeRate   func(uint64, string, string, float64) bool
	Config              *serverconfig.Config
}
