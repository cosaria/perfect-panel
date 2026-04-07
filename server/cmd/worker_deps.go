package cmd

import (
	emailLogic "github.com/perfect-panel/server/worker/email"
	orderLogic "github.com/perfect-panel/server/worker/order"
	"github.com/perfect-panel/server/worker/registry"
	smslogic "github.com/perfect-panel/server/worker/sms"
	"github.com/perfect-panel/server/worker/subscription"
	"github.com/perfect-panel/server/worker/task"
	"github.com/perfect-panel/server/worker/traffic"

	"github.com/perfect-panel/server/svc"
)

func newWorkerRegistryDeps(svcCtx *svc.ServiceContext) registry.Deps {
	deps := registry.Deps{
		Email: emailLogic.Deps{},
		SMS:   smslogic.Deps{},
	}
	if svcCtx == nil {
		return deps
	}

	emailDeps := emailLogic.Deps{
		DB:       svcCtx.DB,
		LogModel: svcCtx.LogModel,
		Config:   &svcCtx.Config,
	}
	smsDeps := smslogic.Deps{
		LogModel: svcCtx.LogModel,
		Config:   &svcCtx.Config,
	}
	orderDeps := orderLogic.Deps{
		OrderModel:     svcCtx.OrderModel,
		PaymentModel:   svcCtx.PaymentModel,
		SubscribeModel: svcCtx.SubscribeModel,
		UserModel:      svcCtx.UserModel,
		CouponModel:    svcCtx.CouponModel,
		LogModel:       svcCtx.LogModel,
		DB:             svcCtx.DB,
		Queue:          svcCtx.Queue,
		Redis:          svcCtx.Redis,
		TelegramBot:    svcCtx.TelegramBot,
		Config:         &svcCtx.Config,
	}
	subscriptionDeps := subscription.Deps{
		UserModel:      svcCtx.UserModel,
		SubscribeModel: svcCtx.SubscribeModel,
		Queue:          svcCtx.Queue,
		Config:         &svcCtx.Config,
	}
	taskDeps := task.Deps{
		DB:             svcCtx.DB,
		SystemModel:    svcCtx.SystemModel,
		SubscribeModel: svcCtx.SubscribeModel,
		UserModel:      svcCtx.UserModel,
		SetExchangeRate: func(rate float64) {
			svcCtx.ExchangeRate = rate
		},
		Config: &svcCtx.Config,
	}
	trafficDeps := traffic.Deps{
		DB:                    svcCtx.DB,
		Redis:                 svcCtx.Redis,
		Queue:                 svcCtx.Queue,
		NodeModel:             svcCtx.NodeModel,
		UserModel:             svcCtx.UserModel,
		SubscribeModel:        svcCtx.SubscribeModel,
		TrafficLogModel:       svcCtx.TrafficLogModel,
		NodeMultiplierManager: svcCtx.NodeMultiplierManager,
		Config:                &svcCtx.Config,
		LogModel:              svcCtx.LogModel,
	}

	deps.Email = emailDeps
	deps.SMS = smsDeps
	deps.DeferCloseOrder = orderLogic.NewDeferCloseOrderLogic(orderDeps)
	deps.ActivateOrder = orderLogic.NewActivateOrderLogic(orderDeps)
	deps.TrafficStatistics = traffic.NewTrafficStatisticsLogic(trafficDeps)
	deps.CheckSubscription = subscription.NewCheckSubscriptionLogic(subscriptionDeps)
	deps.ServerData = traffic.NewServerDataLogic(trafficDeps)
	deps.ResetTraffic = traffic.NewResetTrafficLogic(trafficDeps)
	deps.TrafficStat = traffic.NewStatLogic(trafficDeps)
	deps.ExchangeRate = task.NewRateLogic(taskDeps)
	deps.Quota = task.NewQuotaTaskLogic(taskDeps)
	return deps
}
