// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/system"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetInviteConfigOutput struct {
	Body *types.InviteConfig
}

func GetInviteConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetInviteConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetInviteConfigOutput, error) {
		l := system.NewGetInviteConfigLogic(ctx, svcCtx)
		resp, err := l.GetInviteConfig()
		if err != nil {
			return nil, err
		}
		return &GetInviteConfigOutput{Body: resp}, nil
	}
}
