// huma:migrated
package order

import (
	"context"
	"github.com/perfect-panel/server/services/user/order"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type PreCreateOrderInput struct {
	Body types.PurchaseOrderRequest
}

type PreCreateOrderOutput struct {
	Body *types.PreOrderResponse
}

func PreCreateOrderHandler(svcCtx *svc.ServiceContext) func(context.Context, *PreCreateOrderInput) (*PreCreateOrderOutput, error) {
	return func(ctx context.Context, input *PreCreateOrderInput) (*PreCreateOrderOutput, error) {
		l := order.NewPreCreateOrderLogic(ctx, svcCtx)
		resp, err := l.PreCreateOrder(&input.Body)
		if err != nil {
			return nil, err
		}
		return &PreCreateOrderOutput{Body: resp}, nil
	}
}
