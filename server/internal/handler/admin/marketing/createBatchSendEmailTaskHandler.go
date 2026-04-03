// huma:migrated
package marketing

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/marketing"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CreateBatchSendEmailTaskInput struct {
	Body types.CreateBatchSendEmailTaskRequest
}

func CreateBatchSendEmailTaskHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateBatchSendEmailTaskInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateBatchSendEmailTaskInput) (*struct{}, error) {
		l := marketing.NewCreateBatchSendEmailTaskLogic(ctx, svcCtx)
		if err := l.CreateBatchSendEmailTask(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
