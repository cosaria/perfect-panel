package marketing

import (
	"context"
	"github.com/perfect-panel/server/models/task"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/types"
)

type GetBatchSendEmailTaskStatusInput struct {
	Body types.GetBatchSendEmailTaskStatusRequest
}

type GetBatchSendEmailTaskStatusOutput struct {
	Body *types.GetBatchSendEmailTaskStatusResponse
}

func GetBatchSendEmailTaskStatusHandler(deps Deps) func(context.Context, *GetBatchSendEmailTaskStatusInput) (*GetBatchSendEmailTaskStatusOutput, error) {
	return func(ctx context.Context, input *GetBatchSendEmailTaskStatusInput) (*GetBatchSendEmailTaskStatusOutput, error) {
		l := NewGetBatchSendEmailTaskStatusLogic(ctx, deps)
		resp, err := l.GetBatchSendEmailTaskStatus(&input.Body)
		if err != nil {
			return nil, err
		}
		return &GetBatchSendEmailTaskStatusOutput{Body: resp}, nil
	}
}

type GetBatchSendEmailTaskStatusLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewGetBatchSendEmailTaskStatusLogic Get batch send email task status
func NewGetBatchSendEmailTaskStatusLogic(ctx context.Context, deps Deps) *GetBatchSendEmailTaskStatusLogic {
	return &GetBatchSendEmailTaskStatusLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *GetBatchSendEmailTaskStatusLogic) GetBatchSendEmailTaskStatus(req *types.GetBatchSendEmailTaskStatusRequest) (resp *types.GetBatchSendEmailTaskStatusResponse, err error) {
	tx := l.deps.DB

	var taskInfo *task.Task
	err = tx.Model(&task.Task{}).Where("id = ?", req.Id).First(&taskInfo).Error
	if err != nil {
		l.Errorf("failed to get email task status, error: %v", err)
		return nil, xerr.NewErrCode(xerr.DatabaseQueryError)
	}

	return &types.GetBatchSendEmailTaskStatusResponse{
		Status:  uint8(taskInfo.Status),
		Total:   int64(taskInfo.Total),
		Current: int64(taskInfo.Current),
		Errors:  taskInfo.Errors,
	}, nil
}
