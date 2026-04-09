package portal

import (
	"context"
	"encoding/json"
	"github.com/perfect-panel/server/models/subscribe"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type GetSubscriptionInput struct {
	types.GetSubscriptionRequest
}

type GetSubscriptionOutput struct {
	Body *types.GetSubscriptionResponse
}

func GetSubscriptionHandler(deps Deps) func(context.Context, *GetSubscriptionInput) (*GetSubscriptionOutput, error) {
	return func(ctx context.Context, input *GetSubscriptionInput) (*GetSubscriptionOutput, error) {
		l := NewGetSubscriptionLogic(ctx, deps)
		resp, err := l.GetSubscription(&input.GetSubscriptionRequest)
		if err != nil {
			return nil, err
		}
		return &GetSubscriptionOutput{Body: resp}, nil
	}
}

type GetSubscriptionLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewGetSubscriptionLogic Get Subscription
func NewGetSubscriptionLogic(ctx context.Context, deps Deps) *GetSubscriptionLogic {
	return &GetSubscriptionLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetSubscriptionLogic) GetSubscription(req *types.GetSubscriptionRequest) (resp *types.GetSubscriptionResponse, err error) {
	resp = &types.GetSubscriptionResponse{
		List: make([]types.Subscribe, 0),
	}
	// Get the subscription list
	_, data, err := l.deps.SubscribeModel.FilterList(l.ctx, &subscribe.FilterParams{
		Page:            1,
		Size:            9999,
		Show:            true,
		Language:        req.Language,
		DefaultLanguage: true,
	})
	if err != nil {
		l.Errorw("[Site GetSubscription]", logger.Field("err", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get subscription list error: %v", err.Error())
	}
	list := make([]types.Subscribe, len(data))
	for i, item := range data {
		var sub types.Subscribe
		tool.DeepCopy(&sub, item)
		if item.Discount != "" {
			var discount []types.SubscribeDiscount
			_ = json.Unmarshal([]byte(item.Discount), &discount)
			sub.Discount = discount
			list[i] = sub
		}
		list[i] = sub
	}
	resp.List = list
	return
}
