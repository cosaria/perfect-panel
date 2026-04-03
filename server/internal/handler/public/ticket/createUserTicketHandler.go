// huma:migrated
package ticket

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/ticket"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
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
