// huma:migrated
package ads

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type UpdateAdsInput struct {
	Body types.UpdateAdsRequest
}

func UpdateAdsHandler(svcCtx *svc.ServiceContext) func(context.Context, *UpdateAdsInput) (*struct{}, error) {
	return func(ctx context.Context, input *UpdateAdsInput) (*struct{}, error) {
		l := NewUpdateAdsLogic(ctx, svcCtx)
		if err := l.UpdateAds(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
