package oauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/routers/response"
	"github.com/perfect-panel/server/svc"
	"github.com/perfect-panel/server/types"
)

// AppleLoginCallbackHandler Apple Login Callback
// This handler is kept as a raw Gin handler because the logic layer needs
// direct access to http.Request and http.ResponseWriter for HTTP redirects.
func AppleLoginCallbackHandler(svcCtx *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req types.AppleLoginCallbackRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}
		l := NewAppleLoginCallbackLogic(c.Request.Context(), svcCtx)
		err := l.AppleLoginCallback(&req, c.Request, c.Writer)
		if err != nil {
			response.HttpResult(c, nil, err)
		}
	}
}
