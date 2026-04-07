package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/routers/middleware"
	appruntime "github.com/perfect-panel/server/runtime"
	common "github.com/perfect-panel/server/services/common"
)

func registerCommonRoutes(router *gin.Engine, runtimeDeps *appruntime.Deps, specOnly bool, apis *APIs) {
	commonGroup := router.Group("/api/v1/common")
	if !specOnly {
		commonGroup.Use(middleware.DeviceMiddleware(runtimeDeps))
	}
	commonConfig := governedAPIConfig("Perfect Panel Common API", "1.0.0", "/api/v1/common", "common")
	apis.Common = humagin.NewWithGroup(router, commonGroup, commonConfig)
	configureHumaAPI(apis.Common, compatibilityEnabled(runtimeDeps, specOnly))
	commonDeps := common.Deps{}
	if runtimeDeps != nil {
		commonDeps.AdsModel = runtimeDeps.AdsModel
		commonDeps.AuthModel = runtimeDeps.AuthModel
		commonDeps.ClientModel = runtimeDeps.ClientModel
		commonDeps.SystemModel = runtimeDeps.SystemModel
		commonDeps.UserModel = runtimeDeps.UserModel
		commonDeps.DB = runtimeDeps.DB
		commonDeps.Redis = runtimeDeps.Redis
		commonDeps.AuthLimiter = runtimeDeps.AuthLimiter
		commonDeps.Queue = runtimeDeps.Queue
		commonDeps.Config = runtimeDeps.Config
	}

	registerOperation(apis.Common, huma.Operation{
		OperationID: "getAds",
		Method:      http.MethodGet,
		Path:        "/ads",
		Summary:     "Get Ads",
		Tags:        []string{"common"},
	}, common.GetAdsHandler(commonDeps))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "checkVerificationCode",
		Method:      http.MethodPost,
		Path:        "/check_verification_code",
		Summary:     "Check verification code",
		Tags:        []string{"common"},
	}, common.CheckVerificationCodeHandler(commonDeps))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "getClient",
		Method:      http.MethodGet,
		Path:        "/client",
		Summary:     "Get Client",
		Tags:        []string{"common"},
	}, common.GetClientHandler(commonDeps))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "heartbeat",
		Method:      http.MethodGet,
		Path:        "/heartbeat",
		Summary:     "Heartbeat",
		Tags:        []string{"common"},
	}, common.HeartbeatHandler(commonDeps))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "sendEmailCode",
		Method:      http.MethodPost,
		Path:        "/send_code",
		Summary:     "Get verification code",
		Tags:        []string{"common"},
	}, common.SendEmailCodeHandler(commonDeps))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "sendSmsCode",
		Method:      http.MethodPost,
		Path:        "/send_sms_code",
		Summary:     "Get sms verification code",
		Tags:        []string{"common"},
	}, common.SendSmsCodeHandler(commonDeps))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "getGlobalConfig",
		Method:      http.MethodGet,
		Path:        "/site/config",
		Summary:     "Get global config",
		Tags:        []string{"common"},
	}, common.GetGlobalConfigHandler(commonDeps))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "getPrivacyPolicy",
		Method:      http.MethodGet,
		Path:        "/site/privacy",
		Summary:     "Get Privacy Policy",
		Tags:        []string{"common"},
	}, common.GetPrivacyPolicyHandler(commonDeps))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "getStat",
		Method:      http.MethodGet,
		Path:        "/site/stat",
		Summary:     "Get stat",
		Tags:        []string{"common"},
	}, common.GetStatHandler(commonDeps))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "getTos",
		Method:      http.MethodGet,
		Path:        "/site/tos",
		Summary:     "Get Tos Content",
		Tags:        []string{"common"},
	}, common.GetTosHandler(commonDeps))
}
