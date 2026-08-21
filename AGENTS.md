# AGENTS.md - Agent Coding Guidelines for Dum

This document provides guidelines for agentic coding agents operating in this repository.

## Project Overview

- **Project**: Dum - A CLI tool for managing software installations and configurations
- **Language**: Go 1.26 (requires 1.24.4+)
- **Module**: `github.com/alanktwong/dum`
- **Description**: A hand-rolled CLI tool similar to Ansible, written in Go. Manages software installations and configurations on macOS/Linux laptops.
- **Configuration**: Uses `dum.yml` for defining playbooks, plays, and tasks (legacy `installer.yml` still resolved as fallback)

### Key Concepts

- **Task**: A specific task for installing something (e.g., `BrewTask` runs `brew install ...`)
- **Play**: A group of tasks
- **Playbook**: A grouping of plays

## Build/Lint/Test Commands

### Primary Commands (via Makefile)

The Makefile targets are wrappers around the repository-root `tools/build.sh` dispatcher. Use the Makefile for normal development; run `./tools/build.sh help` to inspect direct subcommands.

```bash
make build           # Local build (includes tidy, fmt, vendor, generate)
make test            # Run tests with coverage
make lint            # Run golangci-lint
make vet             # Go vet
make fmt             # Go fmt
make fix             # Runs go fix to update code to new language features
make generate        # Run go generate and mockery
make tag             # Compute next semver from conventional commits and create git tag
make check-coverage  # Checks coverage meets minimum threshold (default 80%)
make all             # Runs build and check (quality + build)
make clean           # Clean generated files
```

Direct dispatcher commands must be run from the repository root:

```bash
./tools/build.sh help
```


### Code Generation

We are using [go-enum](https://github.com/abice/go-enum) to generate enums from Go struct field tags.
The generated enums, by convention, are suffixed as `_enum.go` and are gitignored.

We are using [mockery](https://vektra.github.io/mockery/v3.0/) to generate
mocks of Go interfaces for testing.
The mockery configurations live in `cfg/mockery.*.yml` and the generated
mocks are gitignored since they should follow the naming convention `*_mocks.go`.

### Versioning and Releases

Versions follow [Conventional Commits](https://www.conventionalcommits.org), enforced locally by
commitlint via pre-commit. [svu](https://github.com/caarlos0/svu) computes the next version from
merged commits; configuration lives in `.svu.yml` (`always: true`, `v0: true`).

Release flow: squash-merge PRs with conventional titles, run `make tag`, then
`git push origin <tag>`. Pushing a tag triggers GoReleaser via `.github/workflows/release.yml`,
which publishes per-platform archives and a SHA256 checksums file to the GitHub Release.

### Coverage Threshold

- Default threshold is **80%**
- Override via environment variable or command line:
    ```bash
    COVERAGE_THRESHOLD=90 make check-coverage
    ```
- `check` target enforces the threshold as part of quality checks

### Single Test Execution

```bash
# Run specific test
go test -v ./internal/playbook/... -run TestBrewTask_Install

# Run tests in specific package
go test -v ./internal/playbook/

# Run with race detector
go test -v -race ./...
```

### Development Setup

```bash
make install    # Install tools (brew deps incl. svu/goreleaser; go-enum, goimports, gocov, golangci-lint, mockery). Node comes from nvm per .nvmrc
make tidy       # Fix go.mod dependencies
make vendor     # Vendor dependencies
make coverage   # View coverage report in browser and check threshold
make check-coverage  # Enforce coverage threshold
```

## Code Style Guidelines

### Formatting

- **Indentation**: Use tabs for Go files and Makefile recipes; respect `.editorconfig`
- **Line endings**: LF (Unix-style)
- **Formatter**: Uses `gofumpt`, `goimports`, and `golines` through `.golangci.yml`

### Import Order

Standard library imports first, then external packages:

```go
import (
    "github.com/alanktwong/dum/pkg/external"
    l "github.com/alanktwong/dum/pkg/logging"
    "context"
    "fmt"

    om "github.com/elliotchance/orderedmap/v3"
)
```

### Naming Conventions

- **Types/Functions**: PascalCase (e.g., `BrewTask`, `NewExecutor`)
- **Variables/Interfaces**: camelCase (e.g., `input`, `task`)
- **Interfaces**: Use descriptive names like `Lister`, `Installer`, `Task`
- **Packages**: Short, lowercase names (e.g., `playbook`, `cmd`, `external`)

### Error Handling

- Use `fmt.Errorf` with wrapped errors: `fmt.Errorf("failed to install task %s: %w", id, err)`
- Return errors rather than logging unless it's the final consumer
- Validate inputs early with descriptive errors

### CLI Framework Boundary

- Keep all interactions with CLI frameworks such as Cobra Viper inside `/cmd/` and its subdirectories.
- Command packages under `/cmd/` own framework wiring, flags, help text, argument parsing, and CLI-specific error handling.
- Packages under `/internal/` must remain framework-agnostic and expose typed structs, services, and errors instead of Cobra or Viper types.

### Types and Interfaces

The codebase uses a task-based pattern:

```go
// Installer can install given an input.
type Installer interface {
    Install(ctx context.Context, input *Input) (*TaskResult, error)
}

// Task can install given an input. Dry runs log disabled tasks instead of
// skipping them so the output shows the full playbook manifest.
type Task interface {
    Installer
}
```

### Context Usage

- All public methods should accept `context.Context` as the first argument
- Use `context.Background()` in tests

### Testing

- Tests use `testify/assert` package
- Mocks are generated with `mockery` (configured in `cfg/mockery.yml`)
- Test files are named `*_test.go`
- Mock files are named `*_mocks_test.go`

### Code Generation

- Enums use `go-enum` with directive: `//go:generate ../../../../bin/go-enum ...`
- Run `make generate` to regenerate code and mocks

### Linters (Enabled in golangci-lint)

- bodyclose, exhaustive, goconst, godot, godox
- gomoddirectives, goprintffuncname, gosec, misspell
- nakedret, nestif, nilerr, noctx, nolintlint
- prealloc, revive, rowserrcheck, sqlclosecheck
- tparallel, unconvert, unparam, whitespace, wrapcheck
- asciicheck, bidichk, errcheck, errname, errorlint
- gocritic, gocyclo, importas, linters, staticcheck, unused

### Import Alias Standards

- **No single-letter aliases**: Aliases like `l`, `i`, `t` are banned
- **Required aliases**: All internal and vendor packages must use standardized aliases:
    - Internal: `cd`, `ext`, `fy`, `lg`, `pb`, `pl`, `plg`, `tk`, `ti`, `ty`, `tyg`
    - Vendor: `clog`, `omv3`, `ca`, `tt`, `asrt`, `mck`, `yamlv3`
- **Auto-fix**: Run `golangci-lint run ./... --fix` to fix importas violations

## Project Structure

```
/cmd/dum/               # CLI entry point and all Cobra-facing command packages
/cmd/linters/            # Custom lint tooling
/internal/factory/       # Typed installer configuration and runtime construction
/internal/playbook/      # Core task execution logic
/internal/external/      # External tool wrappers (brew, git, etc.)
/internal/logging/       # Logging utilities
/internal/rename/        # Cobra-free file rename service
/internal/types/         # Shared domain types
/internal/yaml/          # Typed installer YAML loading and validation
```

## Important Files

- `.golangci.yml` - Linter configuration
- `.svu.yml` - svu versioning config
- `.commitlintrc.yml` - commitlint conventional-commits config
- `.nvmrc` - node version pin for the commitlint pre-commit hook
- `.pre-commit-config.yaml` - local hook enforcement (incl. commitlint at commit-msg)
- `.github/workflows/release.yml` - tag-push release pipeline
- `cfg/mockery.yml` - Mock generation config
- `cfg/goreleaser.yaml` - GoReleaser build config
- `.editorconfig` - Editor formatting rules
- `Makefile` - Build and development commands
- `dum.yml` - Playbook configuration file (legacy `installer.yml` still resolved as fallback)
