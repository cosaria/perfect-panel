package user

import (
	"context"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
	"github.com/pkg/errors"
)

type GetUserSubscribeResetTrafficLogsInput struct {
	types.GetUserSubscribeResetTrafficLogsRequest
}

type GetUserSubscribeResetTrafficLogsOutput struct {
	Body *types.GetUserSubscribeResetTrafficLogsResponse
}

func GetUserSubscribeResetTrafficLogsHandler(deps Deps) func(context.Context, *GetUserSubscribeResetTrafficLogsInput) (*GetUserSubscribeResetTrafficLogsOutput, error) {
	return func(ctx context.Context, input *GetUserSubscribeResetTrafficLogsInput) (*GetUserSubscribeResetTrafficLogsOutput, error) {
		l := NewGetUserSubscribeResetTrafficLogsLogic(ctx, deps)
		resp, err := l.GetUserSubscribeResetTrafficLogs(&input.GetUserSubscribeResetTrafficLogsRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserSubscribeResetTrafficLogsOutput{Body: resp}, nil
	}
}

type GetUserSubscribeResetTrafficLogsLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get user subcribe reset traffic logs
func NewGetUserSubscribeResetTrafficLogsLogic(ctx context.Context, deps Deps) *GetUserSubscribeResetTrafficLogsLogic {
	return &GetUserSubscribeResetTrafficLogsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetUserSubscribeResetTrafficLogsLogic) GetUserSubscribeResetTrafficLogs(req *types.GetUserSubscribeResetTrafficLogsRequest) (resp *types.GetUserSubscribeResetTrafficLogsResponse, err error) {
	data, total, err := l.deps.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:     req.Page,
		Size:     req.Size,
		Type:     log.TypeResetSubscribe.Uint8(),
		ObjectID: req.UserSubscribeId,
	})
	if err != nil {
		l.Errorf("[ResetSubscribeTrafficLog] failed to filter system log: %v", err)
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "FilterSystemLog failed, err: %v", err)
	}

	var list []types.ResetSubscribeTrafficLog

	for _, item := range data {
		var content log.ResetSubscribe
		if err = content.Unmarshal([]byte(item.Content)); err != nil {
			l.Errorf("[ResetSubscribeTrafficLog] failed to unmarshal log: %v", err)
			continue
		}
		list = append(list, types.ResetSubscribeTrafficLog{
			Id:              item.Id,
			Type:            content.Type,
			OrderNo:         content.OrderNo,
			Timestamp:       content.Timestamp,
			UserSubscribeId: item.ObjectID,
		})
	}

	return &types.GetUserSubscribeResetTrafficLogsResponse{
		Total: total,
		List:  list,
	}, nil
}
