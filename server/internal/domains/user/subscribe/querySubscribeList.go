package subscribe

import (
	"context"
	"encoding/json"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/subscribe"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type QuerySubscribeListInput struct {
	types.QuerySubscribeListRequest
}

type QuerySubscribeListOutput struct {
	Body *types.QuerySubscribeListResponse
}

func QuerySubscribeListHandler(deps Deps) func(context.Context, *QuerySubscribeListInput) (*QuerySubscribeListOutput, error) {
	return func(ctx context.Context, input *QuerySubscribeListInput) (*QuerySubscribeListOutput, error) {
		l := NewQuerySubscribeListLogic(ctx, deps)
		resp, err := l.QuerySubscribeList(&input.QuerySubscribeListRequest)
		if err != nil {
			return nil, err
		}
		return &QuerySubscribeListOutput{Body: resp}, nil
	}
}

type QuerySubscribeListLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get subscribe list
func NewQuerySubscribeListLogic(ctx context.Context, deps Deps) *QuerySubscribeListLogic {
	return &QuerySubscribeListLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QuerySubscribeListLogic) QuerySubscribeList(req *types.QuerySubscribeListRequest) (resp *types.QuerySubscribeListResponse, err error) {

	total, data, err := l.deps.SubscribeModel.FilterList(l.ctx, &subscribe.FilterParams{
		Page:            1,
		Size:            9999,
		Language:        req.Language,
		Sell:            true,
		DefaultLanguage: true,
	})
	if err != nil {
		l.Errorw("[QuerySubscribeListLogic] Database Error", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "QuerySubscribeList error: %v", err.Error())
	}

	resp = &types.QuerySubscribeListResponse{
		Total: total,
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
