// huma:migrated
package order

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/order"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateOrderStatusInput struct {
	Body types.UpdateOrderStatusRequest
}

func UpdateOrderStatusHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateOrderStatusInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateOrderStatusInput) (*struct{}, error) {
		l := order.NewUpdateOrderStatusLogic(ctx, svcCtx)
		if err := l.UpdateOrderStatus(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
