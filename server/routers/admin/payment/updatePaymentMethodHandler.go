// huma:migrated
package payment

import (
	"context"
	"github.com/perfect-panel/server/services/admin/payment"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdatePaymentMethodInput struct {
	Body types.UpdatePaymentMethodRequest
}

type UpdatePaymentMethodOutput struct {
	Body *types.PaymentConfig
}

func UpdatePaymentMethodHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdatePaymentMethodInput) (*UpdatePaymentMethodOutput, error) {
	return func(ctx context.Context, input *UpdatePaymentMethodInput) (*UpdatePaymentMethodOutput, error) {
		l := payment.NewUpdatePaymentMethodLogic(ctx, svcCtx)
		resp, err := l.UpdatePaymentMethod(&input.Body)
		if err != nil {
			return nil, err
		}
		return &UpdatePaymentMethodOutput{Body: resp}, nil
	}
}
