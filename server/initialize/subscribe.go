package initialize

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/util/tool"
)

func Subscribe(deps Deps) {
	logger.Debug("Subscribe config initialization")
	configs, err := deps.SystemModel.GetSubscribeConfig(context.Background())
	if err != nil {
		logger.Error("[Init Subscribe Config] Get Subscribe Config Error: ", logger.Field("error", err.Error()))
		return
	}

	var subscribeConfig config.SubscribeConfig
	tool.SystemConfigSliceReflectToStruct(configs, &subscribeConfig)
	if deps.Config != nil {
		deps.Config.Subscribe = subscribeConfig
	}
}
