// huma:migrated
package auth

import (
	"context"
	"github.com/perfect-panel/server/services/auth"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CheckUserTelephoneInput struct {
	types.TelephoneCheckUserRequest
}

type CheckUserTelephoneOutput struct {
	Body *types.TelephoneCheckUserResponse
}

func CheckUserTelephoneHandler(svcCtx *svc.ServiceContext) func(context.Context, *CheckUserTelephoneInput) (*CheckUserTelephoneOutput, error) {
	return func(ctx context.Context, input *CheckUserTelephoneInput) (*CheckUserTelephoneOutput, error) {
		l := auth.NewCheckUserTelephoneLogic(ctx, svcCtx)
		resp, err := l.CheckUserTelephone(&input.TelephoneCheckUserRequest)
		if err != nil {
			return nil, err
		}
		return &CheckUserTelephoneOutput{Body: resp}, nil
	}
}
