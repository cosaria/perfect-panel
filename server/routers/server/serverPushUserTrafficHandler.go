package server

import (
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/routers/response"
	"github.com/perfect-panel/server/services/node"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

// ServerPushUserTrafficHandler Push user Traffic
func ServerPushUserTrafficHandler(svcCtx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.ServerPushUserTrafficRequest
		_ = c.ShouldBind(&req)
		_ = c.ShouldBindQuery(&req.ServerCommon)
		validateErr := svcCtx.Validate(&req)
		if validateErr != nil {
			response.ParamErrorResult(c, validateErr)
			return
		}

		l := server.NewServerPushUserTrafficLogic(c.Request.Context(), svcCtx)
		err := l.ServerPushUserTraffic(&req)
		response.HttpResult(c, nil, err)
	}
}
