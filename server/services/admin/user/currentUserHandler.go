// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CurrentUserOutput struct {
	Body *types.User
}

func CurrentUserHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*CurrentUserOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*CurrentUserOutput, error) {
		l := NewCurrentUserLogic(ctx, svcCtx)
		resp, err := l.CurrentUser()
		if err != nil {
			return nil, err
		}
		return &CurrentUserOutput{Body: resp}, nil
	}
}
