package handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	configinit "github.com/perfect-panel/server/internal/bootstrap/configinit"
	appruntime "github.com/perfect-panel/server/internal/bootstrap/runtime"
	"github.com/perfect-panel/server/internal/platform/persistence/node"
)

func initializeDepsFromRuntimeDeps(runtimeDeps *appruntime.Deps) configinit.Deps {
	deps := configinit.Deps{}
	if runtimeDeps == nil {
		return deps
	}

	return configinit.Deps{
		DB:          runtimeDeps.DB,
		Redis:       runtimeDeps.Redis,
		Config:      runtimeDeps.Config,
		AuthModel:   runtimeDeps.AuthModel,
		SystemModel: runtimeDeps.SystemModel,
		UserModel:   runtimeDeps.UserModel,
		SetExchangeRate: func(rate float64) {
			if runtimeDeps.Live != nil {
				runtimeDeps.Live.SetExchangeRate(rate)
			}
		},
		PrepareExchangeRate: func(from, to string) uint64 {
			if runtimeDeps.Live == nil {
				return 0
			}
			return runtimeDeps.Live.PrepareExchangeRate(from, to)
		},
		StoreExchangeRate: func(version uint64, from, to string, rate float64) bool {
			if runtimeDeps.Live == nil {
				return false
			}
			return runtimeDeps.Live.StoreExchangeRate(version, from, to, rate)
		},
		SetNodeMultiplierManager: func(manager *node.Manager) {
			if runtimeDeps.Live != nil {
				runtimeDeps.Live.SetNodeMultiplierManager(manager)
			}
		},
		SetTelegramBot: func(bot *tgbotapi.BotAPI) {
			if runtimeDeps.Live != nil {
				runtimeDeps.Live.SetTelegramBot(bot)
			}
		},
		SwapTelegramPoller: func(next func()) func() {
			if runtimeDeps.Live == nil {
				return nil
			}
			return runtimeDeps.Live.SwapTelegramPoller(next)
		},
	}
}
