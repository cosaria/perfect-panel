// huma:migrated
package order

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/order"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type RechargeInput struct {
	Body types.RechargeOrderRequest
}

type RechargeOutput struct {
	Body *types.RechargeOrderResponse
}

func RechargeHandler(svcCtx *svc.ServiceContext) func(context.Context, *RechargeInput) (*RechargeOutput, error) {
	return func(ctx context.Context, input *RechargeInput) (*RechargeOutput, error) {
		l := order.NewRechargeLogic(ctx, svcCtx)
		resp, err := l.Recharge(&input.Body)
		if err != nil {
			return nil, err
		}
		return &RechargeOutput{Body: resp}, nil
	}
}
