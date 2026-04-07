// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/services/user/user"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type BindOAuthInput struct {
	Body types.BindOAuthRequest
}

type BindOAuthOutput struct {
	Body *types.BindOAuthResponse
}

func BindOAuthHandler(svcCtx *svc.ServiceContext) func(context.Context, *BindOAuthInput) (*BindOAuthOutput, error) {
	return func(ctx context.Context, input *BindOAuthInput) (*BindOAuthOutput, error) {
		l := user.NewBindOAuthLogic(ctx, svcCtx)
		resp, err := l.BindOAuth(&input.Body)
		if err != nil {
			return nil, err
		}
		return &BindOAuthOutput{Body: resp}, nil
	}
}
