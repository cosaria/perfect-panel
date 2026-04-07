package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/initialize"
	"github.com/perfect-panel/server/routers/middleware"
	appruntime "github.com/perfect-panel/server/runtime"
	adminAds "github.com/perfect-panel/server/services/admin/ads"
	adminAnnouncement "github.com/perfect-panel/server/services/admin/announcement"
	adminApplication "github.com/perfect-panel/server/services/admin/application"
	adminAuthMethod "github.com/perfect-panel/server/services/admin/authMethod"
	adminConsole "github.com/perfect-panel/server/services/admin/console"
	adminCoupon "github.com/perfect-panel/server/services/admin/coupon"
	adminDocument "github.com/perfect-panel/server/services/admin/document"
	adminLog "github.com/perfect-panel/server/services/admin/log"
	adminMarketing "github.com/perfect-panel/server/services/admin/marketing"
	adminOrder "github.com/perfect-panel/server/services/admin/order"
	adminPayment "github.com/perfect-panel/server/services/admin/payment"
	adminServer "github.com/perfect-panel/server/services/admin/server"
	adminSubscribe "github.com/perfect-panel/server/services/admin/subscribe"
	adminSystem "github.com/perfect-panel/server/services/admin/system"
	adminTicket "github.com/perfect-panel/server/services/admin/ticket"
	adminTool "github.com/perfect-panel/server/services/admin/tool"
	adminUser "github.com/perfect-panel/server/services/admin/user"
)

func registerAdminRoutes(router *gin.Engine, runtimeDeps *appruntime.Deps, specOnly bool, apis *APIs) {
	// ===== Admin API =====
	adminGroup := router.Group("/api/v1/admin")
	if !specOnly {
		adminGroup.Use(middleware.AuthMiddleware(runtimeDeps))
	}
	adminConfig := governedAPIConfig(
		"Perfect Panel Admin API",
		"1.0.0",
		"/api/v1/admin",
		"ads",
		"announcement",
		"application",
		"auth-method",
		"console",
		"coupon",
		"document",
		"log",
		"marketing",
		"order",
		"payment",
		"subscribe",
		"system",
		"ticket",
		"tool",
		"user",
	)
	apis.Admin = humagin.NewWithGroup(router, adminGroup, adminConfig)
	configureHumaAPI(apis.Admin, compatibilityEnabled(runtimeDeps, specOnly))
	adminAdsDeps := adminAds.Deps{}
	adminAnnouncementDeps := adminAnnouncement.Deps{}
	adminApplicationDeps := adminApplication.Deps{}
	adminAuthMethodDeps := adminAuthMethod.Deps{}
	adminConsoleDeps := adminConsole.Deps{}
	adminCouponDeps := adminCoupon.Deps{}
	adminDocumentDeps := adminDocument.Deps{}
	adminLogDeps := adminLog.Deps{}
	adminMarketingDeps := adminMarketing.Deps{}
	adminOrderDeps := adminOrder.Deps{}
	adminPaymentDeps := adminPayment.Deps{}
	adminServerDeps := adminServer.Deps{}
	initDeps := initializeDepsFromRuntimeDeps(runtimeDeps)
	adminSystemDeps := newAdminSystemDeps(runtimeDeps, initDeps)
	adminSubscribeDeps := adminSubscribe.Deps{}
	adminTicketDeps := adminTicket.Deps{}
	adminToolDeps := newAdminToolDeps(runtimeDeps)
	adminUserDeps := adminUser.Deps{}
	if runtimeDeps != nil {
		adminAdsDeps.AdsModel = runtimeDeps.AdsModel
		adminAnnouncementDeps.AnnouncementModel = runtimeDeps.AnnouncementModel
		adminApplicationDeps.ClientModel = runtimeDeps.ClientModel
		adminApplicationDeps.NodeModel = runtimeDeps.NodeModel
		adminAuthMethodDeps.AuthModel = runtimeDeps.AuthModel
		adminAuthMethodDeps.Config = runtimeDeps.Config
		adminAuthMethodDeps.ReloadEmail = func() {
			initialize.Email(initDeps)
		}
		adminAuthMethodDeps.ReloadMobile = func() {
			initialize.Mobile(initDeps)
		}
		adminAuthMethodDeps.ReloadDevice = func() {
			initialize.Device(initDeps)
		}
		adminConsoleDeps.OrderModel = runtimeDeps.OrderModel
		adminConsoleDeps.UserModel = runtimeDeps.UserModel
		adminConsoleDeps.NodeModel = runtimeDeps.NodeModel
		adminConsoleDeps.TicketModel = runtimeDeps.TicketModel
		adminConsoleDeps.DB = runtimeDeps.DB
		adminCouponDeps.CouponModel = runtimeDeps.CouponModel
		adminDocumentDeps.DocumentModel = runtimeDeps.DocumentModel
		adminLogDeps.LogModel = runtimeDeps.LogModel
		adminLogDeps.SystemModel = runtimeDeps.SystemModel
		adminLogDeps.DB = runtimeDeps.DB
		adminLogDeps.Config = runtimeDeps.Config
		adminMarketingDeps.DB = runtimeDeps.DB
		adminMarketingDeps.Queue = runtimeDeps.Queue
		adminOrderDeps.OrderModel = runtimeDeps.OrderModel
		adminOrderDeps.PaymentModel = runtimeDeps.PaymentModel
		adminOrderDeps.Queue = runtimeDeps.Queue
		adminPaymentDeps.PaymentModel = runtimeDeps.PaymentModel
		adminPaymentDeps.Config = runtimeDeps.Config
		adminServerDeps.NodeModel = runtimeDeps.NodeModel
		adminServerDeps.UserModel = runtimeDeps.UserModel
		adminServerDeps.DB = runtimeDeps.DB
		adminSubscribeDeps.SubscribeModel = runtimeDeps.SubscribeModel
		adminSubscribeDeps.UserModel = runtimeDeps.UserModel
		adminSubscribeDeps.DB = runtimeDeps.DB
		adminSubscribeDeps.DeviceManager = runtimeDeps.DeviceManager
		adminTicketDeps.TicketModel = runtimeDeps.TicketModel
		adminUserDeps.UserModel = runtimeDeps.UserModel
		adminUserDeps.SubscribeModel = runtimeDeps.SubscribeModel
		adminUserDeps.LogModel = runtimeDeps.LogModel
		adminUserDeps.TrafficLogModel = runtimeDeps.TrafficLogModel
		adminUserDeps.DeviceManager = runtimeDeps.DeviceManager
		adminUserDeps.Config = runtimeDeps.Config
	}

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createAds",
		Method:      http.MethodPost,
		Path:        "/ads",
		Summary:     "Create Ads",
		Tags:        []string{"ads"},
		Security:    bearerSecurity,
	}, adminAds.CreateAdsHandler(adminAdsDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateAds",
		Method:      http.MethodPut,
		Path:        "/ads",
		Summary:     "Update Ads",
		Tags:        []string{"ads"},
		Security:    bearerSecurity,
	}, adminAds.UpdateAdsHandler(adminAdsDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "deleteAds",
		Method:      http.MethodDelete,
		Path:        "/ads",
		Summary:     "Delete Ads",
		Tags:        []string{"ads"},
		Security:    bearerSecurity,
	}, adminAds.DeleteAdsHandler(adminAdsDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getAdsDetail",
		Method:      http.MethodGet,
		Path:        "/ads/detail",
		Summary:     "Get Ads Detail",
		Tags:        []string{"ads"},
		Security:    bearerSecurity,
	}, adminAds.GetAdsDetailHandler(adminAdsDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getAdsList",
		Method:      http.MethodGet,
		Path:        "/ads/list",
		Summary:     "Get Ads List",
		Tags:        []string{"ads"},
		Security:    bearerSecurity,
	}, adminAds.GetAdsListHandler(adminAdsDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createAnnouncement",
		Method:      http.MethodPost,
		Path:        "/announcement",
		Summary:     "Create announcement",
		Tags:        []string{"announcement"},
		Security:    bearerSecurity,
	}, adminAnnouncement.CreateAnnouncementHandler(adminAnnouncementDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateAnnouncement",
		Method:      http.MethodPut,
		Path:        "/announcement",
		Summary:     "Update announcement",
		Tags:        []string{"announcement"},
		Security:    bearerSecurity,
	}, adminAnnouncement.UpdateAnnouncementHandler(adminAnnouncementDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "deleteAnnouncement",
		Method:      http.MethodDelete,
		Path:        "/announcement",
		Summary:     "Delete announcement",
		Tags:        []string{"announcement"},
		Security:    bearerSecurity,
	}, adminAnnouncement.DeleteAnnouncementHandler(adminAnnouncementDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getAnnouncement",
		Method:      http.MethodGet,
		Path:        "/announcement/detail",
		Summary:     "Get announcement",
		Tags:        []string{"announcement"},
		Security:    bearerSecurity,
	}, adminAnnouncement.GetAnnouncementHandler(adminAnnouncementDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getAnnouncementList",
		Method:      http.MethodGet,
		Path:        "/announcement/list",
		Summary:     "Get announcement list",
		Tags:        []string{"announcement"},
		Security:    bearerSecurity,
	}, adminAnnouncement.GetAnnouncementListHandler(adminAnnouncementDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createSubscribeApplication",
		Method:      http.MethodPost,
		Path:        "/application",
		Summary:     "Create subscribe application",
		Tags:        []string{"application"},
		Security:    bearerSecurity,
	}, adminApplication.CreateSubscribeApplicationHandler(adminApplicationDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "previewSubscribeTemplate",
		Method:      http.MethodGet,
		Path:        "/application/preview",
		Summary:     "Preview Template",
		Tags:        []string{"application"},
		Security:    bearerSecurity,
	}, adminApplication.PreviewSubscribeTemplateHandler(adminApplicationDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateSubscribeApplication",
		Method:      http.MethodPut,
		Path:        "/application/subscribe_application",
		Summary:     "Update subscribe application",
		Tags:        []string{"application"},
		Security:    bearerSecurity,
	}, adminApplication.UpdateSubscribeApplicationHandler(adminApplicationDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "deleteSubscribeApplication",
		Method:      http.MethodDelete,
		Path:        "/application/subscribe_application",
		Summary:     "Delete subscribe application",
		Tags:        []string{"application"},
		Security:    bearerSecurity,
	}, adminApplication.DeleteSubscribeApplicationHandler(adminApplicationDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getSubscribeApplicationList",
		Method:      http.MethodGet,
		Path:        "/application/subscribe_application_list",
		Summary:     "Get subscribe application list",
		Tags:        []string{"application"},
		Security:    bearerSecurity,
	}, adminApplication.GetSubscribeApplicationListHandler(adminApplicationDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getAuthMethodConfig",
		Method:      http.MethodGet,
		Path:        "/auth-method/config",
		Summary:     "Get auth method config",
		Tags:        []string{"auth-method"},
		Security:    bearerSecurity,
	}, adminAuthMethod.GetAuthMethodConfigHandler(adminAuthMethodDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateAuthMethodConfig",
		Method:      http.MethodPut,
		Path:        "/auth-method/config",
		Summary:     "Update auth method config",
		Tags:        []string{"auth-method"},
		Security:    bearerSecurity,
	}, adminAuthMethod.UpdateAuthMethodConfigHandler(adminAuthMethodDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getEmailPlatform",
		Method:      http.MethodGet,
		Path:        "/auth-method/email_platform",
		Summary:     "Get email support platform",
		Tags:        []string{"auth-method"},
		Security:    bearerSecurity,
	}, adminAuthMethod.GetEmailPlatformHandler(adminAuthMethodDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getAuthMethodList",
		Method:      http.MethodGet,
		Path:        "/auth-method/list",
		Summary:     "Get auth method list",
		Tags:        []string{"auth-method"},
		Security:    bearerSecurity,
	}, adminAuthMethod.GetAuthMethodListHandler(adminAuthMethodDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getSmsPlatform",
		Method:      http.MethodGet,
		Path:        "/auth-method/sms_platform",
		Summary:     "Get sms support platform",
		Tags:        []string{"auth-method"},
		Security:    bearerSecurity,
	}, adminAuthMethod.GetSmsPlatformHandler(adminAuthMethodDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "testEmailSend",
		Method:      http.MethodPost,
		Path:        "/auth-method/test_email_send",
		Summary:     "Test email send",
		Tags:        []string{"auth-method"},
		Security:    bearerSecurity,
	}, adminAuthMethod.TestEmailSendHandler(adminAuthMethodDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "testSmsSend",
		Method:      http.MethodPost,
		Path:        "/auth-method/test_sms_send",
		Summary:     "Test sms send",
		Tags:        []string{"auth-method"},
		Security:    bearerSecurity,
	}, adminAuthMethod.TestSmsSendHandler(adminAuthMethodDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "queryRevenueStatistics",
		Method:      http.MethodGet,
		Path:        "/console/revenue",
		Summary:     "Query revenue statistics",
		Tags:        []string{"console"},
		Security:    bearerSecurity,
	}, adminConsole.QueryRevenueStatisticsHandler(adminConsoleDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "queryServerTotalData",
		Method:      http.MethodGet,
		Path:        "/console/server",
		Summary:     "Query server total data",
		Tags:        []string{"console"},
		Security:    bearerSecurity,
	}, adminConsole.QueryServerTotalDataHandler(adminConsoleDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "queryTicketWaitReply",
		Method:      http.MethodGet,
		Path:        "/console/ticket",
		Summary:     "Query ticket wait reply",
		Tags:        []string{"console"},
		Security:    bearerSecurity,
	}, adminConsole.QueryTicketWaitReplyHandler(adminConsoleDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "queryUserStatistics",
		Method:      http.MethodGet,
		Path:        "/console/user",
		Summary:     "Query user statistics",
		Tags:        []string{"console"},
		Security:    bearerSecurity,
	}, adminConsole.QueryUserStatisticsHandler(adminConsoleDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createCoupon",
		Method:      http.MethodPost,
		Path:        "/coupon",
		Summary:     "Create coupon",
		Tags:        []string{"coupon"},
		Security:    bearerSecurity,
	}, adminCoupon.CreateCouponHandler(adminCouponDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateCoupon",
		Method:      http.MethodPut,
		Path:        "/coupon",
		Summary:     "Update coupon",
		Tags:        []string{"coupon"},
		Security:    bearerSecurity,
	}, adminCoupon.UpdateCouponHandler(adminCouponDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "deleteCoupon",
		Method:      http.MethodDelete,
		Path:        "/coupon",
		Summary:     "Delete coupon",
		Tags:        []string{"coupon"},
		Security:    bearerSecurity,
	}, adminCoupon.DeleteCouponHandler(adminCouponDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "batchDeleteCoupon",
		Method:      http.MethodDelete,
		Path:        "/coupon/batch",
		Summary:     "Batch delete coupon",
		Tags:        []string{"coupon"},
		Security:    bearerSecurity,
	}, adminCoupon.BatchDeleteCouponHandler(adminCouponDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getCouponList",
		Method:      http.MethodGet,
		Path:        "/coupon/list",
		Summary:     "Get coupon list",
		Tags:        []string{"coupon"},
		Security:    bearerSecurity,
	}, adminCoupon.GetCouponListHandler(adminCouponDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createDocument",
		Method:      http.MethodPost,
		Path:        "/document",
		Summary:     "Create document",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, adminDocument.CreateDocumentHandler(adminDocumentDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateDocument",
		Method:      http.MethodPut,
		Path:        "/document",
		Summary:     "Update document",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, adminDocument.UpdateDocumentHandler(adminDocumentDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "deleteDocument",
		Method:      http.MethodDelete,
		Path:        "/document",
		Summary:     "Delete document",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, adminDocument.DeleteDocumentHandler(adminDocumentDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "batchDeleteDocument",
		Method:      http.MethodDelete,
		Path:        "/document/batch",
		Summary:     "Batch delete document",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, adminDocument.BatchDeleteDocumentHandler(adminDocumentDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getDocumentDetail",
		Method:      http.MethodGet,
		Path:        "/document/detail",
		Summary:     "Get document detail",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, adminDocument.GetDocumentDetailHandler(adminDocumentDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getDocumentList",
		Method:      http.MethodGet,
		Path:        "/document/list",
		Summary:     "Get document list",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, adminDocument.GetDocumentListHandler(adminDocumentDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "filterBalanceLog",
		Method:      http.MethodGet,
		Path:        "/log/balance/list",
		Summary:     "Filter balance log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterBalanceLogHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "filterCommissionLog",
		Method:      http.MethodGet,
		Path:        "/log/commission/list",
		Summary:     "Filter commission log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterCommissionLogHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "filterEmailLog",
		Method:      http.MethodGet,
		Path:        "/log/email/list",
		Summary:     "Filter email log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterEmailLogHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "filterGiftLog",
		Method:      http.MethodGet,
		Path:        "/log/gift/list",
		Summary:     "Filter gift log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterGiftLogHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "filterLoginLog",
		Method:      http.MethodGet,
		Path:        "/log/login/list",
		Summary:     "Filter login log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterLoginLogHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getMessageLogList",
		Method:      http.MethodGet,
		Path:        "/log/message/list",
		Summary:     "Get message log list",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.GetMessageLogListHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "filterMobileLog",
		Method:      http.MethodGet,
		Path:        "/log/mobile/list",
		Summary:     "Filter mobile log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterMobileLogHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "filterRegisterLog",
		Method:      http.MethodGet,
		Path:        "/log/register/list",
		Summary:     "Filter register log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterRegisterLogHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "filterServerTrafficLog",
		Method:      http.MethodGet,
		Path:        "/log/server/traffic/list",
		Summary:     "Filter server traffic log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterServerTrafficLogHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getLogSetting",
		Method:      http.MethodGet,
		Path:        "/log/setting",
		Summary:     "Get log setting",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.GetLogSettingHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateLogSetting",
		Method:      http.MethodPost,
		Path:        "/log/setting",
		Summary:     "Update log setting",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.UpdateLogSettingHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "filterSubscribeLog",
		Method:      http.MethodGet,
		Path:        "/log/subscribe/list",
		Summary:     "Filter subscribe log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterSubscribeLogHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "filterResetSubscribeLog",
		Method:      http.MethodGet,
		Path:        "/log/subscribe/reset/list",
		Summary:     "Filter reset subscribe log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterResetSubscribeLogHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "filterUserSubscribeTrafficLog",
		Method:      http.MethodGet,
		Path:        "/log/subscribe/traffic/list",
		Summary:     "Filter user subscribe traffic log",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterUserSubscribeTrafficLogHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "filterTrafficLogDetails",
		Method:      http.MethodGet,
		Path:        "/log/traffic/details",
		Summary:     "Filter traffic log details",
		Tags:        []string{"log"},
		Security:    bearerSecurity,
	}, adminLog.FilterTrafficLogDetailsHandler(adminLogDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getBatchSendEmailTaskList",
		Method:      http.MethodGet,
		Path:        "/marketing/email/batch/list",
		Summary:     "Get batch send email task list",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.GetBatchSendEmailTaskListHandler(adminMarketingDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getPreSendEmailCount",
		Method:      http.MethodPost,
		Path:        "/marketing/email/batch/pre-send-count",
		Summary:     "Get pre-send email count",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.GetPreSendEmailCountHandler(adminMarketingDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createBatchSendEmailTask",
		Method:      http.MethodPost,
		Path:        "/marketing/email/batch/send",
		Summary:     "Create a batch send email task",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.CreateBatchSendEmailTaskHandler(adminMarketingDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getBatchSendEmailTaskStatus",
		Method:      http.MethodPost,
		Path:        "/marketing/email/batch/status",
		Summary:     "Get batch send email task status",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.GetBatchSendEmailTaskStatusHandler(adminMarketingDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "stopBatchSendEmailTask",
		Method:      http.MethodPost,
		Path:        "/marketing/email/batch/stop",
		Summary:     "Stop a batch send email task",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.StopBatchSendEmailTaskHandler(adminMarketingDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createQuotaTask",
		Method:      http.MethodPost,
		Path:        "/marketing/quota/create",
		Summary:     "Create a quota task",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.CreateQuotaTaskHandler(adminMarketingDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "queryQuotaTaskList",
		Method:      http.MethodGet,
		Path:        "/marketing/quota/list",
		Summary:     "Query quota task list",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.QueryQuotaTaskListHandler(adminMarketingDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "queryQuotaTaskPreCount",
		Method:      http.MethodPost,
		Path:        "/marketing/quota/pre-count",
		Summary:     "Query quota task pre-count",
		Tags:        []string{"marketing"},
		Security:    bearerSecurity,
	}, adminMarketing.QueryQuotaTaskPreCountHandler(adminMarketingDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createOrder",
		Method:      http.MethodPost,
		Path:        "/order",
		Summary:     "Create order",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, adminOrder.CreateOrderHandler(adminOrderDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getOrderList",
		Method:      http.MethodGet,
		Path:        "/order/list",
		Summary:     "Get order list",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, adminOrder.GetOrderListHandler(adminOrderDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateOrderStatus",
		Method:      http.MethodPut,
		Path:        "/order/status",
		Summary:     "Update order status",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, adminOrder.UpdateOrderStatusHandler(adminOrderDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createPaymentMethod",
		Method:      http.MethodPost,
		Path:        "/payment",
		Summary:     "Create Payment Method",
		Tags:        []string{"payment"},
		Security:    bearerSecurity,
	}, adminPayment.CreatePaymentMethodHandler(adminPaymentDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updatePaymentMethod",
		Method:      http.MethodPut,
		Path:        "/payment",
		Summary:     "Update Payment Method",
		Tags:        []string{"payment"},
		Security:    bearerSecurity,
	}, adminPayment.UpdatePaymentMethodHandler(adminPaymentDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "deletePaymentMethod",
		Method:      http.MethodDelete,
		Path:        "/payment",
		Summary:     "Delete Payment Method",
		Tags:        []string{"payment"},
		Security:    bearerSecurity,
	}, adminPayment.DeletePaymentMethodHandler(adminPaymentDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getPaymentMethodList",
		Method:      http.MethodGet,
		Path:        "/payment/list",
		Summary:     "Get Payment Method List",
		Tags:        []string{"payment"},
		Security:    bearerSecurity,
	}, adminPayment.GetPaymentMethodListHandler(adminPaymentDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getPaymentPlatform",
		Method:      http.MethodGet,
		Path:        "/payment/platform",
		Summary:     "Get supported payment platform",
		Tags:        []string{"payment"},
		Security:    bearerSecurity,
	}, adminPayment.GetPaymentPlatformHandler(adminPaymentDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createServer",
		Method:      http.MethodPost,
		Path:        "/server/create",
		Summary:     "Create Server",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.CreateServerHandler(adminServerDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "deleteServer",
		Method:      http.MethodPost,
		Path:        "/server/delete",
		Summary:     "Delete Server",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.DeleteServerHandler(adminServerDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "filterServerList",
		Method:      http.MethodGet,
		Path:        "/server/list",
		Summary:     "Filter Server List",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.FilterServerListHandler(adminServerDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createNode",
		Method:      http.MethodPost,
		Path:        "/server/node/create",
		Summary:     "Create Node",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.CreateNodeHandler(adminServerDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "deleteNode",
		Method:      http.MethodPost,
		Path:        "/server/node/delete",
		Summary:     "Delete Node",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.DeleteNodeHandler(adminServerDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "filterNodeList",
		Method:      http.MethodGet,
		Path:        "/server/node/list",
		Summary:     "Filter Node List",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.FilterNodeListHandler(adminServerDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "resetSortWithNode",
		Method:      http.MethodPost,
		Path:        "/server/node/sort",
		Summary:     "Reset node sort",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.ResetSortWithNodeHandler(adminServerDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "toggleNodeStatus",
		Method:      http.MethodPost,
		Path:        "/server/node/status/toggle",
		Summary:     "Toggle Node Status",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.ToggleNodeStatusHandler(adminServerDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "queryNodeTag",
		Method:      http.MethodGet,
		Path:        "/server/node/tags",
		Summary:     "Query all node tags",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.QueryNodeTagHandler(adminServerDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateNode",
		Method:      http.MethodPost,
		Path:        "/server/node/update",
		Summary:     "Update Node",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.UpdateNodeHandler(adminServerDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getServerProtocols",
		Method:      http.MethodGet,
		Path:        "/server/protocols",
		Summary:     "Get Server Protocols",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.GetServerProtocolsHandler(adminServerDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "resetSortWithServer",
		Method:      http.MethodPost,
		Path:        "/server/server/sort",
		Summary:     "Reset server sort",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.ResetSortWithServerHandler(adminServerDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateServer",
		Method:      http.MethodPost,
		Path:        "/server/update",
		Summary:     "Update Server",
		Tags:        []string{"server"},
		Security:    bearerSecurity,
	}, adminServer.UpdateServerHandler(adminServerDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createSubscribe",
		Method:      http.MethodPost,
		Path:        "/subscribe",
		Summary:     "Create subscribe",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.CreateSubscribeHandler(adminSubscribeDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateSubscribe",
		Method:      http.MethodPut,
		Path:        "/subscribe",
		Summary:     "Update subscribe",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.UpdateSubscribeHandler(adminSubscribeDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "deleteSubscribe",
		Method:      http.MethodDelete,
		Path:        "/subscribe",
		Summary:     "Delete subscribe",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.DeleteSubscribeHandler(adminSubscribeDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "batchDeleteSubscribe",
		Method:      http.MethodDelete,
		Path:        "/subscribe/batch",
		Summary:     "Batch delete subscribe",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.BatchDeleteSubscribeHandler(adminSubscribeDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getSubscribeDetails",
		Method:      http.MethodGet,
		Path:        "/subscribe/details",
		Summary:     "Get subscribe details",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.GetSubscribeDetailsHandler(adminSubscribeDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createSubscribeGroup",
		Method:      http.MethodPost,
		Path:        "/subscribe/group",
		Summary:     "Create subscribe group",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.CreateSubscribeGroupHandler(adminSubscribeDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateSubscribeGroup",
		Method:      http.MethodPut,
		Path:        "/subscribe/group",
		Summary:     "Update subscribe group",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.UpdateSubscribeGroupHandler(adminSubscribeDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "deleteSubscribeGroup",
		Method:      http.MethodDelete,
		Path:        "/subscribe/group",
		Summary:     "Delete subscribe group",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.DeleteSubscribeGroupHandler(adminSubscribeDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "batchDeleteSubscribeGroup",
		Method:      http.MethodDelete,
		Path:        "/subscribe/group/batch",
		Summary:     "Batch delete subscribe group",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.BatchDeleteSubscribeGroupHandler(adminSubscribeDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getSubscribeGroupList",
		Method:      http.MethodGet,
		Path:        "/subscribe/group/list",
		Summary:     "Get subscribe group list",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.GetSubscribeGroupListHandler(adminSubscribeDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getSubscribeList",
		Method:      http.MethodGet,
		Path:        "/subscribe/list",
		Summary:     "Get subscribe list",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.GetSubscribeListHandler(adminSubscribeDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "resetAllSubscribeToken",
		Method:      http.MethodPost,
		Path:        "/subscribe/reset_all_token",
		Summary:     "Reset all subscribe tokens",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.ResetAllSubscribeTokenHandler(adminSubscribeDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "subscribeSort",
		Method:      http.MethodPost,
		Path:        "/subscribe/sort",
		Summary:     "Subscribe sort",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, adminSubscribe.SubscribeSortHandler(adminSubscribeDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getCurrencyConfig",
		Method:      http.MethodGet,
		Path:        "/system/currency_config",
		Summary:     "Get Currency Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetCurrencyConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateCurrencyConfig",
		Method:      http.MethodPut,
		Path:        "/system/currency_config",
		Summary:     "Update Currency Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateCurrencyConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getNodeMultiplier",
		Method:      http.MethodGet,
		Path:        "/system/get_node_multiplier",
		Summary:     "Get Node Multiplier",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetNodeMultiplierHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getInviteConfig",
		Method:      http.MethodGet,
		Path:        "/system/invite_config",
		Summary:     "Get invite config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetInviteConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateInviteConfig",
		Method:      http.MethodPut,
		Path:        "/system/invite_config",
		Summary:     "Update invite config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateInviteConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getModuleConfig",
		Method:      http.MethodGet,
		Path:        "/system/module",
		Summary:     "Get Module Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetModuleConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getNodeConfig",
		Method:      http.MethodGet,
		Path:        "/system/node_config",
		Summary:     "Get node config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetNodeConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateNodeConfig",
		Method:      http.MethodPut,
		Path:        "/system/node_config",
		Summary:     "Update node config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateNodeConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "preViewNodeMultiplier",
		Method:      http.MethodGet,
		Path:        "/system/node_multiplier/preview",
		Summary:     "PreView Node Multiplier",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.PreViewNodeMultiplierHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getPrivacyPolicyConfig",
		Method:      http.MethodGet,
		Path:        "/system/privacy",
		Summary:     "get Privacy Policy Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetPrivacyPolicyConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updatePrivacyPolicyConfig",
		Method:      http.MethodPut,
		Path:        "/system/privacy",
		Summary:     "Update Privacy Policy Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdatePrivacyPolicyConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getRegisterConfig",
		Method:      http.MethodGet,
		Path:        "/system/register_config",
		Summary:     "Get register config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetRegisterConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateRegisterConfig",
		Method:      http.MethodPut,
		Path:        "/system/register_config",
		Summary:     "Update register config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateRegisterConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "setNodeMultiplier",
		Method:      http.MethodPost,
		Path:        "/system/set_node_multiplier",
		Summary:     "Set Node Multiplier",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.SetNodeMultiplierHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "settingTelegramBot",
		Method:      http.MethodPost,
		Path:        "/system/setting_telegram_bot",
		Summary:     "setting telegram bot",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.SettingTelegramBotHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getSiteConfig",
		Method:      http.MethodGet,
		Path:        "/system/site_config",
		Summary:     "Get site config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetSiteConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateSiteConfig",
		Method:      http.MethodPut,
		Path:        "/system/site_config",
		Summary:     "Update site config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateSiteConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getSubscribeConfig",
		Method:      http.MethodGet,
		Path:        "/system/subscribe_config",
		Summary:     "Get subscribe config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetSubscribeConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateSubscribeConfig",
		Method:      http.MethodPut,
		Path:        "/system/subscribe_config",
		Summary:     "Update subscribe config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateSubscribeConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getTosConfig",
		Method:      http.MethodGet,
		Path:        "/system/tos_config",
		Summary:     "Get Team of Service Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetTosConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateTosConfig",
		Method:      http.MethodPut,
		Path:        "/system/tos_config",
		Summary:     "Update Team of Service Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateTosConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getVerifyCodeConfig",
		Method:      http.MethodGet,
		Path:        "/system/verify_code_config",
		Summary:     "Get Verify Code Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetVerifyCodeConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateVerifyCodeConfig",
		Method:      http.MethodPut,
		Path:        "/system/verify_code_config",
		Summary:     "Update Verify Code Config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateVerifyCodeConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getVerifyConfig",
		Method:      http.MethodGet,
		Path:        "/system/verify_config",
		Summary:     "Get verify config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.GetVerifyConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateVerifyConfig",
		Method:      http.MethodPut,
		Path:        "/system/verify_config",
		Summary:     "Update verify config",
		Tags:        []string{"system"},
		Security:    bearerSecurity,
	}, adminSystem.UpdateVerifyConfigHandler(adminSystemDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateTicketStatus",
		Method:      http.MethodPut,
		Path:        "/ticket",
		Summary:     "Update ticket status",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, adminTicket.UpdateTicketStatusHandler(adminTicketDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getTicket",
		Method:      http.MethodGet,
		Path:        "/ticket/detail",
		Summary:     "Get ticket detail",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, adminTicket.GetTicketHandler(adminTicketDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createTicketFollow",
		Method:      http.MethodPost,
		Path:        "/ticket/follow",
		Summary:     "Create ticket follow",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, adminTicket.CreateTicketFollowHandler(adminTicketDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getTicketList",
		Method:      http.MethodGet,
		Path:        "/ticket/list",
		Summary:     "Get ticket list",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, adminTicket.GetTicketListHandler(adminTicketDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "queryIPLocation",
		Method:      http.MethodGet,
		Path:        "/tool/ip/location",
		Summary:     "Query IP Location",
		Tags:        []string{"tool"},
		Security:    bearerSecurity,
	}, adminTool.QueryIPLocationHandler(adminToolDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getSystemLog",
		Method:      http.MethodGet,
		Path:        "/tool/log",
		Summary:     "Get System Log",
		Tags:        []string{"tool"},
		Security:    bearerSecurity,
	}, adminTool.GetSystemLogHandler(adminToolDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "restartSystem",
		Method:      http.MethodGet,
		Path:        "/tool/restart",
		Summary:     "Restart System",
		Tags:        []string{"tool"},
		Security:    bearerSecurity,
	}, adminTool.RestartSystemHandler(adminToolDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getVersion",
		Method:      http.MethodGet,
		Path:        "/tool/version",
		Summary:     "Get Version",
		Tags:        []string{"tool"},
		Security:    bearerSecurity,
	}, adminTool.GetVersionHandler(adminToolDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "deleteUser",
		Method:      http.MethodDelete,
		Path:        "/user",
		Summary:     "Delete user",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.DeleteUserHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createUser",
		Method:      http.MethodPost,
		Path:        "/user",
		Summary:     "Create user",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.CreateUserHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createUserAuthMethod",
		Method:      http.MethodPost,
		Path:        "/user/auth_method",
		Summary:     "Create user auth method",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.CreateUserAuthMethodHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "deleteUserAuthMethod",
		Method:      http.MethodDelete,
		Path:        "/user/auth_method",
		Summary:     "Delete user auth method",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.DeleteUserAuthMethodHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateUserAuthMethod",
		Method:      http.MethodPut,
		Path:        "/user/auth_method",
		Summary:     "Update user auth method",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.UpdateUserAuthMethodHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getUserAuthMethod",
		Method:      http.MethodGet,
		Path:        "/user/auth_method",
		Summary:     "Get user auth method",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserAuthMethodHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateUserBasicInfo",
		Method:      http.MethodPut,
		Path:        "/user/basic",
		Summary:     "Update user basic info",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.UpdateUserBasicInfoHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "batchDeleteUser",
		Method:      http.MethodDelete,
		Path:        "/user/batch",
		Summary:     "Batch delete user",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.BatchDeleteUserHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "currentUser",
		Method:      http.MethodGet,
		Path:        "/user/current",
		Summary:     "Current user",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.CurrentUserHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getUserDetail",
		Method:      http.MethodGet,
		Path:        "/user/detail",
		Summary:     "Get user detail",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserDetailHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateUserDevice",
		Method:      http.MethodPut,
		Path:        "/user/device",
		Summary:     "User device",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.UpdateUserDeviceHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "deleteUserDevice",
		Method:      http.MethodDelete,
		Path:        "/user/device",
		Summary:     "Delete user device",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.DeleteUserDeviceHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "kickOfflineByUserDevice",
		Method:      http.MethodPut,
		Path:        "/user/device/kick_offline",
		Summary:     "kick offline user device",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.KickOfflineByUserDeviceHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getUserList",
		Method:      http.MethodGet,
		Path:        "/user/list",
		Summary:     "Get user list",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserListHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getUserLoginLogs",
		Method:      http.MethodGet,
		Path:        "/user/login/logs",
		Summary:     "Get user login logs",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserLoginLogsHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateUserNotifySetting",
		Method:      http.MethodPut,
		Path:        "/user/notify",
		Summary:     "Update user notify setting",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.UpdateUserNotifySettingHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getUserSubscribe",
		Method:      http.MethodGet,
		Path:        "/user/subscribe",
		Summary:     "Get user subcribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserSubscribeHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "createUserSubscribe",
		Method:      http.MethodPost,
		Path:        "/user/subscribe",
		Summary:     "Create user subcribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.CreateUserSubscribeHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "updateUserSubscribe",
		Method:      http.MethodPut,
		Path:        "/user/subscribe",
		Summary:     "Update user subcribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.UpdateUserSubscribeHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "deleteUserSubscribe",
		Method:      http.MethodDelete,
		Path:        "/user/subscribe",
		Summary:     "Delete user subcribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.DeleteUserSubscribeHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getUserSubscribeById",
		Method:      http.MethodGet,
		Path:        "/user/subscribe/detail",
		Summary:     "Get user subcribe by id",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserSubscribeByIdHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getUserSubscribeDevices",
		Method:      http.MethodGet,
		Path:        "/user/subscribe/device",
		Summary:     "Get user subcribe devices",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserSubscribeDevicesHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getUserSubscribeLogs",
		Method:      http.MethodGet,
		Path:        "/user/subscribe/logs",
		Summary:     "Get user subcribe logs",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserSubscribeLogsHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getUserSubscribeResetTrafficLogs",
		Method:      http.MethodGet,
		Path:        "/user/subscribe/reset/logs",
		Summary:     "Get user subcribe reset traffic logs",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserSubscribeResetTrafficLogsHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "resetUserSubscribeToken",
		Method:      http.MethodPost,
		Path:        "/user/subscribe/reset/token",
		Summary:     "Reset user subscribe token",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.ResetUserSubscribeTokenHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "resetUserSubscribeTraffic",
		Method:      http.MethodPost,
		Path:        "/user/subscribe/reset/traffic",
		Summary:     "Reset user subscribe traffic",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.ResetUserSubscribeTrafficHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "toggleUserSubscribeStatus",
		Method:      http.MethodPost,
		Path:        "/user/subscribe/toggle",
		Summary:     "Stop user subscribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.ToggleUserSubscribeStatusHandler(adminUserDeps))

	registerOperation(apis.Admin, huma.Operation{
		OperationID: "getUserSubscribeTrafficLogs",
		Method:      http.MethodGet,
		Path:        "/user/subscribe/traffic_logs",
		Summary:     "Get user subcribe traffic logs",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, adminUser.GetUserSubscribeTrafficLogsHandler(adminUserDeps))

}
