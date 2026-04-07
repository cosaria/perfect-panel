// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type DeleteUserInput struct {
	Body types.GetDetailRequest
}

func DeleteUserHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteUserInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteUserInput) (*struct{}, error) {
		l := NewDeleteUserLogic(ctx, svcCtx)
		if err := l.DeleteUser(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
