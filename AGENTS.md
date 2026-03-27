# AGENTS.md - Agent Coding Guidelines for Dum

This document provides guidelines for agentic coding agents operating in this repository.

## Project Overview

- **Project**: Dum - A CLI tool for managing software installations and configurations
- **Language**: Go 1.26 (requires 1.24.4+)
- **Module**: `awong/dotfiles`
- **Description**: A hand-rolled CLI tool similar to Ansible, written in Go. Manages software installations and configurations on macOS/Linux laptops.
- **Configuration**: Uses `installer.yml` for defining playbooks, plays, and tasks

### Key Concepts
- **Task**: A specific task for installing something (e.g., `BrewTask` runs `brew install ...`)
- **Play**: A group of tasks
- **Playbook**: A grouping of plays

## Build/Lint/Test Commands

### Primary Commands (via Makefile)
```bash
make build           # Local build (includes tidy, fmt, vendor, generate)
make test            # Run tests with coverage
make lint            # Run golangci-lint
make vet             # Go vet
make fmt             # Go fmt
make fix             # Runs go fix to update code to new language features
make generate        # Run go generate and mockery
make check-coverage  # Checks coverage meets minimum threshold (default 80%)
make all             # Runs build and check (quality + build)
make clean           # Clean generated files
```
### Code Generation

We are using [go-enum](https://github.com/abice/go-enum) to generate enums from Go struct field tags.
The generated enums, by convention, are suffixed as `_enum.go` and are gitignored.

We are using [mockery](https://vektra.github.io/mockery/v3.0/) to generate
mocks of Go interfaces for testing.
The mockery configurations live in `cfg/mockery.*.yml` and the generated
mocks are gitignored since they should follow the naming convention `*_mocks.go`.

### Coverage Threshold
- Default threshold is **80%**
- Override via environment variable or command line:
  ```bash
  COVERAGE_THRESHOLD=90 make coverage-check
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
make install    # Install tools (go-enum, goimports, gocov, golangci-lint, mockery)
make tidy       # Fix go.mod dependencies
make vendor     # Vendor dependencies
make coverage   # View coverage report in browser and check threshold
make coverage-check  # Enforce coverage threshold
```

## Code Style Guidelines

### Formatting
- **Indentation**: Use tabs for Go files (per `.editorconfig`)
- **Line endings**: LF (Unix-style)
- **Formatter**: Uses `gofumpt` and `goimports` (configured in `.golangci.yml`)

### Import Order
Standard library imports first, then external packages:
```go
import (
    "awong/dotfiles/pkg/external"
    l "awong/dotfiles/pkg/logging"
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

### Types and Interfaces

The codebase uses a task-based pattern:
```go
// Lister can list given an input.
type Lister interface {
    List(ctx context.Context, input *Input) (*TaskResult, error)
}

// Installer can install given an input.
type Installer interface {
    Install(ctx context.Context, input *Input) (*TaskResult, error)
}

// Task can list and install given an input.
type Task interface {
    Lister
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
/cmd/dum/main.go       # Entry point
/internal/cmd/              # CLI commands (Cobra-based)
/internal/playbook/         # Core task execution logic
/internal/external/         # External tool wrappers (brew, git, etc.)
/internal/logging/          # Logging utilities
/internal/enums/            # Enum definitions
```

## Important Files

- `.golangci.yml` - Linter configuration
- `cfg/mockery.yml` - Mock generation config
- `cfg/goreleaser.yaml` - GoReleaser build config
- `.editorconfig` - Editor formatting rules
- `Makefile` - Build and development commands
- `installer.yml` - Playbook configuration file
