// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetNodeMultiplierOutput struct {
	Body *types.GetNodeMultiplierResponse
}

func GetNodeMultiplierHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetNodeMultiplierOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetNodeMultiplierOutput, error) {
		l := NewGetNodeMultiplierLogic(ctx, svcCtx)
		resp, err := l.GetNodeMultiplier()
		if err != nil {
			return nil, err
		}
		return &GetNodeMultiplierOutput{Body: resp}, nil
	}
}
