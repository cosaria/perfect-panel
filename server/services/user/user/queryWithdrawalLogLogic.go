package user

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryWithdrawalLogLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewQueryWithdrawalLogLogic Query Withdrawal Log
func NewQueryWithdrawalLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryWithdrawalLogLogic {
	return &QueryWithdrawalLogLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryWithdrawalLogLogic) QueryWithdrawalLog(req *types.QueryWithdrawalLogListRequest) (resp *types.QueryWithdrawalLogListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
