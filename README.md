# perfect-panel

`perfect-panel` 现在按两个并列工程组织：

- `server/` 是独立 Go 后端工程
- `web/` 是 Bun + Turbo 驱动的前端工作区
- 根目录只保留统一命令入口、OpenAPI 契约、Docker 和 CI 编排

默认发布链仍是嵌入式单镜像路径：前端静态产物构建后复制到 `server/web/*-dist`，再由 Go 二进制一起分发。

## Requirements

- Go 1.21+
- Node.js 20+
- Bun 1.3.0+
- golangci-lint

## Quick Start

从仓库根目录执行：

```bash
make bootstrap
make lint
make test
bun run openapi
```

本地联调时，显式选择一个前端应用：

```bash
make dev APP=admin
make dev APP=user
```

`make dev` 会同时启动：

- `server`: `go run . run --config etc/ppanel.yaml`
- `web`: 你选择的一个前端应用

## Root Command Contract

- `make bootstrap`: 下载 server 依赖并安装 web 依赖
- `make lint`: 运行 server 的 `golangci-lint` 和 `go vet`，再运行 web 的 `bun run lint`
- `make test`: 运行 `cd server && go test ./...`
- `make format`: 运行 server 的 `go fmt` / `goimports`，再运行 web 的 `bun run format`
- `make typecheck`: 运行 web 的 `bun run typecheck`
- `make dev APP=admin|user`: 启动 server 和一个前端应用
- `make embed`: 构建前端静态产物并复制到 `server/web/*-dist`
- `make build-all`: 构建带嵌入静态资源的 server 二进制
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

- `.github/workflows/openapi-governance.yml`
- `.github/workflows/docker-build.yml`

默认 canonical 发布路径是根 `Dockerfile`、根 `docker-compose.yml`、`make embed`、`make build-all`。`server`、`admin-web`、`user-web` 镜像链保留为 compatibility path，需要单独审计。

## Subproject Docs

- 后端说明看 `server/README.md`
- 前端说明看 `web/apps/admin/README.md` 与 `web/apps/user/README.md`

但新的贡献者入口以本文件为准。
