package worker

import (
	"time"

	"github.com/perfect-panel/server/modules/infra/logger"

	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/config"
)

type SchedulerService struct {
	server *asynq.Scheduler
}

type scheduledTaskRegistration struct {
	spec     string
	taskType string
	options  []asynq.Option
}

func NewSchedulerService(cfg config.Config) *SchedulerService {
	return &SchedulerService{
		server: initSchedulerService(cfg),
	}
}

func (m *SchedulerService) Start() {
	logger.Infof("start scheduler service")
	for _, registration := range scheduledTasks() {
		task := asynq.NewTask(registration.taskType, nil)
		if _, err := m.server.Register(registration.spec, task, registration.options...); err != nil {
			logger.Errorf("register scheduled task %s failed: %s", registration.taskType, err.Error())
		}
	}

	if err := m.server.Run(); err != nil {
		logger.Errorf("run scheduler failed: %s", err.Error())
	}
}

func (m *SchedulerService) Stop() {
	logger.Info("stop scheduler service")
	m.server.Shutdown()
}

func initSchedulerService(cfg config.Config) *asynq.Scheduler {
	location, _ := time.LoadLocation("Asia/Shanghai")
	return asynq.NewScheduler(
		asynq.RedisClientOpt{Addr: cfg.Redis.Host, Password: cfg.Redis.Pass, DB: 5},
		&asynq.SchedulerOpts{
			Location: location,
		},
	)
}

func scheduledTasks() []scheduledTaskRegistration {
	return []scheduledTaskRegistration{
		{spec: "@every 60s", taskType: SchedulerCheckSubscription},
		{spec: "30 0 * * *", taskType: SchedulerResetTraffic},
		{spec: "0 0 * * *", taskType: SchedulerTrafficStat, options: []asynq.Option{asynq.MaxRetry(3)}},
		{spec: "0 1 * * *", taskType: SchedulerExchangeRate, options: []asynq.Option{asynq.MaxRetry(3)}},
	}
}
