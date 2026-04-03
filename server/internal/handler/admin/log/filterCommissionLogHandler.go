// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/log"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type FilterCommissionLogInput struct {
	types.FilterCommissionLogRequest
}

type FilterCommissionLogOutput struct {
	Body *types.FilterCommissionLogResponse
}

func FilterCommissionLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterCommissionLogInput) (*FilterCommissionLogOutput, error) {
	return func(ctx context.Context, input *FilterCommissionLogInput) (*FilterCommissionLogOutput, error) {
		l := log.NewFilterCommissionLogLogic(ctx, svcCtx)
		resp, err := l.FilterCommissionLog(&input.FilterCommissionLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterCommissionLogOutput{Body: resp}, nil
	}
}
