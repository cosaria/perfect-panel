package bootstrap_test

import (
	"testing"

	"github.com/perfect-panel/server-v2/internal/platform/config"
	"github.com/perfect-panel/server-v2/tests/support"
)

func TestLoadPrefersCLIOverEnvAndFile(t *testing.T) {
	t.Setenv("PPANEL_HTTP_ADDR", ":7001")
	t.Setenv("PPANEL_LOG_LEVEL", "error")

	env := support.NewServiceEnv(t)
	configPath := env.WriteFixture(t, "config.json", config.Config{
		ServiceName: "file-service",
		HTTPAddr:    ":7000",
		LogLevel:    "warn",
	})

	cfg, err := config.Load(config.LoadOptions{
		FilePath: configPath,
		CLI: config.ConfigOverlay{
			HTTPAddr: support.Ptr(":7002"),
		},
	})
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	if cfg.HTTPAddr != ":7002" {
		t.Fatalf("HTTPAddr 应优先使用 CLI，实际: %s", cfg.HTTPAddr)
	}
	if cfg.LogLevel != "error" {
		t.Fatalf("LogLevel 应优先使用 env，实际: %s", cfg.LogLevel)
	}
	if cfg.ServiceName != "file-service" {
		t.Fatalf("ServiceName 应回退到 file，实际: %s", cfg.ServiceName)
	}
}

func TestLoadDistinguishesUnsetAndExplicitZeroValue(t *testing.T) {
	env := support.NewServiceEnv(t)
	configPath := env.WriteFixture(t, "config.json", config.Config{
		ServiceName: "file-service",
		HTTPAddr:    ":7000",
		LogLevel:    "warn",
	})

	cfg, err := config.Load(config.LoadOptions{
		FilePath: configPath,
		CLI: config.ConfigOverlay{
			LogLevel: support.Ptr(""),
		},
	})
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	if cfg.LogLevel != "" {
		t.Fatalf("CLI 显式空字符串应覆盖下层值，实际: %q", cfg.LogLevel)
	}
}
