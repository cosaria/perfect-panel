// huma:migrated
package payment

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/payment"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CreatePaymentMethodInput struct {
	Body types.CreatePaymentMethodRequest
}

type CreatePaymentMethodOutput struct {
	Body *types.PaymentConfig
}

func CreatePaymentMethodHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreatePaymentMethodInput) (*CreatePaymentMethodOutput, error) {
	return func(ctx context.Context, input *CreatePaymentMethodInput) (*CreatePaymentMethodOutput, error) {
		l := payment.NewCreatePaymentMethodLogic(ctx, svcCtx)
		resp, err := l.CreatePaymentMethod(&input.Body)
		if err != nil {
			return nil, err
		}
		return &CreatePaymentMethodOutput{Body: resp}, nil
	}
}
