// huma:migrated
package order

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CloseOrderInput struct {
	Body types.CloseOrderRequest
}

func CloseOrderHandler(svcCtx *svc.ServiceContext) func(context.Context, *CloseOrderInput) (*struct{}, error) {
	return func(ctx context.Context, input *CloseOrderInput) (*struct{}, error) {
		l := NewCloseOrderLogic(ctx, svcCtx)
		if err := l.CloseOrder(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
