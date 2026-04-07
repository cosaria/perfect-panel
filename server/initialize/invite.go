package initialize

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/util/tool"
)

func Invite(deps Deps) {
	// Initialize the system configuration
	logger.Debug("Register config initialization")
	configs, err := deps.SystemModel.GetInviteConfig(context.Background())
	if err != nil {
		logger.Error("[Init Invite Config] Get Invite Config Error: ", logger.Field("error", err.Error()))
		return
	}
	var inviteConfig config.InviteConfig
	tool.SystemConfigSliceReflectToStruct(configs, &inviteConfig)
	if deps.Config != nil {
		deps.Config.Invite = inviteConfig
	}
}
