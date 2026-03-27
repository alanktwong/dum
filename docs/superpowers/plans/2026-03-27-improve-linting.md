# Improve Linting Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add 11 new linters, custom linter for single-letter alias ban, and import alias enforcement via importas.

**Architecture:** Custom go/analysis linter at `cmd/linters/main.go` built to `./dist/linters`, integrated via `.golangci.yml` custom linter config.

**Tech Stack:** golangci-lint, go/analysis, custom linter

---

## File Structure

- `cmd/linters/main.go` - New custom linter (single-char alias detection)
- `.golangci.yml` - Updated with new linters and importas config
- `Makefile` - Updated to build custom linter before linting

---

## Task 1: Create Custom Linter

**Files:**
- Create: `cmd/linters/main.go`

- [ ] **Step 1: Create cmd/linters directory**

```bash
mkdir -p cmd/linters
```

- [ ] **Step 2: Write custom linter**

```go
package main

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/pass"
)

var Analyzer = &analysis.Analyzer{
	Name: "linters",
	Doc:  "Custom linters for dum project",
	Run:  run,
}

func run(pass *pass.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
				for _, spec := range genDecl.Specs {
					importSpec := spec.(*ast.ImportSpec)
					if importSpec.Name != nil {
						name := importSpec.Name.Name
						if len(name) == 1 && name != "_" {
							pass.Reportf(importSpec.Pos(), "import alias %q is a single character", name)
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
```

- [ ] **Step 3: Initialize go module for linters**

```bash
cd cmd/linters && go mod init awong/dotfiles/cmd/linters
```

- [ ] **Step 4: Add dependencies**

```bash
cd cmd/linters && go get golang.org/x/tools/go/analysis
```

- [ ] **Step 5: Test build**

```bash
go build -o ./dist/linters ./cmd/linters
```

- [ ] **Step 6: Commit**

```bash
git add cmd/linters/
git commit -m "feat(lint): add custom linters package"
```

---

## Task 2: Update .golangci.yml

**Files:**
- Modify: `.golangci.yml`

- [ ] **Step 1: Preserve existing linters.exclusions**

The current config has exclusions for generated files and a specific godox rule for expansion.go. Add this at the start of the linters section:
```yaml
linters:
  exclusions:
    generated: lax
    presets:
      - common-false-positives
    rules:
      - path: internal/utilities/expansion.go
        linters:
          - godox
```

- [ ] **Step 2: Update linters enable list**

Add to the enable list:
- `asciicheck`
- `bidichk`
- `errcheck`
- `errname`
- `errorlint`
- `gocritic`
- `gocyclo`
- `importas`
- `linters` (custom)
- `staticcheck`
- `unused`

This adds **11 new linters** (23 + 11 = 34 total, matching the spec).

Remove from current: none (add to existing)

- [ ] **Step 3: Add custom linter config**

Add after linters section:
```yaml
custom:
  linters:
    path: ./dist/linters
```

- [ ] **Step 4: Add importas settings**

Add linters-settings section:
```yaml
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
```

- [ ] **Step 5: Add test file exclusions**

In issues section (should already have `max-issues-per-linter` and `max-same-issues`), add `exclusions`:
```yaml
issues:
  max-issues-per-linter: 0
  max-same-issues: 0
  exclusions:
    rules:
      - path: .*_test\.go
        linters:
          - importas
          - linters
```

- [ ] **Step 6: Preserve formatters section**

Add at the end of the config file (or preserve existing):
```yaml
formatters:
  enable:
    - gofumpt
    - goimports
  exclusions:
    generated: lax
```

- [ ] **Step 7: Commit**

```bash
git add .golangci.yml
git commit -m "feat(lint): enable additional linters and import alias rules"
```

---

## Task 3: Update Makefile

**Files:**
- Modify: `Makefile`

- [ ] **Step 1: Add linters build target**

Add after the build target:
```makefile
.PHONY: linters-build
linters-build: ## build custom linters
	@printf ${COLOR} "Building linters ..."
	@mkdir -p ${OUT_DIR}
	@go build -o ${OUT_DIR}/linters ./cmd/linters
```

- [ ] **Step 2: Update lint target to build linters first**

Change:
```makefile
lint: build ## go linting
```

To:
```makefile
lint: linters-build ## go linting
```

- [ ] **Step 3: Commit**

```bash
git add Makefile
git commit -e "feat(lint): add linters-build to Makefile"
```

---

## Task 4: Run Lint and Fix Violations

**Files:**
- Modify: All `.go` files with import violations

- [ ] **Step 1: Run lint to identify violations**

```bash
make lint
```

Expected: Multiple errors about:
- Single-letter aliases (l, i, t, pl, ty, etc.)
- Missing aliases for vendor packages
- Wrong aliases for internal packages

- [ ] **Step 2: Run with fix to auto-fix importas issues**

```bash
golangci-lint run ./... --fix
```

This should fix importas issues (wrong/missing aliases).

- [ ] **Step 3: Manually fix remaining issues**

For single-letter aliases, need to update imports manually. Common replacements:
- `l` → `lg` (logging)
- `i` → `ti` (tasks/installer)
- `t` → `ty` (types)
- `pl` → `pl` (plays) - already 2-char
- `ty` → `ty` (types) - already 2-char
- `ext` → `ext` (external) - already 2-char

- [ ] **Step 4: Run lint again**

```bash
make lint
```

Expected: All lint errors resolved.

- [ ] **Step 5: Run tests to verify nothing broke**

```bash
make test
```

Expected: All tests pass.

- [ ] **Step 6: Commit all fixes**

```bash
git add .
git commit -m "fix: update imports to use standardized aliases"
```

---

## Task 5: Verify Complete

- [ ] **Step 1: Run full check**

```bash
make check
```

- [ ] **Step 2: Verify lint passes cleanly**

```bash
make lint
```

Expected: No errors.

- [ ] **Step 3: Final commit**

```bash
git commit -m "chore: complete linting improvements"
```
