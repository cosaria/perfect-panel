// huma:migrated
package marketing

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/marketing"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type StopBatchSendEmailTaskInput struct {
	Body types.StopBatchSendEmailTaskRequest
}

func StopBatchSendEmailTaskHandler(svcCtx *svc.ServiceContext) func(context.Context, *StopBatchSendEmailTaskInput) (*struct{}, error) {
	return func(ctx context.Context, input *StopBatchSendEmailTaskInput) (*struct{}, error) {
		l := marketing.NewStopBatchSendEmailTaskLogic(ctx, svcCtx)
		if err := l.StopBatchSendEmailTask(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
