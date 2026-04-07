// huma:migrated
package marketing

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type StopBatchSendEmailTaskInput struct {
	Body types.StopBatchSendEmailTaskRequest
}

func StopBatchSendEmailTaskHandler(svcCtx *svc.ServiceContext) func(context.Context, *StopBatchSendEmailTaskInput) (*struct{}, error) {
	return func(ctx context.Context, input *StopBatchSendEmailTaskInput) (*struct{}, error) {
		l := NewStopBatchSendEmailTaskLogic(ctx, svcCtx)
		if err := l.StopBatchSendEmailTask(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
