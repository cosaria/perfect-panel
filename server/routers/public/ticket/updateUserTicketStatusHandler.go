// huma:migrated
package ticket

import (
	"context"
	"github.com/perfect-panel/server/services/user/ticket"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
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
