<div align="center">

# 0xViva.dev

**A terminal-flavoured personal site that ships as one tiny Go binary.**

[![Build](https://img.shields.io/github/actions/workflow/status/0xViva/0xViva.dev/deploy.yaml?branch=master&label=build&logo=github)](https://github.com/0xViva/0xViva.dev/actions/workflows/deploy.yaml)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![templ](https://img.shields.io/badge/templ-v0.3.1100-ff6b35)](https://templ.guide)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-v4-06B6D4?logo=tailwindcss&logoColor=white)](https://tailwindcss.com)
[![License](https://img.shields.io/badge/license-MIT-22c55e)](https://opensource.org/licenses/MIT)
[![Live](https://img.shields.io/badge/live-0xviva.dev-8b5cf6?logo=firefox&logoColor=white)](https://0xviva.dev)

</div>

## 🖥️ At a Glance

A terminal-inspired homepage for **August Justinus Gran**, built with Go and server-rendered templ components. It pulls a live GitHub commit feed with HTMX and keeps the frontend lean enough to run on almost anything.

> **No SPA. No JS framework. No runtime dependencies.** Just Go, HTML, and a handful of HTMX attributes.

## ⚡ Core Features

| Feature | What it does |
| --- | --- |
| 🖥️ **Terminal-native UI** | CRT grids, scanlines, a live UTC clock, and a keyboard-friendly layout. |
| 📡 **Live commit feed** | Fetches recent repos, branches, languages, and commits from the GitHub API on load. |
| 🪶 **Tiny production image** | Multi-stage build produces a static Go binary on a distroless base. |
| 🔁 **Rootless deploys** | Runs as a Podman Quadlet, with CSS cache-busting by commit SHA. |

## 🧰 Tech Stack

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Echo](https://img.shields.io/badge/Echo-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![templ](https://img.shields.io/badge/templ-ff6b35?style=for-the-badge&logo=go&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-06B6D4?style=for-the-badge&logo=tailwindcss&logoColor=white)
![HTMX](https://img.shields.io/badge/HTMX-3366CC?style=for-the-badge&logo=htmx&logoColor=white)
![Podman](https://img.shields.io/badge/Podman-892CA0?style=for-the-badge&logo=podman&logoColor=white)

## 🚀 Quick Start

> **Prerequisites:** Go 1.26+, [Task](https://taskfile.dev), [templ](https://templ.guide), [Tailwind CSS v4](https://tailwindcss.com), and [Air](https://github.com/air-verse/air).

```sh
# 1. Clone the repo
git clone https://github.com/0xViva/0xViva.dev && cd 0xViva.dev

# 2. Add a read-only GitHub token for the commit feed
echo "GITHUB_TOKEN=ghp_your_token" > .env

# 3. Run with hot reload on http://localhost:8080
task dev
```

Build a production container with:

```sh
task container image_name=0xviva/0xviva.dev image_tag=local
```

## 📄 License & Contact

**MIT** © [August Justinus Gran](https://github.com/0xViva) · [GitHub](https://github.com/0xViva) · [LinkedIn](https://www.linkedin.com/in/august-justinus-gran) · [0xViva@proton.me](mailto:0xViva@proton.me)

<details>
<summary>🔖 Steal my badge</summary>

```html
<script src="https://raw.githubusercontent.com/0xViva/0xViva.dev/master/assets/js/augustg-dev-badge.js"></script>
<augustg-dev-badge></augustg-dev-badge>
```

</details>
