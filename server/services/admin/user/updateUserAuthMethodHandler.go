// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateUserAuthMethodInput struct {
	Body types.UpdateUserAuthMethodRequest
}

func UpdateUserAuthMethodHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateUserAuthMethodInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserAuthMethodInput) (*struct{}, error) {
		l := NewUpdateUserAuthMethodLogic(ctx, svcCtx)
		if err := l.UpdateUserAuthMethod(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
