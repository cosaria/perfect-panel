package cmd

import (
	appruntime "github.com/perfect-panel/server/runtime"
	"github.com/perfect-panel/server/svc"
)

func newRuntimeDeps(svcCtx *svc.ServiceContext) *appruntime.Deps {
	if svcCtx == nil {
		return nil
	}

	deps := &appruntime.Deps{
		DB:                    svcCtx.DB,
		Redis:                 svcCtx.Redis,
		Config:                &svcCtx.Config,
		Queue:                 svcCtx.Queue,
		ExchangeRate:          svcCtx.ExchangeRate,
		AuthModel:             svcCtx.AuthModel,
		AdsModel:              svcCtx.AdsModel,
		LogModel:              svcCtx.LogModel,
		NodeModel:             svcCtx.NodeModel,
		UserModel:             svcCtx.UserModel,
		OrderModel:            svcCtx.OrderModel,
		ClientModel:           svcCtx.ClientModel,
		TicketModel:           svcCtx.TicketModel,
		SystemModel:           svcCtx.SystemModel,
		CouponModel:           svcCtx.CouponModel,
		PaymentModel:          svcCtx.PaymentModel,
		DocumentModel:         svcCtx.DocumentModel,
		SubscribeModel:        svcCtx.SubscribeModel,
		TrafficLogModel:       svcCtx.TrafficLogModel,
		AnnouncementModel:     svcCtx.AnnouncementModel,
		Restart:               svcCtx.Restart,
		TelegramBot:           svcCtx.TelegramBot,
		NodeMultiplierManager: svcCtx.NodeMultiplierManager,
		AuthLimiter:           svcCtx.AuthLimiter,
		DeviceManager:         svcCtx.DeviceManager,
	}
	if svcCtx.GeoIP != nil {
		deps.GeoIPDB = svcCtx.GeoIP.DB
	}
	return deps
}
