// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateUserPasswordInput struct {
	Body types.UpdateUserPasswordRequest
}

func UpdateUserPasswordHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateUserPasswordInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserPasswordInput) (*struct{}, error) {
		l := user.NewUpdateUserPasswordLogic(ctx, svcCtx)
		if err := l.UpdateUserPassword(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
