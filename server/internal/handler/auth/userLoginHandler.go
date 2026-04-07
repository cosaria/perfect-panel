// huma:migrated
package auth

import (
	"context"
	"time"

	"github.com/perfect-panel/server/internal/logic/auth"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
	"github.com/perfect-panel/server/pkg/turnstile"
	"github.com/perfect-panel/server/pkg/xerr"
	"github.com/pkg/errors"
)

type UserLoginInput struct {
	Body      types.UserLoginRequest
	IP        string `header:"X-Original-Forwarded-For" required:"false" doc:"Client IP from proxy"`
	UserAgent string `header:"User-Agent" required:"false" doc:"User agent string"`
	LoginType string `header:"Login-Type" required:"false" doc:"Login type"`
}

type UserLoginOutput struct {
	Body *types.LoginResponse
}

func UserLoginHandler(svcCtx *svc.ServiceContext) func(context.Context, *UserLoginInput) (*UserLoginOutput, error) {
	return func(ctx context.Context, input *UserLoginInput) (*UserLoginOutput, error) {
		input.Body.IP = input.IP
		input.Body.UserAgent = input.UserAgent
		input.Body.LoginType = input.LoginType
		if svcCtx.Config.Verify.LoginVerify && !svcCtx.Config.Debug {
			verifyTurns := turnstile.New(turnstile.Config{
				Secret:  svcCtx.Config.Verify.TurnstileSecret,
				Timeout: 3 * time.Second,
			})
			if verify, err := verifyTurns.Verify(ctx, input.Body.CfToken, input.Body.IP); err != nil || !verify {
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.TooManyRequests), "error: %v, verify: %v", err, verify)
			}
		}
		l := auth.NewUserLoginLogic(ctx, svcCtx)
		resp, err := l.UserLogin(&input.Body)
		if err != nil {
			return nil, err
		}
		return &UserLoginOutput{Body: resp}, nil
	}
}
