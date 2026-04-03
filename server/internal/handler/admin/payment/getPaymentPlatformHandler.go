// huma:migrated
package payment

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/payment"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetPaymentPlatformOutput struct {
	Body *types.PlatformResponse
}

func GetPaymentPlatformHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetPaymentPlatformOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetPaymentPlatformOutput, error) {
		l := payment.NewGetPaymentPlatformLogic(ctx, svcCtx)
		resp, err := l.GetPaymentPlatform()
		if err != nil {
			return nil, err
		}
		return &GetPaymentPlatformOutput{Body: resp}, nil
	}
}
