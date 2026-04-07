// huma:migrated
package tool

import (
	"context"
	"github.com/perfect-panel/server/services/admin/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
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
