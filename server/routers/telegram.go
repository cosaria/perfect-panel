package handler

import (
	"github.com/gin-gonic/gin"
	servicetelegram "github.com/perfect-panel/server/services/telegram"
	"github.com/perfect-panel/server/svc"
)

func RegisterTelegramHandlers(router *gin.Engine, serverCtx *svc.ServiceContext) {
	servicetelegram.RegisterTelegramHandlers(router, serverCtx)
}
