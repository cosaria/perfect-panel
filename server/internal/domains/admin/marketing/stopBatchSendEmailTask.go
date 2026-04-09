package marketing

import (
	"context"

	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/notify/email"
	"github.com/perfect-panel/server/internal/platform/persistence/task"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/xerr"
)

type StopBatchSendEmailTaskInput struct {
	Body types.StopBatchSendEmailTaskRequest
}

func StopBatchSendEmailTaskHandler(deps Deps) func(context.Context, *StopBatchSendEmailTaskInput) (*struct{}, error) {
	return func(ctx context.Context, input *StopBatchSendEmailTaskInput) (*struct{}, error) {
		l := NewStopBatchSendEmailTaskLogic(ctx, deps)
		if err := l.StopBatchSendEmailTask(&input.Body); err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type StopBatchSendEmailTaskLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewStopBatchSendEmailTaskLogic Stop a batch send email task
func NewStopBatchSendEmailTaskLogic(ctx context.Context, deps Deps) *StopBatchSendEmailTaskLogic {
	return &StopBatchSendEmailTaskLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *StopBatchSendEmailTaskLogic) StopBatchSendEmailTask(req *types.StopBatchSendEmailTaskRequest) (err error) {
	if email.Manager != nil {
		email.Manager.RemoveWorker(req.Id)
	} else {
		logger.Error("[StopBatchSendEmailTaskLogic] email.Manager is nil, cannot stop task")
	}
	err = l.deps.DB.Model(&task.Task{}).Where("id = ?", req.Id).Update("status", 2).Error

	if err != nil {
		l.Errorf("failed to stop email task, error: %v", err)
		return xerr.NewErrCode(xerr.DatabaseUpdateError)
	}
	return
}
