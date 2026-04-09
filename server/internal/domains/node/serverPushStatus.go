package server

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/internal/platform/http/response"
	"github.com/perfect-panel/server/internal/platform/http/types"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/modules/infra/logger"
)

func ServerPushStatusHandler(deps Deps) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.ServerPushStatusRequest
		_ = c.ShouldBind(&req)
		_ = c.ShouldBindQuery(&req.ServerCommon)
		validateErr := validateRequest(&req)
		if validateErr != nil {
			response.ParamErrorResult(c, validateErr)
			return
		}

		l := NewServerPushStatusLogic(c.Request.Context(), deps)
		err := l.ServerPushStatus(&req)
		response.HttpResult(c, nil, err)
	}
}

type ServerPushStatusLogic struct {
	logger.Logger
	ctx  context.Context
	deps Deps
}

// NewServerPushStatusLogic Push server status
func NewServerPushStatusLogic(ctx context.Context, deps Deps) *ServerPushStatusLogic {
	return &ServerPushStatusLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		deps:   deps,
	}
}

func (l *ServerPushStatusLogic) ServerPushStatus(req *types.ServerPushStatusRequest) error {
	// Find server info
	serverInfo, err := l.deps.NodeModel.FindOneServer(l.ctx, req.ServerId)
	if err != nil || serverInfo.Id <= 0 {
		l.Errorw("[PushOnlineUsers] FindOne error", logger.Field("error", err))
		return errors.New("server not found")
	}
	err = l.deps.NodeModel.UpdateStatusCache(l.ctx, req.ServerId, &node.Status{
		Cpu:       req.Cpu,
		Mem:       req.Mem,
		Disk:      req.Disk,
		UpdatedAt: req.UpdatedAt,
	})
	if err != nil {
		l.Errorw("[ServerPushStatus] UpdateNodeStatus error", logger.Field("error", err))
		return errors.New("update node status failed")
	}
	now := time.Now()
	serverInfo.LastReportedAt = &now

	err = l.deps.NodeModel.UpdateServer(l.ctx, serverInfo)
	if err != nil {
		l.Errorw("[ServerPushStatus] UpdateServer error", logger.Field("error", err))
		return nil
	}

	return nil
}
