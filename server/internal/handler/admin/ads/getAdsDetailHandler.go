// huma:migrated
package ads

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/ads"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetAdsDetailInput struct {
	types.GetAdsDetailRequest
}

type GetAdsDetailOutput struct {
	Body *types.Ads
}

func GetAdsDetailHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetAdsDetailInput) (*GetAdsDetailOutput, error) {
	return func(ctx context.Context, input *GetAdsDetailInput) (*GetAdsDetailOutput, error) {
		l := ads.NewGetAdsDetailLogic(ctx, svcCtx)
		resp, err := l.GetAdsDetail(&input.GetAdsDetailRequest)
		if err != nil {
			return nil, err
		}
		return &GetAdsDetailOutput{Body: resp}, nil
	}
}
