// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/system"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetSiteConfigOutput struct {
	Body *types.SiteConfig
}

func GetSiteConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetSiteConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetSiteConfigOutput, error) {
		l := system.NewGetSiteConfigLogic(ctx, svcCtx)
		resp, err := l.GetSiteConfig()
		if err != nil {
			return nil, err
		}
		return &GetSiteConfigOutput{Body: resp}, nil
	}
}
