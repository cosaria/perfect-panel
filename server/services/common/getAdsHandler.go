// huma:migrated
package common

import (
	"context"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type GetAdsInput struct {
	types.GetAdsRequest
}

type GetAdsOutput struct {
	Body *types.GetAdsResponse
}

func GetAdsHandler(svcCtx *svc.ServiceContext) func(context.Context, *GetAdsInput) (*GetAdsOutput, error) {
	return func(ctx context.Context, input *GetAdsInput) (*GetAdsOutput, error) {
		l := NewGetAdsLogic(ctx, svcCtx)
		resp, err := l.GetAds(&input.GetAdsRequest)
		if err != nil {
			return nil, err
		}
		return &GetAdsOutput{Body: resp}, nil
	}
}
