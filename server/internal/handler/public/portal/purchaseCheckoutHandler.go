// huma:migrated
package portal

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/portal"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type PurchaseCheckoutInput struct {
	Body types.CheckoutOrderRequest
}

type PurchaseCheckoutOutput struct {
	Body *types.CheckoutOrderResponse
}

func PurchaseCheckoutHandler(svcCtx *svc.ServiceContext) func(context.Context, *PurchaseCheckoutInput) (*PurchaseCheckoutOutput, error) {
	return func(ctx context.Context, input *PurchaseCheckoutInput) (*PurchaseCheckoutOutput, error) {
		l := portal.NewPurchaseCheckoutLogic(ctx, svcCtx)
		resp, err := l.PurchaseCheckout(&input.Body)
		if err != nil {
			return nil, err
		}
		return &PurchaseCheckoutOutput{Body: resp}, nil
	}
}
