package telegram

import (
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/routers/response"
	"github.com/perfect-panel/server/svc"
)

func RegisterTelegramHandlers(router *gin.Engine, serverCtx *svc.ServiceContext) {
	router.POST("/api/v1/telegram/webhook", TelegramHandler(serverCtx))
}

func TelegramHandler(svcCtx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		// auth secret
		secret := c.Query("secret")
		if secret != tool.Md5Encode(svcCtx.Config.Telegram.BotToken, false) {
			logger.WithContext(c.Request.Context()).Error("[TelegramHandler] Secret is wrong", logger.Field("request secret", secret), logger.Field("config secret", tool.Md5Encode(svcCtx.Config.Telegram.BotToken, false)), logger.Field("token", svcCtx.Config.Telegram.BotToken))
			c.Abort()
			response.HttpResult(c, nil, nil)
			return
		}
		var request tgbotapi.Update
		if err := c.BindJSON(&request); err != nil {
			logger.WithContext(c.Request.Context()).Error("[TelegramHandler] Failed to bind request", logger.Field("error", err.Error()))
			c.Abort()
			response.HttpResult(c, nil, err)
		}
		l := NewTelegramLogic(c, svcCtx)
		l.TelegramLogic(&request)
	}
}
