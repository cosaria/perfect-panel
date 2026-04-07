// huma:migrated
package marketing

import (
	"context"
	"github.com/perfect-panel/server/services/admin/marketing"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryQuotaTaskStatusInput struct {
	Body types.QueryQuotaTaskStatusRequest
}

type QueryQuotaTaskStatusOutput struct {
	Body *types.QueryQuotaTaskStatusResponse
}

func QueryQuotaTaskStatusHandler(svcCtx *svc.ServiceContext) func(context.Context, *QueryQuotaTaskStatusInput) (*QueryQuotaTaskStatusOutput, error) {
	return func(ctx context.Context, input *QueryQuotaTaskStatusInput) (*QueryQuotaTaskStatusOutput, error) {
		l := marketing.NewQueryQuotaTaskStatusLogic(ctx, svcCtx)
		resp, err := l.QueryQuotaTaskStatus(&input.Body)
		if err != nil {
			return nil, err
		}
		return &QueryQuotaTaskStatusOutput{Body: resp}, nil
	}
}
