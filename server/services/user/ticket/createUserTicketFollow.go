package ticket

import (
	"context"
	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/ticket"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type CreateUserTicketFollowInput struct {
	Body types.CreateUserTicketFollowRequest
}

func CreateUserTicketFollowHandler(deps Deps) func(context.Context, *CreateUserTicketFollowInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateUserTicketFollowInput) (*struct{}, error) {
		l := NewCreateUserTicketFollowLogic(ctx, deps)
		if err := l.CreateUserTicketFollow(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type CreateUserTicketFollowLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Create ticket follow
func NewCreateUserTicketFollowLogic(ctx context.Context, deps Deps) *CreateUserTicketFollowLogic {
	return &CreateUserTicketFollowLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *CreateUserTicketFollowLogic) CreateUserTicketFollow(req *types.CreateUserTicketFollowRequest) error {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	// query ticket
	t, err := l.deps.TicketModel.FindOne(l.ctx, req.TicketId)
	if err != nil {
		l.Errorw("[CreateUserTicketFollow] Database query error", logger.Field("error", err.Error()), logger.Field("request", req))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "query ticket failed: %v", err.Error())
	}
	// check access
	if u.Id != t.UserId {
		l.Errorw("[CreateUserTicketFollow] Invalid access", logger.Field("user_id", u.Id), logger.Field("ticket_user_id", t.UserId))
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "invalid access")
	}
	// insert follow
	err = l.deps.TicketModel.InsertTicketFollow(l.ctx, &ticket.Follow{
		TicketId: req.TicketId,
		From:     req.From,
		Type:     req.Type,
		Content:  req.Content,
	})
	if err != nil {
		l.Errorw("[CreateUserTicketFollow] Database insert error", logger.Field("error", err.Error()), logger.Field("request", req))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "create ticket follow failed: %v", err.Error())
	}
	err = l.deps.TicketModel.UpdateTicketStatus(l.ctx, req.TicketId, u.Id, ticket.Pending)
	if err != nil {
		l.Errorw("[CreateUserTicketFollow] Database update error", logger.Field("error", err.Error()), logger.Field("status", ticket.Pending))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update ticket status failed: %v", err.Error())
	}
	return nil
}
