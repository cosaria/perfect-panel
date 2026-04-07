package ads

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetAdsDetailInput struct {
	types.GetAdsDetailRequest
}

type GetAdsDetailOutput struct {
	Body *types.Ads
}

func GetAdsDetailHandler(deps Deps) func(context.Context, *GetAdsDetailInput) (*GetAdsDetailOutput, error) {
	return func(ctx context.Context, input *GetAdsDetailInput) (*GetAdsDetailOutput, error) {
		l := NewGetAdsDetailLogic(ctx, deps)
		resp, err := l.GetAdsDetail(&input.GetAdsDetailRequest)
		if err != nil {
			return nil, err
		}
		return &GetAdsDetailOutput{Body: resp}, nil
	}
}

type GetAdsDetailLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get Ads Detail
func NewGetAdsDetailLogic(ctx context.Context, deps Deps) *GetAdsDetailLogic {
	return &GetAdsDetailLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetAdsDetailLogic) GetAdsDetail(req *types.GetAdsDetailRequest) (resp *types.Ads, err error) {
	data, err := l.deps.AdsModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Errorw("find ads error", logger.Field("error", err.Error()), logger.Field("id", req.Id))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find ads error: %v", err.Error())
	}
	resp = new(types.Ads)
	tool.DeepCopy(resp, data)
	return
}
