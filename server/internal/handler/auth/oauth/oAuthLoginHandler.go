// huma:migrated
package oauth

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/auth/oauth"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type OAuthLoginInput struct {
	Body types.OAthLoginRequest
}

type OAuthLoginOutput struct {
	Body *types.OAuthLoginResponse
}

func OAuthLoginHandler(svcCtx *svc.ServiceContext) func(context.Context, *OAuthLoginInput) (*OAuthLoginOutput, error) {
	return func(ctx context.Context, input *OAuthLoginInput) (*OAuthLoginOutput, error) {
		l := oauth.NewOAuthLoginLogic(ctx, svcCtx)
		resp, err := l.OAuthLogin(&input.Body)
		if err != nil {
			return nil, err
		}
		return &OAuthLoginOutput{Body: resp}, nil
	}
}
