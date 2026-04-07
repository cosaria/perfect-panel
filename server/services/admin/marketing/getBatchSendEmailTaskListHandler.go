// huma:migrated
package marketing

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetBatchSendEmailTaskListInput struct {
	Body types.GetBatchSendEmailTaskListRequest
}

type GetBatchSendEmailTaskListOutput struct {
	Body *types.GetBatchSendEmailTaskListResponse
}

func GetBatchSendEmailTaskListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetBatchSendEmailTaskListInput) (*GetBatchSendEmailTaskListOutput, error) {
	return func(ctx context.Context, input *GetBatchSendEmailTaskListInput) (*GetBatchSendEmailTaskListOutput, error) {
		l := NewGetBatchSendEmailTaskListLogic(ctx, svcCtx)
		resp, err := l.GetBatchSendEmailTaskList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetBatchSendEmailTaskListOutput{Body: resp}, nil
	}
}
