package user

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/log"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/pkg/errors"
)

type GetUserLoginLogsInput struct {
	types.GetUserLoginLogsRequest
}

type GetUserLoginLogsOutput struct {
	Body *types.GetUserLoginLogsResponse
}

func GetUserLoginLogsHandler(deps Deps) func(context.Context, *GetUserLoginLogsInput) (*GetUserLoginLogsOutput, error) {
	return func(ctx context.Context, input *GetUserLoginLogsInput) (*GetUserLoginLogsOutput, error) {
		l := NewGetUserLoginLogsLogic(ctx, deps)
		resp, err := l.GetUserLoginLogs(&input.GetUserLoginLogsRequest)
		if err != nil {
			return nil, err
		}
		return &GetUserLoginLogsOutput{Body: resp}, nil
	}
}

type GetUserLoginLogsLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// Get user login logs
func NewGetUserLoginLogsLogic(ctx context.Context, deps Deps) *GetUserLoginLogsLogic {
	return &GetUserLoginLogsLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetUserLoginLogsLogic) GetUserLoginLogs(req *types.GetUserLoginLogsRequest) (resp *types.GetUserLoginLogsResponse, err error) {
	data, total, err := l.deps.LogModel.FilterSystemLog(l.ctx, &log.FilterParams{
		Page:     req.Page,
		Size:     req.Size,
		Type:     log.TypeLogin.Uint8(),
		ObjectID: req.UserId,
	})
	if err != nil {
		l.Errorw("[GetUserLoginLogs] get user login logs failed", logger.Field("error", err.Error()), logger.Field("request", req))
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), "get user login logs failed: %v", err.Error())
	}
	var list []types.UserLoginLog

	for _, datum := range data {
		var content log.Login
		if err = content.Unmarshal([]byte(datum.Content)); err != nil {
			l.Errorf("[GetUserLoginLogs] unmarshal login log content failed: %v", err.Error())
			continue
		}
		list = append(list, types.UserLoginLog{
			Id:        datum.Id,
			UserId:    datum.ObjectID,
			LoginIP:   content.LoginIP,
			UserAgent: content.UserAgent,
			Success:   content.Success,
			Timestamp: datum.CreatedAt.UnixMilli(),
		})
	}

	return &types.GetUserLoginLogsResponse{
		Total: total,
		List:  list,
	}, nil
}
