package initialize

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/modules/infra/logger"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/auth"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
)

func Mobile(ctx *svc.ServiceContext) {
	logger.Debug("Mobile config initialization")
	method, err := ctx.AuthModel.FindOneByMethod(context.Background(), "mobile")
	if err != nil {
		panic(err)
	}
	var cfg config.MobileConfig
	var mobileConfig auth.MobileAuthConfig
	mobileConfig.Unmarshal(method.Config)
	tool.DeepCopy(&cfg, mobileConfig)
	cfg.Enable = *method.Enabled
	value, _ := json.Marshal(mobileConfig.PlatformConfig)
	cfg.PlatformConfig = string(value)
	ctx.Config.Mobile = cfg
}
