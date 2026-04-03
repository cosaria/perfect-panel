// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CreateUserAuthMethodInput struct {
	Body types.CreateUserAuthMethodRequest
}

func CreateUserAuthMethodHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateUserAuthMethodInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateUserAuthMethodInput) (*struct{}, error) {
		l := user.NewCreateUserAuthMethodLogic(ctx, svcCtx)
		if err := l.CreateUserAuthMethod(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
