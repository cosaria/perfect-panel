// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/services/admin/log"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
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
