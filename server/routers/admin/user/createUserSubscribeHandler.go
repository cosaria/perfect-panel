// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/admin/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CreateUserSubscribeInput struct {
	Body types.CreateUserSubscribeRequest
}

func CreateUserSubscribeHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateUserSubscribeInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateUserSubscribeInput) (*struct{}, error) {
		l := user.NewCreateUserSubscribeLogic(ctx, svcCtx)
		if err := l.CreateUserSubscribe(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
