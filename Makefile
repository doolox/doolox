# Doolox build/development Makefile.
#
# The admin interface opens in a native web view (webview/webview_go), which
# links through cgo on every platform and needs the matching native toolkit:
#   Linux   -> gcc/g++ plus GTK 3 and WebKitGTK (libwebkit2gtk-4.0-dev)
#   macOS   -> Xcode CLT (xcode-select --install)
#   Windows -> MinGW-w64 gcc/g++ (see INSTALL.md)
#
# The native web view cannot be cross-compiled, so build natively on each OS
# (see build.ps1 for Windows).

GO      ?= go
APP     := doolox
BIN_DIR := bin
PKG     := .
LDFLAGS := -s -w
GOFLAGS := -trimpath

.PHONY: all build run test vet fmt tidy clean install \
        build-darwin build-deb help

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

build-darwin: ## Build for macOS (must be run on macOS)
	@test "$$(go env GOOS)" = "darwin" || { echo "error: build-darwin requires macOS (tray icon uses cgo/Objective-C)"; exit 1; }
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APP)-darwin $(PKG)

build-deb: ## Build a .deb package for Debian/Ubuntu/Linux Mint
	./build-deb.sh

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'
