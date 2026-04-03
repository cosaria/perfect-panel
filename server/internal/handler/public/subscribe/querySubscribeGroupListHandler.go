// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/subscribe"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type QuerySubscribeGroupListOutput struct {
	Body *types.QuerySubscribeGroupListResponse
}

func QuerySubscribeGroupListHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QuerySubscribeGroupListOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QuerySubscribeGroupListOutput, error) {
		l := subscribe.NewQuerySubscribeGroupListLogic(ctx, svcCtx)
		resp, err := l.QuerySubscribeGroupList()
		if err != nil {
			return nil, err
		}
		return &QuerySubscribeGroupListOutput{Body: resp}, nil
	}
}
