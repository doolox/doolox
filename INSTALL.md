# Installing and building Doolox

Doolox is a Go program. The admin interface opens in a **native web view** on
Windows and macOS (webview/webview_go). On Linux the dashboard opens in the
system browser instead, because embedding WebKitGTK in a Go process is
unreliable.

## Prerequisites

| OS        | Requirements                                                                                                        |
|-----------|----------------------------------------------------------------------------------------------------------------------|
| Linux     | Go 1.26+                                                                                                           |
| macOS     | Go 1.26+ and Xcode Command Line Tools (`xcode-select --install`)                                                     |
| Windows   | Go 1.26+ and MinGW-w64 gcc/g++ (https://www.mingw-w64.org/)                                                          |

## Linux

```sh
make build          # -> bin/doolox
./bin/doolox
```

## Debian / Ubuntu / Linux Mint (.deb package)

Build the installer (requires `dpkg-deb`):

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
plus a `websites/` subfolder with the bundled sites: `simple/` (the default,
served at `/`) and `append/`.
Each website's generated HTML pages and assets are written inside its own
folder, and a single combined `sitemap.xml` is written into the main
website's folder. An empty `templates/` folder is also created. Logs are
written to
`~/doolox/doolox.log`.
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

### macOS disk image (.dmg)

Build a `Doolox.app` bundle and a `.dmg` disk image:

```sh
make build          # -> bin/doolox (native architecture)
./build-dmg.sh      # -> dist/Doolox-<version>.dmg
```

The `.dmg` contains `Doolox.app` — a menu-bar-only app (it hides the Dock
icon, since Doolox lives in the system tray) — and a link to `/Applications`.
Drag the app into `Applications` to install it.

Because the app is unsigned, macOS Gatekeeper blocks it on first launch:
right-click the app in Finder and choose **Open**, or clear the quarantine
flag with:

```sh
xattr -dr com.apple.quarantine /Applications/Doolox.app
```

A signed/notarized build requires an Apple Developer certificate and running
`codesign` + `xcrun notarytool`; that is not automated yet.

### GitHub Actions

The `release` workflow builds the `.dmg` on a `macos-latest` (Apple Silicon)
runner. It compiles both the arm64 and amd64 slices and merges them into a
universal binary with `lipo`, then packages the app with `build-dmg.sh`.

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

The Windows and macOS builds link a native web view through cgo, so those
cannot be cross-compiled — build them natively. macOS additionally cannot be
cross-compiled (Objective-C + cgo).

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
- The main website is set in `config.json` (`"website"`, default `simple`) and
  is served at the site root; any other website is a subfolder of it and is
  served at `/<subfolder>/`. Each website's `url` and `theme` live in that
  website's own `config.json`.
- The tray icon assets (`assets/trayicon.svg`, `assets/trayicon.png`,
  `assets/trayicon.ico`) are embedded into the binary at build time with
  `go:embed`, so no extra runtime files are needed.
- The admin web interface assets (`admin/theme/`) are shipped next to the
  binary by both installers (`/opt/doolox/admin/theme` on Linux,
  `C:\Program Files\Doolox\admin\theme` on Windows) and served by the admin
  server at http://localhost:5555. The HTML templates and logo are embedded
  with `go:embed`; only the CSS/JS/favicon are read from disk at runtime.
- The native web view uses `webview/webview_go` (vendored under
  `third_party/webview_go`; see the `replace` directive in `go.mod`). It is
  only used on Windows and macOS — on Linux the dashboard opens in the system
  browser, because WebKitGTK's JavaScriptCore conflicts with the Go runtime
  over signal handling and the web view cannot be re-created after closing.
- The app version is configured in `version.iss` — the single source of truth
  used by both the Inno Setup (`installer.iss`) and Debian (`build-deb.sh`)
  packaging. Bump it there and rebuild.
