// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryUserInfoOutput struct {
	Body *types.User
}

func QueryUserInfoHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QueryUserInfoOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryUserInfoOutput, error) {
		l := NewQueryUserInfoLogic(ctx, svcCtx)
		resp, err := l.QueryUserInfo()
		if err != nil {
			return nil, err
		}
		return &QueryUserInfoOutput{Body: resp}, nil
	}
}
