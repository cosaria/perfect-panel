# Task 10: 实现 `outbox`、 worker、 scheduler、投影重建与死信恢复面

**状态：** `未开始`  
**前置任务：** Task 09  
**对应总计划：** `2026-04-10-server-v2-implementation-plan.md` 的 `Task 10`

## 目标

把真相层和投影层之间的异步桥接正式落地，包括 `outbox` 分发、 worker 执行、 scheduler 驱动、投影重建、死信记录和死信恢复面。

## 执行边界

- 允许触达：
  - `server-v2/internal/platform/queue/**`
  - `server-v2/internal/domains/subscription/jobs/**`
  - `server-v2/internal/domains/node/jobs/**`
  - `server-v2/internal/domains/system/usecase/{list_dead_letters,replay_dead_letter,discard_dead_letter}.go`
  - `server-v2/internal/domains/system/api/admin_dead_letters.go`
  - `server-v2/internal/app/runtime/{worker,scheduler,outbox_dispatcher}.go`
  - `server-v2/openapi/openapi.yaml`
  - `server-v2/openapi/paths/admin/dead_letters.yaml`
  - `server-v2/tests/integration/runtime/*`
  - `server-v2/tests/contract/dead_letters_api_contract_test.go`

## 关键产物

- `asynq_client.go`
- `dead_letter.go`
- `dead_letter_store.go`
- `entitlement_projection_job.go`
- `output_snapshot_job.go`
- `assignment_rebuild_job.go`
- `usage_rollup_job.go`
- `list_dead_letters.go`
- `replay_dead_letter.go`
- `discard_dead_letter.go`
- `worker.go`
- `scheduler.go`
- `outbox_dispatcher.go`

## 必跑验证

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel/server-v2
go test ./tests/integration/runtime ./tests/contract -run 'TestOutboxDispatchCreatesProjectionTask|TestFailedProjectionMovesToDeadLetterAfterRetries|TestReplayDeadLetterRequeuesTask|TestAdminDeadLettersContract' -count=1
```

## 放行标准

- 事务内只写不可替代事实和 `outbox_events`，不把投影重建硬塞回主事务。
- `outbox` payload 版本化。
- 死信能力至少支持：列表、详情、重放、丢弃、关联对象定位、重放后状态回写。
- worker 和 scheduler 使用统一 runtime 装配，不绕开 bootstrap。

## 默认提交点

```bash
git add server-v2/internal/platform/queue server-v2/internal/domains/subscription/jobs server-v2/internal/domains/node/jobs server-v2/internal/domains/system/usecase server-v2/internal/domains/system/api server-v2/internal/app/runtime server-v2/openapi/openapi.yaml server-v2/openapi/paths/admin/dead_letters.yaml server-v2/tests/integration/runtime server-v2/tests/contract/dead_letters_api_contract_test.go
git commit -m "feat(server-v2): implement outbox and async workflows"
```

## 完成后进入

- [2026-04-10-server-v2-task-11-e2e-release-gates.md](./2026-04-10-server-v2-task-11-e2e-release-gates.md)
