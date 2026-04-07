// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/services/admin/system"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateInviteConfigInput struct {
	Body types.InviteConfig
}

func UpdateInviteConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateInviteConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateInviteConfigInput) (*struct{}, error) {
		l := system.NewUpdateInviteConfigLogic(ctx, svcCtx)
		if err := l.UpdateInviteConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
