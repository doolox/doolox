# Doolox build/development Makefile.
#
# The app links SQLite through cgo (github.com/mattn/go-sqlite3), so a C
# compiler is required on every platform:
#   Linux   -> gcc            (e.g. apt install build-essential)
#   macOS   -> Xcode CLT      (xcode-select --install)
#   Windows -> MinGW-w64 gcc  (see INSTALL.md)
#
# For native Windows builds without make, see build.ps1.

GO      ?= go
APP     := doolox
BIN_DIR := bin
PKG     := .
LDFLAGS := -s -w
GOFLAGS := -trimpath

# Cross C compilers for the cross-compilation targets. Override with
# `make build-windows CC=...` if your cross toolchain has a different name.
CC_LINUX   ?= x86_64-linux-gnu-gcc
CC_WINDOWS ?= x86_64-w64-mingw32-gcc

.PHONY: all build run test vet fmt tidy clean install \
        build-linux build-windows build-darwin build-deb help

all: build ## Build the binary (default target)

build: ## Build for the current OS
	@mkdir -p $(BIN_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APP) $(PKG)

run: ## Run for the current OS
	$(GO) run $(PKG)

test: ## Run tests
	$(GO) test ./...

vet: ## Run go vet
	$(GO) vet ./...

fmt: ## Format all Go sources
	$(GO) fmt ./...

tidy: ## Tidy module dependencies
	$(GO) mod tidy

clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)

install: ## Install the binary with go install
	$(GO) install $(PKG)

build-linux: ## Cross-compile for linux/amd64 (needs a Linux cross C compiler)
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 CC=$(CC_LINUX) GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APP)-linux-amd64 $(PKG)

build-windows: ## Cross-compile for windows/amd64 (needs a MinGW cross C compiler)
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 CC=$(CC_WINDOWS) GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS) -H=windowsgui" -o $(BIN_DIR)/$(APP)-windows-amd64.exe $(PKG)

build-darwin: ## Build for macOS (must be run on macOS)
	@test "$$(go env GOOS)" = "darwin" || { echo "error: build-darwin requires macOS (tray icon uses cgo/Objective-C)"; exit 1; }
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APP)-darwin $(PKG)

build-deb: ## Build a .deb package for Debian/Ubuntu/Linux Mint
	./build-deb.sh

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'
