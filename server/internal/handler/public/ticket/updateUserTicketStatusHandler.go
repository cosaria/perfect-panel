// huma:migrated
package ticket

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/ticket"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateUserTicketStatusInput struct {
	Body types.UpdateUserTicketStatusRequest
}

func UpdateUserTicketStatusHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateUserTicketStatusInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateUserTicketStatusInput) (*struct{}, error) {
		l := ticket.NewUpdateUserTicketStatusLogic(ctx, svcCtx)
		if err := l.UpdateUserTicketStatus(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
