# perfect-panel

perfect-panel 现在以一个产品级根目录协调整个仓库：

- server 是 Go 后端服务。
- web 是 Bun + Turbo 驱动的前端工作区。
- 根目录负责统一开发入口、统一 CI 入口和发布配对信息。

这次 phase one 只建立 monorepo 基线，不改变现有运行时部署拓扑，也不引入 go.work 或 API contract automation。

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
- `make dev APP=admin|user`: 启动 server 和一个前端应用

本地 `make test` 保持轻量。更完整的前端验证放在根级 GitHub Actions workflow 中执行。

## Repository Layout

```text
.
├── Makefile
├── release-manifest.json
├── server/
└── web/
```

- `server/`: 后端服务、API、构建与运行入口
- `web/`: 前端工作区，包含 admin 和 user 两个应用
- `release-manifest.json`: 记录已知可发布的 server/web 配对

## CI

根级 workflow 在 `.github/workflows/monorepo-check.yml` 中：

- `server-validate`
- `web-admin-validate`
- `web-user-validate`
- `monorepo-summary`

迁移窗口内，这个 workflow 先作为观察用入口，不应在分支保护里立即设为 required。升级门槛是连续 2 次 main 绿加至少 1 个功能 PR 绿。

## Release Pairing

根级 [release-manifest.json](release-manifest.json) 是产品级配对文档，不是每个普通开发 PR 都要更新的文件。

只有在形成一个已知可发布的 server/web 组合时，才更新它。

## Subproject Docs

- 后端说明看 `server/README.md`
- 前端说明看 `web/README.md`

但新的贡献者入口以本文件为准。