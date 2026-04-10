package bootstrap

import (
	"log/slog"

	"github.com/perfect-panel/server-v2/internal/app/wiring"
	"github.com/perfect-panel/server-v2/internal/platform/config"
	"github.com/perfect-panel/server-v2/internal/platform/observability"
)

// Options 描述 bootstrap 所需输入。
type Options struct {
	Config config.LoadOptions
	Logger *slog.Logger
}

// Build 执行最小应用装配。
func Build(opts Options) (*wiring.Container, error) {
	cfg, err := config.Load(opts.Config)
	if err != nil {
		return nil, err
	}

	logger := opts.Logger
	if logger == nil {
		logger = observability.NewLogger(cfg.LogLevel)
	}

	return wiring.NewContainer(cfg, logger), nil
}
