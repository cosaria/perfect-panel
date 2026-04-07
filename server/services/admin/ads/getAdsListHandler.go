// huma:migrated
package ads

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetAdsListInput struct {
	Body types.GetAdsListRequest
}

type GetAdsListOutput struct {
	Body *types.GetAdsListResponse
}

func GetAdsListHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetAdsListInput) (*GetAdsListOutput, error) {
	return func(ctx context.Context, input *GetAdsListInput) (*GetAdsListOutput, error) {
		l := NewGetAdsListLogic(ctx, svcCtx)
		resp, err := l.GetAdsList(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetAdsListOutput{Body: resp}, nil
	}
}
