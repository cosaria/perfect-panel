// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/user/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryUserSubscribeOutput struct {
	Body *types.QueryUserSubscribeListResponse
}

func QueryUserSubscribeHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QueryUserSubscribeOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryUserSubscribeOutput, error) {
		l := user.NewQueryUserSubscribeLogic(ctx, svcCtx)
		resp, err := l.QueryUserSubscribe()
		if err != nil {
			return nil, err
		}
		return &QueryUserSubscribeOutput{Body: resp}, nil
	}
}
