// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/services/admin/subscribe"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateSubscribeInput struct {
	Body types.UpdateSubscribeRequest
}

func UpdateSubscribeHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateSubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateSubscribeInput) (*struct{}, error) {
		l := subscribe.NewUpdateSubscribeLogic(ctx, svcCtx)
		if err := l.UpdateSubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
