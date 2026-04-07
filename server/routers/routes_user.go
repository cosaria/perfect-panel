package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/svc"
)

func registerUserRoutes(router *gin.Engine, serverCtx *svc.ServiceContext, specOnly bool, apis *APIs) {
	userAPIs := registerAuthRoutes(router, serverCtx, specOnly)
	userAPIs = append(userAPIs, registerPublicRoutes(router, serverCtx, specOnly)...)
	apis.userAPIs = userAPIs
}
