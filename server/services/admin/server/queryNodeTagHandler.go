// huma:migrated
package server

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryNodeTagOutput struct {
	Body *types.QueryNodeTagResponse
}

func QueryNodeTagHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QueryNodeTagOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryNodeTagOutput, error) {
		l := NewQueryNodeTagLogic(ctx, svcCtx)
		resp, err := l.QueryNodeTag()
		if err != nil {
			return nil, err
		}
		return &QueryNodeTagOutput{Body: resp}, nil
	}
}
