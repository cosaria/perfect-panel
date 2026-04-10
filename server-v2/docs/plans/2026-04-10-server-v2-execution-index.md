# Server-V2 Execution Index

> **用途：** 这是 `server-v2` 的执行总控盘，用于逐项放行、逐项验收、逐项暂停。详细实现步骤仍以 `2026-04-10-server-v2-implementation-plan.md` 为真相源，本文件只负责进度控制和子任务分发。

## 使用方式

1. 先确认当前任务的所有前置任务都已通过验收。
2. 只放行一张任务卡进入执行，避免并行漂移。
3. 每张任务卡完成后，仅验收该卡定义的产物和命令。
4. 验收通过后，再手动勾选本索引和对应任务卡。

## 总控状态

- [ ] Task 01: 建立 Go Module 与 CLI 骨架
- [ ] Task 02: 建立测试支架、配置、装配、日志与健康检查基础
- [ ] Task 03: 建立数据库运行基座、 baseline migration 与 seed 链基础
- [ ] Task 04: 建立 OpenAPI 真相源与合同流水线
- [ ] Task 05: 实现 `auth` 与 `access` 主链
- [ ] Task 06: 实现 `system` 与 `catalog` 主链
- [ ] Task 07: 实现 `billing` 交易主链
- [ ] Task 08: 实现 `subscription` 履约主链
- [ ] Task 09: 实现 `node` 宿主机、协议服务、节点与上报主链
- [ ] Task 10: 实现 `outbox`、 worker、 scheduler、投影重建与死信恢复面
- [ ] Task 11: 做端到端验收、合同门禁与发布前收口

## 任务卡入口

| Task | 前置任务 | 执行卡 |
|---|---|---|
| 01 | 无 | [2026-04-10-server-v2-task-01-cli-foundation.md](./2026-04-10-server-v2-task-01-cli-foundation.md) |
| 02 | 01 | [2026-04-10-server-v2-task-02-bootstrap-config-test-support.md](./2026-04-10-server-v2-task-02-bootstrap-config-test-support.md) |
| 03 | 02 | [2026-04-10-server-v2-task-03-db-baseline-seed.md](./2026-04-10-server-v2-task-03-db-baseline-seed.md) |
| 04 | 03 | [2026-04-10-server-v2-task-04-openapi-contract-pipeline.md](./2026-04-10-server-v2-task-04-openapi-contract-pipeline.md) |
| 05 | 04 | [2026-04-10-server-v2-task-05-auth-access.md](./2026-04-10-server-v2-task-05-auth-access.md) |
| 06 | 05 | [2026-04-10-server-v2-task-06-system-catalog.md](./2026-04-10-server-v2-task-06-system-catalog.md) |
| 07 | 06 | [2026-04-10-server-v2-task-07-billing.md](./2026-04-10-server-v2-task-07-billing.md) |
| 08 | 07 | [2026-04-10-server-v2-task-08-subscription.md](./2026-04-10-server-v2-task-08-subscription.md) |
| 09 | 08 | [2026-04-10-server-v2-task-09-node.md](./2026-04-10-server-v2-task-09-node.md) |
| 10 | 09 | [2026-04-10-server-v2-task-10-outbox-worker-projections.md](./2026-04-10-server-v2-task-10-outbox-worker-projections.md) |
| 11 | 10 | [2026-04-10-server-v2-task-11-e2e-release-gates.md](./2026-04-10-server-v2-task-11-e2e-release-gates.md) |

## 全局放行规则

- 没有完成前置任务，不允许启动后置任务。
- 没有通过该卡的必跑验证，不允许标记完成。
- 没有更新对应的 OpenAPI 根入口和合同链，不允许宣称接口任务完成。
- 没有满足该卡的安全约束和审计约束，不允许以“后续补充”通过验收。
- Task 10 前，不允许把投影重建、输出快照和缓存失效硬塞回主事务。

## 真相源

- 详细实施步骤： [2026-04-10-server-v2-implementation-plan.md](./2026-04-10-server-v2-implementation-plan.md)
- 目录规范： [../specs/2026-04-09-server-v2-directory-design.md](../specs/2026-04-09-server-v2-directory-design.md)
- 数据库规范： [../specs/2026-04-09-server-v2-database-design.md](../specs/2026-04-09-server-v2-database-design.md)
- HTTP 规范： [../specs/2026-04-10-server-v2-http-design.md](../specs/2026-04-10-server-v2-http-design.md)
- 运行时规范： [../specs/2026-04-10-server-v2-runtime-workflow-design.md](../specs/2026-04-10-server-v2-runtime-workflow-design.md)
