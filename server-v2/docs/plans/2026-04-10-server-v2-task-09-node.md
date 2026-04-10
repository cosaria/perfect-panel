# Task 09: 实现 `node` 宿主机、协议服务、节点与上报主链

**状态：** `未开始`  
**前置任务：** Task 08  
**对应总计划：** `2026-04-10-server-v2-implementation-plan.md` 的 `Task 9`

## 目标

实现宿主机、协议服务、用户节点、节点组、授权结果、人工覆盖规则、心跳和 usage 三段式上报链，使节点侧和客户端侧都能读取稳定结果。

## 执行边界

- 允许触达：
  - `server-v2/internal/domains/node/**`
  - `server-v2/internal/platform/http/middleware/node_auth.go`
  - `server-v2/internal/platform/db/migrations/0006_node.sql`
  - `server-v2/openapi/openapi.yaml`
  - `server-v2/openapi/paths/node/registrations.yaml`
  - `server-v2/openapi/paths/node/heartbeats.yaml`
  - `server-v2/openapi/paths/node/usage_reports.yaml`
  - `server-v2/tests/domains/node/usecase/*`
  - `server-v2/tests/contract/node_api_contract_test.go`

## 关键产物

- `Host / HostProtocol / Node`
- `NodeGroup / NodeGroupMember`
- `NodeAssignment / NodeAssignmentOverride`
- `NodeRegistration / NodeHeartbeat / NodeUsageReport`
- `OnlineSession / NodeKey`
- `register_node / record_heartbeat / report_usage / rebuild_assignments`
- `node_auth.go`

## 必跑验证

```bash
cd /Users/admin/Codes/ProxyCode/perfect-panel/server-v2
go test ./tests/domains/node/usecase ./tests/contract -run 'TestRegisterNodeBindsToHostProtocol|TestReportUsagePersistsRawAndBillableValues|TestHeartbeatRequiresValidNodeSignature|TestNodeRegistrationsContract|TestNodeHeartbeatsContract|TestNodeUsageReportsContract' -count=1
```

## 放行标准

- `nodeAuth` 绑定单节点主体，校验 `timestamp` 时间窗和 `nonce` 单次使用。
- 节点签名必须覆盖 `method + path + body digest`。
- usage 同时保存原始值、计费模式、倍率、计费前基数和最终计费结果。
- 在线设备限制先以在线连接记录和活跃连接数判断实现。
- 授权结果服务于用户节点，不直接绑定宿主机或协议服务。

## 默认提交点

```bash
git add server-v2/internal/domains/node server-v2/internal/platform/http/middleware/node_auth.go server-v2/internal/platform/db/migrations/0006_node.sql server-v2/openapi/openapi.yaml server-v2/openapi/paths/node/registrations.yaml server-v2/openapi/paths/node/heartbeats.yaml server-v2/openapi/paths/node/usage_reports.yaml server-v2/tests/domains/node/usecase server-v2/tests/contract/node_api_contract_test.go
git commit -m "feat(server-v2): implement node domain"
```

## 完成后进入

- [2026-04-10-server-v2-task-10-outbox-worker-projections.md](./2026-04-10-server-v2-task-10-outbox-worker-projections.md)
