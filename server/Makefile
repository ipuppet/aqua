# -*- Makefile -*-
# aqua web server — Build & Dev
#
# Public commands: build run gen fmt lint test check clean tools new-service help
# Scope most of them with target=<service>, e.g. make check target=passport

.PHONY: default build run gen fmt lint test check clean tools new-service help print-targets

default: help

# -----------------------------------------------------------------------------
# Config
# -----------------------------------------------------------------------------

BINARY := server

# Auto-discover services with services/*/cmd (package main).
# Override: make build targets="root mono"
targets ?= $(shell ./scripts/list_build_targets.sh)
ifeq ($(strip $(targets)),)
$(error no buildable services found under services/*/cmd)
endif

GOOS        ?= linux
GOARCH      ?= amd64
CGO_ENABLED ?= 0
LDFLAGS     ?= -s -w

# Optional scope: target=<service>
ifdef target
ifeq ($(filter $(target),$(targets)),)
$(error unknown target '$(target)'; known: $(targets))
endif
GO_PKG            := ./services/$(target)/...
GO_FMT_PATH       := ./services/$(target)
AQUA_CHECK_TARGET := $(target)
ifneq ($(wildcard proto/aqua/$(target)/.),)
PROTO_PATH_ARGS := --path proto/aqua/$(target)
HAS_PROTO       := 1
else
PROTO_PATH_ARGS :=
HAS_PROTO       := 0
endif
else
GO_PKG            := ./...
GO_FMT_PATH       := .
AQUA_CHECK_TARGET :=
PROTO_PATH_ARGS   :=
HAS_PROTO         := 1
endif

CHECK_SCRIPTS := $(sort $(wildcard scripts/checks/check_*.sh))

# -----------------------------------------------------------------------------
# Public: codegen / format / lint / test / check
# -----------------------------------------------------------------------------

gen: gen.proto gen.go

fmt: fmt.go fmt.proto

lint: lint.go lint.proto lint.custom

test: test.unit

# Full local/CI gate. With target=, is skipped.
check: gen fmt lint test

# -----------------------------------------------------------------------------
# Public: build / run / clean
# -----------------------------------------------------------------------------

ifeq ($(target),)
build:
	@$(foreach t,$(targets),$(MAKE) build.target target=$(t) &&) true
else
build: build.target
endif

run:
ifndef target
	$(error Usage: make run target=<service>)
endif
	@echo "-> running $(target)"
	@CGO_ENABLED=0 GODEBUG=netdns=go+2 POD_IP=127.0.0.1 \
		go run ./services/$(target)/cmd -config-dir=./config

clean:
	@go clean
	@rm -rf ./bin/*

# -----------------------------------------------------------------------------
# Public: misc
# -----------------------------------------------------------------------------

new-service:
ifndef name
	$(error Usage: make new-service name=<service>)
endif
	@./scaffold/new_service.sh "$(name)"

print-targets:
	@echo $(targets)

tools:
	@echo "-> syncing module tools"
	@go mod download
	@echo "-> installing gopls"
	@env -u GOBIN -u GOOS -u GOARCH -u CGO_ENABLED go install golang.org/x/tools/gopls@latest
	@command -v buf >/dev/null 2>&1 || { echo "buf not in PATH (brew install bufbuild/buf/buf)"; exit 1; }
	@echo "-> tools ready"

help:
	@echo "aqua web server"
	@echo ""
	@echo "  target=<service> scopes gen/fmt/lint/test/check  (known: $(targets))"
	@echo ""
	@echo "  make build [target=]   build binaries (default: all)"
	@echo "  make run target=       run one service on host"
	@echo "  make gen [target=]     buf generate + go generate"
	@echo "  make fmt [target=]     format Go + Proto"
	@echo "  make lint [target=]    golangci-lint + buf lint + custom checks"
	@echo "  make test [target=]    unit tests"
	@echo "  make check [target=]   gen + fmt + lint + test"
	@echo "  make new-service name= scaffold a service"
	@echo "  make clean | tools | help"
	@echo ""
	@echo "  custom checks: scripts/checks/check_*.sh (auto-discovered)"
	@echo "  single check:  AQUA_CHECK_TARGET=passport ./scripts/checks/check_wire_build_order.sh"
	@echo "  build env:     GOOS GOARCH CGO_ENABLED LDFLAGS"

# =============================================================================
# Internal
# =============================================================================

# internal — codegen
gen.proto:
ifeq ($(HAS_PROTO),0)
	@echo "-> buf generate skipped (no proto/aqua/$(target))"
else
	@echo "-> buf generate $(PROTO_PATH_ARGS)"
	@buf generate $(PROTO_PATH_ARGS)
endif

gen.go:
	@echo "-> go generate $(GO_PKG)"
	@go generate $(GO_PKG)

# internal — format
fmt.go:
	@echo "-> goimports $(GO_FMT_PATH)"
	@find $(GO_FMT_PATH) -name '*.go' \
		-not -path '*/gen/*' \
		-not -name '*_gen.go' \
		-not -name 'generated.go' \
		-exec goimports -w -local github.com/ipuppet/aqua/server {} +

fmt.proto:
ifeq ($(HAS_PROTO),0)
	@echo "-> buf format skipped (no proto/aqua/$(target))"
else
	@echo "-> buf format $(PROTO_PATH_ARGS)"
	@buf format -w $(PROTO_PATH_ARGS)
endif

# internal — lint
lint.go:
	@echo "-> golangci-lint $(GO_PKG)"
	@go tool golangci-lint run $(GO_PKG)

lint.proto:
ifeq ($(HAS_PROTO),0)
	@echo "-> buf lint skipped (no proto/aqua/$(target))"
else
	@echo "-> buf lint $(PROTO_PATH_ARGS)"
	@buf lint $(PROTO_PATH_ARGS)
endif

lint.custom:
ifeq ($(strip $(CHECK_SCRIPTS)),)
	@echo "-> custom checks: none"
else
	@$(foreach s,$(CHECK_SCRIPTS), \
		echo "-> custom checks: $$(basename $(s))" && \
		AQUA_CHECK_TARGET="$(AQUA_CHECK_TARGET)" bash "$(s)" &&) true
endif

# internal — test
test.unit:
	@echo "-> go test $(GO_PKG)"
	@go test -v $(GO_PKG) || (echo ""; echo "tests failed"; exit 1)

test.coverage:
	@echo "-> go test -cover $(GO_PKG)"
	@go test -coverprofile=coverage.out $(GO_PKG)
	@go tool cover -func=coverage.out

# internal — build one binary
build.target: clean.target
	@mkdir -p ./bin/$(target)
	@echo "-> build $(target) $(GOOS)/$(GOARCH) CGO_ENABLED=$(CGO_ENABLED)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) \
		go build -ldflags '$(LDFLAGS)' -o ./bin/$(target)/$(BINARY) ./services/$(target)/cmd

clean.target:
	@go clean
	@rm -f ./bin/$(target)/$(BINARY)
