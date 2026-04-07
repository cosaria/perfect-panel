package worker

import (
	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/worker/registry"
)

type ConsumerService struct {
	deps   registry.Deps
	server *asynq.Server
}

func NewConsumerService(cfg config.Config, deps registry.Deps) *ConsumerService {
	return &ConsumerService{
		deps:   deps,
		server: initConsumerService(cfg),
	}
}

func (m *ConsumerService) Start() {
	logger.Infof("start consumer service")
	mux := asynq.NewServeMux()
	// register tasks
	registry.RegisterHandlers(mux, m.deps)
	if err := m.server.Run(mux); err != nil {
		logger.Error("consumer service error", logger.LogField{
			Key:   "error",
			Value: err.Error(),
		})
	}
}

func (m *ConsumerService) Stop() {
	logger.Info("stop consumer service")
	m.server.Stop()
}

func initConsumerService(cfg config.Config) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Redis.Host, Password: cfg.Redis.Pass, DB: 5},
		asynq.Config{
			IsFailure: func(err error) bool {
				logger.Error("consumer service error", logger.Field("error", err.Error()))
				return true
			},
			Concurrency: 20,
		},
	)
}
