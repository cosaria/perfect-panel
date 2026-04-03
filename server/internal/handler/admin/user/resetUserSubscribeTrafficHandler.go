// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type ResetUserSubscribeTrafficInput struct {
	Body types.ResetUserSubscribeTrafficRequest
}

func ResetUserSubscribeTrafficHandler(svcCtx *svc.ServiceContext) func(context.Context, *ResetUserSubscribeTrafficInput) (*struct{}, error) {
	return func(ctx context.Context, input *ResetUserSubscribeTrafficInput) (*struct{}, error) {
		l := user.NewResetUserSubscribeTrafficLogic(ctx, svcCtx)
		if err := l.ResetUserSubscribeTraffic(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
