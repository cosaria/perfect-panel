package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	appruntime "github.com/perfect-panel/server/internal/bootstrap/runtime"
	"github.com/perfect-panel/server/internal/platform/http/middleware"
	publicAnnouncement "github.com/perfect-panel/server/services/user/announcement"
	publicDocument "github.com/perfect-panel/server/services/user/document"
	publicOrder "github.com/perfect-panel/server/services/user/order"
	publicPayment "github.com/perfect-panel/server/services/user/payment"
	publicPortal "github.com/perfect-panel/server/services/user/portal"
	publicSubscribe "github.com/perfect-panel/server/services/user/subscribe"
	publicTicket "github.com/perfect-panel/server/services/user/ticket"
	publicUser "github.com/perfect-panel/server/services/user/user"
)

func registerPublicRoutes(router *gin.Engine, runtimeDeps *appruntime.Deps, specOnly bool) []huma.API {
	publicAnnouncementGroup := router.Group("/api/v1/public/announcement")
	if !specOnly {
		publicAnnouncementGroup.Use(middleware.AuthMiddleware(runtimeDeps))
	}
	publicAnnouncementConfig := governedAPIConfig("Public Announcement API", "1.0.0", "/api/v1/public/announcement", "announcement")
	publicAnnouncementAPI := humagin.NewWithGroup(router, publicAnnouncementGroup, publicAnnouncementConfig)
	configureHumaAPI(publicAnnouncementAPI, compatibilityEnabled(runtimeDeps, specOnly))
	publicAnnouncementDeps := publicAnnouncement.Deps{}
	publicDocumentDeps := publicDocument.Deps{}
	publicOrderDeps := publicOrder.Deps{}
	publicPaymentDeps := publicPayment.Deps{}
	publicPortalDeps := newPublicPortalDeps(runtimeDeps)
	publicSubscribeDeps := publicSubscribe.Deps{}
	publicTicketDeps := publicTicket.Deps{}
	publicUserDeps := newPublicUserDeps(runtimeDeps)
	if runtimeDeps != nil {
		publicAnnouncementDeps.AnnouncementModel = runtimeDeps.AnnouncementModel
		publicDocumentDeps.DocumentModel = runtimeDeps.DocumentModel
		publicOrderDeps.OrderModel = runtimeDeps.OrderModel
		publicOrderDeps.PaymentModel = runtimeDeps.PaymentModel
		publicOrderDeps.SubscribeModel = runtimeDeps.SubscribeModel
		publicOrderDeps.UserModel = runtimeDeps.UserModel
		publicOrderDeps.CouponModel = runtimeDeps.CouponModel
		publicOrderDeps.DB = runtimeDeps.DB
		publicOrderDeps.Queue = runtimeDeps.Queue
		publicOrderDeps.Config = runtimeDeps.Config
		publicPaymentDeps.PaymentModel = runtimeDeps.PaymentModel
		publicSubscribeDeps.SubscribeModel = runtimeDeps.SubscribeModel
		publicSubscribeDeps.UserModel = runtimeDeps.UserModel
		publicSubscribeDeps.NodeModel = runtimeDeps.NodeModel
		publicSubscribeDeps.DB = runtimeDeps.DB
		publicSubscribeDeps.Config = runtimeDeps.Config
		publicTicketDeps.TicketModel = runtimeDeps.TicketModel
	}

	registerOperation(publicAnnouncementAPI, huma.Operation{
		OperationID: "queryAnnouncement",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "Query announcement",
		Tags:        []string{"announcement"},
		Security:    bearerSecurity,
	}, publicAnnouncement.QueryAnnouncementHandler(publicAnnouncementDeps))

	publicDocumentGroup := router.Group("/api/v1/public/document")
	if !specOnly {
		publicDocumentGroup.Use(middleware.AuthMiddleware(runtimeDeps))
	}
	publicDocumentConfig := governedAPIConfig("Public Document API", "1.0.0", "/api/v1/public/document", "document")
	publicDocumentAPI := humagin.NewWithGroup(router, publicDocumentGroup, publicDocumentConfig)
	configureHumaAPI(publicDocumentAPI, compatibilityEnabled(runtimeDeps, specOnly))

	registerOperation(publicDocumentAPI, huma.Operation{
		OperationID: "queryDocumentDetail",
		Method:      http.MethodGet,
		Path:        "/detail",
		Summary:     "Get document detail",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, publicDocument.QueryDocumentDetailHandler(publicDocumentDeps))

	registerOperation(publicDocumentAPI, huma.Operation{
		OperationID: "queryDocumentList",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "Get document list",
		Tags:        []string{"document"},
		Security:    bearerSecurity,
	}, publicDocument.QueryDocumentListHandler(publicDocumentDeps))

	publicOrderGroup := router.Group("/api/v1/public/order")
	if !specOnly {
		publicOrderGroup.Use(middleware.AuthMiddleware(runtimeDeps))
	}
	publicOrderConfig := governedAPIConfig("Public Order API", "1.0.0", "/api/v1/public/order", "order")
	publicOrderAPI := humagin.NewWithGroup(router, publicOrderGroup, publicOrderConfig)
	configureHumaAPI(publicOrderAPI, compatibilityEnabled(runtimeDeps, specOnly))

	registerOperation(publicOrderAPI, huma.Operation{
		OperationID: "closeOrder",
		Method:      http.MethodPost,
		Path:        "/close",
		Summary:     "Close order",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.CloseOrderHandler(publicOrderDeps))

	registerOperation(publicOrderAPI, huma.Operation{
		OperationID: "queryOrderDetail",
		Method:      http.MethodGet,
		Path:        "/detail",
		Summary:     "Get order",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.QueryOrderDetailHandler(publicOrderDeps))

	registerOperation(publicOrderAPI, huma.Operation{
		OperationID: "queryOrderList",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "Get order list",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.QueryOrderListHandler(publicOrderDeps))

	registerOperation(publicOrderAPI, huma.Operation{
		OperationID: "preCreateOrder",
		Method:      http.MethodPost,
		Path:        "/pre",
		Summary:     "Pre create order",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.PreCreateOrderHandler(publicOrderDeps))

	registerOperation(publicOrderAPI, huma.Operation{
		OperationID: "purchase",
		Method:      http.MethodPost,
		Path:        "/purchase",
		Summary:     "purchase Subscription",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.PurchaseHandler(publicOrderDeps))

	registerOperation(publicOrderAPI, huma.Operation{
		OperationID: "recharge",
		Method:      http.MethodPost,
		Path:        "/recharge",
		Summary:     "Recharge",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.RechargeHandler(publicOrderDeps))

	registerOperation(publicOrderAPI, huma.Operation{
		OperationID: "renewal",
		Method:      http.MethodPost,
		Path:        "/renewal",
		Summary:     "Renewal Subscription",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.RenewalHandler(publicOrderDeps))

	registerOperation(publicOrderAPI, huma.Operation{
		OperationID: "resetTraffic",
		Method:      http.MethodPost,
		Path:        "/reset",
		Summary:     "Reset traffic",
		Tags:        []string{"order"},
		Security:    bearerSecurity,
	}, publicOrder.ResetTrafficHandler(publicOrderDeps))

	publicPaymentGroup := router.Group("/api/v1/public/payment")
	if !specOnly {
		publicPaymentGroup.Use(middleware.AuthMiddleware(runtimeDeps))
	}
	publicPaymentConfig := governedAPIConfig("Public Payment API", "1.0.0", "/api/v1/public/payment", "payment")
	publicPaymentAPI := humagin.NewWithGroup(router, publicPaymentGroup, publicPaymentConfig)
	configureHumaAPI(publicPaymentAPI, compatibilityEnabled(runtimeDeps, specOnly))

	registerOperation(publicPaymentAPI, huma.Operation{
		OperationID: "getAvailablePaymentMethods",
		Method:      http.MethodGet,
		Path:        "/methods",
		Summary:     "Get available payment methods",
		Tags:        []string{"payment"},
		Security:    bearerSecurity,
	}, publicPayment.GetAvailablePaymentMethodsHandler(publicPaymentDeps))

	publicSubscribeGroup := router.Group("/api/v1/public/subscribe")
	if !specOnly {
		publicSubscribeGroup.Use(middleware.AuthMiddleware(runtimeDeps))
	}
	publicSubscribeConfig := governedAPIConfig("Public Subscribe API", "1.0.0", "/api/v1/public/subscribe", "subscribe")
	publicSubscribeAPI := humagin.NewWithGroup(router, publicSubscribeGroup, publicSubscribeConfig)
	configureHumaAPI(publicSubscribeAPI, compatibilityEnabled(runtimeDeps, specOnly))

	registerOperation(publicSubscribeAPI, huma.Operation{
		OperationID: "querySubscribeList",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "Get subscribe list",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, publicSubscribe.QuerySubscribeListHandler(publicSubscribeDeps))

	registerOperation(publicSubscribeAPI, huma.Operation{
		OperationID: "queryUserSubscribeNodeList",
		Method:      http.MethodGet,
		Path:        "/node/list",
		Summary:     "Get user subscribe node info",
		Tags:        []string{"subscribe"},
		Security:    bearerSecurity,
	}, publicSubscribe.QueryUserSubscribeNodeListHandler(publicSubscribeDeps))

	publicTicketGroup := router.Group("/api/v1/public/ticket")
	if !specOnly {
		publicTicketGroup.Use(middleware.AuthMiddleware(runtimeDeps))
	}
	publicTicketConfig := governedAPIConfig("Public Ticket API", "1.0.0", "/api/v1/public/ticket", "ticket")
	publicTicketAPI := humagin.NewWithGroup(router, publicTicketGroup, publicTicketConfig)
	configureHumaAPI(publicTicketAPI, compatibilityEnabled(runtimeDeps, specOnly))

	registerOperation(publicTicketAPI, huma.Operation{
		OperationID: "updateUserTicketStatus",
		Method:      http.MethodPut,
		Path:        "/ticket",
		Summary:     "Update ticket status",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, publicTicket.UpdateUserTicketStatusHandler(publicTicketDeps))

	registerOperation(publicTicketAPI, huma.Operation{
		OperationID: "createUserTicket",
		Method:      http.MethodPost,
		Path:        "/ticket",
		Summary:     "Create ticket",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, publicTicket.CreateUserTicketHandler(publicTicketDeps))

	registerOperation(publicTicketAPI, huma.Operation{
		OperationID: "getUserTicketDetails",
		Method:      http.MethodGet,
		Path:        "/detail",
		Summary:     "Get ticket detail",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, publicTicket.GetUserTicketDetailsHandler(publicTicketDeps))

	registerOperation(publicTicketAPI, huma.Operation{
		OperationID: "createUserTicketFollow",
		Method:      http.MethodPost,
		Path:        "/follow",
		Summary:     "Create ticket follow",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, publicTicket.CreateUserTicketFollowHandler(publicTicketDeps))

	registerOperation(publicTicketAPI, huma.Operation{
		OperationID: "getUserTicketList",
		Method:      http.MethodGet,
		Path:        "/list",
		Summary:     "Get ticket list",
		Tags:        []string{"ticket"},
		Security:    bearerSecurity,
	}, publicTicket.GetUserTicketListHandler(publicTicketDeps))

	publicUserGroup := router.Group("/api/v1/public/user")
	if !specOnly {
		publicUserGroup.Use(middleware.AuthMiddleware(runtimeDeps))
	}
	publicUserConfig := governedAPIConfig("Public User API", "1.0.0", "/api/v1/public/user", "user")
	publicUserAPI := humagin.NewWithGroup(router, publicUserGroup, publicUserConfig)
	configureHumaAPI(publicUserAPI, compatibilityEnabled(runtimeDeps, specOnly))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "queryUserAffiliate",
		Method:      http.MethodGet,
		Path:        "/affiliate/count",
		Summary:     "Query User Affiliate Count",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.QueryUserAffiliateHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "queryUserAffiliateList",
		Method:      http.MethodGet,
		Path:        "/affiliate/list",
		Summary:     "Query User Affiliate List",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.QueryUserAffiliateListHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "queryUserBalanceLog",
		Method:      http.MethodGet,
		Path:        "/balance_log",
		Summary:     "Query User Balance Log",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.QueryUserBalanceLogHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "updateBindEmail",
		Method:      http.MethodPut,
		Path:        "/bind_email",
		Summary:     "Update Bind Email",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UpdateBindEmailHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "updateBindMobile",
		Method:      http.MethodPut,
		Path:        "/bind_mobile",
		Summary:     "Update Bind Mobile",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UpdateBindMobileHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "bindOAuth",
		Method:      http.MethodPost,
		Path:        "/bind_oauth",
		Summary:     "Bind OAuth",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.BindOAuthHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "bindOAuthCallback",
		Method:      http.MethodPost,
		Path:        "/bind_oauth/callback",
		Summary:     "Bind OAuth Callback",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.BindOAuthCallbackHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "bindTelegram",
		Method:      http.MethodGet,
		Path:        "/bind_telegram",
		Summary:     "Bind Telegram",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.BindTelegramHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "queryUserCommissionLog",
		Method:      http.MethodGet,
		Path:        "/commission_log",
		Summary:     "Query User Commission Log",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.QueryUserCommissionLogHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "commissionWithdraw",
		Method:      http.MethodPost,
		Path:        "/commission_withdraw",
		Summary:     "Commission Withdraw",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.CommissionWithdrawHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "getDeviceList",
		Method:      http.MethodGet,
		Path:        "/devices",
		Summary:     "Get Device List",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.GetDeviceListHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "queryUserInfo",
		Method:      http.MethodGet,
		Path:        "/info",
		Summary:     "Query User Info",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.QueryUserInfoHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "getLoginLog",
		Method:      http.MethodGet,
		Path:        "/login_log",
		Summary:     "Get Login Log",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.GetLoginLogHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "updateUserNotify",
		Method:      http.MethodPut,
		Path:        "/notify",
		Summary:     "Update User Notify",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UpdateUserNotifyHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "getOAuthMethods",
		Method:      http.MethodGet,
		Path:        "/oauth_methods",
		Summary:     "Get OAuth Methods",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.GetOAuthMethodsHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "updateUserPassword",
		Method:      http.MethodPut,
		Path:        "/password",
		Summary:     "Update User Password",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UpdateUserPasswordHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "updateUserRules",
		Method:      http.MethodPut,
		Path:        "/rules",
		Summary:     "Update User Rules",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UpdateUserRulesHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "queryUserSubscribe",
		Method:      http.MethodGet,
		Path:        "/subscribe",
		Summary:     "Query User Subscribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.QueryUserSubscribeHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "getSubscribeLog",
		Method:      http.MethodGet,
		Path:        "/subscribe_log",
		Summary:     "Get Subscribe Log",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.GetSubscribeLogHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "updateUserSubscribeNote",
		Method:      http.MethodPut,
		Path:        "/subscribe_note",
		Summary:     "Update User Subscribe Note",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UpdateUserSubscribeNoteHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "resetUserSubscribeToken",
		Method:      http.MethodPut,
		Path:        "/subscribe_token",
		Summary:     "Reset User Subscribe Token",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.ResetUserSubscribeTokenHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "unbindDevice",
		Method:      http.MethodPut,
		Path:        "/unbind_device",
		Summary:     "Unbind Device",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UnbindDeviceHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "unbindOAuth",
		Method:      http.MethodPost,
		Path:        "/unbind_oauth",
		Summary:     "Unbind OAuth",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UnbindOAuthHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "unbindTelegram",
		Method:      http.MethodPost,
		Path:        "/unbind_telegram",
		Summary:     "Unbind Telegram",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UnbindTelegramHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "unsubscribe",
		Method:      http.MethodPost,
		Path:        "/unsubscribe",
		Summary:     "Unsubscribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.UnsubscribeHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "preUnsubscribe",
		Method:      http.MethodPost,
		Path:        "/unsubscribe/pre",
		Summary:     "Pre Unsubscribe",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.PreUnsubscribeHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "verifyEmail",
		Method:      http.MethodPost,
		Path:        "/verify_email",
		Summary:     "Verify Email",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.VerifyEmailHandler(publicUserDeps))

	registerOperation(publicUserAPI, huma.Operation{
		OperationID: "queryWithdrawalLog",
		Method:      http.MethodGet,
		Path:        "/withdrawal_log",
		Summary:     "Query Withdrawal Log",
		Tags:        []string{"user"},
		Security:    bearerSecurity,
	}, publicUser.QueryWithdrawalLogHandler(publicUserDeps))

	portalGroup := router.Group("/api/v1/public/portal")
	if !specOnly {
		portalGroup.Use(middleware.DeviceMiddleware(runtimeDeps))
	}
	portalConfig := governedAPIConfig("Portal API", "1.0.0", "/api/v1/public/portal", "portal")
	portalAPI := humagin.NewWithGroup(router, portalGroup, portalConfig)
	configureHumaAPI(portalAPI, compatibilityEnabled(runtimeDeps, specOnly))

	registerOperation(portalAPI, huma.Operation{
		OperationID: "purchaseCheckout",
		Method:      http.MethodPost,
		Path:        "/order/checkout",
		Summary:     "Purchase Checkout",
		Tags:        []string{"portal"},
	}, publicPortal.PurchaseCheckoutHandler(publicPortalDeps))

	registerOperation(portalAPI, huma.Operation{
		OperationID: "queryPurchaseOrder",
		Method:      http.MethodGet,
		Path:        "/order/status",
		Summary:     "Query Purchase Order",
		Tags:        []string{"portal"},
	}, publicPortal.QueryPurchaseOrderHandler(publicPortalDeps))

	registerOperation(portalAPI, huma.Operation{
		OperationID: "portalGetAvailablePaymentMethods",
		Method:      http.MethodGet,
		Path:        "/payment-method",
		Summary:     "Get available payment methods",
		Tags:        []string{"portal"},
	}, publicPortal.GetAvailablePaymentMethodsHandler(publicPortalDeps))

	registerOperation(portalAPI, huma.Operation{
		OperationID: "prePurchaseOrder",
		Method:      http.MethodPost,
		Path:        "/pre",
		Summary:     "Pre Purchase Order",
		Tags:        []string{"portal"},
	}, publicPortal.PrePurchaseOrderHandler(publicPortalDeps))

	registerOperation(portalAPI, huma.Operation{
		OperationID: "portalPurchase",
		Method:      http.MethodPost,
		Path:        "/purchase",
		Summary:     "Purchase subscription",
		Tags:        []string{"portal"},
	}, publicPortal.PurchaseHandler(publicPortalDeps))

	registerOperation(portalAPI, huma.Operation{
		OperationID: "getSubscription",
		Method:      http.MethodGet,
		Path:        "/subscribe",
		Summary:     "Get Subscription",
		Tags:        []string{"portal"},
	}, publicPortal.GetSubscriptionHandler(publicPortalDeps))

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
