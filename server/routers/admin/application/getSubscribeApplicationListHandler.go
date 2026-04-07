// huma:migrated
package application

import (
	"context"
	"github.com/perfect-panel/server/services/admin/application"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetSubscribeApplicationListInput struct {
	types.GetSubscribeApplicationListRequest
}

type GetSubscribeApplicationListOutput struct {
	Body *types.GetSubscribeApplicationListResponse
}

func GetSubscribeApplicationListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetSubscribeApplicationListInput) (*GetSubscribeApplicationListOutput, error) {
	return func(ctx context.Context, input *GetSubscribeApplicationListInput) (*GetSubscribeApplicationListOutput, error) {
		l := application.NewGetSubscribeApplicationListLogic(ctx, svcCtx)
		resp, err := l.GetSubscribeApplicationList(&input.GetSubscribeApplicationListRequest)
		if err != nil {
			return nil, err
		}
		return &GetSubscribeApplicationListOutput{Body: resp}, nil
	}
}
