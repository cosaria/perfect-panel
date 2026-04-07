// huma:migrated
package system

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetSiteConfigOutput struct {
	Body *types.SiteConfig
}

func GetSiteConfigHandler(svcCtx *svc.ServiceContext) func(context.Context, *struct{}) (*GetSiteConfigOutput, error) {
	return func(ctx context.Context, _ *struct{}) (*GetSiteConfigOutput, error) {
		l := NewGetSiteConfigLogic(ctx, svcCtx)
		resp, err := l.GetSiteConfig()
		if err != nil {
			return nil, err
		}
		return &GetSiteConfigOutput{Body: resp}, nil
	}
}
