// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/subscribe"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type QueryUserSubscribeNodeListOutput struct {
	Body *types.QueryUserSubscribeNodeListResponse
}

func QueryUserSubscribeNodeListHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QueryUserSubscribeNodeListOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryUserSubscribeNodeListOutput, error) {
		l := subscribe.NewQueryUserSubscribeNodeListLogic(ctx, svcCtx)
		resp, err := l.QueryUserSubscribeNodeList()
		if err != nil {
			return nil, err
		}
		return &QueryUserSubscribeNodeListOutput{Body: resp}, nil
	}
}
