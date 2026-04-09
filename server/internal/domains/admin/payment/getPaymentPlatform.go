package payment

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/payment"
	"github.com/perfect-panel/server/internal/platform/support/logger"
)

type GetPaymentPlatformOutput struct {
	Body *types.PlatformResponse
}

func GetPaymentPlatformHandler(deps Deps) func(context.Context, *struct{}) (*GetPaymentPlatformOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetPaymentPlatformOutput, error) {
		l := NewGetPaymentPlatformLogic(ctx, deps)
		resp, err := l.GetPaymentPlatform()
		if err != nil {
			return nil, err
		}
		return &GetPaymentPlatformOutput{Body: resp}, nil
	}
}

type GetPaymentPlatformLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get supported payment platform
func NewGetPaymentPlatformLogic(ctx context.Context, deps Deps) *GetPaymentPlatformLogic {
	return &GetPaymentPlatformLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetPaymentPlatformLogic) GetPaymentPlatform() (resp *types.PlatformResponse, err error) {
	resp = &types.PlatformResponse{
		List: payment.GetSupportedPlatforms(),
	}
	return
}
