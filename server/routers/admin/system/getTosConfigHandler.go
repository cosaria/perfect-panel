// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/services/admin/system"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetTosConfigOutput struct {
	Body *types.TosConfig
}

func GetTosConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetTosConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetTosConfigOutput, error) {
		l := system.NewGetTosConfigLogic(ctx, svcCtx)
		resp, err := l.GetTosConfig()
		if err != nil {
			return nil, err
		}
		return &GetTosConfigOutput{Body: resp}, nil
	}
}
