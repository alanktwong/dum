# CLI Boundary Refactor Implementation Plan

> **For agentic workers:** Execute task-by-task with verification after each migration stage.

**Goal:** Move all Cobra integration from `internal/cmd` into `cmd/dum` while preserving the CLI contract.

**Architecture:** `cmd/dum` owns Cobra commands, flags, help, and CLI callbacks. Existing internal domain packages own typed operations and do not import Cobra. Migrate command tests with their implementations, then delete `internal/cmd`.

**Tech Stack:** Go 1.26, Cobra, existing internal factory/playbook/logging/yaml packages, GoReleaser.

---

### Task 1: Inventory command boundaries and establish root package

**Files:**
- Create/modify: `cmd/dum/main.go`, `cmd/dum/root.go`
- Reference: `internal/cmd/dum.go`, `internal/cmd/cmd.go`
- Test: `cmd/dum/root_test.go`

- [ ] Define the root command and typed dependency wiring under `cmd/dum`.
- [ ] Move version and commit linker variables to `cmd/dum`.
- [ ] Preserve default config lookup and shared flag helpers in the CLI layer.
- [ ] Update `main.go` to use the new root constructor.
- [ ] Verify root help and version behavior.

### Task 2: Migrate install and list commands

**Files:**
- Create/modify: `cmd/dum/install/`, `cmd/dum/list/`
- Reference: `internal/cmd/install.go`, `internal/cmd/list.go`
- Tests: corresponding command tests
- Internal APIs: `internal/factory`, `internal/playbook`

- [ ] Move Cobra construction, flags, examples, and argument handling.
- [ ] Convert CLI values into typed factory/playbook options.
- [ ] Remove Cobra command dependencies from the called internal operations.
- [ ] Move and adapt install/list tests.
- [ ] Verify command help and focused tests.

### Task 3: Migrate log commands

**Files:**
- Create/modify: `cmd/dum/log/`
- Reference: `internal/cmd/log.go`, `debug.go`, `info.go`, `success.go`, `warn.go`, `error.go`
- Tests: corresponding log command tests
- Internal API: `internal/logging`

- [ ] Move all log subcommand constructors and help text.
- [ ] Parse prefix and arguments in `cmd/dum/log`.
- [ ] Call logging through typed values and preserve level behavior.
- [ ] Move/adapt tests and verify all levels.

### Task 4: Migrate rename and schema commands

**Files:**
- Create/modify: `cmd/dum/rename/`, `cmd/dum/schema/`
- Reference: `internal/cmd/rename.go`, `schema.go`, `schema_embed.go`
- Tests: corresponding command tests
- Internal APIs: existing rename/schema domain code or focused internal packages

- [ ] Move Cobra construction and flag parsing.
- [ ] Extract rename and schema work behind typed internal calls.
- [ ] Preserve output paths, dry-run, casing, replacement, and schema behavior.
- [ ] Move/adapt tests and verify focused commands.

### Task 5: Remove obsolete internal command package

**Files:**
- Delete: `internal/cmd/*.go`
- Modify: all remaining imports and package references
- Modify: `Makefile`, `cfg/goreleaser.yaml`, `cfg/goreleaser.local.yaml`, `cfg/goreleaser.snapshot.yaml`
- Modify: `cfg/testcoverage.yml`

- [ ] Remove all Cobra imports from `internal/`.
- [ ] Update linker targets from `internal/cmd` to `cmd/dum`.
- [ ] Remove obsolete coverage overrides and references.
- [ ] Delete `internal/cmd` after all callers migrate.
- [ ] Verify no stale references remain.

### Task 6: Full verification and cleanup

**Files:**
- Modify only if verification exposes migration issues.

- [ ] Run `gofmt` and generation as required.
- [ ] Run `go test ./...`.
- [ ] Run `make lint`.
- [ ] Build the binary and verify `dum --help`, `dum --version`, and representative subcommands.
- [ ] Confirm working tree contains only intentional changes.
