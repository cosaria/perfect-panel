// huma:migrated
package auth

import (
	"context"
	"time"

	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/verify/turnstile"
	"github.com/perfect-panel/server/services/auth"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type TelephoneResetPasswordInput struct {
	Body      types.TelephoneResetPasswordRequest
	IP        string `header:"X-Original-Forwarded-For" required:"false" doc:"Client IP from proxy"`
	UserAgent string `header:"User-Agent" required:"false" doc:"User agent string"`
	LoginType string `header:"Login-Type" required:"false" doc:"Login type"`
}

type TelephoneResetPasswordOutput struct {
	Body *types.LoginResponse
}

func TelephoneResetPasswordHandler(svcCtx *svc.ServiceContext) func(context.Context, *TelephoneResetPasswordInput) (*TelephoneResetPasswordOutput, error) {
	return func(ctx context.Context, input *TelephoneResetPasswordInput) (*TelephoneResetPasswordOutput, error) {
		input.Body.IP = input.IP
		input.Body.UserAgent = input.UserAgent
		input.Body.LoginType = input.LoginType
		if svcCtx.Config.Verify.ResetPasswordVerify {
			verifyTurns := turnstile.New(turnstile.Config{
				Secret:  svcCtx.Config.Verify.TurnstileSecret,
				Timeout: 3 * time.Second,
			})
			if verify, err := verifyTurns.Verify(ctx, input.Body.CfToken, input.Body.IP); err != nil || !verify {
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.TooManyRequests), "error: %v, verify: %v", err, verify)
			}
		}
		l := auth.NewTelephoneResetPasswordLogic(ctx, svcCtx)
		resp, err := l.TelephoneResetPassword(&input.Body)
		if err != nil {
			return nil, err
		}
		return &TelephoneResetPasswordOutput{Body: resp}, nil
	}
}
