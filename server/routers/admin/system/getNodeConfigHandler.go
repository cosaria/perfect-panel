// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/services/admin/system"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetNodeConfigOutput struct {
	Body *types.NodeConfig
}

func GetNodeConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetNodeConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetNodeConfigOutput, error) {
		l := system.NewGetNodeConfigLogic(ctx, svcCtx)
		resp, err := l.GetNodeConfig()
		if err != nil {
			return nil, err
		}
		return &GetNodeConfigOutput{Body: resp}, nil
	}
}
