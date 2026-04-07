// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/services/admin/log"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type FilterTrafficLogDetailsInput struct {
	types.FilterTrafficLogDetailsRequest
}

type FilterTrafficLogDetailsOutput struct {
	Body *types.FilterTrafficLogDetailsResponse
}

func FilterTrafficLogDetailsHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterTrafficLogDetailsInput) (*FilterTrafficLogDetailsOutput, error) {
	return func(ctx context.Context, input *FilterTrafficLogDetailsInput) (*FilterTrafficLogDetailsOutput, error) {
		l := log.NewFilterTrafficLogDetailsLogic(ctx, svcCtx)
		resp, err := l.FilterTrafficLogDetails(&input.FilterTrafficLogDetailsRequest)
		if err != nil {
			return nil, err
		}
		return &FilterTrafficLogDetailsOutput{Body: resp}, nil
	}
}
