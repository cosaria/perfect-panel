# Task 08: 实现 `subscription` 履约主链

**状态：** `未开始`  
**前置任务：** Task 07  
**对应总计划：** `2026-04-10-server-v2-implementation-plan.md` 的 `Task 8`

## 目标

把交易结果展开成订阅、周期、 addon 周期、 entitlement 和默认节点组授权关系，为节点分配与输出快照建立真相层。

## 执行边界

- 允许触达：
  - `server-v2/internal/domains/subscription/**`
  - `server-v2/internal/platform/db/migrations/0005_subscription.sql`
  - `server-v2/openapi/openapi.yaml`
  - `server-v2/openapi/paths/user/subscriptions.yaml`
  - `server-v2/openapi/paths/admin/subscriptions.yaml`
  - `server-v2/tests/domains/subscription/usecase/*`
  - `server-v2/tests/contract/subscription_api_contract_test.go`

## 关键产物

- `Subscription`
- `SubscriptionPeriod`
- `SubscriptionAddon`
- `SubscriptionAddonPeriod`
- `Entitlement`
- `EntitlementNodeGroup`
- `SubscriptionEvent`
- `SubscriptionOutputSnapshot`

## 必跑验证

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel/server-v2
go test ./tests/domains/subscription/usecase ./tests/contract -run 'TestActivateSubscriptionCreatesPeriodAndOutboxEvent|TestProjectEntitlementsFromPeriodAndAddons|TestUserSubscriptionsContract|TestAdminSubscriptionsContract' -count=1
```

## 放行标准

- `entitlements` 的粒度固定为“单一来源对象 + 单一权益类型 + 单一生效窗口”。
- 默认授权关系通过 `entitlement_node_groups` 表达，不和 override 混用。
- 订阅激活会产出后续投影链所需的 `outbox` 事件。
- 输出快照只作为派生层，不冒充真相层。

## 默认提交点

```bash
git add server-v2/internal/domains/subscription server-v2/internal/platform/db/migrations/0005_subscription.sql server-v2/openapi/openapi.yaml server-v2/openapi/paths/user/subscriptions.yaml server-v2/openapi/paths/admin/subscriptions.yaml server-v2/tests/domains/subscription/usecase server-v2/tests/contract/subscription_api_contract_test.go
git commit -m "feat(server-v2): implement subscription domain"
```

## 完成后进入

- [2026-04-10-server-v2-task-09-node.md](./2026-04-10-server-v2-task-09-node.md)
