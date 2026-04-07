package handler

import (
	"encoding/json"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	adminAds "github.com/perfect-panel/server/routers/admin/ads"
	adminAnnouncement "github.com/perfect-panel/server/routers/admin/announcement"
	adminApplication "github.com/perfect-panel/server/routers/admin/application"
	adminAuthMethod "github.com/perfect-panel/server/routers/admin/authMethod"
	adminConsole "github.com/perfect-panel/server/routers/admin/console"
	adminCoupon "github.com/perfect-panel/server/routers/admin/coupon"
	adminDocument "github.com/perfect-panel/server/routers/admin/document"
	adminLog "github.com/perfect-panel/server/routers/admin/log"
	adminMarketing "github.com/perfect-panel/server/routers/admin/marketing"
	adminOrder "github.com/perfect-panel/server/routers/admin/order"
	adminPayment "github.com/perfect-panel/server/routers/admin/payment"
	adminServer "github.com/perfect-panel/server/routers/admin/server"
	adminSubscribe "github.com/perfect-panel/server/routers/admin/subscribe"
	adminSystem "github.com/perfect-panel/server/routers/admin/system"
	adminTicket "github.com/perfect-panel/server/routers/admin/ticket"
	adminTool "github.com/perfect-panel/server/routers/admin/tool"
	adminUser "github.com/perfect-panel/server/routers/admin/user"
	auth "github.com/perfect-panel/server/routers/auth"
	authOauth "github.com/perfect-panel/server/routers/auth/oauth"
	common "github.com/perfect-panel/server/routers/common"
	"github.com/perfect-panel/server/routers/middleware"
	publicAnnouncement "github.com/perfect-panel/server/routers/public/announcement"
	publicDocument "github.com/perfect-panel/server/routers/public/document"
	publicOrder "github.com/perfect-panel/server/routers/public/order"
	publicPayment "github.com/perfect-panel/server/routers/public/payment"
	publicPortal "github.com/perfect-panel/server/routers/public/portal"
	publicSubscribe "github.com/perfect-panel/server/routers/public/subscribe"
	publicTicket "github.com/perfect-panel/server/routers/public/ticket"
	publicUser "github.com/perfect-panel/server/routers/public/user"
	server "github.com/perfect-panel/server/routers/server"
	"github.com/perfect-panel/server/svc"
)

var bearerSecurity = []map[string][]string{{"bearer": {}}}

// apiConfig wraps apiConfig with $schema injection disabled.
// huma's default CreateHooks register a SchemaLinkTransformer that injects
// a "$schema" property into every response type — noise for SDK generation.
func apiConfig(title, version string) huma.Config {
	cfg := huma.DefaultConfig(title, version)
	cfg.CreateHooks = nil
	return cfg
}

func securitySchemes() map[string]*huma.SecurityScheme {
	return map[string]*huma.SecurityScheme{
		"bearer": {Type: "http", Scheme: "bearer", BearerFormat: "JWT"},
	}
}

// APIs holds all huma API instances for OpenAPI spec export.
type APIs struct {
	Admin    huma.API
	Common   huma.API
	userAPIs []huma.API // auth + public sub-APIs, merged via UserOpenAPI()
}

// UserOpenAPI merges all auth + public sub-API specs into a single OpenAPI spec.
func (a *APIs) UserOpenAPI() (map[string]interface{}, error) {
	merged := map[string]interface{}{
		"openapi": "3.1.0",
		"info":    map[string]interface{}{"title": "Perfect Panel User API", "version": "1.0.0"},
		"paths":   map[string]interface{}{},
		"components": map[string]interface{}{
			"schemas": map[string]interface{}{},
			"securitySchemes": map[string]interface{}{
				"bearer": map[string]interface{}{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "JWT",
				},
			},
		},
	}

	paths := merged["paths"].(map[string]interface{})
	schemas := merged["components"].(map[string]interface{})["schemas"].(map[string]interface{})

	for _, api := range a.userAPIs {
		data, err := json.Marshal(api.OpenAPI())
		if err != nil {
			return nil, err
		}
		var spec map[string]interface{}
		if err := json.Unmarshal(data, &spec); err != nil {
			return nil, err
		}

		// Extract server prefix for absolute path construction
		prefix := ""
		if servers, ok := spec["servers"].([]interface{}); ok && len(servers) > 0 {
			if s, ok := servers[0].(map[string]interface{}); ok {
				prefix, _ = s["url"].(string)
			}
		}

		if specPaths, ok := spec["paths"].(map[string]interface{}); ok {
			for path, item := range specPaths {
				paths[prefix+path] = item
			}
		}

		if comps, ok := spec["components"].(map[string]interface{}); ok {
			if specSchemas, ok := comps["schemas"].(map[string]interface{}); ok {
				for name, schema := range specSchemas {
					schemas[name] = schema
				}
			}
		}
	}

	return merged, nil
}

func RegisterHandlers(router *gin.Engine, serverCtx *svc.ServiceContext) {
	registerHandlers(router, serverCtx, false)
}

// RegisterHandlersForSpec registers only route metadata (no middleware, no server routes).
// Used by the openapi export command — serverCtx can be nil.
func RegisterHandlersForSpec(router *gin.Engine) *APIs {
	return registerHandlers(router, nil, true)
}

func registerHandlers(router *gin.Engine, serverCtx *svc.ServiceContext, specOnly bool) *APIs {
	apis := &APIs{}

	// ===== Admin API =====
	adminGroup := router.Group("/api/v1/admin")
	if !specOnly {
		adminGroup.Use(middleware.AuthMiddleware(serverCtx))
	}
	adminConfig := apiConfig("Perfect Panel Admin API", "1.0.0")
	adminConfig.Servers = []*huma.Server{{URL: "/api/v1/admin"}}
	adminConfig.Components.SecuritySchemes = securitySchemes()
	apis.Admin = humagin.NewWithGroup(router, adminGroup, adminConfig)

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createAds",
		Method:      http.MethodPost,
		Path:        "/ads",
		Summary:     "Create Ads",
		Tags:        []string{"ads"},
		Security:    bearerSecurity,
	}, adminAds.CreateAdsHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateAds",
		Method:      http.MethodPut,
		Path:        "/ads",
		Summary:     "Update Ads",
		Tags:        []string{"ads"},
		Security:    bearerSecurity,
	}, adminAds.UpdateAdsHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "deleteAds",
		Method:      http.MethodDelete,
		Path:        "/ads",
		Summary:     "Delete Ads",
		Tags:        []string{"ads"},
		Security:    bearerSecurity,
	}, adminAds.DeleteAdsHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getAdsDetail",
		Method:      http.MethodGet,
		Path:        "/ads/detail",
		Summary:     "Get Ads Detail",
		Tags:        []string{"ads"},
		Security:    bearerSecurity,
	}, adminAds.GetAdsDetailHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getAdsList",
		Method:      http.MethodGet,
		Path:        "/ads/list",
		Summary:     "Get Ads List",
		Tags:        []string{"ads"},
		Security:    bearerSecurity,
	}, adminAds.GetAdsListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createAnnouncement",
		Method:      http.MethodPost,
		Path:        "/announcement",
		Summary:     "Create announcement",
		Tags:        []string{"announcement"},
		Security:    bearerSecurity,
	}, adminAnnouncement.CreateAnnouncementHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateAnnouncement",
		Method:      http.MethodPut,
		Path:        "/announcement",
		Summary:     "Update announcement",
		Tags:        []string{"announcement"},
		Security:    bearerSecurity,
	}, adminAnnouncement.UpdateAnnouncementHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "deleteAnnouncement",
		Method:      http.MethodDelete,
		Path:        "/announcement",
		Summary:     "Delete announcement",
		Tags:        []string{"announcement"},
		Security:    bearerSecurity,
	}, adminAnnouncement.DeleteAnnouncementHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getAnnouncement",
		Method:      http.MethodGet,
		Path:        "/announcement/detail",
		Summary:     "Get announcement",
		Tags:        []string{"announcement"},
		Security:    bearerSecurity,
	}, adminAnnouncement.GetAnnouncementHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getAnnouncementList",
		Method:      http.MethodGet,
		Path:        "/announcement/list",
		Summary:     "Get announcement list",
		Tags:        []string{"announcement"},
		Security:    bearerSecurity,
	}, adminAnnouncement.GetAnnouncementListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createSubscribeApplication",
		Method:      http.MethodPost,
		Path:        "/application",
		Summary:     "Create subscribe application",
		Tags:        []string{"application"},
		Security:    bearerSecurity,
	}, adminApplication.CreateSubscribeApplicationHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "previewSubscribeTemplate",
		Method:      http.MethodGet,
		Path:        "/application/preview",
		Summary:     "Preview Template",
		Tags:        []string{"application"},
		Security:    bearerSecurity,
	}, adminApplication.PreviewSubscribeTemplateHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateSubscribeApplication",
		Method:      http.MethodPut,
		Path:        "/application/subscribe_application",
		Summary:     "Update subscribe application",
		Tags:        []string{"application"},
		Security:    bearerSecurity,
	}, adminApplication.UpdateSubscribeApplicationHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "deleteSubscribeApplication",
		Method:      http.MethodDelete,
		Path:        "/application/subscribe_application",
		Summary:     "Delete subscribe application",
		Tags:        []string{"application"},
		Security:    bearerSecurity,
	}, adminApplication.DeleteSubscribeApplicationHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getSubscribeApplicationList",
		Method:      http.MethodGet,
		Path:        "/application/subscribe_application_list",
		Summary:     "Get subscribe application list",
		Tags:        []string{"application"},
		Security:    bearerSecurity,
	}, adminApplication.GetSubscribeApplicationListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getAuthMethodConfig",
		Method:      http.MethodGet,
		Path:        "/auth-method/config",
		Summary:     "Get auth method config",
		Tags:        []string{"auth-method"},
		Security:    bearerSecurity,
	}, adminAuthMethod.GetAuthMethodConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateAuthMethodConfig",
		Method:      http.MethodPut,
		Path:        "/auth-method/config",
		Summary:     "Update auth method config",
		Tags:        []string{"auth-method"},
		Security:    bearerSecurity,
	}, adminAuthMethod.UpdateAuthMethodConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getEmailPlatform",
		Method:      http.MethodGet,
		Path:        "/auth-method/email_platform",
		Summary:     "Get email support platform",
		Tags:        []string{"auth-method"},
		Security:    bearerSecurity,
	}, adminAuthMethod.GetEmailPlatformHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getAuthMethodList",
		Method:      http.MethodGet,
		Path:        "/auth-method/list",
		Summary:     "Get auth method list",
		Tags:        []string{"auth-method"},
		Security:    bearerSecurity,
	}, adminAuthMethod.GetAuthMethodListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getSmsPlatform",
		Method:      http.MethodGet,
		Path:        "/auth-method/sms_platform",
		Summary:     "Get sms support platform",
		Tags:        []string{"auth-method"},
		Security:    bearerSecurity,
	}, adminAuthMethod.GetSmsPlatformHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "testEmailSend",
		Method:      http.MethodPost,
		Path:        "/auth-method/test_email_send",
		Summary:     "Test email send",
		Tags:        []string{"auth-method"},
		Security:    bearerSecurity,
	}, adminAuthMethod.TestEmailSendHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "testSmsSend",
		Method:      http.MethodPost,
		Path:        "/auth-method/test_sms_send",
		Summary:     "Test sms send",
		Tags:        []string{"auth-method"},
		Security:    bearerSecurity,
	}, adminAuthMethod.TestSmsSendHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "queryRevenueStatistics",
		Method:      http.MethodGet,
		Path:        "/console/revenue",
		Summary:     "Query revenue statistics",
		Tags:        []string{"console"},
		Security:    bearerSecurity,
	}, adminConsole.QueryRevenueStatisticsHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "queryServerTotalData",
		Method:      http.MethodGet,
		Path:        "/console/server",
		Summary:     "Query server total data",
		Tags:        []string{"console"},
		Security:    bearerSecurity,
	}, adminConsole.QueryServerTotalDataHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "queryTicketWaitReply",
		Method:      http.MethodGet,
		Path:        "/console/ticket",
		Summary:     "Query ticket wait reply",
		Tags:        []string{"console"},
		Security:    bearerSecurity,
	}, adminConsole.QueryTicketWaitReplyHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "queryUserStatistics",
		Method:      http.MethodGet,
		Path:        "/console/user",
		Summary:     "Query user statistics",
		Tags:        []string{"console"},
		Security:    bearerSecurity,
	}, adminConsole.QueryUserStatisticsHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createCoupon",
		Method:      http.MethodPost,
		Path:        "/coupon",
		Summary:     "Create coupon",
		Tags:        []string{"coupon"},
		Security:    bearerSecurity,
	}, adminCoupon.CreateCouponHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateCoupon",
		Method:      http.MethodPut,
		Path:        "/coupon",
		Summary:     "Update coupon",
		Tags:        []string{"coupon"},
		Security:    bearerSecurity,
	}, adminCoupon.UpdateCouponHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "deleteCoupon",
		Method:      http.MethodDelete,
		Path:        "/coupon",
		Summary:     "Delete coupon",
		Tags:        []string{"coupon"},
		Security:    bearerSecurity,
	}, adminCoupon.DeleteCouponHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "batchDeleteCoupon",
		Method:      http.MethodDelete,
		Path:        "/coupon/batch",
		Summary:     "Batch delete coupon",
		Tags:        []string{"coupon"},
		Security:    bearerSecurity,
	}, adminCoupon.BatchDeleteCouponHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getCouponList",
		Method:      http.MethodGet,
		Path:        "/coupon/list",
		Summary:     "Get coupon list",
		Tags:        []string{"coupon"},
		Security:    bearerSecurity,
	}, adminCoupon.GetCouponListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createDocument",
		Method:      http.MethodPost,
		Path:        "/document",
		Summary:     "Create document",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, adminDocument.CreateDocumentHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateDocument",
		Method:      http.MethodPut,
		Path:        "/document",
		Summary:     "Update document",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, adminDocument.UpdateDocumentHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "deleteDocument",
		Method:      http.MethodDelete,
		Path:        "/document",
		Summary:     "Delete document",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, adminDocument.DeleteDocumentHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "batchDeleteDocument",
		Method:      http.MethodDelete,
		Path:        "/document/batch",
		Summary:     "Batch delete document",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, adminDocument.BatchDeleteDocumentHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getDocumentDetail",
		Method:      http.MethodGet,
		Path:        "/document/detail",
		Summary:     "Get document detail",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, adminDocument.GetDocumentDetailHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getDocumentList",
		Method:      http.MethodGet,
		Path:        "/document/list",
		Summary:     "Get document list",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, adminDocument.GetDocumentListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "filterBalanceLog",
		Method:      http.MethodGet,
		Path:        "/log/balance/list",
		Summary:     "Filter balance log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterBalanceLogHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "filterCommissionLog",
		Method:      http.MethodGet,
		Path:        "/log/commission/list",
		Summary:     "Filter commission log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterCommissionLogHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "filterEmailLog",
		Method:      http.MethodGet,
		Path:        "/log/email/list",
		Summary:     "Filter email log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterEmailLogHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "filterGiftLog",
		Method:      http.MethodGet,
		Path:        "/log/gift/list",
		Summary:     "Filter gift log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterGiftLogHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "filterLoginLog",
		Method:      http.MethodGet,
		Path:        "/log/login/list",
		Summary:     "Filter login log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterLoginLogHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getMessageLogList",
		Method:      http.MethodGet,
		Path:        "/log/message/list",
		Summary:     "Get message log list",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.GetMessageLogListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "filterMobileLog",
		Method:      http.MethodGet,
		Path:        "/log/mobile/list",
		Summary:     "Filter mobile log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterMobileLogHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "filterRegisterLog",
		Method:      http.MethodGet,
		Path:        "/log/register/list",
		Summary:     "Filter register log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterRegisterLogHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "filterServerTrafficLog",
		Method:      http.MethodGet,
		Path:        "/log/server/traffic/list",
		Summary:     "Filter server traffic log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterServerTrafficLogHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getLogSetting",
		Method:      http.MethodGet,
		Path:        "/log/setting",
		Summary:     "Get log setting",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.GetLogSettingHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateLogSetting",
		Method:      http.MethodPost,
		Path:        "/log/setting",
		Summary:     "Update log setting",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.UpdateLogSettingHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "filterSubscribeLog",
		Method:      http.MethodGet,
		Path:        "/log/subscribe/list",
		Summary:     "Filter subscribe log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterSubscribeLogHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "filterResetSubscribeLog",
		Method:      http.MethodGet,
		Path:        "/log/subscribe/reset/list",
		Summary:     "Filter reset subscribe log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterResetSubscribeLogHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "filterUserSubscribeTrafficLog",
		Method:      http.MethodGet,
		Path:        "/log/subscribe/traffic/list",
		Summary:     "Filter user subscribe traffic log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterUserSubscribeTrafficLogHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "filterTrafficLogDetails",
		Method:      http.MethodGet,
		Path:        "/log/traffic/details",
		Summary:     "Filter traffic log details",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterTrafficLogDetailsHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getBatchSendEmailTaskList",
		Method:      http.MethodGet,
		Path:        "/marketing/email/batch/list",
		Summary:     "Get batch send email task list",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.GetBatchSendEmailTaskListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getPreSendEmailCount",
		Method:      http.MethodPost,
		Path:        "/marketing/email/batch/pre-send-count",
		Summary:     "Get pre-send email count",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.GetPreSendEmailCountHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createBatchSendEmailTask",
		Method:      http.MethodPost,
		Path:        "/marketing/email/batch/send",
		Summary:     "Create a batch send email task",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.CreateBatchSendEmailTaskHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getBatchSendEmailTaskStatus",
		Method:      http.MethodPost,
		Path:        "/marketing/email/batch/status",
		Summary:     "Get batch send email task status",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.GetBatchSendEmailTaskStatusHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "stopBatchSendEmailTask",
		Method:      http.MethodPost,
		Path:        "/marketing/email/batch/stop",
		Summary:     "Stop a batch send email task",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.StopBatchSendEmailTaskHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createQuotaTask",
		Method:      http.MethodPost,
		Path:        "/marketing/quota/create",
		Summary:     "Create a quota task",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.CreateQuotaTaskHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "queryQuotaTaskList",
		Method:      http.MethodGet,
		Path:        "/marketing/quota/list",
		Summary:     "Query quota task list",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.QueryQuotaTaskListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "queryQuotaTaskPreCount",
		Method:      http.MethodPost,
		Path:        "/marketing/quota/pre-count",
		Summary:     "Query quota task pre-count",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.QueryQuotaTaskPreCountHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createOrder",
		Method:      http.MethodPost,
		Path:        "/order",
		Summary:     "Create order",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, adminOrder.CreateOrderHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getOrderList",
		Method:      http.MethodGet,
		Path:        "/order/list",
		Summary:     "Get order list",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, adminOrder.GetOrderListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateOrderStatus",
		Method:      http.MethodPut,
		Path:        "/order/status",
		Summary:     "Update order status",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, adminOrder.UpdateOrderStatusHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createPaymentMethod",
		Method:      http.MethodPost,
		Path:        "/payment",
		Summary:     "Create Payment Method",
		Tags:        []string{"payment"},
		Security:    bearerSecurity,
	}, adminPayment.CreatePaymentMethodHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updatePaymentMethod",
		Method:      http.MethodPut,
		Path:        "/payment",
		Summary:     "Update Payment Method",
		Tags:        []string{"payment"},
		Security:    bearerSecurity,
	}, adminPayment.UpdatePaymentMethodHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "deletePaymentMethod",
		Method:      http.MethodDelete,
		Path:        "/payment",
		Summary:     "Delete Payment Method",
		Tags:        []string{"payment"},
		Security:    bearerSecurity,
	}, adminPayment.DeletePaymentMethodHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getPaymentMethodList",
		Method:      http.MethodGet,
		Path:        "/payment/list",
		Summary:     "Get Payment Method List",
		Tags:        []string{"payment"},
		Security:    bearerSecurity,
	}, adminPayment.GetPaymentMethodListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getPaymentPlatform",
		Method:      http.MethodGet,
		Path:        "/payment/platform",
		Summary:     "Get supported payment platform",
		Tags:        []string{"payment"},
		Security:    bearerSecurity,
	}, adminPayment.GetPaymentPlatformHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createServer",
		Method:      http.MethodPost,
		Path:        "/server/create",
		Summary:     "Create Server",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.CreateServerHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "deleteServer",
		Method:      http.MethodPost,
		Path:        "/server/delete",
		Summary:     "Delete Server",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.DeleteServerHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "filterServerList",
		Method:      http.MethodGet,
		Path:        "/server/list",
		Summary:     "Filter Server List",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.FilterServerListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createNode",
		Method:      http.MethodPost,
		Path:        "/server/node/create",
		Summary:     "Create Node",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.CreateNodeHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "deleteNode",
		Method:      http.MethodPost,
		Path:        "/server/node/delete",
		Summary:     "Delete Node",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.DeleteNodeHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "filterNodeList",
		Method:      http.MethodGet,
		Path:        "/server/node/list",
		Summary:     "Filter Node List",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.FilterNodeListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "resetSortWithNode",
		Method:      http.MethodPost,
		Path:        "/server/node/sort",
		Summary:     "Reset node sort",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.ResetSortWithNodeHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "toggleNodeStatus",
		Method:      http.MethodPost,
		Path:        "/server/node/status/toggle",
		Summary:     "Toggle Node Status",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.ToggleNodeStatusHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "queryNodeTag",
		Method:      http.MethodGet,
		Path:        "/server/node/tags",
		Summary:     "Query all node tags",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.QueryNodeTagHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateNode",
		Method:      http.MethodPost,
		Path:        "/server/node/update",
		Summary:     "Update Node",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.UpdateNodeHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getServerProtocols",
		Method:      http.MethodGet,
		Path:        "/server/protocols",
		Summary:     "Get Server Protocols",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.GetServerProtocolsHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "resetSortWithServer",
		Method:      http.MethodPost,
		Path:        "/server/server/sort",
		Summary:     "Reset server sort",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.ResetSortWithServerHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateServer",
		Method:      http.MethodPost,
		Path:        "/server/update",
		Summary:     "Update Server",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.UpdateServerHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createSubscribe",
		Method:      http.MethodPost,
		Path:        "/subscribe",
		Summary:     "Create subscribe",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.CreateSubscribeHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateSubscribe",
		Method:      http.MethodPut,
		Path:        "/subscribe",
		Summary:     "Update subscribe",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.UpdateSubscribeHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "deleteSubscribe",
		Method:      http.MethodDelete,
		Path:        "/subscribe",
		Summary:     "Delete subscribe",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.DeleteSubscribeHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "batchDeleteSubscribe",
		Method:      http.MethodDelete,
		Path:        "/subscribe/batch",
		Summary:     "Batch delete subscribe",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.BatchDeleteSubscribeHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getSubscribeDetails",
		Method:      http.MethodGet,
		Path:        "/subscribe/details",
		Summary:     "Get subscribe details",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.GetSubscribeDetailsHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createSubscribeGroup",
		Method:      http.MethodPost,
		Path:        "/subscribe/group",
		Summary:     "Create subscribe group",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.CreateSubscribeGroupHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateSubscribeGroup",
		Method:      http.MethodPut,
		Path:        "/subscribe/group",
		Summary:     "Update subscribe group",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.UpdateSubscribeGroupHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "deleteSubscribeGroup",
		Method:      http.MethodDelete,
		Path:        "/subscribe/group",
		Summary:     "Delete subscribe group",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.DeleteSubscribeGroupHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "batchDeleteSubscribeGroup",
		Method:      http.MethodDelete,
		Path:        "/subscribe/group/batch",
		Summary:     "Batch delete subscribe group",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.BatchDeleteSubscribeGroupHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getSubscribeGroupList",
		Method:      http.MethodGet,
		Path:        "/subscribe/group/list",
		Summary:     "Get subscribe group list",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.GetSubscribeGroupListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getSubscribeList",
		Method:      http.MethodGet,
		Path:        "/subscribe/list",
		Summary:     "Get subscribe list",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.GetSubscribeListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "resetAllSubscribeToken",
		Method:      http.MethodPost,
		Path:        "/subscribe/reset_all_token",
		Summary:     "Reset all subscribe tokens",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.ResetAllSubscribeTokenHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "subscribeSort",
		Method:      http.MethodPost,
		Path:        "/subscribe/sort",
		Summary:     "Subscribe sort",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.SubscribeSortHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getCurrencyConfig",
		Method:      http.MethodGet,
		Path:        "/system/currency_config",
		Summary:     "Get Currency Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetCurrencyConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateCurrencyConfig",
		Method:      http.MethodPut,
		Path:        "/system/currency_config",
		Summary:     "Update Currency Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateCurrencyConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getNodeMultiplier",
		Method:      http.MethodGet,
		Path:        "/system/get_node_multiplier",
		Summary:     "Get Node Multiplier",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetNodeMultiplierHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getInviteConfig",
		Method:      http.MethodGet,
		Path:        "/system/invite_config",
		Summary:     "Get invite config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetInviteConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateInviteConfig",
		Method:      http.MethodPut,
		Path:        "/system/invite_config",
		Summary:     "Update invite config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateInviteConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getModuleConfig",
		Method:      http.MethodGet,
		Path:        "/system/module",
		Summary:     "Get Module Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetModuleConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getNodeConfig",
		Method:      http.MethodGet,
		Path:        "/system/node_config",
		Summary:     "Get node config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetNodeConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateNodeConfig",
		Method:      http.MethodPut,
		Path:        "/system/node_config",
		Summary:     "Update node config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateNodeConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "preViewNodeMultiplier",
		Method:      http.MethodGet,
		Path:        "/system/node_multiplier/preview",
		Summary:     "PreView Node Multiplier",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.PreViewNodeMultiplierHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getPrivacyPolicyConfig",
		Method:      http.MethodGet,
		Path:        "/system/privacy",
		Summary:     "get Privacy Policy Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetPrivacyPolicyConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updatePrivacyPolicyConfig",
		Method:      http.MethodPut,
		Path:        "/system/privacy",
		Summary:     "Update Privacy Policy Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdatePrivacyPolicyConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getRegisterConfig",
		Method:      http.MethodGet,
		Path:        "/system/register_config",
		Summary:     "Get register config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetRegisterConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateRegisterConfig",
		Method:      http.MethodPut,
		Path:        "/system/register_config",
		Summary:     "Update register config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateRegisterConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "setNodeMultiplier",
		Method:      http.MethodPost,
		Path:        "/system/set_node_multiplier",
		Summary:     "Set Node Multiplier",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.SetNodeMultiplierHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "settingTelegramBot",
		Method:      http.MethodPost,
		Path:        "/system/setting_telegram_bot",
		Summary:     "setting telegram bot",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.SettingTelegramBotHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getSiteConfig",
		Method:      http.MethodGet,
		Path:        "/system/site_config",
		Summary:     "Get site config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetSiteConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateSiteConfig",
		Method:      http.MethodPut,
		Path:        "/system/site_config",
		Summary:     "Update site config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateSiteConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getSubscribeConfig",
		Method:      http.MethodGet,
		Path:        "/system/subscribe_config",
		Summary:     "Get subscribe config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetSubscribeConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateSubscribeConfig",
		Method:      http.MethodPut,
		Path:        "/system/subscribe_config",
		Summary:     "Update subscribe config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateSubscribeConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getTosConfig",
		Method:      http.MethodGet,
		Path:        "/system/tos_config",
		Summary:     "Get Team of Service Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetTosConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateTosConfig",
		Method:      http.MethodPut,
		Path:        "/system/tos_config",
		Summary:     "Update Team of Service Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateTosConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getVerifyCodeConfig",
		Method:      http.MethodGet,
		Path:        "/system/verify_code_config",
		Summary:     "Get Verify Code Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetVerifyCodeConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateVerifyCodeConfig",
		Method:      http.MethodPut,
		Path:        "/system/verify_code_config",
		Summary:     "Update Verify Code Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateVerifyCodeConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getVerifyConfig",
		Method:      http.MethodGet,
		Path:        "/system/verify_config",
		Summary:     "Get verify config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetVerifyConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateVerifyConfig",
		Method:      http.MethodPut,
		Path:        "/system/verify_config",
		Summary:     "Update verify config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateVerifyConfigHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateTicketStatus",
		Method:      http.MethodPut,
		Path:        "/ticket",
		Summary:     "Update ticket status",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, adminTicket.UpdateTicketStatusHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getTicket",
		Method:      http.MethodGet,
		Path:        "/ticket/detail",
		Summary:     "Get ticket detail",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, adminTicket.GetTicketHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createTicketFollow",
		Method:      http.MethodPost,
		Path:        "/ticket/follow",
		Summary:     "Create ticket follow",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, adminTicket.CreateTicketFollowHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getTicketList",
		Method:      http.MethodGet,
		Path:        "/ticket/list",
		Summary:     "Get ticket list",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, adminTicket.GetTicketListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "queryIPLocation",
		Method:      http.MethodGet,
		Path:        "/tool/ip/location",
		Summary:     "Query IP Location",
		Tags:        []string{"tool"},
		Security:    bearerSecurity,
	}, adminTool.QueryIPLocationHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getSystemLog",
		Method:      http.MethodGet,
		Path:        "/tool/log",
		Summary:     "Get System Log",
		Tags:        []string{"tool"},
		Security:    bearerSecurity,
	}, adminTool.GetSystemLogHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "restartSystem",
		Method:      http.MethodGet,
		Path:        "/tool/restart",
		Summary:     "Restart System",
		Tags:        []string{"tool"},
		Security:    bearerSecurity,
	}, adminTool.RestartSystemHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getVersion",
		Method:      http.MethodGet,
		Path:        "/tool/version",
		Summary:     "Get Version",
		Tags:        []string{"tool"},
		Security:    bearerSecurity,
	}, adminTool.GetVersionHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "deleteUser",
		Method:      http.MethodDelete,
		Path:        "/user",
		Summary:     "Delete user",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.DeleteUserHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createUser",
		Method:      http.MethodPost,
		Path:        "/user",
		Summary:     "Create user",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.CreateUserHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createUserAuthMethod",
		Method:      http.MethodPost,
		Path:        "/user/auth_method",
		Summary:     "Create user auth method",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.CreateUserAuthMethodHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "deleteUserAuthMethod",
		Method:      http.MethodDelete,
		Path:        "/user/auth_method",
		Summary:     "Delete user auth method",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.DeleteUserAuthMethodHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateUserAuthMethod",
		Method:      http.MethodPut,
		Path:        "/user/auth_method",
		Summary:     "Update user auth method",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.UpdateUserAuthMethodHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getUserAuthMethod",
		Method:      http.MethodGet,
		Path:        "/user/auth_method",
		Summary:     "Get user auth method",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserAuthMethodHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateUserBasicInfo",
		Method:      http.MethodPut,
		Path:        "/user/basic",
		Summary:     "Update user basic info",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.UpdateUserBasicInfoHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "batchDeleteUser",
		Method:      http.MethodDelete,
		Path:        "/user/batch",
		Summary:     "Batch delete user",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.BatchDeleteUserHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "currentUser",
		Method:      http.MethodGet,
		Path:        "/user/current",
		Summary:     "Current user",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.CurrentUserHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getUserDetail",
		Method:      http.MethodGet,
		Path:        "/user/detail",
		Summary:     "Get user detail",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserDetailHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateUserDevice",
		Method:      http.MethodPut,
		Path:        "/user/device",
		Summary:     "User device",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.UpdateUserDeviceHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "deleteUserDevice",
		Method:      http.MethodDelete,
		Path:        "/user/device",
		Summary:     "Delete user device",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.DeleteUserDeviceHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "kickOfflineByUserDevice",
		Method:      http.MethodPut,
		Path:        "/user/device/kick_offline",
		Summary:     "kick offline user device",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.KickOfflineByUserDeviceHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getUserList",
		Method:      http.MethodGet,
		Path:        "/user/list",
		Summary:     "Get user list",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserListHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getUserLoginLogs",
		Method:      http.MethodGet,
		Path:        "/user/login/logs",
		Summary:     "Get user login logs",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserLoginLogsHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateUserNotifySetting",
		Method:      http.MethodPut,
		Path:        "/user/notify",
		Summary:     "Update user notify setting",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.UpdateUserNotifySettingHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getUserSubscribe",
		Method:      http.MethodGet,
		Path:        "/user/subscribe",
		Summary:     "Get user subcribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserSubscribeHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "createUserSubscribe",
		Method:      http.MethodPost,
		Path:        "/user/subscribe",
		Summary:     "Create user subcribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.CreateUserSubscribeHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "updateUserSubscribe",
		Method:      http.MethodPut,
		Path:        "/user/subscribe",
		Summary:     "Update user subcribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.UpdateUserSubscribeHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "deleteUserSubscribe",
		Method:      http.MethodDelete,
		Path:        "/user/subscribe",
		Summary:     "Delete user subcribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.DeleteUserSubscribeHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getUserSubscribeById",
		Method:      http.MethodGet,
		Path:        "/user/subscribe/detail",
		Summary:     "Get user subcribe by id",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserSubscribeByIdHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getUserSubscribeDevices",
		Method:      http.MethodGet,
		Path:        "/user/subscribe/device",
		Summary:     "Get user subcribe devices",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserSubscribeDevicesHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getUserSubscribeLogs",
		Method:      http.MethodGet,
		Path:        "/user/subscribe/logs",
		Summary:     "Get user subcribe logs",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserSubscribeLogsHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getUserSubscribeResetTrafficLogs",
		Method:      http.MethodGet,
		Path:        "/user/subscribe/reset/logs",
		Summary:     "Get user subcribe reset traffic logs",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserSubscribeResetTrafficLogsHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "resetUserSubscribeToken",
		Method:      http.MethodPost,
		Path:        "/user/subscribe/reset/token",
		Summary:     "Reset user subscribe token",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.ResetUserSubscribeTokenHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "resetUserSubscribeTraffic",
		Method:      http.MethodPost,
		Path:        "/user/subscribe/reset/traffic",
		Summary:     "Reset user subscribe traffic",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.ResetUserSubscribeTrafficHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "toggleUserSubscribeStatus",
		Method:      http.MethodPost,
		Path:        "/user/subscribe/toggle",
		Summary:     "Stop user subscribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.ToggleUserSubscribeStatusHandler(serverCtx))

	huma.Register(apis.Admin, huma.Operation{
		OperationID: "getUserSubscribeTrafficLogs",
		Method:      http.MethodGet,
		Path:        "/user/subscribe/traffic_logs",
		Summary:     "Get user subcribe traffic logs",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserSubscribeTrafficLogsHandler(serverCtx))

	// ===== Common API =====
	commonGroup := router.Group("/api/v1/common")
	if !specOnly {
		commonGroup.Use(middleware.DeviceMiddleware(serverCtx))
	}
	commonConfig := apiConfig("Perfect Panel Common API", "1.0.0")
	commonConfig.Servers = []*huma.Server{{URL: "/api/v1/common"}}
	apis.Common = humagin.NewWithGroup(router, commonGroup, commonConfig)

	huma.Register(apis.Common, huma.Operation{
		OperationID: "getAds",
		Method:      http.MethodGet,
		Path:        "/ads",
		Summary:     "Get Ads",
		Tags:        []string{"common"},
	}, common.GetAdsHandler(serverCtx))

	huma.Register(apis.Common, huma.Operation{
		OperationID: "checkVerificationCode",
		Method:      http.MethodPost,
		Path:        "/check_verification_code",
		Summary:     "Check verification code",
		Tags:        []string{"common"},
	}, common.CheckVerificationCodeHandler(serverCtx))

	huma.Register(apis.Common, huma.Operation{
		OperationID: "getClient",
		Method:      http.MethodGet,
		Path:        "/client",
		Summary:     "Get Client",
		Tags:        []string{"common"},
	}, common.GetClientHandler(serverCtx))

	huma.Register(apis.Common, huma.Operation{
		OperationID: "heartbeat",
		Method:      http.MethodGet,
		Path:        "/heartbeat",
		Summary:     "Heartbeat",
		Tags:        []string{"common"},
	}, common.HeartbeatHandler(serverCtx))

	huma.Register(apis.Common, huma.Operation{
		OperationID: "sendEmailCode",
		Method:      http.MethodPost,
		Path:        "/send_code",
		Summary:     "Get verification code",
		Tags:        []string{"common"},
	}, common.SendEmailCodeHandler(serverCtx))

	huma.Register(apis.Common, huma.Operation{
		OperationID: "sendSmsCode",
		Method:      http.MethodPost,
		Path:        "/send_sms_code",
		Summary:     "Get sms verification code",
		Tags:        []string{"common"},
	}, common.SendSmsCodeHandler(serverCtx))

	huma.Register(apis.Common, huma.Operation{
		OperationID: "getGlobalConfig",
		Method:      http.MethodGet,
		Path:        "/site/config",
		Summary:     "Get global config",
		Tags:        []string{"common"},
	}, common.GetGlobalConfigHandler(serverCtx))

	huma.Register(apis.Common, huma.Operation{
		OperationID: "getPrivacyPolicy",
		Method:      http.MethodGet,
		Path:        "/site/privacy",
		Summary:     "Get Privacy Policy",
		Tags:        []string{"common"},
	}, common.GetPrivacyPolicyHandler(serverCtx))

	huma.Register(apis.Common, huma.Operation{
		OperationID: "getStat",
		Method:      http.MethodGet,
		Path:        "/site/stat",
		Summary:     "Get stat",
		Tags:        []string{"common"},
	}, common.GetStatHandler(serverCtx))

	huma.Register(apis.Common, huma.Operation{
		OperationID: "getTos",
		Method:      http.MethodGet,
		Path:        "/site/tos",
		Summary:     "Get Tos Content",
		Tags:        []string{"common"},
	}, common.GetTosHandler(serverCtx))

	// ===== User API (auth + auth/oauth + public) =====
	// Auth routes use DeviceMiddleware; public routes use Auth+Device.
	// Each middleware group uses its own gin.Group + humagin.NewWithGroup.
	// Sub-API specs are merged into one via APIs.UserOpenAPI().

	// Auth routes
	authGroup := router.Group("/api/v1/auth")
	if !specOnly {
		authGroup.Use(middleware.DeviceMiddleware(serverCtx))
	}
	authConfig := apiConfig("Auth API", "1.0.0")
	authConfig.Servers = []*huma.Server{{URL: "/api/v1/auth"}}
	authAPI := humagin.NewWithGroup(router, authGroup, authConfig)

	huma.Register(authAPI, huma.Operation{
		OperationID: "checkUser",
		Method:      http.MethodGet,
		Path:        "/check",
		Summary:     "Check user is exist",
		Tags:        []string{"auth"},
	}, auth.CheckUserHandler(serverCtx))

	huma.Register(authAPI, huma.Operation{
		OperationID: "checkUserTelephone",
		Method:      http.MethodGet,
		Path:        "/check/telephone",
		Summary:     "Check user telephone is exist",
		Tags:        []string{"auth"},
	}, auth.CheckUserTelephoneHandler(serverCtx))

	huma.Register(authAPI, huma.Operation{
		OperationID: "userLogin",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "User login",
		Tags:        []string{"auth"},
	}, auth.UserLoginHandler(serverCtx))

	huma.Register(authAPI, huma.Operation{
		OperationID: "deviceLogin",
		Method:      http.MethodPost,
		Path:        "/login/device",
		Summary:     "Device Login",
		Tags:        []string{"auth"},
	}, auth.DeviceLoginHandler(serverCtx))

	huma.Register(authAPI, huma.Operation{
		OperationID: "telephoneLogin",
		Method:      http.MethodPost,
		Path:        "/login/telephone",
		Summary:     "User Telephone login",
		Tags:        []string{"auth"},
	}, auth.TelephoneLoginHandler(serverCtx))

	huma.Register(authAPI, huma.Operation{
		OperationID: "userRegister",
		Method:      http.MethodPost,
		Path:        "/register",
		Summary:     "User register",
		Tags:        []string{"auth"},
	}, auth.UserRegisterHandler(serverCtx))

	huma.Register(authAPI, huma.Operation{
		OperationID: "telephoneUserRegister",
		Method:      http.MethodPost,
		Path:        "/register/telephone",
		Summary:     "User Telephone register",
		Tags:        []string{"auth"},
	}, auth.TelephoneUserRegisterHandler(serverCtx))

	huma.Register(authAPI, huma.Operation{
		OperationID: "resetPassword",
		Method:      http.MethodPost,
		Path:        "/reset",
		Summary:     "Reset password",
		Tags:        []string{"auth"},
	}, auth.ResetPasswordHandler(serverCtx))

	huma.Register(authAPI, huma.Operation{
		OperationID: "telephoneResetPassword",
		Method:      http.MethodPost,
		Path:        "/reset/telephone",
		Summary:     "Reset password",
		Tags:        []string{"auth"},
	}, auth.TelephoneResetPasswordHandler(serverCtx))

	// Auth OAuth routes (no middleware)
	authOauthGroup := router.Group("/api/v1/auth/oauth")
	authOauthConfig := apiConfig("Auth OAuth API", "1.0.0")
	authOauthConfig.Servers = []*huma.Server{{URL: "/api/v1/auth/oauth"}}
	authOauthAPI := humagin.NewWithGroup(router, authOauthGroup, authOauthConfig)

	// AppleLoginCallback is registered as a raw Gin handler because the logic
	// layer needs direct access to http.Request and http.ResponseWriter for
	// HTTP redirects, which cannot be expressed in huma's handler signature.
	authOauthGroup.POST("/callback/apple", authOauth.AppleLoginCallbackHandler(serverCtx))

	huma.Register(authOauthAPI, huma.Operation{
		OperationID: "oAuthLogin",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "OAuth login",
		Tags:        []string{"oauth"},
	}, authOauth.OAuthLoginHandler(serverCtx))

	huma.Register(authOauthAPI, huma.Operation{
		OperationID: "oAuthLoginGetToken",
		Method:      http.MethodPost,
		Path:        "/login/token",
		Summary:     "OAuth login get token",
		Tags:        []string{"oauth"},
	}, authOauth.OAuthLoginGetTokenHandler(serverCtx))

	publicAnnouncementGroup := router.Group("/api/v1/public/announcement")
	if !specOnly {
		publicAnnouncementGroup.Use(middleware.AuthMiddleware(serverCtx))
	}
	publicAnnouncementConfig := apiConfig("Public Announcement API", "1.0.0")
	publicAnnouncementConfig.Servers = []*huma.Server{{URL: "/api/v1/public/announcement"}}
	publicAnnouncementConfig.Components.SecuritySchemes = securitySchemes()
	publicAnnouncementAPI := humagin.NewWithGroup(router, publicAnnouncementGroup, publicAnnouncementConfig)

	huma.Register(publicAnnouncementAPI, huma.Operation{
		OperationID: "queryAnnouncement",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "Query announcement",
		Tags:        []string{"announcement"},
		Security:    bearerSecurity,
	}, publicAnnouncement.QueryAnnouncementHandler(serverCtx))

	publicDocumentGroup := router.Group("/api/v1/public/document")
	if !specOnly {
		publicDocumentGroup.Use(middleware.AuthMiddleware(serverCtx))
	}
	publicDocumentConfig := apiConfig("Public Document API", "1.0.0")
	publicDocumentConfig.Servers = []*huma.Server{{URL: "/api/v1/public/document"}}
	publicDocumentConfig.Components.SecuritySchemes = securitySchemes()
	publicDocumentAPI := humagin.NewWithGroup(router, publicDocumentGroup, publicDocumentConfig)

	huma.Register(publicDocumentAPI, huma.Operation{
		OperationID: "queryDocumentDetail",
		Method:      http.MethodGet,
		Path:        "/detail",
		Summary:     "Get document detail",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, publicDocument.QueryDocumentDetailHandler(serverCtx))

	huma.Register(publicDocumentAPI, huma.Operation{
		OperationID: "queryDocumentList",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "Get document list",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, publicDocument.QueryDocumentListHandler(serverCtx))

	publicOrderGroup := router.Group("/api/v1/public/order")
	if !specOnly {
		publicOrderGroup.Use(middleware.AuthMiddleware(serverCtx))
	}
	publicOrderConfig := apiConfig("Public Order API", "1.0.0")
	publicOrderConfig.Servers = []*huma.Server{{URL: "/api/v1/public/order"}}
	publicOrderConfig.Components.SecuritySchemes = securitySchemes()
	publicOrderAPI := humagin.NewWithGroup(router, publicOrderGroup, publicOrderConfig)

	huma.Register(publicOrderAPI, huma.Operation{
		OperationID: "closeOrder",
		Method:      http.MethodPost,
		Path:        "/close",
		Summary:     "Close order",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.CloseOrderHandler(serverCtx))

	huma.Register(publicOrderAPI, huma.Operation{
		OperationID: "queryOrderDetail",
		Method:      http.MethodGet,
		Path:        "/detail",
		Summary:     "Get order",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.QueryOrderDetailHandler(serverCtx))

	huma.Register(publicOrderAPI, huma.Operation{
		OperationID: "queryOrderList",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "Get order list",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.QueryOrderListHandler(serverCtx))

	huma.Register(publicOrderAPI, huma.Operation{
		OperationID: "preCreateOrder",
		Method:      http.MethodPost,
		Path:        "/pre",
		Summary:     "Pre create order",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.PreCreateOrderHandler(serverCtx))

	huma.Register(publicOrderAPI, huma.Operation{
		OperationID: "purchase",
		Method:      http.MethodPost,
		Path:        "/purchase",
		Summary:     "purchase Subscription",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.PurchaseHandler(serverCtx))

	huma.Register(publicOrderAPI, huma.Operation{
		OperationID: "recharge",
		Method:      http.MethodPost,
		Path:        "/recharge",
		Summary:     "Recharge",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.RechargeHandler(serverCtx))

	huma.Register(publicOrderAPI, huma.Operation{
		OperationID: "renewal",
		Method:      http.MethodPost,
		Path:        "/renewal",
		Summary:     "Renewal Subscription",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.RenewalHandler(serverCtx))

	huma.Register(publicOrderAPI, huma.Operation{
		OperationID: "resetTraffic",
		Method:      http.MethodPost,
		Path:        "/reset",
		Summary:     "Reset traffic",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.ResetTrafficHandler(serverCtx))

	publicPaymentGroup := router.Group("/api/v1/public/payment")
	if !specOnly {
		publicPaymentGroup.Use(middleware.AuthMiddleware(serverCtx))
	}
	publicPaymentConfig := apiConfig("Public Payment API", "1.0.0")
	publicPaymentConfig.Servers = []*huma.Server{{URL: "/api/v1/public/payment"}}
	publicPaymentConfig.Components.SecuritySchemes = securitySchemes()
	publicPaymentAPI := humagin.NewWithGroup(router, publicPaymentGroup, publicPaymentConfig)

	huma.Register(publicPaymentAPI, huma.Operation{
		OperationID: "getAvailablePaymentMethods",
		Method:      http.MethodGet,
		Path:        "/methods",
		Summary:     "Get available payment methods",
		Tags:        []string{"payment"},
		Security:    bearerSecurity,
	}, publicPayment.GetAvailablePaymentMethodsHandler(serverCtx))

	publicSubscribeGroup := router.Group("/api/v1/public/subscribe")
	if !specOnly {
		publicSubscribeGroup.Use(middleware.AuthMiddleware(serverCtx))
	}
	publicSubscribeConfig := apiConfig("Public Subscribe API", "1.0.0")
	publicSubscribeConfig.Servers = []*huma.Server{{URL: "/api/v1/public/subscribe"}}
	publicSubscribeConfig.Components.SecuritySchemes = securitySchemes()
	publicSubscribeAPI := humagin.NewWithGroup(router, publicSubscribeGroup, publicSubscribeConfig)

	huma.Register(publicSubscribeAPI, huma.Operation{
		OperationID: "querySubscribeList",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "Get subscribe list",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, publicSubscribe.QuerySubscribeListHandler(serverCtx))

	huma.Register(publicSubscribeAPI, huma.Operation{
		OperationID: "queryUserSubscribeNodeList",
		Method:      http.MethodGet,
		Path:        "/node/list",
		Summary:     "Get user subscribe node info",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, publicSubscribe.QueryUserSubscribeNodeListHandler(serverCtx))

	publicTicketGroup := router.Group("/api/v1/public/ticket")
	if !specOnly {
		publicTicketGroup.Use(middleware.AuthMiddleware(serverCtx))
	}
	publicTicketConfig := apiConfig("Public Ticket API", "1.0.0")
	publicTicketConfig.Servers = []*huma.Server{{URL: "/api/v1/public/ticket"}}
	publicTicketConfig.Components.SecuritySchemes = securitySchemes()
	publicTicketAPI := humagin.NewWithGroup(router, publicTicketGroup, publicTicketConfig)

	huma.Register(publicTicketAPI, huma.Operation{
		OperationID: "updateUserTicketStatus",
		Method:      http.MethodPut,
		Path:        "/ticket",
		Summary:     "Update ticket status",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, publicTicket.UpdateUserTicketStatusHandler(serverCtx))

	huma.Register(publicTicketAPI, huma.Operation{
		OperationID: "createUserTicket",
		Method:      http.MethodPost,
		Path:        "/ticket",
		Summary:     "Create ticket",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, publicTicket.CreateUserTicketHandler(serverCtx))

	huma.Register(publicTicketAPI, huma.Operation{
		OperationID: "getUserTicketDetails",
		Method:      http.MethodGet,
		Path:        "/detail",
		Summary:     "Get ticket detail",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, publicTicket.GetUserTicketDetailsHandler(serverCtx))

	huma.Register(publicTicketAPI, huma.Operation{
		OperationID: "createUserTicketFollow",
		Method:      http.MethodPost,
		Path:        "/follow",
		Summary:     "Create ticket follow",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, publicTicket.CreateUserTicketFollowHandler(serverCtx))

	huma.Register(publicTicketAPI, huma.Operation{
		OperationID: "getUserTicketList",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "Get ticket list",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, publicTicket.GetUserTicketListHandler(serverCtx))

	publicUserGroup := router.Group("/api/v1/public/user")
	if !specOnly {
		publicUserGroup.Use(middleware.AuthMiddleware(serverCtx))
	}
	publicUserConfig := apiConfig("Public User API", "1.0.0")
	publicUserConfig.Servers = []*huma.Server{{URL: "/api/v1/public/user"}}
	publicUserConfig.Components.SecuritySchemes = securitySchemes()
	publicUserAPI := humagin.NewWithGroup(router, publicUserGroup, publicUserConfig)

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "queryUserAffiliate",
		Method:      http.MethodGet,
		Path:        "/affiliate/count",
		Summary:     "Query User Affiliate Count",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.QueryUserAffiliateHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "queryUserAffiliateList",
		Method:      http.MethodGet,
		Path:        "/affiliate/list",
		Summary:     "Query User Affiliate List",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.QueryUserAffiliateListHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "queryUserBalanceLog",
		Method:      http.MethodGet,
		Path:        "/balance_log",
		Summary:     "Query User Balance Log",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.QueryUserBalanceLogHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "updateBindEmail",
		Method:      http.MethodPut,
		Path:        "/bind_email",
		Summary:     "Update Bind Email",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UpdateBindEmailHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "updateBindMobile",
		Method:      http.MethodPut,
		Path:        "/bind_mobile",
		Summary:     "Update Bind Mobile",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UpdateBindMobileHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "bindOAuth",
		Method:      http.MethodPost,
		Path:        "/bind_oauth",
		Summary:     "Bind OAuth",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.BindOAuthHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "bindOAuthCallback",
		Method:      http.MethodPost,
		Path:        "/bind_oauth/callback",
		Summary:     "Bind OAuth Callback",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.BindOAuthCallbackHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "bindTelegram",
		Method:      http.MethodGet,
		Path:        "/bind_telegram",
		Summary:     "Bind Telegram",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.BindTelegramHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "queryUserCommissionLog",
		Method:      http.MethodGet,
		Path:        "/commission_log",
		Summary:     "Query User Commission Log",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.QueryUserCommissionLogHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "commissionWithdraw",
		Method:      http.MethodPost,
		Path:        "/commission_withdraw",
		Summary:     "Commission Withdraw",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.CommissionWithdrawHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "getDeviceList",
		Method:      http.MethodGet,
		Path:        "/devices",
		Summary:     "Get Device List",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.GetDeviceListHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "queryUserInfo",
		Method:      http.MethodGet,
		Path:        "/info",
		Summary:     "Query User Info",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.QueryUserInfoHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "getLoginLog",
		Method:      http.MethodGet,
		Path:        "/login_log",
		Summary:     "Get Login Log",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.GetLoginLogHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "updateUserNotify",
		Method:      http.MethodPut,
		Path:        "/notify",
		Summary:     "Update User Notify",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UpdateUserNotifyHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "getOAuthMethods",
		Method:      http.MethodGet,
		Path:        "/oauth_methods",
		Summary:     "Get OAuth Methods",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.GetOAuthMethodsHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "updateUserPassword",
		Method:      http.MethodPut,
		Path:        "/password",
		Summary:     "Update User Password",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UpdateUserPasswordHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "updateUserRules",
		Method:      http.MethodPut,
		Path:        "/rules",
		Summary:     "Update User Rules",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UpdateUserRulesHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "queryUserSubscribe",
		Method:      http.MethodGet,
		Path:        "/subscribe",
		Summary:     "Query User Subscribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.QueryUserSubscribeHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "getSubscribeLog",
		Method:      http.MethodGet,
		Path:        "/subscribe_log",
		Summary:     "Get Subscribe Log",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.GetSubscribeLogHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "updateUserSubscribeNote",
		Method:      http.MethodPut,
		Path:        "/subscribe_note",
		Summary:     "Update User Subscribe Note",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UpdateUserSubscribeNoteHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "resetUserSubscribeToken",
		Method:      http.MethodPut,
		Path:        "/subscribe_token",
		Summary:     "Reset User Subscribe Token",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.ResetUserSubscribeTokenHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "unbindDevice",
		Method:      http.MethodPut,
		Path:        "/unbind_device",
		Summary:     "Unbind Device",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UnbindDeviceHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "unbindOAuth",
		Method:      http.MethodPost,
		Path:        "/unbind_oauth",
		Summary:     "Unbind OAuth",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UnbindOAuthHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "unbindTelegram",
		Method:      http.MethodPost,
		Path:        "/unbind_telegram",
		Summary:     "Unbind Telegram",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UnbindTelegramHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "unsubscribe",
		Method:      http.MethodPost,
		Path:        "/unsubscribe",
		Summary:     "Unsubscribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UnsubscribeHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "preUnsubscribe",
		Method:      http.MethodPost,
		Path:        "/unsubscribe/pre",
		Summary:     "Pre Unsubscribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.PreUnsubscribeHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "verifyEmail",
		Method:      http.MethodPost,
		Path:        "/verify_email",
		Summary:     "Verify Email",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.VerifyEmailHandler(serverCtx))

	huma.Register(publicUserAPI, huma.Operation{
		OperationID: "queryWithdrawalLog",
		Method:      http.MethodGet,
		Path:        "/withdrawal_log",
		Summary:     "Query Withdrawal Log",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.QueryWithdrawalLogHandler(serverCtx))

	// Portal routes (DeviceMiddleware only)
	portalGroup := router.Group("/api/v1/public/portal")
	if !specOnly {
		portalGroup.Use(middleware.DeviceMiddleware(serverCtx))
	}
	portalConfig := apiConfig("Portal API", "1.0.0")
	portalConfig.Servers = []*huma.Server{{URL: "/api/v1/public/portal"}}
	portalAPI := humagin.NewWithGroup(router, portalGroup, portalConfig)

	huma.Register(portalAPI, huma.Operation{
		OperationID: "purchaseCheckout",
		Method:      http.MethodPost,
		Path:        "/order/checkout",
		Summary:     "Purchase Checkout",
		Tags:        []string{"portal"},
	}, publicPortal.PurchaseCheckoutHandler(serverCtx))

	huma.Register(portalAPI, huma.Operation{
		OperationID: "queryPurchaseOrder",
		Method:      http.MethodGet,
		Path:        "/order/status",
		Summary:     "Query Purchase Order",
		Tags:        []string{"portal"},
	}, publicPortal.QueryPurchaseOrderHandler(serverCtx))

	huma.Register(portalAPI, huma.Operation{
		OperationID: "portalGetAvailablePaymentMethods",
		Method:      http.MethodGet,
		Path:        "/payment-method",
		Summary:     "Get available payment methods",
		Tags:        []string{"portal"},
	}, publicPortal.GetAvailablePaymentMethodsHandler(serverCtx))

	huma.Register(portalAPI, huma.Operation{
		OperationID: "prePurchaseOrder",
		Method:      http.MethodPost,
		Path:        "/pre",
		Summary:     "Pre Purchase Order",
		Tags:        []string{"portal"},
	}, publicPortal.PrePurchaseOrderHandler(serverCtx))

	huma.Register(portalAPI, huma.Operation{
		OperationID: "portalPurchase",
		Method:      http.MethodPost,
		Path:        "/purchase",
		Summary:     "Purchase subscription",
		Tags:        []string{"portal"},
	}, publicPortal.PurchaseHandler(serverCtx))

	huma.Register(portalAPI, huma.Operation{
		OperationID: "getSubscription",
		Method:      http.MethodGet,
		Path:        "/subscribe",
		Summary:     "Get Subscription",
		Tags:        []string{"portal"},
	}, publicPortal.GetSubscriptionHandler(serverCtx))

	// ===== Server routes (raw Gin, no OpenAPI) =====
	if !specOnly {
		v1_serverGroup := router.Group("/api/v1/server")
		v1_serverGroup.Use(middleware.ServerMiddleware(serverCtx))

		v1_serverGroup.GET("/config", server.GetServerConfigHandler(serverCtx))
		v1_serverGroup.POST("/online", server.PushOnlineUsersHandler(serverCtx))
		v1_serverGroup.POST("/push", server.ServerPushUserTrafficHandler(serverCtx))
		v1_serverGroup.POST("/status", server.ServerPushStatusHandler(serverCtx))
		v1_serverGroup.GET("/user", server.GetServerUserListHandler(serverCtx))

	}

	apis.userAPIs = []huma.API{
		authAPI, authOauthAPI,
		publicAnnouncementAPI, publicDocumentAPI, publicOrderAPI,
		publicPaymentAPI, publicSubscribeAPI, publicTicketAPI,
		publicUserAPI, portalAPI,
	}

	return apis
}
