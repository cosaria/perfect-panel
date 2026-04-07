// huma:migrated
package auth

import (
	"context"
	"github.com/perfect-panel/server/services/auth"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CheckUserInput struct {
	types.CheckUserRequest
}

type CheckUserOutput struct {
	Body *types.CheckUserResponse
}

func CheckUserHandler(svcCtx *svc.ServiceContext) func(context.Context, *CheckUserInput) (*CheckUserOutput, error) {
	return func(ctx context.Context, input *CheckUserInput) (*CheckUserOutput, error) {
		l := auth.NewCheckUserLogic(ctx, svcCtx)
		resp, err := l.CheckUser(&input.CheckUserRequest)
		if err != nil {
			return nil, err
		}
		return &CheckUserOutput{Body: resp}, nil
	}
}
