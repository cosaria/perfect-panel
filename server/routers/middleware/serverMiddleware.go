package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/perfect-panel/server/svc"
)

func ServerMiddleware(svc *svc.ServiceContext) func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.GetHeader("X-Node-Secret") == svc.Config.Node.NodeSecret && svc.Config.Node.NodeSecret != "" {
			c.Next()
			return
		}
		c.String(http.StatusForbidden, "Forbidden")
		c.Abort()
	}
}
