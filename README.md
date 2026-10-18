# Doolox CMS

Doolox is a simple CMS and static website builder for webmasters. It runs a local
web server for **live preview**, lives in your **system tray**, and **generates a
static website** you can deploy anywhere — no database, no build pipeline, just
files on disk.

> This file is a hands-on tutorial for using Doolox. To install or build the
> program itself, see [INSTALL.md](INSTALL.md).

---

## Table of contents

1. [How Doolox works](#how-doolox-works)
2. [Quick start](#quick-start)
3. [The data directory](#the-data-directory)
4. [Configuration](#configuration)
5. [Core concepts](#core-concepts)
6. [Creating your first page](#creating-your-first-page)
7. [Page front matter](#page-front-matter)
8. [Layouts](#layouts)
9. [Components (partials)](#components-partials)
10. [Template variables & helpers](#template-variables--helpers)
11. [Assets](#assets)
12. [Multiple websites](#multiple-websites)
13. [Previewing and generating the static site](#previewing-and-generating-the-static-site)
14. [Deploying the generated site](#deploying-the-generated-site)
15. [Reference](#reference)

---

## How Doolox works

Doolox reads your website from plain files on disk and serves it locally so you
can edit and see changes immediately. When you are happy, you ask Doolox to
"generate": it renders every page into a static HTML file, copies your assets,
and writes a `sitemap.xml`. The result is a folder of static files you can
upload to any web host.

There is no admin panel for editing content — you edit files directly in your
favourite editor. The bundled **Doolox Dashboard** is used to trigger site
generation and to check the result.

While running, Doolox exposes two local addresses:

| URL                     | What it is                                                         |
| ----------------------- | ------------------------------------------------------------------ |
| `http://localhost:5000` | Your website, served **live** (the preview).                       |
| `http://localhost:5555` | The **Doolox Dashboard** (AdminLTE interface with a *Generate* button). |

A **system tray icon** gives you the same controls without a browser:

- **Generate** — generate the static site.
- **Publish** — generate the static site and push it with git.
- **Doolox Dashboard** — open the dashboard.
- **Quit** — stop Doolox.

Clicking the tray icon opens your live website.

---

## Quick start

Prerequisites: **Go 1.26+** (see [INSTALL.md](INSTALL.md) for OS-specific notes).

```sh
# from the repository root
make build        # -> bin/doolox
./bin/doolox      # or: make run
```

Doolox starts, opens your site in the browser, and drops into the system tray.
On Linux the dashboard opens in your system browser; on Windows/macOS it opens
in a native web view.

Stop it from the tray menu (**Quit**), or with `Ctrl+C` when running from a
terminal.

---

## The data directory

On first launch Doolox creates a `doolox` folder in your home directory and
works entirely out of it:

```
~/doolox/
├── config.json          # your configuration
├── doolox.log           # log output
├── templates/           # reserved (currently unused)
└── websites/            # one subfolder per website
    └── append/          # the default website ("append")
        ├── src/         # the website source (your content)
        │   ├── pages/           # one .html file per page
        │   ├── layout.html      # layout for inner pages
        │   ├── layout_home.html # layout for the home page
        │   ├── components/      # reusable partials (header, menu, footer)
        │   └── assets/          # css, js, images, vendor files
        └── ...          # generated .html files + assets/ after generation
```

You can override the base directory with the `DOOLOX_HOME` environment
variable (mainly useful for testing).

---

## Configuration

Everything is configured in `~/doolox/config.json`:

```json
{
    "theme": "simple",
    "url": "https://example.com",
    "website": "append"
}
```

| Key       | Purpose                                                                                                    |
| --------- | ---------------------------------------------------------------------------------------------------------- |
| `website` | Name of the **main** website (its folder under `websites/`). Served at the site root `/`. Defaults to `append`. |
| `url`     | The public URL of your site. Used to build `sitemap.xml`. Required for generation to write the sitemap.    |
| `theme`   | Reserved for future use — it is loaded but not currently used by the renderer.                             |

---

## Core concepts

A **website** is a folder under `~/doolox/websites/`. It uses a simple source
layout (the modern `src/` directory):

```
websites/<name>/src/
├── pages/           # .html files — each file is one page
├── layout.html      # wraps every inner page
├── layout_home.html # wraps the home page only
├── components/      # partials you include with {{template}}
└── assets/          # static files, copied verbatim on generation
```

Doolox also still supports the **legacy** layout, where pages live in
`content/pages/` and assets in `content/theme/assets/`. If a `src/` directory
exists it is used; otherwise the `content/` directory is used.

The pieces fit together like this:

- **Pages** hold your content (`pages/*.html`).
- **Layouts** provide the surrounding HTML shell and inject the page content
  where the `{{ yield }}` placeholder appears.
- **Components** are small, reusable templates (`header`, `main_menu`, `footer`).

---

## Creating your first page

A page is a single `.html` file in the `pages/` directory. The **URL of the
page is its file name without the extension**:

| File               | URL        |
| ------------------ | ---------- |
| `pages/index.html` | `/` (home) |
| `pages/about.html` | `/about`   |
| `pages/pricing.html` | `/pricing` |

> Pages are **flat**: only files directly inside `pages/` are discovered.
> `index.html` is the home page (a legacy `home.html` is also accepted and maps
> to `/`; `index.html` wins if both exist).

Each page is a Go `html/template` definition named after its URL. The minimum
page looks like this:

`websites/append/src/pages/about.html`

```html
<!--
title: About Us
subtitle: About
type: page
-->
{{define "about"}}
<section>
  <h2>About Us</h2>
  <p>We build websites people love.</p>
</section>
{{end}}
```

After saving, visit `http://localhost:5000/about` — the page is rendered inside
the `layout.html` shell immediately.

---

## Page front matter

The optional comment at the very top of a page file sets the page metadata:

```html
<!--
title: About Us
subtitle: About
description: A short summary used by search engines.
type: page
-->
```

| Key           | Stored on            | Purpose                                         |
| ------------- | -------------------- | ----------------------------------------------- |
| `title`       | `.Page.Title`        | Page title (used in `<title>`, headings).       |
| `subtitle`    | `.Page.Subtitle`     | Short secondary label (breadcrumbs, headings).  |
| `description` | `.Page.Description`  | Meta description for SEO.                       |
| `type`        | `.Page.Type`         | `page` or `post` (constants `PAGE` / `POST`).   |

All keys are optional — a page with no front matter simply renders without
metadata.

---

## Layouts

Layouts wrap page content with the HTML shell (doctype, `<head>`, header, menu,
footer). The page content is injected where the layout contains `{{ yield }}`.

There are two layouts:

| File               | Used for                              |
| ------------------ | ------------------------------------- |
| `layout_home.html` | The home page (`index.html`, URL `/`) |
| `layout.html`      | Every other page                      |

Example inner layout:

```html
<!DOCTYPE html>
<html lang="en">
{{template "header" .}}
<body>
  {{template "main_menu" .}}
  <main>
    <h1>{{ .Page.Title }}</h1>
    {{ yield }}
  </main>
  {{template "footer" .}}
</body>
</html>
```

To give the home page a different hero/header, override `layout_home.html`
instead — Doolox picks the correct layout automatically.

---

## Components (partials)

Components live in `components/` and are defined with `{{define "name"}}` and
included with `{{template "name" .}}`:

`components/header.html`

```html
{{define "header"}}
<head>
  <meta charset="utf-8">
  <title>{{ .Page.Title }} | My Site</title>
</head>
{{end}}
```

`components/main_menu.html`

```html
{{define "main_menu"}}
<nav>
  <a href="./"{{if .Page.IsActive "/"}} class="active"{{end}}>Home</a>
  <a href="about"{{if .Page.IsActive "/about"}} class="active"{{end}}>About</a>
</nav>
{{end}}
```

Because the renderer compiles `.tmpl` and `.html` files, components can live in
any file; the bundled theme keeps one component per file for clarity.

---

## Template variables & helpers

Templates are Go `html/template` templates (via Macaron). Inside a page or
component you have access to:

| Expression                        | Description                                                        |
| --------------------------------- | ------------------------------------------------------------------ |
| `{{ .Page.Title }}`               | Page title from front matter.                                      |
| `{{ .Page.Subtitle }}`            | Page subtitle from front matter.                                   |
| `{{ .Page.Description }}`         | Page description from front matter.                                |
| `{{ .Page.Type }}`                | `page` or `post`.                                                  |
| `{{ .Page.URL }}`                 | The page slug (e.g. `about`).                                      |
| `{{ .Page.Lastmod }}`             | File modification time (used in the sitemap).                      |
| `{{ .Page.IsActive "/about" }}`   | `true` when the link points at this page — handy for menus.        |
| `{{ yield }}`                     | In a layout: where the page content is inserted.                   |
| `{{ obfuscate "me@example.com" }}`| Renders an email with a hidden span to deter scrapers.             |

Standard Go template syntax (`{{if}}`, `{{range}}`, `{{with}}`, `{{$var :=}}`)
works throughout.

---

## Assets

Everything under `src/assets/` is **copied verbatim** into the website's
`assets/` output folder when you generate, and is served live at `/assets/...`.

Reference them with a root-relative path in your templates and pages:

```html
<link href="assets/css/main.css" rel="stylesheet">
<img src="assets/img/logo.png" alt="Logo">
```

The bundled `append` theme keeps a `vendor/` folder inside `assets/` for
Bootstrap, icons, and other libraries.

---

## Multiple websites

The **main** website (from `config.json` → `website`) is served at `/`. Any
**subfolder** of the main website that contains its own `src/` (or `content/`)
source layout becomes an additional website, served under its own sub-path:

```
websites/append/            # main site  -> /
├── src/ ...
└── cms/                    # nested site -> /cms
    └── src/ ...
```

So a page `cms/src/pages/about.html` is reachable at `/cms/about`, and its
assets live under `/cms/assets/...`. Generation renders every website into its
own folder and writes a single combined `sitemap.xml` into the main website's
folder.

This is how the bundled theme ships a second, simpler site at `/cms`.

---

## Previewing and generating the static site

While Doolox is running, the live preview at `http://localhost:5000` renders
pages on the fly. The preview also injects a tiny helper script so you can
regenerate without leaving the page.

**Generate the static site** in any of these ways:

1. **Dashboard** — open `http://localhost:5555` and click *Generate website*
   (or the sidebar link). You'll be redirected back with a success/error
   message.
2. **Tray menu** — choose *Generate*.
3. **Keyboard shortcut** — while viewing the live site, press `Ctrl+G`
   (`Cmd+G` on macOS).
4. **Endpoint** — visit or `curl http://localhost:5000/generate`.

Generation:

- renders every page of every website into a static `.html` file inside the
  website's own folder (next to `src/`),
- copies each website's `assets/` into the output,
- writes one combined `sitemap.xml` into the main website's folder (so it ships
  with the deployable output),
- shows a native system notification with the result (e.g. *"Site generated: 4 pages"*).

Before writing output, each website's previously generated files are cleared.
The source (`src/` / `content/`), `config.json`, `.git`, nested website folders,
`CNAME`, and `README` files are always preserved.

> The sitemap needs `url` set in `config.json`; without it generation reports
> *"no url set in the configuration"*.

### Publishing with git

If your website's folder is a git repository, you can generate **and** publish
in one step. While viewing the live site, press `Ctrl+P` (`Cmd+P` on macOS):

1. Doolox generates the static files (exactly like `Ctrl+G`), then
2. runs `git commit -a -m "Doolox update."` and `git push` inside the
   website's folder (`~/doolox/websites/<website>/`).

A native system notification reports the result. The same action is available
from the tray menu (*Publish*) or at `http://localhost:5000/publish`.

---

## Deploying the generated site

Generation writes static files directly into each website's folder:

```
~/doolox/websites/append/
├── src/           # your source (kept)
├── index.html     # generated
├── about.html     # generated
├── pricing.html   # generated
└── assets/        # copied from src/assets
```

Upload those generated `.html` files and the `assets/` folder (plus a `CNAME`
file if you use a custom domain) to any static host — GitHub Pages, Netlify,
Cloudflare Pages, S3, or a plain web server. Doolox is a static site generator,
so the output needs no runtime.

---

## Reference

### Ports

| Port   | Service                  |
| ------ | ------------------------ |
| `5000` | Live website preview     |
| `5555` | Doolox Dashboard         |

### HTTP routes (live preview)

| Route              | Purpose                                                  |
| ------------------ | -------------------------------------------------------- |
| `/`                | Home page of the main website.                           |
| `/generate`        | Generate the static site (returns plain text).           |
| `/publish`         | Generate, then commit and push the site with git.        |
| `/doolox.js`       | Live-preview helper (Ctrl+G / Ctrl+P shortcuts).         |
| `/assets/*`        | Assets of the main website.                              |
| `/:seg`            | Page of the main website, or home of a nested website.   |
| `/:seg/:page`      | Page of the website named by the first segment.          |
| `/:seg/assets/*`   | Assets of the website named by the first segment.        |

### Environment variables

| Variable      | Effect                                          |
| ------------- | ----------------------------------------------- |
| `DOOLOX_HOME` | Overrides the base data directory (`~/doolox`). |

### Make targets

```sh
make run        # run for the current OS
make build      # build the binary -> bin/doolox
make test       # run the tests
make vet        # run go vet
make fmt        # format all Go sources
make tidy       # tidy module dependencies
make build-deb  # build a .deb package for Debian/Ubuntu/Linux Mint
make clean      # remove bin/
make help       # list all targets
```

### Command-line flag

| Flag            | Default | Effect                                                                 |
| --------------- | ------- | ---------------------------------------------------------------------- |
| `-open-website` | `true`  | Open the live website in the browser on startup (`-open-website=false` to disable). |

---

## Further reading

- **[INSTALL.md](INSTALL.md)** — building and installing Doolox on Linux,
  macOS, and Windows, plus `.deb` and Windows installer packaging.
