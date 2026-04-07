package user

import (
	"github.com/perfect-panel/server/config"
	modellog "github.com/perfect-panel/server/models/log"
	modelsubscribe "github.com/perfect-panel/server/models/subscribe"
	modeltraffic "github.com/perfect-panel/server/models/traffic"
	modeluser "github.com/perfect-panel/server/models/user"
	verifydevice "github.com/perfect-panel/server/modules/verify/device"
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
