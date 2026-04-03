// huma:migrated
package ads

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/ads"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type DeleteAdsInput struct {
	Body types.DeleteAdsRequest
}

func DeleteAdsHandler(svcCtx *svc.ServiceContext) func(context.Context, *DeleteAdsInput) (*struct{}, error) {
	return func(ctx context.Context, input *DeleteAdsInput) (*struct{}, error) {
		l := ads.NewDeleteAdsLogic(ctx, svcCtx)
		if err := l.DeleteAds(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
