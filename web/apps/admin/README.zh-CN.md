<a name="readme-top"></a>

<div align="center">

<img width="160" src="https://raw.githubusercontent.com/perfect-panel/ppanel-assets/refs/heads/main/logo.svg">

<h1>PPanel 管理后台</h1>

这是由 PPanel 提供支持的 PPanel 管理后台

[英文](./README.md)
·
中文
·
[更新日志](../../CHANGELOG.md)
·
[报告问题][issues-link]
·
[请求功能][issues-link]

<!-- SHIELD GROUP -->

[![][github-release-shield]][github-release-link]
[![][github-releasedate-shield]][github-releasedate-link]
[![][github-action-test-shield]][github-action-test-link]
[![][github-action-release-shield]][github-action-release-link]<br/>
[![][github-contributors-shield]][github-contributors-link]
[![][github-forks-shield]][github-forks-link]
[![][github-stars-shield]][github-stars-link]
[![][github-issues-shield]][github-issues-link]
[![][github-license-shield]][github-license-link]

![](https://urlscan.io/liveshot/?width=1920&height=1080&url=https://admin.ppanel.dev)

</div>

<details>
<summary><kbd>目录</kbd></summary>

#### 目录

- [⌨️ 本地开发](#️-本地开发)
- [📦 构建与嵌入](#-构建与嵌入)
- [🤝 贡献](#-贡献)
- [📝 许可证](#-许可证)

####

</details>

## ⌨️ 本地开发

管理端现在以 `Vite` 单页应用（SPA）的形式运行在 monorepo 中。

先克隆仓库并安装依赖：

```bash
git clone https://github.com/cosaria/perfect-panel.git
cd perfect-panel

# 安装仓库依赖
make bootstrap

# 仅启动管理端前端
cd web/apps/admin
bun run dev
```

在浏览器中打开 <http://localhost:3000> 查看结果。

如果需要同时启动 Go 服务和管理端前端，使用仓库根命令：

```bash
make dev APP=admin
```

## 📦 构建与嵌入

官方发布链会把管理端静态资源嵌入 Go 服务：

```bash
# 构建 admin 到 apps/admin/dist，并复制到 server/web/admin-dist
make embed-admin

# 构建带嵌入前端的最终 Go 二进制
make build-all
```

容器构建以仓库根 `Dockerfile` 为准。

## 🤝 贡献

欢迎各种类型的贡献，
如果您有兴趣贡献代码，请随时查看我们的 GitHub
[问题][github-issues-link] 来展示您的能力。

[![][pr-welcome-shield]][pr-welcome-link]

[![][contributors-contrib]][contributors-url]

<div align="right">

[![][back-to-top]](#readme-top)

</div>

---

## 📝 许可证

版权所有 © 2024 [PPanel][profile-link]。<br />
本项目使用 [GNU](./LICENSE) 许可证。

<!-- LINK GROUP -->

[back-to-top]: https://img.shields.io/badge/-BACK_TO_TOP-151515?style=flat-square
[contributors-contrib]: https://contrib.rocks/image?repo=perfect-panel/ppanel-web
[contributors-url]: https://github.com/perfect-panel/ppanel-web/graphs/contributors
[github-action-release-link]: https://github.com/perfect-panel/ppanel-web/actions/workflows/release.yml
[github-action-release-shield]: https://img.shields.io/github/actions/workflow/status/perfect-panel/ppanel-web/release.yml?label=release&labelColor=black&logo=githubactions&logoColor=white&style=flat-square
[github-action-test-link]: https://github.com/perfect-panel/ppanel-web/actions/workflows/test.yml
[github-action-test-shield]: https://img.shields.io/github/actions/workflow/status/perfect-panel/ppanel-web/test.yml?label=test&labelColor=black&logo=githubactions&logoColor=white&style=flat-square
[github-contributors-link]: https://github.com/perfect-panel/ppanel-web/graphs/contributors
[github-contributors-shield]: https://img.shields.io/github/contributors/perfect-panel/ppanel-web?color=c4f042&labelColor=black&style=flat-square
[github-forks-link]: https://github.com/perfect-panel/ppanel-web/network/members
[github-forks-shield]: https://img.shields.io/github/forks/perfect-panel/ppanel-web?color=8ae8ff&labelColor=black&style=flat-square
[github-issues-link]: https://github.com/perfect-panel/ppanel-web/issues
[github-issues-shield]: https://img.shields.io/github/issues/perfect-panel/ppanel-web?color=ff80eb&labelColor=black&style=flat-square
[github-license-link]: https://github.com/perfect-panel/ppanel-web/blob/master/LICENSE
[github-license-shield]: https://img.shields.io/github/license/perfect-panel/ppanel-web?color=white&labelColor=black&style=flat-square
[github-release-link]: https://github.com/perfect-panel/ppanel-web/releases
[github-release-shield]: https://img.shields.io/github/v/release/perfect-panel/ppanel-web?style=flat-square&sort=semver&logo=github
[github-releasedate-link]: https://github.com/perfect-panel/ppanel-web/releases
[github-releasedate-shield]: https://img.shields.io/github/release-date/perfect-panel/ppanel-web?labelColor=black&style=flat-square
[github-stars-link]: https://github.com/perfect-panel/ppanel-web/network/stargazers
[github-stars-shield]: https://img.shields.io/github/stars/perfect-panel/ppanel-web?color=ffcb47&labelColor=black&style=flat-square
[issues-link]: https://github.com/perfect-panel/ppanel-web/issues/new/choose
[pr-welcome-link]: https://github.com/perfect-panel/ppanel-web/pulls
[pr-welcome-shield]: https://img.shields.io/badge/🤯_pr_welcome-%E2%86%92-ffcb47?labelColor=black&style=for-the-badge
[profile-link]: https://github.com/perfect-panel
