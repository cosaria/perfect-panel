// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/services/admin/subscribe"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type DeleteSubscribeInput struct {
	Body types.DeleteSubscribeRequest
}

func DeleteSubscribeHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteSubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteSubscribeInput) (*struct{}, error) {
		l := subscribe.NewDeleteSubscribeLogic(ctx, svcCtx)
		if err := l.DeleteSubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
