// huma:migrated
package ticket

import (
	"context"
	"github.com/perfect-panel/server/services/admin/ticket"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CreateTicketFollowInput struct {
	Body types.CreateTicketFollowRequest
}

func CreateTicketFollowHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateTicketFollowInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateTicketFollowInput) (*struct{}, error) {
		l := ticket.NewCreateTicketFollowLogic(ctx, svcCtx)
		if err := l.CreateTicketFollow(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
