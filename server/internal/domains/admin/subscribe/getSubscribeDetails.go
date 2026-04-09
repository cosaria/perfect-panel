package subscribe

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type GetSubscribeDetailsInput struct {
	types.GetSubscribeDetailsRequest
}

type GetSubscribeDetailsOutput struct {
	Body *types.Subscribe
}

func GetSubscribeDetailsHandler(deps Deps) func(context.Context, *GetSubscribeDetailsInput) (*GetSubscribeDetailsOutput, error) {
	return func(ctx context.Context, input *GetSubscribeDetailsInput) (*GetSubscribeDetailsOutput, error) {
		l := NewGetSubscribeDetailsLogic(ctx, deps)
		resp, err := l.GetSubscribeDetails(&input.GetSubscribeDetailsRequest)
		if err != nil {
			return nil, err
		}
		return &GetSubscribeDetailsOutput{Body: resp}, nil
	}
}

type GetSubscribeDetailsLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get subscribe details
func NewGetSubscribeDetailsLogic(ctx context.Context, deps Deps) *GetSubscribeDetailsLogic {
	return &GetSubscribeDetailsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetSubscribeDetailsLogic) GetSubscribeDetails(req *types.GetSubscribeDetailsRequest) (resp *types.Subscribe, err error) {
	sub, err := l.deps.SubscribeModel.FindOne(l.ctx, req.Id)
	if err != nil {
		l.Error("[GetSubscribeDetailsLogic] get subscribe details failed: ", logger.Field("error", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get subscribe details failed: %v", err.Error())
	}
	resp = &types.Subscribe{}
	tool.DeepCopy(resp, sub)
	if sub.Discount != "" {
		err = json.Unmarshal([]byte(sub.Discount), &resp.Discount)
		if err != nil {
			l.Error("[GetSubscribeDetailsLogic] JSON unmarshal failed: ", logger.Field("error", err.Error()), logger.Field("discount", sub.Discount))
		}
	}
	resp.Nodes = tool.StringToInt64Slice(sub.Nodes)
	resp.NodeTags = strings.Split(sub.NodeTags, ",")
	return resp, nil
}
