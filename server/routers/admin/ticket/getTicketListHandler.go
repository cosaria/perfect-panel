// huma:migrated
package ticket

import (
	"context"
	"github.com/perfect-panel/server/services/admin/ticket"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetTicketListInput struct {
	Body types.GetTicketListRequest
}

type GetTicketListOutput struct {
	Body *types.GetTicketListResponse
}

func GetTicketListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetTicketListInput) (*GetTicketListOutput, error) {
	return func(ctx context.Context, input *GetTicketListInput) (*GetTicketListOutput, error) {
		l := ticket.NewGetTicketListLogic(ctx, svcCtx)
		resp, err := l.GetTicketList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetTicketListOutput{Body: resp}, nil
	}
}
