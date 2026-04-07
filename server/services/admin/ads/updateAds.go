package ads

import (
	"context"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"time"
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

type UpdateAdsLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Update Ads
func NewUpdateAdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAdsLogic {
	return &UpdateAdsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAdsLogic) UpdateAds(req *types.UpdateAdsRequest) error {
	data, err := l.svcCtx.AdsModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Errorw("find ads error", logger.Field("error", err.Error()), logger.Field("id", req.Id))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "find ads error: %v", err.Error())
	}
	tool.DeepCopy(data, req)
	data.StartTime = time.UnixMilli(req.StartTime)
	data.EndTime = time.UnixMilli(req.EndTime)
	if err := l.svcCtx.AdsModel.Update(l.ctx, data); err != nil {
		l.Errorw("update ads error", logger.Field("error", err.Error()), logger.Field("req", req))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseUpdateError), "update ads error: %v", err.Error())
	}
	return nil
}
