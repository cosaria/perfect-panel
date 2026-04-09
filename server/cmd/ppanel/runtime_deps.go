package ppanel

import (
	appruntime "github.com/perfect-panel/server/runtime"
	"github.com/perfect-panel/server/svc"
)

func newLiveState(svcCtx *svc.ServiceContext) *appruntime.LiveState {
	live := appruntime.NewLiveState()
	if svcCtx == nil {
		return live
	}

	if svcCtx.Config.Currency.Unit != "" {
		live.PrepareExchangeRate(svcCtx.Config.Currency.Unit, "CNY")
	}
	live.SetExchangeRate(svcCtx.ExchangeRate)
	live.SetRestart(svcCtx.Restart)
	live.SetTelegramBot(svcCtx.TelegramBot)
	live.SetNodeMultiplierManager(svcCtx.NodeMultiplierManager)
	return live
}

func newRuntimeDeps(svcCtx *svc.ServiceContext, live *appruntime.LiveState) *appruntime.Deps {
	if svcCtx == nil {
		return nil
	}
	if live == nil {
		live = newLiveState(svcCtx)
	}

	deps := &appruntime.Deps{
		DB:                svcCtx.DB,
		Redis:             svcCtx.Redis,
		Config:            &svcCtx.Config,
		Queue:             svcCtx.Queue,
		Live:              live,
		AuthModel:         svcCtx.AuthModel,
		LogModel:          svcCtx.LogModel,
		NodeModel:         svcCtx.NodeModel,
		UserModel:         svcCtx.UserModel,
		OrderModel:        svcCtx.OrderModel,
		ClientModel:       svcCtx.ClientModel,
		TicketModel:       svcCtx.TicketModel,
		SystemModel:       svcCtx.SystemModel,
		CouponModel:       svcCtx.CouponModel,
		PaymentModel:      svcCtx.PaymentModel,
		DocumentModel:     svcCtx.DocumentModel,
		SubscribeModel:    svcCtx.SubscribeModel,
		TrafficLogModel:   svcCtx.TrafficLogModel,
		AnnouncementModel: svcCtx.AnnouncementModel,
		AuthLimiter:       svcCtx.AuthLimiter,
		DeviceManager:     svcCtx.DeviceManager,
	}
	if svcCtx.GeoIP != nil {
		deps.GeoIPDB = svcCtx.GeoIP.DB
	}
	return deps
}
