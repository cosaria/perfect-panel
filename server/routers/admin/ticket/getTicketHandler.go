// huma:migrated
package ticket

import (
	"context"
	"github.com/perfect-panel/server/services/admin/ticket"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetTicketInput struct {
	types.GetTicketRequest
}

type GetTicketOutput struct {
	Body *types.Ticket
}

func GetTicketHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetTicketInput) (*GetTicketOutput, error) {
	return func(ctx context.Context, input *GetTicketInput) (*GetTicketOutput, error) {
		l := ticket.NewGetTicketLogic(ctx, svcCtx)
		resp, err := l.GetTicket(&input.GetTicketRequest)
		if err != nil {
			return nil, err
		}
		return &GetTicketOutput{Body: resp}, nil
	}
}
