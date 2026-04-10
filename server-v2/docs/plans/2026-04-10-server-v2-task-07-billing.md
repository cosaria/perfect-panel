# Task 07: 实现 `billing` 交易主链

**状态：** `未开始`  
**前置任务：** Task 06  
**对应总计划：** `2026-04-10-server-v2-implementation-plan.md` 的 `Task 7`

## 目标

建立订单、订单项快照、支付、支付回调、退款和幂等保护主链，为订阅履约提供稳定的交易真相。

## 执行边界

- 允许触达：
  - `server-v2/internal/domains/billing/**`
  - `server-v2/internal/platform/http/idempotency/guard.go`
  - `server-v2/internal/platform/db/migrations/0004_billing.sql`
  - `server-v2/openapi/openapi.yaml`
  - `server-v2/openapi/paths/user/orders.yaml`
  - `server-v2/openapi/paths/public/payment_callbacks.yaml`
  - `server-v2/openapi/paths/admin/refunds.yaml`
  - `server-v2/tests/domains/billing/usecase/*`
  - `server-v2/tests/contract/billing_api_contract_test.go`

## 关键产物

- `Order / OrderItem / Payment / PaymentEvent / Refund / RefundItem`
- `create_order / record_payment / record_payment_callback / create_refund`
- `idempotency/guard.go`

## 必跑验证

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel/server-v2
go test ./tests/domains/billing/usecase ./tests/contract -run 'TestCreateOrderPersistsItemSnapshots|TestCreateOrderRejectsIdempotencyKeyWithDifferentPayload|TestRecordPaymentCallbackDeduplicatesProviderEvent|TestRecordPaymentCallbackRejectsInvalidSignature|TestUserOrdersContract|TestPaymentCallbacksContract|TestAdminRefundsContract' -count=1
```

## 放行标准

- 订单项快照是历史交易真相，不回查当前商品价格。
- 幂等作用域固定为 `subject + method + path template + body digest`，窗口固定 24 小时。
- 同 key 同摘要重放返回首次成功响应；同 key 不同摘要返回 `409 Conflict`。
- 支付回调验签失败必须拒绝，不能入账、不能履约、不能写成功事件。
- `(provider, provider_event_key)` 只允许首次入账。

## 默认提交点

```bash
git add server-v2/internal/domains/billing server-v2/internal/platform/http/idempotency server-v2/internal/platform/db/migrations/0004_billing.sql server-v2/openapi/openapi.yaml server-v2/openapi/paths/user/orders.yaml server-v2/openapi/paths/public/payment_callbacks.yaml server-v2/openapi/paths/admin/refunds.yaml server-v2/tests/domains/billing/usecase server-v2/tests/contract/billing_api_contract_test.go
git commit -m "feat(server-v2): implement billing domain"
```

## 完成后进入

- [2026-04-10-server-v2-task-08-subscription.md](./2026-04-10-server-v2-task-08-subscription.md)
