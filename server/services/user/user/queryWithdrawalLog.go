package user

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryWithdrawalLogInput struct {
	types.QueryWithdrawalLogListRequest
}

type QueryWithdrawalLogOutput struct {
	Body *types.QueryWithdrawalLogListResponse
}

func QueryWithdrawalLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *QueryWithdrawalLogInput) (*QueryWithdrawalLogOutput, error) {
	return func(ctx context.Context, input *QueryWithdrawalLogInput) (*QueryWithdrawalLogOutput, error) {
		l := NewQueryWithdrawalLogLogic(ctx, svcCtx)
		resp, err := l.QueryWithdrawalLog(&input.QueryWithdrawalLogListRequest)
		if err != nil {
			return nil, err
		}
		return &QueryWithdrawalLogOutput{Body: resp}, nil
	}
}

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
