package ticket

import (
	"context"
	"github.com/perfect-panel/server/models/ticket"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type CreateTicketFollowInput struct {
	Body types.CreateTicketFollowRequest
}

func CreateTicketFollowHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateTicketFollowInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateTicketFollowInput) (*struct{}, error) {
		l := NewCreateTicketFollowLogic(ctx, svcCtx)
		if err := l.CreateTicketFollow(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type CreateTicketFollowLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Create ticket follow
func NewCreateTicketFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTicketFollowLogic {
	return &CreateTicketFollowLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTicketFollowLogic) CreateTicketFollow(req *types.CreateTicketFollowRequest) (err error) {
	// find ticket
	_, err = l.svcCtx.TicketModel.FindOne(l.ctx, req.TicketId)
	if err != nil {
		l.Errorw("[CreateTicketFollow] FindOne error", logger.Field("error", err.Error()), logger.Field("ticketId", req.TicketId))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find ticket failed: %v", err.Error())
	}
	err = l.svcCtx.TicketModel.InsertTicketFollow(l.ctx, &ticket.Follow{
		TicketId: req.TicketId,
		From:     req.From,
		Type:     req.Type,
		Content:  req.Content,
	})
	if err != nil {
		l.Errorw("[CreateTicketFollow] Database insert error", logger.Field("error", err.Error()), logger.Field("request", req))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "create ticket follow failed: %v", err.Error())
	}
	err = l.svcCtx.TicketModel.UpdateTicketStatus(l.ctx, req.TicketId, 0, ticket.Waiting)
	if err != nil {
		l.Errorw("[CreateTicketFollow] Database update error", logger.Field("error", err.Error()), logger.Field("status", ticket.Waiting))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update ticket status failed: %v", err.Error())
	}
	return
}
