// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/system"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateSubscribeConfigInput struct {
	Body types.SubscribeConfig
}

func UpdateSubscribeConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateSubscribeConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateSubscribeConfigInput) (*struct{}, error) {
		l := system.NewUpdateSubscribeConfigLogic(ctx, svcCtx)
		if err := l.UpdateSubscribeConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
