package user

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
	"github.com/pkg/errors"
)

type GetUserSubscribeTrafficLogsInput struct {
	types.GetUserSubscribeTrafficLogsRequest
}

type GetUserSubscribeTrafficLogsOutput struct {
	Body *types.GetUserSubscribeTrafficLogsResponse
}

func GetUserSubscribeTrafficLogsHandler(deps Deps) func(context.Context, *GetUserSubscribeTrafficLogsInput) (*GetUserSubscribeTrafficLogsOutput, error) {
	return func(ctx context.Context, input *GetUserSubscribeTrafficLogsInput) (*GetUserSubscribeTrafficLogsOutput, error) {
		l := NewGetUserSubscribeTrafficLogsLogic(ctx, deps)
		resp, err := l.GetUserSubscribeTrafficLogs(&input.GetUserSubscribeTrafficLogsRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserSubscribeTrafficLogsOutput{Body: resp}, nil
	}
}

type GetUserSubscribeTrafficLogsLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get user subcribe traffic logs
func NewGetUserSubscribeTrafficLogsLogic(ctx context.Context, deps Deps) *GetUserSubscribeTrafficLogsLogic {
	return &GetUserSubscribeTrafficLogsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetUserSubscribeTrafficLogsLogic) GetUserSubscribeTrafficLogs(req *types.GetUserSubscribeTrafficLogsRequest) (resp *types.GetUserSubscribeTrafficLogsResponse, err error) {
	list, total, err := l.deps.TrafficLogModel.QueryTrafficLogPageList(l.ctx, req.UserId, req.SubscribeId, req.Page, req.Size)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "GetUserSubscribeTrafficLogs failed: %v", err.Error())
	}
	userRespList := make([]types.TrafficLog, 0)
	tool.DeepCopy(&userRespList, list)
	return &types.GetUserSubscribeTrafficLogsResponse{
		Total: total,
		List:  userRespList,
	}, nil
}
