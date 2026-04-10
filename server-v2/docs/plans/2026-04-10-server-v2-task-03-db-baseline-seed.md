# Task 03: 建立数据库运行基座、 baseline migration 与 seed 链基础

**状态：** `未开始`  
**前置任务：** Task 02  
**对应总计划：** `2026-04-10-server-v2-implementation-plan.md` 的 `Task 3`

## 目标

建立 PostgreSQL 连接、事务包装、 baseline migration、 `schema_revisions` 门禁和 `seed-required / seed-demo` 的最小链路，让 `server-v2` 能从空库稳定起步。

## 执行边界

- 只允许触达：
  - `server-v2/internal/platform/db/*`
  - `server-v2/internal/platform/db/migrations/*`
  - `server-v2/internal/platform/db/seeds/*`
  - `server-v2/internal/app/bootstrap/bootstrap.go`
  - `server-v2/cmd/server/migrate.go`
  - `server-v2/cmd/server/seed_required.go`
  - `server-v2/cmd/server/seed_demo.go`
  - `server-v2/tests/integration/db/*`
- 不允许提前创建 `auth/access` 之后才会落地的权限表行为逻辑。

## 关键产物

- `connect.go`
- `transaction.go`
- `migrate.go`
- `schema_version.go`
- `seeds/required.go`
- `seeds/demo.go`
- `migrations/0001_baseline.sql`

## 必跑验证

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel/server-v2
go test ./tests/integration/db -run 'TestMigrateAppliesBaselineTables|TestServeFailsWhenSchemaVersionMismatches' -count=1
```

## 放行标准

- baseline schema 至少包含 `users`、`roles`、`system_settings`、`outbox_events`、`schema_revisions`。
- `migrate` 成功后会写入目标 schema version。
- `serve-api / serve-worker / serve-scheduler` 在 schema version 不匹配时明确失败。
- `seed-required` 当前只负责 baseline 可启动最小集，不提前偷塞后续领域种子。

## 默认提交点

```bash
git add server-v2/internal/platform/db server-v2/cmd/server/migrate.go server-v2/cmd/server/seed_required.go server-v2/cmd/server/seed_demo.go server-v2/tests/integration/db
git commit -m "feat(server-v2): add migration and seed foundation"
```

## 完成后进入

- [2026-04-10-server-v2-task-04-openapi-contract-pipeline.md](./2026-04-10-server-v2-task-04-openapi-contract-pipeline.md)
