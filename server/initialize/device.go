package initialize

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/auth"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
)

func Device(ctx *svc.ServiceContext) {
	logger.Debug("device config initialization")
	method, err := ctx.AuthModel.FindOneByMethod(context.Background(), "device")
	if err != nil {
		panic(err)
	}
	var cfg config.DeviceConfig
	var deviceConfig auth.DeviceConfig
	if err = deviceConfig.Unmarshal(method.Config); err != nil {
		panic(err)
	}
	tool.DeepCopy(&cfg, deviceConfig)
	cfg.Enable = *method.Enabled
	ctx.Config.Device = cfg
}
