# Server-V2 HTTP 规范

## 背景

`server-v2/` 已经完成了目录规范与数据库规范，但仍缺少一份稳定的 HTTP 契约规范，来回答：

- 系统如何对外暴露 `public / user / admin / node` 四个调用面
- OpenAPI 如何成为文档、客户端生成和 AI 修改的单一真相源
- 请求、响应、错误、认证、异步和兼容性如何保持一致

这份规范默认建立在以下文档之上：

- [2026-04-09-server-v2-directory-design.md](/Users/admin/Codes/ProxyCode/perfect-panel/server-v2/docs/specs/2026-04-09-server-v2-directory-design.md)
- [2026-04-09-server-v2-database-design.md](/Users/admin/Codes/ProxyCode/perfect-panel/server-v2/docs/specs/2026-04-09-server-v2-database-design.md)

它不追求“理论上最通用的 API 风格”，而是优先服务 4 件事：

1. `OpenAPI 3.1` 真相源
2. [Redocly](https://redocly.com/) 文档生成
3. `@hey-api/openapi-ts` 客户端生成
4. AI 在后续实现与维护中的稳定识别

## 目标

1. 定义 `server-v2` 第一版 HTTP API 的统一路径、命名、认证、响应和错误规则。
2. 为 `public / user / admin / node` 四个调用面提供稳定的边界。
3. 定义 OpenAPI 源码组织方式，使其可直接支撑 Redocly 与 `@hey-api/openapi-ts`。
4. 明确 `v1` 的兼容边界，避免后续实现阶段再临场发明 API 规则。

## 非目标

1. 本文不列出每个业务域的完整接口清单。
2. 本文不展开 OpenAPI 列级 schema、每个字段的完整定义或逐接口示例正文。
3. 本文不设计 Web 前端页面流或节点协议内部实现细节。
4. 本文不取代实现计划；它只定义 HTTP 契约规则。

## 核心判断

这份 HTTP 规范采用：

**单一 OpenAPI 真相源 + 四个调用面 + 资源化路径 + 统一成功响应体 + 结构化错误 + 显式兼容规则**

它对应 10 条工作原则：

1. 维护一份 OpenAPI 3.1 总规范，而不是多份彼此独立的 API 文档。
2. 路径先按调用面分流，再按资源组织，不走 RPC 风格。
3. 所有带响应体的成功响应统一使用 `data + meta` 外层结构，`v1` 不使用 `204 No Content`。
4. 所有错误统一使用 `Problem Details`，校验错误追加统一字段级错误数组。
5. `user / admin` 共享 `sessionAuth`，`node` 使用独立机器认证方案，权限要求结构化标注。
6. 所有非幂等写操作统一支持 `Idempotency-Key`。
7. 更新统一使用 `PATCH` 业务型局部更新对象，不使用 `PUT`、`JSON Merge Patch` 或 `JSON Patch`。
8. 少数异步写操作统一返回 `202 Accepted` 与标准任务受理对象。
9. `v1` 内只允许向后兼容演进，破坏性变更必须进入 `v2`。
10. 规范优先追求可预测性，而不是为了少写几行 schema 牺牲可读性。

## OpenAPI 真相源

### 总体要求

`server-v2` 第一版 HTTP API 以 **一份 OpenAPI 3.1 规范** 作为唯一真相源。

这份真相源同时服务于：

- Redocly 文档输出
- `@hey-api/openapi-ts` 客户端生成
- 实现阶段的 handler / DTO 契约对齐
- AI 后续修改接口时的稳定参考

不允许维护四份彼此独立的 `public / user / admin / node` 规范。

### 真相源与派生产物流水线

`server-v2` 的 HTTP 契约遵守固定流水线：

1. 人工维护 `server-v2/openapi/openapi.yaml` 及其 `$ref` 子文件，作为唯一可编辑真相源。
2. 使用 Redocly 对真相源执行 lint 与 bundle。
3. 基于 bundle 产出文档站可消费的规范文件。
4. 基于同一份 bundle 产出 `@hey-api/openapi-ts` 客户端。
5. `server-v2` 的 Go 实现对齐这份规范，但不反向生成新的 OpenAPI 真相源。

这条规则的含义是：

- 文档和客户端都从 OpenAPI 真相源派生
- 服务实现服从 OpenAPI 契约
- 不允许同时维护“Go 代码是真相源”和“OpenAPI 文件是真相源”两套体系

如果未来工具链完整支持 `OpenAPI 3.2`，可以在不改变这条流水线的前提下升级；第一版基线固定为 `3.1`，以确保 Redocly 和客户端生成链路稳定。

### 源码组织

OpenAPI 源码采用：

- 一份总入口
- `paths/` 按调用面拆分
- `components/` 按业务域拆分

推荐结构：

```text
server-v2/openapi/
├── openapi.yaml
├── paths/
│   ├── public/
│   ├── user/
│   ├── admin/
│   └── node/
└── components/
    ├── schemas/
    │   ├── auth/
    │   ├── access/
    │   ├── billing/
    │   ├── catalog/
    │   ├── subscription/
    │   ├── node/
    │   ├── system/
    │   └── common/
    ├── parameters/
    ├── responses/
    ├── headers/
    └── security/
```

这个组织方式的判断标准是：

- path 第一跳先回答“谁在调用”
- schema 第一跳再回答“这是什么业务对象”

## 调用面与路径骨架

### 四个调用面

第一版固定 4 个 HTTP 调用面：

1. `public`
   - 匿名读取
   - 少量匿名创建
   - 例如注册、登录、验证码发起、密码重置发起
2. `user`
   - 当前登录用户自助操作
3. `admin`
   - 后台运营和管理操作
4. `node`
   - 宿主机 / 节点到控制面的机器接口

### 路径前缀

所有路径统一采用版本前缀：

- `/api/v1/public/...`
- `/api/v1/user/...`
- `/api/v1/admin/...`
- `/api/v1/node/...`

第一版不允许混入：

- `/public/...`
- `/admin/...`
- `/node/...`
- 无版本前缀路径

### 资源命名

路径与 schema 命名统一采用：

- path 使用复数资源名
- schema 使用单数资源名

例如：

- path：`/api/v1/user/subscriptions`
- schema：`SubscriptionDetail`

### 路径风格

HTTP 规范默认采用资源化路径，而不是动作化路径。

优先表达：

- 创建资源
- 读取资源
- 更新资源
- 删除资源
- 查询子资源集合

不推荐使用：

- `/sign-in`
- `/create-order`
- `/report-usage`

而应优先表达成：

- `POST /api/v1/public/sessions`
- `POST /api/v1/user/orders`
- `POST /api/v1/node/usage-reports`

### 嵌套深度

路径嵌套最多允许一层。

允许：

- `/api/v1/user/me/subscriptions`
- `/api/v1/admin/users/{userId}/subscriptions`

不鼓励：

- `/api/v1/admin/users/{userId}/subscriptions/{subscriptionId}/addons/{addonId}`

如果关系过深，优先改为：

- 独立资源
- 或查询参数筛选

### `user` 面与 `me`

`user` 面显式使用 `me` 作为当前登录用户资源入口。

例如：

- `/api/v1/user/me`
- `/api/v1/user/me/sessions`
- `/api/v1/user/me/subscriptions`

### `admin` 面

`admin` 面一律使用显式资源 ID，不使用 `me`。

例如：

- `/api/v1/admin/users/{userId}`
- `/api/v1/admin/nodes/{nodeId}`
- `/api/v1/admin/orders/{orderId}`

### `node` 面

`node` 面也遵守资源化 JSON over HTTP 规则。

例如：

- `POST /api/v1/node/usage-reports`
- `POST /api/v1/node/heartbeats`
- `POST /api/v1/node/registrations`

它不拥有单独的“动作路径哲学”。

## 认证、权限与安全方案

### 安全方案

OpenAPI 中显式维护两套安全方案：

- `sessionAuth`
- `nodeAuth`

其中：

- `user / admin` 共享同一套 Bearer Token 会话语义
- `node` 使用独立机器认证语义

第一版不允许把 `node` 也压成笼统的 Bearer Token。

`sessionAuth` 只表达“这是一个已登录主体的会话认证”；它不表达调用面权限差异。`admin` 能否访问某个接口，仍由：

- 路径调用面
- `x-requiredPermissions`
- 必要时 `x-requiredRoles`

共同决定。

在 OpenAPI `securitySchemes` 中：

- `sessionAuth` 固定为 `type: http`、`scheme: bearer`
- `nodeAuth` 固定为 `type: apiKey`、`in: header`、`name: X-Node-Key-Id`

### `nodeAuth`

`nodeAuth` 在协议上固定为一套签名式机器认证，而不是普通 Bearer Token。

第一版至少包含以下 header：

- `X-Node-Key-Id`
- `X-Node-Timestamp`
- `X-Node-Nonce`
- `X-Node-Signature`

其中 `X-Node-Key-Id` 由 `nodeAuth` security scheme 表达，其余 3 个 header 作为共享参数组件挂载到所有 `node` 写接口与受保护读接口。

其中：

- `X-Node-Key-Id` 标识当前节点或宿主机使用的可吊销密钥
- `X-Node-Timestamp` 提供有限时间窗
- `X-Node-Nonce` 提供单次请求随机性
- `X-Node-Signature` 对请求摘要签名

规范层冻结以下最小安全语义：

- 密钥必须绑定到单一节点主体或单一宿主机主体
- 时间窗默认不超过 5 分钟
- `nonce` 在有效时间窗内不得重复使用
- 节点密钥必须支持单独吊销，不影响其他节点主体
- 签名覆盖 method、规范化 path、时间戳、nonce 和请求体摘要

### 权限标注

认证只回答“你是谁”，不回答“你需要什么权限”。

因此，需要权限控制的 operation 必须额外携带结构化扩展字段，例如：

- `x-requiredPermissions`
- 必要时 `x-requiredRoles`

不允许只在自然语言描述里写“需要某权限”。

### 幂等键

所有非幂等写操作统一支持：

- `Idempotency-Key`

它应作为公共 header 组件存在于 OpenAPI 中，并至少覆盖：

- 创建订单
- 支付相关动作
- 订阅续费
- addon 加购
- 高风险后台写操作
- 其他可重试的创建型请求

幂等语义冻结如下：

- 幂等判定作用域至少包含：调用主体、HTTP method、路径模板、请求摘要
- 相同作用域内，`Idempotency-Key` 的默认去重时间窗为 24 小时
- 同一 `Idempotency-Key` 在时间窗内如果命中完全相同的请求摘要，应返回第一次成功受理的响应
- 同一 `Idempotency-Key` 在时间窗内如果命中不同请求摘要，必须返回 `409 Conflict`
- `409` 响应仍使用 Problem Details，并携带稳定的业务错误码，例如 `idempotency_conflict`

## 请求与响应契约

### 成功响应 envelope

所有带响应体的成功响应统一采用：

```json
{
  "data": {},
  "meta": {}
}
```

规则如下：

- `data` 始终存在
- `meta` 始终存在
- 不因为“这个接口很简单”就省略 envelope

但这条规则不意味着实现一个全局泛型 `SuccessEnvelope` schema。为了保证 `@hey-api/openapi-ts` 的生成类型精度，每个 operation 仍应拥有自己显式命名的成功响应 schema，例如：

- `UserSubscriptionListResponse`
- `AdminNodeDetailResponse`
- `PublicSessionCreateResponse`

公共组件只抽取真正稳定复用的片段，例如：

- `PaginationMeta`
- `Money`
- `Problem`
- `ValidationProblem`

### 列表响应

列表响应也遵守统一 envelope，但 `data` 通常是数组或列表对象，`meta` 使用统一分页元数据。

统一分页字段：

- `page`
- `pageSize`
- `total`
- `totalPages`
- `hasNext`
- `hasPrev`

对应的列表请求参数统一使用：

- `page`
- `pageSize`

第一版不引入 cursor 分页。

### 错误响应

所有错误统一使用 Problem Details 风格对象。

至少包含：

- `type`
- `title`
- `status`
- `detail`
- 可选 `code`
- 可选 `instance`

### 校验错误

业务校验失败采用：

**Problem Details + 统一字段级错误数组**

推荐字段级错误对象至少包含：

- `field`
- `code`
- `message`

### 字段命名

所有 JSON 请求与响应字段统一使用 `camelCase`。

### 时间字段

所有时间字段统一使用 RFC 3339 / ISO 8601 UTC 字符串。

例如：

- `createdAt`
- `updatedAt`
- `acceptedAt`

### 金额字段

所有对外金额统一暴露为结构化金额对象，而不是裸整数或字符串金额。

推荐形态：

```json
{
  "amount": 1299,
  "currency": "CNY"
}
```

其中 `amount` 仍表示最小货币单位。

### 资源标识

所有资源主标识统一为字符串 UUID。

例如：

- `id`

业务编号单独暴露，例如：

- `orderNo`
- `paymentNo`
- `subscriptionNo`

### 空值策略

默认最小化 `null`。

规则：

- 优先使用字段缺席，而不是 `null`
- 列表返回空数组，不返回 `null`
- 只有当“字段缺席”和“显式空值”语义不同，才允许 `nullable`

### 状态字段

所有对外状态字段必须是显式 enum。

不允许：

- `status: string`

而必须在 OpenAPI 中冻结：

- `subscriptionStatus`
- `paymentStatus`
- `nodeStatus`
- 等其他状态字段的可选值

### 请求与响应模型分离

请求模型和响应模型强制分离。

例如：

- `OrderCreateRequest`
- `OrderUpdateRequest`
- `OrderListItem`
- `OrderDetail`

不要让同一个 schema 同时承担输入和输出语义。

### 列表项与详情项

默认区分：

- `XxxListItem`
- `XxxDetail`

只有极简单资源才允许列表和详情共用同一表示。

## 方法语义

### 创建

创建使用 `POST`。

### 更新

更新统一使用 `PATCH`，并采用业务型局部更新对象。

例如：

```json
{
  "displayName": "香港 02 - IEPL 专线",
  "status": "disabled"
}
```

第一版不使用：

- `PUT`
- `JSON Merge Patch`
- `JSON Patch`

### 删除

`DELETE` 只用于真正删除资源。

大多数生命周期变更，例如：

- 归档
- 禁用
- 取消
- 停用

统一使用 `PATCH` 修改资源状态，而不是滥用 `DELETE`。

### 写操作成功后的返回

创建和更新默认都返回最新资源表示。

例如：

- `POST` 创建成功后返回新资源
- `PATCH` 更新成功后返回更新后资源

只有少数异步动作允许返回 `202 Accepted`。

`public` 调用面第一版不返回 `202`。如果匿名请求内部需要异步副作用，例如发邮件，它仍应对外表现为同步受理结果，而不是暴露匿名任务查询资源。

## 异步受理与任务资源

### `202 Accepted`

异步写操作统一返回标准受理对象。

至少包含：

- `jobId`
- `status`
- `acceptedAt`
- `statusUrl`

### 任务资源

异步任务状态查询资源按调用面分开：

- `/api/v1/user/jobs/{jobId}`
- `/api/v1/admin/jobs/{jobId}`
- `/api/v1/node/jobs/{jobId}`

第一版不使用全局混合任务路径 `/api/v1/jobs/{jobId}`，也不定义 `public` 面任务资源。

### 任务表示

统一任务资源应至少表达：

- `id`
- `status`
- `acceptedAt`
- `startedAt`
- `finishedAt`
- `result` 或结果引用
- `error` 或错误摘要

## 查询参数规范

全局统一的查询参数命名至少包括：

- `page`
- `pageSize`
- `q`
- `sort`
- `order`

其中：

- `q` 代表通用搜索关键词
- `sort` 代表排序字段
- `order` 代表排序方向

领域特有筛选项可以在此基础上增加，但不应重新发明：

- `keyword`
- `query`
- `page_size`

等近义字段。

## 认证资源化

### 会话

认证链默认资源化表达：

- 登录：创建 session
- 登出：删除 session
- 查看会话：读取 session 资源

推荐方向：

- `POST /api/v1/public/sessions`
- `GET /api/v1/user/me/sessions`
- `DELETE /api/v1/user/me/sessions/{sessionId}`

### 验证与重置

邮箱验证、找回密码等辅助认证流程同样资源化。

推荐方向：

- `POST /api/v1/public/verification-tokens`
- `POST /api/v1/public/password-reset-requests`
- `POST /api/v1/public/password-resets`

不使用：

- `/send-code`
- `/reset-password`

之类动作路径。

## `operationId` 与 tags

### `operationId`

所有 operation 强制显式定义 `operationId`。

命名模板收紧为：

`{surface}{domain}{resource}{action}`

其中：

- `surface` 固定为 `public | user | admin | node`
- `domain` 固定为 `Auth | Access | Catalog | Billing | Subscription | Node | System`
- `resource` 使用单数 PascalCase 资源名，必要时可以是父子资源拼接
- `action` 使用受控动词集：`List / Get / Create / Update / Delete / Cancel / Refresh / Retry / Approve / Reject`

机械规则如下：

- 集合读取：`{surface}{domain}{resource}List`
- 详情读取：`{surface}{domain}{resource}Get`
- 创建：`{surface}{domain}{resource}Create`
- 更新：`{surface}{domain}{resource}Update`
- 删除：`{surface}{domain}{resource}Delete`
- 嵌套资源：`{surface}{domain}{parentResource}{childResource}{action}`

例如：

- `publicAuthSessionCreate`
- `userSubscriptionList`
- `adminNodeUpdate`
- `nodeUsageReportCreate`
- `adminUserSubscriptionList`

所有 `operationId` 必须全局唯一，不允许不同调用面依赖同名方法再由客户端二次分组。

不允许依赖工具按 path 自动推导。

### tags

tags 按业务域分组，而不是按调用面分组。

推荐：

- `Auth`
- `Access`
- `Catalog`
- `Billing`
- `Subscription`
- `Node`
- `System`

调用面已经体现在路径中，不再让 tag 重复承担这层职责。

## 组件复用与 schema 复杂度

### 公共组件

以下高复用对象应作为公共组件抽出：

- `Problem`
- `ValidationProblem`
- `PaginationMeta`
- `Money`
- `AsyncJobAccepted`
- `PageParam`
- `PageSizeParam`
- `SortParam`
- `OrderParam`
- `IdempotencyKeyHeader`
- `sessionAuth` 与 `nodeAuth` 两套 `securitySchemes`

成功响应 envelope 是一条契约规则，不是一个必须在所有 operation 上复用的单一泛型 schema。

### 组合类型

`oneOf / anyOf / allOf` 默认严格限制。

原则：

- 默认优先使用平铺 schema
- 只有少数必要场景允许组合，例如判别联合或错误扩展
- 不为了追求 DRY 而把 schema 结构做得难读

## 状态码最小规范

第一版至少冻结以下状态码语义：

- `200`：读取成功 / 一般成功
- `201`：创建成功
- `202`：异步受理
- `400`：请求格式错误
- `401`：未认证
- `403`：已认证但无权限
- `404`：资源不存在
- `409`：冲突、重复、状态冲突或幂等冲突
- `422`：业务校验失败
- `429`：限流
- `500`：内部错误

状态码不允许在各领域中各自重新解释。

第一版不使用 `204 No Content`，以避免和统一成功响应体规则冲突。原本可能使用 `204` 的接口，统一返回 `200` 与显式成功响应体。

## 示例要求

关键接口必须提供 request / response example。

至少覆盖：

- `public` 面认证接口
- `user` 面订阅与会话接口
- `billing` 的订单 / 支付 / 退款接口
- `admin` 的套餐 / 节点 / 用户接口
- `node` 的 usage report / heartbeat / registration 接口
- `Problem` 与字段级校验错误示例

普通接口可以不为每一个都强制提供示例，但关键主链必须有。

## 向后兼容规则

### `v1` 内允许的变更

`v1` 内只允许向后兼容演进，例如：

- 新增 endpoint
- 新增可选字段
- 新增可选 query 参数
- 新增不破坏旧客户端的响应字段
- 在明确定义兼容策略下新增 enum 值
- 在不改变语义的前提下为文档补充 example 或 `x-` 扩展字段

### `v1` 内禁止的变更

以下变更视为 breaking change，必须进入 `v2`：

- 删除字段
- 重命名字段
- 修改字段类型
- 修改响应 envelope 外层结构
- 修改已发布 path 语义
- 修改 `operationId`
- 修改已有状态码语义
- 在 `v1` 内把 OpenAPI 真相源从 `3.1` 切到不被现有工具链稳定消费的版本

## 允许的简化

这份规范允许第一版做 4 类受控简化：

1. 不在第一版定义所有业务域的完整接口清单
2. 不在第一版为每个普通 CRUD 接口都强制写 example
3. 不在第一版为每个资源单独设计复杂过滤 DSL
4. 不在第一版引入 cursor 分页或 GraphQL 风格查询能力

## 明确禁止事项

以下做法在第一版 HTTP 规范中明确禁止：

1. 为不同调用面分别维护彼此独立的 OpenAPI 真相源
2. 让不同接口各自发明不同的成功 envelope
3. 让错误响应脱离 Problem Details 体系
4. 在 `v1` 中随意修改 `operationId` 或字段类型
5. 让 `node` 面脱离 JSON over HTTP 与统一路径规则
6. 广泛使用 RPC 风格路径，例如 `/create-*`、`/get-*`、`/report-*`
7. 在 JSON 中混用 `camelCase` 与 `snake_case`
8. 把写接口是否支持幂等留给实现阶段自由判断
9. 让认证规则只靠 prose 描述而没有结构化 `security` 与权限扩展字段
10. 在 `v1` 中返回 `204 No Content` 破坏统一成功响应体规则

## 决策结果

`server-v2` 第一版 HTTP 规范采用：

- 单一 OpenAPI 3.1 真相源
- Redocly 文档输出
- `@hey-api/openapi-ts` 客户端生成
- `paths` 按调用面拆分，`components` 按业务域拆分
- `/api/v1/public|user|admin|node/...` 四个调用面
- 资源化路径、复数 path、单数 schema、显式 `operationId`
- 带响应体的成功响应统一 `data + meta`
- 错误统一 Problem Details + 字段级错误数组
- 共享 `sessionAuth` 的 `user/admin` 会话认证 + 签名式 `nodeAuth`
- 统一 `Idempotency-Key`
- 统一 `202` 受理对象与按调用面分开的任务资源，`public` 面不暴露任务资源
- `v1` 明确的向后兼容边界

一句话总结：

**这份规范把 OpenAPI 从“接口文档”提升为 `server-v2` 的 HTTP 契约真相源，使文档、客户端生成、实现和 AI 修改共享同一套可预测规则。**
