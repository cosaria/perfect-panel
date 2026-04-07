// huma:migrated
package order

import (
	"context"
	"github.com/perfect-panel/server/services/admin/order"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CreateOrderInput struct {
	Body types.CreateOrderRequest
}

func CreateOrderHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateOrderInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateOrderInput) (*struct{}, error) {
		l := order.NewCreateOrderLogic(ctx, svcCtx)
		if err := l.CreateOrder(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
