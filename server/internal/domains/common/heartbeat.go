package common

import (
	"context"
	"time"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
)

type HeartbeatOutput struct {
	Body *types.HeartbeatResponse
}

func HeartbeatHandler(deps Deps) func(context.Context, *struct{}) (*HeartbeatOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*HeartbeatOutput, error) {
		l := NewHeartbeatLogic(ctx, deps)
		resp, err := l.Heartbeat()
		if err != nil {
			return nil, err
		}
		return &HeartbeatOutput{Body: resp}, nil
	}
}

type HeartbeatLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewHeartbeatLogic Heartbeat
func NewHeartbeatLogic(ctx context.Context, deps Deps) *HeartbeatLogic {
	return &HeartbeatLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *HeartbeatLogic) Heartbeat() (resp *types.HeartbeatResponse, err error) {
	return &types.HeartbeatResponse{
		Status:    true,
		Message:   "service is alive",
		Timestamp: time.Now().Unix(),
	}, nil
}
