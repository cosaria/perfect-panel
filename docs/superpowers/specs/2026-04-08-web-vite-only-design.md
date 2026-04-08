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
5. 保持 `dist` 产物、Go 嵌入链和线上访问路径不变

## 目标目录结构

每个前端应用都收敛到类似结构：

```text
apps/<app>/
  index.html
  vite.config.ts
  tsconfig.json
  src/
    main.tsx
    router.tsx
    routes.tsx
    layouts/
    pages/
    features/
    components/
    providers/
    i18n/
    theme/
    lib/
```

对应约束：

- `app/` 目录在迁移完成后整体删除
- `src/compat/` 目录在迁移完成后整体删除
- 页面懒加载入口只允许从 `src/pages/**` 暴露
- 布局与 provider 不再通过 Next App Router 目录结构表达

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

迁移完成后，不再保留任何针对 Next Router 的路径别名和兼容函数。

### 2. 页面与布局迁移

现有 `app/**` 中的页面、布局和与页面强耦合的子组件整体迁移到 `src/` 下：

- 页面文件迁到 `src/pages/**`
- 共享布局迁到 `src/layouts/**`
- 与单个业务域高度相关的子组件迁到 `src/features/**`
- 通用组件继续保留在 `components/**` 或 `packages/ui`

迁移原则：

- 优先保留现有业务边界，不在迁移过程中做无关的领域重构
- 仅在目录和导入路径层面重组
- 同时清理仅为 Next App Router 服务的包装层

### 3. 图片与静态资源

删除 `next/image` 与 `next/legacy/image` 兼容层。

替换策略：

- 默认直接使用原生 `img`
- 如果多个页面确实需要统一的 `loading`、`decoding`、回退样式，可引入一个极薄的本地 `AppImage` 组件
- 不再保留模拟 Next 图片优化行为的 API

### 4. 国际化（i18n）

删除 `next-intl`，改为本地轻量 i18n 实现。

约束如下：

- 继续复用现有 `locales/**` JSON 文件
- 在各应用内提供本地 i18n 模块，例如：
  - `I18nProvider`
  - `useLocale`
  - `useTranslations`
  - `setLocale`
- 保持当前 locale 持久化与 HTML `lang` / `dir` 同步能力
- 尽量兼容现有调用心智，减少页面改写噪音

这样可以保留现有翻译资产，同时彻底去掉对 Next runtime 的依赖。

### 5. 主题

删除 `next-themes`，改为本地主题 provider。

最小职责：

- 读取和写入主题偏好
- 处理系统主题跟随
- 在 `document.documentElement` 上同步主题 class
- 对外暴露 `useTheme` / `setTheme`

`packages/ui` 中依赖 `useTheme` 的组件同步改为依赖本地主题上下文，不再从 `next-themes` 取值。

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
- 清理 `web/package.json`、`web/turbo.json` 中 `.next`、`out` 等输出声明
- 清理 `web/tests` 中“同时保留 Vite 与 Next 输出”的历史断言
- 清理文档中对 Next 构建和 `NEXT_PUBLIC_*` 的描述

迁移完成后，`web/` 应只剩 Vite 相关入口、脚本、输出和校验。

## 服务端嵌入链设计

虽然前端内部全面去 Next，但发布链保持不变：

- `apps/admin/dist` 继续复制到 `server/web/admin-dist`
- `apps/user/dist` 继续复制到 `server/web/user-dist`
- Dockerfile 继续从前端构建阶段复制 `dist` 结果
- Go server 继续负责：
  - 注入 `window.__ENV`
  - 返回静态资源
  - 保留 `/admin` 与用户端现有路径入口

需要变化的是注入内容与测试断言从 `NEXT_PUBLIC_*` 改为 `VITE_*`，而不是整体推翻嵌入模式。

## 风险与约束

本次改造的主要风险点如下：

1. 页面目录大迁移导致导入路径大面积变化，容易出现漏改
2. `next-intl` 与 `next-themes` 替换后，运行时上下文若不兼容，可能影响大量页面
3. `NEXT_PUBLIC_*` 改名会同时波及前端、Go 注入、测试、文档与开发环境变量
4. 若把业务逻辑重构混入迁移，会显著放大回归面

因此执行时必须遵守以下约束：

- 先迁基础设施，再迁页面
- 只做与“纯 Vite 化”直接相关的重构
- 每一步都要配套类型检查、构建和关键测试验证

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

## 实施顺序

推荐按以下顺序执行：

1. 替换 env 读取层与 Go 注入层，建立 `VITE_*` 新基线
2. 替换主题与 i18n 基础设施
3. 清理 Next compat 层与路由 API，改为纯 React Router
4. 把 `app/` 页面与布局迁移到 `src/` 新结构
5. 清理依赖、配置、测试与文档残留
6. 完整跑通类型检查、构建、嵌入与关键测试

这样可以先把底层约束固定，再推进高噪音的目录迁移，降低回归定位成本。

## 结论

本次方案不是在现有 Vite 外壳上继续兼容 Next，而是把 `web/` 完整收敛为：

- Vite 作为唯一构建源
- React Router 作为唯一前端路由源
- 本地 i18n 与主题 provider 作为唯一运行时上下文源
- `VITE_*` 作为唯一公开环境变量命名源

同时保留现有 URL、`dist` 产物位置与 Go 嵌入发布模式，确保对外行为稳定、对内架构彻底收口。
