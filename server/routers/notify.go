package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/routers/middleware"
	appruntime "github.com/perfect-panel/server/runtime"
	"github.com/perfect-panel/server/services/notify"
)

func RegisterNotifyHandlers(router *gin.Engine, runtimeDeps *appruntime.Deps) {
	group := router.Group("/api/v1/notify/")
	group.Use(middleware.NotifyMiddleware(runtimeDeps))
	notifyDeps := notify.Deps{}
	if runtimeDeps != nil {
		notifyDeps.OrderModel = runtimeDeps.OrderModel
		notifyDeps.Queue = runtimeDeps.Queue
		notifyDeps.Config = runtimeDeps.Config
	}
	{
		group.Any("/:platform/:token", notify.PaymentNotifyHandler(notifyDeps))
	}

}
