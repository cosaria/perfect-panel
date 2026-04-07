// huma:migrated
package ticket

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetUserTicketListInput struct {
	Body types.GetUserTicketListRequest
}

type GetUserTicketListOutput struct {
	Body *types.GetUserTicketListResponse
}

func GetUserTicketListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetUserTicketListInput) (*GetUserTicketListOutput, error) {
	return func(ctx context.Context, input *GetUserTicketListInput) (*GetUserTicketListOutput, error) {
		l := NewGetUserTicketListLogic(ctx, svcCtx)
		resp, err := l.GetUserTicketList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetUserTicketListOutput{Body: resp}, nil
	}
}
