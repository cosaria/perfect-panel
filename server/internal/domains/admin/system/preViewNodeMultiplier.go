package system

import (
	"context"
	"time"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
)

type PreViewNodeMultiplierOutput struct {
	Body *types.PreViewNodeMultiplierResponse
}

func PreViewNodeMultiplierHandler(deps Deps) func(context.Context, *struct{}) (*PreViewNodeMultiplierOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*PreViewNodeMultiplierOutput, error) {
		l := NewPreViewNodeMultiplierLogic(ctx, deps)
		resp, err := l.PreViewNodeMultiplier()
		if err != nil {
			return nil, err
		}
		return &PreViewNodeMultiplierOutput{Body: resp}, nil
	}
}

type PreViewNodeMultiplierLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// PreView Node Multiplier
func NewPreViewNodeMultiplierLogic(ctx context.Context, deps Deps) *PreViewNodeMultiplierLogic {
	return &PreViewNodeMultiplierLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *PreViewNodeMultiplierLogic) PreViewNodeMultiplier() (resp *types.PreViewNodeMultiplierResponse, err error) {
	now := time.Now()
	ratio := float32(1)
	if manager := l.deps.CurrentNodeMultiplierManager(); manager != nil {
		ratio = manager.GetMultiplier(now)
	}
	return &types.PreViewNodeMultiplierResponse{
		Ratio:       ratio,
		CurrentTime: now.Format("2006-01-02 15:04:05"),
	}, nil
}
