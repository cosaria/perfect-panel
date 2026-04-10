# Task 04: 建立 OpenAPI 真相源与合同流水线

**状态：** `未开始`  
**前置任务：** Task 03  
**对应总计划：** `2026-04-10-server-v2-implementation-plan.md` 的 `Task 4`

## 目标

建立 `OpenAPI 3.1` 单一真相源、 `Redocly` 合同流水线和 `@hey-api/openapi-ts` 生成链，冻结 `operationId`、公共参数和公共 schema 的最小基线。

## 执行边界

- 只允许触达：
  - `server-v2/Makefile`
  - `server-v2/openapi/**`
  - `server-v2/tests/contract/openapi_bundle_test.go`
  - `server-v2/tests/contract/openapi_operation_id_test.go`
- 不允许在这一卡里写具体领域业务实现。

## 关键产物

- `openapi/openapi.yaml`
- `openapi/paths/public/sessions.yaml`
- `openapi/paths/public/verification_tokens.yaml`
- `openapi/paths/public/password_reset_requests.yaml`
- `openapi/paths/public/password_resets.yaml`
- `openapi/paths/user/me_sessions.yaml`
- `openapi/paths/admin/users.yaml`
- `openapi/paths/node/usage_reports.yaml`
- `openapi/components/parameters/{page,page_size,q,sort,order}.yaml`
- `openapi/components/schemas/common/{problem,validation_problem,pagination_meta,money}.yaml`
- `openapi/components/headers/{idempotency_key,node_key_id,node_timestamp,node_nonce,node_signature}.yaml`
- `openapi/components/security/security.yaml`
- `openapi/redocly.yaml`

## 必跑验证

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel/server-v2
make contract
go test ./tests/contract -run 'TestContractPipelinePasses|TestOperationIDsAreExplicit' -count=1
```

## 放行标准

- `make contract` 是唯一合同入口，必须完成 `lint -> bundle -> openapi-ts`。
- 所有 operation 都有显式、全局唯一的 `operationId`。
- `Problem`、`ValidationProblem`、`PaginationMeta`、`Money` 已冻结为公共组件。
- `sessionAuth` 与 `nodeAuth` 已在合同层声明。

## 默认提交点

```bash
git add server-v2/Makefile server-v2/openapi server-v2/tests/contract/openapi_bundle_test.go server-v2/tests/contract/openapi_operation_id_test.go
git commit -m "feat(server-v2): add openapi source and contract pipeline"
```

## 完成后进入

- [2026-04-10-server-v2-task-05-auth-access.md](./2026-04-10-server-v2-task-05-auth-access.md)
