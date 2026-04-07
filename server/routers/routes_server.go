package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/routers/middleware"
	appruntime "github.com/perfect-panel/server/runtime"
	server "github.com/perfect-panel/server/services/node"
)

func registerServerRoutes(router *gin.Engine, runtimeDeps *appruntime.Deps, specOnly bool) {
	if specOnly {
		return
	}

	serverDeps := server.Deps{}
	if runtimeDeps != nil {
		serverDeps.NodeModel = runtimeDeps.NodeModel
		serverDeps.SubscribeModel = runtimeDeps.SubscribeModel
		serverDeps.UserModel = runtimeDeps.UserModel
		serverDeps.Redis = runtimeDeps.Redis
		serverDeps.Queue = runtimeDeps.Queue
		serverDeps.Config = runtimeDeps.Config
	}

	serverGroup := router.Group("/api/v1/server")
	serverGroup.Use(middleware.ServerMiddleware(runtimeDeps))

	serverGroup.GET("/config", server.GetServerConfigHandler(serverDeps))
	serverGroup.POST("/online", server.PushOnlineUsersHandler(serverDeps))
	serverGroup.POST("/push", server.ServerPushUserTrafficHandler(serverDeps))
	serverGroup.POST("/status", server.ServerPushStatusHandler(serverDeps))
	serverGroup.GET("/user", server.GetServerUserListHandler(serverDeps))
}
