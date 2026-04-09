package ticket

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type GetTicketListInput struct {
	Body types.GetTicketListRequest
}

type GetTicketListOutput struct {
	Body *types.GetTicketListResponse
}

func GetTicketListHandler(deps Deps) func(context.Context, *GetTicketListInput) (*GetTicketListOutput, error) {
	return func(ctx context.Context, input *GetTicketListInput) (*GetTicketListOutput, error) {
		l := NewGetTicketListLogic(ctx, deps)
		resp, err := l.GetTicketList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetTicketListOutput{Body: resp}, nil
	}
}

type GetTicketListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get ticket list
func NewGetTicketListLogic(ctx context.Context, deps Deps) *GetTicketListLogic {
	return &GetTicketListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetTicketListLogic) GetTicketList(req *types.GetTicketListRequest) (resp *types.GetTicketListResponse, err error) {
	total, list, err := l.deps.TicketModel.QueryTicketList(l.ctx, int(req.Page), int(req.Size), req.UserId, req.Status, req.Search)
	if err != nil {
		l.Errorw("[GetTicketList] Query Database Error: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "QueryTicketList error: %v", err)
	}
	resp = &types.GetTicketListResponse{
		Total: total,
		List:  make([]types.Ticket, 0),
	}
	tool.DeepCopy(&resp.List, list)
	return
}
