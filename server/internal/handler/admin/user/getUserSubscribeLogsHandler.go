// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetUserSubscribeLogsInput struct {
	types.GetUserSubscribeLogsRequest
}

type GetUserSubscribeLogsOutput struct {
	Body *types.GetUserSubscribeLogsResponse
}

func GetUserSubscribeLogsHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserSubscribeLogsInput) (*GetUserSubscribeLogsOutput, error) {
	return func(ctx context.Context, input *GetUserSubscribeLogsInput) (*GetUserSubscribeLogsOutput, error) {
		l := user.NewGetUserSubscribeLogsLogic(ctx, svcCtx)
		resp, err := l.GetUserSubscribeLogs(&input.GetUserSubscribeLogsRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserSubscribeLogsOutput{Body: resp}, nil
	}
}
