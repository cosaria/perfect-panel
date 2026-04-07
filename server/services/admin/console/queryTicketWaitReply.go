package console

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryTicketWaitReplyOutput struct {
	Body *types.TicketWaitRelpyResponse
}

func QueryTicketWaitReplyHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QueryTicketWaitReplyOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryTicketWaitReplyOutput, error) {
		l := NewQueryTicketWaitReplyLogic(ctx, svcCtx)
		resp, err := l.QueryTicketWaitReply()
		if err != nil {
			return nil, err
		}
		return &QueryTicketWaitReplyOutput{Body: resp}, nil
	}
}

type QueryTicketWaitReplyLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewQueryTicketWaitReplyLogic Query ticket wait reply
func NewQueryTicketWaitReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryTicketWaitReplyLogic {
	return &QueryTicketWaitReplyLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryTicketWaitReplyLogic) QueryTicketWaitReply() (resp *types.TicketWaitRelpyResponse, err error) {
	count, err := l.svcCtx.TicketModel.QueryWaitReplyTotal(l.ctx)
	if err != nil {
		l.Errorw("[QueryTicketWaitReply] Query Database Error: ", logger.Field("error", err.Error()))
		return nil, err
	}
	return &types.TicketWaitRelpyResponse{
		Count: count,
	}, nil
}
