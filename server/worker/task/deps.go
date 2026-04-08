package task

import (
	serverconfig "github.com/perfect-panel/server/config"
	modelsubscribe "github.com/perfect-panel/server/models/subscribe"
	modelsystem "github.com/perfect-panel/server/models/system"
	modeluser "github.com/perfect-panel/server/models/user"
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
