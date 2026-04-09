package registry

import (
	"github.com/hibiken/asynq"
	smslogic "github.com/perfect-panel/server/internal/jobs/sms"
	"github.com/perfect-panel/server/internal/jobs/spec"

	emailLogic "github.com/perfect-panel/server/internal/jobs/email"
)

type Deps struct {
	Email             emailLogic.Deps
	SMS               smslogic.Deps
	DeferCloseOrder   asynq.Handler
	ActivateOrder     asynq.Handler
	TrafficStatistics asynq.Handler
	CheckSubscription asynq.Handler
	ServerData        asynq.Handler
	ResetTraffic      asynq.Handler
	TrafficStat       asynq.Handler
	ExchangeRate      asynq.Handler
	Quota             asynq.Handler
}

func RegisterHandlers(mux *asynq.ServeMux, deps Deps) {
	// Send email task
	mux.Handle(spec.ForthwithSendEmail, emailLogic.NewSendEmailLogic(deps.Email))
	// Send sms task
	mux.Handle(spec.ForthwithSendSms, smslogic.NewSendSmsLogic(deps.SMS))
	// Defer close order task
	if deps.DeferCloseOrder != nil {
		mux.Handle(spec.DeferCloseOrder, deps.DeferCloseOrder)
	}
	// Forthwith activate order task
	if deps.ActivateOrder != nil {
		mux.Handle(spec.ForthwithActivateOrder, deps.ActivateOrder)
	}

	// Forthwith traffic statistics
	if deps.TrafficStatistics != nil {
		mux.Handle(spec.ForthwithTrafficStatistics, deps.TrafficStatistics)
	}

	// Schedule check subscription
	if deps.CheckSubscription != nil {
		mux.Handle(spec.SchedulerCheckSubscription, deps.CheckSubscription)
	}

	// Schedule total server data
	if deps.ServerData != nil {
		mux.Handle(spec.SchedulerTotalServerData, deps.ServerData)
	}

	// Schedule reset traffic
	if deps.ResetTraffic != nil {
		mux.Handle(spec.SchedulerResetTraffic, deps.ResetTraffic)
	}

	// ScheduledBatchSendEmail
	mux.Handle(spec.ScheduledBatchSendEmail, emailLogic.NewBatchEmailLogic(deps.Email))

	// ScheduledTrafficStat
	if deps.TrafficStat != nil {
		mux.Handle(spec.SchedulerTrafficStat, deps.TrafficStat)
	}

	// SchedulerExchangeRate
	if deps.ExchangeRate != nil {
		mux.Handle(spec.SchedulerExchangeRate, deps.ExchangeRate)
	}

	// ForthwithQuotaTask
	if deps.Quota != nil {
		mux.Handle(spec.ForthwithQuotaTask, deps.Quota)
	}
}
