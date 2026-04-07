package ticket

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type UpdateTicketStatusInput struct {
	Body types.UpdateTicketStatusRequest
}

func UpdateTicketStatusHandler(deps Deps) func(context.Context, *UpdateTicketStatusInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateTicketStatusInput) (*struct{}, error) {
		l := NewUpdateTicketStatusLogic(ctx, deps)
		if err := l.UpdateTicketStatus(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type UpdateTicketStatusLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Update ticket status
func NewUpdateTicketStatusLogic(ctx context.Context, deps Deps) *UpdateTicketStatusLogic {
	return &UpdateTicketStatusLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *UpdateTicketStatusLogic) UpdateTicketStatus(req *types.UpdateTicketStatusRequest) error {

	err := l.deps.TicketModel.UpdateTicketStatus(l.ctx, req.Id, 0, *req.Status)
	if err != nil {
		l.Errorw("[UpdateTicketStatus] Update Database Error: ", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update ticket error: %v", err.Error())
	}
	return nil
}
