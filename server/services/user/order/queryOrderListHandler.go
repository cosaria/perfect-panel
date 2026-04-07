// huma:migrated
package order

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryOrderListInput struct {
	types.QueryOrderListRequest
}

type QueryOrderListOutput struct {
	Body *types.QueryOrderListResponse
}

func QueryOrderListHandler(svcCtx *svc.ServiceContext) func(context.Context, *QueryOrderListInput) (*QueryOrderListOutput, error) {
	return func(ctx context.Context, input *QueryOrderListInput) (*QueryOrderListOutput, error) {
		l := NewQueryOrderListLogic(ctx, svcCtx)
		resp, err := l.QueryOrderList(&input.QueryOrderListRequest)
		if err != nil {
			return nil, err
		}
		return &QueryOrderListOutput{Body: resp}, nil
	}
}
