// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateNodeConfigInput struct {
	Body types.NodeConfig
}

func UpdateNodeConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateNodeConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateNodeConfigInput) (*struct{}, error) {
		l := NewUpdateNodeConfigLogic(ctx, svcCtx)
		if err := l.UpdateNodeConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
