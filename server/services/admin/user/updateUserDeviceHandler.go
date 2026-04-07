// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateUserDeviceInput struct {
	Body types.UserDevice
}

func UpdateUserDeviceHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateUserDeviceInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserDeviceInput) (*struct{}, error) {
		l := NewUpdateUserDeviceLogic(ctx, svcCtx)
		if err := l.UpdateUserDevice(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
