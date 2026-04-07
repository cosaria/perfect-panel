package server

import (
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/pkg/result"
	"github.com/perfect-panel/server/services/node"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

// Push server status
func ServerPushStatusHandler(svcCtx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.ServerPushStatusRequest
		_ = c.ShouldBind(&req)
		_ = c.ShouldBindQuery(&req.ServerCommon)
		validateErr := svcCtx.Validate(&req)
		if validateErr != nil {
			result.ParamErrorResult(c, validateErr)
			return
		}

		l := server.NewServerPushStatusLogic(c.Request.Context(), svcCtx)
		err := l.ServerPushStatus(&req)
		result.HttpResult(c, nil, err)
	}
}
