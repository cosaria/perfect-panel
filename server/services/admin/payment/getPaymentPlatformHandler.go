// huma:migrated
package payment

import (
	"context"
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
