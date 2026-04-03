// huma:migrated
package console

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/console"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type QueryUserStatisticsOutput struct {
	Body *types.UserStatisticsResponse
}

func QueryUserStatisticsHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QueryUserStatisticsOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryUserStatisticsOutput, error) {
		l := console.NewQueryUserStatisticsLogic(ctx, svcCtx)
		resp, err := l.QueryUserStatistics()
		if err != nil {
			return nil, err
		}
		return &QueryUserStatisticsOutput{Body: resp}, nil
	}
}
