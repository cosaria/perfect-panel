package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/routers/middleware"
	publicAnnouncement "github.com/perfect-panel/server/services/user/announcement"
	publicDocument "github.com/perfect-panel/server/services/user/document"
	publicOrder "github.com/perfect-panel/server/services/user/order"
	publicPayment "github.com/perfect-panel/server/services/user/payment"
	publicPortal "github.com/perfect-panel/server/services/user/portal"
	publicSubscribe "github.com/perfect-panel/server/services/user/subscribe"
	publicTicket "github.com/perfect-panel/server/services/user/ticket"
	publicUser "github.com/perfect-panel/server/services/user/user"
	"github.com/perfect-panel/server/svc"
)

func registerPublicRoutes(router *gin.Engine, serverCtx *svc.ServiceContext, specOnly bool) []huma.API {
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

	return []huma.API{
		publicAnnouncementAPI,
		publicDocumentAPI,
		publicOrderAPI,
		publicPaymentAPI,
		publicSubscribeAPI,
		publicTicketAPI,
		publicUserAPI,
		portalAPI,
	}
}
