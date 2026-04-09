package ticket

import (
	"context"

	"github.com/perfect-panel/server/config"
	"github.com/perfect-panel/server/models/ticket"
	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type UpdateUserTicketStatusInput struct {
	Body types.UpdateUserTicketStatusRequest
}

func UpdateUserTicketStatusHandler(deps Deps) func(context.Context, *UpdateUserTicketStatusInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserTicketStatusInput) (*struct{}, error) {
		l := NewUpdateUserTicketStatusLogic(ctx, deps)
		if err := l.UpdateUserTicketStatus(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateUserTicketStatusLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Update ticket status
func NewUpdateUserTicketStatusLogic(ctx context.Context, deps Deps) *UpdateUserTicketStatusLogic {
	return &UpdateUserTicketStatusLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateUserTicketStatusLogic) UpdateUserTicketStatus(req *types.UpdateUserTicketStatusRequest) error {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	if req.Id <= 0 || req.Status == nil || *req.Status != ticket.Closed {
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidParams), "invalid ticket status update request")
	}

	err := l.deps.TicketModel.UpdateTicketStatus(l.ctx, req.Id, u.Id, *req.Status)
	if err != nil {
		l.Errorw("[UpdateUserTicketStatusLogic] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update ticket error: %v", err.Error())
	}
	return nil
}
