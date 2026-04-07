// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateUserPasswordInput struct {
	Body types.UpdateUserPasswordRequest
}

func UpdateUserPasswordHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateUserPasswordInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserPasswordInput) (*struct{}, error) {
		l := NewUpdateUserPasswordLogic(ctx, svcCtx)
		if err := l.UpdateUserPassword(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
