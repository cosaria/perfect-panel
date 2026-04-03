---
date: 2026-04-03
topic: api-client-codegen-migration
---

# API Client 代码生成迁移：@umijs/openapi → @hey-api/openapi-ts

## Problem Frame

perfect-panel 的前端 API client 当前由 `@umijs/openapi`（openapi2ts）从远端 Swagger JSON URL 生成，依赖外部仓库 `ppanel-docs`。monorepo 化后前后端已在同一仓库，该依赖链路不合理。spec 生成链路（goctl-openapi）已在 feat/monorepo-baseline 分支中落地。

本文档聚焦**第二步**：用 `@hey-api/openapi-ts` 替换 `@umijs/openapi`，生成 TypeScript SDK + TanStack Query hooks，消除手动 `queryFn` / `queryKey` 管理。

**两个 app 均已使用 `@tanstack/react-query` v5**。影响范围：
- `apps/admin/services/`（19 个生成文件）和 `apps/user/services/`（15 个生成文件）
- 47 个组件文件（102 处 `useQuery` / `useMutation` 调用）
- 29 个 ProTable `request` prop 使用处（全部在 admin app）
- 2 个 SSR layout 文件（`apps/admin/app/layout.tsx`、`apps/user/app/layout.tsx`）
- 10 处 `skipErrorHandler` 使用

## 核心技术发现

### 发现 1：@hey-api/openapi-ts 原生支持多 spec 配置

PR #2602 已合并，支持数组配置和 Job Matrix 两种方式：

```typescript
// 数组配置（推荐）— 每个 spec 独立输出目录
export default [
  { input: '../../docs/openapi/admin.json', output: './services/admin', plugins: [...] },
  { input: '../../docs/openapi/common.json', output: './services/common', plugins: [...] },
];
```

每个输出目录生成独立的 `client.gen.ts`、`sdk.gen.ts`、`types.gen.ts`。类型隔离天然存在，无冲突风险。但需确保多个 client 实例共享同一套拦截器配置。

### 发现 2：类型系统从全局 namespace 变为具名导出

当前 `@umijs/openapi` 生成 `declare namespace API { ... }` 全局类型，组件中直接用 `API.Ads`、`API.GetAdsListParams` 等。

`@hey-api/openapi-ts` 生成具名导出：`export type Ads = { ... }`、`export type GetAdsListData = { body: ..., query: ... }`。

**影响**：所有使用 `API.*` 的地方需改为具名导入。这是迁移体量最大的变更之一，涉及几乎每个组件文件。

### 发现 3：SDK 函数参数结构变化

当前模式：
```typescript
// params 和 options 分开传
getAdsList(params: API.GetAdsListParams, options?: { [key: string]: any })
```

hey-api 模式：
```typescript
// 统一为一个 options 对象
getAdsList(options?: { query?: GetAdsListData['query'] })
```

返回值也不同：当前返回完整 axios response（`{ data: { code, data: { list, total } } }`），hey-api 默认返回 unwrapped data。这影响所有调用点的数据解构方式。

### 发现 4：SSR 中 `skipErrorHandler` 和 `Authorization` 的传递机制

当前 SSR layout 通过 `options` 参数同时传递 `skipErrorHandler` 和 `Authorization`：
```typescript
await currentUser({ skipErrorHandler: true, Authorization }).then(...)
```

这些自定义字段通过 `...(options || {})` spread 到 axios config 中，拦截器从 `response.config` 读取。hey-api 的 SDK 函数不再支持这种透传。

SSR 场景分析：`handleError` 中 `!isBrowser()` 会提前 return（不 toast），所以 SSR 调用实际上已经天然跳过了 toast。`skipErrorHandler` 在 SSR 中的唯一作用是跳过 toast——而这已经由 `!isBrowser()` 守护覆盖。

### 发现 5：ProTable request prop 是非 useQuery 的直接调用模式

29 个 ProTable 使用处全部是 `request={async (pagination, filters) => { ... }}` 模式，内部直接 `await` SDK 函数，不走 `useQuery`。ProTable 自身管理加载状态和分页。

这些调用可继续使用 `sdk.gen.ts` 中的直接函数（R15 覆盖），但需适配参数结构和返回值解构的变化。

### 发现 6：hey-api 支持 `createClient()` 创建独立实例

```typescript
import { createClient } from './client/client.gen';
const myClient = createClient({ baseURL: '...', auth: () => token });
const response = await getFoo({ client: myClient });
```

这为 SSR 场景的 token 隔离提供了原生解决方案。

## Requirements

**Client 生成器替换**
- R3. 用 `@hey-api/openapi-ts` 替换 `@umijs/openapi`，从本地 JSON spec 生成 TypeScript SDK
- R4. 启用 `@hey-api/client-axios` 作为 HTTP client（与现有 axios 依赖一致）
- R5. 启用 `@tanstack/react-query` 插件，生成 `queryOptions` / `mutationOptions` 工厂函数
- R6. 每个 app 使用数组配置，admin 和 common spec 分别输出到 `services/admin/` 和 `services/common/`（user app 同理为 `services/user/` 和 `services/common/`），保持现有目录结构
- R7. 每个 app 维护独立的 `openapi-ts.config.ts`

**拦截器迁移**
- R8. CSR 环境：全局 client 单例 + axios interceptors，行为与当前完全一致（JWT 注入、40002-40005 登出、业务层 toast）
- R9. SSR 环境：layout.tsx 中使用 `createClient()` 创建 per-request client 实例，从 cookies 读取 token 并显式传入 `auth` 字段，不注册 toast 相关拦截器。SSR client 仅保留 40002-40005 登出逻辑（throw error 即可，layout 已有 try-catch）
- R10. `skipErrorHandler` 在 CSR 中的实现：通过 TypeScript module augmentation 扩展 `AxiosRequestConfig` 添加 `skipErrorHandler` 字段，拦截器从 `response.config.skipErrorHandler` 读取。CSR 中约 8 处使用需适配为 per-request axios config 传入方式

**组件调用层迁移**
- R11. 将 `useQuery({ queryFn: () => getSomeApi() })` 迁移到 `useQuery({ ...getSomeApiOptions() })`
- R12. 将 `useMutation({ mutationFn: createSomeApi })` 迁移到 `useMutation({ ...createSomeApiMutation() })`
- R13. ProTable `request` prop 内部的直接 SDK 调用继续使用 `sdk.gen.ts` 函数，适配参数结构（`params` → `{ query: params }`）和返回值解构变化
- R14. 全局 `API.*` 类型引用迁移为具名导入：`import type { Ads, GetAdsListData } from '@/services/admin'`
- R15. 现有不通过 TanStack Query 的直接调用（如 `await createAds({...}}`）可继续使用 `sdk.gen.ts` 中的函数

**工具链集成**
- R16. 根 `package.json` 的 `openapi` script 在 spec 生成后调用 `@hey-api/openapi-ts`（替代 `openapi2ts`）
- R17. 安装 `@hey-api/openapi-ts`、`@hey-api/client-axios` 作为 devDependencies；安装 `@hey-api/client-axios` 作为 apps 的 dependencies

## 迁移对照

```
迁移前                                   迁移后
─────────────────────────────────────    ─────────────────────────────────────
import { getAdsList }                    import { getAdsList }
  from '@/services/admin/ads'              from '@/services/admin'
                                           // 或 '@/services/admin/sdk.gen'

API.Ads                                  import type { Ads }
                                           from '@/services/admin'

useQuery({                               useQuery({
  queryKey: ['adsList', params],            ...getAdsListOptions({
  queryFn: () => getAdsList(params),          query: params,
})                                         }),
                                         })

useMutation({                            useMutation({
  mutationFn: createAds,                   ...createAdsMutation(),
})                                       })

// ProTable request (直接调用)           // ProTable request (直接调用)
const { data } = await getAdsList({      const response = await getAdsList({
  ...pagination, ...filters,               query: { ...pagination, ...filters },
});                                      });
return {                                 return {
  list: data.data?.list || [],             list: response.data?.list || [],
  total: data.data?.total || 0,            total: response.data?.total || 0,
};                                       };
                                         // 注意：实际解构方式取决于 hey-api
                                         // 的返回值结构，见下方注意事项

// SSR layout (skipErrorHandler)         // SSR layout (createClient)
await getGlobalConfig({                  const ssrClient = createClient({
  skipErrorHandler: true,                  baseURL: NEXT_PUBLIC_API_URL,
})                                       });
                                         await getGlobalConfig({
                                           client: ssrClient,
                                         })
```

> **注意**：ProTable 迁移后的返回值解构取决于 hey-api SDK 的实际返回值结构（是否自动 unwrap axios response.data）。上表中保持了 `data.data?.list` 的写法作为保守估计，实际可能简化为 `data.list`。需要在生成 SDK 后验证。

## 迁移影响矩阵

| 变更类型 | 文件数 | 调用点数 | 复杂度 | 可批量处理 |
|---|---|---|---|---|
| import 路径变更（SDK 函数） | ~47 | ~80 | 低 | 是（find-replace） |
| `API.*` 类型 → 具名导入 | ~47 | ~200+ | 中 | 部分（需确认类型名映射） |
| useQuery → xxxOptions() | ~18 | ~35 | 中 | 是（模式统一） |
| useMutation → xxxMutation() | ~29 | ~32 | 低 | 是（模式统一） |
| ProTable request 参数适配 | ~29 | ~29 | 低 | 是（模式极统一） |
| SSR layout createClient | 2 | ~6 | 高 | 否（需手工处理） |
| skipErrorHandler 适配 | ~5 | ~8 | 中 | 否（需逐个确认） |
| 返回值解构调整 | ~47 | ~100+ | 中 | 待验证 |

> **注意**：useQuery (~35) + useMutation (~32) + ProTable request 内直接调用 (~29) + SSR/其他直接调用 (~6) ≈ 102，与 Problem Frame 中的总数一致。

## Success Criteria

1. `bun run typecheck` 零错误
2. `bun run lint` 零错误
3. 所有 47 个组件文件编译无错误
4. request.ts 拦截器行为回归：JWT 注入、40002-40005 登出、业务层 toast
5. SSR layout 正常渲染，无跨请求 token 污染
6. `useQuery({ ...xxxOptions() })` 替代手动 queryFn 模式覆盖率 100%
7. ProTable 列表页数据正常加载、分页正常工作

## Scope Boundaries

- 不涉及 spec 生成链路变更（已在 feat/monorepo-baseline 完成）
- 不涉及 server 端代码变更
- 不涉及 ProTable 组件本身的重构（仅适配 SDK 调用变化）
- 不涉及新增 API 端点
- 不新增 `@hey-api/client-next`（继续使用 `@hey-api/client-axios`，避免同时引入两个 client 的额外复杂度）

## Key Decisions

- **CSR/SSR 分离策略**：CSR 用全局 client 单例 + interceptors；SSR 用 `createClient()` per-request。理由：SSR 并发场景下全局单例会导致 token 污染，`createClient()` 是 hey-api 原生支持的安全模式。
- **保持 `@hey-api/client-axios` 而非 `@hey-api/client-next`**：减少迁移变量，axios interceptor 与当前代码结构兼容度最高。
- **多 spec 使用数组配置、分目录输出**：与当前 `services/admin/` + `services/common/` 结构一致，类型天然隔离，无需处理合并冲突。
- **skipErrorHandler 通过 axios module augmentation 实现**：最小侵入性方案，保持拦截器逻辑不变，仅需声明类型扩展。

## Dependencies / Assumptions

- spec 生成链路（goctl-openapi → redocly lint）已稳定运行
- `@hey-api/openapi-ts` 对 OpenAPI 3.0.3 spec 的支持与 3.1.x 同等稳定（已确认支持 3.0.x）
- `@hey-api/openapi-ts` 的数组配置功能稳定可用（PR #2602 已合并）
- hey-api SDK 函数的返回值结构需要在实际生成后验证（是否自动 unwrap `response.data`、业务层 `{ code, data }` 的处理方式）
- `createClient()` 创建的实例支持 per-request auth 注入且不影响全局单例

## Outstanding Questions

### Resolve Before Planning

（无阻断性问题，所有关键技术问题已通过代码库分析和文档研究解决）

### Deferred to Planning

- [影响 R13][Needs research] hey-api SDK 函数的返回值结构：是返回完整 axios response 还是 unwrapped data？这决定了 ProTable request 和所有直接调用点的解构方式。需要生成一份样例 SDK 后验证。
- [影响 R10][Technical] `skipErrorHandler` 通过 axios config 传递时，hey-api SDK 的 per-request options 是否支持透传自定义 axios config 字段？如果不支持，需要在 SDK 调用层用 `{ client: createClient({ axios: customAxiosInstance }) }` 或 try-catch 替代。
- [影响 R14][Technical] `API.*` 全局 namespace 到具名导入的映射关系：hey-api 生成的类型名是否与当前 `API.Ads`、`API.GetAdsListParams` 等一一对应？需要生成后对照确认，并考虑是否需要一个兼容层（re-export as namespace）减少迁移量。
- [影响 R8][Technical] 多 spec（admin + common）各自生成独立 `client.gen.ts` 时，CSR 全局拦截器是否需要配置两次？如果是，考虑封装一个 `setupClient()` 工具函数统一初始化。
- [影响 R9][Needs research] `@hey-api/client-axios` 是否导出 `createClient()` 函数？hey-api 文档中 `createClient` 示例主要出现在 fetch/next client 文档中，需验证 axios client 是否同样支持。
- [影响 R9][Technical] SSR 环境中 `Logout()` 的行为定义：当前 `Logout()` 可能操作客户端状态（cookies、路由重定向），在 Server Component 中调用可能无效或抛异常。SSR client 的 40002-40005 处理应改为 throw error 而非调用 `Logout()`。
- [影响 R10][Needs research] `@hey-api/client-axios` 的 per-request options 类型是否基于 `AxiosRequestConfig`？module augmentation 的目标接口需确认，否则 `skipErrorHandler` 字段可能无法通过类型检查传递到拦截器。

## Next Steps

`→ /ce:plan` 规划 hey-api 迁移实施（R3-R17），建议分 3 个子任务：
1. 生成器配置 + 首次生成 SDK（验证输出结构和类型映射）
2. 拦截器迁移 + SSR client 分离
3. 47 个组件文件批量迁移

参考设计文档：`~/.gstack/projects/cosaria-perfect-panel/admin-feat-monorepo-baseline-design-20260403-141807.md`
