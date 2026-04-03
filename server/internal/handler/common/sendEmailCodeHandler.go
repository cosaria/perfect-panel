// huma:migrated
package common

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/common"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type SendEmailCodeInput struct {
	Body types.SendCodeRequest
}

type SendEmailCodeOutput struct {
	Body *types.SendCodeResponse
}

func SendEmailCodeHandler(svcCtx *svc.ServiceContext) func(context.Context, *SendEmailCodeInput) (*SendEmailCodeOutput, error) {
	return func(ctx context.Context, input *SendEmailCodeInput) (*SendEmailCodeOutput, error) {
		l := common.NewSendEmailCodeLogic(ctx, svcCtx)
		resp, err := l.SendEmailCode(&input.Body)
		if err != nil {
			return nil, err
		}
		return &SendEmailCodeOutput{Body: resp}, nil
	}
}
