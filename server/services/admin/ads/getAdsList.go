package ads

import (
	"context"
	"github.com/perfect-panel/server/models/ads"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
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

type GetAdsListLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Ads List
func NewGetAdsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdsListLogic {
	return &GetAdsListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdsListLogic) GetAdsList(req *types.GetAdsListRequest) (resp *types.GetAdsListResponse, err error) {
	total, data, err := l.svcCtx.AdsModel.GetAdsListByPage(l.ctx, req.Page, req.Size, ads.Filter{
		Search: req.Search,
		Status: req.Status,
	})
	if err != nil {
		l.Errorw("get ads list error", logger.Field("error", err.Error()), logger.Field("req", req))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get ads list error: %v", err.Error())
	}
	resp = &types.GetAdsListResponse{
		Total: total,
		List:  make([]types.Ads, len(data)),
	}
	tool.DeepCopy(&resp.List, data)
	return
}
