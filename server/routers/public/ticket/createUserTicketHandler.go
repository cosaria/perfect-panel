// huma:migrated
package ticket

import (
	"context"
	"github.com/perfect-panel/server/services/user/ticket"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type CreateUserTicketInput struct {
	Body types.CreateUserTicketRequest
}

func CreateUserTicketHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateUserTicketInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateUserTicketInput) (*struct{}, error) {
		l := ticket.NewCreateUserTicketLogic(ctx, svcCtx)
		if err := l.CreateUserTicket(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
