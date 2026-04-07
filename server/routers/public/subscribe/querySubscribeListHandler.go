// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/services/user/subscribe"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QuerySubscribeListInput struct {
	types.QuerySubscribeListRequest
}

type QuerySubscribeListOutput struct {
	Body *types.QuerySubscribeListResponse
}

func QuerySubscribeListHandler(svcCtx *svc.ServiceContext) func(context.Context, *QuerySubscribeListInput) (*QuerySubscribeListOutput, error) {
	return func(ctx context.Context, input *QuerySubscribeListInput) (*QuerySubscribeListOutput, error) {
		l := subscribe.NewQuerySubscribeListLogic(ctx, svcCtx)
		resp, err := l.QuerySubscribeList(&input.QuerySubscribeListRequest)
		if err != nil {
			return nil, err
		}
		return &QuerySubscribeListOutput{Body: resp}, nil
	}
}
