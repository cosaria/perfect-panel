package worker

import (
	"reflect"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/internal/jobs/registry"
	"github.com/perfect-panel/server/internal/jobs/task"
)

func TestRegisterHandlersIncludesExchangeRateTask(t *testing.T) {
	mux := asynq.NewServeMux()
	registry.RegisterHandlers(mux, registry.Deps{
		ExchangeRate: task.NewRateLogic(task.Deps{}),
		Quota:        task.NewQuotaTaskLogic(task.Deps{}),
	})

	rateHandler, ratePattern := mux.Handler(asynq.NewTask(SchedulerExchangeRate, nil))
	if ratePattern != SchedulerExchangeRate {
		t.Fatalf("expected exchange rate task to be registered, got pattern %q", ratePattern)
	}

	quotaHandler, quotaPattern := mux.Handler(asynq.NewTask(ForthwithQuotaTask, nil))
	if quotaPattern != ForthwithQuotaTask {
		t.Fatalf("expected quota task to remain registered, got pattern %q", quotaPattern)
	}

	if reflect.TypeOf(rateHandler) != reflect.TypeOf(task.NewRateLogic(task.Deps{})) {
		t.Fatalf("expected exchange rate task to use RateLogic, got %T", rateHandler)
	}

	if reflect.TypeOf(quotaHandler) != reflect.TypeOf(task.NewQuotaTaskLogic(task.Deps{})) {
		t.Fatalf("expected quota task to use QuotaTaskLogic, got %T", quotaHandler)
	}
}

func TestScheduledTasksIncludeExchangeRateJob(t *testing.T) {
	registrations := scheduledTasks()
	for _, registration := range registrations {
		if registration.spec == "0 1 * * *" {
			if registration.taskType != SchedulerExchangeRate {
				t.Fatalf("expected 01:00 scheduler to enqueue %q, got %q", SchedulerExchangeRate, registration.taskType)
			}
			return
		}
	}

	t.Fatal("expected exchange rate scheduler registration at 01:00")
}
