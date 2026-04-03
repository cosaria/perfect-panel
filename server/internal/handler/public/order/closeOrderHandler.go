// huma:migrated
package order

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/order"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CloseOrderInput struct {
	Body types.CloseOrderRequest
}

func CloseOrderHandler(svcCtx *svc.ServiceContext) func(context.Context, *CloseOrderInput) (*struct{}, error) {
	return func(ctx context.Context, input *CloseOrderInput) (*struct{}, error) {
		l := order.NewCloseOrderLogic(ctx, svcCtx)
		if err := l.CloseOrder(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
