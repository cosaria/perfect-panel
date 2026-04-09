package common

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/services/report"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type GetGlobalConfigOutput struct {
	Body *types.GetGlobalConfigResponse
}

func GetGlobalConfigHandler(deps Deps) func(context.Context, *struct{}) (*GetGlobalConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetGlobalConfigOutput, error) {
		l := NewGetGlobalConfigLogic(ctx, deps)
		resp, err := l.GetGlobalConfig()
		if err != nil {
			return nil, err
		}
		return &GetGlobalConfigOutput{Body: resp}, nil
	}
}

type GetGlobalConfigLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get global config
func NewGetGlobalConfigLogic(ctx context.Context, deps Deps) *GetGlobalConfigLogic {
	return &GetGlobalConfigLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetGlobalConfigLogic) GetGlobalConfig() (resp *types.GetGlobalConfigResponse, err error) {
	resp = new(types.GetGlobalConfigResponse)
	cfg := l.deps.currentConfig()

	currencyCfg, err := l.deps.SystemModel.GetCurrencyConfig(l.ctx)
	if err != nil {
		l.Error("[GetGlobalConfigLogic] GetCurrencyConfig error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetCurrencyConfig error: %v", err.Error())
	}
	verifyCodeCfg, err := l.deps.SystemModel.GetVerifyCodeConfig(l.ctx)
	if err != nil {
		l.Error("[GetGlobalConfigLogic] GetVerifyCodeConfig error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetVerifyCodeConfig error: %v", err.Error())
	}

	tool.DeepCopy(&resp.Site, cfg.Site)
	tool.DeepCopy(&resp.Subscribe, cfg.Subscribe)
	tool.DeepCopy(&resp.Auth.Email, cfg.Email)
	tool.DeepCopy(&resp.Auth.Mobile, cfg.Mobile)
	tool.DeepCopy(&resp.Auth.Register, cfg.Register)
	tool.DeepCopy(&resp.Verify, cfg.Verify)
	tool.DeepCopy(&resp.Invite, cfg.Invite)
	tool.SystemConfigSliceReflectToStruct(currencyCfg, &resp.Currency)
	tool.SystemConfigSliceReflectToStruct(verifyCodeCfg, &resp.VerifyCode)

	if report.IsGatewayMode() {
		resp.Subscribe.SubscribePath = "/sub" + cfg.Subscribe.SubscribePath
	}

	resp.Verify = types.VeifyConfig{
		TurnstileSiteKey:          cfg.Verify.TurnstileSiteKey,
		EnableLoginVerify:         cfg.Verify.LoginVerify,
		EnableRegisterVerify:      cfg.Verify.RegisterVerify,
		EnableResetPasswordVerify: cfg.Verify.ResetPasswordVerify,
	}
	var methods []string

	// auth methods
	authMethods, err := l.deps.AuthModel.FindAll(l.ctx)
	if err != nil {
		l.Error("[GetGlobalConfigLogic] FindAll error: ", logger.Field("error", err.Error()))
	}

	for _, method := range authMethods {
		if *method.Enabled {
			methods = append(methods, method.Method)
			if method.Method == "device" {
				_ = json.Unmarshal([]byte(method.Config), &resp.Auth.Device)
				resp.Auth.Device.Enable = true
			}
		}
	}
	resp.OAuthMethods = methods

	webAds, err := l.deps.SystemModel.FindOneByKey(l.ctx, "WebAD")
	if err != nil {
		l.Error("[GetGlobalConfigLogic] FindOneByKey error: ", logger.Field("error", err.Error()), logger.Field("key", "WebAD"))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FindOneByKey error: %v", err.Error())
	}
	// web ads config
	resp.WebAd = webAds.Value == "true"
	return
}
