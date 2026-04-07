package system

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"time"
)

type PreViewNodeMultiplierOutput struct {
	Body *types.PreViewNodeMultiplierResponse
}

func PreViewNodeMultiplierHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*PreViewNodeMultiplierOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*PreViewNodeMultiplierOutput, error) {
		l := NewPreViewNodeMultiplierLogic(ctx, svcCtx)
		resp, err := l.PreViewNodeMultiplier()
		if err != nil {
			return nil, err
		}
		return &PreViewNodeMultiplierOutput{Body: resp}, nil
	}
}

type PreViewNodeMultiplierLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// PreView Node Multiplier
func NewPreViewNodeMultiplierLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreViewNodeMultiplierLogic {
	return &PreViewNodeMultiplierLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreViewNodeMultiplierLogic) PreViewNodeMultiplier() (resp *types.PreViewNodeMultiplierResponse, err error) {
	now := time.Now()
	ratio := l.svcCtx.NodeMultiplierManager.GetMultiplier(now)
	return &types.PreViewNodeMultiplierResponse{
		Ratio:       ratio,
		CurrentTime: now.Format("2006-01-02 15:04:05"),
	}, nil
}
