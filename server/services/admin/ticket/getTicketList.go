package ticket

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetTicketListInput struct {
	Body types.GetTicketListRequest
}

type GetTicketListOutput struct {
	Body *types.GetTicketListResponse
}

func GetTicketListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetTicketListInput) (*GetTicketListOutput, error) {
	return func(ctx context.Context, input *GetTicketListInput) (*GetTicketListOutput, error) {
		l := NewGetTicketListLogic(ctx, svcCtx)
		resp, err := l.GetTicketList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetTicketListOutput{Body: resp}, nil
	}
}

type GetTicketListLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get ticket list
func NewGetTicketListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTicketListLogic {
	return &GetTicketListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTicketListLogic) GetTicketList(req *types.GetTicketListRequest) (resp *types.GetTicketListResponse, err error) {
	total, list, err := l.svcCtx.TicketModel.QueryTicketList(l.ctx, int(req.Page), int(req.Size), req.UserId, req.Status, req.Search)
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
