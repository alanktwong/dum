# Viper Configuration Integration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use `subagent-driven-development` (recommended) or `executing-plans` to implement this plan task-by-task. Steps use checkbox syntax for tracking.

**Goal:** Move installer-path and log-level environment resolution to isolated Viper helpers under `cmd/dum`, preserving existing CLI behavior and keeping `internal/` Viper-free.

**Architecture:** `cmd/dum/cli` owns Viper setup and converts environment values into typed path/log-level results. `cmd/dum/root.go` initializes the logger with resolved level; command packages consume resolved file paths through existing typed factory options. `internal/yaml` continues typed YAML parsing; `internal/logging` provides logger construction only.

**Tech Stack:** Go 1.26, Cobra, Viper, Charmbracelet log, Testify.

---

## File map

- Modify: `go.mod`, `go.sum`, `vendor/modules.txt`, and vendored Viper dependencies through the repository dependency workflow.
- Modify: `cmd/dum/cli/helpers.go` — isolated Viper configuration resolver and config-path helpers.
- Create or modify: `cmd/dum/cli/log_level.go` — Viper-backed `ZSH_LOG_LEVEL` parsing into `clog.Level`.
- Modify: `cmd/dum/root.go` — resolve log level before constructing logger.
- Modify: `cmd/dum/log/log.go` — remove `internal/logging.EnvLevel()` lookup.
- Modify: `internal/logging/logging.go` — retain logger construction API without environment behavior if needed.
- Delete or modify: `internal/logging/env_level.go` and `internal/logging/logging_test.go` — remove obsolete env ownership while preserving applicable level behavior tests at CLI boundary.
- Modify: `cmd/dum/flags_test.go`, `cmd/dum/root_test.go`, and relevant `cmd/dum/log/*_test.go` — test path resolution, log-level parsing, root wiring, and CLI behavior.

## Task 1: Add Viper dependency

**Files:** `go.mod`, `go.sum`, `vendor/`

- [ ] Add `github.com/spf13/viper` using the project’s existing dependency/vendor workflow.
- [ ] Confirm `go.mod` uses the selected Viper version compatible with Go 1.26 and existing dependencies.
- [ ] Run the repository dependency/vendor command required by `Makefile`.
- [ ] Commit dependency changes:

```sh
git add go.mod go.sum vendor
git commit -m "build: add Viper dependency"
```

## Task 2: Characterize resolver behavior

**Files:** `cmd/dum/flags_test.go`, `cmd/dum/install/install_test.go`, `cmd/dum/list/list_test.go`, new `cmd/dum/cli/config_test.go` if package boundaries require it

- [ ] Add characterization tests proving `INSTALLER_CONFIG` wins over `XDG_CONFIG_HOME`.
- [ ] Add characterization tests proving `XDG_CONFIG_HOME` builds `<xdg>/dum/installer.yml`.
- [ ] Add characterization tests proving empty values fall back to `~/.config/dum/installer.yml`.
- [ ] Add characterization tests proving repeated resolver calls do not leak state after `t.Setenv` changes.
- [ ] Add command-consumer coverage in install/list tests: set `INSTALLER_CONFIG`, execute with `--file`, and assert mock factory receives explicit path.
- [ ] Run focused characterization tests and confirm PASS before implementation:

```sh
go test ./cmd/dum/... -run 'Test(GetDefaultConfig|AddFileFlag|.*Config|.*File)'
```

## Task 3: Implement isolated config-path resolver

**Files:** `cmd/dum/cli/helpers.go` or focused `cmd/dum/cli/config.go`

- [ ] Create/configure fresh Viper instance per resolver operation; no package-global Viper state.
- [ ] Bind/read `INSTALLER_CONFIG` and `XDG_CONFIG_HOME` through Viper.
- [ ] Preserve precedence: non-empty `INSTALLER_CONFIG`, then non-empty `XDG_CONFIG_HOME`, then literal `~/.config` fallback.
- [ ] Preserve returned path strings, including `~` rather than expanding it.
- [ ] Keep `AddFileFlag` API and Cobra behavior unchanged.
- [ ] Run focused path and command tests; expect PASS.
- [ ] Commit:

```sh
git add cmd/dum/cli cmd/dum/flags_test.go
git commit -m "refactor: resolve config path with Viper"
```

## Task 4: Characterize and move log-level parsing

**Files:** `internal/logging/logging_test.go`, `internal/logging/env_level.go`, new `cmd/dum/cli/log_level.go`, `cmd/dum/cli/log_level_test.go`

- [ ] Extend existing `internal/logging` characterization tests with named values, exact numeric boundaries `10`, `20`, `30`, `40`, `50`, negative and over-50 values, invalid and empty fallback to `warn`.
- [ ] Run current parser characterization tests before moving code:

```sh
go test ./internal/logging -run 'Test.*EnvLevel'
```

- [ ] Add equivalent CLI-boundary tests for named/numeric values and a two-call `t.Setenv` isolation case; keep these tests initially disabled only by placing them beside the new API implementation step, not by asserting nonexistent symbols.
- [ ] Implement Viper-backed lookup/conversion under `cmd/dum/cli`; fresh Viper state per call; preserve existing `internal/logging.EnvLevel` semantics.
- [ ] Run CLI log-level tests after API introduction; expect PASS:

```sh
go test ./cmd/dum/... -run 'Test.*LogLevel'
```

## Task 5: Wire root logger and remove internal env lookup

**Files:** `cmd/dum/root.go`, `cmd/dum/root_test.go`, `cmd/dum/log/log.go`, `cmd/dum/log/*_test.go`, `internal/logging/env_level.go`, `internal/logging/logging_test.go`

- [ ] Update `NewDum()` to resolve `ZSH_LOG_LEVEL` through `cmd/dum/cli` before calling `lg.NewLogger`.
- [ ] Remove `lg.EnvLevel()` call from `cmd/dum/log.executeLogging`.
- [ ] Remove `internal/logging.EnvLevel` and obsolete internal tests; retain logger API tests and migrated CLI-level tests.
- [ ] Add self-exec subprocess coverage in `cmd/dum/root_test.go`: select helper via `-test.run`, launch fresh children with `ZSH_LOG_LEVEL=debug` and `ZSH_LOG_LEVEL=warn`, invoke non-terminating `log debug`/`log warn`, capture stderr, assert debug output appears only in debug child.
- [ ] Run focused command/log tests:

```sh
go test ./cmd/dum/... -run 'Test(NewDum|.*Log|.*Level)'
```

- [ ] Commit:

```sh
git add cmd/dum internal/logging
git commit -m "refactor: move log env config to cmd"
```

## Task 6: Full verification and cleanup

**Files:** temporary fixture/files plus any files exposed by tests, lint, or smoke checks

- [ ] Create temporary minimal valid `installer.yml` fixture with one harmless enabled play/task; remove it after smoke tests.
- [ ] Run package tests:

```sh
go test ./cmd/dum/...
```

- [ ] Run full tests:

```sh
go test ./...
```

- [ ] Run lint:

```sh
make lint
```

- [ ] Build binary:

```sh
make build
```

- [ ] Run `dum --help`; expect root help and all command names.
- [ ] Run `dum list --file "$fixture"`; expect fixture play/task output.
- [ ] Run `INSTALLER_CONFIG="$fixture" dum list`; expect same fixture output.
- [ ] Run `INSTALLER_CONFIG="$fixture" dum install --dry-run`; expect task preview without mutations.
- [ ] Run `ZSH_LOG_LEVEL=debug dum log debug "probe"`; expect `probe` output.
- [ ] Run rename against temporary files with `--dry-run`; expect preview only.
- [ ] Run `dum schema --output "$tmp/installer.schema.json"`; expect schema file creation.
- [ ] Remove temporary fixture/files.
- [ ] Confirm no `internal/` package imports Viper and no obsolete `EnvLevel` references remain.
- [ ] Commit any required cleanup with a focused Conventional Commit message.

## Delivery

Single cohesive change; use one PR targeting `main`. Before handoff:

- [ ] Push branch.
- [ ] Open PR using `stacked-prs` skill if repository workflow requires it.
