package marketing

import (
	"context"

	"github.com/perfect-panel/server/models/task"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/infra/xerr"
	"github.com/perfect-panel/server/modules/notify/email"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

type StopBatchSendEmailTaskLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewStopBatchSendEmailTaskLogic Stop a batch send email task
func NewStopBatchSendEmailTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StopBatchSendEmailTaskLogic {
	return &StopBatchSendEmailTaskLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StopBatchSendEmailTaskLogic) StopBatchSendEmailTask(req *types.StopBatchSendEmailTaskRequest) (err error) {
	if email.Manager != nil {
		email.Manager.RemoveWorker(req.Id)
	} else {
		logger.Error("[StopBatchSendEmailTaskLogic] email.Manager is nil, cannot stop task")
	}
	err = l.svcCtx.DB.Model(&task.Task{}).Where("id = ?", req.Id).Update("status", 2).Error

	if err != nil {
		l.Errorf("failed to stop email task, error: %v", err)
		return xerr.NewErrCode(xerr.DatabaseUpdateError)
	}
	return
}
