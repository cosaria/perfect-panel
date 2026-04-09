# Server 全新数据库重构设计

## 背景

当前 `server/` 虽然已经完成目录重构，但数据库层本质上仍是旧项目的延续：

- `server/internal/platform/persistence/` 是旧 `models/` 的路径迁移版本，不是为新项目重新建模的结果。
- `server/internal/platform/persistence/migrate/` 仍然依赖历史 SQL migration 文件。
- 启动链、运行时依赖和业务服务默认假设数据库仍然遵循旧表结构和旧状态语义。

这次需求不是“继续演进旧库”，而是把数据库当作一个**全新项目（greenfield project）**来重做。旧数据全部作废，不迁移；旧表结构、旧字段名、旧 SQL migration 历史都不再构成约束；旧 API 和旧业务语义也不要求兼容。

同时，产品主线不变：它仍然是一个代理面板，覆盖用户、鉴权、套餐、订阅、节点、订单、支付、工单、公告、文档和后台运维这些全量业务域。

## 已确认的决策边界

1. 这是一次全新数据库重构，不做旧数据迁移。
2. 旧数据库中的历史数据全部作废，按空库启动。
3. 不保留旧 API、旧表结构、旧字段语义兼容。
4. 产品方向不变，仍然是代理面板，而不是改做别的产品。
5. `V1` 覆盖接近当前产品的全量业务域，而不是只做最小闭环。

## 目标

1. 用新的领域模型重新定义数据库，而不是继续沿用旧 `models` 思路。
2. 消除对历史 SQL migration 文件的依赖。
3. 建立明确的数据库 bootstrap、seed、reset 边界，让它更像一个新项目。
4. 把身份、售卖目录、计费、订阅、节点运行态、支持内容、系统运维拆成清晰领域。
5. 为后续新的 API 层和新的服务层重写提供稳定的数据边界。

## 非目标

1. 不迁移任何旧数据库记录。
2. 不保留旧表名、旧字段名、旧索引命名和旧 migration 版本号。
3. 不要求现有 `internal/domains/**` 服务实现继续直接复用。
4. 不在本次设计中承诺兼容现有前端生成的 OpenAPI 合同。
5. 不把新数据库设计成事件溯源（event sourcing）或多租户系统。

## 设计原则

1. 新数据库优先服务新的业务边界，而不是服务当前代码的迁移便利。
2. 每个业务域只保留自己的核心聚合根和必要外键，不做跨域“万能表”。
3. 当前状态和历史事件分离，不能只依赖主表最终值进行审计。
4. 静态配置与运行态指标分离，尤其是节点和订阅。
5. 所有金额、流量、时长都使用强约束字段，不把展示格式写进数据库。
6. 数据库初始化和服务启动分离，服务不自动建库。
7. 不再保留历史 SQL migration 工作流；新空库通过 schema bootstrap 建立。

## 领域地图

### `identity`

负责身份与认证。

- 终端用户
- 后台管理员
- 登录身份
- 会话
- 设备
- 验证码与短期凭证
- 安全事件

关键决策：

- `admin_users` 与 `end_users` 分离，不再共用一张用户表。
- 认证标识统一抽到独立表，不把邮箱、手机号、OAuth、Telegram 等直接塞回主体表。

### `catalog`

负责“卖什么”。

- 套餐
- 价格和计费周期
- 套餐能力项
- 节点组
- 套餐与节点组关系
- 促销活动定义

关键决策：

- `catalog` 只定义可售卖能力，不直接保存某个用户已经购买后的结果。

### `billing`

负责“钱”和交易事实。

- 订单
- 订单明细
- 支付单
- 支付尝试与网关回执
- 退款
- 优惠券
- 优惠券核销
- 账本

关键决策：

- 支付状态与订单状态分离。
- 订单只描述交易流程，不直接承担订阅可用性。

### `subscription`

负责“买完后用户真正拥有什么”。

- 订阅主体
- 订阅周期
- 订阅快照版本
- 订阅 token
- 订阅事件
- 配额快照

关键决策：

- 订阅状态机与订单状态机分离。
- 套餐配置和用户获得的订阅快照分离，避免套餐后改影响历史订阅解释。

### `node_runtime`

负责节点世界的静态配置与运行态。

- 节点
- 节点协议能力
- 节点组成员关系
- 节点心跳
- 节点指标
- 节点流量上报
- 订阅与节点分配关系

关键决策：

- 节点静态配置和节点实时状态分表。
- 用户流量使用记录来自节点上报，不直接混在订单或套餐表里。

### `support_content`

负责内容和支持。

- 工单
- 工单消息
- 公告
- 知识库文档
- 消息模板

### `system_ops`

负责平台运行。

- 站点设置
- 注册策略
- 计费策略
- 节点策略
- 审计日志
- 管理员操作日志
- 后台任务元数据
- outbox 事件

## 第一版核心表清单

### `identity`

- `end_users`
- `admin_users`
- `auth_identities`
- `user_sessions`
- `user_devices`
- `verification_tokens`
- `security_events`

### `catalog`

- `plans`
- `plan_prices`
- `plan_features`
- `node_groups`
- `plan_node_group_rules`
- `promotions`

### `billing`

- `orders`
- `order_items`
- `payments`
- `payment_attempts`
- `refunds`
- `coupons`
- `coupon_redemptions`
- `billing_ledger`

### `subscription`

- `subscriptions`
- `subscription_periods`
- `subscription_revisions`
- `subscription_tokens`
- `subscription_events`
- `quota_snapshots`

### `node_runtime`

- `nodes`
- `node_protocol_profiles`
- `node_group_members`
- `node_heartbeats`
- `node_metrics`
- `node_usage_reports`
- `subscription_node_assignments`

### `support_content`

- `tickets`
- `ticket_messages`
- `announcements`
- `knowledge_articles`
- `message_templates`

### `system_ops`

- `site_settings`
- `registration_policies`
- `billing_policies`
- `node_policies`
- `audit_logs`
- `admin_operation_logs`
- `background_jobs`
- `outbox_events`

## 关键状态机

### 订单

`draft -> pending_payment -> paid -> fulfilled -> closed`

允许：

- `pending_payment -> expired`
- `paid -> refunded`

### 支付

`initiated -> pending_gateway -> succeeded / failed / canceled`

如存在回单对账，再进入：

- `succeeded -> reconciled`

### 订阅

`pending_activation -> active -> grace_period -> expired -> canceled`

允许：

- `active -> suspended`
- `active -> exhausted`

### 节点健康

`provisioning -> active -> degraded -> offline -> retired`

### 工单

`open -> in_progress -> waiting_user -> resolved -> closed`

## 数据库约束策略

### 主键与公共字段

1. 所有核心表统一使用 `id` 主键，主键类型为 `UUID` 或 `ULID`。
2. 所有核心表默认具有：
   - `created_at`
   - `updated_at`
3. 核心交易表和订阅表默认不用软删除。
4. 需要记录操作者时，统一使用：
   - `created_by_admin_id`
   - `updated_by_admin_id`
   - 或 `actor_type + actor_id`

### 强类型字段

1. 金额统一使用最小货币单位整数，如 `amount_minor`。
2. 币种独立字段保存，如 `currency`。
3. 流量统一使用整数，不保存展示字符串。
4. 状态和类型字段使用显式枚举值，不使用魔法数字。

### 唯一键

1. `admin_users.email` 唯一。
2. `auth_identities(provider, identifier)` 唯一。
3. `plans.code` 唯一。
4. `node_groups.code` 唯一。
5. `nodes.code` 或 `nodes.slug` 唯一。
6. `orders.order_no` 唯一。
7. `payments.payment_no` 唯一。
8. `subscriptions.subscription_no` 唯一。
9. `subscription_tokens.token_hash` 唯一。
10. `coupons.code` 唯一。
11. `verification_tokens(token_hash, purpose)` 唯一。

### 索引

至少需要以下索引模式：

1. 生命周期索引：
   - `orders(status, created_at)`
   - `payments(status, created_at)`
   - `subscriptions(status, current_period_end_at)`
2. 归属索引：
   - `orders(user_id)`
   - `subscriptions(user_id)`
   - `tickets(user_id)`
   - `user_devices(user_id)`
3. 节点时序索引：
   - `node_heartbeats(node_id, reported_at desc)`
   - `node_usage_reports(node_id, reported_at desc)`
4. 配置关系索引：
   - `plan_prices(plan_id, billing_period)`
   - `plan_node_group_rules(plan_id, node_group_id)`

### 外键

1. 核心交易与订阅表保留真实外键，防止脏数据漂移。
2. 高吞吐记录表可以在性能需要时弱化外键，但必须保留索引和应用层约束。
3. 不允许内容域或设置域跨域直连交易表。

## 审计策略

1. 所有状态机核心表必须有事件或尝试记录表。
2. 管理员后台修改进入 `admin_operation_logs`。
3. 安全相关操作进入 `security_events`。
4. 对外通知或内部异步投递统一经 `outbox_events`。
5. 不能依赖主表当前字段值回推历史。

## schema 与代码落地结构

新的数据层结构目标如下：

```text
server/internal/platform/persistence/
├── schema/
│   ├── bootstrap.go
│   ├── registry.go
│   ├── naming.go
│   ├── constraints.go
│   └── seed/
│       ├── admin_seed.go
│       └── site_seed.go
├── identity/
├── catalog/
├── billing/
├── subscription/
├── node_runtime/
├── support_content/
└── system_ops/
```

说明：

1. `schema/` 只负责建库、建表、建索引、加约束和执行种子初始化。
2. 各业务域目录只承载该域的持久化模型与相关装配。
3. 不再保留旧 `migrate/database/*.sql` 作为数据库演进主路径。
4. 不再维持旧 `models` 风格的平铺目录。

## 初始化命令设计

数据库初始化必须与服务启动分离。

建议新增命令：

1. `go run . db bootstrap --config etc/ppanel.yaml`
   - 用于新空库创建 schema、索引、约束和基础系统记录
2. `go run . db seed --config etc/ppanel.yaml`
   - 用于开发环境或演示环境注入样例数据
3. `go run . db reset --config etc/ppanel.yaml --force`
   - 仅用于本地和 CI，清空后重建数据库

明确规则：

1. `go run . run` 不自动建库。
2. 生产环境禁用 `db reset`。
3. 服务启动时如果 schema 未初始化，应直接报错并退出。

## seed 策略

### `bootstrap seed`

必须存在，且幂等。

包含：

- 初始管理员
- 站点基础设置
- 默认注册策略
- 默认计费策略
- 默认节点策略

### `sample seed`

仅用于开发和演示。

可包含：

- 示例套餐
- 示例节点组
- 示例节点
- 示例公告与知识库文档
- 示例订单与订阅

## 测试策略

至少需要以下测试层次：

1. `schema smoke tests`
   - 验证核心表、索引、唯一键、外键存在
2. `bootstrap idempotency tests`
   - 连续执行两次 `db bootstrap` 不报错，也不重复插入基础记录
3. `domain model tests`
   - 用户身份唯一性
   - 非法订单状态跳转
   - 非法订阅状态跳转
   - 节点激活前缺协议配置的拒绝逻辑
4. `seed tests`
   - `bootstrap seed` 幂等
   - `sample seed` 可选关闭
5. `integration tests`
   - 用户注册 -> 下单 -> 支付成功 -> 订阅生效
   - 节点心跳 -> 流量上报 -> 配额扣减
   - 工单创建 -> 回复 -> 关闭
   - 管理员修改套餐或节点策略后的系统行为

## 对现有代码库的影响

1. 当前 `server/internal/platform/persistence/` 将被新的领域化持久层替换，而不是继续在现有平铺模型上演进。
2. 当前 `server/internal/platform/persistence/migrate/` 和其下 SQL migration 文件将整体退场。
3. 当前启动链要从“默认依赖旧 schema”切换为“显式依赖 bootstrap 完成”。
4. 当前 `internal/domains/**` 的服务实现只能作为业务词汇参考，不能默认复用其状态机和查询逻辑。
5. 当前运行时对旧模型的平铺注入会被新的领域依赖装配替换。

## 推荐实施顺序

1. 先冻结旧数据库层，不再继续给旧持久层加新字段或新表。
2. 搭建新的 `db bootstrap / db seed / db reset` 命令骨架。
3. 先实现 `identity + system_ops`。
4. 再实现 `catalog + billing + subscription`。
5. 再实现 `node_runtime`。
6. 最后实现 `support_content`。
7. 待新 schema 稳定后，再重写新的 API 层和新的服务层。

## 接受的破坏范围

本次重构接受以下破坏：

1. 旧数据库无法继续使用。
2. 旧 migration 版本链全部失效。
3. 旧服务层和旧 API 层需要大范围重写。
4. 旧测试中依赖旧 schema 的部分需要全部替换。

本次重构不接受以下结果：

1. 新空库无法通过单一 bootstrap 流程建立。
2. `bootstrap seed` 不是幂等的。
3. 交易、订阅、节点运行态继续混在同一张万能表或同一批平铺模型里。
4. 服务启动仍然隐式建库。

## 风险

1. 这是一次真正的产品内核重构，不是“改几张表”，业务服务层和 API 层重写成本很高。
2. 因为不兼容旧 API，前端最终也需要按新契约重新接线。
3. 如果实现阶段没有强约束，很容易在“先跑起来”的压力下重新长回旧式平铺模型。
4. 如果不尽早删除旧 migration 工作流，团队容易继续沿用旧思路，导致双轨并存。

## 决策结果

采用“全新数据库 + 全新业务模型”的方案，把代理面板作为一个新的全量产品来重建数据库层：

- 不迁移旧数据
- 不兼容旧 schema
- 不保留历史 migration
- 以 `identity / catalog / billing / subscription / node_runtime / support_content / system_ops` 七个领域重建数据库
- 用显式 `db bootstrap` 工作流代替历史 SQL migration 链
