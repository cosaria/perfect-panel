// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type QueryUserCommissionLogInput struct {
	types.QueryUserCommissionLogListRequest
}

type QueryUserCommissionLogOutput struct {
	Body *types.QueryUserCommissionLogListResponse
}

func QueryUserCommissionLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *QueryUserCommissionLogInput) (*QueryUserCommissionLogOutput, error) {
	return func(ctx context.Context, input *QueryUserCommissionLogInput) (*QueryUserCommissionLogOutput, error) {
		l := user.NewQueryUserCommissionLogLogic(ctx, svcCtx)
		resp, err := l.QueryUserCommissionLog(&input.QueryUserCommissionLogListRequest)
		if err != nil {
			return nil, err
		}
		return &QueryUserCommissionLogOutput{Body: resp}, nil
	}
}
