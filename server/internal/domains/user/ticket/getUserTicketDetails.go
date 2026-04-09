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

type GetUserTicketDetailsInput struct {
	types.GetUserTicketDetailRequest
}

type GetUserTicketDetailsOutput struct {
	Body *types.Ticket
}

func GetUserTicketDetailsHandler(deps Deps) func(context.Context, *GetUserTicketDetailsInput) (*GetUserTicketDetailsOutput, error) {
	return func(ctx context.Context, input *GetUserTicketDetailsInput) (*GetUserTicketDetailsOutput, error) {
		l := NewGetUserTicketDetailsLogic(ctx, deps)
		resp, err := l.GetUserTicketDetails(&input.GetUserTicketDetailRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserTicketDetailsOutput{Body: resp}, nil
	}
}

type GetUserTicketDetailsLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get ticket detail
func NewGetUserTicketDetailsLogic(ctx context.Context, deps Deps) *GetUserTicketDetailsLogic {
	return &GetUserTicketDetailsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetUserTicketDetailsLogic) GetUserTicketDetails(req *types.GetUserTicketDetailRequest) (resp *types.Ticket, err error) {

	data, err := l.deps.TicketModel.QueryTicketDetail(l.ctx, req.Id)
	if err != nil {
		l.Errorw("[GetUserTicketDetailsLogic] Database Error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get ticket detail failed: %v", err.Error())
	}
	// check access
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	if data.UserId != u.Id {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "invalid access")
	}
	resp = &types.Ticket{}
	tool.DeepCopy(resp, data)
	return
}
