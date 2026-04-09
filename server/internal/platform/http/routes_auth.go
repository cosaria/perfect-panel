package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	appruntime "github.com/perfect-panel/server/internal/bootstrap/runtime"
	"github.com/perfect-panel/server/internal/platform/http/middleware"
	auth "github.com/perfect-panel/server/services/auth"
	authOauth "github.com/perfect-panel/server/services/auth/oauth"
)

func registerAuthRoutes(router *gin.Engine, runtimeDeps *appruntime.Deps, specOnly bool) []huma.API {
	authGroup := router.Group("/api/v1/auth")
	if !specOnly {
		authGroup.Use(middleware.DeviceMiddleware(runtimeDeps))
	}
	authConfig := governedAPIConfig("Auth API", "1.0.0", "/api/v1/auth", "auth")
	authAPI := humagin.NewWithGroup(router, authGroup, authConfig)
	configureHumaAPI(authAPI, compatibilityEnabled(runtimeDeps, specOnly))
	authDeps := auth.Deps{}
	if runtimeDeps != nil {
		authDeps.UserModel = runtimeDeps.UserModel
		authDeps.LogModel = runtimeDeps.LogModel
		authDeps.SubscribeModel = runtimeDeps.SubscribeModel
		authDeps.Redis = runtimeDeps.Redis
		authDeps.Config = runtimeDeps.Config
	}

	registerOperation(authAPI, huma.Operation{
		OperationID: "checkUser",
		Method:      http.MethodGet,
		Path:        "/check",
		Summary:     "Check user is exist",
		Tags:        []string{"auth"},
	}, auth.CheckUserHandler(authDeps))

	registerOperation(authAPI, huma.Operation{
		OperationID: "checkUserTelephone",
		Method:      http.MethodGet,
		Path:        "/check/telephone",
		Summary:     "Check user telephone is exist",
		Tags:        []string{"auth"},
	}, auth.CheckUserTelephoneHandler(authDeps))

	registerOperation(authAPI, huma.Operation{
		OperationID: "userLogin",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "User login",
		Tags:        []string{"auth"},
	}, auth.UserLoginHandler(authDeps))

	registerOperation(authAPI, huma.Operation{
		OperationID: "deviceLogin",
		Method:      http.MethodPost,
		Path:        "/login/device",
		Summary:     "Device Login",
		Tags:        []string{"auth"},
	}, auth.DeviceLoginHandler(authDeps))

	registerOperation(authAPI, huma.Operation{
		OperationID: "telephoneLogin",
		Method:      http.MethodPost,
		Path:        "/login/telephone",
		Summary:     "User Telephone login",
		Tags:        []string{"auth"},
	}, auth.TelephoneLoginHandler(authDeps))

	registerOperation(authAPI, huma.Operation{
		OperationID: "userRegister",
		Method:      http.MethodPost,
		Path:        "/register",
		Summary:     "User register",
		Tags:        []string{"auth"},
	}, auth.UserRegisterHandler(authDeps))

	registerOperation(authAPI, huma.Operation{
		OperationID: "telephoneUserRegister",
		Method:      http.MethodPost,
		Path:        "/register/telephone",
		Summary:     "User Telephone register",
		Tags:        []string{"auth"},
	}, auth.TelephoneUserRegisterHandler(authDeps))

	registerOperation(authAPI, huma.Operation{
		OperationID: "resetPassword",
		Method:      http.MethodPost,
		Path:        "/reset",
		Summary:     "Reset password",
		Tags:        []string{"auth"},
	}, auth.ResetPasswordHandler(authDeps))

	registerOperation(authAPI, huma.Operation{
		OperationID: "telephoneResetPassword",
		Method:      http.MethodPost,
		Path:        "/reset/telephone",
		Summary:     "Reset password",
		Tags:        []string{"auth"},
	}, auth.TelephoneResetPasswordHandler(authDeps))

	authOauthGroup := router.Group("/api/v1/auth/oauth")
	authOauthConfig := governedAPIConfig("Auth OAuth API", "1.0.0", "/api/v1/auth/oauth", "oauth")
	authOauthAPI := humagin.NewWithGroup(router, authOauthGroup, authOauthConfig)
	configureHumaAPI(authOauthAPI, compatibilityEnabled(runtimeDeps, specOnly))
	authOauthDeps := authOauth.Deps{}
	if runtimeDeps != nil {
		authOauthDeps.AuthModel = runtimeDeps.AuthModel
		authOauthDeps.UserModel = runtimeDeps.UserModel
		authOauthDeps.LogModel = runtimeDeps.LogModel
		authOauthDeps.SubscribeModel = runtimeDeps.SubscribeModel
		authOauthDeps.Redis = runtimeDeps.Redis
		authOauthDeps.Config = runtimeDeps.Config
	}

	// AppleLoginCallback stays raw Gin because it needs direct redirect primitives.
	authOauthGroup.POST("/callback/apple", authOauth.AppleLoginCallbackHandler(authOauthDeps))

	registerOperation(authOauthAPI, huma.Operation{
		OperationID: "oAuthLogin",
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "OAuth login",
		Tags:        []string{"oauth"},
	}, authOauth.OAuthLoginHandler(authOauthDeps))

	registerOperation(authOauthAPI, huma.Operation{
		OperationID: "oAuthLoginGetToken",
		Method:      http.MethodPost,
		Path:        "/login/token",
		Summary:     "OAuth login get token",
		Tags:        []string{"oauth"},
	}, authOauth.OAuthLoginGetTokenHandler(authOauthDeps))

	return []huma.API{authAPI, authOauthAPI}
}
