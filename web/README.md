# PPanel Web

PPanel Web 是 perfect-panel monorepo 里的前端工作区，包含两个应用：

- `apps/admin`: 管理端
- `apps/user`: 用户端

如果你是第一次进入这个仓库，先从仓库根目录的 `README.md` 开始。根 README 是唯一 onboarding 入口。

## Frontend-Only Commands

当你只需要处理前端工作区时，可以在 `web/` 目录下执行：

```bash
bun install
bun run lint
bun run build
```

按应用单独开发时：

```bash
bun run dev --filter=ppanel-admin-web
bun run dev --filter=ppanel-user-web
```

## Applications

- `ppanel-admin-web`: Next.js 管理后台
- `ppanel-user-web`: Next.js 用户端

## Monorepo Policy

- 仓库级开发入口以根目录 `README.md` 和根级 `Makefile` 为准
- 这个文件只描述前端子项目，不再承担整仓 onboarding
- 问题反馈和仓库元数据统一归到 `cosaria/perfect-panel`

## Related Files

- `../README.md`: 仓库入口
- `./package.json`: 前端工作区入口
- `./apps/admin/package.json`: 管理端脚本
- `./apps/user/package.json`: 用户端脚本
