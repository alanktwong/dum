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
.DEFAULT_GOAL := help
DISPATCH = OUT_DIR="$(OUT_DIR)" DUM_VERSION="$(DUM_VERSION)" DUM_APP="$(DUM_APP)" DUM_SOURCE="$(DUM_SOURCE)" GIT_COMMIT="$(GIT_COMMIT)" DUM_EXECUTABLE="$(DUM_EXECUTABLE)" ./tools/build.sh

## Tools
.PHONY: print-go-version
print-go-version: ## prints the go version
	@$(DISPATCH) print-go-version

.PHONY: rhard
rhard: ## git reset hard
	@$(DISPATCH) rhard

.PHONY: install
install: ## Installs tools
	@$(DISPATCH) install

.PHONY: check_tools
check_tools: ## Check tools
	@$(DISPATCH) check_tools

.PHONY: ci-setup
ci-setup: ## install tools needed for CI (Go tools only)
	@$(DISPATCH) ci-setup

## Compile
make_out:
	@$(DISPATCH) make_out

.PHONY: clean
clean: ## cleans binary and other generated files
	@$(DISPATCH) clean

.PHONY: tidy
tidy: ## runs tidy to fix go.mod dependencies
	@$(DISPATCH) tidy

.PHONY: vendor
vendor: tidy ## all packages required to support builds and tests in the /vendor directory
	@$(DISPATCH) vendor

.PHONY: generate
generate: ## generate code
	@$(DISPATCH) generate

.PHONY: build
build: fmt vendor generate ## local build
	@$(DISPATCH) make_out
	@$(DISPATCH) build

.PHONY: linters-build
linters-build: ## build custom linters
	@$(DISPATCH) linters-build

.PHONY: release
release: ## compile release binaries for darwin/linux on arm64/amd64
	@$(DISPATCH) release

release-build: build  ## compile snapshot binaries without archives
	@$(DISPATCH) release-build

## Quality
.PHONY: check
check: build fmt fix lint vet test check-coverage ## runs code quality checks
	@$(DISPATCH) check

# Append || true below if blocking local development
.PHONY: lint
lint: build ## go linting. Update and use specific lint tool and options
	@$(DISPATCH) lint

.PHONY: vet
vet: ## go vet
	@$(DISPATCH) vet

.PHONY: fmt
fmt: ## runs go formatter
	@$(DISPATCH) fmt

.PHONY: fix
fix: ## runs go fix to update code to use new language features
	@$(DISPATCH) fix

## Test
.PHONY: test
test: build ## runs tests and create generates coverage report
	@$(DISPATCH) test

.PHONY: coverage
coverage: test ## displays test coverage report and checks threshold
	@$(DISPATCH) coverage

.PHONY: check-coverage
check-coverage: test ## checks that test coverage meets the minimum threshold
	@$(DISPATCH) check-coverage

.PHONY: profile
profile: ## profiles
	@$(DISPATCH) profile

.PHONY: bench
bench: ## benchmarks
	@$(DISPATCH) bench

## All
.PHONY: all
all: build check ## runs setup, quality checks and builds
	@$(DISPATCH) all

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
