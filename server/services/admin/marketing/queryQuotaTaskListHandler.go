// huma:migrated
package marketing

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryQuotaTaskListInput struct {
	Body types.QueryQuotaTaskListRequest
}

type QueryQuotaTaskListOutput struct {
	Body *types.QueryQuotaTaskListResponse
}

func QueryQuotaTaskListHandler(svcCtx *svc.ServiceContext) func(context.Context, *QueryQuotaTaskListInput) (*QueryQuotaTaskListOutput, error) {
	return func(ctx context.Context, input *QueryQuotaTaskListInput) (*QueryQuotaTaskListOutput, error) {
		l := NewQueryQuotaTaskListLogic(ctx, svcCtx)
		resp, err := l.QueryQuotaTaskList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &QueryQuotaTaskListOutput{Body: resp}, nil
	}
}
