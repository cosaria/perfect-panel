package ticket

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type GetTicketInput struct {
	types.GetTicketRequest
}

type GetTicketOutput struct {
	Body *types.Ticket
}

func GetTicketHandler(deps Deps) func(context.Context, *GetTicketInput) (*GetTicketOutput, error) {
	return func(ctx context.Context, input *GetTicketInput) (*GetTicketOutput, error) {
		l := NewGetTicketLogic(ctx, deps)
		resp, err := l.GetTicket(&input.GetTicketRequest)
		if err != nil {
			return nil, err
		}
		return &GetTicketOutput{Body: resp}, nil
	}
}

type GetTicketLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get ticket detail
func NewGetTicketLogic(ctx context.Context, deps Deps) *GetTicketLogic {
	return &GetTicketLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetTicketLogic) GetTicket(req *types.GetTicketRequest) (resp *types.Ticket, err error) {
	data, err := l.deps.TicketModel.QueryTicketDetail(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[GetTicket] Query Database Error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get ticket detail failed: %v", err.Error())
	}
	resp = &types.Ticket{}
	tool.DeepCopy(resp, data)
	return
}
