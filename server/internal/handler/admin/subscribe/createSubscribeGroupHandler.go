// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/subscribe"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CreateSubscribeGroupInput struct {
	Body types.CreateSubscribeGroupRequest
}

func CreateSubscribeGroupHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateSubscribeGroupInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateSubscribeGroupInput) (*struct{}, error) {
		l := subscribe.NewCreateSubscribeGroupLogic(ctx, svcCtx)
		if err := l.CreateSubscribeGroup(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
