package user

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type GetUserSubscribeInput struct {
	types.GetUserSubscribeListRequest
}

type GetUserSubscribeOutput struct {
	Body *types.GetUserSubscribeListResponse
}

func GetUserSubscribeHandler(deps Deps) func(context.Context, *GetUserSubscribeInput) (*GetUserSubscribeOutput, error) {
	return func(ctx context.Context, input *GetUserSubscribeInput) (*GetUserSubscribeOutput, error) {
		l := NewGetUserSubscribeLogic(ctx, deps)
		resp, err := l.GetUserSubscribe(&input.GetUserSubscribeListRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserSubscribeOutput{Body: resp}, nil
	}
}

type GetUserSubscribeLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get user subcribe
func NewGetUserSubscribeLogic(ctx context.Context, deps Deps) *GetUserSubscribeLogic {
	return &GetUserSubscribeLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetUserSubscribeLogic) GetUserSubscribe(req *types.GetUserSubscribeListRequest) (resp *types.GetUserSubscribeListResponse, err error) {
	data, err := l.deps.UserModel.QueryUserSubscribe(l.ctx, req.UserId, 0, 1, 2, 3, 4)
	if err != nil {
		l.Errorw("[GetUserSubscribeLogs] Get User Subscribe Error:", logger.Field("err", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Get User Subscribe Error")
	}

	resp = &types.GetUserSubscribeListResponse{
		List:  make([]types.UserSubscribe, 0),
		Total: int64(len(data)),
	}

	for _, item := range data {
		var sub types.UserSubscribe
		tool.DeepCopy(&sub, item)
		sub.Short, _ = tool.FixedUniqueString(item.Token, 8, "")
		resp.List = append(resp.List, sub)
	}
	return
}
