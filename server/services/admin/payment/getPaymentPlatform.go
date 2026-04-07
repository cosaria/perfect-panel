package payment

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/payment"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetPaymentPlatformOutput struct {
	Body *types.PlatformResponse
}

func GetPaymentPlatformHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetPaymentPlatformOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetPaymentPlatformOutput, error) {
		l := NewGetPaymentPlatformLogic(ctx, svcCtx)
		resp, err := l.GetPaymentPlatform()
		if err != nil {
			return nil, err
		}
		return &GetPaymentPlatformOutput{Body: resp}, nil
	}
}

type GetPaymentPlatformLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get supported payment platform
func NewGetPaymentPlatformLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentPlatformLogic {
	return &GetPaymentPlatformLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaymentPlatformLogic) GetPaymentPlatform() (resp *types.PlatformResponse, err error) {
	resp = &types.PlatformResponse{
		List: payment.GetSupportedPlatforms(),
	}
	return
}
