package ppanel

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/perfect-panel/server/initialize"
	"github.com/perfect-panel/server/models/node"
	appruntime "github.com/perfect-panel/server/runtime"
	"github.com/perfect-panel/server/svc"
)

func newInitializeDeps(svcCtx *svc.ServiceContext, live *appruntime.LiveState) initialize.Deps {
	deps := initialize.Deps{}
	if svcCtx == nil {
		return deps
	}
	if live == nil {
		live = newLiveState(svcCtx)
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
			live.SetExchangeRate(rate)
		},
		PrepareExchangeRate: func(from, to string) uint64 {
			version := live.PrepareExchangeRate(from, to)
			svcCtx.ExchangeRate = live.ExchangeRate()
			return version
		},
		StoreExchangeRate: func(version uint64, from, to string, rate float64) bool {
			if !live.StoreExchangeRate(version, from, to, rate) {
				svcCtx.ExchangeRate = live.ExchangeRate()
				return false
			}
			svcCtx.ExchangeRate = rate
			return true
		},
		SetNodeMultiplierManager: func(manager *node.Manager) {
			svcCtx.NodeMultiplierManager = manager
			live.SetNodeMultiplierManager(manager)
		},
		SetTelegramBot: func(bot *tgbotapi.BotAPI) {
			svcCtx.TelegramBot = bot
			live.SetTelegramBot(bot)
		},
		SwapTelegramPoller: live.SwapTelegramPoller,
	}
}
