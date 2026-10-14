# Installing and building Doolox

Doolox is a Go program that links **SQLite through cgo**
(`github.com/mattn/go-sqlite3`), so a **C compiler is required on every
platform**. The menu-bar icon also uses cgo/Objective-C on macOS.

## Prerequisites

| OS        | Requirements                                                                          |
|-----------|----------------------------------------------------------------------------------------|
| Linux     | Go 1.26+ and gcc (`sudo apt install build-essential`)                                  |
| macOS     | Go 1.26+ and Xcode Command Line Tools (`xcode-select --install`)                       |
| Windows   | Go 1.26+ and MinGW-w64 gcc (https://www.mingw-w64.org/)                                |

## Linux

```sh
make build          # -> bin/doolox
./bin/doolox
```

## Debian / Ubuntu / Linux Mint (.deb package)

Build the installer (requires `dpkg-deb` and a C compiler for cgo):

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
a menu entry and icons, plus the bundled templates/database. On first launch a
`doolox` folder is created in your home directory, with a `websites/` subfolder
containing `content/`, an empty `docs/` folder, `config.json` and `doolox.db`,
plus an empty `templates/` folder. Logs are written to `~/doolox/doolox.log`.
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

Option B — Git Bash / MSYS2 with `make` and MinGW-w64 gcc:

```sh
make build-windows  # -> bin/doolox-windows-amd64.exe
```

Option C — Inno Setup installer (requires Inno Setup 6+):

```powershell
.\build.ps1               # first build the Windows binary
ISCC.exe installer.iss    # -> dist\doolox-setup-<version>.exe
```

## Cross-compilation (optional)

Cross-compiling requires the matching **cross C compiler** (because of cgo):

```sh
# on Linux/macOS, build a Windows binary (needs x86_64-w64-mingw32-gcc):
make build-windows CC=x86_64-w64-mingw32-gcc

# on macOS/Windows, build a Linux binary (needs x86_64-linux-gnu-gcc):
make build-linux CC=x86_64-linux-gnu-gcc
```

macOS cannot be cross-compiled from other OSes (Objective-C + cgo); build it
on a Mac.

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
- `build-windows` / `build.ps1` pass `-H=windowsgui` so the app does not open
  a console window when launched.
- All platforms build with cgo enabled (`CGO_ENABLED=1`) because of SQLite.
- The tray icon assets (`assets/trayicon.svg`, `assets/trayicon.png`,
  `assets/trayicon.ico`) are embedded into the binary at build time with
  `go:embed`, so no extra runtime files are needed.
- The app version is configured in `version.iss` — the single source of truth
  used by both the Inno Setup (`installer.iss`) and Debian (`build-deb.sh`)
  packaging. Bump it there and rebuild.
