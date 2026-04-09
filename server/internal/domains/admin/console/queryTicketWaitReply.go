package console

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
)

type QueryTicketWaitReplyOutput struct {
	Body *types.TicketWaitRelpyResponse
}

func QueryTicketWaitReplyHandler(deps Deps) func(context.Context, *struct{}) (*QueryTicketWaitReplyOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryTicketWaitReplyOutput, error) {
		l := NewQueryTicketWaitReplyLogic(ctx, deps)
		resp, err := l.QueryTicketWaitReply()
		if err != nil {
			return nil, err
		}
		return &QueryTicketWaitReplyOutput{Body: resp}, nil
	}
}

type QueryTicketWaitReplyLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewQueryTicketWaitReplyLogic Query ticket wait reply
func NewQueryTicketWaitReplyLogic(ctx context.Context, deps Deps) *QueryTicketWaitReplyLogic {
	return &QueryTicketWaitReplyLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryTicketWaitReplyLogic) QueryTicketWaitReply() (resp *types.TicketWaitRelpyResponse, err error) {
	count, err := l.deps.TicketModel.QueryWaitReplyTotal(l.ctx)
	if err != nil {
		l.Errorw("[QueryTicketWaitReply] Query Database Error: ", logger.Field("error", err.Error()))
		return nil, err
	}
	return &types.TicketWaitRelpyResponse{
		Count: count,
	}, nil
}
