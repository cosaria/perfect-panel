package handler

import (
	"github.com/perfect-panel/server/initialize"
	appruntime "github.com/perfect-panel/server/runtime"
	adminSystem "github.com/perfect-panel/server/services/admin/system"
	adminTool "github.com/perfect-panel/server/services/admin/tool"
	servicetelegram "github.com/perfect-panel/server/services/telegram"
	publicPortal "github.com/perfect-panel/server/services/user/portal"
	publicUser "github.com/perfect-panel/server/services/user/user"
)

func newPublicPortalDeps(runtimeDeps *appruntime.Deps) publicPortal.Deps {
	deps := publicPortal.Deps{}
	if runtimeDeps == nil {
		return deps
	}

	deps.PaymentModel = runtimeDeps.PaymentModel
	deps.SubscribeModel = runtimeDeps.SubscribeModel
	deps.CouponModel = runtimeDeps.CouponModel
	deps.OrderModel = runtimeDeps.OrderModel
	deps.UserModel = runtimeDeps.UserModel
	deps.DB = runtimeDeps.DB
	deps.Redis = runtimeDeps.Redis
	deps.Queue = runtimeDeps.Queue
	deps.Config = runtimeDeps.Config
	if runtimeDeps.Live != nil {
		deps.GetExchangeRate = runtimeDeps.Live.ExchangeRate
		deps.SetExchangeRate = runtimeDeps.Live.SetExchangeRate
		deps.GetExchangeRateSnapshot = func() publicPortal.ExchangeRateSnapshot {
			quote := runtimeDeps.Live.ExchangeRateQuote()
			return publicPortal.ExchangeRateSnapshot{
				Version: quote.Version,
				From:    quote.From,
				To:      quote.To,
				Rate:    quote.Rate,
			}
		}
		deps.PrepareExchangeRate = runtimeDeps.Live.PrepareExchangeRate
		deps.StoreExchangeRate = runtimeDeps.Live.StoreExchangeRate
	}
	return deps
}

func newPublicUserDeps(runtimeDeps *appruntime.Deps) publicUser.Deps {
	deps := publicUser.Deps{}
	if runtimeDeps == nil {
		return deps
	}

	deps.UserModel = runtimeDeps.UserModel
	deps.LogModel = runtimeDeps.LogModel
	deps.AuthModel = runtimeDeps.AuthModel
	deps.OrderModel = runtimeDeps.OrderModel
	deps.SubscribeModel = runtimeDeps.SubscribeModel
	deps.Redis = runtimeDeps.Redis
	deps.Config = runtimeDeps.Config
	deps.DB = runtimeDeps.DB
	if runtimeDeps.Live != nil {
		deps.TelegramBot = runtimeDeps.Live.TelegramBot
	}
	return deps
}

func newAdminSystemDeps(runtimeDeps *appruntime.Deps, initDeps initialize.Deps) adminSystem.Deps {
	deps := adminSystem.Deps{}
	if runtimeDeps == nil {
		return deps
	}

	deps.SystemModel = runtimeDeps.SystemModel
	deps.Redis = runtimeDeps.Redis
	deps.Config = runtimeDeps.Config
	if runtimeDeps.Live != nil {
		deps.NodeMultiplierManager = runtimeDeps.Live.NodeMultiplierManager
		deps.Restart = func() error {
			restart := runtimeDeps.Live.Restart()
			if restart == nil {
				return nil
			}
			return restart()
		}
	}
	deps.ReloadVerify = func() { initialize.Verify(initDeps) }
	deps.ReloadNode = func() { initialize.Node(initDeps) }
	deps.ReloadCurrency = func() { initialize.Currency(initDeps) }
	deps.ReloadInvite = func() { initialize.Invite(initDeps) }
	deps.ReloadRegister = func() { initialize.Register(initDeps) }
	deps.ReloadSite = func() { initialize.Site(initDeps) }
	deps.ReloadSubscribe = func() { initialize.Subscribe(initDeps) }
	deps.ReloadTelegram = func() { initialize.Telegram(initDeps) }
	return deps
}

func newAdminToolDeps(runtimeDeps *appruntime.Deps) adminTool.Deps {
	deps := adminTool.Deps{}
	if runtimeDeps == nil {
		return deps
	}

	deps.Config = runtimeDeps.Config
	deps.GeoIPDB = runtimeDeps.GeoIPDB
	if runtimeDeps.Live != nil {
		deps.Restart = func() error {
			restart := runtimeDeps.Live.Restart()
			if restart == nil {
				return nil
			}
			return restart()
		}
	}
	return deps
}

func newTelegramServiceDeps(runtimeDeps *appruntime.Deps) servicetelegram.Deps {
	deps := servicetelegram.Deps{}
	if runtimeDeps == nil {
		return deps
	}

	deps.AuthModel = runtimeDeps.AuthModel
	deps.SystemModel = runtimeDeps.SystemModel
	deps.UserModel = runtimeDeps.UserModel
	deps.Redis = runtimeDeps.Redis
	deps.DB = runtimeDeps.DB
	deps.Config = runtimeDeps.Config
	if runtimeDeps.Live != nil {
		deps.TelegramBot = runtimeDeps.Live.TelegramBot
	}
	return deps
}
