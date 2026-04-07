package tool

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/svc"
)

func RestartSystemHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*struct{}, error) {
	return func(ctx context.Context, _ *struct{}) (*struct{}, error) {
		l := NewRestartSystemLogic(ctx, svcCtx)
		if err := l.RestartSystem(); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type RestartSystemLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Restart System
func NewRestartSystemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RestartSystemLogic {
	return &RestartSystemLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RestartSystemLogic) RestartSystem() error {
	l.Logger.Info("[RestartSystem]", logger.Field("info", "Restarting system"))
	go func() {
		err := l.svcCtx.Restart()
		if err != nil {
			l.Errorw("[RestartSystem]", logger.Field("error", err.Error()))
		}
		l.Logger.Info("[RestartSystem]", logger.Field("info", "System restarted"))
	}()
	return nil
}
