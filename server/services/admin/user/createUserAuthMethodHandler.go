// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CreateUserAuthMethodInput struct {
	Body types.CreateUserAuthMethodRequest
}

func CreateUserAuthMethodHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateUserAuthMethodInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateUserAuthMethodInput) (*struct{}, error) {
		l := NewCreateUserAuthMethodLogic(ctx, svcCtx)
		if err := l.CreateUserAuthMethod(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
