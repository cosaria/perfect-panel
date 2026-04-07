// huma:migrated
package payment

import (
	"context"
	"github.com/perfect-panel/server/services/admin/payment"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type DeletePaymentMethodInput struct {
	Body types.DeletePaymentMethodRequest
}

func DeletePaymentMethodHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeletePaymentMethodInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeletePaymentMethodInput) (*struct{}, error) {
		l := payment.NewDeletePaymentMethodLogic(ctx, svcCtx)
		if err := l.DeletePaymentMethod(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
