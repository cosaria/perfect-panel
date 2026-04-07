// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateSubscribeGroupInput struct {
	Body types.UpdateSubscribeGroupRequest
}

func UpdateSubscribeGroupHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateSubscribeGroupInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateSubscribeGroupInput) (*struct{}, error) {
		l := NewUpdateSubscribeGroupLogic(ctx, svcCtx)
		if err := l.UpdateSubscribeGroup(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
