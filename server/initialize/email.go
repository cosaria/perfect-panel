package initialize

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/auth"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
)

// Email get email smtp config
func Email(ctx *svc.ServiceContext) {
	logger.Debug("Email config initialization")
	method, err := ctx.AuthModel.FindOneByMethod(context.Background(), "email")
	if err != nil {
		panic(fmt.Sprintf("[Error] Initialization Failed to find email auth method: %v", err.Error()))
	}
	var cfg config.EmailConfig
	var emailConfig = new(auth.EmailAuthConfig)
	emailConfig.Unmarshal(method.Config)
	tool.DeepCopy(&cfg, emailConfig)
	cfg.Enable = *method.Enabled
	value, _ := json.Marshal(emailConfig.PlatformConfig)
	cfg.PlatformConfig = string(value)
	ctx.Config.Email = cfg
}
