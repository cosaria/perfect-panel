// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/admin/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type ResetUserSubscribeTokenInput struct {
	Body types.ResetUserSubscribeTokenRequest
}

func ResetUserSubscribeTokenHandler(svcCtx *svc.ServiceContext) func(context.Context, *ResetUserSubscribeTokenInput) (*struct{}, error) {
	return func(ctx context.Context, input *ResetUserSubscribeTokenInput) (*struct{}, error) {
		l := user.NewResetUserSubscribeTokenLogic(ctx, svcCtx)
		if err := l.ResetUserSubscribeToken(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
