// huma:migrated
package ticket

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetUserTicketDetailsInput struct {
	types.GetUserTicketDetailRequest
}

type GetUserTicketDetailsOutput struct {
	Body *types.Ticket
}

func GetUserTicketDetailsHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserTicketDetailsInput) (*GetUserTicketDetailsOutput, error) {
	return func(ctx context.Context, input *GetUserTicketDetailsInput) (*GetUserTicketDetailsOutput, error) {
		l := NewGetUserTicketDetailsLogic(ctx, svcCtx)
		resp, err := l.GetUserTicketDetails(&input.GetUserTicketDetailRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserTicketDetailsOutput{Body: resp}, nil
	}
}
