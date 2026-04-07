// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/services/admin/subscribe"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetSubscribeListInput struct {
	types.GetSubscribeListRequest
}

type GetSubscribeListOutput struct {
	Body *types.GetSubscribeListResponse
}

func GetSubscribeListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetSubscribeListInput) (*GetSubscribeListOutput, error) {
	return func(ctx context.Context, input *GetSubscribeListInput) (*GetSubscribeListOutput, error) {
		l := subscribe.NewGetSubscribeListLogic(ctx, svcCtx)
		resp, err := l.GetSubscribeList(&input.GetSubscribeListRequest)
		if err != nil {
			return nil, err
		}
		return &GetSubscribeListOutput{Body: resp}, nil
	}
}
