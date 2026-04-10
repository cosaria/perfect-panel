# Server-V2 运行时与工作流规范

## 背景

`server-v2/` 已经完成了目录规范、数据库规范和 HTTP 规范，但仍缺少一份稳定的运行时规范，来回答：

- 单一 `cmd/server` 入口如何承载不同运行角色
- HTTP、CLI、worker、scheduler 如何共享同一套工作流主轴
- 哪些状态必须同事务提交，哪些结果必须异步重建
- 缓存、输出快照、死信恢复和健康检查如何形成统一运行纪律

这份规范默认建立在以下文档之上：

- [2026-04-09-server-v2-directory-design.md](/Users/admin/Codes/ProxyCode/perfect-panel/server-v2/docs/specs/2026-04-09-server-v2-directory-design.md)
- [2026-04-09-server-v2-database-design.md](/Users/admin/Codes/ProxyCode/perfect-panel/server-v2/docs/specs/2026-04-09-server-v2-database-design.md)
- [2026-04-10-server-v2-http-design.md](/Users/admin/Codes/ProxyCode/perfect-panel/server-v2/docs/specs/2026-04-10-server-v2-http-design.md)

它不追求“第一版就做成一个平台化运行时框架”，而是优先服务 4 件事：

1. 建立单一、可预测的运行主轴
2. 冻结事务边界与异步重建边界
3. 保证开发、测试、生产共享同一套工作流
4. 为实现计划和实际编码提供稳定的运行时契约

## 目标

1. 定义 `server-v2` 第一版运行时的主轴、进程模型和角色边界。
2. 明确事实层、真相投影层、派生输出层之间的责任分层。
3. 固化 `outbox`、worker、scheduler、migration、seed、缓存与死信恢复的工作流规则。
4. 为后续实现计划、CI 验收和部署装配提供统一的运行时前提。

## 非目标

1. 本文不设计完整部署拓扑、容器编排或云资源拓扑。
2. 本文不替代数据库规范或 HTTP 规范中的数据与接口契约定义。
3. 本文不展开每一个任务类型的详细 payload 字段。
4. 本文不要求第一版就建设完整的可视化运维平台。

## 核心判断

这份运行时规范采用：

**单入口、多角色、真相驱动、事务内写事实、事务后重建投影**

它对应 8 条工作原则：

1. `cmd/server` 是唯一运行入口，不在第一版拆成多二进制。
2. HTTP、CLI、worker、scheduler 共享同一套 app bootstrap 和配置装配。
3. 不可替代事实同事务提交，并在同事务内写入 `outbox_events`。
4. 可重建投影、输出快照、缓存失效与通知副作用统一异步执行。
5. `scheduler` 只负责按规则发起任务，不直接写主业务状态。
6. 开发、测试、生产遵守最小差异原则，不允许各环境各玩一套。
7. 失败恢复必须通过有限重试、死信区和人工重放面完成，不能退化成“翻日志 + 手工修库”。
8. 缓存和输出快照永远是派生层，不得反向成为业务真相源。

## 运行主轴

### 总体工作流

`server-v2` 的运行主轴固定为：

1. HTTP / CLI / worker 接收请求或任务
2. 写入不可替代事实
3. 在同一事务内写入 `outbox_events`
4. 事务提交后，由 worker 异步重建投影、输出快照和缓存

这里的“不可替代事实”包括：

- 订单主状态
- 支付主状态
- 订阅与订阅周期事实
- addon 购买事实
- override 规则
- 审计事件

这里的“可重建结果”包括：

- `entitlements`
- `node_assignments`
- `subscription_output_snapshots`
- 缓存失效
- 通知类副作用

### 明确禁止的实现方式

第一版运行时明确禁止：

1. 在 handler 中一边写主状态，一边直接刷新缓存，把工作流混成一团。
2. 让输出快照或缓存反向定义业务状态。
3. 让可重建投影在没有事实层更新的情况下独立漂移。
4. 让 `outbox_events` 充当业务历史真相表。

## 进程模型与运行角色

### 单一入口

`server-v2` 第一版只保留一个入口：

- `cmd/server`

它通过运行模式承载不同角色，而不是一开始就拆成多二进制。

### 运行模式

第一版固定以下运行模式：

- `serve-api`
- `serve-worker`
- `serve-scheduler`
- `migrate`
- `seed-required`
- `seed-demo`

### 角色职责

#### `serve-api`

只负责：

- 暴露 `public / user / admin / node` 四个 HTTP 调用面
- 接收请求
- 执行事务内事实写入
- 写入 `outbox_events`

它不负责：

- 直接重建投影
- 直接刷新输出快照
- 直接执行通知类副作用

#### `serve-worker`

只负责：

- 消费显式任务
- 重建 `entitlements`
- 重建 `node_assignments`
- 重建 `subscription_output_snapshots`
- 执行缓存失效
- 执行通知类副作用

它不负责：

- 暴露 HTTP 面
- 直接承接外部匿名请求
- 自己决定新的周期性调度计划

#### `serve-scheduler`

只负责：

- 按时间或周期创建任务
- 发起受控工作流

它不负责：

- 直接修改主业务状态
- 跳过 worker 直接重算投影
- 形成第二套独立业务写入口

#### `migrate`

只负责：

- baseline schema
- revisions

它不负责：

- 自动执行 seed
- 隐式启动服务

#### `seed-required`

只负责系统必需初始数据，例如：

- 角色
- 权限
- 必需系统配置初值

它必须幂等。

#### `seed-demo`

只负责开发和演示数据，例如：

- 示例用户
- 示例套餐
- 示例节点
- 示例订阅

它不能作为生产前置步骤，也不能和 `required seed` 混为一谈。

## 事务边界与异步重建

### 事务内必须完成的内容

事务内只允许写两类对象：

1. 不可替代事实
2. `outbox_events`

事务内必须完成的典型对象包括：

- 订单主状态
- 支付主状态
- 订阅主状态
- 订阅周期事实
- addon 购买事实
- override 规则
- 审计事件
- `outbox_events`

### 事务后异步完成的内容

以下内容统一通过 worker 异步完成：

- `entitlements` 重建
- `node_assignments` 重建
- `subscription_output_snapshots` 重建
- 缓存失效
- 邮件、通知等副作用

### 不可替代事实与可重建结果

运行时必须始终维持这条判断：

- 不能靠重算恢复的，属于事实层，必须同事务提交
- 可以从事实层重新生成的，属于投影层或派生层，不应被硬塞回主事务

这意味着：

- 投影失败不会推翻事实提交成功
- 但系统必须通过任务、死信和人工恢复面保证投影最终可补齐

## `outbox`、worker 与 scheduler 分层

### `outbox_events`

`outbox_events` 的职责只有一个：

**记录事务成功后必须继续发生的后续工作。**

它不是：

- 业务真相表
- 通用审计历史表
- 任意消息总线的替代品

### worker

worker 只消费显式任务，并执行：

- 投影重建
- 输出快照重建
- 缓存失效
- 通知类副作用

worker 不得绕开事实层直接凭空制造主状态。

### scheduler

scheduler 只负责按规则创建任务，例如：

- 周期性 usage 汇总
- 过期在线会话清理
- 超时订单关闭
- 订阅或 addon 到期扫描

scheduler 不直接写主业务状态。  
如果需要产生状态变化，它必须通过标准任务链触发受控工作流。

## 配置优先级与装配

### 配置优先级

`server-v2` 的配置优先级固定为：

1. 命令行显式参数
2. 环境变量
3. 本地配置文件
4. 代码默认值

### 配置装配边界

配置读取统一进入：

- `internal/platform/config`

`cmd/server` 只负责接入配置来源，不负责解释业务配置语义。

### 配置纪律

第一版固定以下纪律：

1. 生产运行不依赖某个必须存在的本地配置文件。
2. 本地开发和 demo 可以使用配置文件提升便利性。
3. 所有运行模式共享同一套配置优先级和装配逻辑。
4. 敏感配置必须有明确的环境变量映射，不允许只藏在 demo 文件中。

## Bootstrap、migration 与 seed 启动顺序

### 标准顺序

从空库到可运行状态的标准顺序固定为：

1. `cmd/server migrate`
2. `cmd/server seed-required`
3. `cmd/server seed-demo`（仅开发 / 演示）
4. `cmd/server serve-api | serve-worker | serve-scheduler`

### 明确规则

#### `serve-*` 不隐式迁移

`serve-api`、`serve-worker`、`serve-scheduler` 不允许：

- 偷偷执行 migration
- 偷偷写入 seed

如果发现 schema 版本不满足，进程应明确失败，而不是暗中修库。

#### `required seed` 必须幂等

`seed-required` 必须可重复执行，不因重复执行而产生脏重复数据。

#### `demo seed` 必须隔离

`seed-demo` 与 `seed-required` 的职责必须完全隔离。  
`demo seed` 不能承载任何系统成立所必需的数据。

## 缓存、输出快照与重建触发

### 三层模型

运行时必须区分 3 层：

1. **事实层**
   - 订单、支付、订阅、周期、addon、override、审计事件
2. **真相投影层**
   - `entitlements`
   - `node_assignments`
3. **派生输出层**
   - `subscription_output_snapshots`
   - 各类缓存

### 固定重建顺序

一旦事实层发生变化，系统必须沿固定顺序推进：

1. 事务内写事实 + `outbox_events`
2. worker 重建 `entitlements` 或 `node_assignments`
3. worker 重建 `subscription_output_snapshots`
4. 最后做缓存失效

### 必须触发重建的变更

以下变更必须进入统一重建链：

- 订阅创建、续费、取消、到期
- addon 购买、到期
- override 规则变更
- 节点组成员变更
- 节点状态或展示属性变更
- 影响输出的系统配置变更

### 派生层纪律

缓存和输出快照只承担加速职责：

- 允许缺失
- 允许重建
- 不允许反向定义业务状态

## 失败处理、重试与死信恢复

### 有限重试

第一版所有 `outbox` 投递和 worker 任务都采用有限重试，而不是无限重试。

### 死信区

超过重试阈值的对象必须进入死信区。  
死信对象至少保留：

- 任务类型
- 关联对象
- 最近错误
- 重试次数
- 最近尝试时间

### 人工恢复面

死信区必须支持：

- 查询
- 手工重放
- 手工丢弃
- 关联对象定位

### 明确禁止事项

第一版不允许用以下方式代替死信恢复面：

- 只打日志，不留结构化失败对象
- 让开发者直接去数据库里手工猜修复顺序
- 让 worker 无限重试到队列被毒死

## 最小观测面与健康检查

### 必备观测能力

第一版至少要求：

- 结构化日志
- 健康检查端点
- 最小任务运行态

### 健康检查

健康检查至少要区分：

- 进程存活
- 数据库连通
- 队列 / worker 可用

### 最小运行态

运行时至少要能定位：

- 当前运行模式
- 当前 schema 版本
- 最近一次 `required seed` 状态
- 任务创建数、成功数、失败数、死信数

### 暴露边界

健康检查和运行态接口只暴露最小必要信息，不泄露：

- secret
- 配置正文
- 敏感业务数据

日志、健康检查和运行态接口都不能替代业务审计真相表。

## 环境差异原则

### 最小差异

开发、测试、生产共享同一套：

- `cmd/server` 命令体系
- 配置优先级
- 事务主轴
- `outbox` / worker / scheduler 工作流
- schema bootstrap / revisions 机制

### 允许差异

环境之间只允许在以下方面不同：

- 数据源连接信息
- 日志与观测输出级别
- 是否允许 `seed-demo`
- 少量开发便利能力

### 禁止差异

以下差异在第一版明确禁止：

1. 开发环境偷偷自动迁移，生产环境不这样做。
2. 测试环境绕过 `outbox`，只为了省实现成本。
3. 本地读取链绕过正式输出快照和重建工作流。
4. 不同环境使用不同事务主轴。

## 允许的简化

这份规范允许第一版做 4 类受控简化：

1. 不在第一版设计完整的可视化运维控制台。
2. 不在第一版定义所有任务类型的完整 payload schema。
3. 不在第一版把所有观测数据都做成指标平台级产品。
4. 不在第一版把 worker 和 scheduler 拆成独立二进制。

## 明确禁止事项

以下做法在第一版运行时规范中明确禁止：

1. 让缓存或输出快照反向成为业务真相源。
2. 让 `serve-*` 进程偷偷执行 migration 或 seed。
3. 让 `scheduler` 直接写主业务状态。
4. 让 `outbox_events` 充当业务历史真相表。
5. 让无限重试替代死信恢复面。
6. 让开发、测试、生产使用三套不同的运行主轴。

## 决策结果

`server-v2` 第一版运行时与工作流采用：

- 单一 `cmd/server` 入口
- `serve-api / serve-worker / serve-scheduler / migrate / seed-required / seed-demo` 多角色模式
- 不可替代事实同事务提交
- 同事务写 `outbox_events`
- worker 负责投影重建、输出快照重建、缓存失效和通知副作用
- scheduler 只按规则创建任务
- 统一配置优先级：命令行参数 > 环境变量 > 配置文件 > 默认值
- 标准启动顺序：`migrate -> seed-required -> seed-demo -> serve-*`
- 缓存与输出快照作为派生层
- 有限重试 + 死信区 + 人工重放
- 开发、测试、生产遵守最小差异原则

一句话总结：

**这份规范把 `server-v2` 的运行方式冻结成一条单入口、真相驱动、事务内写事实、事务后重建投影的标准工作流，使实现阶段不再临场发明“到底谁该做什么”。**
