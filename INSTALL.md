# Installing and building Doolox

Doolox is a Go program. The admin interface opens in a **native web view**
(webview/webview_go), which links through cgo on every platform and requires
the matching native toolkit (GTK/WebKitGTK on Linux, WebView2 on Windows,
WKWebView on macOS).

## Prerequisites

| OS        | Requirements                                                                                                        |
|-----------|----------------------------------------------------------------------------------------------------------------------|
| Linux     | Go 1.26+, gcc/g++, GTK 3 + WebKitGTK (`sudo apt install build-essential libgtk-3-dev libwebkit2gtk-4.1-dev`)         |
| macOS     | Go 1.26+ and Xcode Command Line Tools (`xcode-select --install`)                                                     |
| Windows   | Go 1.26+ and MinGW-w64 gcc/g++ (https://www.mingw-w64.org/)                                                          |

## Linux

```sh
make build          # -> bin/doolox
./bin/doolox
```

## Debian / Ubuntu / Linux Mint (.deb package)

Build the installer (requires `dpkg-deb` and the GTK/WebKitGTK build deps above):

```sh
make build-deb       # -> dist/doolox_<version>_amd64.deb
# or
./build-deb.sh
```

Install it:

```sh
sudo apt install ./dist/doolox_*.deb
```

This installs the binary to `/opt/doolox`, a launcher to `/usr/bin/doolox`,
a menu entry and icons, plus the bundled templates and config. On first launch
a `doolox` folder is created in your home directory, containing `config.json`,
plus a `websites/` subfolder with `content/` and an empty `docs/` folder, and
an empty `templates/` folder. Logs are written to `~/doolox/doolox.log`.
Run `doolox` from a terminal or launch it from the application menu.

Remove it with:

```sh
sudo apt remove doolox
```

(User data in `~/doolox` is left untouched on removal.)

## macOS

The macOS build must happen **on a Mac**, because the menu-bar icon depends
on Objective-C (cgo) and cannot be cross-compiled from Linux or Windows.

```sh
make build          # -> bin/doolox
# or
make build-darwin   # -> bin/doolox-darwin

./bin/doolox
```

## Windows

Option A — PowerShell (no `make` required):

```powershell
.\build.ps1         # -> bin\doolox-windows-amd64.exe
```

Option B — Git Bash / MSYS2 with `make`:

```sh
make build         # -> bin/doolox.exe
```

Option C — Inno Setup installer (requires Inno Setup 6+):

```powershell
.\build.ps1               # first build the Windows binary
ISCC.exe installer.iss    # -> dist\doolox-setup-<version>.exe
```

## Cross-compilation

The native web view links through cgo, so cross-compiling is **not supported**
— build natively on each OS. macOS additionally cannot be cross-compiled
(Objective-C + cgo); build it on a Mac.

## Other Makefile targets

```sh
make run        # run for the current OS
make test       # run the tests
make vet        # run go vet
make fmt        # format all Go sources
make tidy       # tidy module dependencies
make clean      # remove bin/
make install    # go install the binary
make build-deb  # build a .deb package for Debian/Ubuntu/Linux Mint
make help       # list all targets
```

## Notes

- `go.mod` declares `go 1.26.5`.
- `build.ps1` passes `-H=windowsgui` so the app does not open
  a console window when launched.
- Pages are plain files in `content/pages/*.html`. Each page's metadata
  (`title`, `subtitle`, `description`, `type`) lives in an optional
  front-matter comment at the top of the file:

  ```html
  <!--
  title: About Doolox
  subtitle: About
  type: page
  -->
  {{define "about"}} ... {{end}}
  ```

  The page URL is the file name without the extension (`about.html` ->
  `/about`); `index.html` is the home page.
- The website name is set in `config.json` (`"website"`) and defaults to
  `doolox`; it names the folder under `~/doolox/websites/`.
- The tray icon assets (`assets/trayicon.svg`, `assets/trayicon.png`,
  `assets/trayicon.ico`) are embedded into the binary at build time with
  `go:embed`, so no extra runtime files are needed.
- The admin web interface assets (`admin/theme/`) are shipped next to the
  binary by both installers (`/opt/doolox/admin/theme` on Linux,
  `C:\Program Files\Doolox\admin\theme` on Windows) and served by the admin
  server at http://localhost:5555. The HTML templates and logo are embedded
  with `go:embed`; only the CSS/JS/favicon are read from disk at runtime.
- The native web view uses `webview/webview_go`, vendored under
  `third_party/webview_go` with a one-line patch to build against WebKitGTK
  4.1 (Ubuntu 24.04+ / Debian 13); see the `replace` directive in `go.mod`.
- The app version is configured in `version.iss` — the single source of truth
  used by both the Inno Setup (`installer.iss`) and Debian (`build-deb.sh`)
  packaging. Bump it there and rebuild.
