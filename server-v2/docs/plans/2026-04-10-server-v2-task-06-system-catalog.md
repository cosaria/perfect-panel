# Task 06: 实现 `system` 与 `catalog` 主链

**状态：** `未开始`  
**前置任务：** Task 05  
**对应总计划：** `2026-04-10-server-v2-implementation-plan.md` 的 `Task 6`

## 目标

建立系统配置、后台操作审计、公开套餐列表和后台套餐维护主链，为后续交易与履约提供可售商品和可运营配置。

## 执行边界

- 允许触达：
  - `server-v2/internal/domains/system/**`
  - `server-v2/internal/domains/catalog/**`
  - `server-v2/internal/platform/db/migrations/0003_system_catalog.sql`
  - `server-v2/openapi/openapi.yaml`
  - `server-v2/openapi/paths/public/plans.yaml`
  - `server-v2/openapi/paths/admin/plans.yaml`
  - `server-v2/tests/domains/catalog/usecase/*`
  - `server-v2/tests/domains/system/usecase/*`
  - `server-v2/tests/contract/catalog_api_contract_test.go`

## 关键产物

- `SystemSetting`
- `AdminOperationLog`
- `Plan / PlanVariant / PlanAddon`
- `get_settings / update_settings / record_admin_operation`
- `list_public_plans / upsert_plan`

## 必跑验证

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel/server-v2
go test ./tests/domains/catalog/usecase ./tests/domains/system/usecase ./tests/contract -run 'TestListPublicPlansReturnsActiveVariants|TestUpdateSettingsRecordsAdminOperationLog|TestPublicPlansContract|TestAdminPlansContract' -count=1
```

## 放行标准

- 公开套餐列表只返回可售、有效的变体。
- 所有 `admin` 写接口都必须记录 `admin_operation_log`。
- `system_settings` 继续遵守 `scope + key` 单表模型。
- `catalog` 不越界承担交易或履约语义。

## 默认提交点

```bash
git add server-v2/internal/domains/system server-v2/internal/domains/catalog server-v2/internal/platform/db/migrations/0003_system_catalog.sql server-v2/openapi/openapi.yaml server-v2/openapi/paths/public/plans.yaml server-v2/openapi/paths/admin/plans.yaml server-v2/tests/domains/catalog/usecase server-v2/tests/domains/system/usecase server-v2/tests/contract/catalog_api_contract_test.go
git commit -m "feat(server-v2): implement system and catalog domains"
```

## 完成后进入

- [2026-04-10-server-v2-task-07-billing.md](./2026-04-10-server-v2-task-07-billing.md)
