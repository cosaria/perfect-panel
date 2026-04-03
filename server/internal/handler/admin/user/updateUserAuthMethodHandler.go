// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateUserAuthMethodInput struct {
	Body types.UpdateUserAuthMethodRequest
}

func UpdateUserAuthMethodHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateUserAuthMethodInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserAuthMethodInput) (*struct{}, error) {
		l := user.NewUpdateUserAuthMethodLogic(ctx, svcCtx)
		if err := l.UpdateUserAuthMethod(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
