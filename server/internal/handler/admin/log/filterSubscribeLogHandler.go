// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/log"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type FilterSubscribeLogInput struct {
	types.FilterSubscribeLogRequest
}

type FilterSubscribeLogOutput struct {
	Body *types.FilterSubscribeLogResponse
}

func FilterSubscribeLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterSubscribeLogInput) (*FilterSubscribeLogOutput, error) {
	return func(ctx context.Context, input *FilterSubscribeLogInput) (*FilterSubscribeLogOutput, error) {
		l := log.NewFilterSubscribeLogLogic(ctx, svcCtx)
		resp, err := l.FilterSubscribeLog(&input.FilterSubscribeLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterSubscribeLogOutput{Body: resp}, nil
	}
}
