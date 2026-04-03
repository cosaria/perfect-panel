// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateUserNotifySettingInput struct {
	Body types.UpdateUserNotifySettingRequest
}

func UpdateUserNotifySettingHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateUserNotifySettingInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserNotifySettingInput) (*struct{}, error) {
		l := user.NewUpdateUserNotifySettingLogic(ctx, svcCtx)
		if err := l.UpdateUserNotifySetting(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
