package subscribe

import (
	"github.com/perfect-panel/server/config"
	modelclient "github.com/perfect-panel/server/models/client"
	modellog "github.com/perfect-panel/server/models/log"
	modelnode "github.com/perfect-panel/server/models/node"
	modelsubscribe "github.com/perfect-panel/server/models/subscribe"
	modeluser "github.com/perfect-panel/server/models/user"
)

// Deps holds the narrow subscribe dependencies while Phase 6 removes direct
// ServiceContext usage from service packages.
type Deps struct {
	ClientModel    modelclient.Model
	LogModel       modellog.Model
	NodeModel      modelnode.Model
	SubscribeModel modelsubscribe.Model
	UserModel      modeluser.Model
	Config         *config.Config
}

func (d Deps) currentConfig() config.Config {
	if d.Config == nil {
		return config.Config{}
	}
	return *d.Config
}
