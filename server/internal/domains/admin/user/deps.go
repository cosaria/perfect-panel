package user

import (
	"github.com/perfect-panel/server/config"
	modellog "github.com/perfect-panel/server/internal/platform/persistence/log"
	modelsubscribe "github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	modeltraffic "github.com/perfect-panel/server/internal/platform/persistence/traffic"
	modeluser "github.com/perfect-panel/server/internal/platform/persistence/user"
	verifydevice "github.com/perfect-panel/server/internal/platform/support/verify/device"
)

// Deps holds the narrow admin user dependencies while Phase 6 removes direct
// ServiceContext usage from service packages.
type Deps struct {
	UserModel       modeluser.Model
	SubscribeModel  modelsubscribe.Model
	LogModel        modellog.Model
	TrafficLogModel modeltraffic.Model
	DeviceManager   *verifydevice.DeviceManager
	Config          *config.Config
}
