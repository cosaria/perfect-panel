# Server 目录一次性重构 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 `server/` 一次性重排为 `cmd + internal` 骨架，同时保住当前仓库的编译、测试、OpenAPI 导出和根目录构建链。

**Architecture:** 保留 `server/` 作为独立 Go module，新增 `server/main.go` 作为根兼容入口，把实现代码收口到 `server/internal/bootstrap`、`server/internal/platform`、`server/internal/domains`、`server/internal/jobs`。目录重排优先于包名重命名；除 `main` 包和新增 `openapi` 导出器外，优先保留现有 package 名，避免把“路径迁移”升级成“语义重写”。

**Tech Stack:** Go 1.25、Cobra、Gin、Huma、GORM、Asynq、Make

---

## 目录与责任锁定

- `server/main.go`
  责任：保住 `cd server && go run . run --config ...` 和 `cd server && go build .` 这条根入口链路，内部只委托给 `server/cmd/ppanel`。
- `server/cmd/ppanel/*.go`
  责任：承载根 CLI 命令树，包括 `run`、`version`、`openapi` 子命令和运行时依赖拼装。
- `server/cmd/openapi/main.go`
  责任：提供可单独构建的 OpenAPI 导出二进制，但不替代根 CLI 子命令。
- `server/internal/bootstrap/app/*.go`
  责任：承接 `server/svc/*.go`，保留 `ServiceContext` 与运行时组合根职责。
- `server/internal/bootstrap/configinit/*.go`
  责任：承接 `server/initialize/*.go`，保留初始化配置、模板、安装引导。
- `server/internal/bootstrap/runtime/*.go`
  责任：承接 `server/runtime/*.go`，保留 live state 和运行时依赖聚合。
- `server/internal/platform/http/*.go`
  责任：承接 `server/routers/*.go`、middleware、response、OpenAPI 导出器和支付通知 HTTP 入口。
- `server/internal/platform/http/types/*.go`
  责任：承接 `server/types/*.go`，仅做路径迁移，不手改生成体内容。
- `server/internal/jobs/*.go`
  责任：承接 `server/worker/*.go` 以及其子目录。
- `server/internal/domains/{admin,auth,common,node,subscribe,telegram,user}/`
  责任：承接 `server/services/*` 主业务目录。
- `server/internal/domains/common/report/`
  责任：承接 `server/services/report/`。
- `server/internal/platform/persistence/*`
  责任：承接 `server/models/*`。
- `server/internal/platform/{cache,crypto,notify,payment}/`
  责任：承接 `server/modules/cache/*`、`server/modules/crypto/*`、`server/modules/notify/*`、`server/modules/payment/*`。
- `server/internal/platform/support/{adapter,auth,traffic,verify}/`
  责任：承接 `server/adapter/*`、`server/modules/auth/*`、`server/modules/traffic/*`、`server/modules/verify/*`。
- `server/internal/platform/support/*`
  责任：继续承接 `server/modules/infra/*` 和 `server/modules/util/*` 下的散项能力。
- `server/cache/`、`server/bin/`、`server/script/`
  责任：保持顶层不动，它们是运行资产和脚本，不是这次 Go 包重排的目标。

## 迁移约束

- 不手改 `server/types/types.go` 的生成内容；只允许移动文件路径、保留 `package types`、修 import。
- 先迁目录，再统一跑 `goimports -w .`，最后再做小范围人工修 import 别名。
- 优先保住已有阶段性护栏测试：`server/cmd/phase1_paths_test.go`、`server/cmd/phase34_structure_test.go`、`server/routers/phase3_route_split_test.go`、`server/worker/phase4_worker_contract_test.go`、`server/services/notify/phase5_protocol_surface_test.go`。
- `server/modules/auth/*` 明确迁到 `server/internal/platform/support/auth/*`。
- `server/modules/traffic/convert.go` 明确迁到 `server/internal/platform/support/traffic/convert.go`。
- `server/services/notify/*` 明确迁到 `server/internal/platform/http/notify/*`。

### Task 1: 建立骨架并拆出根 CLI 入口

**Files:**
- Create: `server/main.go`
- Create: `server/cmd/ppanel/root.go`
- Create: `server/cmd/ppanel/run.go`
- Create: `server/cmd/ppanel/version.go`
- Create: `server/cmd/ppanel/initialize_deps.go`
- Create: `server/cmd/ppanel/runtime_deps.go`
- Create: `server/cmd/ppanel/server_service.go`
- Create: `server/cmd/ppanel/worker_deps.go`
- Create: `server/cmd/phase7_entry_layout_test.go`
- Modify: `server/cmd/phase1_paths_test.go`
- Delete: `server/ppanel.go`
- Delete: `server/cmd/root.go`
- Delete: `server/cmd/run.go`
- Delete: `server/cmd/version.go`
- Delete: `server/cmd/initialize_deps.go`
- Delete: `server/cmd/runtime_deps.go`
- Delete: `server/cmd/server_service.go`
- Delete: `server/cmd/worker_deps.go`
- Test: `server/cmd/phase7_entry_layout_test.go`

- [ ] **Step 1: 写入口布局护栏测试**

```go
package cmd_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPhase7RootMainDelegatesToCmdPpanel(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "main.go"))
	if err != nil {
		t.Fatalf("read main.go: %v", err)
	}
	source := string(data)
	if !strings.Contains(source, "\"github.com/perfect-panel/server/cmd/ppanel\"") {
		t.Fatalf("expected root main to import cmd/ppanel, got:\n%s", source)
	}
	if !strings.Contains(source, "ppanel.Execute()") {
		t.Fatalf("expected root main to delegate to ppanel.Execute, got:\n%s", source)
	}
}

func TestPhase7PpanelCommandPackageExists(t *testing.T) {
	if _, err := os.Stat(filepath.Join("ppanel", "root.go")); err != nil {
		t.Fatalf("expected cmd/ppanel/root.go to exist: %v", err)
	}
}
```

- [ ] **Step 2: 运行护栏测试并确认先失败**

Run: `cd server && go test ./cmd -run 'TestPhase7RootMainDelegatesToCmdPpanel|TestPhase7PpanelCommandPackageExists' -count=1`
Expected: FAIL，报 `open ../main.go: no such file or directory` 或 `cmd/ppanel/root.go` 不存在。

- [ ] **Step 3: 搬迁根 CLI 并补上根兼容入口**

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel
mkdir -p server/cmd/ppanel
git mv server/ppanel.go server/main.go
git mv server/cmd/root.go server/cmd/ppanel/root.go
git mv server/cmd/run.go server/cmd/ppanel/run.go
git mv server/cmd/version.go server/cmd/ppanel/version.go
git mv server/cmd/initialize_deps.go server/cmd/ppanel/initialize_deps.go
git mv server/cmd/runtime_deps.go server/cmd/ppanel/runtime_deps.go
git mv server/cmd/server_service.go server/cmd/ppanel/server_service.go
git mv server/cmd/worker_deps.go server/cmd/ppanel/worker_deps.go
```

```go
// server/main.go
package main

import "github.com/perfect-panel/server/cmd/ppanel"

func main() {
	ppanel.Execute()
}
```

- [ ] **Step 4: 修 import 并回跑入口测试**

Run: `cd server && goimports -w main.go cmd/ppanel && go test ./cmd -run 'TestPhase7RootMainDelegatesToCmdPpanel|TestPhase7PpanelCommandPackageExists' -count=1`
Expected: PASS

- [ ] **Step 5: 提交入口拆分检查点**

```bash
git add server/main.go server/cmd/ppanel server/cmd/phase7_entry_layout_test.go server/cmd/phase1_paths_test.go
git commit -m "refactor(server): split root cli package"
```

### Task 2: 迁移 bootstrap 组合根

**Files:**
- Create: `server/internal/bootstrap/app/serviceContext.go`
- Create: `server/internal/bootstrap/app/asynq.go`
- Create: `server/internal/bootstrap/app/device.go`
- Create: `server/internal/bootstrap/app/logger.go`
- Create: `server/internal/bootstrap/app/mmdb.go`
- Create: `server/internal/bootstrap/app/validate.go`
- Create: `server/internal/bootstrap/configinit/*.go`
- Create: `server/internal/bootstrap/configinit/templates/index.html`
- Create: `server/internal/bootstrap/runtime/deps.go`
- Create: `server/internal/bootstrap/runtime/live_state.go`
- Create: `server/cmd/phase7_bootstrap_layout_test.go`
- Modify: `server/cmd/ppanel/run.go`
- Modify: `server/cmd/ppanel/runtime_deps.go`
- Modify: `server/cmd/ppanel/server_service.go`
- Modify: `server/cmd/ppanel/worker_deps.go`
- Delete: `server/svc/*.go`
- Delete: `server/initialize/*.go`
- Delete: `server/initialize/templates/index.html`
- Delete: `server/runtime/*.go`
- Test: `server/cmd/phase7_bootstrap_layout_test.go`

- [ ] **Step 1: 写 bootstrap 布局护栏测试**

```go
package cmd_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPhase7BootstrapDirectoriesExist(t *testing.T) {
	targets := []string{
		filepath.Join("..", "internal", "bootstrap", "app", "serviceContext.go"),
		filepath.Join("..", "internal", "bootstrap", "configinit", "config.go"),
		filepath.Join("..", "internal", "bootstrap", "runtime", "live_state.go"),
	}
	for _, target := range targets {
		if _, err := os.Stat(target); err != nil {
			t.Fatalf("expected bootstrap target %s to exist: %v", target, err)
		}
	}
}
```

- [ ] **Step 2: 运行护栏测试并确认先失败**

Run: `cd server && go test ./cmd -run TestPhase7BootstrapDirectoriesExist -count=1`
Expected: FAIL，报 `internal/bootstrap/...` 路径不存在。

- [ ] **Step 3: 搬迁 `svc`、`initialize`、`runtime` 到新目录**

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel
mkdir -p server/internal/bootstrap/app server/internal/bootstrap/configinit/templates server/internal/bootstrap/runtime
git mv server/svc/asynq.go server/internal/bootstrap/app/asynq.go
git mv server/svc/devce.go server/internal/bootstrap/app/device.go
git mv server/svc/logger.go server/internal/bootstrap/app/logger.go
git mv server/svc/mmdb.go server/internal/bootstrap/app/mmdb.go
git mv server/svc/serviceContext.go server/internal/bootstrap/app/serviceContext.go
git mv server/svc/validate.go server/internal/bootstrap/app/validate.go
git mv server/initialize/*.go server/internal/bootstrap/configinit/
git mv server/initialize/templates/index.html server/internal/bootstrap/configinit/templates/index.html
git mv server/runtime/deps.go server/internal/bootstrap/runtime/deps.go
git mv server/runtime/live_state.go server/internal/bootstrap/runtime/live_state.go
```

```go
// server/cmd/ppanel/run.go 的关键 import 方向
import (
	appbootstrap "github.com/perfect-panel/server/internal/bootstrap/app"
	configinit "github.com/perfect-panel/server/internal/bootstrap/configinit"
	appruntime "github.com/perfect-panel/server/internal/bootstrap/runtime"
)
```

- [ ] **Step 4: 修引用并回跑 bootstrap 测试**

Run: `cd server && goimports -w cmd/ppanel internal/bootstrap && go test ./cmd -run TestPhase7BootstrapDirectoriesExist -count=1`
Expected: PASS

- [ ] **Step 5: 提交 bootstrap 检查点**

```bash
git add server/internal/bootstrap server/cmd/ppanel server/cmd/phase7_bootstrap_layout_test.go
git commit -m "refactor(server): move bootstrap packages under internal"
```

### Task 3: 迁移 HTTP 层、支付通知和 OpenAPI 导出器

**Files:**
- Create: `server/internal/platform/http/*.go`
- Create: `server/internal/platform/http/middleware/*.go`
- Create: `server/internal/platform/http/response/*.go`
- Create: `server/internal/platform/http/notify/*.go`
- Create: `server/internal/platform/http/openapi/export.go`
- Create: `server/internal/platform/http/types/types.go`
- Create: `server/internal/platform/http/types/subscribe.go`
- Create: `server/cmd/openapi/main.go`
- Modify: `server/cmd/ppanel/openapi.go`
- Create: `server/cmd/phase7_http_layout_test.go`
- Delete: `server/routers/*.go`
- Delete: `server/routers/middleware/*.go`
- Delete: `server/routers/response/*.go`
- Delete: `server/services/notify/*.go`
- Delete: `server/types/types.go`
- Delete: `server/types/subscribe.go`
- Test: `server/cmd/phase7_http_layout_test.go`
- Test: `server/internal/platform/http/phase3_route_split_test.go`
- Test: `server/internal/platform/http/notify/phase5_protocol_surface_test.go`

- [ ] **Step 1: 写 HTTP 布局护栏测试**

```go
package cmd_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPhase7HTTPDirectoriesExist(t *testing.T) {
	targets := []string{
		filepath.Join("..", "internal", "platform", "http", "routes.go"),
		filepath.Join("..", "internal", "platform", "http", "middleware", "authMiddleware.go"),
		filepath.Join("..", "internal", "platform", "http", "response", "response.go"),
		filepath.Join("..", "internal", "platform", "http", "types", "types.go"),
		filepath.Join("..", "cmd", "openapi", "main.go"),
	}
	for _, target := range targets {
		if _, err := os.Stat(target); err != nil {
			t.Fatalf("expected HTTP target %s to exist: %v", target, err)
		}
	}
}
```

- [ ] **Step 2: 运行护栏测试并确认先失败**

Run: `cd server && go test ./cmd -run TestPhase7HTTPDirectoriesExist -count=1`
Expected: FAIL，至少报 `internal/platform/http/routes.go` 与 `cmd/openapi/main.go` 不存在。

- [ ] **Step 3: 搬迁 `routers`、`services/notify`、`types` 并提取 OpenAPI 导出器**

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel
mkdir -p server/internal/platform/http/middleware server/internal/platform/http/response server/internal/platform/http/notify server/internal/platform/http/openapi server/internal/platform/http/types server/cmd/openapi
git mv server/routers/huma_problem.go server/internal/platform/http/huma_problem.go
git mv server/routers/initialize_deps.go server/internal/platform/http/initialize_deps.go
git mv server/routers/notify.go server/internal/platform/http/notify.go
git mv server/routers/openapi_conventions.go server/internal/platform/http/openapi_conventions.go
git mv server/routers/openapi_conventions_test.go server/internal/platform/http/openapi_conventions_test.go
git mv server/routers/phase3_route_split_test.go server/internal/platform/http/phase3_route_split_test.go
git mv server/routers/phase5_huma_error_contract_test.go server/internal/platform/http/phase5_huma_error_contract_test.go
git mv server/routers/phase6_live_runtime_deps_test.go server/internal/platform/http/phase6_live_runtime_deps_test.go
git mv server/routers/routes.go server/internal/platform/http/routes.go
git mv server/routers/routes_admin.go server/internal/platform/http/routes_admin.go
git mv server/routers/routes_auth.go server/internal/platform/http/routes_auth.go
git mv server/routers/routes_common.go server/internal/platform/http/routes_common.go
git mv server/routers/routes_public.go server/internal/platform/http/routes_public.go
git mv server/routers/routes_server.go server/internal/platform/http/routes_server.go
git mv server/routers/routes_user.go server/internal/platform/http/routes_user.go
git mv server/routers/runtime_dynamic_deps.go server/internal/platform/http/runtime_dynamic_deps.go
git mv server/routers/subscribe.go server/internal/platform/http/subscribe.go
git mv server/routers/telegram.go server/internal/platform/http/telegram.go
git mv server/routers/middleware server/internal/platform/http/middleware
git mv server/routers/response server/internal/platform/http/response
git mv server/services/notify/*.go server/internal/platform/http/notify/
git mv server/types/types.go server/internal/platform/http/types/types.go
git mv server/types/subscribe.go server/internal/platform/http/types/subscribe.go
git mv server/cmd/openapi.go server/cmd/ppanel/openapi.go
```

```go
// server/internal/platform/http/openapi/export.go
package openapi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	handler "github.com/perfect-panel/server/internal/platform/http"
)

func Export(outputDir string) error {
	if outputDir == "" {
		outputDir = "docs/openapi"
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	apis := handler.RegisterHandlersForSpec(router)

	userSpec, err := apis.UserOpenAPI()
	if err != nil {
		return fmt.Errorf("merge user spec: %w", err)
	}

	specs := map[string]any{
		"admin":  apis.Admin.OpenAPI(),
		"common": apis.Common.OpenAPI(),
		"user":   userSpec,
	}

	for name, spec := range specs {
		data, err := json.MarshalIndent(spec, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal %s spec: %w", name, err)
		}
		path := filepath.Join(outputDir, name+".json")
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
	}
	return nil
}
```

```go
// server/cmd/openapi/main.go
package main

import (
	"flag"
	"log"

	httpopenapi "github.com/perfect-panel/server/internal/platform/http/openapi"
)

func main() {
	outputDir := flag.String("o", "docs/openapi", "Output directory for spec files")
	flag.Parse()
	if err := httpopenapi.Export(*outputDir); err != nil {
		log.Fatal(err)
	}
}
```

- [ ] **Step 4: 改根 CLI `openapi` 子命令并回跑 HTTP 测试**

Run: `cd server && goimports -w cmd/ppanel cmd/openapi internal/platform/http && go test ./cmd -run TestPhase7HTTPDirectoriesExist -count=1 && go test ./internal/platform/http/... -count=1`
Expected: PASS

- [ ] **Step 5: 提交 HTTP 检查点**

```bash
git add server/internal/platform/http server/cmd/ppanel/openapi.go server/cmd/openapi server/cmd/phase7_http_layout_test.go
git commit -m "refactor(server): move http and openapi packages under internal"
```

### Task 4: 迁移异步任务与调度器

**Files:**
- Create: `server/internal/jobs/*.go`
- Create: `server/internal/jobs/email/*.go`
- Create: `server/internal/jobs/order/*.go`
- Create: `server/internal/jobs/registry/routes.go`
- Create: `server/internal/jobs/sms/*.go`
- Create: `server/internal/jobs/spec/*.go`
- Create: `server/internal/jobs/subscription/*.go`
- Create: `server/internal/jobs/task/*.go`
- Create: `server/internal/jobs/traffic/*.go`
- Create: `server/cmd/phase7_jobs_layout_test.go`
- Modify: `server/cmd/ppanel/run.go`
- Modify: `server/cmd/ppanel/worker_deps.go`
- Modify: `server/internal/bootstrap/runtime/deps.go`
- Delete: `server/worker/*.go`
- Delete: `server/worker/email/*.go`
- Delete: `server/worker/order/*.go`
- Delete: `server/worker/registry/routes.go`
- Delete: `server/worker/sms/*.go`
- Delete: `server/worker/spec/*.go`
- Delete: `server/worker/subscription/*.go`
- Delete: `server/worker/task/*.go`
- Delete: `server/worker/traffic/*.go`
- Test: `server/cmd/phase7_jobs_layout_test.go`
- Test: `server/internal/jobs/phase4_worker_contract_test.go`

- [ ] **Step 1: 写 jobs 布局护栏测试**

```go
package cmd_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPhase7JobsDirectoryExists(t *testing.T) {
	targets := []string{
		filepath.Join("..", "internal", "jobs", "consumer_service.go"),
		filepath.Join("..", "internal", "jobs", "scheduler_service.go"),
		filepath.Join("..", "internal", "jobs", "registry", "routes.go"),
	}
	for _, target := range targets {
		if _, err := os.Stat(target); err != nil {
			t.Fatalf("expected jobs target %s to exist: %v", target, err)
		}
	}
}
```

- [ ] **Step 2: 运行护栏测试并确认先失败**

Run: `cd server && go test ./cmd -run TestPhase7JobsDirectoryExists -count=1`
Expected: FAIL，报 `internal/jobs/...` 不存在。

- [ ] **Step 3: 搬迁 `worker` 到 `internal/jobs` 并修依赖**

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel
mkdir -p server/internal/jobs
git mv server/worker/* server/internal/jobs/
```

```go
// server/cmd/ppanel/run.go 的关键替换
import (
	serverjobs "github.com/perfect-panel/server/internal/jobs"
)

services.Add(serverjobs.NewConsumerService(ctx.Config, newWorkerRegistryDeps(ctx, live)))
services.Add(serverjobs.NewSchedulerService(ctx.Config))
```

- [ ] **Step 4: 格式化并回跑 jobs 测试**

Run: `cd server && goimports -w cmd/ppanel internal/jobs internal/bootstrap && go test ./cmd -run TestPhase7JobsDirectoryExists -count=1 && go test ./internal/jobs/... -count=1`
Expected: PASS

- [ ] **Step 5: 提交 jobs 检查点**

```bash
git add server/internal/jobs server/cmd/ppanel server/cmd/phase7_jobs_layout_test.go server/internal/bootstrap/runtime
git commit -m "refactor(server): move worker packages under internal jobs"
```

### Task 5: 迁移业务领域目录

**Files:**
- Create: `server/internal/domains/admin/`
- Create: `server/internal/domains/auth/`
- Create: `server/internal/domains/common/`
- Create: `server/internal/domains/common/report/`
- Create: `server/internal/domains/node/`
- Create: `server/internal/domains/subscribe/`
- Create: `server/internal/domains/telegram/`
- Create: `server/internal/domains/user/`
- Modify: `server/internal/platform/http/routes*.go`
- Modify: `server/internal/jobs/**/*.go`
- Modify: `server/cmd/phase34_structure_test.go`
- Create: `server/cmd/phase7_domain_layout_test.go`
- Delete: `server/services/admin/`
- Delete: `server/services/auth/`
- Delete: `server/services/common/`
- Delete: `server/services/node/`
- Delete: `server/services/report/`
- Delete: `server/services/subscribe/`
- Delete: `server/services/telegram/`
- Delete: `server/services/user/`
- Test: `server/cmd/phase7_domain_layout_test.go`
- Test: `server/cmd/phase34_structure_test.go`

- [ ] **Step 1: 写 domain 布局护栏测试**

```go
package cmd_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPhase7DomainDirectoriesExist(t *testing.T) {
	targets := []string{
		filepath.Join("..", "internal", "domains", "admin"),
		filepath.Join("..", "internal", "domains", "auth"),
		filepath.Join("..", "internal", "domains", "common"),
		filepath.Join("..", "internal", "domains", "node"),
		filepath.Join("..", "internal", "domains", "subscribe"),
		filepath.Join("..", "internal", "domains", "telegram"),
		filepath.Join("..", "internal", "domains", "user"),
	}
	for _, target := range targets {
		if _, err := os.Stat(target); err != nil {
			t.Fatalf("expected domain directory %s to exist: %v", target, err)
		}
	}
}
```

- [ ] **Step 2: 运行护栏测试并确认先失败**

Run: `cd server && go test ./cmd -run TestPhase7DomainDirectoriesExist -count=1`
Expected: FAIL，报 `internal/domains/...` 不存在。

- [ ] **Step 3: 搬迁 `services/*` 到 `internal/domains/*`**

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel
mkdir -p server/internal/domains/common
git mv server/services/admin server/internal/domains/admin
git mv server/services/auth server/internal/domains/auth
git mv server/services/common server/internal/domains/common
git mv server/services/node server/internal/domains/node
git mv server/services/report server/internal/domains/common/report
git mv server/services/subscribe server/internal/domains/subscribe
git mv server/services/telegram server/internal/domains/telegram
git mv server/services/user server/internal/domains/user
```

```bash
cd server
rg -l --glob '*.go' 'github.com/perfect-panel/server/services/(admin|auth|common|node|report|subscribe|telegram|user)' . | \
	xargs perl -0pi -e '
		s#github.com/perfect-panel/server/services/admin#github.com/perfect-panel/server/internal/domains/admin#g;
		s#github.com/perfect-panel/server/services/auth#github.com/perfect-panel/server/internal/domains/auth#g;
		s#github.com/perfect-panel/server/services/common#github.com/perfect-panel/server/internal/domains/common#g;
		s#github.com/perfect-panel/server/services/node#github.com/perfect-panel/server/internal/domains/node#g;
		s#github.com/perfect-panel/server/services/report#github.com/perfect-panel/server/internal/domains/common/report#g;
		s#github.com/perfect-panel/server/services/subscribe#github.com/perfect-panel/server/internal/domains/subscribe#g;
		s#github.com/perfect-panel/server/services/telegram#github.com/perfect-panel/server/internal/domains/telegram#g;
		s#github.com/perfect-panel/server/services/user#github.com/perfect-panel/server/internal/domains/user#g;
	'
```

- [ ] **Step 4: 格式化并回跑领域测试**

Run: `cd server && goimports -w internal/domains internal/platform/http internal/jobs && go test ./cmd -run 'TestPhase7DomainDirectoriesExist|TestPhase3ServicesNoLongerKeepHandlerLogicPairs' -count=1`
Expected: PASS

- [ ] **Step 5: 提交 domains 检查点**

```bash
git add server/internal/domains server/internal/platform/http server/internal/jobs server/cmd/phase34_structure_test.go server/cmd/phase7_domain_layout_test.go
git commit -m "refactor(server): move service packages under internal domains"
```

### Task 6: 迁移平台能力与历史支持包

**Files:**
- Create: `server/internal/platform/persistence/*`
- Create: `server/internal/platform/cache/*`
- Create: `server/internal/platform/crypto/*`
- Create: `server/internal/platform/notify/*`
- Create: `server/internal/platform/payment/*`
- Create: `server/internal/platform/support/*`
- Create: `server/internal/platform/support/auth/*`
- Create: `server/internal/platform/support/traffic/convert.go`
- Create: `server/internal/platform/support/verify/*`
- Create: `server/internal/platform/support/adapter/*`
- Modify: `server/cmd/phase1_paths_test.go`
- Create: `server/cmd/phase7_legacy_imports_test.go`
- Delete: `server/models/*`
- Delete: `server/modules/auth/*`
- Delete: `server/modules/cache/*`
- Delete: `server/modules/crypto/*`
- Delete: `server/modules/infra/*`
- Delete: `server/modules/notify/*`
- Delete: `server/modules/payment/*`
- Delete: `server/modules/traffic/*`
- Delete: `server/modules/util/*`
- Delete: `server/modules/verify/*`
- Delete: `server/adapter/*`
- Test: `server/cmd/phase1_paths_test.go`
- Test: `server/cmd/phase7_legacy_imports_test.go`

- [ ] **Step 1: 写 legacy import 清理护栏测试**

```go
package cmd_test

import (
	"os/exec"
	"strings"
	"testing"
)

func TestPhase7NoLegacyPackageImportsRemain(t *testing.T) {
	cmd := exec.Command("rg", "-n", "--glob", "*.go", "github.com/perfect-panel/server/(models|modules|services|routers|svc|initialize|runtime|worker|types|adapter)", ".")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected no legacy imports, found:\n%s", string(out))
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return
	}
	if strings.TrimSpace(string(out)) != "" {
		t.Fatalf("unexpected rg output:\n%s", string(out))
	}
}
```

- [ ] **Step 2: 运行护栏测试并确认先失败**

Run: `cd server && go test ./cmd -run TestPhase7NoLegacyPackageImportsRemain -count=1`
Expected: FAIL，输出仍能搜到旧 import 路径。

- [ ] **Step 3: 搬迁 `models`、`modules`、`adapter` 并做全量 import 重写**

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel
mkdir -p server/internal/platform/persistence server/internal/platform/cache server/internal/platform/crypto server/internal/platform/notify server/internal/platform/payment server/internal/platform/support/auth server/internal/platform/support/traffic server/internal/platform/support/verify server/internal/platform/support/adapter
git mv server/models/* server/internal/platform/persistence/
git mv server/modules/cache/* server/internal/platform/cache/
git mv server/modules/crypto/* server/internal/platform/crypto/
git mv server/modules/notify/* server/internal/platform/notify/
git mv server/modules/payment/* server/internal/platform/payment/
git mv server/modules/auth/* server/internal/platform/support/auth/
git mv server/modules/traffic/convert.go server/internal/platform/support/traffic/convert.go
git mv server/modules/verify/* server/internal/platform/support/verify/
git mv server/modules/infra/* server/internal/platform/support/
git mv server/modules/util/* server/internal/platform/support/
git mv server/adapter/* server/internal/platform/support/adapter/
```

```bash
cd server
rg -l --glob '*.go' 'github.com/perfect-panel/server/(models|modules|adapter|types)' . | \
	xargs perl -0pi -e '
		s#github.com/perfect-panel/server/models/#github.com/perfect-panel/server/internal/platform/persistence/#g;
		s#github.com/perfect-panel/server/modules/cache#github.com/perfect-panel/server/internal/platform/cache#g;
		s#github.com/perfect-panel/server/modules/crypto#github.com/perfect-panel/server/internal/platform/crypto#g;
		s#github.com/perfect-panel/server/modules/notify#github.com/perfect-panel/server/internal/platform/notify#g;
		s#github.com/perfect-panel/server/modules/payment#github.com/perfect-panel/server/internal/platform/payment#g;
		s#github.com/perfect-panel/server/modules/auth#github.com/perfect-panel/server/internal/platform/support/auth#g;
		s#github.com/perfect-panel/server/modules/traffic#github.com/perfect-panel/server/internal/platform/support/traffic#g;
		s#github.com/perfect-panel/server/modules/verify#github.com/perfect-panel/server/internal/platform/support/verify#g;
		s#github.com/perfect-panel/server/modules/infra#github.com/perfect-panel/server/internal/platform/support#g;
		s#github.com/perfect-panel/server/modules/util#github.com/perfect-panel/server/internal/platform/support#g;
		s#github.com/perfect-panel/server/adapter#github.com/perfect-panel/server/internal/platform/support/adapter#g;
		s#github.com/perfect-panel/server/types#github.com/perfect-panel/server/internal/platform/http/types#g;
	'
```

- [ ] **Step 4: 格式化并回跑路径烟测**

Run: `cd server && goimports -w . && go test ./cmd -run 'TestPhase1TopLevelPathsExist|TestPhase7NoLegacyPackageImportsRemain' -count=1`
Expected: PASS

- [ ] **Step 5: 提交 platform 检查点**

```bash
git add server/internal/platform server/cmd/phase1_paths_test.go server/cmd/phase7_legacy_imports_test.go
git commit -m "refactor(server): move platform packages under internal"
```

### Task 7: 更新文档并完成全量验证

**Files:**
- Modify: `server/README.md`
- Modify: `AGENTS.md`
- Modify: `docs/api-governance.md`
- Modify: `Makefile`（仅在根链路失配时）
- Test: `server/cmd/phase1_paths_test.go`
- Test: `server/cmd/phase34_structure_test.go`
- Test: `server/cmd/phase5_openapi_contract_test.go`
- Test: `server/cmd/phase6_openapi_command_test.go`
- Test: `server/internal/platform/http/...`
- Test: `server/internal/jobs/...`

- [ ] **Step 1: 更新仓库导航文档**

```markdown
在 `server/README.md`、`AGENTS.md`、`docs/api-governance.md` 中统一替换：
- `server/routers/` -> `server/internal/platform/http/`
- `server/services/...` -> `server/internal/domains/...`
- `server/models/...` -> `server/internal/platform/persistence/...`
- `server/svc/` / `server/initialize/` / `server/runtime/` -> `server/internal/bootstrap/...`
- `server/worker/` -> `server/internal/jobs/`
- `server/types/` -> `server/internal/platform/http/types/`
```

- [ ] **Step 2: 运行格式化和单模块测试**

Run: `cd server && go fmt ./... && goimports -w . && go test ./...`
Expected: PASS

- [ ] **Step 3: 运行关键命令回归**

Run: `cd server && go build ./... && go run . openapi -o ../docs/openapi`
Expected: PASS，且 `../docs/openapi/admin.json`、`../docs/openapi/common.json`、`../docs/openapi/user.json` 被更新。

- [ ] **Step 4: 运行仓库级验证**

Run: `cd /Users/admin/Codes/ProxyCode/perfect-panel && make test && make build-all`
Expected: PASS，`make build-all` 能产出 `server/bin/ppanel`。

- [ ] **Step 5: 如本地配置可用，再做启动验证并提交最终结果**

Run: `cd server && go run . run --config etc/ppanel.yaml`
Expected: 服务成功启动，出现首批初始化日志后手动 `Ctrl+C` 退出；随后执行：

```bash
git add server/README.md AGENTS.md docs/api-governance.md Makefile docs/openapi
git commit -m "refactor(server): collapse legacy backend layout into cmd-internal"
```

## 自检

### Spec coverage

- `cmd/` 与 `internal/` 骨架：由 Task 1 到 Task 6 覆盖。
- `bootstrap` 收口：由 Task 2 覆盖。
- `platform/http`、支付通知和 OpenAPI 导出：由 Task 3 覆盖。
- `jobs`：由 Task 4 覆盖。
- `domains`：由 Task 5 覆盖。
- `platform`：由 Task 6 覆盖。
- `go test ./...`、`go run . openapi ...`、`go build ./...`、`make test`、`make build-all`：由 Task 7 覆盖。
- 设计文档里没写死但实施必须落锤的目录归属：
  `modules/auth`、`modules/traffic`、`services/notify`、`types`、`cache/bin/script` 已在“目录与责任锁定”中补齐。

### Placeholder scan

- 已确认文档中没有计划占位词或“后面再补”的表述。
- 所有任务都给出了明确的目标路径、命令和预期输出。

### Type consistency

- 根入口固定为 `server/main.go -> server/cmd/ppanel.Execute()`
- OpenAPI 导出固定由 `server/internal/platform/http/openapi.Export` 复用
- `types` 统一落在 `server/internal/platform/http/types`
- `modules/auth` 统一落在 `server/internal/platform/support/auth`
