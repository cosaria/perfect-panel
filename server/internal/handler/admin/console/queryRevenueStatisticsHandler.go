// huma:migrated
package console

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/console"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type QueryRevenueStatisticsOutput struct {
	Body *types.RevenueStatisticsResponse
}

func QueryRevenueStatisticsHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QueryRevenueStatisticsOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryRevenueStatisticsOutput, error) {
		l := console.NewQueryRevenueStatisticsLogic(ctx, svcCtx)
		resp, err := l.QueryRevenueStatistics()
		if err != nil {
			return nil, err
		}
		return &QueryRevenueStatisticsOutput{Body: resp}, nil
	}
}
