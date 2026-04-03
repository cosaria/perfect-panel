// huma:migrated
package tool

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/tool"
	"github.com/perfect-panel/server/internal/svc"
)

func RestartSystemHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*struct{}, error) {
	return func(ctx context.Context, _ *struct{}) (*struct{}, error) {
		l := tool.NewRestartSystemLogic(ctx, svcCtx)
		if err := l.RestartSystem(); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
