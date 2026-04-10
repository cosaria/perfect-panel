package bootstrap_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/perfect-panel/server-v2/internal/platform/config"
)

func TestLoadPrefersCLIOverEnvAndFile(t *testing.T) {
	t.Setenv("PPANEL_HTTP_ADDR", ":7001")
	t.Setenv("PPANEL_LOG_LEVEL", "error")

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")
	fileConfig := config.Config{
		ServiceName: "file-service",
		HTTPAddr:    ":7000",
		LogLevel:    "warn",
	}
	fileContent, err := json.Marshal(fileConfig)
	if err != nil {
		t.Fatalf("序列化配置文件失败: %v", err)
	}
	if err := os.WriteFile(configPath, fileContent, 0o600); err != nil {
		t.Fatalf("写入配置文件失败: %v", err)
	}

	cfg, err := config.Load(config.LoadOptions{
		FilePath: configPath,
		CLI: config.Config{
			HTTPAddr: ":7002",
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
