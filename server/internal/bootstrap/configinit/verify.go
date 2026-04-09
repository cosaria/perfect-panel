package configinit

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/support/logger"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/support/tool"
)

type verifyConfig struct {
	TurnstileSiteKey          string
	TurnstileSecret           string
	EnableLoginVerify         bool
	EnableRegisterVerify      bool
	EnableResetPasswordVerify bool
}

func Verify(deps Deps) {
	logger.Debug("Verify config initialization")
	configs, err := deps.SystemModel.GetVerifyConfig(context.Background())
	if err != nil {
		logger.Error("[Init Verify Config] Get Verify Config Error: ", logger.Field("error", err.Error()))
		return
	}
	var verify verifyConfig
	tool.SystemConfigSliceReflectToStruct(configs, &verify)
	verifyCfg := config.Verify{
		TurnstileSiteKey:    verify.TurnstileSiteKey,
		TurnstileSecret:     verify.TurnstileSecret,
		LoginVerify:         verify.EnableLoginVerify,
		RegisterVerify:      verify.EnableRegisterVerify,
		ResetPasswordVerify: verify.EnableResetPasswordVerify,
	}
	if deps.Config != nil {
		deps.Config.Verify = verifyCfg
	}

	logger.Debug("Verify code config initialization")

	var verifyCodeConfig config.VerifyCode
	cfg, err := deps.SystemModel.GetVerifyCodeConfig(context.Background())
	if err != nil {
		logger.Errorf("[Init Verify Config] Get Verify Code Config Error: %s", err.Error())
		return
	}
	tool.SystemConfigSliceReflectToStruct(cfg, &verifyCodeConfig)
	if deps.Config != nil {
		deps.Config.VerifyCode = verifyCodeConfig
	}
}
