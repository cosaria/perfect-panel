package ticket

import (
	"context"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetUserTicketListInput struct {
	Body types.GetUserTicketListRequest
}

type GetUserTicketListOutput struct {
	Body *types.GetUserTicketListResponse
}

func GetUserTicketListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserTicketListInput) (*GetUserTicketListOutput, error) {
	return func(ctx context.Context, input *GetUserTicketListInput) (*GetUserTicketListOutput, error) {
		l := NewGetUserTicketListLogic(ctx, svcCtx)
		resp, err := l.GetUserTicketList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetUserTicketListOutput{Body: resp}, nil
	}
}

type GetUserTicketListLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get ticket list
func NewGetUserTicketListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserTicketListLogic {
	return &GetUserTicketListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserTicketListLogic) GetUserTicketList(req *types.GetUserTicketListRequest) (resp *types.GetUserTicketListResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	l.Logger.Debugf("Current user: %v", u.Id)
	total, list, err := l.svcCtx.TicketModel.QueryTicketList(l.ctx, req.Page, req.Size, u.Id, req.Status, req.Search)
	if err != nil {
		l.Errorw("[GetUserTicketListLogic] Database Error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "QueryTicketList error: %v", err)
	}
	resp = &types.GetUserTicketListResponse{
		Total: total,
		List:  make([]types.Ticket, 0),
	}
	tool.DeepCopy(&resp.List, list)
	return
}
