// huma:migrated
package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/verify/turnstile"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type TelephoneLoginInput struct {
	Body      types.TelephoneLoginRequest
	IP        string `header:"X-Original-Forwarded-For" required:"false" doc:"Client IP from proxy"`
	UserAgent string `header:"User-Agent" required:"false" doc:"User agent string"`
	LoginType string `header:"Login-Type" required:"false" doc:"Login type"`
}

type TelephoneLoginOutput struct {
	Body *types.LoginResponse
}

func TelephoneLoginHandler(svcCtx *svc.ServiceContext) func(context.Context, *TelephoneLoginInput) (*TelephoneLoginOutput, error) {
	return func(ctx context.Context, input *TelephoneLoginInput) (*TelephoneLoginOutput, error) {
		input.Body.IP = input.IP
		input.Body.UserAgent = input.UserAgent
		input.Body.LoginType = input.LoginType
		if svcCtx.Config.Verify.LoginVerify {
			verifyTurns := turnstile.New(turnstile.Config{
				Secret:  svcCtx.Config.Verify.TurnstileSecret,
				Timeout: 3 * time.Second,
			})
			if verify, err := verifyTurns.Verify(ctx, input.Body.CfToken, input.Body.IP); err != nil || !verify {
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.TooManyRequests), "error: %v, verify: %v", err, verify)
			}
		}
		// Construct a minimal *http.Request with User-Agent header for the logic layer
		r := &http.Request{Header: http.Header{}}
		r.Header.Set("User-Agent", input.UserAgent)
		l := NewTelephoneLoginLogic(ctx, svcCtx)
		resp, err := l.TelephoneLogin(&input.Body, r, input.IP)
		if err != nil {
			return nil, err
		}
		return &TelephoneLoginOutput{Body: resp}, nil
	}
}
