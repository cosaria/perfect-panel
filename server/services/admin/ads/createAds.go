package ads

import (
	"context"
	"github.com/perfect-panel/server/models/ads"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
	"time"
)

type CreateAdsInput struct {
	Body types.CreateAdsRequest
}

func CreateAdsHandler(deps Deps) func(context.Context, *CreateAdsInput) (*struct{}, error) {
	return func(ctx context.Context, input *CreateAdsInput) (*struct{}, error) {
		l := NewCreateAdsLogic(ctx, deps)
		if err := l.CreateAds(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type CreateAdsLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Create Ads
func NewCreateAdsLogic(ctx context.Context, deps Deps) *CreateAdsLogic {
	return &CreateAdsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *CreateAdsLogic) CreateAds(req *types.CreateAdsRequest) error {
	if err := l.deps.AdsModel.Insert(l.ctx, &ads.Ads{
		Title:     req.Title,
		Type:      req.Type,
		Content:   req.Content,
		TargetURL: req.TargetURL,
		StartTime: time.UnixMilli(req.StartTime),
		EndTime:   time.UnixMilli(req.EndTime),
		Status:    req.Status,
	}); err != nil {
		l.Errorw("insert ads error: %v", logger.Field("error", err.Error()), logger.Field("req", req))
		return errors.Wrapf(xerr.NewErrCode(xerr.DatabaseInsertError), "insert ads error: %v", err.Error())
	}
	return nil
}
