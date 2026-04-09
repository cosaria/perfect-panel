# Web 前端彻底 Vite 化设计

## 背景

当前 `web/` 已经以 Vite 作为实际运行入口，但仓库内部仍保留明显的 Next.js 遗留：

- `apps/admin` 与 `apps/user` 仍有大量 `app/` 目录源码
- 页面与组件仍广泛依赖 `next/link`、`next/navigation`、`next/image`
- `src/compat/*` 通过兼容层把 Next API 伪装成可在 Vite 中运行
- i18n、主题、环境变量命名仍带有明显的 Next 语义
- Turbo、测试、Docker、Go 注入链里仍残留 `.next`、`NEXT_PUBLIC_*` 等 Next 痕迹

这导致当前仓库处于“Vite 负责启动，但 Next 语义仍主导源码和配置”的混合状态，不符合“Vite 作为唯一源”的目标。

## 目标

本次改造的目标是把 `web/` 收敛为真正的纯 Vite React Monorepo：

- 删除所有 Next.js 运行时、构建链、兼容层与目录语义
- 让 `apps/admin` 与 `apps/user` 都只通过 Vite 构建与启动
- 把页面源码从 `app/` 目录迁移到常规 React 结构
- 把 `NEXT_PUBLIC_*` 全量迁移到 `VITE_*`
- 保持现有产品 URL、构建产物位置和 Go 嵌入发布链不变

## 非目标

以下内容不在本次改造范围内：

- 不修改后端 API 合同（contract）
- 不改变管理端 `/admin/*` 与用户端现有 URL 结构
- 不顺手重构与迁移无关的业务逻辑
- 不替换 Vite、Bun、Turbo 这一套现有前端基座

## 决策摘要

采用“全量重构，但不改产品外部行为”的方案：

1. 保留当前 Vite 作为唯一构建与运行入口
2. 删除所有 Next 运行时依赖、兼容层和专用配置
3. 把 `app/` 页面源码整体迁移到 `src/` 下的普通 React 目录
4. 将环境变量与服务端注入统一改为 `VITE_*`
5. 目标运行时明确为纯客户端 SPA，不保留 SSR / hydration 残留层
6. 通用基础设施优先收敛到共享层，避免 `admin` / `user` 各维护一套
7. 迁移按阶段推进，每个阶段都必须满足独立 DoD 后才能进入下一阶段
8. 保持 `dist` 产物、Go 嵌入链和线上访问路径不变

## 目标目录结构

每个前端应用都收敛到类似结构，但这是一组推荐落点，不是必须把所有目录都搭出来的模板：

```text
apps/<app>/
  index.html
  vite.config.ts
  tsconfig.json
  src/
    main.tsx
    router.tsx
    routes.tsx
    pages/
    components/
    ...
```

对应约束：

- `app/` 目录在迁移完成后整体删除
- `src/compat/` 目录在迁移完成后整体删除
- `src/main.tsx`、`src/router.tsx`、`src/routes.tsx` 是硬约束
- 页面懒加载入口必须从 `src/` 内的正式页面模块暴露
- 目录按需落位，只有在代码边界清晰时才新增 `layouts/`、`features/`、`providers/`、`lib/` 等目录
- 布局与 provider 不再通过 Next App Router 目录结构表达

## 运行时数据流

目标运行时是一条纯客户端链路，不再夹带任何 SSR / hydration 语义：

```text
Go Server
  |
  | inject window.__ENV (VITE_*)
  v
index.html / route html
  |
  v
Vite app bootstrap (src/main.tsx)
  |
  v
Router basename + route matching
  |
  v
Shared providers
  |-- i18n adapter
  |-- theme adapter
  |-- query client
  v
Lazy page chunks
  |
  v
API client + static assets + toast/theme behavior
```

## 核心设计

### 1. 路由

路由层统一只使用 `react-router-dom`。

具体规则：

- `src/router.tsx` 与 `src/routes.tsx` 作为唯一正式路由入口
- 页面组件统一通过 React Router 懒加载
- 所有 `next/navigation` 调用替换为：
  - `useNavigate`
  - `useLocation`
  - `useParams`
  - `useSearchParams`
- 所有 `next/link` 调用替换为 React Router 的 `Link` 或 `NavLink`
- 页面迁移后必须继续保留路由级懒加载，不允许因为目录迁移把所有页面改成首屏静态导入
- 目录迁移开始前，必须先加固 `admin` / `user` 的路由完整性测试，锁住预期路由图

迁移完成后，不再保留任何针对 Next Router 的路径别名和兼容函数。

### 2. 页面与布局迁移

现有 `app/**` 中的页面、布局和与页面强耦合的子组件整体迁移到 `src/` 下：

- 页面文件迁到 `src/pages/**`
- 共享布局按需迁到 `src/layouts/**`
- 与单个业务域高度相关的子组件按需迁到 `src/features/**`
- 通用组件继续保留在 `components/**` 或 `packages/ui`

迁移原则：

- 优先保留现有业务边界，不在迁移过程中做无关的领域重构
- 仅在目录和导入路径层面重组
- 同时清理仅为 Next App Router 服务的包装层
- `next/*` import 的替换可以先在原 `app/**` 位置完成，不要求必须先搬目录才能删 compat
- 不为了满足目录名而拆文件，先保留清晰边界，再决定是否抽目录

### 3. 图片与静态资源

删除 `next/image` 与 `next/legacy/image` 兼容层。

替换策略：

- 默认直接使用原生 `img`
- 如果多个页面确实需要统一的 `loading`、`decoding`、回退样式，可引入一个极薄的共享 `AppImage` 组件
- 不再保留模拟 Next 图片优化行为的 API
- 图片替换不能牺牲现有加载体验，至少要保留显式的 `loading`、`decoding` 与必要的尺寸约束

### 4. 国际化（i18n）

删除 `next-intl`，但不自研一套新的 i18n runtime，而是采用成熟 React i18n 库作为底座，再封一层极薄的本地适配。

约束如下：

- 继续复用现有 `locales/**` JSON 文件
- 共享层对外暴露统一接口，例如：
  - `I18nProvider`
  - `useLocale`
  - `useTranslations`
  - `setLocale`
- 保持当前 locale 持久化与 HTML `lang` / `dir` 同步能力
- 尽量兼容现有调用心智，减少页面改写噪音
- 共享 i18n 底座优先放到共享模块，`admin` / `user` 只保留应用特定配置和薄适配
- locale 替换必须带行为测试，至少验证文案刷新与 reload 后持久化

这样可以保留现有翻译资产，同时避免把“删 Next”扩大成“顺手重写整套翻译系统”。

### 5. 主题

删除 `next-themes`，但不引入新的主题框架，只保留一层极薄的共享浏览器主题适配。

最小职责：

- 读取和写入主题偏好
- 处理系统主题跟随
- 在 `document.documentElement` 上同步主题 class
- 对外暴露 `useTheme` / `setTheme`

`packages/ui` 中依赖 `useTheme` 的组件同步改为依赖共享主题上下文，不再从 `next-themes` 取值。
主题替换必须带行为测试，至少验证 DOM class 与 toast 主题联动。

### 6. 环境变量与运行时注入

环境变量命名从 `NEXT_PUBLIC_*` 全量迁移到 `VITE_*`。

统一读取优先级：

1. `window.__ENV.VITE_*`，用于 Go 嵌入后的运行时注入
2. `import.meta.env.VITE_*`，用于本地开发和静态构建

这意味着以下链路必须同步改造：

- `apps/admin/config/constants.ts`
- `apps/user/config/constants.ts`
- `vite.config.ts` 中的公开变量注入策略
- `server/cmd/server_service.go` 中的前端环境变量注入
- `server/web/*` 中校验注入内容的测试
- 依赖这些常量的前端测试和文档
- 运行时 basename、API baseUrl 和默认 locale 相关 smoke 校验

建议映射关系如下：

| 旧变量 | 新变量 |
| --- | --- |
| `NEXT_PUBLIC_API_URL` | `VITE_API_URL` |
| `NEXT_PUBLIC_SITE_URL` | `VITE_SITE_URL` |
| `NEXT_PUBLIC_DEFAULT_LANGUAGE` | `VITE_DEFAULT_LANGUAGE` |
| `NEXT_PUBLIC_ADMIN_PATH` | `VITE_ADMIN_PATH` |
| `NEXT_PUBLIC_DEFAULT_USER_EMAIL` | `VITE_DEFAULT_USER_EMAIL` |
| `NEXT_PUBLIC_DEFAULT_USER_PASSWORD` | `VITE_DEFAULT_USER_PASSWORD` |
| `NEXT_PUBLIC_CDN_URL` | `VITE_CDN_URL` |
| `NEXT_PUBLIC_EMAIL` | `VITE_EMAIL` |
| `NEXT_PUBLIC_TELEGRAM_LINK` | `VITE_TELEGRAM_LINK` |
| `NEXT_PUBLIC_DISCORD_LINK` | `VITE_DISCORD_LINK` |
| `NEXT_PUBLIC_GITHUB_LINK` | `VITE_GITHUB_LINK` |
| `NEXT_PUBLIC_LINKEDIN_LINK` | `VITE_LINKEDIN_LINK` |
| `NEXT_PUBLIC_TWITTER_LINK` | `VITE_TWITTER_LINK` |
| `NEXT_PUBLIC_INSTAGRAM_LINK` | `VITE_INSTAGRAM_LINK` |
| `NEXT_PUBLIC_HOME_USER_COUNT` | `VITE_HOME_USER_COUNT` |
| `NEXT_PUBLIC_HOME_SERVER_COUNT` | `VITE_HOME_SERVER_COUNT` |
| `NEXT_PUBLIC_HOME_LOCATION_COUNT` | `VITE_HOME_LOCATION_COUNT` |

## 构建链与仓库配置调整

为保证 Vite 成为唯一源，需要同时清理以下残留：

- 删除 `next`、`next-intl`、`next-themes`、`@tanstack/react-query-next-experimental`、`@netlify/plugin-nextjs` 等依赖
- 删除 `packages/typescript-config/nextjs.json`
- 删除 `tsconfig` 与 `vite.config` 中针对 `next/*` 的别名映射
- 更新 `packages/ui/tailwind.config.ts` 及应用侧 re-export 配置，让 Tailwind content globs 覆盖 `src/pages/**`、`src/features/**`、`src/layouts/**`、`src/components/**` 等新落点，不再只依赖 `app/**` 与 `components/**`
- 清理 `web/package.json`、`web/turbo.json` 中 `.next`、`out` 等输出声明
- 清理 `web/tests` 中“同时保留 Vite 与 Next 输出”的历史断言
- 清理文档中对 Next 构建和 `NEXT_PUBLIC_*` 的描述

迁移完成后，`web/` 应只剩 Vite 相关入口、脚本、输出和校验。

## 共享层边界

这次迁移的 DRY 原则必须写清楚，不然实现时很容易让 `admin` / `user` 各自复制一套新代码。

优先进入共享层的内容：

- i18n 底座与薄适配接口
- 主题底座与 `useTheme`
- 图片薄包装
- env 读取工具与公共常量解析逻辑

允许保留在应用内的内容：

- `admin` / `user` 各自的 locale 资源装配
- 各自的默认路由、basename、业务特定 provider glue
- 各自的业务页面与业务组件

判断标准只有一个：如果两端行为与接口一致，就优先共享；如果只是为了“看起来对称”而抽象，则不要新增共享层。

## 服务端嵌入链设计

虽然前端内部全面去 Next，但发布链保持不变：

- `apps/admin/dist` 继续复制到 `server/web/admin-dist`
- `apps/user/dist` 继续复制到 `server/web/user-dist`
- Dockerfile 继续从前端构建阶段复制 `dist` 结果
- Go server 继续负责：
  - 注入 `window.__ENV`
  - 返回静态资源
  - 保留 `/admin` 与用户端现有路径入口

需要变化的不只是注入变量名，还包括 Go 静态层的历史 Next 语义必须一起迁走：

- `window.__ENV` 注入内容与测试断言从 `NEXT_PUBLIC_*` 改为 `VITE_*`
- `server/web/static.go` 中 `_next` 相关路径与缓存假设改为纯 Vite 产物语义
- admin HTML base path 重写逻辑不再依赖 `_next` 路径约定
- 相关 Go 测试必须一并更新，避免“前端是 Vite，静态层还是半个 Next”的状态

这条嵌入链是本次迁移的关键回归面，不能只改前端，不改 Go 静态层。

### Admin 自定义路径与静态资产策略

管理端保留“运行时 `adminPath` 可变”的要求，因此必须把静态资产 URL 策略写死，避免实现时各自猜测。

- admin 构建继续以 `/admin/` 作为编译期基准产出绝对资产路径
- Go 在返回 admin HTML 时，将 HTML 中的 `/admin/...` 绝对链接统一重写到运行时 `adminPath`
- 重写范围只覆盖 HTML 中的链接与资源引用，不修改 Vite 生成的哈希文件名和磁盘产物结构
- 回归验证至少覆盖：
  - 自定义 admin path 下首页可打开
  - 自定义 admin path 下嵌套路由刷新可打开
  - 自定义 admin path 下至少一个 `assets/*.js` 与一个 `assets/*.css` 请求成功
  - 返回给浏览器的 HTML 中不再泄露 `/admin/assets/` 绝对路径

这样可以继续兼容 `/manage` 等运行时路径，又不把这次迁移扩大成另一轮静态资源路径机制重写。

## 风险与约束

本次改造的主要风险点如下：

1. 页面目录大迁移导致导入路径大面积变化，容易出现漏改
2. i18n / theme 适配层替换后，若行为测试不足，容易出现无报错但用户可感知的静默回归
3. `VITE_*` 改名会同时波及前端、Go 注入、测试、文档与开发环境变量
4. Go 嵌入静态层若保留 `_next` 假设，会让缓存和路径规则落在半迁移状态
5. 若把业务逻辑重构混入迁移，会显著放大回归面

因此执行时必须遵守以下约束：

- 先迁基础设施，再迁页面
- 只做与“纯 Vite 化”直接相关的重构
- 每一步都要配套类型检查、构建和关键测试验证
- 每一阶段未满足 DoD 前，不进入下一阶段

## 验收标准

迁移完成后，应满足以下条件：

- `web/` 内不再存在 `next`、`next-*` 运行时依赖
- `web/apps/admin/app` 与 `web/apps/user/app` 目录被删除
- `web/apps/*/src/compat` 目录被删除
- 前端源码不再引用 `next/link`、`next/navigation`、`next/image`、`next-intl`、`next-themes`
- 前端运行时环境变量不再引用 `NEXT_PUBLIC_*`
- `bun run lint`、`bun run typecheck`、`bun run build` 通过
- 根级 `make embed` 成功
- 服务端与前端中校验静态资源注入、路由与构建链的测试更新后通过
- 管理端 `/admin/*` 与用户端既有路由继续可访问
- 至少有最小浏览器 smoke 覆盖：
  - `/admin`
  - 自定义 admin path 下的嵌套路由刷新
  - `/auth`
  - `/dashboard`
- locale 替换后验证文案刷新与持久化
- theme 替换后验证 DOM class 与 toast 主题联动
- 路由完整性测试在页面迁移前后都通过
- 路由级懒加载仍然存在，未把大量页面卷入首屏入口
- `src/` 目录下的新页面与组件样式仍被 Tailwind 扫描，不因 content globs 只盯 `app/**` 而丢样式
- Go 静态层缓存规则已切换到 Vite 资产语义并有回归测试
- 本次 PR 内补上前端 bundle / 静态资产缓存的持续回归监控，并把它接入现有 `bun:test` + `go test` + `repo-contracts` gate，不把它留成迁移后的 TODO

## 分阶段实施与阶段 DoD

迁移按阶段推进。每个阶段都必须在独立构建、关键测试和 smoke 满足后，才能进入下一阶段。

### 阶段 1：建立 `VITE_*` 与 Go 静态层新基线

内容：

- 前端常量与 env 读取逻辑切到 `VITE_*`
- Go 注入与 `server/web` 相关测试切到 `VITE_*`
- 清理 `_next` 缓存与路径语义，改为 Vite 资产规则

DoD：

- `window.__ENV` 注入测试通过
- `client baseUrl` 与 `api-base` 测试通过
- Go 静态层测试通过，且不再依赖 `_next` 语义
- 自定义 admin path smoke 不回退
- 用户端 `/auth` 首屏、HTML fallback 与 `/api` / `/v1` 不被 SPA fallback 吞掉的测试通过

阻断条件：

- 任一 admin path / env 注入测试失败
- Go 静态缓存规则仍保留 Next 专用假设
- 用户端静态路由回退规则被破坏，导致页面兜底吞掉 API 404 或公共路由不可达

### 阶段 2：替换共享基础设施

内容：

- 用成熟 i18n 库 + 薄适配替换 `next-intl`
- 用共享主题薄适配替换 `next-themes`
- 视需要引入共享 `AppImage`

DoD：

- locale 行为测试通过，文案刷新与持久化有效
- theme 行为测试通过，DOM class 与 toast 主题联动有效
- 不引入新的应用内双份基础设施实现

阻断条件：

- locale / theme 行为出现静默回归
- `admin` / `user` 各自复制一套共享底座实现

### 阶段 3：清理 compat 层并固定纯客户端 SPA 运行时

内容：

- 清理 `next/*` compat 层与别名
- 删除 `ReactQueryStreamedHydration` 等 hydration 残留
- 明确纯客户端 SPA 启动链
- 在仍位于 `app/**` 的页面中先完成 `next/*` import 替换，确保 compat 清理与目录迁移解耦

DoD：

- `admin` / `user` 都以纯客户端方式启动
- 不再存在 `next/*` 路径别名
- 即使 `app/**` 页面尚未全部迁移，源码也不再依赖 compat alias 才能编译
- 最小浏览器 smoke 可通过

阻断条件：

- 仍有 hydration / SSR 兼容层残留
- 浏览器 smoke 失败

### 阶段 4：先加固路由测试，再迁移页面目录

内容：

- 先更新 `admin` / `user` 路由完整性测试
- 再做 `app/ -> src/` 页面与布局迁移
- 同步更新 Tailwind content globs，确保新 `src/**` 页面与组件继续参与样式生成

DoD：

- 路由完整性测试在迁移前后都通过
- 既有导航路径无丢失
- 页面仍保留路由级懒加载
- `src/pages`、`src/layouts`、`src/features` 等新目录中的类名能在构建产物中正确出样式

阻断条件：

- 路由完整性测试未加固
- 迁移后页面被卷入首屏静态入口
- Tailwind content 仍只扫描 `app/**` / `components/**`，导致迁移后出现静默样式缺失

### 阶段 5：构建链、文档、监控与性能收口

内容：

- 清理依赖、配置、测试与文档残留
- 增加 bundle / 静态资产缓存回归监控
- 收口性能 DoD

持续回归监控的落点必须具体，而不是抽象描述：

- 前端侧复用现有 `web/tests/*.test.ts` 的 `bun:test` 合同测试，至少覆盖：
  - `admin-build-chain`
  - `user-build-chain`
  - 关键路由图与 admin path 合同
  - 与 `VITE_*`、静态资产 URL、Tailwind 扫描边界相关的构建断言
- 服务端侧复用 `server/web/static*.go` 的 `go test`，至少覆盖：
  - `window.__ENV` 注入
  - user API 404 / SPA fallback 边界
  - admin 自定义路径 HTML 重写
  - immutable asset cache 规则
- CI 接入现有 `repo:contracts` / `repo-contracts` workflow gate，保证这些回归断言不是“本地手动跑一次”

DoD：

- `bun run lint`、`bun run typecheck`、`bun run build`、`make embed` 全通过
- build chain 测试只验证 Vite 输出
- 已建立本次 PR 内的 bundle / 缓存回归检查
- chunk 分裂与缓存规则未出现明显回退

阻断条件：

- 仍保留 `.next` / `out` 输出语义
- 未建立回归监控
- 构建通过但首屏 chunk 或缓存规则明显退步

这样可以先把底层约束固定，再推进高噪音的目录迁移，降低回归定位成本。

## 性能要求

这次不是性能专项重构，但不能接受明显的性能回退。

必须显式检查：

- 路由级懒加载仍然存在
- 首屏入口没有因为目录迁移变成大而全静态包
- Vite 产物仍保持合理的 chunk 分裂
- 静态哈希资产命中 immutable cache
- admin / user 首屏在迁移后没有因基础设施替换明显增重
- 持续回归检查结果会进入现有 CI gate，而不是停留在一次性人工对比

这里不要求复杂 benchmark 平台，但必须在本次 PR 中补上最小可持续回归检查，避免迁移后性能慢慢长回去。

## 结论

本次方案不是在现有 Vite 外壳上继续兼容 Next，而是把 `web/` 完整收敛为：

- Vite 作为唯一构建源
- React Router 作为唯一前端路由源
- 共享层中的成熟 i18n 底座 + 薄适配，以及共享主题薄适配，作为唯一运行时上下文源
- `VITE_*` 作为唯一公开环境变量命名源

同时保留现有 URL、`dist` 产物位置与 Go 嵌入发布模式，确保对外行为稳定、对内架构彻底收口。
