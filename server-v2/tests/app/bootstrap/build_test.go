package bootstrap_test

import (
	"testing"

	"github.com/perfect-panel/server-v2/internal/app/bootstrap"
	"github.com/perfect-panel/server-v2/internal/platform/config"
	"github.com/perfect-panel/server-v2/tests/support"
)

func TestBuildInjectsDefaultLoggerAndHealthHandler(t *testing.T) {
	env := support.NewServiceEnv(t)
	if !env.DB.IsReady() {
		t.Fatal("测试服务环境中的 DB stub 应该是可用状态")
	}

	configPath := env.WriteFixture(t, "bootstrap.json", config.Config{
		ServiceName: "bootstrap-service",
		HTTPAddr:    ":7788",
		LogLevel:    "warn",
	})

	container, err := bootstrap.Build(bootstrap.Options{
		Config: config.LoadOptions{
			FilePath: configPath,
		},
		Logger: nil,
	})
	if err != nil {
		t.Fatalf("bootstrap.Build 返回错误: %v", err)
	}

	if container.Config.HTTPAddr != ":7788" {
		t.Fatalf("容器配置未正确加载，HTTPAddr=%s", container.Config.HTTPAddr)
	}
	if container.Logger == nil {
		t.Fatal("未注入默认 logger")
	}
	if container.HealthHandler == nil {
		t.Fatal("未装配 health handler")
	}
}
