// huma:migrated
package order

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type ResetTrafficInput struct {
	Body types.ResetTrafficOrderRequest
}

type ResetTrafficOutput struct {
	Body *types.ResetTrafficOrderResponse
}

func ResetTrafficHandler(svcCtx *svc.ServiceContext) func(context.Context, *ResetTrafficInput) (*ResetTrafficOutput, error) {
	return func(ctx context.Context, input *ResetTrafficInput) (*ResetTrafficOutput, error) {
		l := NewResetTrafficLogic(ctx, svcCtx)
		resp, err := l.ResetTraffic(&input.Body)
		if err != nil {
			return nil, err
		}
		return &ResetTrafficOutput{Body: resp}, nil
	}
}
