package ticket

import (
	"context"

	"github.com/perfect-panel/server/config"

	"github.com/perfect-panel/server/models/user"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type UpdateUserTicketStatusLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Update ticket status
func NewUpdateUserTicketStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserTicketStatusLogic {
	return &UpdateUserTicketStatusLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserTicketStatusLogic) UpdateUserTicketStatus(req *types.UpdateUserTicketStatusRequest) error {
	u, ok := l.ctx.Value(config.CtxKeyUser).(*user.User)
	if !ok {
		logger.Error("current user is not found in context")
		return errors.Wrapf(xerr.NewErrCode(xerr.InvalidAccess), "Invalid Access")
	}
	err := l.svcCtx.TicketModel.UpdateTicketStatus(l.ctx, req.Id, u.Id, *req.Status)
	if err != nil {
		l.Errorw("[UpdateUserTicketStatusLogic] Database Error", logger.Field("error", err.Error()))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update ticket error: %v", err.Error())
	}
	return nil
}
