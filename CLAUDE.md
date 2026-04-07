# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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
make bootstrap          # 安装 server + web 全部依赖 + lefthook install
make lint               # 运行 server golangci-lint/go vet + web biome lint
make test               # 运行 server go test ./...
make build              # 构建 server 二进制 + web 应用
make format             # server go fmt/goimports + web biome format
make typecheck          # web TypeScript 类型检查
make clean              # 清理构建产物（bin/、.next/、.turbo/）
make dev APP=admin      # 同时启动 server + admin 前端
make dev APP=user       # 同时启动 server + user 前端
```

### Server (Go)

```bash
cd server
go test ./...                         # 全部测试
go test -race ./...                   # 带竞态检测
go test ./services/admin/...          # 单个包测试
go test -run TestFoo ./pkg/tool/...   # 单个测试函数
golangci-lint run                     # lint
go vet ./...                          # 静态分析
go run . run --config etc/ppanel.yaml # 本地启动
```

### Web (Bun + Turbo)

```bash
bun install                           # 安装依赖（从根目录）
bun run lint                          # biome lint（通过 turbo）
bun run format                        # biome format --write
bun run dev:admin                     # 单独启动 admin
bun run dev:user                      # 单独启动 user
bun run build                         # 全量构建
bun run typecheck                     # TypeScript 类型检查
bun run openapi                       # 从远端 Swagger JSON 重新生成 API client
```

## Architecture

Monorepo 包含两个独立子项目，不共享运行时，通过根 Makefile 统一开发入口。

### Server — Go / Gin / huma v2

当前 server 已进入目录重构后的骨架阶段：

1. **路由层** (`server/routers/`) — huma/Gin 路由注册与 HTTP 入口
   - `routers/routes.go` — 当前总路由注册入口
   - `routers/{admin,auth,public,server,...}/` — 现阶段仍保留按域拆分的 handler wrapper
   - `routers/middleware/` — HTTP 中间件
2. **业务层** (`server/services/{admin,auth,user,node,common,...}/`) — 业务逻辑
3. **数据层** (`server/models/{user,subscribe,order,node,...}/`) — GORM v2 模型，内建 Redis 缓存
4. **类型层** (`server/types/`) — 请求/响应结构体，当前仍保留集中定义
5. **依赖注入** (`server/svc/ServiceContext`) — model、config、redis、asynq client 的持有者

关键约定：
- 这轮重构后，旧 `internal/{handler,logic,model,svc,config,types}` 路径已提升到 `routers/services/models/svc/config/types`
- 所有 HTTP 响应均为 200，错误通过 JSON body 中 `{code, msg}` 表达（`pkg/result/`）
- 错误码体系在 `pkg/xerr/` — 10xxx DB、20xxx 用户、30xxx 节点、40xxx 鉴权
- 异步任务使用 **Asynq**（Redis DB 5），包括邮件发送、订单关闭、流量统计等
- 定时任务在 `scheduler/`，时区固定为 `Asia/Shanghai`

路由中间件分组：
| 前缀 | 鉴权 |
|---|---|
| `/v1/admin/*` | JWT AuthMiddleware |
| `/v1/public/*` | AuthMiddleware + DeviceMiddleware |
| `/v1/auth/*` | DeviceMiddleware |
| `/v1/server/*` | ServerMiddleware（节点通信） |
| `/v1/common/*` | DeviceMiddleware |

`server/adapter/` 是订阅配置生成器，将节点转换为各代理协议格式（Vmess/Vless/Trojan/SS/Hysteria2 等）。

### Web — Next.js 15 / Turbo / Bun

两个 Next.js 15 App Router 应用 + 共享包：

- `apps/admin/` — 管理后台（port 3000）
- `apps/user/` — 用户面板（port 3001）
- `packages/ui/` — 共享组件库（shadcn/ui + 自定义 pro-table/editor 等），导入路径 `@workspace/ui/*`

关键技术栈：
- **状态管理**: `@tanstack/react-query` v5（服务端状态）+ `zustand` v5（客户端状态）
- **i18n**: `next-intl`，23 种语言，翻译由 `lobe-i18n` + GPT 自动生成
- **Lint/Format**: Biome（已替代 ESLint+Prettier），配置在 `biome.json`
- **API Client**: 由 `@umijs/openapi` 从远端 Swagger JSON 生成到 `apps/*/services/`，**生成文件不要手动编辑**

`request.ts` 拦截器约定：自动附加 JWT header，错误码 40002–40005 触发登出。

### 前后端解耦

前端的 OpenAPI 规范来源是 `ppanel-docs` 仓库的 Swagger JSON，不是直接从 server 源码生成。修改 server API 后需同步更新 docs 仓库，再在根目录执行 `bun run openapi` 重新生成 client。

## Git Hooks

根目录 `lefthook.yml` 统管 server 和 web 的 git hooks：

- **pre-commit**（并行）：server 侧执行 go fmt/goimports/golangci-lint/go vet/go test；web 侧执行 biome check --write（仅暂存文件，自动 re-stage）
- **commit-msg**：commitlint 校验提交格式，统一使用 **Conventional Commits** 规范（`feat:` / `fix:` / `chore:` 等）

`make bootstrap` 会自动执行 `lefthook install`。

## Dependencies

- Go 1.23+、golangci-lint、goimports
- Node.js 20+、Bun 1.3.0+
- lefthook（`brew install lefthook` 或 `go install github.com/evilmartians/lefthook@latest`）
- MySQL、Redis（server 运行时依赖）
- 配置文件：复制 `server/etc/ppanel.yaml.example` 为 `ppanel.yaml`，填入 DB/Redis 连接信息。首次运行若配置为空则从 `PPANEL_DB` / `PPANEL_REDIS` 环境变量自动生成
