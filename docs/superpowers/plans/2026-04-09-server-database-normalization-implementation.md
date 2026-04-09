# Server 数据库规范化重构 Implementation Plan

> **面向 agentic workers：** 必选子技能：使用 `superpowers:subagent-driven-development`（推荐）或 `superpowers:executing-plans` 按任务执行本计划。步骤使用 checkbox (`- [ ]`) 语法跟踪。

**Goal:** 在不变更当前 API、OpenAPI 契约、节点协议和支付回调协议的前提下，用规范化 schema、关系表、兼容 repository 和 revision 工作流替换旧 migration 驱动的平铺持久层。

**Architecture:** 先建立 `schema/bootstrap + revisions + db CLI` 骨架，再按 `identity + system -> catalog + node relation -> billing + subscription -> node ingestion + async compatibility -> content + cleanup` 的顺序分 phase 落地。每个 phase 都以“现有 HTTP 面不变、核心集成测试仍通过”为验收标准，旧 `user/auth/order/payment/subscribe/...` 包只作为兼容 façade 继续对上层服务暴露既有接口。

**Tech Stack:** Go 1.25、Cobra、Gin、Huma、GORM、MySQL、Redis、Asynq、Testify、Make

---

## 范围说明

这份 spec 覆盖多个子系统，不能按“一次大重写”执行。这个计划是主实施计划，但执行时必须按 phase 落地，每一段都能独立编译、独立测试、独立回滚。

## 文件与职责锁定

- `server/cmd/ppanel/root.go`
  责任：注册新的 `db` 命令树，同时保住现有 `run`、`version`、`openapi` 行为。
- `server/cmd/ppanel/db*.go`
  责任：承载 `db bootstrap`、`db seed`、`db reset`、`db revisions` 等数据库工作流。
- `server/internal/platform/persistence/schema/*.go`
  责任：定义 schema registry、bootstrap runner、revision runner、命名规则和基础 seed。
- `server/internal/platform/persistence/schema/revisions/*.go`
  责任：声明 forward-only schema revisions，替代旧 SQL migration 历史链。
- `server/internal/platform/persistence/identity/*.go`
  责任：承载规范化后的用户主体、认证标识、会话、设备、验证码和安全事件。
- `server/internal/platform/persistence/catalog/*.go`
  责任：承载套餐、价格、能力项、节点组和套餐到节点组关系。
- `server/internal/platform/persistence/billing/*.go`
  责任：承载订单、支付、支付回调事实、支付网关配置、退款和账本。
- `server/internal/platform/persistence/subscription/*.go`
  责任：承载用户订阅、周期、token、usage snapshot、assignment 和订阅事件。
- `server/internal/platform/persistence/node/*.go`
  责任：承载 `servers`、`nodes`、协议 profile、运行态事实、节点流量写入和授权读模型。
- `server/internal/platform/persistence/content/*.go`
  责任：承载工单、公告、文档和消息模板的新规范化实现。
- `server/internal/platform/persistence/system/*.go`
  责任：承载站点配置、认证策略、计费策略、节点策略、后台任务和 outbox。
- `server/internal/platform/persistence/{user,auth,order,payment,subscribe,announcement,document,ticket}/*.go`
  责任：作为兼容 façade，继续为现有 service layer 暴露旧接口，但内部委托给新模块，禁止继续直接承载新 schema 设计。
- `server/internal/bootstrap/configinit/*.go`
  责任：从旧 migration 初始化切到 `schema/bootstrap` 工作流，并保住现有 init 页面和配置写入能力。
- `server/internal/bootstrap/app/serviceContext.go`
  责任：逐步把旧 `Model` 注入切成“兼容 façade + 新模块 runner”的混合装配。
- `server/internal/domains/**`
  责任：保持 handler 和逻辑入口不变，但逐步从旧平铺模型切到规范化持久层。
- `server/internal/jobs/**`
  责任：补齐异步幂等、回放审计、节点流量归属和订单激活的兼容行为。
- `server/cmd/*_test.go`
  责任：继续作为跨 phase 护栏，补充 db workflow、compatibility 和 legacy migration 退场测试。

## 迁移约束

- 现有 `docs/openapi/*.json` 和对应路由 contract 不得因为数据库重构而变化。
- 现有节点上报协议、节点轮询协议、支付回调协议不得换字段或换语义。
- 所有新 schema 变化必须走 `schema/revisions/`，不再向 `persistence/migrate/database/*.sql` 添加任何内容。
- 所有阶段都优先保住 `cd server && go test ./...`、`cd server && go build ./...`、`cd server && go run . openapi -o ../docs/openapi`。
- 只有当兼容 façade 已经接住旧 service layer 时，才能删除旧 migration 调用或旧关系字段读取逻辑。

### Task 1: 建立 `schema/bootstrap + revisions + db CLI` 骨架

**Files:**
- Create: `server/cmd/ppanel/db.go`
- Create: `server/cmd/ppanel/db_bootstrap.go`
- Create: `server/cmd/ppanel/db_seed.go`
- Create: `server/cmd/ppanel/db_reset.go`
- Create: `server/cmd/ppanel/db_revisions.go`
- Create: `server/internal/platform/persistence/schema/bootstrap.go`
- Create: `server/internal/platform/persistence/schema/registry.go`
- Create: `server/internal/platform/persistence/schema/naming.go`
- Create: `server/internal/platform/persistence/schema/revisions.go`
- Create: `server/internal/platform/persistence/schema/seed/admin_seed.go`
- Create: `server/internal/platform/persistence/schema/seed/site_seed.go`
- Create: `server/internal/platform/persistence/schema/revisions/0001_baseline.go`
- Create: `server/internal/platform/persistence/schema/bootstrap_test.go`
- Create: `server/internal/platform/persistence/schema/revisions_test.go`
- Create: `server/cmd/ppanel/db_command_test.go`
- Modify: `server/cmd/ppanel/root.go`
- Modify: `server/internal/bootstrap/configinit/config.go`
- Modify: `server/internal/bootstrap/configinit/version.go`
- Test: `server/cmd/ppanel/db_command_test.go`
- Test: `server/internal/platform/persistence/schema/bootstrap_test.go`

- [ ] **Step 1: 写 `db` 命令与 bootstrap workflow 的失败测试**

Run: `cd server && go test ./cmd/ppanel -run 'TestDBCommandTree|TestDBBootstrapRejectsUnknownRevisionSource' -count=1`
Expected: FAIL，提示 `db` 子命令未注册，或 schema bootstrap runner 不存在。

- [ ] **Step 2: 搭建 `schema` 骨架，并让 `db bootstrap` 能在空库上创建 registry 和 baseline revision**

实现要求：
`bootstrap.go` 负责建立 schema registry 表并执行 baseline。
`revisions.go` 负责 forward-only revision runner。
`0001_baseline.go` 只做“接管未来工作流”的最小 schema 记录，不直接复制旧 SQL 文件。

- [ ] **Step 3: 注册 `db bootstrap`、`db seed`、`db reset`、`db revisions` 子命令**

Run: `cd server && go test ./cmd/ppanel -run TestDBCommandTree -count=1`
Expected: PASS，命令树包含 `db bootstrap`、`db seed`、`db reset`、`db revisions`。

- [ ] **Step 4: 把 init/bootstrap 链路从旧 migration 调用切到 `schema/bootstrap`**

重点修改：
`server/internal/bootstrap/configinit/config.go`
`server/internal/bootstrap/configinit/version.go`

要求：
初始化页面仍能写配置、测试 MySQL、创建管理员，但底层不再调用 `persistence/migrate.Migrate(...).Up()`。

- [ ] **Step 5: 回跑 schema 与 CLI 测试**

Run: `cd server && go test ./cmd/ppanel ./internal/platform/persistence/schema -count=1`
Expected: PASS

- [ ] **Step 6: 提交骨架阶段**

```bash
git add server/cmd/ppanel/db*.go server/internal/platform/persistence/schema server/internal/bootstrap/configinit/config.go server/internal/bootstrap/configinit/version.go
git commit -m "feat(server): add schema bootstrap and revision workflow"
```

### Task 2: 规范化 `identity + system`，并保住现有 `user/auth/system` 接口

**Files:**
- Create: `server/internal/platform/persistence/identity/models.go`
- Create: `server/internal/platform/persistence/identity/repository.go`
- Create: `server/internal/platform/persistence/identity/security_events.go`
- Create: `server/internal/platform/persistence/identity/verification.go`
- Create: `server/internal/platform/persistence/identity/compat_test.go`
- Create: `server/internal/platform/persistence/system/policies.go`
- Create: `server/internal/platform/persistence/system/settings.go`
- Create: `server/internal/platform/persistence/schema/revisions/0002_identity_system.go`
- Modify: `server/internal/platform/persistence/user/default.go`
- Modify: `server/internal/platform/persistence/user/authMethod.go`
- Modify: `server/internal/platform/persistence/user/device.go`
- Modify: `server/internal/platform/persistence/auth/auth.go`
- Modify: `server/internal/platform/persistence/system/model.go`
- Modify: `server/internal/bootstrap/app/serviceContext.go`
- Modify: `server/internal/domains/auth/*.go`
- Modify: `server/internal/domains/admin/authMethod/*.go`
- Modify: `server/internal/domains/admin/system/*.go`
- Test: `server/internal/platform/persistence/identity/compat_test.go`
- Test: `server/internal/domains/admin/authMethod/phase6_runtime_reload_test.go`

- [ ] **Step 1: 写 identity 兼容测试**

覆盖点：
旧 `UserModel.FindOneByEmail`、`FindUserAuthMethods`、`InsertUserAuthMethods`、`FindOneDevice` 在新 schema 下仍返回相同语义。

Run: `cd server && go test ./internal/platform/persistence/identity -run TestIdentityCompatibility -count=1`
Expected: FAIL，identity 模块和 façade 尚未接通。

- [ ] **Step 2: 用 revision 建立 `users`、`user_auth_identities`、`user_sessions`、`user_devices`、`verification_tokens`、`auth_providers`、`auth_provider_configs`、`verification_policies`、`verification_deliveries`、`security_events`**

要求：
当前保持单一 `users` 主体，不拆 `admin_users` / `end_users`。
`is_admin`、现有登录方式和现有策略读取行为必须仍可表达。

- [ ] **Step 3: 让旧 `user/auth/system` 包退化为兼容 façade**

要求：
旧 service layer 不改函数签名。
新 schema 读写落到 `identity` 与 `system` 新 repository。
不再向旧 `system` 类万能结构里堆新的 auth config 语义。

- [ ] **Step 4: 回跑认证与系统配置相关测试**

Run: `cd server && go test ./internal/platform/persistence/identity ./internal/domains/auth ./internal/domains/admin/authMethod ./internal/domains/admin/system -count=1`
Expected: PASS

- [ ] **Step 5: 提交 `identity + system` 阶段**

```bash
git add server/internal/platform/persistence/identity server/internal/platform/persistence/system server/internal/platform/persistence/user server/internal/platform/persistence/auth server/internal/bootstrap/app/serviceContext.go server/internal/domains/auth server/internal/domains/admin/authMethod server/internal/domains/admin/system
git commit -m "refactor(server): normalize identity and system persistence"
```

### Task 3: 规范化 `catalog + node relation`，去掉逗号字段关系

**Files:**
- Create: `server/internal/platform/persistence/catalog/models.go`
- Create: `server/internal/platform/persistence/catalog/repository.go`
- Create: `server/internal/platform/persistence/catalog/relation_test.go`
- Create: `server/internal/platform/persistence/node/assignment_repository.go`
- Create: `server/internal/platform/persistence/schema/revisions/0003_catalog_node_relations.go`
- Modify: `server/internal/platform/persistence/subscribe/model.go`
- Modify: `server/internal/platform/persistence/node/model.go`
- Modify: `server/internal/platform/persistence/node/server.go`
- Modify: `server/internal/domains/node/getServerUserList.go`
- Modify: `server/internal/domains/admin/server/*.go`
- Modify: `server/internal/domains/admin/subscribe/*.go`
- Test: `server/internal/platform/persistence/catalog/relation_test.go`
- Test: `server/internal/domains/node/phase8_assignment_contract_test.go`

- [ ] **Step 1: 写节点授权解析失败测试**

覆盖点：
套餐到节点组、节点组到节点、订阅到节点 assignment 都以关系表驱动。
`GetServerUserList` 不再依赖 `FIND_IN_SET` 或空结果伪用户兜底。

Run: `cd server && go test ./internal/domains/node -run TestAssignmentDrivenServerUserList -count=1`
Expected: FAIL

- [ ] **Step 2: 用 revision 建立 `node_groups`、`node_group_nodes`、`plan_node_group_rules`、`subscription_node_assignments`**

要求：
先保住当前 `subscribe`、`server`、`node` 词汇和 API，不新增对外名词。

- [ ] **Step 3: 让 `subscribe` 和 `node` 旧包只做兼容读取**

要求：
`FilterList`、`FilterNodeList`、`GetServerUserList` 的底层关系解析改为关系表和 assignment 读模型。
禁止继续新增逗号字段关系。

- [ ] **Step 4: 回跑 catalog/node 相关测试**

Run: `cd server && go test ./internal/platform/persistence/catalog ./internal/platform/persistence/node ./internal/domains/node ./internal/domains/admin/server ./internal/domains/admin/subscribe -count=1`
Expected: PASS

- [ ] **Step 5: 提交 `catalog + node relation` 阶段**

```bash
git add server/internal/platform/persistence/catalog server/internal/platform/persistence/node server/internal/platform/persistence/subscribe server/internal/domains/node server/internal/domains/admin/server server/internal/domains/admin/subscribe
git commit -m "refactor(server): normalize catalog and node relations"
```

### Task 4: 规范化 `billing + subscription`，保住现有下单、支付、订阅查询接口

**Files:**
- Create: `server/internal/platform/persistence/billing/models.go`
- Create: `server/internal/platform/persistence/billing/repository.go`
- Create: `server/internal/platform/persistence/billing/payment_callback.go`
- Create: `server/internal/platform/persistence/billing/compat_test.go`
- Create: `server/internal/platform/persistence/subscription/models.go`
- Create: `server/internal/platform/persistence/subscription/repository.go`
- Create: `server/internal/platform/persistence/schema/revisions/0004_billing_subscription.go`
- Modify: `server/internal/platform/persistence/order/model.go`
- Modify: `server/internal/platform/persistence/payment/model.go`
- Modify: `server/internal/platform/persistence/subscribe/model.go`
- Modify: `server/internal/platform/persistence/user/model.go`
- Modify: `server/internal/domains/user/portal/*.go`
- Modify: `server/internal/domains/admin/order/*.go`
- Modify: `server/internal/domains/admin/payment/*.go`
- Modify: `server/internal/domains/admin/user/*.go`
- Test: `server/internal/platform/persistence/billing/compat_test.go`
- Test: `server/internal/domains/user/portal/queryPurchaseOrder_test.go`

- [ ] **Step 1: 写下单与订阅兼容测试**

覆盖点：
现有下单、支付方式获取、订单详情、订阅 token 查询和管理员侧订单查询语义不变。

Run: `cd server && go test ./internal/platform/persistence/billing ./internal/domains/user/portal ./internal/domains/admin/order ./internal/domains/admin/payment -count=1`
Expected: FAIL

- [ ] **Step 2: 用 revision 建立 `orders`、`order_items`、`payments`、`payment_callbacks`、`payment_gateways`、`payment_gateway_secrets`、`refunds`、`billing_ledgers`、`subscriptions`、`subscription_periods`、`subscription_tokens`、`subscription_usage_snapshots`、`subscription_events`**

要求：
支付方式配置与支付事实拆开。
订阅目录与用户持有结果拆开。

- [ ] **Step 3: 让旧 `order/payment/subscribe/user` 包转成兼容 façade**

要求：
现有 handler / job 不改函数签名。
旧包内部查询与写入都落到新 `billing` / `subscription` 模块。

- [ ] **Step 4: 回跑订单、支付、订阅测试**

Run: `cd server && go test ./internal/platform/persistence/billing ./internal/platform/persistence/subscription ./internal/domains/user/portal ./internal/domains/admin/order ./internal/domains/admin/payment ./internal/domains/admin/user -count=1`
Expected: PASS

- [ ] **Step 5: 提交 `billing + subscription` 阶段**

```bash
git add server/internal/platform/persistence/billing server/internal/platform/persistence/subscription server/internal/platform/persistence/order server/internal/platform/persistence/payment server/internal/platform/persistence/subscribe server/internal/platform/persistence/user server/internal/domains/user/portal server/internal/domains/admin/order server/internal/domains/admin/payment server/internal/domains/admin/user
git commit -m "refactor(server): normalize billing and subscription persistence"
```

### Task 5: 补齐外部写入口信任、异步幂等和节点流量归属

**Files:**
- Create: `server/internal/platform/persistence/system/external_trust.go`
- Create: `server/internal/platform/persistence/node/usage_ingest.go`
- Create: `server/internal/platform/persistence/node/usage_ingest_test.go`
- Create: `server/internal/platform/persistence/billing/payment_idempotency_test.go`
- Create: `server/internal/platform/persistence/schema/revisions/0005_async_trust_and_usage.go`
- Modify: `server/internal/platform/http/middleware/notifyMiddleware.go`
- Modify: `server/internal/platform/http/notify/paymentNotify.go`
- Modify: `server/internal/platform/http/notify/stripeNotify.go`
- Modify: `server/internal/platform/http/notify/ePayNotify.go`
- Modify: `server/internal/domains/node/serverPushUserTraffic.go`
- Modify: `server/internal/jobs/order/activateOrderLogic.go`
- Modify: `server/internal/jobs/task/quotaLogic.go`
- Modify: `server/internal/jobs/traffic/*.go`
- Test: `server/internal/platform/persistence/billing/payment_idempotency_test.go`
- Test: `server/internal/platform/persistence/node/usage_ingest_test.go`
- Test: `server/internal/platform/http/notify/phase5_protocol_surface_test.go`

- [ ] **Step 1: 写支付回调与节点上报的幂等失败测试**

覆盖点：
支付回调重复投递只生效一次。
订单激活重复消费不重复履约。
节点流量重复上报不重复记账。

Run: `cd server && go test ./internal/platform/http/notify ./internal/platform/persistence/billing ./internal/platform/persistence/node ./internal/jobs/order ./internal/jobs/task ./internal/jobs/traffic -count=1`
Expected: FAIL

- [ ] **Step 2: 建立外部写入口信任与事实表**

要求：
支付回调、节点上报都必须记录认证结果、原始事实、幂等键和处理结果。
不允许“写日志然后 return nil”继续作为默认失败处理模型。

- [ ] **Step 3: 改造回调和上报入口为“先记事实，再做幂等处理”**

重点修改：
`notifyMiddleware.go`
`paymentNotify.go`
`stripeNotify.go`
`ePayNotify.go`
`serverPushUserTraffic.go`
`activateOrderLogic.go`
`quotaLogic.go`

- [ ] **Step 4: 回跑 notify / jobs / node ingestion 测试**

Run: `cd server && go test ./internal/platform/http/notify ./internal/platform/persistence/billing ./internal/platform/persistence/node ./internal/jobs/order ./internal/jobs/task ./internal/jobs/traffic -count=1`
Expected: PASS

- [ ] **Step 5: 提交幂等与信任阶段**

```bash
git add server/internal/platform/persistence/system server/internal/platform/persistence/node server/internal/platform/persistence/billing server/internal/platform/http/middleware/notifyMiddleware.go server/internal/platform/http/notify server/internal/domains/node/serverPushUserTraffic.go server/internal/jobs/order server/internal/jobs/task server/internal/jobs/traffic
git commit -m "refactor(server): add async idempotency and external trust contracts"
```

### Task 6: 整理 `content`、切走旧 migration 调用，并做全链路兼容验收

**Files:**
- Create: `server/internal/platform/persistence/content/models.go`
- Create: `server/internal/platform/persistence/content/repository.go`
- Create: `server/internal/platform/persistence/content/compat_test.go`
- Create: `server/internal/platform/persistence/schema/revisions/0006_content_cleanup.go`
- Create: `server/cmd/phase8_database_normalization_contract_test.go`
- Modify: `server/internal/platform/persistence/announcement/*.go`
- Modify: `server/internal/platform/persistence/document/*.go`
- Modify: `server/internal/platform/persistence/ticket/*.go`
- Modify: `server/internal/bootstrap/app/serviceContext.go`
- Modify: `server/internal/bootstrap/configinit/*.go`
- Modify: `server/internal/domains/admin/announcement/*.go`
- Modify: `server/internal/domains/admin/document/*.go`
- Modify: `server/internal/domains/admin/ticket/*.go`
- Modify: `server/internal/domains/user/announcement/*.go`
- Modify: `server/internal/domains/user/document/*.go`
- Modify: `server/internal/domains/user/ticket/*.go`
- Modify: `server/README.md`
- Modify: `docs/superpowers/specs/2026-04-09-server-new-database-design.md`
- Test: `server/cmd/phase8_database_normalization_contract_test.go`

- [ ] **Step 1: 写数据库规范化总护栏测试**

覆盖点：
`persistence/migrate` 不再被 `configinit` 和 `run` 主链依赖。
现有 OpenAPI 导出仍通过。
现有核心路径注册、下单、订阅、节点、工单相关 contract 仍在。

Run: `cd server && go test ./cmd -run TestPhase8DatabaseNormalizationContracts -count=1`
Expected: FAIL

- [ ] **Step 2: 规范化 `announcement/document/ticket`，并让旧包退成 content façade**

要求：
工单消息与工单主表分离。
公告、文档、模板不再继续绑定万能配置结构。

- [ ] **Step 3: 删除主链上对旧 migration runner 的依赖**

要求：
`run`、`configinit`、`version` 主链都只依赖 `schema/bootstrap + revisions`。
旧 `persistence/migrate` 包保留只用于过渡测试时，后续单独清退。

- [ ] **Step 4: 进行全链路兼容验收**

Run:
`cd server && go test ./... -count=1`
`cd server && go build ./...`
`cd server && go run . openapi -o ../docs/openapi`

Expected:
全部通过，且 `docs/openapi/*.json` 没有非预期 contract 变化。

- [ ] **Step 5: 更新文档并提交收口阶段**

```bash
git add server/internal/platform/persistence/content server/internal/platform/persistence/announcement server/internal/platform/persistence/document server/internal/platform/persistence/ticket server/internal/bootstrap/app/serviceContext.go server/internal/bootstrap/configinit server/cmd/phase8_database_normalization_contract_test.go server/README.md docs/superpowers/specs/2026-04-09-server-new-database-design.md
git commit -m "docs(server): finalize database normalization rollout"
```

## 自检

### Spec 覆盖

- `bootstrap / seed / reset / revisions`：Task 1
- `identity + system` 规范化：Task 2
- `catalog + node relation` 规范化：Task 3
- `billing + subscription` 规范化：Task 4
- 外部写入口信任、异步幂等、节点授权解析：Task 5
- `content`、兼容验收、migration 主链退场：Task 6

### Placeholder 扫描

- 本计划没有使用 `TBD`、`TODO`、`implement later`、`similar to task N` 一类占位语句。
- 所有 phase 都给出了明确文件路径、测试命令、预期结果和提交检查点。

### 类型与命名一致性

- 新规范化模块命名统一使用 `identity / catalog / billing / subscription / node / content / system / schema`。
- 旧兼容 façade 命名继续沿用 `user / auth / order / payment / subscribe / announcement / document / ticket`，避免上层 service 签名漂移。

## 执行交接

Plan complete and saved to `docs/superpowers/plans/2026-04-09-server-database-normalization-implementation.md`. Two execution options:

**1. Subagent-Driven（推荐）** - 我按 phase 或 task 派发新 subagent，逐段实现、逐段回看、逐段验收

**2. Inline Execution** - 我在当前会话里按这个计划连续执行，中间按检查点停下来汇报

Which approach?
