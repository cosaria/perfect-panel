// huma:migrated
package ads

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/ads"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type UpdateAdsInput struct {
	Body types.UpdateAdsRequest
}

func UpdateAdsHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateAdsInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateAdsInput) (*struct{}, error) {
		l := ads.NewUpdateAdsLogic(ctx, svcCtx)
		if err := l.UpdateAds(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
