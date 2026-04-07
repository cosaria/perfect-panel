// huma:migrated
package order

import (
	"context"
	"github.com/perfect-panel/server/services/user/order"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type PurchaseInput struct {
	Body types.PurchaseOrderRequest
}

type PurchaseOutput struct {
	Body *types.PurchaseOrderResponse
}

func PurchaseHandler(svcCtx *svc.ServiceContext) func(context.Context, *PurchaseInput) (*PurchaseOutput, error) {
	return func(ctx context.Context, input *PurchaseInput) (*PurchaseOutput, error) {
		l := order.NewPurchaseLogic(ctx, svcCtx)
		resp, err := l.Purchase(&input.Body)
		if err != nil {
			return nil, err
		}
		return &PurchaseOutput{Body: resp}, nil
	}
}
