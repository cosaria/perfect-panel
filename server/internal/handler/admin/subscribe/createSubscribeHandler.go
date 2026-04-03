// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/subscribe"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
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
