// huma:migrated
package portal

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryPurchaseOrderInput struct {
	types.QueryPurchaseOrderRequest
}

type QueryPurchaseOrderOutput struct {
	Body *types.QueryPurchaseOrderResponse
}

func QueryPurchaseOrderHandler(svcCtx *svc.ServiceContext) func(context.Context, *QueryPurchaseOrderInput) (*QueryPurchaseOrderOutput, error) {
	return func(ctx context.Context, input *QueryPurchaseOrderInput) (*QueryPurchaseOrderOutput, error) {
		l := NewQueryPurchaseOrderLogic(ctx, svcCtx)
		resp, err := l.QueryPurchaseOrder(&input.QueryPurchaseOrderRequest)
		if err != nil {
			return nil, err
		}
		return &QueryPurchaseOrderOutput{Body: resp}, nil
	}
}
