package server

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/models/node"
	"github.com/perfect-panel/server/modules/infra/logger"
	"github.com/perfect-panel/server/routers/response"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
	"time"
)

func ServerPushStatusHandler(svcCtx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.ServerPushStatusRequest
		_ = c.ShouldBind(&req)
		_ = c.ShouldBindQuery(&req.ServerCommon)
		validateErr := svcCtx.Validate(&req)
		if validateErr != nil {
			response.ParamErrorResult(c, validateErr)
			return
		}

		l := NewServerPushStatusLogic(c.Request.Context(), svcCtx)
		err := l.ServerPushStatus(&req)
		response.HttpResult(c, nil, err)
	}
}

type ServerPushStatusLogic struct {
	logger.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewServerPushStatusLogic Push server status
func NewServerPushStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServerPushStatusLogic {
	return &ServerPushStatusLogic{
		Logger: logger.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ServerPushStatusLogic) ServerPushStatus(req *types.ServerPushStatusRequest) error {
	// Find server info
	serverInfo, err := l.svcCtx.NodeModel.FindOneServer(l.ctx, req.ServerId)
	if err != nil || serverInfo.Id <= 0 {
		l.Errorw("[PushOnlineUsers] FindOne error", logger.Field("error", err))
		return errors.New("server not found")
	}
	err = l.svcCtx.NodeModel.UpdateStatusCache(l.ctx, req.ServerId, &node.Status{
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

	err = l.svcCtx.NodeModel.UpdateServer(l.ctx, serverInfo)
	if err != nil {
		l.Errorw("[ServerPushStatus] UpdateServer error", logger.Field("error", err))
		return nil
	}

	return nil
}
