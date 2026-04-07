// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/admin/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type DeleteUserDeviceInput struct {
	Body types.DeleteUserDeivceRequest
}

func DeleteUserDeviceHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteUserDeviceInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteUserDeviceInput) (*struct{}, error) {
		l := user.NewDeleteUserDeviceLogic(ctx, svcCtx)
		if err := l.DeleteUserDevice(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
