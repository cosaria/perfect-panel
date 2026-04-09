# AGENTS.md

This file provides guidance to coding agents working with this repository.

## Skill routing

When the user's request matches an available skill, ALWAYS invoke it using the Skill
tool as your FIRST action. Do NOT answer directly, do NOT use other tools first.
The skill has specialized workflows that produce better results than ad-hoc answers.

Key routing rules:
- Product ideas, "is this worth building", brainstorming -> invoke office-hours
- Bugs, errors, "why is this broken", 500 errors -> invoke investigate
- Ship, deploy, push, create PR -> invoke ship
- QA, test the site, find bugs -> invoke qa
- Code review, check my diff -> invoke review
- Update docs after shipping -> invoke document-release
- Weekly retro -> invoke retro
- Design system, brand -> invoke design-consultation
- Visual audit, design polish -> invoke design-review
- Architecture review -> invoke plan-eng-review

## Commands

### Root (monorepo)

```bash
make bootstrap          # 安装 repo-local Go 工具 + server/web 依赖 +（若可用）lefthook install
make tools              # 安装 repo-local golangci-lint / goimports 到 .tools/bin
make lint               # 运行 server golangci-lint/go vet + web biome lint
make test               # 运行 server go test ./...
make build              # 本地构建 server 二进制 + web 应用，不是官方发布入口
make format             # server go fmt/goimports + web biome format
make typecheck          # web TypeScript 类型检查
make clean              # 清理构建产物（bin/、.next/、.turbo/）
make dev APP=admin      # 同时启动 server + admin 前端
make dev APP=user       # 同时启动 server + user 前端
make embed              # 构建前端静态产物并复制到 server/web/*-dist
make build-all          # 构建带嵌入静态资源的 server 二进制
bun run repo:contracts  # 根命令合同与 OpenAPI lint；canonical image smoke 由 `.github/workflows/repo-contracts.yml` 额外执行
```

### Server (Go)

```bash
cd server
go test ./...                         # 全部测试
go test -race ./...                   # 带竞态检测
go test ./internal/domains/admin/...  # 单个包测试
go test -run TestFoo ./internal/platform/support/tool/... # 单个测试函数
../.tools/bin/golangci-lint run       # lint
go vet ./...                          # 静态分析
go run . run --config etc/ppanel.yaml # 本地启动
```

### Web (Bun + Turbo)

```bash
cd web
bun install                           # 安装 web 工作区依赖
bun run lint                          # biome lint（通过 turbo）
bun run format                        # biome format --write
bun run dev:admin                     # 单独启动 admin
bun run dev:user                      # 单独启动 user
bun run build                         # 全量构建
bun run typecheck                     # TypeScript 类型检查
bun run openapi                       # 导出本地 huma spec 并重新生成 API client
```

## Architecture

Monorepo 包含两个独立子项目，不共享运行时，通过根 Makefile 统一开发入口。

默认官方发布链是根 `Dockerfile`、根 `docker-compose.yml`、`make embed`、`make build-all`。`server/`、`admin-web`、`user-web` 镜像链只保留为 compatibility path，用于兼容现有镜像流水线和单独审计，不是新贡献者默认入口。

### Server — Go / Gin / huma v2

当前 server 已完成 Phase 0-4 的目录重构主线：

1. **路由层** (`server/internal/platform/http/`) — huma/Gin 路由注册与 HTTP 入口
   - `internal/platform/http/routes.go` — 当前总路由注册入口
   - `routes_{admin,auth,common,public,server,user}.go` — 按域拆分的路由注册器
   - `subscribe.go` / `telegram.go` / `notify.go` — 非 huma 入口桥接
   - `internal/platform/http/middleware/` — HTTP 中间件
2. **业务层** (`server/internal/domains/{admin,auth,user,node,common,...}/`) — 业务逻辑 + huma handler shim
3. **数据层** (`server/internal/platform/persistence/{user,subscribe,order,node,...}/`) — GORM v2 模型，内建 Redis 缓存
4. **类型层** (`server/internal/platform/http/types/`) — 请求/响应结构体定义
5. **依赖注入** (`server/internal/bootstrap/app/ServiceContext`) — model、config、redis、asynq client 的持有者
6. **异步任务** (`server/internal/jobs/`) — Asynq consumer、scheduler、task handler、task payload

关键约定：
- 这轮重构后，server 目录已统一收敛到 `internal/bootstrap`、`internal/domains`、`internal/jobs`、`internal/platform`
- 所有 HTTP 响应均为 200，错误通过 JSON body 中 `{code, msg}` 表达（`internal/platform/http/response/`）
- 错误码体系在 `internal/platform/support/xerr/` — 10xxx DB、20xxx 用户、30xxx 节点、40xxx 鉴权
- 异步任务使用 **Asynq**（Redis DB 5），包括邮件发送、订单关闭、流量统计等
- 定时任务已并入 `internal/jobs/`，时区固定为 `Asia/Shanghai`

路由中间件分组：
| 前缀 | 鉴权 |
|---|---|
| `/api/v1/admin/*` | JWT AuthMiddleware |
| `/api/v1/public/*` | AuthMiddleware + DeviceMiddleware |
| `/api/v1/auth/*` | DeviceMiddleware |
| `/api/v1/server/*` | ServerMiddleware（节点通信） |
| `/api/v1/common/*` | DeviceMiddleware |

`server/internal/platform/support/adapter/` 是订阅配置生成器，将节点转换为各代理协议格式（Vmess/Vless/Trojan/SS/Hysteria2 等）。

### Web — Next.js 15 / Turbo / Bun

两个 Next.js 15 App Router 应用 + 共享包：

- `web/apps/admin/` — 管理后台（port 3000）
- `web/apps/user/` — 用户面板（port 3001）
- `web/packages/ui/` — 共享组件库（shadcn/ui + 自定义 pro-table/editor 等），导入路径 `@workspace/ui/*`

关键技术栈：
- **状态管理**: `@tanstack/react-query` v5（服务端状态）+ `zustand` v5（客户端状态）
- **i18n**: `next-intl`，23 种语言，翻译由 `lobe-i18n` + GPT 自动生成
- **Lint/Format**: Biome（已替代 ESLint+Prettier），配置在 `web/biome.json`
- **API Client**: 由 `@hey-api/openapi-ts` 从 `docs/openapi/*.json` 生成到 `web/apps/*/services/`，**生成文件不要手动编辑**

`request.ts` 拦截器约定：自动附加 JWT header，错误码 40002–40005 触发登出。

### 前后端解耦

前端 OpenAPI 规范现在直接由 server 的 huma 路由导出。修改 server API 后在根目录执行 `bun run openapi`，会依次：
1. 运行 `server/cmd/openapi/main.go` 导出 `docs/openapi/{admin,common,user}.json`
2. 用 Redocly 校验 spec
3. 用 `@hey-api/openapi-ts` 重新生成 `web/apps/*/services/*-api`
4. 最后执行一次格式化

仓库通过 `.github/workflows/monorepo-boundary-guardrail.yml` 的边界护栏（guardrail）阻止 `apps`、`packages`、`scripts`、`bun.lock`、`turbo.json`、`biome.json`、`tsconfig.json` 等前端 workspace 关键根路径重新回流到仓库根目录。

## Git Hooks

根目录 `lefthook.yml` 统管 server 和 web 的 git hooks：

- **pre-commit**（并行）：server 侧执行 go fmt/goimports/golangci-lint/go vet/go test；web 侧执行 biome check --write（仅暂存文件，自动 re-stage）
- server 侧的 `goimports` 与 `golangci-lint` 统一走 repo-local `.tools/bin/*`，缺失时先执行根目录 `make tools`
- **commit-msg**：commitlint 校验提交格式，统一使用 **Conventional Commits** 规范（`feat:` / `fix:` / `chore:` 等）

`make bootstrap` 会先安装 repo-local Go 工具，再自动执行 `lefthook install`（若本机存在 `lefthook`）。

## Dependencies

- Go 1.25+
- Node.js 20+、Bun 1.3.0+
- repo-local `golangci-lint`、`goimports` 通过根目录 `make tools` 安装到 `.tools/bin/`
- lefthook（`brew install lefthook` 或 `go install github.com/evilmartians/lefthook@latest`）
- MySQL、Redis（server 运行时依赖）
- 配置文件：复制 `server/etc/ppanel.yaml.example` 为 `ppanel.yaml`，填入 DB/Redis 连接信息。首次运行若配置为空则从 `PPANEL_DB` / `PPANEL_REDIS` 环境变量自动生成
