package handler

import (
	"github.com/gin-gonic/gin"
	appruntime "github.com/perfect-panel/server/internal/bootstrap/runtime"
	servicetelegram "github.com/perfect-panel/server/internal/domains/telegram"
)

func RegisterTelegramHandlers(router *gin.Engine, runtimeDeps *appruntime.Deps) {
	deps := newTelegramServiceDeps(runtimeDeps)
	servicetelegram.RegisterTelegramHandlers(router, deps)
}
