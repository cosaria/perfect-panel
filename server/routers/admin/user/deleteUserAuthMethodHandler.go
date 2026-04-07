// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/admin/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type DeleteUserAuthMethodInput struct {
	Body types.DeleteUserAuthMethodRequest
}

func DeleteUserAuthMethodHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteUserAuthMethodInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteUserAuthMethodInput) (*struct{}, error) {
		l := user.NewDeleteUserAuthMethodLogic(ctx, svcCtx)
		if err := l.DeleteUserAuthMethod(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
