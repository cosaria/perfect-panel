package wiring

import (
	"log/slog"
	"net/http"

	"github.com/perfect-panel/server-v2/internal/platform/config"
	"github.com/perfect-panel/server-v2/internal/platform/http/health"
)

// Container 承载应用启动所需的基础依赖。
type Container struct {
	Config        config.Config
	Logger        *slog.Logger
	HealthHandler http.Handler
}

// NewContainer 仅负责依赖装配，不承载业务语义。
func NewContainer(cfg config.Config, logger *slog.Logger) *Container {
	return &Container{
		Config:        cfg,
		Logger:        logger,
		HealthHandler: health.NewHandler(),
	}
}
