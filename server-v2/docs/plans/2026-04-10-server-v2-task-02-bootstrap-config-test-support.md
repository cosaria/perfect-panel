# Task 02: 建立测试支架、配置、装配、日志与健康检查基础

**状态：** `未开始`  
**前置任务：** Task 01  
**对应总计划：** `2026-04-10-server-v2-implementation-plan.md` 的 `Task 2`

## 目标

把后续所有任务都会依赖的测试支架、配置优先级、容器装配、日志基础和健康检查骨架一次立稳，避免后续任务重复发明测试 helper 和启动逻辑。

## 执行边界

- 只允许触达：
  - `server-v2/tests/support/*`
  - `server-v2/internal/platform/config/*`
  - `server-v2/internal/platform/observability/*`
  - `server-v2/internal/app/wiring/*`
  - `server-v2/internal/app/bootstrap/*`
  - `server-v2/internal/platform/http/health/*`
  - `server-v2/tests/app/bootstrap/*`
  - `server-v2/tests/smoke/health_handler_test.go`
- 不允许在这一卡里引入数据库 schema 或领域业务规则。

## 关键产物

- `tests/support/db.go`
- `tests/support/services.go`
- `tests/support/runtime.go`
- `tests/support/e2e.go`
- `tests/support/fixtures.go`
- `internal/platform/config/config.go`
- `internal/platform/config/load.go`
- `internal/platform/observability/logger.go`
- `internal/app/wiring/container.go`
- `internal/app/bootstrap/bootstrap.go`
- `internal/platform/http/health/handler.go`

## 必跑验证

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel/server-v2
go test ./tests/app/bootstrap ./tests/smoke -run 'TestLoadPrefersCLIOverEnvAndFile|TestHealthHandlerReturnsEnvelope' -count=1
```

## 放行标准

- 配置优先级固定为 `CLI > env > file`。
- `tests/support` 已可被后续各领域测试直接复用。
- 健康检查返回符合 HTTP 规范的 `data + meta` envelope。
- 容器与 bootstrap 只做装配，不吞业务语义。

## 默认提交点

```bash
git add server-v2/tests/support server-v2/internal/platform/config server-v2/internal/platform/observability server-v2/internal/app/wiring server-v2/internal/app/bootstrap server-v2/internal/platform/http/health server-v2/tests/app/bootstrap server-v2/tests/smoke/health_handler_test.go
git commit -m "feat(server-v2): add config and bootstrap foundation"
```

## 完成后进入

- [2026-04-10-server-v2-task-03-db-baseline-seed.md](./2026-04-10-server-v2-task-03-db-baseline-seed.md)
