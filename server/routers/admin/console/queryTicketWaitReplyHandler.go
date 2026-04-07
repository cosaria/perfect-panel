// huma:migrated
package console

import (
	"context"
	"github.com/perfect-panel/server/services/admin/console"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type QueryTicketWaitReplyOutput struct {
	Body *types.TicketWaitRelpyResponse
}

func QueryTicketWaitReplyHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QueryTicketWaitReplyOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryTicketWaitReplyOutput, error) {
		l := console.NewQueryTicketWaitReplyLogic(ctx, svcCtx)
		resp, err := l.QueryTicketWaitReply()
		if err != nil {
			return nil, err
		}
		return &QueryTicketWaitReplyOutput{Body: resp}, nil
	}
}
