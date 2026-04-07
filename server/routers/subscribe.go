package handler

import (
	"github.com/gin-gonic/gin"
	servicesubscribe "github.com/perfect-panel/server/services/subscribe"
	"github.com/perfect-panel/server/svc"
)

func RegisterSubscribeHandlers(router *gin.Engine, serverCtx *svc.ServiceContext) {
	servicesubscribe.RegisterSubscribeHandlers(router, serverCtx)
}
