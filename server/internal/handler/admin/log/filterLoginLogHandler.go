// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/log"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type FilterLoginLogInput struct {
	types.FilterLoginLogRequest
}

type FilterLoginLogOutput struct {
	Body *types.FilterLoginLogResponse
}

func FilterLoginLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterLoginLogInput) (*FilterLoginLogOutput, error) {
	return func(ctx context.Context, input *FilterLoginLogInput) (*FilterLoginLogOutput, error) {
		l := log.NewFilterLoginLogLogic(ctx, svcCtx)
		resp, err := l.FilterLoginLog(&input.FilterLoginLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterLoginLogOutput{Body: resp}, nil
	}
}
