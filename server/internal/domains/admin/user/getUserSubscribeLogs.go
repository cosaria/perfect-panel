package user

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/persistence/log"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type GetUserSubscribeLogsInput struct {
	types.GetUserSubscribeLogsRequest
}

type GetUserSubscribeLogsOutput struct {
	Body *types.GetUserSubscribeLogsResponse
}

func GetUserSubscribeLogsHandler(deps Deps) func(context.Context, *GetUserSubscribeLogsInput) (*GetUserSubscribeLogsOutput, error) {
	return func(ctx context.Context, input *GetUserSubscribeLogsInput) (*GetUserSubscribeLogsOutput, error) {
		l := NewGetUserSubscribeLogsLogic(ctx, deps)
		resp, err := l.GetUserSubscribeLogs(&input.GetUserSubscribeLogsRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserSubscribeLogsOutput{Body: resp}, nil
	}
}

type GetUserSubscribeLogsLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get user subcribe logs
func NewGetUserSubscribeLogsLogic(ctx context.Context, deps Deps) *GetUserSubscribeLogsLogic {
	return &GetUserSubscribeLogsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetUserSubscribeLogsLogic) GetUserSubscribeLogs(req *types.GetUserSubscribeLogsRequest) (resp *types.GetUserSubscribeLogsResponse, err error) {
	data, total, err := l.deps.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{})

	if err != nil {
		l.Errorw("[GetUserSubscribeLogs] Get User Subscribe Logs Error:", logger.Field("err", err.Error()))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "Get User Subscribe Logs Error")
	}
	var list []types.UserSubscribeLog
	tool.DeepCopy(&list, data)

	return &types.GetUserSubscribeLogsResponse{
		List:  list,
		Total: total,
	}, err
}
