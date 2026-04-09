# Server 数据库规范化重构设计

## 背景

当前 `server/` 的目录已经完成一次结构收口，但数据库层仍然保留了很多旧项目形态：

- `server/internal/platform/persistence/` 里的模型边界大多沿用了历史 `models` 拆法。
- 用户、订单、订阅、节点、支付配置之间仍然存在跨域耦合和字段职责混杂。
- 节点与套餐关系里仍有逗号字段、`FIND_IN_SET` 一类不利于约束和查询优化的结构。
- 支付方式配置、认证方式配置、运行时策略和交易事实没有被明确分层。
- 现有服务和 API 默认依赖这些旧表语义，因此一旦表结构激进重建，代价会直接传导到 API、前端和运行时。

这次重构的目标不再是“把数据库当作全新产品重新发明一遍”，而是：

**在不变更当前 API 合同和主要业务语义的前提下，尽可能把数据库规范化。**

也就是说，这次工作优先解决的是存储层混乱、约束不清、关系不规范、运行态与配置态混放这些问题，而不是同时重写 API、前端契约和整个服务层。

## 已确认的决策边界

1. 当前公开 HTTP 接口、OpenAPI 契约和主要业务语义保持不变。
2. 当前前端不需要因为这次数据库重构而重新接线。
3. 这次重构优先做数据库规范化，不做产品级 API 重设计。
4. 延续此前前提：不再保留旧 SQL migration 历史作为主工作流，可按空库初始化。
5. 允许为了适配新 schema 在服务内部增加 repository、adapter 或 read model，但不通过改 API 来转移复杂度。

## 目标

1. 在保持现有 API 合同不变的前提下，重构数据库为更清晰的规范化结构。
2. 拆开静态配置、运行态指标、交易事实、审计事件和外部回调记录。
3. 用关系表替换逗号字段和跨表现算逻辑，减少 `FIND_IN_SET` 一类历史查询模式。
4. 给现有服务层提供更稳定、可约束、可审计的底层数据边界。
5. 建立明确的 `bootstrap / seed / reset / revisions` 工作流，替代历史 migration 链。

## 非目标

1. 不修改现有 OpenAPI 契约和前端生成 client 的使用方式。
2. 不把这次工作扩展成新的 API 层或新的前端契约重写。
3. 不强行把所有业务域改造成全新的产品模型。
4. 不在本次设计里引入事件溯源（event sourcing）、多租户或跨产品平台化抽象。
5. 不要求一次性重写全部 `internal/domains/**` 逻辑。

## 兼容性原则

### API 兼容

1. `docs/openapi/*.json` 对应的路由、请求字段、响应字段和状态语义保持稳定。
2. 节点轮询、节点上报、支付回调、Telegram webhook 这类非 OpenAPI 面仍保持当前协议形态。
3. 后端可以为了适配新 schema 重写查询实现，但不能通过修改 API 来规避数据库设计问题。

### 业务语义兼容

1. 当前订单、支付、订阅、节点、工单的主状态语义保持兼容。
2. 当前用户模型仍以单一 `users` 主体为中心，`is_admin` 等现有行为边界保持可表达。
3. 当前支付回调、节点拉配置、节点上报流量、工单流转等核心流程保持兼容。

### 实现兼容

1. 如果服务层暂时无法直接读取新 schema，可以增加内部 adapter、聚合查询或 read model。
2. 兼容层是为了保护现有 API，不是为了把旧表原样续命。
3. 一旦新 schema 稳定，兼容层应逐步收敛，而不是成为新的永久负担。

## 当前主要问题

1. 用户、订阅、设备、认证方式、邀请码等语义在旧模型中交织，主体边界不清。
2. 订单同时承载交易、订阅关联、支付方式、优惠券和一部分履约语义。
3. 支付“方式配置”和“支付事实”耦合，导致网关配置与回调记录难以独立治理。
4. 节点授权解析依赖运行时跨表拼装，缺少稳定的读模型。
5. 流量、配额、节点上报、排行榜和运营统计之间缺乏明确分层。
6. 系统配置、认证配置、节点策略、计费策略分散，容易继续长成万能配置表。

## 设计原则

1. 先保 API，再规范库。
2. 优先规范关系、约束、索引和审计，不优先追求概念上的“全新世界”。
3. 用实体表和关系表表达业务关系，不再依赖字符串集合字段承载结构化关系。
4. 当前状态和历史事件分离，主表不承担全部审计职责。
5. 静态配置和运行态数据分离，尤其是支付、节点、认证、订阅。
6. 所有金额、流量、时长和状态都使用强约束字段。
7. 所有外部写入口都必须有信任模型、幂等规则和失败审计。

## 规范化后的模块边界

这次不是按“全新产品七大域”重造一套世界，而是在保持现有 API 语义的前提下，把持久层收口为下面几组更务实的模块：

```text
server/internal/platform/persistence/
├── schema/
│   ├── bootstrap.go
│   ├── registry.go
│   ├── naming.go
│   ├── constraints.go
│   ├── revisions/
│   └── seed/
├── identity/
├── catalog/
├── billing/
├── subscription/
├── node/
├── content/
└── system/
```

说明：

1. `schema/` 负责建库、建表、索引、约束、基础 seed 和后续 schema revisions。
2. `identity/` 负责用户主体、认证标识、会话、设备、验证码和安全事件。
3. `catalog/` 负责套餐、价格、能力项、节点组和套餐到节点组的关系。
4. `billing/` 负责订单、支付、退款、优惠券、账本和支付回调事实。
5. `subscription/` 负责用户实际拥有的订阅、周期、token、配额与订阅事件。
6. `node/` 负责 server、node、节点协议配置、节点运行态、节点授权解析和流量上报事实。
7. `content/` 与 `system/` 负责工单、公告、文档、模板及各类平台配置和策略。

## 表设计策略

## `identity`

目标：在不改变当前登录与权限语义的前提下，把“用户主体”和“认证方式”拆开。

建议核心表：

- `users`
- `user_auth_identities`
- `user_sessions`
- `user_devices`
- `verification_tokens`
- `auth_providers`
- `auth_provider_configs`
- `verification_policies`
- `verification_deliveries`
- `security_events`

关键决策：

1. 当前阶段不拆 `admin_users` / `end_users` 两张主体表，保留单一 `users` 主体和 `is_admin` 兼容语义。
2. 邮箱、手机号、Telegram、OAuth 标识统一收敛到 `user_auth_identities`。
3. 登录方式启停、配置和验证策略单独建模，不再继续散落在 `system` 或 JSON 配置里。

## `catalog`

目标：把“卖什么”和“用户实际拥有什么”彻底分开，但不改变现有购买语义。

建议核心表：

- `plans`
- `plan_prices`
- `plan_features`
- `node_groups`
- `node_group_nodes`
- `plan_node_group_rules`
- `promotions`

关键决策：

1. 套餐定义只承载可售能力，不直接承载用户履约结果。
2. 节点与套餐关系使用关系表，不再依赖逗号字段和 `FIND_IN_SET`。

## `billing`

目标：把“交易事实”“支付方式配置”“回调记录”拆开。

建议核心表：

- `orders`
- `order_items`
- `payments`
- `payment_callbacks`
- `payment_gateways`
- `payment_gateway_secrets`
- `refunds`
- `coupons`
- `coupon_redemptions`
- `billing_ledgers`

关键决策：

1. `payment_gateways` 负责可用支付方式和展示配置，`payments` 负责用户交易事实。
2. 支付回调原始事实进入 `payment_callbacks`，不能只依赖主交易表当前状态。
3. 订单状态和支付状态分离，但状态语义保持兼容当前 API 行为。

## `subscription`

目标：把套餐目录和用户实际订阅拆开，同时保持现有订阅查询和 token 语义。

建议核心表：

- `subscriptions`
- `subscription_periods`
- `subscription_tokens`
- `subscription_usage_snapshots`
- `subscription_events`
- `subscription_node_assignments`

关键决策：

1. 保留当前订阅 token 与用户订阅查询能力，但底层拆成更清晰的主表、周期表和事件表。
2. 订阅流量、到期时间、状态迁移和手工调整都要可审计。
3. 节点授权结果最终落到 `subscription_node_assignments`，避免运行时现算。

## `node`

目标：把 server / node 配置、节点运行态、节点授权解析和节点上报事实分开。

建议核心表：

- `servers`
- `nodes`
- `node_protocol_profiles`
- `node_status_reports`
- `node_usage_reports`
- `node_metrics_daily`

关键决策：

1. `servers` / `nodes` 保留现有概念，避免为了规范化再额外引入新的 API 词汇。
2. 节点拉取用户授权时优先读取 assignment 结果，而不是跨套餐和订阅实时拼装。
3. 节点上报流量、在线用户、状态变化进入独立事实表，不直接覆盖业务主表。

## `content` 与 `system`

目标：把内容数据和平台配置数据分离，避免继续堆成万能系统表。

建议核心表：

- `tickets`
- `ticket_messages`
- `announcements`
- `documents`
- `message_templates`
- `site_settings`
- `registration_policies`
- `billing_policies`
- `node_policies`
- `admin_operation_logs`
- `background_jobs`
- `outbox_events`

关键决策：

1. 工单消息与工单主表分离。
2. 公告、文档、模板不继续混在同一类配置记录中。
3. 系统配置按职责拆分，而不是只依赖 `key/value/category` 万能结构。

## 关键状态语义

本次重构以**保持现有主状态语义兼容**为目标，不主动引入新的 API 可见状态。

### 订单

兼容当前主流程：

`pending -> paid -> finished`

兼容当前异常分支：

- `pending -> close`
- `paid -> failed`

### 支付

保持当前网关通知驱动的成功确认语义，至少覆盖：

`pending -> succeeded / failed / canceled`

并额外记录：

- 回调接收事实
- 验签结果
- 幂等处理结果

### 订阅

保持当前主语义兼容：

`pending -> active -> finished / expired / canceled`

### 节点健康

保留当前运行态监控语义：

`active / degraded / offline / retired`

### 工单

`open -> in_progress -> waiting_user -> resolved -> closed`

## 数据库约束策略

### 主键与公共字段

1. 核心表统一使用 `id` 主键，优先 `UUID` 或 `ULID`。
2. 核心表默认具有 `created_at`、`updated_at`。
3. 核心交易表、回调事实表和订阅表默认不用软删除。
4. 需要记录操作者时，统一使用 `created_by_admin_id`、`updated_by_admin_id` 或 `actor_type + actor_id`。

### 强类型字段

1. 金额统一使用最小货币单位整数，如 `amount_minor`。
2. 币种独立字段保存，如 `currency`。
3. 流量统一使用整数，不保存展示字符串。
4. 状态和类型字段使用显式枚举值，不使用魔法数字。

### 唯一键

至少包括：

1. `user_auth_identities(provider, identifier)` 唯一。
2. `plans.code` 唯一。
3. `node_groups.code` 唯一。
4. `nodes.code` 或 `nodes.slug` 唯一。
5. `orders.order_no` 唯一。
6. `payments.payment_no` 唯一。
7. `subscription_tokens.token_hash` 唯一。
8. `payment_gateways.token` 唯一。
9. `verification_tokens(token_hash, purpose)` 唯一。

### 索引

至少需要以下索引模式：

1. 生命周期索引：
   - `orders(status, created_at)`
   - `payments(status, created_at)`
   - `subscriptions(status, expired_at)`
2. 归属索引：
   - `orders(user_id)`
   - `subscriptions(user_id)`
   - `tickets(user_id)`
   - `user_devices(user_id)`
3. 节点时序索引：
   - `node_status_reports(node_id, reported_at desc)`
   - `node_usage_reports(node_id, reported_at desc)`
4. 配置关系索引：
   - `plan_prices(plan_id, billing_period)`
   - `plan_node_group_rules(plan_id, node_group_id)`
   - `subscription_node_assignments(subscription_id, node_id)`

### 外键

1. 交易、订阅、支付回调和节点授权关系保留真实外键。
2. 高吞吐运行态记录可按性能需要弱化部分外键，但必须保留索引和应用层约束。
3. 内容域和设置域不允许直接侵入交易域主表。

## 审计与事件策略

1. 订单、支付、订阅、节点授权变更必须有事件或事实记录。
2. 管理员后台修改进入 `admin_operation_logs`。
3. 安全相关操作进入 `security_events`。
4. 对外通知和内部异步副作用统一进入 `outbox_events`。
5. 不能依赖主表当前字段值回推完整历史。

## 外部写入口信任合同

所有会把外部世界的数据写进系统的入口，统一遵循以下规则：

1. 支付回调、节点上报、机器人 webhook 等都必须有独立的认证模型。
2. 所有凭据都必须支持轮换、禁用和审计。
3. 所有签名校验失败、token 失效、secret 不匹配都必须记录失败事实。
4. 外部写入口必须定义回放窗口和重放拦截策略。
5. 不允许某些入口靠 URL token，另一些入口靠隐式配置，再来一个入口靠手写逻辑，各自为政。

## 异步幂等与重放合同

支付通知、订单激活、节点流量上报、批量配额任务统一遵循以下规则：

1. 每类异步事件都必须定义事件唯一键。
2. 重复提交允许发生，但业务效果只能生效一次。
3. 重试策略、死信策略和人工排障入口必须可定义。
4. 成功、失败、跳过、重复消费都必须可审计。
5. 不能再依赖“记了日志但返回 nil”来假装问题已处理。

## 节点授权解析合同

`node` 相关流程统一按下面的权威顺序工作：

1. 套餐目录定义节点访问规则。
2. 用户订阅在购买、续费、调整后生成或刷新节点授权结果。
3. 节点拉用户权限时优先读取 `subscription_node_assignments` 这类明确读模型。
4. 节点上报流量和状态时，优先按 assignment 或订阅主键做幂等归属。
5. 不再把“节点可见哪些用户”完全留给运行时跨表现算。

## schema 初始化与后续演进

### 初始化命令

建议新增命令：

1. `go run . db bootstrap --config etc/ppanel.yaml`
   - 创建 schema、索引、约束和基础系统记录
2. `go run . db seed --config etc/ppanel.yaml`
   - 注入开发或演示数据
3. `go run . db reset --config etc/ppanel.yaml --force`
   - 仅用于本地和 CI

明确规则：

1. `go run . run` 不自动建库。
2. 生产环境禁用 `db reset`。
3. 服务启动时如果 schema 未初始化，应直接报错并退出。

### 后续 schema 演进机制

1. 初始空库使用 `bootstrap` 建立。
2. 进入新 schema 时代后，后续结构变化通过 `schema/revisions/` 管理。
3. revision 采用 forward-only 策略，不回到旧项目那种历史 migration 链治理方式。
4. 新增字段、索引、约束、表拆分都必须通过 revision 明确记录。

## Runtime Cutover Contract

这次 spec 明确约束运行时切换边界：

1. 在新 schema 尚未完成兼容适配前，不承诺旧运行时可直接跑在新 schema 上。
2. 允许先完成 `bootstrap / seed / reset / revisions` 和 schema tests，再逐步接入现有服务。
3. 切换期间保护的是 API 和协议，不是保护旧表结构。
4. 真正完成切换的最低标准至少包括：
   - 注册 / 登录
   - 下单 / 支付回调 / 订单激活
   - 订阅查询与 token 获取
   - 节点拉配置 / 上报流量
   - 工单主流程

## 测试策略

至少需要以下测试层次：

1. `schema smoke tests`
   - 验证核心表、索引、唯一键、外键存在
2. `bootstrap idempotency tests`
   - 连续执行两次 `db bootstrap` 不报错，也不重复插入基础记录
3. `compatibility tests`
   - 现有 API request/response 形状不变
   - 现有 webhook / node 协议不变
4. `domain tests`
   - 用户认证方式唯一性
   - 非法订单状态跳转
   - 非法订阅状态跳转
   - 节点授权解析结果正确
5. `async idempotency tests`
   - 支付回调重复投递只生效一次
   - 订单激活重复消费不重复履约
   - 节点流量重复上报不重复记账
6. `integration tests`
   - 用户注册 -> 下单 -> 支付成功 -> 订阅生效
   - 节点拉权限 -> 节点上报流量 -> 配额扣减
   - 工单创建 -> 回复 -> 关闭

## 对现有代码库的影响

1. `server/internal/platform/persistence/` 需要按规范化模块重组，但目标是支撑现有 API，不是推翻现有 HTTP 面。
2. `server/internal/platform/persistence/migrate/` 和其下 SQL migration 文件退出主工作流。
3. 当前服务层可以继续沿用现有 handler / route / job 名称，但底层查询与装配会逐步改写。
4. `internal/domains/**` 的业务词汇继续保留参考价值，但底层不能继续依赖旧式平铺表结构。
5. 当前支付回调、节点协议和后台任务需要补上幂等、审计和信任合同，而不是只换表名。

## 推荐实施顺序

1. 冻结旧持久层，不再继续往旧平铺模型里加新字段和新关系。
2. 搭建 `schema/bootstrap`、`seed`、`revisions` 骨架。
3. 先实现 `identity + system` 的规范化表和兼容查询。
4. 再实现 `catalog + billing + subscription`。
5. 再实现 `node` 侧授权解析和流量事实表。
6. 最后整理 `content` 与剩余配置域。
7. 每完成一段，都以“现有 API 不变”为验收标准做回归。

## 接受的破坏范围

本次重构接受以下变化：

1. 旧 migration 历史不再继续沿用。
2. 旧表结构和旧字段组织方式可以被重构。
3. 服务内部 repository、adapter、read model 可以重写。
4. 旧测试里直接依赖旧表结构的部分需要替换。

本次重构不接受以下结果：

1. OpenAPI 契约或现有 HTTP 路由因为数据库重构而发生破坏。
2. 前端必须配合改接口才能接上新数据库。
3. 节点关系、套餐关系继续停留在逗号字段和 `FIND_IN_SET` 上。
4. 支付方式配置继续和支付交易事实混在一起。
5. 外部回调和节点上报仍然没有幂等和审计边界。
6. 服务启动仍然隐式建库。

## 风险

1. 这是一次“兼容 API 的数据库规范化”，实现难点在兼容层，不在画一个新世界。
2. 如果一味追求纯理论上的最漂亮模型，很容易反向增加服务层复杂度。
3. 如果不尽早收掉旧 migration 工作流，团队会继续双轨并存。
4. 如果规范化没有同步补上幂等、信任和审计，换了新表也只是旧问题换壳。

## 决策结果

采用“保持现有 API 不变、尽可能规范化数据库”的方案：

- 保持现有 API 和主要业务语义兼容
- 不再延续旧 SQL migration 历史作为主工作流
- 用规范化实体表和关系表替代旧式平铺模型、逗号字段和运行时现算
- 用 `bootstrap + revisions` 取代旧 migration 链
- 用兼容层保护现有 API，而不是让 API 为数据库问题背锅
