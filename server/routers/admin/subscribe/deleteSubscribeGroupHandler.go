// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/services/admin/subscribe"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type DeleteSubscribeGroupInput struct {
	Body types.DeleteSubscribeGroupRequest
}

func DeleteSubscribeGroupHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteSubscribeGroupInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteSubscribeGroupInput) (*struct{}, error) {
		l := subscribe.NewDeleteSubscribeGroupLogic(ctx, svcCtx)
		if err := l.DeleteSubscribeGroup(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
