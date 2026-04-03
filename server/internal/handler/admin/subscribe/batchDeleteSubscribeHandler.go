// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/subscribe"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type BatchDeleteSubscribeInput struct {
	Body types.BatchDeleteSubscribeRequest
}

func BatchDeleteSubscribeHandler(svcCtx *svc.ServiceContext) func(context.Context, *BatchDeleteSubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *BatchDeleteSubscribeInput) (*struct{}, error) {
		l := subscribe.NewBatchDeleteSubscribeLogic(ctx, svcCtx)
		if err := l.BatchDeleteSubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
