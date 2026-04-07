// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetUserSubscribeResetTrafficLogsInput struct {
	types.GetUserSubscribeResetTrafficLogsRequest
}

type GetUserSubscribeResetTrafficLogsOutput struct {
	Body *types.GetUserSubscribeResetTrafficLogsResponse
}

func GetUserSubscribeResetTrafficLogsHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserSubscribeResetTrafficLogsInput) (*GetUserSubscribeResetTrafficLogsOutput, error) {
	return func(ctx context.Context, input *GetUserSubscribeResetTrafficLogsInput) (*GetUserSubscribeResetTrafficLogsOutput, error) {
		l := NewGetUserSubscribeResetTrafficLogsLogic(ctx, svcCtx)
		resp, err := l.GetUserSubscribeResetTrafficLogs(&input.GetUserSubscribeResetTrafficLogsRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserSubscribeResetTrafficLogsOutput{Body: resp}, nil
	}
}
