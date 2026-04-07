// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/services/admin/log"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type FilterResetSubscribeLogInput struct {
	types.FilterResetSubscribeLogRequest
}

type FilterResetSubscribeLogOutput struct {
	Body *types.FilterResetSubscribeLogResponse
}

func FilterResetSubscribeLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterResetSubscribeLogInput) (*FilterResetSubscribeLogOutput, error) {
	return func(ctx context.Context, input *FilterResetSubscribeLogInput) (*FilterResetSubscribeLogOutput, error) {
		l := log.NewFilterResetSubscribeLogLogic(ctx, svcCtx)
		resp, err := l.FilterResetSubscribeLog(&input.FilterResetSubscribeLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterResetSubscribeLogOutput{Body: resp}, nil
	}
}
