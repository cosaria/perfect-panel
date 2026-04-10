# Web-V2 Admin 设计规范

## 背景

`web-v2/` 被定义为一个全新的前端 monorepo。  
它不复用现有 [web](/Users/admin/Codes/ProxyCode/perfect-panel/web) 的目录组织、运行时约束或生成链路，也不以兼容旧 admin 实现为目标。

这份规范首先服务两个目标：

1. 为 `server-v2` 提供一套可以稳定消费 `OpenAPI` 合同的新 admin 前端
2. 为后续追加 `user` 端保留清晰、可扩展的 monorepo 骨架

`web-v2` 的第一阶段只考虑 `admin` 调用面，不讨论 `user` 端迁移与并存策略。

## 目标

1. 定义 `web-v2/` 的标准 monorepo 结构与技术栈。
2. 定义 `admin` 第一阶段的页面范围、认证流和 UI 骨架。
3. 定义 `server-v2` `OpenAPI` 到 `web-v2` SDK 的单一路径。
4. 为后续 `web-v2` 实施计划提供稳定前提。

## 非目标

1. 本文不重写现有 [web](/Users/admin/Codes/ProxyCode/perfect-panel/web)。
2. 本文不实现完整 admin 业务模块，只覆盖第一阶段壳子与认证页范围。
3. 本文不定义 `user` 端信息架构和页面实现。
4. 本文不替代 `server-v2` 的后端真相源，不在前端复制后端业务逻辑。

## 核心判断

`web-v2` 采用**独立 monorepo + TanStack Start React 应用（底层开发与构建基于 Vite）+ shadcn/ui 共享 UI 包**的方案。

不选“继续扩展现有 `web/`”的原因：

- 旧前端已经带有自身历史结构和工具链约束
- `web-v2` 需要以新项目方式建立更稳定的目录和合同边界
- admin 第一阶段只需要一套更干净的底盘，不需要兼容全部旧页面

不选“纯 Vite SPA 手工拼装 TanStack Router”的原因：

- `TanStack Start` 能自然承接 `TanStack Router` 生态
- 后续如需引入更强的路由、loader、server-side 能力，不需要再迁框架
- 仍然保持 Vite 生态的开发体验

## 技术栈

- 应用基座：`TanStack Start`（React，底层开发与构建基于 `Vite`）
- 包管理：`pnpm workspace`
- 任务编排：`turbo`
- 样式：`Tailwind CSS 4`
- UI：`shadcn/ui`
- 路由：`@tanstack/react-router`
- 服务端状态：`@tanstack/react-query`
- 表单：`@tanstack/react-form`
- 表格：`@tanstack/react-table`
- 本地薄状态：`@tanstack/store`
- 合同 SDK：`@hey-api/openapi-ts`
- 单元测试：`vitest`
- 代码规范：`biome`

## 推荐目录结构

```text
web-v2/
├── apps/
│   └── admin/
│       ├── src/
│       │   ├── app/
│       │   ├── pages/
│       │   ├── features/
│       │   ├── widgets/
│       │   └── shared/
│       ├── public/
│       ├── app.config.ts
│       ├── vite.config.ts
│       ├── tsconfig.json
│       └── package.json
├── packages/
│   ├── ui/
│   ├── api-client/
│   ├── typescript-config/
│   └── biome-config/
├── package.json
├── pnpm-workspace.yaml
├── turbo.json
└── biome.json
```

## 顶层目录职责

- `apps/admin`
  - `web-v2` 第一阶段唯一前端应用
- `packages/ui`
  - `shadcn/ui` 基础组件、主题 token、共享布局原子
- `packages/api-client`
  - 由 `@hey-api/openapi-ts` 生成的共享 SDK 与类型
- `packages/typescript-config`
  - monorepo 统一 TypeScript 配置
- `packages/biome-config`
  - monorepo 统一 `Biome` 规则

## Admin 应用结构

`apps/admin/src` 按应用职责分层：

- `app/`
  - 应用入口、providers、router、layout shell
- `pages/`
  - 路由页面壳，只做页面级装配
- `features/`
  - 登录、忘记密码、重置密码、Dashboard、用户管理、系统设置等业务模块
- `widgets/`
  - sidebar、topbar、状态卡片、空状态等跨页面组合
- `shared/`
  - 通用 hooks、store、constants、utils、auth guard

### 分层约束

- 页面不得直接手写 fetch 逻辑，统一经过 `packages/api-client`
- `shared/store` 只承载本地薄状态，不承载服务端真相
- 业务表单 schema 与提交逻辑放入对应 `features/*`
- `widgets/` 只负责跨页面组合，不持有领域业务规则

## UI 与设计系统

### 基础约束

- 主布局以 [shadcn/ui `blocks/sidebar`](https://ui.shadcn.com/blocks/sidebar) 为基础
- 认证页以 [shadcn/ui `blocks/login`](https://ui.shadcn.com/blocks/login) 为基础
- `shadcn/ui` 组件统一落在 `packages/ui`
- admin 页面级组合组件仍落在 `apps/admin`

### 主题策略

- 使用 `shadcn/ui` 官方 CSS 变量体系
- 第一阶段默认主打浅色 admin 主题
- 结构上保留暗色模式切换能力
- 不允许直接使用“默认模板色板”而不做项目级语义 token 命名

## 第一阶段页面范围

### 认证页

第一阶段认证页固定为：

- 登录页
- 忘记密码页
- 重置密码页

第一阶段**不实现 admin 自注册页**。

原因：

- 管理后台不应默认开放自注册
- `server-v2` 当前也未定义完整注册主链

### 登录后页面

第一阶段只实现 3 个菜单：

- `Dashboard`
- `用户管理`
- `系统设置`

### Layout 结构

- 左侧：sidebar navigation
- 顶部：轻量 topbar
- 主体：内容区
- 移动端：sidebar 折叠为抽屉

topbar 第一阶段承载：

- 页面标题
- 搜索占位
- 主题切换
- 当前管理员菜单

## 第一阶段信息架构

### Dashboard

`Dashboard` 第一阶段采用**系统概览型**首页，而不是业务 KPI 首页。

展示内容固定为：

- 当前管理员身份信息
- `server-v2` API 连通状态
- 当前环境 / 版本信息
- 进入 `用户管理`、`系统设置` 的快捷入口

### 用户管理

第一阶段只做结构化页面骨架：

- 页面标题
- 搜索框占位
- “新建用户”按钮占位
- 表格区域骨架
- 空状态说明

不在第一阶段实现完整用户 CRUD。

### 系统设置

第一阶段只做结构化页面骨架：

- 页面标题与说明
- `站点设置`
- `认证设置`

以上仅提供卡片分组骨架，不实现真实保存逻辑。

## 认证与状态流

### 认证来源

`web-v2/admin` 直接调用 `server-v2` 已存在的真实 auth 接口：

- `POST /api/v1/public/sessions`
- `POST /api/v1/public/password-reset-requests`
- `POST /api/v1/public/password-resets`

### 登录态策略

- `accessToken` 存在 `TanStack Store`
- 登录态同步持久化到 `localStorage`
- 应用启动时恢复最小认证壳状态
- API client 自动注入 `Authorization: Bearer <token>`
- 路由层统一使用 `RequireAuth`
- 收到 `401` 时统一清理登录态并跳回登录页

### Store 承重边界

`@tanstack/store` 在第一阶段按**薄状态、非承重**能力使用，只承担：

- 当前登录状态
- 当前管理员最小展示信息
- sidebar 展开/折叠状态
- 主题模式

`@tanstack/store` 不承担：

- 远端查询缓存
- 业务实体真相
- 分页列表数据真相

这些职责仍归 `TanStack Query`。

同时，`@tanstack/store` 当前不作为远端数据真相容器，也不承接复杂业务状态迁移。

## TanStack 生态职责划分

- `TanStack Router`
  - 路由树、layout 嵌套、受保护页面、URL 状态
- `TanStack Query`
  - 远端数据获取、缓存、失效、请求状态
- `TanStack Form`
  - 登录、忘记密码、重置密码等表单状态与校验
- `TanStack Table`
  - 用户管理页表格骨架
- `TanStack Store`
  - 薄壳状态

## OpenAPI 与 SDK 生成链

### 单一路径

`server-v2` 仍然是 `OpenAPI` 真相源，但 SDK 生成动作放在 `web-v2` 一侧执行。

固定流水线：

1. `server-v2` 维护 `openapi/openapi.yaml`
2. `server-v2` 产出可消费的 bundle / json 产物（例如 `openapi/dist/openapi.json`）
3. `web-v2` 读取该 bundle 作为输入
4. `web-v2/packages/api-client` 内运行 `@hey-api/openapi-ts`
5. `apps/admin` 只消费 `packages/api-client`

### 禁止事项

- 不允许在 `apps/admin` 里手写第二套 API 类型
- 不允许在 `apps/admin` 内部各页面自己生成 SDK
- 不允许把生成动作塞回 `server-v2`，再把生成产物当作前端源码直接提交到别处

## 工程命令约定

第一阶段至少提供：

- `pnpm install`
- `pnpm dev`
- `pnpm build`
- `pnpm test`
- `pnpm lint`
- `pnpm format`
- `pnpm typecheck`
- `pnpm openapi`

其中：

- `pnpm openapi`
  - 在 `web-v2` 内执行 `@hey-api/openapi-ts`
  - 输入为 `server-v2` 导出的 `OpenAPI` 产物

## 第一阶段质量门禁

第一阶段至少要求：

1. `pnpm lint`
2. `pnpm typecheck`
3. `pnpm test`
4. `pnpm build`
5. `pnpm openapi`

## 明确禁止事项

- 不把旧 [web](/Users/admin/Codes/ProxyCode/perfect-panel/web) 的 admin 页面整包迁入 `web-v2`
- 不把 `server-v2` 逻辑复制到前端 `server functions`
- 不在第一阶段提前挂出大量空菜单
- 不在 `apps/admin` 内直接存放 `shadcn/ui` 基础组件副本
- 不让 `TanStack Store` 和 `TanStack Query` 争夺远端数据真相

## 允许的弹性

- 第一阶段只有 `apps/admin`，后续允许新增 `apps/user`
- `packages/ui` 内部组件按需创建，不要求预建所有组件
- `Dashboard / 用户管理 / 系统设置` 在第一阶段允许先以骨架页存在

## 决策结果

`web-v2` 第一阶段采用：

- 独立 monorepo
- `TanStack Start` React 应用
- `pnpm workspace + turbo`
- `Tailwind CSS 4 + shadcn/ui`
- `TanStack Router / Query / Form / Table / Store`
- `@hey-api/openapi-ts` 在 `web-v2` 一侧生成 SDK
- `admin` 先行，范围限定为认证页、admin shell、Dashboard、用户管理骨架、系统设置骨架

这份规范确立的是：

- `web-v2` 怎么组织
- `admin` 第一阶段怎么收口
- 前后端合同怎么衔接

它不要求第一阶段就把 admin 全业务模块做完。
