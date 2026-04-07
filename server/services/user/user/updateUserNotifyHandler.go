// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateUserNotifyInput struct {
	Body types.UpdateUserNotifyRequest
}

func UpdateUserNotifyHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateUserNotifyInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserNotifyInput) (*struct{}, error) {
		l := NewUpdateUserNotifyLogic(ctx, svcCtx)
		if err := l.UpdateUserNotify(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
