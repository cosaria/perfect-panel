package common

import (
	"context"
	"github.com/perfect-panel/server/models/ads"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
)

type GetAdsInput struct {
	types.GetAdsRequest
}

type GetAdsOutput struct {
	Body *types.GetAdsResponse
}

func GetAdsHandler(deps Deps) func(context.Context, *GetAdsInput) (*GetAdsOutput, error) {
	return func(ctx context.Context, input *GetAdsInput) (*GetAdsOutput, error) {
		l := NewGetAdsLogic(ctx, deps)
		resp, err := l.GetAds(&input.GetAdsRequest)
		if err != nil {
			return nil, err
		}
		return &GetAdsOutput{Body: resp}, nil
	}
}

type GetAdsLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get Ads
func NewGetAdsLogic(ctx context.Context, deps Deps) *GetAdsLogic {
	return &GetAdsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetAdsLogic) GetAds(req *types.GetAdsRequest) (resp *types.GetAdsResponse, err error) {
	// todo: add ads position and device
	status := 1
	_, data, err := l.deps.AdsModel.GetAdsListByPage(l.ctx, 1, 200, ads.Filter{
		Status: &status,
	})
	if err != nil {
		return nil, err
	}
	resp = &types.GetAdsResponse{
		List: make([]types.Ads, len(data)),
	}
	tool.DeepCopy(&resp.List, data)
	return
}
