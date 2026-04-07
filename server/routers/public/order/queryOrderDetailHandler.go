// huma:migrated
package order

import (
	"context"
	"github.com/perfect-panel/server/services/user/order"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryOrderDetailInput struct {
	types.QueryOrderDetailRequest
}

type QueryOrderDetailOutput struct {
	Body *types.OrderDetail
}

func QueryOrderDetailHandler(svcCtx *svc.ServiceContext) func(context.Context, *QueryOrderDetailInput) (*QueryOrderDetailOutput, error) {
	return func(ctx context.Context, input *QueryOrderDetailInput) (*QueryOrderDetailOutput, error) {
		l := order.NewQueryOrderDetailLogic(ctx, svcCtx)
		resp, err := l.QueryOrderDetail(&input.QueryOrderDetailRequest)
		if err != nil {
			return nil, err
		}
		return &QueryOrderDetailOutput{Body: resp}, nil
	}
}
