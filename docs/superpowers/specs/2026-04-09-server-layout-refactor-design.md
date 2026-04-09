# Server 目录一次性重构设计

## 背景

`server/` 当前结构同时存在 `routers`、`services`、`models`、`modules`、`svc`、`initialize`、`runtime`、`worker` 等多个横切目录。

这套结构对长期人工维护还能勉强成立，但对 AI 开发不友好。一个功能改动通常需要在多个顶层目录之间跳转，先猜“这是路由问题、服务问题、模型问题，还是基础设施问题”，然后才能开始改代码。

这次重构的目标不是模仿某个 Go 模板仓库，也不是把目录改得更时髦，而是让 `server/` 更接近 Go 官方对服务型项目的推荐方向，同时显著降低 AI 定位代码的成本。

## 目标

1. 将 `server/` 收敛为更少的顶层概念，减少横切目录。
2. 保留 `server/` 作为独立 Go module 的事实。
3. 接受一次性大重排，允许打破旧 import 路径，只要求仓库自身编译、测试和运行正确。
4. 为后续按领域继续演进提供清晰骨架。

## 非目标

1. 不在本次重构中修改业务语义。
2. 不在本次重构中重写 API 合同、错误码体系或数据库 schema。
3. 不在本次重构中顺手重设计前端或根目录结构。
4. 不承诺保留 `github.com/perfect-panel/server/...` 下旧子目录的外部兼容导入路径。

## 设计原则

1. `cmd/` 只放命令入口，不放业务实现。
2. `internal/` 承载服务实现，外部不应依赖。
3. 按“领域（domain）”和“平台能力（platform）”组织代码，而不是长期维持 `services/models/modules` 这种横切分层。
4. 启动、依赖拼装、运行时状态管理属于 bootstrap 问题，不应散落在 `cmd`、`svc`、`initialize`、`runtime` 里。
5. 异步任务和调度器单独建模，不和 HTTP 入口混放。

## 目标目录树

```text
server/
├── cmd/
│   ├── ppanel/
│   │   └── main.go
│   └── openapi/
│       └── main.go
├── internal/
│   ├── bootstrap/
│   │   ├── app/
│   │   ├── configinit/
│   │   └── runtime/
│   ├── domains/
│   │   ├── admin/
│   │   ├── auth/
│   │   ├── common/
│   │   ├── node/
│   │   ├── subscribe/
│   │   ├── telegram/
│   │   └── user/
│   ├── jobs/
│   └── platform/
│       ├── cache/
│       ├── config/
│       ├── crypto/
│       ├── http/
│       ├── notify/
│       ├── payment/
│       ├── persistence/
│       └── support/
├── config/
├── doc/
├── etc/
├── tests/
└── web/
```

## 旧目录到新目录的映射

### 保留顶层

- `config/` 保留顶层
- `doc/` 保留顶层
- `etc/` 保留顶层
- `tests/` 保留顶层
- `web/` 保留顶层

### 命令入口

- `ppanel.go` -> `cmd/ppanel/main.go`
- `cmd/openapi.go` -> `cmd/openapi/main.go`
- `cmd` 下其余运行相关命令，按职责拆分后并入 `cmd/ppanel/`

### 启动与拼装

- `svc/*` -> `internal/bootstrap/app/`
- `initialize/*` -> `internal/bootstrap/configinit/`
- `runtime/*` -> `internal/bootstrap/runtime/`

### HTTP 层

- `routers/*` -> `internal/platform/http/`
- `routers/middleware/*` -> `internal/platform/http/middleware/`
- `routers/response/*` -> `internal/platform/http/response/`

### 业务领域

- `services/admin/*` -> `internal/domains/admin/*`
- `services/auth/*` -> `internal/domains/auth/*`
- `services/auth/oauth/*` -> `internal/domains/auth/oauth/*`
- `services/common/*` -> `internal/domains/common/*`
- `services/node/*` -> `internal/domains/node/*`
- `services/subscribe/*` -> `internal/domains/subscribe/*`
- `services/telegram/*` -> `internal/domains/telegram/*`
- `services/user/*` -> `internal/domains/user/*`
- `services/report/*` 先并入最接近的领域；若无法自然归属，暂放 `internal/domains/common/report/`

### 异步任务

- `worker/*` -> `internal/jobs/*`

### 平台能力

- `models/*` -> `internal/platform/persistence/*`
- `modules/cache/*` -> `internal/platform/cache/*`
- `modules/notify/*` -> `internal/platform/notify/*`
- `modules/payment/*` -> `internal/platform/payment/*`
- `modules/crypto/*` -> `internal/platform/crypto/*`
- `modules/infra/*` -> `internal/platform/support/*`
- `modules/util/*` -> `internal/platform/support/*`
- `modules/verify/*` -> `internal/platform/support/verify/*`
- `adapter/*` 暂放 `internal/platform/support/adapter/*`

## 架构意图

```text
CLI / OpenAPI
    |
    v
bootstrap
    |
    +--> platform/http --------> domains/*
    |          |                     |
    |          v                     v
    |      middleware          persistence / cache / notify / payment
    |
    +--> jobs ----------------> domains/* / platform/*
```

说明：

- `bootstrap` 负责把配置、数据库、缓存、运行时状态、任务队列拼起来。
- `platform/http` 只承载路由、middleware、响应契约和 OpenAPI 注册。
- `domains/*` 承载业务用例，不负责底层连接初始化。
- `platform/*` 承载数据库、缓存、通知、支付、辅助工具等平台实现。
- `jobs/*` 承载异步消费与调度。

## 为什么不是继续保留 `services/models/modules`

因为这三类目录本质上都是横切维度。

对人类来说，横切维度意味着“我知道模型大概在 models，服务大概在 services”。对 AI 来说，这意味着它需要在多个层之间猜测同一个功能的落点，再重新组装上下文。仓库越大，这个成本越高。

这次重构把“改订阅功能去哪找”“改登录去哪找”“改后台节点管理去哪找”变成稳定答案：

- 先找领域目录 `internal/domains/...`
- 再找配套平台实现 `internal/platform/...`
- 路由层固定在 `internal/platform/http/...`

## 风险

1. import 路径会大面积变化，属于高机械改动风险。
2. 测试、脚本、嵌入式静态资源和 OpenAPI 导出命令都需要回归。
3. `models` 向 `persistence` 迁移时，局部命名可能需要二次收敛，不能完全机械替换。
4. 某些 `modules/util` 里的历史杂项可能暂时找不到最优归属，需要先放入 `platform/support/` 过渡。

## 接受的破坏范围

本次重构接受以下破坏：

1. 旧包导入路径失效。
2. 外部依赖旧子包路径的代码无法直接编译。
3. 大量文件移动导致 git blame 历史碎片化。

本次重构不接受以下结果：

1. 仓库自身无法通过 `go test ./...`
2. `go run . run --config etc/ppanel.yaml` 无法启动
3. `go run . openapi -o ../docs/openapi` 无法导出 spec
4. `make test`、`make embed`、`make build-all` 的 server 相关链路失效

## 实施策略

一次性大重排，但仍按下面的内部顺序执行，避免纯粹乱搬：

1. 建立新骨架目录。
2. 迁移命令入口到 `cmd/ppanel` 与 `cmd/openapi`。
3. 迁移 bootstrap 相关目录，先消灭 `svc`。
4. 迁移 HTTP 层到 `internal/platform/http`。
5. 迁移 jobs。
6. 迁移 domains。
7. 迁移 platform。
8. 全量修复 import、命名和编译错误。
9. 运行测试与关键命令验证。

## 验证标准

至少通过以下验证：

```bash
cd server && go test ./...
cd server && go run . openapi -o ../docs/openapi
cd server && go build ./...
make test
```

如果运行时配置可用，再补：

```bash
cd server && go run . run --config etc/ppanel.yaml
```

## 决策结果

采用这次设计，直接执行“一次性大重排”，不保留旧 import 路径兼容层。
