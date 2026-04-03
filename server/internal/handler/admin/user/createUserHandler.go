// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CreateUserInput struct {
	Body types.CreateUserRequest
}

func CreateUserHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateUserInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateUserInput) (*struct{}, error) {
		l := user.NewCreateUserLogic(ctx, svcCtx)
		if err := l.CreateUser(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
