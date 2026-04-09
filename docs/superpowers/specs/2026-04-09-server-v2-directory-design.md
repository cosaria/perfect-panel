# Server-V2 目录设计

## 背景

当前仓库里的 [server/README.md](/Users/admin/Codes/ProxyCode/perfect-panel/server/README.md) 已经把 `server/` 收敛到 `cmd + internal` 外骨架，但内部仍承载多轮重构遗留的复杂语义。

这次已经明确决定新建 `server-v2/`，目标不是继续在旧 `server/` 上修补，而是为一个全新的控制面后端建立稳定、可预测、适合 AI Vibe Coding 的目录结构。

目录设计的优先目标不是传统团队分层美观，而是让 AI 以尽量少的搜索和跳转成本定位代码、理解边界并安全修改。

## 目标

1. 遵循 Go 服务项目常见外骨架，保持 `cmd/` 和 `internal/` 的稳定组织方式。
2. 让 AI 在绝大多数改动里以稳定路径找到代码，减少全仓库横跳。
3. 将业务语义与平台基础设施明确隔离。
4. 为后续 `server-v2` 的 spec、implementation plan 和实现阶段提供可直接沿用的目录契约。

## 非目标

1. 本文不设计仓库根目录未来的整体布局。
2. 本文不设计 `server-v2/` 与现有 `server/` 的并存、迁移或 cutover 策略。
3. 本文不展开数据库 schema、HTTP API、节点协议或任务系统细节。
4. 本文不要求先解决所有未来安全、运维和部署增强项。

## 核心判断

从 AI 视角看，最差的目录不是“目录多”，而是“落点不可预测”。  
因此，`server-v2/` 不采用纯技术横切结构，也不采用完全自由生长的领域切片，而采用：

**Go 官方风格外壳 + 稳定领域切片 + 固定子目录模板**

这个判断对应 4 条工作原则：

1. 顶层目录尽量少，且语义稳定。
2. 业务代码先按领域定位，再按固定子目录细分。
3. 基础设施只放进 `platform/`，不能吞业务。
4. 测试统一进入 `tests/`，并镜像实现路径。

## 推荐目录树

```text
server-v2/
├── cmd/
│   ├── ppanel/
│   │   ├── main.go
│   │   ├── run.go
│   │   ├── db_bootstrap.go
│   │   ├── db_seed.go
│   │   └── openapi.go
│   └── openapi/
│       └── main.go
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
│   │   ├── auth/
│   │   ├── billing/
│   │   ├── catalog/
│   │   ├── node/
│   │   ├── subscription/
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

## 顶层目录职责

### `cmd/`

只放命令入口。  
`cmd/` 里的代码负责把 CLI 形状暴露出来，不负责承载业务流程或数据库读写。

### `internal/app/`

只负责系统装配。

- `bootstrap/`：启动顺序与初始化流程
- `runtime/`：运行时容器、生命周期管理
- `routing/`：总路由树拼装
- `wiring/`：依赖注入与模块装配

`app/` 绝不能成为“先放这里以后再整理”的临时业务目录。

### `internal/platform/`

只放跨域基础设施。

- `cache/`：缓存与 Redis 适配
- `config/`：配置读取与解析
- `db/`：数据库连接、事务包装、通用持久层基础
- `http/`：通用 HTTP 中间件、响应封装、无业务语义 transport helper
- `observability/`：日志、trace、metrics
- `queue/`：异步任务执行基础设施
- `support/`：少量真正跨域的通用帮助能力

判断标准很简单：如果一个包带有明确业务语义，就不应放在 `platform/`。

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

这条规则优先于 Go 常见的“测试文件跟实现同目录”做法，因为当前目标是提高 AI 搜索稳定性，而不是优先遵循传统手感。

## 领域目录模板

每个领域都遵循固定子目录模板，不允许自由发散命名。

### `api/`

放该领域自己的 handler、DTO、transport adapter。  
它负责输入输出边界，不负责复杂业务流程。

### `usecase/`

放业务流程编排。  
这是领域行为的主要承载处，也是 AI 改需求时最常进入的目录。

### `store/`

放该领域自己的数据读写。  
它只负责持久化访问，不负责规则判断。

### `model/`

放该领域核心结构、值对象、枚举和领域内聚合需要的基础表示。  
禁止出现全局横切 `models/` 顶层目录。

### `policy/`

放规则、校验、状态迁移、授权判断等逻辑。  
这类代码如果混进 `usecase/`，文件会很快变得过重。

### `jobs/`

只有该领域确实存在异步消费或后台任务时才出现。  
不是所有域都必须有 `jobs/`。

## AI 友好性说明

这套目录不是为了“理论最优分层”，而是为了让 AI 的路径判断尽量固定。

典型定位路径如下：

- 改登录流程：先看 `internal/domains/auth/usecase/`
- 改支付回调：先看 `internal/domains/billing/api/`，再看 `internal/domains/billing/usecase/`
- 改订阅开通规则：先看 `internal/domains/subscription/policy/`
- 改节点下发：先看 `internal/domains/node/usecase/` 或 `internal/domains/node/api/`
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

## 允许的弹性

这份目录契约允许两类受控弹性：

1. 某些领域没有 `jobs/`
2. 某些领域未来只有在修订目录规范时，才允许增加极少量补充子目录，并且必须满足两个前提：
   - 现有 `api / usecase / store / model / policy / jobs` 无法合理承载
   - 新名称可以在所有相似领域稳定复用，而不是一次性例外

## 决策结果

`server-v2/` 采用：

- Go 服务项目常见的 `cmd/ + internal/` 外骨架
- `internal/domains/` 作为业务主视角
- 每个业务域内部使用固定子目录模板
- `tests/` 统一镜像实现路径
- `platform/` 严格限制为跨域基础设施

一句话总结：

**顶层少而稳，域内模板固定，测试镜像实现，平台层绝不吞业务。**
