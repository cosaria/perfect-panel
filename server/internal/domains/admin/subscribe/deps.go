package subscribe

import (
	modelsubscribe "github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	modeluser "github.com/perfect-panel/server/internal/platform/persistence/user"
	verifydevice "github.com/perfect-panel/server/internal/platform/support/verify/device"
	"gorm.io/gorm"
)

// Deps holds the narrow admin subscribe dependencies while Phase 6 removes
// direct ServiceContext usage from service packages.
type Deps struct {
	SubscribeModel modelsubscribe.Model
	UserModel      modeluser.Model
	DB             *gorm.DB
	DeviceManager  *verifydevice.DeviceManager
}
