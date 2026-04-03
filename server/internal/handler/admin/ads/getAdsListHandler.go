// huma:migrated
package ads

import (
	"context"
	"github.com/perfect-panel/server/internal/logic/admin/ads"
	"github.com/perfect-panel/server/internal/svc"
	"github.com/perfect-panel/server/internal/types"
)

type GetAdsListInput struct {
	Body types.GetAdsListRequest
}

type GetAdsListOutput struct {
	Body *types.GetAdsListResponse
}

func GetAdsListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetAdsListInput) (*GetAdsListOutput, error) {
	return func(ctx context.Context, input *GetAdsListInput) (*GetAdsListOutput, error) {
		l := ads.NewGetAdsListLogic(ctx, svcCtx)
		resp, err := l.GetAdsList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetAdsListOutput{Body: resp}, nil
	}
}
