// huma:migrated
package subscribe

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/subscribe"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetSubscribeGroupListOutput struct {
	Body *types.GetSubscribeGroupListResponse
}

func GetSubscribeGroupListHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetSubscribeGroupListOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetSubscribeGroupListOutput, error) {
		l := subscribe.NewGetSubscribeGroupListLogic(ctx, svcCtx)
		resp, err := l.GetSubscribeGroupList()
		if err != nil {
			return nil, err
		}
		return &GetSubscribeGroupListOutput{Body: resp}, nil
	}
}
