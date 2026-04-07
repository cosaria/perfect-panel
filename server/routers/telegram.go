package handler

import (
	"github.com/gin-gonic/gin"
	appruntime "github.com/perfect-panel/server/runtime"
	servicetelegram "github.com/perfect-panel/server/services/telegram"
)

func RegisterTelegramHandlers(router *gin.Engine, runtimeDeps *appruntime.Deps) {
	deps := servicetelegram.Deps{}
	if runtimeDeps != nil {
		deps.AuthModel = runtimeDeps.AuthModel
		deps.SystemModel = runtimeDeps.SystemModel
		deps.UserModel = runtimeDeps.UserModel
		deps.Redis = runtimeDeps.Redis
		deps.DB = runtimeDeps.DB
		deps.TelegramBot = runtimeDeps.TelegramBot
		deps.Config = runtimeDeps.Config
	}
	servicetelegram.RegisterTelegramHandlers(router, deps)
}
