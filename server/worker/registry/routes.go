package registry

import (
	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/svc"
	orderLogic "github.com/perfect-panel/server/worker/order"
	smslogic "github.com/perfect-panel/server/worker/sms"
	"github.com/perfect-panel/server/worker/spec"
	"github.com/perfect-panel/server/worker/subscription"
	"github.com/perfect-panel/server/worker/task"
	"github.com/perfect-panel/server/worker/traffic"

	emailLogic "github.com/perfect-panel/server/worker/email"
)

func RegisterHandlers(mux *asynq.ServeMux, serverCtx *svc.ServiceContext) {
	// Send email task
	mux.Handle(spec.ForthwithSendEmail, emailLogic.NewSendEmailLogic(serverCtx))
	// Send sms task
	mux.Handle(spec.ForthwithSendSms, smslogic.NewSendSmsLogic(serverCtx))
	// Defer close order task
	mux.Handle(spec.DeferCloseOrder, orderLogic.NewDeferCloseOrderLogic(serverCtx))
	// Forthwith activate order task
	mux.Handle(spec.ForthwithActivateOrder, orderLogic.NewActivateOrderLogic(serverCtx))

	// Forthwith traffic statistics
	mux.Handle(spec.ForthwithTrafficStatistics, traffic.NewTrafficStatisticsLogic(serverCtx))

	// Schedule check subscription
	mux.Handle(spec.SchedulerCheckSubscription, subscription.NewCheckSubscriptionLogic(serverCtx))

	// Schedule total server data
	mux.Handle(spec.SchedulerTotalServerData, traffic.NewServerDataLogic(serverCtx))

	// Schedule reset traffic
	mux.Handle(spec.SchedulerResetTraffic, traffic.NewResetTrafficLogic(serverCtx))

	// ScheduledBatchSendEmail
	mux.Handle(spec.ScheduledBatchSendEmail, emailLogic.NewBatchEmailLogic(serverCtx))

	// ScheduledTrafficStat
	mux.Handle(spec.SchedulerTrafficStat, traffic.NewStatLogic(serverCtx))

	// SchedulerExchangeRate
	mux.Handle(spec.SchedulerExchangeRate, task.NewRateLogic(serverCtx))

	// ForthwithQuotaTask
	mux.Handle(spec.ForthwithQuotaTask, task.NewQuotaTaskLogic(serverCtx))
}
