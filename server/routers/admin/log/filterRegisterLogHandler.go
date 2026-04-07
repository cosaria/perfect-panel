// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/services/admin/log"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type FilterRegisterLogInput struct {
	types.FilterRegisterLogRequest
}

type FilterRegisterLogOutput struct {
	Body *types.FilterRegisterLogResponse
}

func FilterRegisterLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterRegisterLogInput) (*FilterRegisterLogOutput, error) {
	return func(ctx context.Context, input *FilterRegisterLogInput) (*FilterRegisterLogOutput, error) {
		l := log.NewFilterRegisterLogLogic(ctx, svcCtx)
		resp, err := l.FilterRegisterLog(&input.FilterRegisterLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterRegisterLogOutput{Body: resp}, nil
	}
}
