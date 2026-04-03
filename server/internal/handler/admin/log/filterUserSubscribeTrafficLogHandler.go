// huma:migrated
package log

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/log"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type FilterUserSubscribeTrafficLogInput struct {
	types.FilterSubscribeTrafficRequest
}

type FilterUserSubscribeTrafficLogOutput struct {
	Body *types.FilterSubscribeTrafficResponse
}

func FilterUserSubscribeTrafficLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *FilterUserSubscribeTrafficLogInput) (*FilterUserSubscribeTrafficLogOutput, error) {
	return func(ctx context.Context, input *FilterUserSubscribeTrafficLogInput) (*FilterUserSubscribeTrafficLogOutput, error) {
		l := log.NewFilterUserSubscribeTrafficLogLogic(ctx, svcCtx)
		resp, err := l.FilterUserSubscribeTrafficLog(&input.FilterSubscribeTrafficRequest)
		if err != nil {
			return nil, err
		}
		return &FilterUserSubscribeTrafficLogOutput{Body: resp}, nil
	}
}
