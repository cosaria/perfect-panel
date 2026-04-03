// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type QueryUserBalanceLogOutput struct {
	Body *types.QueryUserBalanceLogListResponse
}

func QueryUserBalanceLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QueryUserBalanceLogOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryUserBalanceLogOutput, error) {
		l := user.NewQueryUserBalanceLogLogic(ctx, svcCtx)
		resp, err := l.QueryUserBalanceLog()
		if err != nil {
			return nil, err
		}
		return &QueryUserBalanceLogOutput{Body: resp}, nil
	}
}
