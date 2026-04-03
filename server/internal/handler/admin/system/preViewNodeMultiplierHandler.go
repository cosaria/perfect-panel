// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/system"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type PreViewNodeMultiplierOutput struct {
	Body *types.PreViewNodeMultiplierResponse
}

func PreViewNodeMultiplierHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*PreViewNodeMultiplierOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*PreViewNodeMultiplierOutput, error) {
		l := system.NewPreViewNodeMultiplierLogic(ctx, svcCtx)
		resp, err := l.PreViewNodeMultiplier()
		if err != nil {
			return nil, err
		}
		return &PreViewNodeMultiplierOutput{Body: resp}, nil
	}
}
