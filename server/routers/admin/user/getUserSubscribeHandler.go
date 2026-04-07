// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/admin/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetUserSubscribeInput struct {
	types.GetUserSubscribeListRequest
}

type GetUserSubscribeOutput struct {
	Body *types.GetUserSubscribeListResponse
}

func GetUserSubscribeHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserSubscribeInput) (*GetUserSubscribeOutput, error) {
	return func(ctx context.Context, input *GetUserSubscribeInput) (*GetUserSubscribeOutput, error) {
		l := user.NewGetUserSubscribeLogic(ctx, svcCtx)
		resp, err := l.GetUserSubscribe(&input.GetUserSubscribeListRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserSubscribeOutput{Body: resp}, nil
	}
}
