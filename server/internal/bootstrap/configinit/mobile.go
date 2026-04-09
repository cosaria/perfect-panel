package configinit

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/internal/platform/support/logger"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/persistence/auth"
	"github.com/perfect-panel/server/internal/platform/support/tool"
)

func Mobile(deps Deps) {
	logger.Debug("Mobile config initialization")
	method, err := deps.AuthModel.FindOneByMethod(context.Background(), "mobile")
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
	if deps.Config != nil {
		deps.Config.Mobile = cfg
	}
}
