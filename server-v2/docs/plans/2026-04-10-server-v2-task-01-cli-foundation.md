# Task 01: 建立 Go Module 与 CLI 骨架

**状态：** `未开始`  
**前置任务：** 无  
**对应总计划：** `2026-04-10-server-v2-implementation-plan.md` 的 `Task 1`

## 目标

建立 `server-v2` 的独立 Go module 和唯一 `cmd/server` 入口，形成 `serve-api / serve-worker / serve-scheduler / migrate / seed-required / seed-demo` 这 6 个运行模式的最小骨架。

## 执行边界

- 只允许触达：
  - `server-v2/go.mod`
  - `server-v2/cmd/server/*`
  - `server-v2/internal/app/runtime/*`
  - `server-v2/tests/smoke/cli_bootstrap_test.go`
- 不允许提前引入业务域实现。

## 关键产物

- `go.mod`
- `cmd/server/main.go`
- `cmd/server/root.go`
- `cmd/server/serve_api.go`
- `cmd/server/serve_worker.go`
- `cmd/server/serve_scheduler.go`
- `cmd/server/migrate.go`
- `cmd/server/seed_required.go`
- `cmd/server/seed_demo.go`
- `internal/app/runtime/modes.go`

## 必跑验证

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel/server-v2
go test ./tests/smoke -run 'TestCliBootstrapFilesExist|TestCliHelpBoots' -count=1
```

## 放行标准

- `go run ./cmd/server --help` 成功。
- 六个运行模式命令都已注册到 root command。
- `internal/app/runtime/modes.go` 中的模式常量齐全。
- 没有业务逻辑混入 CLI 入口。

## 默认提交点

```bash
git add server-v2/go.mod server-v2/cmd/server server-v2/internal/app/runtime server-v2/tests/smoke/cli_bootstrap_test.go
git commit -m "feat(server-v2): scaffold module and cli skeleton"
```

## 完成后进入

- [2026-04-10-server-v2-task-02-bootstrap-config-test-support.md](./2026-04-10-server-v2-task-02-bootstrap-config-test-support.md)
