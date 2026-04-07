package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/perfect-panel/server/initialize"
	"github.com/perfect-panel/server/models/node"
	appruntime "github.com/perfect-panel/server/runtime"
)

func initializeDepsFromRuntimeDeps(runtimeDeps *appruntime.Deps) initialize.Deps {
	deps := initialize.Deps{}
	if runtimeDeps == nil {
		return deps
	}

	return initialize.Deps{
		DB:          runtimeDeps.DB,
		Redis:       runtimeDeps.Redis,
		Config:      runtimeDeps.Config,
		AuthModel:   runtimeDeps.AuthModel,
		SystemModel: runtimeDeps.SystemModel,
		UserModel:   runtimeDeps.UserModel,
		SetExchangeRate: func(rate float64) {
			runtimeDeps.ExchangeRate = rate
		},
		SetNodeMultiplierManager: func(manager *node.Manager) {
			runtimeDeps.NodeMultiplierManager = manager
		},
		SetTelegramBot: func(bot *tgbotapi.BotAPI) {
			runtimeDeps.TelegramBot = bot
		},
	}
}
