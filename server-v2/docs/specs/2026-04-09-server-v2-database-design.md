# Server-V2 数据库设计

## 背景

`server-v2/` 是一个从零建立的独立 Go 服务工程。  
它的数据库设计不继承旧系统的数据模型，也不以兼容历史表结构、历史命名或历史迁移链为目标。

这份规范建立在现有目录规范之上，默认遵循 [2026-04-09-server-v2-directory-design.md](/Users/admin/Codes/ProxyCode/perfect-panel/server-v2/docs/specs/2026-04-09-server-v2-directory-design.md) 中已经确认的领域边界：

- `auth`
- `access`
- `catalog`
- `billing`
- `subscription`
- `node`
- `system`

数据库设计的优先目标不是“做成一个理论上最通用的商品引擎”，而是为 `server-v2` 的商业主链提供：

1. 清晰的主体模型
2. 稳定的交易与订阅履约链
3. 可重算的节点授权链
4. 可审计、可追溯、可重建的真相层

## 目标

1. 定义 `server-v2` 第一版数据库的核心实体、主关系和领域归属。
2. 为 `auth / access / catalog / billing / subscription / node / system` 七个域提供稳定的数据边界。
3. 明确交易、订阅、权益、节点授权、usage、审计和缓存同步的真相来源。
4. 为后续 schema bootstrap、实现计划和代码落地提供统一的数据前提。

## 非目标

1. 本文不展开完整 DDL、列级索引清单或供应商语法细节。
2. 本文不定义完整 HTTP API、节点协议出参格式或前端展示结构。
3. 本文不设计多币种、复杂促销引擎或通用电商商品规格系统。
4. 本文不把 `SQLite` 纳入第一版正式数据库目标。

## 数据库目标

第一版数据库以 **`PostgreSQL`** 作为正式真源数据库。  
`MySQL 8+` 被视为后续兼容目标，但不反向约束第一版主模型。  
`SQLite` 不纳入第一版主设计。

这意味着：

- 第一版 schema 以 PostgreSQL 语义为主
- 设计时避免无必要地依赖极少数 PostgreSQL 独占能力
- 未来如果增加 `MySQL 8+` 兼容，应通过 bootstrap / revision / 契约测试收口，而不是回头推翻领域模型

## 核心判断

这份数据库设计采用：

**主体统一 + 交易与履约分离 + 权益快照 + 节点授权预计算 + 真相层与派生层分离**

它对应 8 条工作原则：

1. `users` 是唯一主体表，认证和授权附着其上，不再拆分前后台主体表。
2. 商品定义、交易事实、订阅履约、节点消费分别建模，不互相混用。
3. 历史交易一律以订单项快照为准，不回查当前商品价格。
4. 权益必须展开为稳定快照，节点和客户端不直接回头临时拼装订阅规则。
5. 节点授权采用“默认权益 + 覆盖规则 + 最终授权结果”的三层模型。
6. usage 采用“原始事实 + 计费规则快照 + 最终计费结果”的三段式模型。
7. 核心业务状态在事务内提交，邮件、缓存刷新、输出重建等副作用通过 `outbox_events` 异步同步。
8. 业务历史查询依赖数据库事件表和操作日志，不依赖应用运行日志作为真相来源。

## 核心关系图

```text
users
-> user_identities
-> user_sessions

users
-> user_roles
-> roles
-> role_permissions
-> permissions

plans
-> plan_variants
-> plan_addons

orders
-> order_items
-> payments
-> payment_events
-> refunds
-> refund_items

subscriptions
-> subscription_periods
-> subscription_addons
-> subscription_addon_periods
-> entitlements

hosts
-> host_protocols
-> nodes
-> node_groups
-> node_assignments

node_assignment_overrides
node_usage_reports
online_sessions

system_settings
verification_tokens
outbox_events
admin_operation_logs
subscription_output_snapshots
```

这张图的核心含义是：

- `catalog` 定义卖什么
- `billing` 记录买了什么、付了什么、退了什么
- `subscription` 负责把交易履约成可被系统消费的权益
- `node` 负责宿主机、协议服务、用户节点和最终节点授权结果
- `system` 负责配置和平台性支撑数据

## 主体与认证模型

### 主体模型

第一版采用：

- `users`
- `user_identities`
- `user_sessions`

其中：

- `users` 是唯一主体表
- `user_identities` 承接登录身份
- `user_sessions` 承接会话

### 身份模型

第一版只支持邮箱身份，但不把邮箱直接绑定成唯一主体属性。

推荐结构：

- `users`
  - `id`
  - `status`
  - `archived_at`
  - 基础档案字段

- `user_identities`
  - `id`
  - `user_id`
  - `provider`
  - `identifier`
  - `secret_hash` 或对应认证材料
  - `verified_at`
  - `status`

唯一性原则：

- `user_identities(provider, identifier)` 唯一
- 不依赖 `users.email` 作为唯一真源

### 会话模型

采用独立 `user_sessions` 表，而不是把 token 直接塞进 `users`。

推荐职责：

- 记录登录后的一次独立会话
- 支持多会话并存
- 支持撤销单个会话
- 记录来源 IP、最后活跃时间、过期时间

### 验证令牌模型

验证码、找回密码、确认操作不落在 `users` 上，统一进入：

- `verification_tokens`

它至少应支持：

- 归属用户或目标标识
- `purpose`
- token / code 哈希
- `expires_at`
- `used_at`
- `status`

## 权限模型

第一版采用标准 RBAC（Role-Based Access Control）：

- `roles`
- `permissions`
- `role_permissions`
- `user_roles`

职责划分：

- `auth` 负责“是谁”
- `access` 负责“能做什么”

第一版不引入资源级授权，不做更重的 ACL 或策略引擎。

## 商品目录模型

第一版不做通用商品规格引擎，而采用更适合当前业务的轻量目录模型：

- `plans`
- `plan_variants`
- `plan_addons`

### `plans`

套餐模板。  
表示“这类套餐是什么”，而不是最终成交配置本身。

### `plan_variants`

可售变体。  
一条变体就是一个真正可买的组合，直接承载：

- 套餐周期
- 流量额度定义
- 在线连接数上限定义
- 该变体当前价格

第一版明确选择“可售变体”，而不是动态规格维度引擎。

### `plan_addons`

附加权益定义。  
它表示“哪些额外能力可以被售卖或附加到订阅上”。

第一版 addon 不是套餐表上的零散附加字段，而是独立可售定义。

## 交易模型

第一版交易链采用：

- `orders`
- `order_items`
- `payments`
- `payment_events`
- `refunds`
- `refund_items`

### 订单模型

订单采用“主表 + 明细项”模型，而不是“一张订单只对应一个动作”。

这意味着一张订单可以包含：

- 一个主套餐变体购买项
- 零到多个 addon 购买项
- 或只包含 addon
- 或只包含续费项

### 订单项模型

`order_items` 统一承载不同购买项，并通过 `item_type` 区分：

- `plan_purchase`
- `subscription_renewal`
- `addon_purchase`

不为不同 item 类型拆多张明细表。

### 价格快照原则

这是第一版数据库的硬规则：

**历史交易一律以订单项快照为准，不回查当前商品价格。**

每个 `order_item` 至少应保存：

- 类型快照
- 名称快照
- 单价快照
- 数量
- 小计
- 折扣分摊结果
- 最终应付额
- `currency`

### 优惠券模型

第一版优惠券作用于整单，但数据库必须保留订单项级别的折扣分摊结果。

这意味着：

- 优惠券使用记录可以挂在订单层
- 订单项仍需保存本项实际折扣快照

这样后续退款才有稳定依据。

### 支付模型

支付链选择：

- `orders`
- `payments`
- `payment_events`

不在第一版引入更重的 `payment_intents`。

职责：

- `orders` 代表业务交易
- `payments` 代表支付记录
- `payment_events` 保存第三方回调和支付事件轨迹

### 退款模型

退款采用：

- `refunds`
- `refund_items`

`refund_items` 不只引用订单项，还要保存自己的退款快照，例如：

- 本次退款数量
- 本次退款金额
- 本次退款对应的折扣份额
- 退款原因或类型

## 订阅与履约模型

第一版订阅链采用：

- `subscriptions`
- `subscription_periods`
- `subscription_addons`
- `subscription_addon_periods`
- `entitlements`

### `subscriptions`

表示一个用户拥有的一条订阅主记录。  
一个用户可以同时拥有多条订阅。

### `subscription_periods`

订阅的每个生效周期。

这层必须存在，因为系统需要清楚表达：

- 首购周期
- 续费周期
- 到期历史
- 周期快照

流量额度、在线连接数、周期等限制不会只留在 `plan_variants` 定义里，而要在 `subscription_periods` 中保留当期配置快照。

### `subscription_addons`

表示某条订阅实际购买了哪些 addon。

### `subscription_addon_periods`

表示 addon 的独立有效期。

这样可同时支持：

- 新购套餐时加购 addon
- 现有订阅上后加购 addon
- addon 与主订阅不同步过期

### `entitlements`

`entitlements` 是展开后的稳定权益快照，不是临时计算结果。

它的来源可以包括：

- 主订阅周期
- addon 周期

它至少应表达：

- 来源
- 生效期
- 状态
- 最终有效限制与权益快照

系统和节点消费的是真正展开后的 `entitlements`，而不是直接回头读取套餐定义。

## 节点与授权模型

第一版节点主模型采用：

- `hosts`
- `host_protocols`
- `nodes`
- `node_groups`
- `node_assignment_overrides`
- `node_assignments`

### `hosts`

宿主机。  
承载运维属性，例如：

- SSH 能力
- 主机状态
- 资源信息
- 承载关系

### `host_protocols`

宿主机上实际运行的协议服务 / 协议实例。  
它是技术承载对象，不直接等于用户看到的节点。

### `nodes`

真正下发给用户、带展示名称的用户节点。

例如：

- `香港 02 - CN2 专线`
- `香港 02 - IEPL 专线`

它们是两个不同节点，不是“一个节点的多个端点别名”。

第一版明确采用：

- 一个宿主机可以承载多个节点
- 一个用户节点绑定一个协议服务

倍率和流量计算方式挂在 `nodes`，而不是挂在宿主机或协议服务层。

### `node_groups`

节点组挂在“用户节点”这一层，而不是宿主机或协议服务层。

授权链采用：

`entitlements -> node_groups -> nodes -> node_assignments`

### 覆盖规则

为支持“给某个客户额外加指定节点”或“屏蔽某些节点”，第一版采用统一覆盖规则表：

- `node_assignment_overrides`

它应支持：

- allow group
- allow node
- deny group
- deny node

### 最终授权结果

节点侧和客户端不直接读默认权益或覆盖规则，而是读最终授权结果：

- `node_assignments`

这是授权真相层的一部分。

## usage 与在线连接模型

### usage 三段式

usage 采用：

1. 原始流量事实
2. 计费规则快照
3. 最终计费结果

推荐原始事实表：

- `node_usage_reports`

它至少应保留：

- `raw_upload_bytes`
- `raw_download_bytes`
- 节点
- 归属对象
- 上报时间

计费侧应保留：

- `billing_mode`
  - `max_of_up_down`
  - `sum_of_up_down`
  - `upload_only`
  - `download_only`
- `multiplier`
- `billable_base_bytes`
- `billable_bytes`

边界原则：

- `node` 域负责原始上报事实
- `subscription` 域负责最终归集与消耗判定

### 在线连接模型

第一版不做重设备指纹模型，而采用更轻的在线连接记录：

- `online_sessions`

它主要服务于：

- 记录在线 IP
- 统计当前活跃连接数
- 执行在线连接数限制
- 清理超时连接

在线连接数限制不靠单一计数字段，而靠活跃在线记录计算。

## 授权真相层与输出快照层

为避免“套餐改了节点但客户端没刷新”这类问题，数据库规范明确区分两层：

1. **授权真相层**
   - `entitlements`
   - `node_assignment_overrides`
   - `node_assignments`

2. **输出快照层**
   - `subscription_output_snapshots`
   - 或等价的面向客户端订阅输出快照对象

原则：

- 真相层变更后，输出快照必须可重建
- 缓存只做加速，不承担真相职责
- 客户端看到的订阅输出属于派生层，不是规则真相本身

## 系统配置模型

第一版配置不拆多张配置表，而采用一张：

- `system_settings`

但它不是纯前缀字符串垃圾桶，而采用：

- `scope`
- `key`

双维度分域。

示例：

- `scope = site`
- `scope = auth`
- `scope = billing`
- `scope = node`
- `scope = mail`

建议字段方向：

- `scope`
- `key`
- `value_type`
- `value_text`
- `value_int`
- `value_bool`
- `value_json`
- `updated_by_user_id`
- `updated_at`

## 审计与历史模型

### 业务历史查询原则

业务历史应查询数据库表，而不是依赖应用运行日志。

需要区分：

1. **当前状态**
   - 查主表
2. **状态历史**
   - 查事件表
3. **后台人工操作**
   - 查操作日志

### 推荐事件与日志对象

第一版建议至少具备：

- `payment_events`
- `subscription_events`
- `auth_events`
- `admin_operation_logs`

其中 `admin_operation_logs` 应至少保留：

- `user_id`
- `session_id`
- actor context
- 目标对象
- 动作类型
- 来源 IP
- 时间
- 简要变更说明

## 生命周期与删除策略

核心业务表不默认采用软删除。

第一版采用：

- `status`
- `archived_at`
- `ended_at`

等生命周期字段表达状态，而不是一律加 `deleted_at`。

适用对象包括但不限于：

- 套餐
- 节点
- 订阅
- 订单

只有极少数明确需要逻辑删除的辅助对象，才考虑单独使用 `deleted_at`。

## 主键、编号与货币策略

### 主键策略

核心主表使用 UUID 作为主键。

### 业务编号

第一版明确保留以下业务编号：

- `order_no`
- `payment_no`
- `subscription_no`

UUID 是数据库主键，业务编号用于对外识别和运营排查，两者职责分离。

### 货币策略

第一版只支持单币种，但金额相关表仍保留 `currency` 字段。

金额统一以最小货币单位保存，例如分（cents）。

## 初始化与演进策略

第一版采用：

- baseline schema
- revisions

而不是纯历史 migration 链，也不依赖 ORM auto-migrate 作为正式 schema 真源。

### bootstrap

负责建立当前基线 schema。

### revisions

负责基线之后的小步演进。

### seed 分层

seed 分成两类：

1. `required`
   - 系统必需初始数据
   - 例如默认角色、基础权限、必要配置初值

2. `demo`
   - 开发 / 演示数据
   - 例如样例套餐、样例节点、样例用户

真实业务数据不属于 seed。

## 事务与异步同步原则

第一版采用：

**核心状态同事务，副作用异步化**

也就是：

- 订单主状态
- 支付状态
- 订阅开通
- 周期创建
- 权益展开
- 节点授权结果更新

这些核心数据库状态，尽量在同一事务内完成。

而以下副作用通过异步机制完成：

- 邮件发送
- 缓存失效
- 订阅输出重建
- 次级任务投递

### 推荐支撑对象

- `outbox_events`
- `cache_invalidation_jobs`
- `assignment_rebuild_jobs`

规范意图不是要求它们都直接表现为同名表，而是要求：

- 真相层变更与异步同步之间必须存在可靠桥接
- 缓存刷新和输出重建不能依赖“代码顺手做一下”

## 允许的简化

这份数据库规范允许第一版做 3 类受控简化：

1. 不在第一版支持多币种
2. 不在第一版支持资源级权限模型
3. 不在第一版引入通用商品规格引擎或完整事件溯源

## 明确禁止事项

以下做法在第一版数据库设计中明确禁止：

1. 把邮箱唯一性直接绑定成 `users.email` 真相，而绕开 `user_identities`
2. 把历史订单金额建立在“回头查当前商品价格”之上
3. 让节点、客户端直接读取套餐规则而不经过 `entitlements` 与 `node_assignments`
4. 把缓存当作真相层
5. 把宿主机、协议服务、用户节点混成一张万能节点表
6. 用一张无结构 `settings(key,value)` 大表承接所有配置
7. 把 usage 简化成不可追溯的单一累计字段

## 决策结果

`server-v2` 第一版数据库采用：

- `PostgreSQL` 作为正式真源数据库
- `users + user_identities + user_sessions + RBAC`
- `plans + plan_variants + plan_addons`
- `orders + order_items + payments + payment_events + refunds + refund_items`
- `subscriptions + subscription_periods + subscription_addons + subscription_addon_periods + entitlements`
- `hosts + host_protocols + nodes + node_groups + node_assignment_overrides + node_assignments`
- `node_usage_reports + online_sessions`
- `system_settings + verification_tokens + outbox_events + admin_operation_logs + subscription_output_snapshots`

一句话总结：

**交易产生订阅，订阅展开权益，权益生成节点授权，真相层通过事件驱动输出层与缓存层同步。**
