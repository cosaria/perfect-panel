// huma:migrated
package order

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetOrderListInput struct {
	types.GetOrderListRequest
}

type GetOrderListOutput struct {
	Body *types.GetOrderListResponse
}

func GetOrderListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetOrderListInput) (*GetOrderListOutput, error) {
	return func(ctx context.Context, input *GetOrderListInput) (*GetOrderListOutput, error) {
		l := NewGetOrderListLogic(ctx, svcCtx)
		resp, err := l.GetOrderList(&input.GetOrderListRequest)
		if err != nil {
			return nil, err
		}
		return &GetOrderListOutput{Body: resp}, nil
	}
}
