// huma:migrated
package user

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/public/user"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetLoginLogInput struct {
	types.GetLoginLogRequest
}

type GetLoginLogOutput struct {
	Body *types.GetLoginLogResponse
}

func GetLoginLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetLoginLogInput) (*GetLoginLogOutput, error) {
	return func(ctx context.Context, input *GetLoginLogInput) (*GetLoginLogOutput, error) {
		l := user.NewGetLoginLogLogic(ctx, svcCtx)
		resp, err := l.GetLoginLog(&input.GetLoginLogRequest)
		if err != nil {
			return nil, err
		}
		return &GetLoginLogOutput{Body: resp}, nil
	}
}
