// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/system"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetModuleConfigOutput struct {
	Body *types.ModuleConfig
}

func GetModuleConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetModuleConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetModuleConfigOutput, error) {
		l := system.NewGetModuleConfigLogic(ctx, svcCtx)
		resp, err := l.GetModuleConfig()
		if err != nil {
			return nil, err
		}
		return &GetModuleConfigOutput{Body: resp}, nil
	}
}
