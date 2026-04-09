package configinit

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/support/logger"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/support/tool"
)

func Register(deps Deps) {
	logger.Debug("Register config initialization")
	configs, err := deps.SystemModel.GetRegisterConfig(context.Background())
	if err != nil {
		logger.Errorf("[Init Register Config] Get Register Config Error: %s", err.Error())
		return
	}
	var registerConfig config.RegisterConfig
	tool.SystemConfigSliceReflectToStruct(configs, &registerConfig)
	if deps.Config != nil {
		deps.Config.Register = registerConfig
	}
}
