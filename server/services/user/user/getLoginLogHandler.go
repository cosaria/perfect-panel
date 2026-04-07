// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetLoginLogInput struct {
	types.GetLoginLogRequest
}

type GetLoginLogOutput struct {
	Body *types.GetLoginLogResponse
}

func GetLoginLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetLoginLogInput) (*GetLoginLogOutput, error) {
	return func(ctx context.Context, input *GetLoginLogInput) (*GetLoginLogOutput, error) {
		l := NewGetLoginLogLogic(ctx, svcCtx)
		resp, err := l.GetLoginLog(&input.GetLoginLogRequest)
		if err != nil {
			return nil, err
		}
		return &GetLoginLogOutput{Body: resp}, nil
	}
}
