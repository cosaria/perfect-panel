// huma:migrated
package marketing

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CreateBatchSendEmailTaskInput struct {
	Body types.CreateBatchSendEmailTaskRequest
}

func CreateBatchSendEmailTaskHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateBatchSendEmailTaskInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateBatchSendEmailTaskInput) (*struct{}, error) {
		l := NewCreateBatchSendEmailTaskLogic(ctx, svcCtx)
		if err := l.CreateBatchSendEmailTask(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
