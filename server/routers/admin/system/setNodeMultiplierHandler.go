// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/services/admin/system"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type SetNodeMultiplierInput struct {
	Body types.SetNodeMultiplierRequest
}

func SetNodeMultiplierHandler(svcCtx *svc.ServiceContext) func(context.Context, *SetNodeMultiplierInput) (*struct{}, error) {
	return func(ctx context.Context, input *SetNodeMultiplierInput) (*struct{}, error) {
		l := system.NewSetNodeMultiplierLogic(ctx, svcCtx)
		if err := l.SetNodeMultiplier(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
