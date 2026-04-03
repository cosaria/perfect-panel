// huma:migrated
package console

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/console"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type QueryServerTotalDataOutput struct {
	Body *types.ServerTotalDataResponse
}

func QueryServerTotalDataHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*QueryServerTotalDataOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*QueryServerTotalDataOutput, error) {
		l := console.NewQueryServerTotalDataLogic(ctx, svcCtx)
		resp, err := l.QueryServerTotalData()
		if err != nil {
			return nil, err
		}
		return &QueryServerTotalDataOutput{Body: resp}, nil
	}
}
