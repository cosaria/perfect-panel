<a name="readme-top"></a>

<div align="center">

<img width="160" src="https://raw.githubusercontent.com/perfect-panel/ppanel-assets/refs/heads/main/logo.svg">

<h1>PPanel admin web</h1>

This is a PPanel admin web powered by PPanel

English
·
[Chinese](./README.zh-CN.md)
·
[Changelog](../../CHANGELOG.md)
·
[Report Bug][issues-link]
·
[Request Feature][issues-link]

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
<summary><kbd>Table of contents</kbd></summary>

#### TOC

- [⌨️ Local Development](#️-local-development)
- [📦 Build and Embed](#-build-and-embed)
- [🤝 Contributing](#-contributing)
- [📝 License](#-license)

####

</details>

## ⌨️ Local Development

This admin app now runs as a Vite-powered SPA inside the monorepo.

Clone the repository and install dependencies:

```bash
git clone https://github.com/cosaria/perfect-panel.git
cd perfect-panel

# Install repo dependencies
make bootstrap

# Run the admin frontend only
cd web/apps/admin
bun run dev
```

Open <http://localhost:3000> with your browser to see the result.

If you want the Go server and admin frontend together, use the repo root command:

```bash
make dev APP=admin
```

## 📦 Build and Embed

The canonical release path embeds the admin static build into the Go server:

```bash
# Build admin to apps/admin/dist and copy it into server/web/admin-dist
make embed-admin

# Build the final Go binary with embedded frontends
make build-all
```

For container builds, the repo root `Dockerfile` is the source of truth.

## 🤝 Contributing

Contributions of all types are more than welcome,
if you're interested in contributing code, feel free to check out our GitHub
[Issues][github-issues-link] to get stuck in to show us what you’re made of.

[![][pr-welcome-shield]][pr-welcome-link]

[![][contributors-contrib]][contributors-url]

<div align="right">

[![][back-to-top]](#readme-top)

</div>

---

## 📝 License

Copyright © 2024 [PPanel][profile-link]. <br />
This project is [GNU](../../LICENSE) licensed.

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
