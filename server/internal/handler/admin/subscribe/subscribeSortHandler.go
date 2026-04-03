// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/subscribe"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type SubscribeSortInput struct {
	Body types.SubscribeSortRequest
}

func SubscribeSortHandler(svcCtx *svc.ServiceContext) func(context.Context, *SubscribeSortInput) (*struct{}, error) {
	return func(ctx context.Context, input *SubscribeSortInput) (*struct{}, error) {
		l := subscribe.NewSubscribeSortLogic(ctx, svcCtx)
		if err := l.SubscribeSort(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
