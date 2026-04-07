package initialize

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
)

func Register(ctx *svc.ServiceContext) {
	logger.Debug("Register config initialization")
	configs, err := ctx.SystemModel.GetRegisterConfig(context.Background())
	if err != nil {
		logger.Errorf("[Init Register Config] Get Register Config Error: %s", err.Error())
		return
	}
	var registerConfig config.RegisterConfig
	tool.SystemConfigSliceReflectToStruct(configs, &registerConfig)
	ctx.Config.Register = registerConfig
}
