// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/user/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryWithdrawalLogInput struct {
	types.QueryWithdrawalLogListRequest
}

type QueryWithdrawalLogOutput struct {
	Body *types.QueryWithdrawalLogListResponse
}

func QueryWithdrawalLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *QueryWithdrawalLogInput) (*QueryWithdrawalLogOutput, error) {
	return func(ctx context.Context, input *QueryWithdrawalLogInput) (*QueryWithdrawalLogOutput, error) {
		l := user.NewQueryWithdrawalLogLogic(ctx, svcCtx)
		resp, err := l.QueryWithdrawalLog(&input.QueryWithdrawalLogListRequest)
		if err != nil {
			return nil, err
		}
		return &QueryWithdrawalLogOutput{Body: resp}, nil
	}
}
