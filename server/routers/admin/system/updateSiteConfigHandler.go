// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/services/admin/system"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateSiteConfigInput struct {
	Body types.SiteConfig
}

func UpdateSiteConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateSiteConfigInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateSiteConfigInput) (*struct{}, error) {
		l := system.NewUpdateSiteConfigLogic(ctx, svcCtx)
		if err := l.UpdateSiteConfig(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
