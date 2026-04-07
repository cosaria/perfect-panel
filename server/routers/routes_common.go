package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	common "github.com/perfect-panel/server/services/common"
	"github.com/perfect-panel/server/routers/middleware"
	"github.com/perfect-panel/server/svc"
)

func registerCommonRoutes(router *gin.Engine, serverCtx *svc.ServiceContext, specOnly bool, apis *APIs) {
	commonGroup := router.Group("/api/v1/common")
	if !specOnly {
		commonGroup.Use(middleware.DeviceMiddleware(serverCtx))
	}
	commonConfig := apiConfig("Perfect Panel Common API", "1.0.0")
	commonConfig.Servers = []*huma.Server{{URL: "/api/v1/common"}}
	apis.Common = humagin.NewWithGroup(router, commonGroup, commonConfig)
	configureHumaAPI(apis.Common, compatibilityEnabled(serverCtx, specOnly))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "getAds",
		Method:      http.MethodGet,
		Path:        "/ads",
		Summary:     "Get Ads",
		Tags:        []string{"common"},
	}, common.GetAdsHandler(serverCtx))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "checkVerificationCode",
		Method:      http.MethodPost,
		Path:        "/check_verification_code",
		Summary:     "Check verification code",
		Tags:        []string{"common"},
	}, common.CheckVerificationCodeHandler(serverCtx))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "getClient",
		Method:      http.MethodGet,
		Path:        "/client",
		Summary:     "Get Client",
		Tags:        []string{"common"},
	}, common.GetClientHandler(serverCtx))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "heartbeat",
		Method:      http.MethodGet,
		Path:        "/heartbeat",
		Summary:     "Heartbeat",
		Tags:        []string{"common"},
	}, common.HeartbeatHandler(serverCtx))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "sendEmailCode",
		Method:      http.MethodPost,
		Path:        "/send_code",
		Summary:     "Get verification code",
		Tags:        []string{"common"},
	}, common.SendEmailCodeHandler(serverCtx))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "sendSmsCode",
		Method:      http.MethodPost,
		Path:        "/send_sms_code",
		Summary:     "Get sms verification code",
		Tags:        []string{"common"},
	}, common.SendSmsCodeHandler(serverCtx))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "getGlobalConfig",
		Method:      http.MethodGet,
		Path:        "/site/config",
		Summary:     "Get global config",
		Tags:        []string{"common"},
	}, common.GetGlobalConfigHandler(serverCtx))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "getPrivacyPolicy",
		Method:      http.MethodGet,
		Path:        "/site/privacy",
		Summary:     "Get Privacy Policy",
		Tags:        []string{"common"},
	}, common.GetPrivacyPolicyHandler(serverCtx))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "getStat",
		Method:      http.MethodGet,
		Path:        "/site/stat",
		Summary:     "Get stat",
		Tags:        []string{"common"},
	}, common.GetStatHandler(serverCtx))

	registerOperation(apis.Common, huma.Operation{
		OperationID: "getTos",
		Method:      http.MethodGet,
		Path:        "/site/tos",
		Summary:     "Get Tos Content",
		Tags:        []string{"common"},
	}, common.GetTosHandler(serverCtx))
}
