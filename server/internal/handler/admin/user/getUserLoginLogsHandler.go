// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetUserLoginLogsInput struct {
	types.GetUserLoginLogsRequest
}

type GetUserLoginLogsOutput struct {
	Body *types.GetUserLoginLogsResponse
}

func GetUserLoginLogsHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserLoginLogsInput) (*GetUserLoginLogsOutput, error) {
	return func(ctx context.Context, input *GetUserLoginLogsInput) (*GetUserLoginLogsOutput, error) {
		l := user.NewGetUserLoginLogsLogic(ctx, svcCtx)
		resp, err := l.GetUserLoginLogs(&input.GetUserLoginLogsRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserLoginLogsOutput{Body: resp}, nil
	}
}
