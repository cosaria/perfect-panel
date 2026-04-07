// huma:migrated
package tool

import (
	"context"
	"github.com/perfect-panel/server/services/admin/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryIPLocationInput struct {
	types.QueryIPLocationRequest
}

type QueryIPLocationOutput struct {
	Body *types.QueryIPLocationResponse
}

func QueryIPLocationHandler(svcCtx *svc.ServiceContext) func(context.Context, *QueryIPLocationInput) (*QueryIPLocationOutput, error) {
	return func(ctx context.Context, input *QueryIPLocationInput) (*QueryIPLocationOutput, error) {
		l := tool.NewQueryIPLocationLogic(ctx, svcCtx)
		resp, err := l.QueryIPLocation(&input.QueryIPLocationRequest)
		if err != nil {
			return nil, err
		}
		return &QueryIPLocationOutput{Body: resp}, nil
	}
}
