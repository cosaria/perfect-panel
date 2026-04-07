// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateTosConfigInput struct {
	Body types.TosConfig
}

func UpdateTosConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateTosConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateTosConfigInput) (*struct{}, error) {
		l := NewUpdateTosConfigLogic(ctx, svcCtx)
		if err := l.UpdateTosConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
