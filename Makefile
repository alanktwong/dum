# update app name. this is the name of binary
# See https://www.mohitkhare.com/blog/go-makefile/
OUT_DIR="./dist"

DUM_APP=dum
DUM_SOURCE=./cmd/${DUM_APP}
GIT_TAG := $(shell git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
GIT_DIRTY := $(shell git diff --quiet HEAD 2>/dev/null || echo "-dirty")
DUM_VERSION ?= $(if $(GIT_TAG),$(GIT_TAG)-$(GIT_COMMIT)$(GIT_DIRTY),$(GIT_COMMIT)$(GIT_DIRTY))
DUM_EXECUTABLE="$(OUT_DIR)/$(DUM_APP)"

ALL_PACKAGES=$(shell go list ./... | grep -v /vendor)
SHELL := /bin/bash # Use bash syntax

COVERAGE_THRESHOLD ?= 80

# Optional colors to beautify output
BLACK   := $(shell tput -Txterm setaf 0)
RED     := $(shell tput -Txterm setaf 1)
GREEN   := $(shell tput -Txterm setaf 2)
YELLOW  := $(shell tput -Txterm setaf 3)
BLUE    := $(shell tput -Txterm setaf 4)
MAGENTA := $(shell tput -Txterm setaf 5)
CYAN    := $(shell tput -Txterm setaf 6)
WHITE   := $(shell tput -Txterm setaf 7)
RESET   := $(shell tput -Txterm sgr0)
COLOR   := "\e[1;36m%s\e[0m\n"

## Tools
.PHONY: print-go-version
print-go-version: ## prints the go version
	echo "Go Version: $(shell go version)"

.PHONY: rhard
rhard: ## git reset hard
	@printf ${COLOR} "git resetting hard ..."
	@git reset --hard HEAD
	@for dir in vim/vim.d/plugged/*/; do \
		if [ -d "$$dir/.git" ]; then \
			echo "Resetting $$dir ..."; \
			git -C "$$dir" reset --hard HEAD; \
			git -C "$$dir" clean -fd; \
		fi; \
	done

.PHONY: install
install: ## Installs tools
	@printf ${COLOR} "Installing tools ..."
	@bash ./tools/dev-prep.sh

.PHONY: check_tools
check_tools: ## Check tools
	@printf ${COLOR} "Checking tools ..."
	@bash ./tools/dev-env-check.sh

.PHONY: ci-setup
ci-setup: ## install tools needed for CI (Go tools only)
	@printf ${COLOR} "Installing CI tools ..."
	@go install github.com/abice/go-enum@latest
	@go install github.com/vektra/mockery/v3@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install mvdan.cc/gofumpt@latest

## Compile
make_out:
	@mkdir -p ${OUT_DIR}

clean: ## cleans binary and other generated files
	@printf ${COLOR} "Cleaning up ..."
	@go clean -cache
	@rm -rf $(OUT_DIR)
	@rm -f ./internal/**/*_mocks.go
	@rm -f ./internal/**/*_enum.go

.PHONY: tidy
tidy: ## runs tidy to fix go.mod dependencies
	@printf ${COLOR} "Tidying ..."
	@go mod tidy

.PHONY: vendor
vendor: tidy ## all packages required to support builds and tests in the /vendor directory
	@printf ${COLOR} "Vendoring ..."
	@go mod vendor

.PHONY: generate
generate: ## generate code
	@printf ${COLOR} "Generating ..."
	@go generate ./...
	@for config in cfg/mockery*.yml; do \
		mockery --config "$$config"; \
	done

.PHONY: build
build: make_out fmt vendor generate ## local build
	@printf ${COLOR} "Building ..."
	@go build -v -ldflags "-X main.version=$(DUM_VERSION) -X main.commit=$(GIT_COMMIT)" -o $(DUM_EXECUTABLE)-${DUM_VERSION} $(DUM_SOURCE)

.PHONY: linters-build
linters-build: ## build custom linters
	@printf ${COLOR} "Building linters ..."
	@mkdir -p ${OUT_DIR}
	@cd cmd/linters && go build -o ../dist/linters .

.PHONY: release
release: ## compile release binaries for all platforms (darwin/linux, arm64/amd64/386)
	@printf ${COLOR} "Compiling releases for every OS and Platform ..."
	@goreleaser build --clean --snapshot --config cfg/goreleaser.local.yaml

release-build: build  ## compile snapshot binaries without archives
	@printf ${COLOR} "Compiling builds for every OS and Platform ..."
	@goreleaser build --clean --snapshot --config cfg/goreleaser.snapshot.yaml

## Quality
.PHONY: check
check: build fmt fix lint vet test check-coverage ## runs code quality checks
	@printf ${COLOR} "Checking code quality ..."

# Append || true below if blocking local development
.PHONY: lint
lint: build ## go linting. Update and use specific lint tool and options
	@printf ${COLOR} "Linting ..."
	@golangci-lint run cmd/...
	@golangci-lint run internal/...

.PHONY: vet
vet: ## go vet
	@printf ${COLOR} "Vetting ..."
	@go vet ./...

.PHONY: fmt
fmt: ## runs go formatter
	@printf ${COLOR} "Formatting ..."
	@go fmt ./...
	@gofumpt -w $$(go list -f '{{.Dir}}' ./...)

.PHONY: fix
fix: ## runs go fix to update code to use new language features
	@printf ${COLOR} "Fixing ..."
	@go fix ./...

## Test
.PHONY: test
test: build ## runs tests and create generates coverage report
	@printf ${COLOR} "Testing ..."
	@go test -v -cover \
	    -coverprofile=${OUT_DIR}/cover.out \
	    -covermode=atomic \
	    -coverpkg=./... \
	    -race ./...

.PHONY: coverage
coverage: test ## displays test coverage report and checks threshold
	@printf ${COLOR} "Checking coverage of tests ..."
	@go-test-coverage --config=./cfg/testcoverage.yml || true
	@printf ${COLOR} "Generating HTML report ..."
	@go tool cover -html=${OUT_DIR}/cover.out -o ${OUT_DIR}/coverage.html
	@printf ${COLOR} "HTML report: ${OUT_DIR}/coverage.html"

.PHONY: check-coverage
check-coverage: test ## checks that test coverage meets the minimum threshold
	@printf ${COLOR} "Checking coverage threshold ..."
	@go-test-coverage --config=./cfg/testcoverage.yml

.PHONY: profile
profile: ## profiles
	@printf ${COLOR} "Profiling ..."
	@go test -cpuprofile cpu.prof -memprofile mem.prof -v -bench ./...

.PHONY: bench
bench: ## benchmarks
	@printf ${COLOR} "Benchmarking ..."
	@go test -v -benchmem -bench=. ./...

## All
.PHONY: all
all: build check ## runs setup, quality checks and builds
	@printf ${COLOR} "running all ..."

## Help
.PHONY: help
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)