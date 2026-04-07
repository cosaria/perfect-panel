package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/routers/middleware"
	server "github.com/perfect-panel/server/services/node"
	"github.com/perfect-panel/server/svc"
)

func registerServerRoutes(router *gin.Engine, serverCtx *svc.ServiceContext, specOnly bool) {
	if specOnly {
		return
	}

	serverGroup := router.Group("/api/v1/server")
	serverGroup.Use(middleware.ServerMiddleware(serverCtx))

	serverGroup.GET("/config", server.GetServerConfigHandler(serverCtx))
	serverGroup.POST("/online", server.PushOnlineUsersHandler(serverCtx))
	serverGroup.POST("/push", server.ServerPushUserTrafficHandler(serverCtx))
	serverGroup.POST("/status", server.ServerPushStatusHandler(serverCtx))
	serverGroup.GET("/user", server.GetServerUserListHandler(serverCtx))
}
