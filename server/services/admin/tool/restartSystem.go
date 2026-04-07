package tool

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
)

func RestartSystemHandler(deps Deps) func(context.Context, *struct{}) (*struct{}, error) {
	return func(ctx context.Context, _ *struct{}) (*struct{}, error) {
		l := NewRestartSystemLogic(ctx, deps)
		if err := l.RestartSystem(); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type RestartSystemLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Restart System
func NewRestartSystemLogic(ctx context.Context, deps Deps) *RestartSystemLogic {
	return &RestartSystemLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *RestartSystemLogic) RestartSystem() error {
	l.Logger.Info("[RestartSystem]", logger.Field("info", "Restarting system"))
	go func() {
		if l.deps.Restart == nil {
			return
		}
		err := l.deps.Restart()
		if err != nil {
			l.Errorw("[RestartSystem]", logger.Field("error", err.Error()))
		}
		l.Logger.Info("[RestartSystem]", logger.Field("info", "System restarted"))
	}()
	return nil
}
