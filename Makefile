# Doolox build/development Makefile.
#
# The app is pure Go on Linux and Windows. A C compiler is only required on
# macOS, where the menu-bar icon uses cgo/Objective-C:
#   macOS -> Xcode CLT (xcode-select --install)
#
# For native Windows builds without make, see build.ps1.

GO      ?= go
APP     := doolox
BIN_DIR := bin
PKG     := .
LDFLAGS := -s -w
GOFLAGS := -trimpath

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

build-linux: ## Cross-compile for linux/amd64
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APP)-linux-amd64 $(PKG)

build-windows: ## Cross-compile for windows/amd64
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS) -H=windowsgui" -o $(BIN_DIR)/$(APP)-windows-amd64.exe $(PKG)

build-darwin: ## Build for macOS (must be run on macOS)
	@test "$$(go env GOOS)" = "darwin" || { echo "error: build-darwin requires macOS (tray icon uses cgo/Objective-C)"; exit 1; }
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(APP)-darwin $(PKG)

build-deb: ## Build a .deb package for Debian/Ubuntu/Linux Mint
	./build-deb.sh

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'
