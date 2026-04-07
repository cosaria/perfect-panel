package handler

import (
	"github.com/gin-gonic/gin"
	appruntime "github.com/perfect-panel/server/runtime"
)

func registerUserRoutes(router *gin.Engine, runtimeDeps *appruntime.Deps, specOnly bool, apis *APIs) {
	userAPIs := registerAuthRoutes(router, runtimeDeps, specOnly)
	userAPIs = append(userAPIs, registerPublicRoutes(router, runtimeDeps, specOnly)...)
	apis.userAPIs = userAPIs
}
