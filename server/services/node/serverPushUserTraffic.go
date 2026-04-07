package server

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/modules/util/tool"
	"github.com/perfect-panel/server/routers/response"
	"github.com/perfect-panel/server/types"
	task "github.com/perfect-panel/server/worker"
	"github.com/pkg/errors"
	"time"
)

func ServerPushUserTrafficHandler(deps Deps) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.ServerPushUserTrafficRequest
		_ = c.ShouldBind(&req)
		_ = c.ShouldBindQuery(&req.ServerCommon)
		validateErr := validateRequest(&req)
		if validateErr != nil {
			response.ParamErrorResult(c, validateErr)
			return
		}

		l := NewServerPushUserTrafficLogic(c.Request.Context(), deps)
		err := l.ServerPushUserTraffic(&req)
		response.HttpResult(c, nil, err)
	}
}

type ServerPushUserTrafficLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewServerPushUserTrafficLogic Push user Traffic
func NewServerPushUserTrafficLogic(ctx context.Context, deps Deps) *ServerPushUserTrafficLogic {
	return &ServerPushUserTrafficLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *ServerPushUserTrafficLogic) ServerPushUserTraffic(req *types.ServerPushUserTrafficRequest) error {
	// Find server info
	serverInfo, err := l.deps.NodeModel.FindOneServer(l.ctx, req.ServerId)
	if err != nil {
		l.Errorw("[PushOnlineUsers] FindOne error", logger.Field("error", err))
		return errors.New("server not found")
	}

	// Create traffic task
	var request task.TrafficStatistics
	request.ServerId = serverInfo.Id
	request.Protocol = req.Protocol
	tool.DeepCopy(&request.Logs, req.Traffic)

	// Push traffic task
	val, _ := json.Marshal(request)
	t := asynq.NewTask(task.ForthwithTrafficStatistics, val, asynq.MaxRetry(3))
	info, err := l.deps.Queue.EnqueueContext(l.ctx, t)
	if err != nil {
		l.Errorw("[ServerPushUserTraffic] Push traffic task error", logger.Field("error", err.Error()), logger.Field("task", t))
	} else {
		l.Infow("[ServerPushUserTraffic] Push traffic task success", logger.Field("task", t.Type()), logger.Field("info", string(info.Payload)))
	}

	// Update server last reported time
	now := time.Now()
	serverInfo.LastReportedAt = &now

	err = l.deps.NodeModel.UpdateServer(l.ctx, serverInfo)
	if err != nil {
		l.Errorw("[ServerPushUserTraffic] UpdateServer error", logger.Field("error", err))
		return nil
	}
	return nil
}
