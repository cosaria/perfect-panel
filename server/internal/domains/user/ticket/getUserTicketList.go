package ticket

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/user"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type GetUserTicketListInput struct {
	Body types.GetUserTicketListRequest
}

type GetUserTicketListOutput struct {
	Body *types.GetUserTicketListResponse
}

func GetUserTicketListHandler(deps Deps) func(context.Context, *GetUserTicketListInput) (*GetUserTicketListOutput, error) {
	return func(ctx context.Context, input *GetUserTicketListInput) (*GetUserTicketListOutput, error) {
		l := NewGetUserTicketListLogic(ctx, deps)
		resp, err := l.GetUserTicketList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetUserTicketListOutput{Body: resp}, nil
	}
}

type GetUserTicketListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get ticket list
func NewGetUserTicketListLogic(ctx context.Context, deps Deps) *GetUserTicketListLogic {
	return &GetUserTicketListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetUserTicketListLogic) GetUserTicketList(req *types.GetUserTicketListRequest) (resp *types.GetUserTicketListResponse, err error) {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	l.Debugf("Current user: %v", u.Id)
	total, list, err := l.deps.TicketModel.QueryTicketList(l.ctx, req.Page, req.Size, u.Id, req.Status, req.Search)
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
