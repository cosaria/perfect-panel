// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetUserSubscribeTrafficLogsInput struct {
	types.GetUserSubscribeTrafficLogsRequest
}

type GetUserSubscribeTrafficLogsOutput struct {
	Body *types.GetUserSubscribeTrafficLogsResponse
}

func GetUserSubscribeTrafficLogsHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserSubscribeTrafficLogsInput) (*GetUserSubscribeTrafficLogsOutput, error) {
	return func(ctx context.Context, input *GetUserSubscribeTrafficLogsInput) (*GetUserSubscribeTrafficLogsOutput, error) {
		l := user.NewGetUserSubscribeTrafficLogsLogic(ctx, svcCtx)
		resp, err := l.GetUserSubscribeTrafficLogs(&input.GetUserSubscribeTrafficLogsRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserSubscribeTrafficLogsOutput{Body: resp}, nil
	}
}
