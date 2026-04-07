// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type VerifyEmailInput struct {
	Body types.VerifyEmailRequest
}

func VerifyEmailHandler(svcCtx *svc.ServiceContext) func(context.Context, *VerifyEmailInput) (*struct{}, error) {
	return func(ctx context.Context, input *VerifyEmailInput) (*struct{}, error) {
		l := NewVerifyEmailLogic(ctx, svcCtx)
		if err := l.VerifyEmail(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
