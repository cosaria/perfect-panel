// huma:migrated
package portal

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type PrePurchaseOrderInput struct {
	Body types.PrePurchaseOrderRequest
}

type PrePurchaseOrderOutput struct {
	Body *types.PrePurchaseOrderResponse
}

func PrePurchaseOrderHandler(svcCtx *svc.ServiceContext) func(context.Context, *PrePurchaseOrderInput) (*PrePurchaseOrderOutput, error) {
	return func(ctx context.Context, input *PrePurchaseOrderInput) (*PrePurchaseOrderOutput, error) {
		l := NewPrePurchaseOrderLogic(ctx, svcCtx)
		resp, err := l.PrePurchaseOrder(&input.Body)
		if err != nil {
			return nil, err
		}
		return &PrePurchaseOrderOutput{Body: resp}, nil
	}
}
