# API Generation Pipeline Best Practices

**Date**: 2026-04-03
**Status**: Draft
**Stack**: huma v2 + Redocly CLI + @hey-api/openapi-ts + oasdiff

---

## 1. Architecture Overview

```
Go Source Code (huma v2)
        │
        ▼
  go run . openapi          ← Code-First: Go types → OpenAPI 3.1 spec
        │
        ▼
  docs/openapi/*.json       ← 3 specs: admin.json, user.json, common.json
        │
        ├─→ redocly lint    ← Spec 质量校验
        ├─→ oasdiff         ← Breaking change 检测
        ├─→ redocly build-docs ← API 文档站点
        │
        ▼
  @hey-api/openapi-ts       ← Spec → TypeScript SDK + TanStack Query hooks
        │
        ▼
  apps/*/services/*-api/    ← 生成的 SDK（types.gen.ts, sdk.gen.ts, @tanstack/react-query.gen.ts）
```

**单向数据流**：Go types 是唯一的 source of truth。所有下游产物（spec、SDK、文档）都是派生的。

## 2. Server: huma v2 OpenAPI Generation

### 2.1 禁用 $schema 注入

当前状态：admin.json 有 187 个 `$schema` 字段，user.json 有 91 个。这些是 huma 的 `SchemaLinkTransformer` 注入的，对前端 SDK 生成是噪音。

```go
// server/internal/handler/routes.go — 每个 huma.API 初始化处
config := huma.DefaultConfig("Admin API", "1.0.0")
config.CreateHooks = nil  // 禁用 SchemaLinkTransformer
api := humagin.New(router, config)
```

**影响**：spec 文件体积减少约 5-8%，生成的 TypeScript 类型更干净（不含 `$schema?: string` 字段）。

### 2.2 GET 请求 Query Parameters

**反模式**（当前部分 handler 的做法）：

```go
// 错误：GET 请求不应有 Body
type GetAdsListInput struct {
    Body struct {
        Page int `json:"page"`
        Size int `json:"size"`
    }
}
```

**正确模式**：

```go
type GetAdsListInput struct {
    Page   int    `query:"page" default:"1" minimum:"1" doc:"页码"`
    Size   int    `query:"size" default:"20" minimum:"1" maximum:"100" doc:"每页数量"`
    Search string `query:"search,omitempty" doc:"搜索关键词"`
}
```

**迁移策略**：
1. 从 `types.go` 中的 `form` tag 提取字段名
2. 将 `form:"xxx"` 替换为 `query:"xxx"`
3. handler 注册处的 `http.MethodGet` 对应的 input struct 不再需要 Body wrapper
4. 可编写脚本批量转换（约 107 个 GET handler）

### 2.3 Required/Optional 字段控制

huma 默认：Body 中的字段 **required**，Query/Header/Cookie 参数 **optional**。

```go
// Body 字段
type CreateUserBody struct {
    Name  string `json:"name" minLength:"1"`              // required（默认）
    Email string `json:"email" format:"email"`              // required（默认）
    Bio   string `json:"bio,omitempty"`                     // optional（omitempty）
    Age   *int   `json:"age,omitempty" minimum:"0"`         // optional（指针 + omitempty）
}
```

**规则**：
- 必填字段：不加 `omitempty`，非指针类型
- 可选字段：加 `omitempty` 或使用指针类型
- 显式控制：`required:"false"` 或 `required:"true"`

**当前问题**：`types.go` 中 17 个字段仅有 `form` tag 无 `json` tag（FilterLogParams 等），需要补齐。

### 2.4 可复用的输入输出模式

```go
// 分页参数 — 所有列表接口复用
type PaginationParams struct {
    Page int `query:"page" default:"1" minimum:"1" doc:"页码"`
    Size int `query:"size" default:"10" minimum:"1" maximum:"100" doc:"每页数量"`
}

// 分页响应 — 泛型复用
type PaginatedBody[T any] struct {
    List  []T   `json:"list"`
    Total int64 `json:"total"`
}

// 使用
type ListUsersInput struct {
    PaginationParams          // 嵌入分页参数
    Search string `query:"search,omitempty" doc:"搜索关键词"`
}

type ListUsersOutput struct {
    Body PaginatedBody[User]
}
```

### 2.5 UserOpenAPI() 合并优化

当前实现通过 JSON marshal/unmarshal 合并 10 个子 API spec。这个方案可以工作但有几个改进点：

1. **添加 info.description**：每个 spec 应包含描述
2. **添加 server 信息**：spec 中应包含 base URL
3. **排序 paths**：确保输出稳定，便于 diff

## 3. Redocly: Spec Validation & Documentation

### 3.1 安装

```bash
bun add -D @redocly/cli  # 根目录 devDependency
```

### 3.2 配置文件

```yaml
# redocly.yaml — 项目根目录
extends:
  - recommended

apis:
  admin:
    root: docs/openapi/admin.json
  user:
    root: docs/openapi/user.json
  common:
    root: docs/openapi/common.json

rules:
  # 严格模式
  operation-operationId-unique: error
  operation-operationId-url-safe: error
  operation-summary: error
  no-unresolved-refs: error
  no-enum-type-mismatch: error
  security-defined: error

  # 放宽（渐进式收紧）
  no-unused-components: warn
  operation-4xx-response: warn
  tag-description: warn

  # $schema 字段的 readOnly 标记不算 lint 问题，
  # 但禁用 CreateHooks 后这条就不存在了
```

### 3.3 命令

```bash
# Lint 所有 spec
npx redocly lint docs/openapi/*.json

# 预览文档（本地开发）
npx redocly preview-docs docs/openapi/admin.json

# 生成静态文档
npx redocly build-docs docs/openapi/admin.json -o docs/api/admin.html
```

### 3.4 Makefile 集成

```makefile
# 追加到现有 Makefile
spec-lint:
	npx redocly lint docs/openapi/*.json

spec-preview:
	npx redocly preview-docs docs/openapi/admin.json
```

## 4. oasdiff: Breaking Change Detection

Redocly 不内置 breaking change 检测。使用 **oasdiff** 填补。

### 4.1 本地使用

```bash
# 安装
brew install oasdiff

# 对比当前 spec 与 main 分支
git show main:docs/openapi/admin.json > /tmp/admin-base.json
oasdiff breaking /tmp/admin-base.json docs/openapi/admin.json
```

### 4.2 CI 集成 (GitHub Actions)

```yaml
# .github/workflows/api-check.yml
name: API Spec Validation
on:
  pull_request:
    paths:
      - 'server/**'
      - 'docs/openapi/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      # 1. 从 Go 代码重新生成 spec，检查是否与提交的 spec 一致
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - name: Check spec drift
        run: |
          cd server && go run . openapi --output ../docs/openapi/
          if ! git diff --exit-code docs/openapi/; then
            echo "::error::OpenAPI specs 与代码不同步，运行 bun run openapi:spec 后提交"
            exit 1
          fi

      # 2. Redocly lint
      - run: npx @redocly/cli@latest lint docs/openapi/*.json --format=github-actions

      # 3. Breaking change 检测
      - uses: oasdiff/oasdiff-action/breaking@v0.0.37
        with:
          base: 'origin/${{ github.base_ref }}:docs/openapi/admin.json'
          revision: 'HEAD:docs/openapi/admin.json'
          fail-on: ERR
      - uses: oasdiff/oasdiff-action/breaking@v0.0.37
        with:
          base: 'origin/${{ github.base_ref }}:docs/openapi/user.json'
          revision: 'HEAD:docs/openapi/user.json'
          fail-on: ERR
```

## 5. Client: @hey-api/openapi-ts Configuration

### 5.1 合并配置文件

当前 admin app 有 3 个独立配置文件，user app 有 2 个。合并为数组配置：

**apps/admin/openapi-ts.config.ts**：
```typescript
import { defaultPlugins, defineConfig } from "@hey-api/openapi-ts";

const sharedPlugins = [
  ...defaultPlugins,
  "@hey-api/client-fetch",
  {
    name: "@hey-api/typescript",
    enums: "javascript",
  },
  {
    name: "@hey-api/sdk",
    auth: false,
  },
] as const;

export default [
  defineConfig({
    input: "../../docs/openapi/admin.json",
    output: { path: "./services/admin-api", clean: true },
    plugins: [...sharedPlugins],
  }),
  defineConfig({
    input: "../../docs/openapi/common.json",
    output: { path: "./services/common-api", clean: true },
    plugins: [...sharedPlugins],
  }),
  defineConfig({
    input: "../../docs/openapi/user.json",
    output: { path: "./services/user-api", clean: true },
    plugins: [...sharedPlugins],
  }),
];
```

**apps/user/openapi-ts.config.ts**：
```typescript
import { defaultPlugins, defineConfig } from "@hey-api/openapi-ts";

const sharedPlugins = [
  ...defaultPlugins,
  "@hey-api/client-fetch",
  {
    name: "@hey-api/typescript",
    enums: "javascript",
  },
  {
    name: "@hey-api/sdk",
    auth: false,
  },
] as const;

export default [
  defineConfig({
    input: "../../docs/openapi/user.json",
    output: { path: "./services/user-api", clean: true },
    plugins: [...sharedPlugins],
  }),
  defineConfig({
    input: "../../docs/openapi/common.json",
    output: { path: "./services/common-api", clean: true },
    plugins: [...sharedPlugins],
  }),
];
```

这样每个 app 只需要 1 个配置文件，`openapi` 脚本也简化为一条命令：

```json
{
  "openapi": "openapi-ts"
}
```

### 5.2 Client 初始化

保持现有 `setup-clients.ts` 的拦截器模式（这已经是最佳实践）：

```typescript
// apps/*/utils/setup-clients.ts
import { client as adminClient } from "@/services/admin-api/client.gen";
import { client as commonClient } from "@/services/common-api/client.gen";
import { client as userClient } from "@/services/user-api/client.gen";

function setupClient(client: typeof adminClient, serverPrefix: string) {
  client.setConfig({ baseUrl: `${baseUrl}${serverPrefix}` });

  client.interceptors.request.use((request) => {
    const token = getAuthorization();
    if (token) request.headers.set("Authorization", token);
    return request;
  });

  client.interceptors.response.use(async (response) => {
    if (response.ok) return response;
    if (response.status === 401) { Logout(); return response; }
    // error toast...
    return response;
  });
}
```

### 5.3 Enum 策略

使用 `enums: "javascript"` 生成 `as const` 对象：

```typescript
// 生成结果
export const SubscribeType = { cycle: 'cycle', traffic: 'traffic' } as const;
export type SubscribeType = (typeof SubscribeType)[keyof typeof SubscribeType];

// 使用：既有类型安全，又有运行时值
const options = Object.values(SubscribeType).map(v => ({ label: v, value: v }));
```

## 6. End-to-End Pipeline Commands

### 6.1 开发者工作流

```bash
# 完整管道（一条命令）
bun run openapi
# 等价于：
#   1. cd server && go run . openapi -o ../docs/openapi    # 生成 spec
#   2. turbo run openapi                                     # 生成 SDK
#   3. bun run format                                        # 格式化

# 加入 redocly lint 后
bun run openapi
# 等价于：
#   1. cd server && go run . openapi -o ../docs/openapi
#   2. npx redocly lint docs/openapi/*.json
#   3. turbo run openapi
#   4. bun run format
```

### 6.2 根 package.json 脚本更新

```json
{
  "scripts": {
    "openapi:spec": "cd server && go run . openapi -o ../docs/openapi",
    "openapi:lint": "redocly lint docs/openapi/*.json",
    "openapi:client": "turbo run openapi",
    "openapi": "bun run openapi:spec && bun run openapi:lint && bun run openapi:client && bun run format"
  }
}
```

### 6.3 Makefile 更新

```makefile
# 追加
spec-gen:
	cd server && go run . openapi -o ../docs/openapi

spec-lint:
	npx redocly lint docs/openapi/*.json

spec-breaking:
	@echo "Checking breaking changes against main..."
	@for f in admin common user; do \
		git show main:docs/openapi/$$f.json > /tmp/$$f-base.json 2>/dev/null && \
		oasdiff breaking /tmp/$$f-base.json docs/openapi/$$f.json || true; \
	done

openapi: spec-gen spec-lint
	bun run openapi:client
	bun run format
```

## 7. Implementation Priority

| # | Task | Impact | Effort |
|---|------|--------|--------|
| P0 | 禁用 $schema 注入 (`CreateHooks = nil`) | 清理 292 个噪音字段 | 10 min |
| P0 | 补齐 17 个 form tag 缺 json tag 的字段 | 修复 spec 中字段名大写 | 10 min |
| P1 | 合并 openapi-ts 配置文件为数组格式 | 减少 5 个文件为 2 个 | 15 min |
| P1 | 添加 redocly.yaml + lint 脚本 | Spec 质量保障 | 15 min |
| P1 | 更新根 package.json openapi 脚本加入 lint | 管道完整性 | 5 min |
| P2 | GET 请求迁移：form tag → query tag | 修复 ~107 个 handler 的参数传递方式 | 2-4 hours |
| P2 | CI: spec drift + redocly lint + oasdiff | 自动化保障 | 30 min |
| P3 | Redocly 文档站点部署 | API 文档可视化 | 1 hour |

## 8. Spec 版本管理

**将 `docs/openapi/*.json` 纳入 Git 版本控制**（当前已纳入），原因：
1. PR diff 直接可见 API 变更
2. 前端可离线开发（不依赖后端运行）
3. oasdiff 基于 git ref 比较
4. 回滚时 spec 和代码保持一致

**不将生成的 SDK 纳入 Git**（`apps/*/services/*-api/` 加入 .gitignore），原因：
1. SDK 是 spec 的完全派生物，可随时重新生成
2. 避免 PR 中大量生成代码干扰 review
3. CI 中 `bun run openapi:client` 保证一致性

> **注意**：如果项目有离线开发需求或 CI 不执行 codegen，则保留 SDK 在 Git 中也是合理选择。根据团队实际情况决定。

## 9. Security Considerations

- Spec 文件不包含真实的 API keys、密码或 PII
- `redocly.yaml` 的 `security-defined` 规则确保所有端点声明了安全需求
- `setup-clients.ts` 的拦截器负责运行时 auth，spec 中 `auth: false` 是正确的

## 10. Decision Log

| Decision | Rationale |
|----------|-----------|
| Code-First (not Spec-First) | 团队是后端驱动，Go types 已是 source of truth |
| Redocly + oasdiff (not Redocly alone) | Redocly 不内置 breaking change 检测 |
| `enums: "javascript"` | `as const` 兼顾类型安全和运行时值 |
| `auth: false` in SDK plugin | 通过拦截器统一管理，而非每次调用传 token |
| Spec 文件入 Git | 可见性 + 离线开发 + diff 基础 |
| 数组配置合并 | 减少维护负担，共享 plugin 配置 |
