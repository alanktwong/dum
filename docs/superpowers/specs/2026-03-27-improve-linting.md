# Design: Improve Linting

**Date**: 2026-03-27
**Status**: Approved

## Background

This design improves the project's linting by:
1. Enabling additional high-value linters from popular Go projects
2. Adding rules to enforce consistent import alias patterns

The current config (based on charmbracelet/gum) enables 24 linters. This adds 10 more recommended linters plus import alias enforcement.

---

## Part 1: Additional Linters

Based on the golangci-lint golden config (maratori/golangci-lint-config, 477 stars), enable these linters:

| Linter | Purpose |
|--------|---------|
| `errcheck` | Checks for unchecked errors |
| `govet` | Standard vet checks |
| `staticcheck` | Advanced static analysis |
| `unused` | Checks for unused code |
| `asciicheck` | Non-ASCII identifiers |
| `bidichk` | Dangerous unicode sequences |
| `errorlint` | Error wrapping issues |
| `gocritic` | Comprehensive style/bug checks |
| `gocyclo` | Cyclomatic complexity |
| `errname` | Error naming conventions |

---

## Part 2: Import Alias Rules

### Requirements

1. **Single-letter alias ban**: No import alias may be a single character (e.g., `x`, `y`, `z`)
2. **All non-stdlib imports require aliases**: Every import except Go standard library must use a defined alias
3. **Test file exemptions**: Both rules exclude `*_test.go` files

### Standardized Aliases

All non-stdlib imports must use these aliases:

**Internal packages:**
| Package | Alias |
|---------|-------|
| `awong/dotfiles/internal/cmd` | `cd` |
| `awong/dotfiles/internal/external` | `ext` |
| `awong/dotfiles/internal/factory` | `fy` |
| `awong/dotfiles/internal/logging` | `lg` |
| `awong/dotfiles/internal/playbook` | `pb` |
| `awong/dotfiles/internal/plays` | `pl` |
| `awong/dotfiles/internal/plays/gen` | `plg` |
| `awong/dotfiles/internal/tasks` | `tk` |
| `awong/dotfiles/internal/tasks/installer` | `ti` |
| `awong/dotfiles/internal/types` | `ty` |
| `awong/dotfiles/internal/types/gen` | `tyg` |

**Vendor packages:**
| Package | Alias |
|---------|-------|
| `github.com/charmbracelet/log` | `clog` |
| `github.com/elliotchance/orderedmap/v3` | `omv3` |
| `github.com/spf13/cobra` | `ca` |
| `github.com/stretchr/testify` | `tt` |
| `github.com/stretchr/testify/assert` | `asrt` |
| `github.com/stretchr/testify/mock` | `mck` |
| `gopkg.in/yaml.v3` | `yamlv3` |

Note: All aliases use 2+ characters to comply with single-letter ban.

### Custom Linters

**Location:** `cmd/linters/main.go`

A custom go/analysis linter that detects and reports single-character import aliases.

**Behavior:**
- Scans all `.go` files (except `*_test.go`)
- Reports error for any import alias that is exactly 1 character long
- Example: `x "example.com/pkg"` → Error

Additional custom linters can be added to this package.

---

## Implementation

### `.golangci.yml` changes:

```yaml
version: "2"
run:
  tests: false
linters:
  enable:
    - asciicheck
    - bidichk
    - bodyclose
    - errcheck
    - errname
    - errorlint
    - exhaustive
    - goconst
    - gocritic
    - godot
    - godox
    - gocyclo
    - gomoddirectives
    - goprintffuncname
    - gosec
    - govet
    - importas
    - misspell
    - nakedret
    - nestif
    - nilerr
    - noctx
    - nolintlint
    - prealloc
    - revive
    - rowserrcheck
    - sqlclosecheck
    - staticcheck
    - tparallel
    - unconvert
    - unparam
    - unused
    - whitespace
    - wrapcheck
    - linters
  exclusions:
    generated: lax
    presets:
      - common-false-positives
    rules:
      - path: internal/utilities/expansion.go
        linters:
          - godox

custom:
  linters:
    path: ./dist/linters

linters-settings:
  importas:
    no-unaliased: true
    no-extra-aliases: true
    alias:
      # Internal packages
      - prefix: awong/dotfiles/internal/cmd
        alias: cd
      - prefix: awong/dotfiles/internal/external
        alias: ext
      - prefix: awong/dotfiles/internal/factory
        alias: fy
      - prefix: awong/dotfiles/internal/logging
        alias: lg
      - prefix: awong/dotfiles/internal/playbook
        alias: pb
      - prefix: awong/dotfiles/internal/plays
        alias: pl
      - prefix: awong/dotfiles/internal/plays/gen
        alias: plg
      - prefix: awong/dotfiles/internal/tasks
        alias: tk
      - prefix: awong/dotfiles/internal/tasks/installer
        alias: ti
      - prefix: awong/dotfiles/internal/types
        alias: ty
      - prefix: awong/dotfiles/internal/types/gen
        alias: tyg
      # Vendor packages
      - prefix: github.com/charmbracelet/log
        alias: clog
      - prefix: github.com/elliotchance/orderedmap/v3
        alias: omv3
      - prefix: github.com/spf13/cobra
        alias: ca
      - prefix: github.com/stretchr/testify
        alias: tt
      - prefix: github.com/stretchr/testify/assert
        alias: asrt
      - prefix: github.com/stretchr/testify/mock
        alias: mck
      - prefix: gopkg.in/yaml.v3
        alias: yamlv3

issues:
  max-issues-per-linter: 0
  max-same-issues: 0
  exclusions:
    rules:
      - path: .*_test\.go
        linters:
          - importas
          - linters

formatters:
  enable:
    - gofumpt
    - goimports
  exclusions:
    generated: lax
```

### How Each Requirement Is Addressed

| Requirement | Tool | Configuration |
|-------------|------|---------------|
| Additional linters | 10 new built-in linters | Added to `enable` list |
| Ban single-letter aliases | custom `linters` | N/A - linter logic |
| Require specific aliases for all non-stdlib packages | built-in `importas` | `alias:` mappings for internal + vendor + `no-unaliased: true` + `no-extra-aliases: true` |
| Exempt test files | Both | `issues.exclusions.rules` |

---

## Testing

1. **Verify new linters work:** Run `make lint` with new linters enabled
2. **Verify custom linter detects single-letter aliases:** Create a test file with `x "example.com/pkg"` and run lint
3. **Verify importas enforces internal aliases:** Create a file importing `awong/dotfiles/internal/types` without alias and run lint
4. **Verify exemptions:** Ensure `*_test.go` files are not flagged
5. **Fix all violations:** Update codebase to comply with new rules

---

## Migration Plan

1. Create custom linter at `cmd/linters/main.go`
2. Build the linter: `go build -o ./dist/linters ./cmd/linters`
3. Update `.golangci.yml` with all changes
4. Run `make lint` to identify violations
5. Fix all violations in the codebase
6. Commit changes
