package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appruntime "github.com/perfect-panel/server/internal/bootstrap/runtime"
)

func ServerMiddleware(runtimeDeps *appruntime.Deps) func(c *gin.Context) {
	return func(c *gin.Context) {
		if runtimeDeps != nil && runtimeDeps.Config != nil && c.GetHeader("X-Node-Secret") == runtimeDeps.Config.Node.NodeSecret && runtimeDeps.Config.Node.NodeSecret != "" {
			c.Next()
			return
		}
		c.String(http.StatusForbidden, "Forbidden")
		c.Abort()
	}
}
