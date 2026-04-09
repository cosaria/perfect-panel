package server

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	task "github.com/perfect-panel/server/internal/jobs"
	modelnode "github.com/perfect-panel/server/internal/platform/persistence/node"
	"github.com/perfect-panel/server/internal/platform/http/response"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/internal/platform/support/logger"
	"github.com/perfect-panel/server/internal/platform/support/tool"
	"github.com/pkg/errors"
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
	rawPayload, err := json.Marshal(request)
	if err != nil {
		return err
	}
	ingestRepo := modelnode.NewUsageIngestRepository(l.deps.DB)
	if ingestRepo != nil && ingestRepo.Available() {
		sum := sha256.Sum256(rawPayload)
		decision, err := ingestRepo.Ingest(l.ctx, &modelnode.UsageIngestInput{
			ServerID:       serverInfo.Id,
			Protocol:       req.Protocol,
			IdempotencyKey: hex.EncodeToString(sum[:]),
			AuthStatus:     "verified",
			RawPayload:     string(rawPayload),
			LogCount:       len(request.Logs),
		})
		if err != nil {
			return err
		}
		if decision != nil {
			if !decision.Accepted {
				l.Infow("[ServerPushUserTraffic] Duplicate traffic report ignored",
					logger.Field("server_id", serverInfo.Id),
					logger.Field("report_id", decision.ReportID),
					logger.Field("state", decision.ProcessingState))
				return nil
			}
			request.ReportID = decision.ReportID
		}
	}

	// Push traffic task
	val := rawPayload
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
