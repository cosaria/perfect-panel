// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/services/admin/log"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type FilterMobileLogInput struct {
	types.FilterLogParams
}

type FilterMobileLogOutput struct {
	Body *types.FilterMobileLogResponse
}

func FilterMobileLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterMobileLogInput) (*FilterMobileLogOutput, error) {
	return func(ctx context.Context, input *FilterMobileLogInput) (*FilterMobileLogOutput, error) {
		l := log.NewFilterMobileLogLogic(ctx, svcCtx)
		resp, err := l.FilterMobileLog(&input.FilterLogParams)
		if err != nil {
			return nil, err
		}
		return &FilterMobileLogOutput{Body: resp}, nil
	}
}
