package user

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type PreUnsubscribeLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewPreUnsubscribeLogic Pre Unsubscribe
func NewPreUnsubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreUnsubscribeLogic {
	return &PreUnsubscribeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreUnsubscribeLogic) PreUnsubscribe(req *types.PreUnsubscribeRequest) (resp *types.PreUnsubscribeResponse, err error) {
	remainingAmount, err := CalculateRemainingAmount(l.ctx, l.svcCtx, req.Id)
	if err != nil {
		l.Errorw("[PreUnsubscribeLogic] Calculate Remaining Amount Error:", logger.Field("err", err.Error()))
		return nil, err
	}
	return &types.PreUnsubscribeResponse{
		DeductionAmount: remainingAmount,
	}, nil
}
