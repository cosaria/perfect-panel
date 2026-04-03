// huma:migrated
package common

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/common"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CheckVerificationCodeInput struct {
	Body types.CheckVerificationCodeRequest
}

type CheckVerificationCodeOutput struct {
	Body *types.CheckVerificationCodeRespone
}

func CheckVerificationCodeHandler(svcCtx *svc.ServiceContext) func(context.Context, *CheckVerificationCodeInput) (*CheckVerificationCodeOutput, error) {
	return func(ctx context.Context, input *CheckVerificationCodeInput) (*CheckVerificationCodeOutput, error) {
		l := common.NewCheckVerificationCodeLogic(ctx, svcCtx)
		resp, err := l.CheckVerificationCode(&input.Body)
		if err != nil {
			return nil, err
		}
		return &CheckVerificationCodeOutput{Body: resp}, nil
	}
}
