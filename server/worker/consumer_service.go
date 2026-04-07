package worker

import (
	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/worker/registry"
)

type ConsumerService struct {
	svc    *svc.ServiceContext
	server *asynq.Server
}

func NewConsumerService(svc *svc.ServiceContext) *ConsumerService {
	return &ConsumerService{
		svc:    svc,
		server: initConsumerService(svc),
	}
}

func (m *ConsumerService) Start() {
	logger.Infof("start consumer service")
	mux := asynq.NewServeMux()
	// register tasks
	registry.RegisterHandlers(mux, m.svc)
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

func initConsumerService(svc *svc.ServiceContext) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: svc.Config.Redis.Host, Password: svc.Config.Redis.Pass, DB: 5},
		asynq.Config{
			IsFailure: func(err error) bool {
				logger.Error("consumer service error", logger.Field("error", err.Error()))
				return true
			},
			Concurrency: 20,
		},
	)
}
