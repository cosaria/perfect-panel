// huma:migrated
package ticket

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/ticket"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetUserTicketListInput struct {
	Body types.GetUserTicketListRequest
}

type GetUserTicketListOutput struct {
	Body *types.GetUserTicketListResponse
}

func GetUserTicketListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserTicketListInput) (*GetUserTicketListOutput, error) {
	return func(ctx context.Context, input *GetUserTicketListInput) (*GetUserTicketListOutput, error) {
		l := ticket.NewGetUserTicketListLogic(ctx, svcCtx)
		resp, err := l.GetUserTicketList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetUserTicketListOutput{Body: resp}, nil
	}
}
