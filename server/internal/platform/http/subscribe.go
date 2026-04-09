package handler

import (
	"github.com/gin-gonic/gin"
	appruntime "github.com/perfect-panel/server/internal/bootstrap/runtime"
	servicesubscribe "github.com/perfect-panel/server/internal/domains/subscribe"
)

func RegisterSubscribeHandlers(router *gin.Engine, runtimeDeps *appruntime.Deps) {
	deps := servicesubscribe.Deps{}
	if runtimeDeps != nil {
		deps.ClientModel = runtimeDeps.ClientModel
		deps.LogModel = runtimeDeps.LogModel
		deps.NodeModel = runtimeDeps.NodeModel
		deps.SubscribeModel = runtimeDeps.SubscribeModel
		deps.UserModel = runtimeDeps.UserModel
		deps.Config = runtimeDeps.Config
	}
	servicesubscribe.RegisterSubscribeHandlers(router, deps)
}
