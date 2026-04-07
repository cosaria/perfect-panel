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

type CreateUserTicketInput struct {
	Body types.CreateUserTicketRequest
}

func CreateUserTicketHandler(deps Deps) func(context.Context, *CreateUserTicketInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateUserTicketInput) (*struct{}, error) {
		l := NewCreateUserTicketLogic(ctx, deps)
		if err := l.CreateUserTicket(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type CreateUserTicketLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Create ticket
func NewCreateUserTicketLogic(ctx context.Context, deps Deps) *CreateUserTicketLogic {
	return &CreateUserTicketLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *CreateUserTicketLogic) CreateUserTicket(req *types.CreateUserTicketRequest) error {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	err := l.deps.TicketModel.Insert(l.ctx, &ticket.Ticket{
		Title:       req.Title,
		Description: req.Description,
		UserId:      u.Id,
		Status:      ticket.Pending,
	})
	if err != nil {
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "insert ticket error: %v", err.Error())
	}
	return nil
}
