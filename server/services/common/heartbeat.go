package common

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"time"
)

type HeartbeatOutput struct {
	Body *types.HeartbeatResponse
}

func HeartbeatHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*HeartbeatOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*HeartbeatOutput, error) {
		l := NewHeartbeatLogic(ctx, svcCtx)
		resp, err := l.Heartbeat()
		if err != nil {
			return nil, err
		}
		return &HeartbeatOutput{Body: resp}, nil
	}
}

type HeartbeatLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewHeartbeatLogic Heartbeat
func NewHeartbeatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HeartbeatLogic {
	return &HeartbeatLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HeartbeatLogic) Heartbeat() (resp *types.HeartbeatResponse, err error) {
	return &types.HeartbeatResponse{
		Status:    true,
		Message:   "service is alive",
		Timestamp: time.Now().Unix(),
	}, nil
}
