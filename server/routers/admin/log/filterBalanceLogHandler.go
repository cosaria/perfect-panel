// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/services/admin/log"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type FilterBalanceLogInput struct {
	types.FilterBalanceLogRequest
}

type FilterBalanceLogOutput struct {
	Body *types.FilterBalanceLogResponse
}

func FilterBalanceLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterBalanceLogInput) (*FilterBalanceLogOutput, error) {
	return func(ctx context.Context, input *FilterBalanceLogInput) (*FilterBalanceLogOutput, error) {
		l := log.NewFilterBalanceLogLogic(ctx, svcCtx)
		resp, err := l.FilterBalanceLog(&input.FilterBalanceLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterBalanceLogOutput{Body: resp}, nil
	}
}
