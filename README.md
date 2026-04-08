# perfect-panel

`perfect-panel` 现在按两个并列工程组织：

- `server/` 是独立 Go 后端工程
- `web/` 是 Bun + Turbo 驱动的前端工作区
- 根目录只保留统一命令入口、OpenAPI 契约、Docker 和 CI 编排

默认打包与部署入口仍是嵌入式单镜像路径：前端静态产物先构建后复制到 `server/web/*-dist`，再由 Go 二进制一起分发。

## Support Matrix

| 路径 | 角色 | 新贡献者默认入口 | 说明 |
|---|---|---|---|
| 根 `Dockerfile` / `docker-compose.yml` / `make embed` / `make build-all` | canonical | 是 | 仓库默认官方发布链。优先维护这里的构建、打包和发布说明。 |
| `server/Dockerfile`、`docker-build.yml` 中的 `server` / `admin-web` / `user-web` 镜像链 | compatibility | 否 | 保留给兼容现有镜像流水线和单独审计场景，不作为新贡献者的默认入口。 |

## Requirements

- Go 1.25+
- Node.js 20+
- Bun 1.3.0+
- lefthook（可选，用于本地 git hooks）

## Quick Start

从仓库根目录执行：

```bash
make bootstrap
make lint
make test
bun run openapi
```

`make bootstrap` 会安装 repo-local Go 工具（`golangci-lint`、`goimports`）到 `.tools/bin/`，再下载 `server` 依赖、安装 `web` 依赖，并在本机存在 `lefthook` 时自动安装 hooks。若只想补齐 Go 工具链，可单独执行 `make tools`。

本地联调时，显式选择一个前端应用：

```bash
make dev APP=admin
make dev APP=user
```

`make dev` 会同时启动：

- `server`: `go run . run --config etc/ppanel.yaml`
- `web`: 你选择的一个前端应用

## Root Command Contract

- `make bootstrap`: 安装 repo-local Go 工具、下载 server 依赖并安装 web 依赖
- `make tools`: 安装 repo-local Go 工具到 `.tools/bin/`
- `make lint`: 运行 server 的 `golangci-lint` 和 `go vet`，再运行 web 的 `bun run lint`
- `make test`: 运行 `cd server && go test ./...`
- `make format`: 运行 server 的 `go fmt` / `goimports`，再运行 web 的 `bun run format`
- `make typecheck`: 运行 web 的 `bun run typecheck`
- `make dev APP=admin|user`: 启动 server 和一个前端应用
- `make embed`: 构建前端静态产物并复制到 `server/web/*-dist`
- `make build-all`: 构建带嵌入静态资源的 server 二进制
- `bun run repo:contracts`: 运行根命令合同和 OpenAPI lint；canonical image smoke 由 `.github/workflows/repo-contracts.yml` 额外执行 `docker build -f Dockerfile .`
- `bun run openapi`: 导出 spec、校验 spec、更新前端 generated clients

## Repository Layout

```text
.
├── Makefile
├── package.json
├── docs/openapi/
├── server/
└── web/
```

- `server/`: 后端服务、API、构建与运行入口
- `web/`: 前端工作区，包含 admin 和 user 两个应用
- `package.json`: repo 级 orchestration shell，不再是前端 workspace 根
- `docs/openapi/`: server 与 web 共享的 API 契约产物

## CI

当前根级 workflow：

- `.github/workflows/lint.yml`
- `.github/workflows/openapi-governance.yml`
- `.github/workflows/docker-build.yml`
- `.github/workflows/monorepo-boundary-guardrail.yml`
- `.github/workflows/repo-contracts.yml`

默认 canonical 打包与部署入口是根 `Dockerfile`、根 `docker-compose.yml`、`make embed`、`make build-all`。`.github/workflows/docker-build.yml` 只覆盖 `server`、`admin-web`、`user-web` 这些 compatibility image lanes，不覆盖根 all-in-one 镜像。

仓库通过 `.github/workflows/monorepo-boundary-guardrail.yml` 的边界护栏（guardrail）阻止 `apps`、`packages`、`scripts`、`bun.lock`、`turbo.json`、`biome.json`、`tsconfig.json` 等前端 workspace 关键根路径重新回流到仓库根目录。

## Subproject Docs

- 后端说明看 `server/README.md`
- 前端说明看 `web/apps/admin/README.md` 与 `web/apps/user/README.md`

但新的贡献者入口以本文件为准。
