package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/routers/middleware"
	auth "github.com/perfect-panel/server/services/auth"
	authOauth "github.com/perfect-panel/server/services/auth/oauth"
	"github.com/perfect-panel/server/svc"
)

func registerAuthRoutes(router *gin.Engine, serverCtx *svc.ServiceContext, specOnly bool) []huma.API {
	authGroup := router.Group("/api/v1/auth")
	if !specOnly {
		authGroup.Use(middleware.DeviceMiddleware(serverCtx))
	}
	authConfig := governedAPIConfig("Auth API", "1.0.0", "/api/v1/auth", "auth")
	authAPI := humagin.NewWithGroup(router, authGroup, authConfig)
	configureHumaAPI(authAPI, compatibilityEnabled(serverCtx, specOnly))

	registerOperation(authAPI, huma.Operation{
		OperationID: "checkUser",
		Method:      http.MethodGet,
		Path:        "/check",
		Summary:     "Check user is exist",
		Tags:        []string{"auth"},
	}, auth.CheckUserHandler(serverCtx))

	registerOperation(authAPI, huma.Operation{
		OperationID: "checkUserTelephone",
		Method:      http.MethodGet,
		Path:        "/check/telephone",
		Summary:     "Check user telephone is exist",
		Tags:        []string{"auth"},
	}, auth.CheckUserTelephoneHandler(serverCtx))

	registerOperation(authAPI, huma.Operation{
		OperationID: "userLogin",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "User login",
		Tags:        []string{"auth"},
	}, auth.UserLoginHandler(serverCtx))

	registerOperation(authAPI, huma.Operation{
		OperationID: "deviceLogin",
		Method:      http.MethodPost,
		Path:        "/login/device",
		Summary:     "Device Login",
		Tags:        []string{"auth"},
	}, auth.DeviceLoginHandler(serverCtx))

	registerOperation(authAPI, huma.Operation{
		OperationID: "telephoneLogin",
		Method:      http.MethodPost,
		Path:        "/login/telephone",
		Summary:     "User Telephone login",
		Tags:        []string{"auth"},
	}, auth.TelephoneLoginHandler(serverCtx))

	registerOperation(authAPI, huma.Operation{
		OperationID: "userRegister",
		Method:      http.MethodPost,
		Path:        "/register",
		Summary:     "User register",
		Tags:        []string{"auth"},
	}, auth.UserRegisterHandler(serverCtx))

	registerOperation(authAPI, huma.Operation{
		OperationID: "telephoneUserRegister",
		Method:      http.MethodPost,
		Path:        "/register/telephone",
		Summary:     "User Telephone register",
		Tags:        []string{"auth"},
	}, auth.TelephoneUserRegisterHandler(serverCtx))

	registerOperation(authAPI, huma.Operation{
		OperationID: "resetPassword",
		Method:      http.MethodPost,
		Path:        "/reset",
		Summary:     "Reset password",
		Tags:        []string{"auth"},
	}, auth.ResetPasswordHandler(serverCtx))

	registerOperation(authAPI, huma.Operation{
		OperationID: "telephoneResetPassword",
		Method:      http.MethodPost,
		Path:        "/reset/telephone",
		Summary:     "Reset password",
		Tags:        []string{"auth"},
	}, auth.TelephoneResetPasswordHandler(serverCtx))

	authOauthGroup := router.Group("/api/v1/auth/oauth")
	authOauthConfig := governedAPIConfig("Auth OAuth API", "1.0.0", "/api/v1/auth/oauth", "oauth")
	authOauthAPI := humagin.NewWithGroup(router, authOauthGroup, authOauthConfig)
	configureHumaAPI(authOauthAPI, compatibilityEnabled(serverCtx, specOnly))

	// AppleLoginCallback stays raw Gin because it needs direct redirect primitives.
	authOauthGroup.POST("/callback/apple", authOauth.AppleLoginCallbackHandler(serverCtx))

	registerOperation(authOauthAPI, huma.Operation{
		OperationID: "oAuthLogin",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "OAuth login",
		Tags:        []string{"oauth"},
	}, authOauth.OAuthLoginHandler(serverCtx))

	registerOperation(authOauthAPI, huma.Operation{
		OperationID: "oAuthLoginGetToken",
		Method:      http.MethodPost,
		Path:        "/login/token",
		Summary:     "OAuth login get token",
		Tags:        []string{"oauth"},
	}, authOauth.OAuthLoginGetTokenHandler(serverCtx))

	return []huma.API{authAPI, authOauthAPI}
}
