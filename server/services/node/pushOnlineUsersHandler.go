package server

import (
	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/routers/response"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

// Push online users
func PushOnlineUsersHandler(svcCtx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.OnlineUsersRequest
		_ = c.ShouldBind(&req)
		_ = c.ShouldBindQuery(&req.ServerCommon)
		validateErr := svcCtx.Validate(&req)
		if validateErr != nil {
			response.ParamErrorResult(c, validateErr)
			return
		}

		l := NewPushOnlineUsersLogic(c.Request.Context(), svcCtx)
		err := l.PushOnlineUsers(&req)
		response.HttpResult(c, nil, err)
	}
}
