# Server-V2 目录设计

## 背景

`server-v2/` 是一个从零建立的独立 Go 服务工程。目录设计不再以历史结构为默认参考，也不以兼容旧路径、旧命名或旧分层方式为目标。

在 AI 驱动开发场景里，目录一旦同时存在多套组织规则，AI 就需要反复猜测代码落点，容易增加搜索成本、上下文切换和误改风险。
相比“分层是否传统”或“结构是否对称”，`server-v2/` 更优先追求三件事：落点可预测、职责边界清晰、搜索路径稳定。

## 目标

1. 定义 `server-v2/` 的标准目录结构，作为后续实现的唯一目录约束。
2. 降低 AI 定位代码的搜索成本，让业务代码、平台代码和装配代码都有稳定落点。
3. 明确业务域、平台基础设施、应用装配和测试目录之间的职责边界。
4. 为后续数据库设计、API 设计和 implementation plan 提供统一的目录前提。

## 非目标

1. 本文不设计仓库根目录未来的整体布局。
2. 本文不设计数据库 schema、HTTP API、节点协议或异步任务模型。
3. 本文不设计部署拓扑、运行环境、发布流程或运维规范。
4. 本文不讨论安全策略、观测体系或其他非目录层面的系统设计。

## 核心判断

`server-v2/` 默认采用 Go 服务项目常见的外骨架：`cmd/server/main.go + internal/`。内部组织不走纯技术横切，也不允许领域目录自由生长。

原因有两个：

1. 纯技术横切会把同一个需求拆散到多个顶层语义里，放大 AI 的搜索面。
2. 完全自由的领域切片会让每个领域逐渐长出不同结构，削弱路径可预测性。

因此，这份设计选择：

**稳定外骨架 + 领域优先定位 + 固定子目录模板**

它对应 4 条工作原则：

1. 顶层目录尽量少，且语义稳定。
2. 业务代码先按领域定位，再按固定子目录细分。
3. 基础设施只放进 `platform/`，不能吞业务。
4. 测试默认进入 `tests/`，只有需要访问未导出符号的包内单测才允许与实现同目录。

## 推荐目录树

```text
server-v2/
├── cmd/
│   └── server/
│       ├── main.go
│       ├── root.go
│       ├── run.go
│       ├── db_bootstrap.go
│       ├── db_seed.go
│       └── openapi.go
├── internal/
│   ├── app/
│   │   ├── bootstrap/
│   │   ├── runtime/
│   │   ├── routing/
│   │   └── wiring/
│   ├── platform/
│   │   ├── cache/
│   │   ├── config/
│   │   ├── db/
│   │   ├── http/
│   │   ├── observability/
│   │   ├── queue/
│   │   └── support/
│   └── domains/
│       ├── access/
│       │   ├── api/
│       │   ├── model/
│       │   ├── policy/
│       │   ├── store/
│       │   └── usecase/
│       ├── auth/
│       │   ├── api/
│       │   ├── model/
│       │   ├── policy/
│       │   ├── store/
│       │   └── usecase/
│       ├── billing/
│       │   ├── api/
│       │   ├── jobs/
│       │   ├── model/
│       │   ├── policy/
│       │   ├── store/
│       │   └── usecase/
│       ├── catalog/
│       │   ├── api/
│       │   ├── model/
│       │   ├── policy/
│       │   ├── store/
│       │   └── usecase/
│       ├── node/
│       │   ├── api/
│       │   ├── jobs/
│       │   ├── model/
│       │   ├── policy/
│       │   ├── store/
│       │   └── usecase/
│       ├── subscription/
│       │   ├── api/
│       │   ├── jobs/
│       │   ├── model/
│       │   ├── policy/
│       │   ├── store/
│       │   └── usecase/
│       └── system/
│           ├── api/
│           ├── model/
│           ├── policy/
│           ├── store/
│           └── usecase/
├── tests/
│   ├── app/
│   │   ├── bootstrap/
│   │   ├── routing/
│   │   └── wiring/
│   ├── domains/
│   │   ├── access/
│   │   │   ├── api/
│   │   │   ├── store/
│   │   │   └── usecase/
│   │   ├── auth/
│   │   │   ├── api/
│   │   │   ├── store/
│   │   │   └── usecase/
│   │   ├── billing/
│   │   │   ├── api/
│   │   │   ├── jobs/
│   │   │   ├── store/
│   │   │   └── usecase/
│   │   ├── catalog/
│   │   ├── node/
│   │   │   ├── api/
│   │   │   ├── jobs/
│   │   │   ├── store/
│   │   │   └── usecase/
│   │   ├── subscription/
│   │   │   ├── api/
│   │   │   ├── jobs/
│   │   │   ├── policy/
│   │   │   ├── store/
│   │   │   └── usecase/
│   │   └── system/
│   ├── platform/
│   │   ├── cache/
│   │   ├── db/
│   │   ├── http/
│   │   ├── queue/
│   │   └── support/
│   ├── contract/
│   ├── fixtures/
│   ├── integration/
│   └── smoke/
├── docs/
├── etc/
├── go.mod
└── go.sum
```

说明：

- 上面的目录树是规范骨架，不要求一次性建出所有按需目录。
- `tests/` 里的镜像目录只展示代表性示例，不表示完整枚举；完整落点以本文后面的测试规则为准。

## 顶层目录职责

### `cmd/server/`

当前唯一的可执行入口目录。  
`cmd/server/main.go` 负责启动根 CLI，`run`、`db bootstrap`、`db seed`、`openapi` 等命令也都定义在这个目录下。

这样做有两个直接好处：

1. 更贴近 Go 服务项目常见的 `cmd/<binary>/main.go` 组织方式。
2. 未来如果拆分额外二进制或微服务，可以自然扩展为新的 `cmd/<service>/` 兄弟目录，而不需要回头重排入口结构。

`cmd/server/` 只承载命令定义、子命令注册和参数绑定，不承载业务流程、数据库读写或运行时装配。

### `internal/app/`

只负责系统装配。

- `bootstrap/`：启动顺序与初始化流程
- `runtime/`：运行时容器、生命周期管理
- `routing/`：总路由树拼装与域级路由挂载
- `wiring/`：依赖注入与模块装配

`app/` 绝不能成为临时业务目录。

其中 `routing/` 只负责顶层路由树、挂载顺序、入口分面和跨域公共中间件装配。  
具体业务域的子路由、路径参数、HTTP 方法和接口落点仍然归各自 `domains/*/api/` 所有。

路由规则采用**单一路径原则**：

- `domains/*/api/` 定义本域子路由
- `internal/app/routing/` 只负责挂载这些子路由到顶层入口

`internal/app/routing/` 不直接声明业务接口路径，也不重复注册领域内的具体 handler。

### `internal/platform/`

只放跨域基础设施。

- `cache/`：缓存与 Redis 适配
- `config/`：配置读取与解析
- `db/`：数据库连接、事务包装、通用持久层基础
- `http/`：通用 HTTP 中间件、响应封装、无业务语义 transport helper
- `observability/`：日志、trace、metrics
- `queue/`：异步任务执行基础设施
- `support/`：少量真正跨域的通用帮助能力与无业务语义原语

判断标准很简单：如果一个包带有明确业务语义，就不应放在 `platform/`。

对 `support/` 还有一条额外硬约束：

- 只允许进入不带领域词汇、业务状态、业务错误码、业务 DTO 的纯通用代码
- 任何引用 `auth / billing / catalog / node / subscription / system / access` 等业务语义的代码，都不得进入 `platform/support/`

### `internal/domains/`

这是 `server-v2` 的主活动区。  
业务代码按领域放在这里，AI 首先应该猜“这是什么业务域”，而不是先猜“这是 handler、service 还是 repository”。

### `tests/`

统一承接所有测试，并镜像实现路径。

- `tests/domains/`：镜像业务域测试
- `tests/platform/`：镜像平台层测试
- `tests/app/`：镜像装配与启动测试
- `tests/integration/`：真实依赖联动
- `tests/smoke/`：关键路径冒烟
- `tests/contract/`：契约测试
- `tests/fixtures/`：共享测试数据

测试放置采用明确优先级，而不是并列自由选择：

1. 如果测试主要对应某个具体实现目录，应优先镜像到 `tests/app/`、`tests/platform/` 或 `tests/domains/`
2. 如果测试跨越多个实现目录、真实依赖或多个业务域协作，再进入 `tests/integration/`
3. 如果测试验证外部契约、生成物或稳定接口，再进入 `tests/contract/`
4. 如果测试只验证启动成功和关键主链可用，再进入 `tests/smoke/`

默认规则仍然是“测试先放 `tests/` 并镜像实现路径”，但这里有一个明确例外：

- 如果测试必须访问未导出符号、未导出构造器或包内状态，允许把该类 `_test.go` 放回实现目录，与被测包同目录存在

这意味着：

- 黑盒测试、契约测试、集成测试、绝大多数用例测试，优先进入 `tests/`
- 只有包内白盒单测，才允许留在实现目录
- 不得因为更方便组织 fixture、helper、mock，就把本可镜像到 `tests/` 的测试放回实现目录

目录树里的测试目录只展示代表性镜像示例，而不是完整列出所有镜像目录；真正的落点以这组优先级规则为准。

### `docs/`

只放 `server-v2` 自身文档资产。  
它不承载运行配置、样例数据或实现代码。

### `etc/`

只放配置样例、默认配置和本地运行所需的静态配置资产。  
它不承载文档说明，也不承载业务 seed 数据。

## 顶层领域职责

### `access`

负责角色、权限、授权关系和后台访问控制语义。  
它不负责登录、注册、验证码和会话。

### `auth`

负责用户身份、注册、登录、验证、会话和安全事件。  
它不负责 RBAC 权限模型。

### `billing`

负责订单、支付意图、支付事件、退款、交易状态流转，以及优惠券在交易中的应用结果、核销记录和折扣快照。  
它不负责套餐目录、价格目录和优惠券定义本身，也不负责订阅履约结果。

### `catalog`

负责套餐、价格、优惠券定义和面向商业售卖的目录规则。  
它不负责下单支付本身，也不负责优惠券在交易中的应用结果、核销记录或折扣快照。

### `node`

负责节点库存、协议入口、协议配置、节点读模型、节点注册和节点侧原始交互。  
它不负责订阅授权生成，也不负责权益配额扣减结果。

### `subscription`

负责订阅、权益、周期、配额、授权分配、使用量归集和最终消耗判定。  
它不负责订单支付采集，也不负责节点侧原始 usage 上报。

### `system`

负责站点级配置、各领域可运营策略配置、邮件和其他系统级设置。  
它不直接承载各业务域的运行时规则实现，也不负责具体交易或订阅生命周期。

## 跨域主责裁决规则

当一个能力同时涉及多个领域时，不允许按“谁先用到就放谁那里”的方式随意落点。必须按下面的顺序裁决主责域：

1. **定义归属**：谁拥有该对象的定义、可配置规则或售卖规则，谁拥有它的规范定义。
2. **事实归属**：谁记录外部交互、交易发生或原始输入事实，谁拥有该事实记录。
3. **生命周期归属**：谁拥有最终状态、权益结果或持续生命周期，谁拥有结果判定。

套用到当前领域上：

- 优惠券定义归 `catalog`，优惠券在订单中的应用结果和核销事实归 `billing`
- 节点侧原始 usage 上报归 `node`，配额归集与最终消耗归 `subscription`
- 用户身份与会话归 `auth`，角色与权限关系归 `access`
- 站点级可运营策略配置归 `system`，具体运行时规则实现仍归各领域自己的 `policy/`

## 领域目录模板

每个领域都遵循固定子目录模板，不允许自由发散命名。

### `api/`

放该领域自己的 handler、route 绑定、请求 DTO、响应 DTO 和 transport adapter。  
它负责具体接口的输入输出边界，不负责复杂业务流程。

### `usecase/`

放业务流程编排。  
这是领域行为的主要承载处，也是 AI 改需求时最常进入的目录。  
如果一段逻辑需要协调多个步骤、多个 `store / jobs / api` 或存在显式 I/O 编排，它应进入 `usecase/`。

### `store/`

放该领域自己的数据读写。  
它只负责持久化访问，不负责规则判断。  
具体领域 repository、查询对象、读模型映射都应归该领域 `store/` 所有。

### `model/`

放该领域核心结构、值对象、枚举和领域内聚合需要的基础表示。  
禁止出现全局横切 `models/` 顶层目录。

### `policy/`

放规则、校验、状态迁移、授权判断等逻辑。  
这类代码如果混进 `usecase/`，文件会很快变得过重。  
如果一段逻辑是纯规则、可被多个流程复用、且不依赖 I/O，它应进入 `policy/`。

### `jobs/`

只有该领域确实存在异步消费或后台任务时才出现。  
不是所有域都必须有 `jobs/`。

## `store` 与 `platform/db` 的边界

这条边界必须稳定，不允许把 repository 和事务管理来回漂移。

### 放进 `domains/*/store/` 的内容

- 具体领域 repository
- 领域表或读模型的查询与写入
- 与该领域数据结构绑定的 ORM/SQL 映射
- 该领域的持久化接口与读写适配

### 放进 `platform/db/` 的内容

- 数据库连接池
- 事务执行器与事务包装
- 通用查询辅助、分页基础、SQL helper
- 测试数据库 harness 和无业务语义的 DB 原语

### 明确规则

- `platform/db/` 不拥有任何具体业务域 repository
- 跨 store 的事务由发起该流程的 `usecase/` 通过 `platform/db/` 提供的事务执行器协调
- `store/` 不自行创建全局数据库入口，只消费 `platform/db/` 提供的基础能力

## `jobs` 与 `platform/queue` 的边界

这条边界必须稳定，不允许任务契约和队列基础设施混在一起。

### 放进 `domains/*/jobs/` 的内容

- 该领域的任务名称
- 该领域的 task payload
- enqueue 入口
- consumer handler
- 幂等、重试判定和任务级业务规则

### 放进 `platform/queue/` 的内容

- 队列客户端
- worker runtime
- 通用重试/backoff 原语
- 无业务语义的任务执行包装与观测

### 明确规则

- `platform/queue/` 不定义任何具体业务 task payload
- `domains/*/jobs/` 不重复实现通用队列客户端或 worker 基础设施
- 如果一个任务由某个领域状态变化触发，则该领域拥有 enqueue 契约

## `api` 与 `platform/http` 的边界

这条边界必须稳定，不允许实现阶段自由漂移。

### 放进 `domains/*/api/` 的内容

- 与具体业务域绑定的 handler
- 与具体业务域绑定的请求 DTO / 响应 DTO
- 该业务域自己的 route 声明与子路由注册
- 带业务语义的 HTTP 落点、输入输出映射和序列化约束

### 放进 `platform/http/` 的内容

- 全局中间件
- 通用错误映射
- 通用响应包裹
- 与业务无关的分页、绑定、上下文、HTTP helper
- 无业务语义的 HTTP 原语与框架适配层

### 明确禁止

- 不允许把某个业务域的请求 DTO 或响应 DTO 放进 `platform/http/`
- 不允许把具体路由 handler 放进 `platform/http/`
- 不允许把跨域中间件塞回 `domains/*/api/`
- 不允许把顶层路由树装配回 `domains/*/api/`
- 不允许在 `internal/app/routing/` 里重新声明领域内的具体 HTTP 落点

## AI 友好性说明

这套目录不是为了“理论最优分层”，而是为了让 AI 的路径判断尽量固定。

典型定位路径如下：

- 改登录流程：先看 `internal/domains/auth/usecase/`
- 改支付回调：先看 `internal/domains/billing/api/`，再看 `internal/domains/billing/usecase/`
- 改订阅开通规则：先看 `internal/domains/subscription/policy/`
- 改节点协议下发：先看 `internal/domains/node/usecase/` 或 `internal/domains/node/api/`
- 改权益授权分配：先看 `internal/domains/subscription/usecase/` 或 `internal/domains/subscription/policy/`
- 找对应测试：直接镜像到 `tests/domains/...`

这使 AI 的搜索变成固定两步：

1. 先猜业务域
2. 再猜固定子目录

相比全局 `services / handlers / repositories / models` 横切分层，这种方式更少出现全仓库横跳。

## 明确禁止事项

以下目录或倾向在 `server-v2/` 中明确禁止：

1. 新增全局横切业务目录：
   - `services/`
   - `repositories/`
   - `handlers/`
   - `models/`
   - `types/`
2. 在 `internal/app/` 写业务规则
3. 在 `internal/platform/` 写带业务语义的 repository 或 usecase
4. 在 `domains/*/api/` 承载复杂业务流程
5. 在 `domains/*/store/` 写规则、策略和状态机
6. 新增模糊职责的 `common/` 目录来承接无法命名的代码

## 禁止项的正向替代

当实现者觉得某段代码“好像只能放进 `types/`、`common/` 或 `shared/`”时，必须先按下面的顺序判断：

1. 如果它带有明确业务语义，必须归属一个具体领域，由其他领域通过适配、转换或显式依赖使用。
2. 如果它跨域复用且不带业务语义，应进入 `internal/platform/` 下合适的稳定子目录，例如 `platform/http/`、`platform/support/`、`platform/db/`。
3. 如果它只服务于系统装配或启动链路，应进入 `internal/app/`，而不是伪装成“公共业务代码”。

默认原则是：宁可先明确归属到一个领域，也不要创建模糊的全局承接目录。

## 允许的弹性

这份目录契约允许三类受控弹性：

1. 目录模板是受控词汇表，不要求一开始把所有子目录都建出来。  
   任何模板子目录都按需创建；如果当前没有实现内容，就不应预建空目录。
2. 某些领域没有 `jobs/`
3. 某些领域未来只有在修订目录规范时，才允许增加极少量补充子目录，并且必须满足两个前提：
   - 现有 `api / usecase / store / model / policy / jobs` 无法合理承载
   - 新名称可以在所有相似领域稳定复用，而不是一次性例外

## 决策结果

`server-v2/` 采用：

- 单一服务入口的 `cmd/server/main.go + internal/` 外骨架
- `cmd/` 在初始状态只保留 `server/`，未来只有出现新的真实可执行服务时才允许新增 `cmd/<service>/`
- `internal/domains/` 作为业务主视角
- 每个业务域内部使用固定子目录模板
- `tests/` 作为默认测试落点，包内白盒单测为受控例外
- `platform/` 严格限制为跨域基础设施

一句话总结：

**顶层少而稳，域内模板固定，测试镜像实现，平台层绝不吞业务。**
