// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CurrentUserOutput struct {
	Body *types.User
}

func CurrentUserHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*CurrentUserOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*CurrentUserOutput, error) {
		l := user.NewCurrentUserLogic(ctx, svcCtx)
		resp, err := l.CurrentUser()
		if err != nil {
			return nil, err
		}
		return &CurrentUserOutput{Body: resp}, nil
	}
}
