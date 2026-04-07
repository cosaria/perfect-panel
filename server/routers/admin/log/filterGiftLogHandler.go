// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/services/admin/log"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type FilterGiftLogInput struct {
	types.FilterGiftLogRequest
}

type FilterGiftLogOutput struct {
	Body *types.FilterGiftLogResponse
}

func FilterGiftLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterGiftLogInput) (*FilterGiftLogOutput, error) {
	return func(ctx context.Context, input *FilterGiftLogInput) (*FilterGiftLogOutput, error) {
		l := log.NewFilterGiftLogLogic(ctx, svcCtx)
		resp, err := l.FilterGiftLog(&input.FilterGiftLogRequest)
		if err != nil {
			return nil, err
		}
		return &FilterGiftLogOutput{Body: resp}, nil
	}
}
