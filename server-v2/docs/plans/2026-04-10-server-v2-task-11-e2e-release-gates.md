# Task 11: 做端到端验收、合同门禁与发布前收口

**状态：** `未开始`  
**前置任务：** Task 10  
**对应总计划：** `2026-04-10-server-v2-implementation-plan.md` 的 `Task 11`

## 目标

在全部主链完成后，以端到端测试、竞态测试、 OpenAPI 合同和发布前清理作为最终放行门，确认 `server-v2` 第一版达到可运行、可验证、可继续演进的状态。

## 执行边界

- 允许触达：
  - `server-v2/tests/smoke/api_boot_test.go`
  - `server-v2/tests/integration/e2e/order_to_subscription_test.go`
  - `server-v2/tests/integration/e2e/node_projection_refresh_test.go`
  - `server-v2/.gitignore`
- 不允许再新增新的核心领域范围；这里只做验收、补洞和收口。

## 关键产物

- `api_boot_test.go`
- `order_to_subscription_test.go`
- `node_projection_refresh_test.go`
- 最终全量验证记录

## 必跑验证

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel/server-v2
go test ./tests/integration/e2e -run TestPaidOrderActivatesSubscriptionAndProjectsNodes -count=1
go test ./... -count=1
go test ./... -race
make contract
```

## 放行标准

- `paid order -> subscription activation -> entitlement projection -> node assignments` 主链闭环成立。
- 全量 Go 测试通过。
- 竞态测试通过。
- `make contract` 通过，且 `openapi-ts` 生成无漂移。
- 没有把新的未规划领域或运行时复杂度偷偷塞进最后一张卡。

## 默认提交点

```bash
git add server-v2
git commit -m "feat(server-v2): deliver first runnable v2 backend"
```

## 结束条件

- 本卡通过后，`server-v2` 第一版进入可执行实施阶段，可按任务卡逐项派发执行。
