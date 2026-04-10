package bootstrap

import (
	"context"
	"log/slog"

	"github.com/perfect-panel/server-v2/internal/app/runtime"
	"github.com/perfect-panel/server-v2/internal/app/wiring"
	"github.com/perfect-panel/server-v2/internal/platform/config"
	ppdb "github.com/perfect-panel/server-v2/internal/platform/db"
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

// BuildForMode 在 Build 的基础上，为 serve 模式增加 schema version 门禁。
func BuildForMode(mode string, opts Options) (*wiring.Container, error) {
	container, err := Build(opts)
	if err != nil {
		return nil, err
	}
	if !isServeMode(mode) {
		return container, nil
	}

	database, err := ppdb.OpenFromEnv()
	if err != nil {
		return nil, err
	}
	defer database.Close()

	if err := ppdb.EnsureSchemaVersion(context.Background(), database); err != nil {
		return nil, err
	}
	return container, nil
}

func isServeMode(mode string) bool {
	switch mode {
	case runtime.ModeServeAPI, runtime.ModeServeWorker, runtime.ModeServeScheduler:
		return true
	default:
		return false
	}
}
