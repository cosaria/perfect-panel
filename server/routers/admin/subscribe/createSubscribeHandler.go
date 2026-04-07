// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/services/admin/subscribe"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CreateSubscribeInput struct {
	Body types.CreateSubscribeRequest
}

func CreateSubscribeHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateSubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateSubscribeInput) (*struct{}, error) {
		l := subscribe.NewCreateSubscribeLogic(ctx, svcCtx)
		if err := l.CreateSubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
