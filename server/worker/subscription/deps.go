package subscription

import (
	"github.com/hibiken/asynq"
	serverconfig "github.com/perfect-panel/server/config"
	modelsubscribe "github.com/perfect-panel/server/models/subscribe"
	modeluser "github.com/perfect-panel/server/models/user"
)

type Deps struct {
	UserModel      modeluser.Model
	SubscribeModel modelsubscribe.Model
	Queue          *asynq.Client
	Config         *serverconfig.Config
}

func (d Deps) currentConfig() serverconfig.Config {
	if d.Config == nil {
		return serverconfig.Config{}
	}
	return *d.Config
}
