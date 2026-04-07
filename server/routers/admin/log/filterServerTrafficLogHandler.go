// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/services/admin/log"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type FilterServerTrafficLogInput struct {
	types.FilterServerTrafficLogRequest
}

type FilterServerTrafficLogOutput struct {
	Body *types.FilterServerTrafficLogResponse
}

func FilterServerTrafficLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterServerTrafficLogInput) (*FilterServerTrafficLogOutput, error) {
	return func(ctx context.Context, input *FilterServerTrafficLogInput) (*FilterServerTrafficLogOutput, error) {
		l := log.NewFilterServerTrafficLogLogic(ctx, svcCtx)
		resp, err := l.FilterServerTrafficLog(&input.FilterServerTrafficLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterServerTrafficLogOutput{Body: resp}, nil
	}
}
