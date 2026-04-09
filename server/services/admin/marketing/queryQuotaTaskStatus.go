package marketing

import (
	"context"
	"github.com/perfect-panel/server/models/task"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/pkg/errors"
)

type QueryQuotaTaskStatusInput struct {
	Body types.QueryQuotaTaskStatusRequest
}

type QueryQuotaTaskStatusOutput struct {
	Body *types.QueryQuotaTaskStatusResponse
}

func QueryQuotaTaskStatusHandler(deps Deps) func(context.Context, *QueryQuotaTaskStatusInput) (*QueryQuotaTaskStatusOutput, error) {
	return func(ctx context.Context, input *QueryQuotaTaskStatusInput) (*QueryQuotaTaskStatusOutput, error) {
		l := NewQueryQuotaTaskStatusLogic(ctx, deps)
		resp, err := l.QueryQuotaTaskStatus(&input.Body)
		if err != nil {
			return nil, err
		}
		return &QueryQuotaTaskStatusOutput{Body: resp}, nil
	}
}

type QueryQuotaTaskStatusLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewQueryQuotaTaskStatusLogic Query quota task status
func NewQueryQuotaTaskStatusLogic(ctx context.Context, deps Deps) *QueryQuotaTaskStatusLogic {
	return &QueryQuotaTaskStatusLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *QueryQuotaTaskStatusLogic) QueryQuotaTaskStatus(req *types.QueryQuotaTaskStatusRequest) (resp *types.QueryQuotaTaskStatusResponse, err error) {
	var data *task.Task
	err = l.deps.DB.Model(&task.Task{}).Where("id = ? AND `type` = ?", req.Id, task.TypeQuota).First(&data).Error
	if err != nil {
		l.Errorf("[QueryQuotaTaskStatus] failed to get quota task: %v", err.Error())
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DatabaseQueryError), " failed to get quota task: %v", err.Error())
	}
	return &types.QueryQuotaTaskStatusResponse{
		Status:  uint8(data.Status),
		Current: int64(data.Current),
		Total:   int64(data.Total),
		Errors:  data.Errors,
	}, nil
}
