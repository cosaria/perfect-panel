// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CommissionWithdrawInput struct {
	Body types.CommissionWithdrawRequest
}

type CommissionWithdrawOutput struct {
	Body *types.WithdrawalLog
}

func CommissionWithdrawHandler(svcCtx *svc.ServiceContext) func(context.Context, *CommissionWithdrawInput) (*CommissionWithdrawOutput, error) {
	return func(ctx context.Context, input *CommissionWithdrawInput) (*CommissionWithdrawOutput, error) {
		l := NewCommissionWithdrawLogic(ctx, svcCtx)
		resp, err := l.CommissionWithdraw(&input.Body)
		if err != nil {
			return nil, err
		}
		return &CommissionWithdrawOutput{Body: resp}, nil
	}
}
