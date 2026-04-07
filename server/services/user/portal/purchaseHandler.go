// huma:migrated
package portal

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type PurchaseInput struct {
	Body types.PortalPurchaseRequest
}

type PurchaseOutput struct {
	Body *types.PortalPurchaseResponse
}

func PurchaseHandler(svcCtx *svc.ServiceContext) func(context.Context, *PurchaseInput) (*PurchaseOutput, error) {
	return func(ctx context.Context, input *PurchaseInput) (*PurchaseOutput, error) {
		l := NewPurchaseLogic(ctx, svcCtx)
		resp, err := l.Purchase(&input.Body)
		if err != nil {
			return nil, err
		}
		return &PurchaseOutput{Body: resp}, nil
	}
}
