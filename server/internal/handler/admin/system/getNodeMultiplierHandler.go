// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/system"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetNodeMultiplierOutput struct {
	Body *types.GetNodeMultiplierResponse
}

func GetNodeMultiplierHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetNodeMultiplierOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetNodeMultiplierOutput, error) {
		l := system.NewGetNodeMultiplierLogic(ctx, svcCtx)
		resp, err := l.GetNodeMultiplier()
		if err != nil {
			return nil, err
		}
		return &GetNodeMultiplierOutput{Body: resp}, nil
	}
}
