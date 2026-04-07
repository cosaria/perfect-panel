// huma:migrated
package marketing

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryQuotaTaskPreCountInput struct {
	Body types.QueryQuotaTaskPreCountRequest
}

type QueryQuotaTaskPreCountOutput struct {
	Body *types.QueryQuotaTaskPreCountResponse
}

func QueryQuotaTaskPreCountHandler(svcCtx *svc.ServiceContext) func(context.Context, *QueryQuotaTaskPreCountInput) (*QueryQuotaTaskPreCountOutput, error) {
	return func(ctx context.Context, input *QueryQuotaTaskPreCountInput) (*QueryQuotaTaskPreCountOutput, error) {
		l := NewQueryQuotaTaskPreCountLogic(ctx, svcCtx)
		resp, err := l.QueryQuotaTaskPreCount(&input.Body)
		if err != nil {
			return nil, err
		}
		return &QueryQuotaTaskPreCountOutput{Body: resp}, nil
	}
}
