package cmd

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/perfect-panel/server/initialize"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/svc"
)

func newInitializeDeps(svcCtx *svc.ServiceContext) initialize.Deps {
	deps := initialize.Deps{}
	if svcCtx == nil {
		return deps
	}

	return initialize.Deps{
		DB:          svcCtx.DB,
		Redis:       svcCtx.Redis,
		Config:      &svcCtx.Config,
		AuthModel:   svcCtx.AuthModel,
		SystemModel: svcCtx.SystemModel,
		UserModel:   svcCtx.UserModel,
		SetExchangeRate: func(rate float64) {
			svcCtx.ExchangeRate = rate
		},
		SetNodeMultiplierManager: func(manager *node.Manager) {
			svcCtx.NodeMultiplierManager = manager
		},
		SetTelegramBot: func(bot *tgbotapi.BotAPI) {
			svcCtx.TelegramBot = bot
		},
	}
}
