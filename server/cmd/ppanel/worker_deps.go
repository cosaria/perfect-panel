package ppanel

import (
	"context"
	"encoding/json"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	appbootstrap "github.com/perfect-panel/server/internal/bootstrap/app"
	appruntime "github.com/perfect-panel/server/internal/bootstrap/runtime"
	"github.com/perfect-panel/server/models/node"
	emailLogic "github.com/perfect-panel/server/worker/email"
	orderLogic "github.com/perfect-panel/server/worker/order"
	"github.com/perfect-panel/server/worker/registry"
	smslogic "github.com/perfect-panel/server/worker/sms"
	"github.com/perfect-panel/server/worker/subscription"
	"github.com/perfect-panel/server/worker/task"
	"github.com/perfect-panel/server/worker/traffic"
)

func newWorkerRegistryDeps(svcCtx *appbootstrap.ServiceContext, live *appruntime.LiveState) registry.Deps {
	deps := registry.Deps{
		Email: emailLogic.Deps{},
		SMS:   smslogic.Deps{},
	}
	if svcCtx == nil {
		return deps
	}

	deps.Email = newEmailWorkerDeps(svcCtx)
	deps.SMS = newSMSWorkerDeps(svcCtx)

	orderDeps := newOrderWorkerDeps(svcCtx, live)
	subscriptionDeps := newSubscriptionWorkerDeps(svcCtx)
	taskDeps := newTaskWorkerDeps(svcCtx, live)
	trafficDeps := newTrafficWorkerDeps(svcCtx, live)

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

func newEmailWorkerDeps(svcCtx *appbootstrap.ServiceContext) emailLogic.Deps {
	return emailLogic.Deps{
		DB:       svcCtx.DB,
		LogModel: svcCtx.LogModel,
		Config:   &svcCtx.Config,
	}
}

func newSMSWorkerDeps(svcCtx *appbootstrap.ServiceContext) smslogic.Deps {
	return smslogic.Deps{
		LogModel: svcCtx.LogModel,
		Config:   &svcCtx.Config,
	}
}

func newOrderWorkerDeps(svcCtx *appbootstrap.ServiceContext, live *appruntime.LiveState) orderLogic.Deps {
	return orderLogic.Deps{
		OrderModel:     svcCtx.OrderModel,
		PaymentModel:   svcCtx.PaymentModel,
		SubscribeModel: svcCtx.SubscribeModel,
		UserModel:      svcCtx.UserModel,
		CouponModel:    svcCtx.CouponModel,
		LogModel:       svcCtx.LogModel,
		DB:             svcCtx.DB,
		Queue:          svcCtx.Queue,
		Redis:          svcCtx.Redis,
		TelegramBot: func() *tgbotapi.BotAPI {
			if live != nil {
				return live.TelegramBot()
			}
			return svcCtx.TelegramBot
		},
		Config: &svcCtx.Config,
	}
}

func newSubscriptionWorkerDeps(svcCtx *appbootstrap.ServiceContext) subscription.Deps {
	return subscription.Deps{
		UserModel:      svcCtx.UserModel,
		SubscribeModel: svcCtx.SubscribeModel,
		Queue:          svcCtx.Queue,
		Config:         &svcCtx.Config,
	}
}

func newTaskWorkerDeps(svcCtx *appbootstrap.ServiceContext, live *appruntime.LiveState) task.Deps {
	return task.Deps{
		DB:             svcCtx.DB,
		SystemModel:    svcCtx.SystemModel,
		SubscribeModel: svcCtx.SubscribeModel,
		UserModel:      svcCtx.UserModel,
		SetExchangeRate: func(rate float64) {
			svcCtx.ExchangeRate = rate
			if live != nil {
				live.SetExchangeRate(rate)
			}
		},
		PrepareExchangeRate: func(from, to string) uint64 {
			if live == nil {
				return 0
			}
			return live.PrepareExchangeRate(from, to)
		},
		StoreExchangeRate: func(version uint64, from, to string, rate float64) bool {
			if live == nil {
				svcCtx.ExchangeRate = rate
				return true
			}
			if !live.StoreExchangeRate(version, from, to, rate) {
				svcCtx.ExchangeRate = live.ExchangeRate()
				return false
			}
			svcCtx.ExchangeRate = rate
			return true
		},
		Config: &svcCtx.Config,
	}
}

func newTrafficWorkerDeps(svcCtx *appbootstrap.ServiceContext, live *appruntime.LiveState) traffic.Deps {
	return traffic.Deps{
		DB:              svcCtx.DB,
		Redis:           svcCtx.Redis,
		Queue:           svcCtx.Queue,
		NodeModel:       svcCtx.NodeModel,
		UserModel:       svcCtx.UserModel,
		SubscribeModel:  svcCtx.SubscribeModel,
		TrafficLogModel: svcCtx.TrafficLogModel,
		NodeMultiplierManager: func() *node.Manager {
			if live != nil {
				return live.NodeMultiplierManager()
			}
			return svcCtx.NodeMultiplierManager
		},
		LoadNodeMultiplierManager: func(ctx context.Context) (*node.Manager, error) {
			return loadNodeMultiplierManager(ctx, svcCtx, live)
		},
		Config:   &svcCtx.Config,
		LogModel: svcCtx.LogModel,
	}
}

func loadNodeMultiplierManager(ctx context.Context, svcCtx *appbootstrap.ServiceContext, live *appruntime.LiveState) (*node.Manager, error) {
	if svcCtx == nil {
		return nil, nil
	}
	if live != nil {
		if manager := live.NodeMultiplierManager(); manager != nil {
			return manager, nil
		}
	}
	if svcCtx.NodeMultiplierManager != nil {
		return svcCtx.NodeMultiplierManager, nil
	}

	manager := node.NewNodeMultiplierManager(nil)
	if svcCtx.SystemModel != nil {
		data, err := svcCtx.SystemModel.FindNodeMultiplierConfig(ctx)
		if err != nil {
			return nil, err
		}
		if data != nil && data.Value != "" {
			var periods []node.TimePeriod
			if err := json.Unmarshal([]byte(data.Value), &periods); err != nil {
				return nil, err
			}
			manager = node.NewNodeMultiplierManager(periods)
		}
	}

	svcCtx.NodeMultiplierManager = manager
	if live != nil {
		live.SetNodeMultiplierManager(manager)
	}
	return manager, nil
}
