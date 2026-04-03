// huma:migrated
package ticket

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/ticket"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CreateUserTicketFollowInput struct {
	Body types.CreateUserTicketFollowRequest
}

func CreateUserTicketFollowHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateUserTicketFollowInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateUserTicketFollowInput) (*struct{}, error) {
		l := ticket.NewCreateUserTicketFollowLogic(ctx, svcCtx)
		if err := l.CreateUserTicketFollow(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
