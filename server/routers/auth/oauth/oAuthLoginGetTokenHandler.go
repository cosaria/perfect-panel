// huma:migrated
package oauth

import (
	"context"

	"github.com/perfect-panel/server/services/auth/oauth"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type OAuthLoginGetTokenInput struct {
	Body      types.OAuthLoginGetTokenRequest
	IP        string `header:"X-Original-Forwarded-For" required:"false" doc:"Client IP from proxy"`
	UserAgent string `header:"User-Agent" required:"false" doc:"User agent string"`
}

type OAuthLoginGetTokenOutput struct {
	Body *types.LoginResponse
}

func OAuthLoginGetTokenHandler(svcCtx *svc.ServiceContext) func(context.Context, *OAuthLoginGetTokenInput) (*OAuthLoginGetTokenOutput, error) {
	return func(ctx context.Context, input *OAuthLoginGetTokenInput) (*OAuthLoginGetTokenOutput, error) {
		l := oauth.NewOAuthLoginGetTokenLogic(ctx, svcCtx)
		resp, err := l.OAuthLoginGetToken(&input.Body, input.IP, input.UserAgent)
		if err != nil {
			return nil, err
		}
		return &OAuthLoginGetTokenOutput{Body: resp}, nil
	}
}
