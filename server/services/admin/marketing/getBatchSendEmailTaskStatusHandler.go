// huma:migrated
package marketing

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetBatchSendEmailTaskStatusInput struct {
	Body types.GetBatchSendEmailTaskStatusRequest
}

type GetBatchSendEmailTaskStatusOutput struct {
	Body *types.GetBatchSendEmailTaskStatusResponse
}

func GetBatchSendEmailTaskStatusHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetBatchSendEmailTaskStatusInput) (*GetBatchSendEmailTaskStatusOutput, error) {
	return func(ctx context.Context, input *GetBatchSendEmailTaskStatusInput) (*GetBatchSendEmailTaskStatusOutput, error) {
		l := NewGetBatchSendEmailTaskStatusLogic(ctx, svcCtx)
		resp, err := l.GetBatchSendEmailTaskStatus(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetBatchSendEmailTaskStatusOutput{Body: resp}, nil
	}
}
