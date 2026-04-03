// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/subscribe"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type BatchDeleteSubscribeGroupInput struct {
	Body types.BatchDeleteSubscribeGroupRequest
}

func BatchDeleteSubscribeGroupHandler(svcCtx *svc.ServiceContext) func(context.Context, *BatchDeleteSubscribeGroupInput) (*struct{}, error) {
	return func(ctx context.Context, input *BatchDeleteSubscribeGroupInput) (*struct{}, error) {
		l := subscribe.NewBatchDeleteSubscribeGroupLogic(ctx, svcCtx)
		if err := l.BatchDeleteSubscribeGroup(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
