package ads

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetAdsDetailLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Get Ads Detail
func NewGetAdsDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAdsDetailLogic {
	return &GetAdsDetailLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdsDetailLogic) GetAdsDetail(req *types.GetAdsDetailRequest) (resp *types.Ads, err error) {
	data, err := l.svcCtx.AdsModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Errorw("find ads error", logger.Field("error", err.Error()), logger.Field("id", req.Id))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find ads error: %v", err.Error())
	}
	resp = new(types.Ads)
	tool.DeepCopy(resp, data)
	return
}
