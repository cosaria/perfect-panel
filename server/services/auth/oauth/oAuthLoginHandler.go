// huma:migrated
package oauth

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type OAuthLoginInput struct {
	Body types.OAthLoginRequest
}

type OAuthLoginOutput struct {
	Body *types.OAuthLoginResponse
}

func OAuthLoginHandler(svcCtx *svc.ServiceContext) func(context.Context, *OAuthLoginInput) (*OAuthLoginOutput, error) {
	return func(ctx context.Context, input *OAuthLoginInput) (*OAuthLoginOutput, error) {
		l := NewOAuthLoginLogic(ctx, svcCtx)
		resp, err := l.OAuthLogin(&input.Body)
		if err != nil {
			return nil, err
		}
		return &OAuthLoginOutput{Body: resp}, nil
	}
}
