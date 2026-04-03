// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateBindEmailInput struct {
	Body types.UpdateBindEmailRequest
}

func UpdateBindEmailHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateBindEmailInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateBindEmailInput) (*struct{}, error) {
		l := user.NewUpdateBindEmailLogic(ctx, svcCtx)
		if err := l.UpdateBindEmail(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
