// huma:migrated
package common

import (
	"context"
	"github.com/perfect-panel/server/services/common"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type SendSmsCodeInput struct {
	Body types.SendSmsCodeRequest
}

type SendSmsCodeOutput struct {
	Body *types.SendCodeResponse
}

func SendSmsCodeHandler(svcCtx *svc.ServiceContext) func(context.Context, *SendSmsCodeInput) (*SendSmsCodeOutput, error) {
	return func(ctx context.Context, input *SendSmsCodeInput) (*SendSmsCodeOutput, error) {
		l := common.NewSendSmsCodeLogic(ctx, svcCtx)
		resp, err := l.SendSmsCode(&input.Body)
		if err != nil {
			return nil, err
		}
		return &SendSmsCodeOutput{Body: resp}, nil
	}
}
