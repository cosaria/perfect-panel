// huma:migrated
package order

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type RenewalInput struct {
	Body types.RenewalOrderRequest
}

type RenewalOutput struct {
	Body *types.RenewalOrderResponse
}

func RenewalHandler(svcCtx *svc.ServiceContext) func(context.Context, *RenewalInput) (*RenewalOutput, error) {
	return func(ctx context.Context, input *RenewalInput) (*RenewalOutput, error) {
		l := NewRenewalLogic(ctx, svcCtx)
		resp, err := l.Renewal(&input.Body)
		if err != nil {
			return nil, err
		}
		return &RenewalOutput{Body: resp}, nil
	}
}
