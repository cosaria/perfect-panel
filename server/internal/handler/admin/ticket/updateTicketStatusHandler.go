// huma:migrated
package ticket

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/ticket"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateTicketStatusInput struct {
	Body types.UpdateTicketStatusRequest
}

func UpdateTicketStatusHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateTicketStatusInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateTicketStatusInput) (*struct{}, error) {
		l := ticket.NewUpdateTicketStatusLogic(ctx, svcCtx)
		if err := l.UpdateTicketStatus(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
