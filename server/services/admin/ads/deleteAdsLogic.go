package ads

import (
	"context"

	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type DeleteAdsLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Delete Ads
func NewDeleteAdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAdsLogic {
	return &DeleteAdsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAdsLogic) DeleteAds(req *types.DeleteAdsRequest) error {
	if err := l.svcCtx.AdsModel.Delete(l.ctx, req.Id); err != nil {
		l.Errorw("delete ads error", logger.Field("error", err.Error()), logger.Field("id", req.Id))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseDeletedError), "delete ads error: %v", err.Error())
	}
	return nil
}
