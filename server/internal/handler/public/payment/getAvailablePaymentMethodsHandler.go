// huma:migrated
package payment

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/payment"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetAvailablePaymentMethodsOutput struct {
	Body *types.GetAvailablePaymentMethodsResponse
}

func GetAvailablePaymentMethodsHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetAvailablePaymentMethodsOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetAvailablePaymentMethodsOutput, error) {
		l := payment.NewGetAvailablePaymentMethodsLogic(ctx, svcCtx)
		resp, err := l.GetAvailablePaymentMethods()
		if err != nil {
			return nil, err
		}
		return &GetAvailablePaymentMethodsOutput{Body: resp}, nil
	}
}
