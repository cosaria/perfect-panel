// huma:migrated
package ads

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/ads"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type CreateAdsInput struct {
	Body types.CreateAdsRequest
}

func CreateAdsHandler(svcCtx *svc.ServiceContext) func(context.Context, *CreateAdsInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateAdsInput) (*struct{}, error) {
		l := ads.NewCreateAdsLogic(ctx, svcCtx)
		if err := l.CreateAds(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
