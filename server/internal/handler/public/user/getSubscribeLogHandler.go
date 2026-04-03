// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetSubscribeLogInput struct {
	types.GetSubscribeLogRequest
}

type GetSubscribeLogOutput struct {
	Body *types.GetSubscribeLogResponse
}

func GetSubscribeLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetSubscribeLogInput) (*GetSubscribeLogOutput, error) {
	return func(ctx context.Context, input *GetSubscribeLogInput) (*GetSubscribeLogOutput, error) {
		l := user.NewGetSubscribeLogLogic(ctx, svcCtx)
		resp, err := l.GetSubscribeLog(&input.GetSubscribeLogRequest)
		if err != nil {
			return nil, err
		}
		return &GetSubscribeLogOutput{Body: resp}, nil
	}
}
