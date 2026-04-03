// huma:migrated
package marketing

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/marketing"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetBatchSendEmailTaskListInput struct {
	Body types.GetBatchSendEmailTaskListRequest
}

type GetBatchSendEmailTaskListOutput struct {
	Body *types.GetBatchSendEmailTaskListResponse
}

func GetBatchSendEmailTaskListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetBatchSendEmailTaskListInput) (*GetBatchSendEmailTaskListOutput, error) {
	return func(ctx context.Context, input *GetBatchSendEmailTaskListInput) (*GetBatchSendEmailTaskListOutput, error) {
		l := marketing.NewGetBatchSendEmailTaskListLogic(ctx, svcCtx)
		resp, err := l.GetBatchSendEmailTaskList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetBatchSendEmailTaskListOutput{Body: resp}, nil
	}
}
