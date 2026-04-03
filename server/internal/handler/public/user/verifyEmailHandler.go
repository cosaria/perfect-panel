// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type VerifyEmailInput struct {
	Body types.VerifyEmailRequest
}

func VerifyEmailHandler(svcCtx *svc.ServiceContext) func(context.Context, *VerifyEmailInput) (*struct{}, error) {
	return func(ctx context.Context, input *VerifyEmailInput) (*struct{}, error) {
		l := user.NewVerifyEmailLogic(ctx, svcCtx)
		if err := l.VerifyEmail(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
