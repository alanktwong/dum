# CLI Boundary Refactor Design

**Date**: 2026-08-13

## Goal

Move all Cobra-facing code from `internal/cmd` into `cmd/dum`, while keeping
installer behavior unchanged. The executable layer will call existing internal
domain packages through typed structs and errors.

## Boundaries

`cmd/dum` owns root and subcommand construction, flags, positional arguments,
help text, Cobra callbacks, and CLI-specific error/log handling. Internal
packages must not import Cobra or expose Cobra command types.

Existing internal packages retain domain ownership:

- `internal/factory`: typed installer options and configuration conversion
- `internal/playbook`: install/list orchestration
- `internal/logging`: logging operations
- `internal/yaml`: typed configuration loading and validation
- existing or narrowly scoped schema/rename packages for those operations

No centralized service package will be introduced unless an operation has no
appropriate existing domain owner.

## Layout

```text
cmd/dum/
  main.go
  root.go
  install/
  list/
  log/
  rename/
  schema/
```

Each command package exposes a Cobra constructor and converts CLI inputs into
internal request structs. Internal calls return typed results or errors.

## Migration

1. Map each `internal/cmd` file to its command package or domain package.
2. Move Cobra construction, flags, and help text to `cmd/dum`.
3. Extract callbacks into typed internal calls.
4. Remove Cobra dependencies from internal code.
5. Move and adapt tests with the commands.
6. Update `cmd/dum/main.go` to construct the new root command.
7. Migrate Makefile and all GoReleaser linker targets from `internal/cmd` to
   the new `cmd/dum` package, preserving version and commit injection.
8. Remove `internal/cmd` after all references migrate.

## Compatibility and verification

Command names, aliases, flags, help text, configuration behavior, logging,
rename behavior, schema output, and release version reporting remain unchanged.
Verify with package tests, `go test ./...`, `make lint`, a binary build, a
release-style build/`dum --version`, and CLI smoke tests for root help and
representative install/list/log/rename/schema commands.
