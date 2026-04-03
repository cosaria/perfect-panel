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

type UserRegisterInput struct {
	Body      types.UserRegisterRequest
	IP        string `header:"X-Original-Forwarded-For" required:"false" doc:"Client IP from proxy"`
	UserAgent string `header:"User-Agent" required:"false" doc:"User agent string"`
}

type UserRegisterOutput struct {
	Body *types.LoginResponse
}

func UserRegisterHandler(svcCtx *svc.ServiceContext) func(context.Context, *UserRegisterInput) (*UserRegisterOutput, error) {
	return func(ctx context.Context, input *UserRegisterInput) (*UserRegisterOutput, error) {
		input.Body.IP = input.IP
		input.Body.UserAgent = input.UserAgent
		if svcCtx.Config.Verify.RegisterVerify {
			verifyTurns := turnstile.New(turnstile.Config{
				Secret:  svcCtx.Config.Verify.TurnstileSecret,
				Timeout: 3 * time.Second,
			})
			if verify, err := verifyTurns.Verify(ctx, input.Body.CfToken, input.Body.IP); err != nil || !verify {
				return nil, errors.Wrapf(xerr.NewErrCode(xerr.TooManyRequests), "verify error: %v", err)
			}
		}
		l := auth.NewUserRegisterLogic(ctx, svcCtx)
		resp, err := l.UserRegister(&input.Body)
		if err != nil {
			return nil, err
		}
		return &UserRegisterOutput{Body: resp}, nil
	}
}
