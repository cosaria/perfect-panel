---
date: 2026-04-03
topic: hey-api-openapi-ts-research
---

# @hey-api/openapi-ts 开源项目用法研究报告

## 1. @hey-api/openapi-ts 基本用法

### 1.1 openapi-ts.config.ts 配置

@hey-api/openapi-ts 支持 `openapi-ts.config.ts`、`.cjs`、`.mjs` 多种格式，通过 `jiti` 加载器解析。推荐使用 `defineConfig` 辅助函数获得类型提示：

```typescript
import { defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: './docs/openapi/admin.json',   // 文件路径、URL 或 OpenAPI spec 对象
  output: 'src/client',                 // 输出目录（视为依赖，不要手动修改内容）
  plugins: [
    '@hey-api/client-axios',            // HTTP client（v0.63.0 起从顶层 client 字段移入 plugins）
    '@hey-api/typescript',              // 类型生成（默认插件之一）
    '@hey-api/sdk',                     // SDK 函数生成（默认插件之一）
    '@tanstack/react-query',            // TanStack Query hooks
  ],
});
```

**默认插件（defaultPlugins）**：`@hey-api/typescript` + `@hey-api/sdk`。当自定义 `plugins` 数组时，默认插件不会自动包含，需手动添加或使用 `defaultPlugins` 导出：

```typescript
import { defaultPlugins, defineConfig } from '@hey-api/openapi-ts';

export default defineConfig({
  input: './spec.json',
  output: 'src/client',
  plugins: [
    ...defaultPlugins,
    '@hey-api/client-axios',
    '@tanstack/react-query',
  ],
});
```

### 1.2 多 spec 文件处理

v0.73+ 支持三种多 spec 模式：

**数组配置（推荐，适合 perfect-panel）**：每个 spec 独立配置、独立输出：

```typescript
export default [
  { input: './docs/openapi/admin.json', output: 'services/admin', plugins: [...] },
  { input: './docs/openapi/common.json', output: 'services/common', plugins: [...] },
];
```

**Job Matrix**：共享配置、一一对应输出：

```typescript
export default {
  input: ['foo.yaml', 'bar.yaml'],
  output: ['src/foo', 'src/bar'],
  plugins: [...],
};
```

**合并输入**：多 spec 合并为单一输出（类型可能冲突，不推荐用于独立 API 域）：

```typescript
export default {
  input: ['foo.yaml', 'bar.yaml'],
  output: 'src/client',
};
```

社区中微服务架构项目常用**编程式方案**处理大量 spec：

```typescript
import { createClient, defaultPlugins } from '@hey-api/openapi-ts';

const services = { one: 'url1', two: 'url2', three: 'url3' };
Object.entries(services).forEach(([name, url]) => {
  createClient({
    input: url,
    output: { path: `src/http-client/types/${name}` },
    plugins: defaultPlugins.filter(p => p !== '@hey-api/sdk'),
  });
});
```

### 1.3 输出目录结构

默认生成以下文件（`.gen.ts` 后缀）：

```
src/client/
├── client.gen.ts              # HTTP client 实例（由 client plugin 生成）
├── sdk.gen.ts                 # SDK 函数（按 operation 生成 typed 函数）
├── types.gen.ts               # TypeScript 类型定义
├── @tanstack/react-query.gen.ts  # TanStack Query hooks（如果启用）
├── index.ts                   # barrel 导出文件（可配置关闭）
├── client/                    # client 脚手架（bundled，与 spec 无关）
└── core/                      # core 脚手架
```

**输出配置选项**：

```typescript
output: {
  path: 'src/client',
  clean: true,                    // 生成前清空目录（默认 true）
  fileName: {
    case: 'snake_case',           // 文件名风格
    suffix: '.gen',               // 后缀（设为 null 可去除）
  },
  postProcess: ['biome:format'],  // 后处理：支持 biome:format, biome:lint, prettier, eslint, oxlint
  // format 和 lint 选项已被 postProcess 替代
}
```

### 1.4 plugins 配置详解

#### @hey-api/typescript

```typescript
{
  name: '@hey-api/typescript',
  enums: 'javascript',    // 'javascript' | 'typescript' | 'typescript+namespace' | false
  // 推荐 'javascript'（普通对象），避免 TS enum 的类型问题
  dates: 'types',          // boolean | 'types' | 'types+transform'
  comments: true,          // 生成 JSDoc 注释
}
```

#### @hey-api/sdk

```typescript
{
  name: '@hey-api/sdk',
  auth: true,                       // 启用内置 auth 处理（仅支持 bearer / basic）
  operations: {
    strategy: 'flat',               // 'flat'（tree-shakeable 函数）| 'single'（class）
    // containerName: 'PetStore',   // class 模式下的类名
    // nesting: (operation) => [...], // 自定义方法分组和嵌套
  },
  // throwOnError 已移至 client plugin
}
```

#### @tanstack/react-query

```typescript
{
  name: '@tanstack/react-query',
  queryOptions: true,               // 生成 queryOptions 工厂函数
  queryKeys: { tags: true },        // 生成 queryKey，可含 operation tags
  // infiniteQueryOptions: true,    // 分页端点自动生成 infiniteQueryOptions
  // queryOptions: { meta: (op) => ({ id: op.id }) },  // 自定义 meta
}
```

## 2. 与 @tanstack/react-query 的集成

### 2.1 queryOptions / mutationOptions 工厂函数

生成的文件（`@tanstack/react-query.gen.ts`）导出以下类型的函数：

**Query 操作**（GET）生成 `xxxOptions()`：

```typescript
import { useQuery } from '@tanstack/react-query';
import { getPetByIdOptions } from './client/@tanstack/react-query.gen';

// 直接展开到 useQuery
const query = useQuery({
  ...getPetByIdOptions({
    path: { petId: 1 },
  }),
});
```

**Mutation 操作**（POST/PUT/DELETE）生成 `xxxMutation()`：

```typescript
import { useMutation } from '@tanstack/react-query';
import { addPetMutation } from './client/@tanstack/react-query.gen';

const addPet = useMutation({
  ...addPetMutation(),
  onError: (error) => console.log(error),
});

addPet.mutate({
  body: { name: 'Kitty' },
});
```

### 2.2 queryKey 管理策略

生成的 queryKey 包含规范化的参数和元数据：

```typescript
const queryKey = [
  {
    _id: 'getPetById',
    baseUrl: 'https://app.heyapi.dev',
    path: { petId: 1 },
  },
];
```

**两种访问方式**：

```typescript
// 方式 1：从 options 结果上的 .queryKey 属性
const options = getPetByIdOptions({ path: { petId: 1 } });
queryClient.invalidateQueries({ queryKey: options.queryKey });

// 方式 2：独立的 queryKey 函数
import { getPetByIdQueryKey } from './client/@tanstack/react-query.gen';
const key = getPetByIdQueryKey({ path: { petId: 1 } });
```

**启用 tags 后的 queryKey**（适合按域批量失效缓存）：

```typescript
// 配置 queryKeys: { tags: true } 后，queryKey 包含 tags
// e.g. ['pets', 'one', 'get', { _id: 'getPetById', ... }]
```

### 2.3 在 Next.js App Router 中的使用方式

perfect-panel 已使用 `@tanstack/react-query-next-experimental` 的 `ReactQueryStreamedHydration`，这是最先进的集成模式。

**当前模式（保持不变）**：

```tsx
// providers.tsx
"use client";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ReactQueryStreamedHydration } from "@tanstack/react-query-next-experimental";

export default function Providers({ children }) {
  const [queryClient] = useState(() => new QueryClient({...}));
  return (
    <QueryClientProvider client={queryClient}>
      <ReactQueryStreamedHydration>{children}</ReactQueryStreamedHydration>
    </QueryClientProvider>
  );
}
```

**迁移后的组件用法变化**：

```tsx
// 迁移前
const { data } = useQuery({
  queryKey: ['userList', params],
  queryFn: () => getUserList(params),
});

// 迁移后
const { data } = useQuery({
  ...getUserListOptions({ query: params }),
});
```

### 2.4 SSR/RSC 环境注意事项

- **Streaming Hydration**：`ReactQueryStreamedHydration` 已处理 SSR → 客户端的数据流转，不需要手动 `prefetchQuery` + `dehydrate`
- **Server Component 直接调用**：layout.tsx 中的 SSR 调用（如 `currentUser()`、`getGlobalConfig()`）不走 TanStack Query，直接调用 SDK 函数即可
- **QueryClient 隔离**：当前每个请求通过 `useState` 创建新 `QueryClient`，已是正确做法
- **不要在 Server Component 中使用 `useQuery`**：hooks 只能在客户端组件中使用

## 3. 与 @hey-api/client-axios 的集成

### 3.1 client 单例初始化

```typescript
// openapi-ts.config.ts
export default defineConfig({
  input: './spec.json',
  output: 'src/client',
  plugins: [
    {
      name: '@hey-api/client-axios',
      runtimeConfigPath: './src/hey-api.ts',  // 推荐：自定义运行时配置文件
    },
  ],
});
```

**runtimeConfigPath 方案（推荐）**：

```typescript
// src/hey-api.ts
import type { CreateClientConfig } from './client/client.gen';

export const createClientConfig: CreateClientConfig = (config) => ({
  ...config,
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'https://example.com',
});
```

**setConfig 方案（简单但有时序风险）**：

```typescript
import { client } from './client/client.gen';
client.setConfig({
  baseURL: 'https://example.com',
});
```

### 3.2 拦截器配置

```typescript
import { client } from './client/client.gen';

// 请求拦截器
client.instance.interceptors.request.use(
  (config) => {
    const token = getToken();
    if (token) config.headers.set('Authorization', `Bearer ${token}`);
    return config;
  },
  (error) => Promise.reject(error),
);

// 响应拦截器
client.instance.interceptors.response.use(
  (response) => {
    const { code } = response.data;
    if (code !== 200) {
      handleError(response);
      throw response;
    }
    return response;
  },
  (error) => {
    handleError(error);
    return Promise.reject(error);
  },
);
```

**关键发现**：`client.instance` 暴露底层 axios 实例，拦截器 API 与原生 axios 完全一致。这意味着 perfect-panel 当前的 `request.ts` 拦截器逻辑可以几乎原样迁移。

### 3.3 JWT token 注入最佳实践

**方式 1：auth 配置（推荐用于简单场景）**：

```typescript
client.setConfig({
  auth: () => getToken(),  // 自动附加到需要 auth 的请求
  baseURL: 'https://example.com',
});
```

**方式 2：请求拦截器（推荐用于 perfect-panel，控制力更强）**：

```typescript
client.instance.interceptors.request.use((config) => {
  const Authorization = getAuthorization(config.Authorization);
  if (Authorization) config.headers.Authorization = Authorization;
  return config;
});
```

**方式 3：SSR per-request client**：

```typescript
import { createClient } from './client/client.gen';

// 每个 SSR 请求创建独立 client，避免 token 污染
const ssrClient = createClient({
  baseURL: NEXT_PUBLIC_API_URL,
  auth: () => tokenFromCookies,
});

const data = await getGlobalConfig({ client: ssrClient });
```

### 3.4 错误处理和 toast 通知

perfect-panel 当前的错误处理模式：
1. 业务码 40002-40005 → 触发 `Logout()`
2. `skipErrorHandler` 为 true → 跳过 toast
3. 非浏览器环境 → 跳过 toast（`!isBrowser()`）
4. 其他错误 → `toast.error(message)`

**迁移方案**：拦截器逻辑保持不变，从 `response.config` 读取自定义字段。

### 3.5 skipErrorHandler 的实现方式

hey-api SDK 函数不支持直接透传自定义 axios config 字段。有以下替代方案：

**方案 A：TypeScript module augmentation（推荐）**

```typescript
// types/axios.d.ts
import 'axios';
declare module 'axios' {
  interface InternalAxiosRequestConfig {
    skipErrorHandler?: boolean;
  }
}
```

在调用时通过 per-request axios config 传入：

```typescript
// 对于直接 SDK 调用
const response = await someApi({
  axios: { skipErrorHandler: true },  // 如果 hey-api 支持透传
});

// 如果不支持透传，用 try-catch 替代
try {
  const response = await someApi();
} catch (e) {
  // 手动处理，不走全局 toast
}
```

**方案 B：createClient 创建无 toast 的专用 client**

```typescript
const silentClient = createClient({
  baseURL: NEXT_PUBLIC_API_URL,
  // 不注册 toast 相关拦截器
});

const data = await someApi({ client: silentClient });
```

## 4. 与 Redocly 的配合

### 4.1 redocly.yaml 配置

```yaml
extends:
  - recommended          # 推荐规则集（生产环境推荐 recommended-strict）

apis:
  admin:
    root: ./docs/openapi/admin.json
  common:
    root: ./docs/openapi/common.json
  user:
    root: ./docs/openapi/user.json

rules:
  # 自定义规则覆盖
  no-ambiguous-paths: error
  operation-operationId: error
  path-segment-plural:
    severity: warn
```

### 4.2 OpenAPI 3.0.3 spec 的 lint 最佳实践

- **起步用 `recommended` 规则集**：覆盖命名规范、路径唯一性、operationId 等基础规则
- **CI 中用 `recommended-strict`**：将所有 warning 升级为 error，确保 CI 不遗漏问题
- **legacy 遗留问题**：使用 `redocly.yaml` 的 ignore file 机制标记已知且暂不修复的违规
- **多 API 配置**：每个 API 可以有独立的规则覆盖

### 4.3 CI/CD 集成

```yaml
# .github/workflows/api-lint.yml
name: API Lint
on: [pull_request]
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: 20
      - run: npx @redocly/cli lint --format=github-actions
```

`--format=github-actions` 会将 lint 结果直接标注到 PR 的代码变更行上。

### 4.4 与 @hey-api/openapi-ts 的配合流程

```
goctl-openapi → OpenAPI 3.0.3 JSON → redocly lint → @hey-api/openapi-ts → TypeScript SDK
```

Redocly 负责验证 spec 质量，hey-api 负责消费 spec 生成代码。两者通过 spec 文件解耦，互不依赖。

## 5. 大型项目中的使用案例

### 5.1 知名用户

@hey-api/openapi-ts（GitHub 4.4k+ stars）被以下知名项目/公司使用：
- **Vercel**（Guillermo Rauch 评价："OpenAPI codegen that just works"）
- **PayPal**
- **OpenCode**

### 5.2 Monorepo 中的使用模式

**模式 A：每个 app 独立配置（适合 perfect-panel）**

```
monorepo/
├── apps/
│   ├── admin/
│   │   ├── openapi-ts.config.ts     # admin + common spec
│   │   └── services/
│   │       ├── admin/               # admin spec 输出
│   │       └── common/              # common spec 输出
│   └── user/
│       ├── openapi-ts.config.ts     # user + common spec
│       └── services/
│           ├── user/
│           └── common/
├── docs/openapi/                    # 共享 spec 文件
└── package.json
```

**模式 B：共享 packages（NX/大型 monorepo）**

```
monorepo/
├── apps/
│   └── web/
├── packages/
│   ├── api-types/                   # @hey-api/typescript 输出
│   └── api-react-query/             # @tanstack/react-query 输出
└── openapi-ts.config.ts
```

GitHub Issue #1578 讨论了按 plugin 分目录输出的需求，但目前尚未原生支持。workaround 是运行多次生成，每次指定不同 plugins 和 output。

### 5.3 FastAPI + Next.js Monorepo 实践

Vinta Software 的博客案例展示了完整的 monorepo 工作流：
- OpenAPI spec 由后端自动导出
- `chokidar` 监控 spec 文件变化，自动重新生成 client
- 生成的 SDK 直接在 Next.js Server Component 中使用
- `fork-ts-checker-webpack-plugin` 提供实时类型检查反馈

## 6. 潜在陷阱和注意事项

### 6.1 版本兼容性

- **v0.91.0 起为 ESM-only**，不再支持 CommonJS。Bun 环境无影响
- **v0.86.0 起要求 Node.js 20.19+**。perfect-panel 已满足（engines: node >= 20）
- **v0.73.0 起 client 自动 bundled**，不再需要单独安装 `@hey-api/client-axios` 包。但如果要 `import { createClient }` 用于 SSR per-request client，仍需确认 bundled 版本是否导出此函数
- **v0.63.0 起 `client` 配置字段移至 `plugins`**：新项目直接在 plugins 中配置，无需顶层 `client` 字段

### 6.2 类型系统变化

- **全局 namespace → 具名导出**：`API.Ads` 变为 `import type { Ads } from './types.gen'`
  - 这是迁移体量最大的变更（200+ 处引用）
  - 可考虑生成兼容层：`export namespace API { export type { Ads, ... } }` 减少初始迁移量
- **参数类型命名**：`GetAdsListParams` 可能变为 `GetAdsListData`，包含 `{ query?: ..., path?: ..., body?: ... }` 结构
- **enum 处理**：推荐 `enums: 'javascript'`（普通对象 + as const），避免 TypeScript enum 的 tree-shaking 和类型兼容问题

### 6.3 SDK 函数签名变化

```typescript
// 旧模式（@umijs/openapi）
getAdsList(params: API.GetAdsListParams, options?: { [key: string]: any })

// 新模式（@hey-api/openapi-ts）
getAdsList(options?: { query?: GetAdsListData['query'] })
```

参数从多个位置参数变为统一的 options 对象。所有调用点需要适配。

### 6.4 返回值结构

@hey-api/client-axios 的 SDK 函数返回值取决于配置：
- 默认返回 `{ data, error, request, response }` 结构体（`responseStyle: 'fields'`）
- 配置 `throwOnError: true` 后，`data` 不再是 `T | undefined`，而是确定的 `T`（error 时抛出异常）
- **与当前 perfect-panel 的差异**：当前代码返回完整 AxiosResponse，通过 `response.data.data` 访问业务数据。迁移后可能变为 `result.data`（即 SDK 已 unwrap 一层 axios response）

**建议**：在正式迁移前，先生成一份样例 SDK 验证返回值结构。

### 6.5 多 spec 生成时的 client 初始化

每个 spec 输出目录会生成独立的 `client.gen.ts`，拥有独立的 client 实例。CSR 环境中，如果两个 spec（admin + common）的 client 都需要相同的拦截器，需要：

```typescript
// utils/setup-client.ts
import { client as adminClient } from '@/services/admin/client.gen';
import { client as commonClient } from '@/services/common/client.gen';

function setupInterceptors(client) {
  client.instance.interceptors.request.use(/* ... */);
  client.instance.interceptors.response.use(/* ... */);
}

setupInterceptors(adminClient);
setupInterceptors(commonClient);
```

或者使用 `runtimeConfigPath` 指向同一个配置文件，但需验证多个 spec 的 `runtimeConfigPath` 能否指向同一文件。

### 6.6 SSR 中的 Logout() 行为

当前 `Logout()` 可能操作客户端状态（cookies、路由重定向）。在 Server Component 中调用可能无效或抛异常。SSR client 的 40002-40005 处理应改为 throw error 而非调用 `Logout()`。

## 7. 对 perfect-panel 项目迁移的具体建议

### 7.1 推荐配置

**apps/admin/openapi-ts.config.ts**：

```typescript
import { defaultPlugins, defineConfig } from '@hey-api/openapi-ts';

const sharedPlugins = [
  ...defaultPlugins,
  {
    name: '@hey-api/client-axios',
    runtimeConfigPath: './utils/hey-api-config.ts',
  },
  {
    name: '@hey-api/typescript',
    enums: 'javascript',
  },
  {
    name: '@hey-api/sdk',
    auth: true,
    operations: { strategy: 'flat' },  // tree-shakeable 函数
  },
  {
    name: '@tanstack/react-query',
    queryOptions: true,
  },
];

export default [
  {
    input: '../../docs/openapi/admin.json',
    output: {
      path: './services/admin',
      postProcess: ['biome:format'],
    },
    plugins: sharedPlugins,
  },
  {
    input: '../../docs/openapi/common.json',
    output: {
      path: './services/common',
      postProcess: ['biome:format'],
    },
    plugins: sharedPlugins,
  },
];
```

**apps/admin/utils/hey-api-config.ts**：

```typescript
import type { CreateClientConfig } from '@/services/admin/client.gen';
import { NEXT_PUBLIC_API_URL, NEXT_PUBLIC_SITE_URL } from '@/config/constants';

export const createClientConfig: CreateClientConfig = (config) => ({
  ...config,
  baseURL: NEXT_PUBLIC_API_URL || NEXT_PUBLIC_SITE_URL,
});
```

### 7.2 拦截器迁移方案

**CSR 全局拦截器**（`utils/setup-clients.ts`）：

```typescript
import { client as adminClient } from '@/services/admin/client.gen';
import { client as commonClient } from '@/services/common/client.gen';
import { isBrowser } from '@workspace/ui/utils';
import { toast } from 'sonner';
import { getAuthorization, Logout } from './common';

function setupClient(client: typeof adminClient) {
  client.instance.interceptors.request.use(async (config) => {
    const Authorization = getAuthorization();
    if (Authorization) config.headers.Authorization = Authorization;
    return config;
  });

  client.instance.interceptors.response.use(
    async (response) => {
      const { code } = response.data;
      if (code !== 200) {
        if ([40002, 40003, 40004, 40005].includes(code)) return Logout();
        if (isBrowser()) toast.error(response.data.message);
        throw response;
      }
      return response;
    },
    async (error) => {
      if (isBrowser()) toast.error(error.message);
      return Promise.reject(error);
    },
  );
}

setupClient(adminClient);
setupClient(commonClient);
```

**SSR per-request client**（`utils/ssr-client.ts`）：

```typescript
import { createClient } from '@/services/admin/client.gen';
import { NEXT_PUBLIC_API_URL } from '@/config/constants';

export function createSSRClient(token: string) {
  return createClient({
    baseURL: NEXT_PUBLIC_API_URL,
    auth: () => token,
    // 不注册 toast 拦截器，错误直接抛出由 layout try-catch 处理
  });
}
```

### 7.3 package.json scripts 更新

```json
{
  "scripts": {
    "openapi": "openapi-ts"
  }
}
```

turbo 通过 `turbo run openapi` 并行执行每个 app 的生成任务。

### 7.4 redocly.yaml 配置建议

```yaml
extends:
  - recommended

apis:
  admin:
    root: ./docs/openapi/admin.json
  common:
    root: ./docs/openapi/common.json
  user:
    root: ./docs/openapi/user.json

rules:
  operation-operationId: error      # 确保每个 operation 有 ID
  no-ambiguous-paths: error         # 禁止歧义路径
  operation-summary: warn           # 建议添加摘要
  tag-description: warn             # 建议添加 tag 描述
```

### 7.5 迁移步骤建议

1. **Phase 1：生成器配置 + 首次 SDK 生成**
   - 安装 `@hey-api/openapi-ts`（devDependency）
   - 创建 `openapi-ts.config.ts`（admin + user 各一个）
   - 首次运行生成，验证输出结构和类型映射
   - 验证关键问题：返回值结构、`createClient` 是否可用、类型名对照

2. **Phase 2：拦截器迁移 + SSR client 分离**
   - 创建 `setup-clients.ts` 统一初始化
   - 创建 `ssr-client.ts` 用于 SSR 场景
   - 迁移 2 个 layout.tsx 的 SSR 调用
   - 验证 JWT 注入、40002-40005 登出、toast 行为

3. **Phase 3：47 个组件文件批量迁移**
   - `API.*` 类型 → 具名导入（200+ 处，可批量 find-replace）
   - `useQuery` → `xxxOptions()`（~35 处）
   - `useMutation` → `xxxMutation()`（~32 处）
   - ProTable `request` 参数适配（~29 处）
   - `skipErrorHandler` 适配（~8 处）

## 参考来源

- [Hey API 官方文档](https://heyapi.dev/)
- [Hey API Configuration](https://heyapi.dev/openapi-ts/configuration)
- [Hey API TanStack Query Plugin](https://heyapi.dev/openapi-ts/plugins/tanstack-query)
- [Hey API Axios Client](https://heyapi.dev/openapi-ts/clients/axios)
- [Hey API Next.js Client](https://heyapi.dev/openapi-ts/clients/next-js)
- [Hey API SDK Plugin](https://heyapi.dev/openapi-ts/plugins/sdk)
- [Hey API Output Configuration](https://heyapi.dev/openapi-ts/configuration/output)
- [Hey API Migration Guide](https://heyapi.dev/openapi-ts/migrating)
- [GitHub: hey-api/openapi-ts](https://github.com/hey-api/openapi-ts)
- [GitHub Issue #1056: Multiple OpenAPI Schemas](https://github.com/hey-api/openapi-ts/issues/1056)
- [GitHub Issue #1578: Multiple Output Directories](https://github.com/hey-api/openapi-ts/issues/1578)
- [GitHub Issue #914: Throw on Error](https://github.com/hey-api/openapi-ts/issues/914)
- [GitHub Issue #1703: Axios Interceptors Setup](https://github.com/hey-api/openapi-ts/issues/1703)
- [DeepWiki: hey-api/openapi-ts Architecture](https://deepwiki.com/hey-api/openapi-ts)
- [Vinta Software: FastAPI + Next.js Monorepo](https://www.vintasoftware.com/blog/nextjs-fastapi-monorepo)
- [Redocly CLI Configuration](https://redocly.com/docs/cli/configuration)
- [Redocly: Consistent APIs with GitHub Actions](https://redocly.com/blog/consistent-apis-redocly-github-actions)
- [TanStack Query Advanced SSR](https://tanstack.com/query/latest/docs/framework/react/guides/advanced-ssr)
