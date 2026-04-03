// huma:migrated
package tool

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/tool"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetSystemLogOutput struct {
	Body *types.LogResponse
}

func GetSystemLogHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetSystemLogOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetSystemLogOutput, error) {
		l := tool.NewGetSystemLogLogic(ctx, svcCtx)
		resp, err := l.GetSystemLog()
		if err != nil {
			return nil, err
		}
		return &GetSystemLogOutput{Body: resp}, nil
	}
}
