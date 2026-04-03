// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/log"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type FilterEmailLogInput struct {
	types.FilterLogParams
}

type FilterEmailLogOutput struct {
	Body *types.FilterEmailLogResponse
}

func FilterEmailLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterEmailLogInput) (*FilterEmailLogOutput, error) {
	return func(ctx context.Context, input *FilterEmailLogInput) (*FilterEmailLogOutput, error) {
		l := log.NewFilterEmailLogLogic(ctx, svcCtx)
		resp, err := l.FilterEmailLog(&input.FilterLogParams)
		if err != nil {
			return nil, err
		}
		return &FilterEmailLogOutput{Body: resp}, nil
	}
}
