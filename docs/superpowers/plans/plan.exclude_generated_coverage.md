# Plan: Exclude Generated Code from Test Coverage

**Date**: 2026-03-25
**Design**: `../designs/design.exclude_generated_coverage.md`
**Approach**: B - Separate Packages
**Base PR**: [#5](https://github.com/alanktwong/dum/pull/5)

## Goal

Move generated code (`*_enum.go`, `*_mocks.go`) and their source type definitions into `gen/` subpackages. This allows excluding entire packages from test coverage at collection time.

## Stacked PRs Strategy

Each phase creates a PR that stacks onto the previous one, ultimately stacking on PR #5.
- Phase 1 → PR stacking on #5
- Phase 2 → PR stacking on Phase 1's PR
- And so on...

---

## TODOs

### Phase 1: Types Package (Stacked on #5)

**PR**: [#6](https://github.com/alanktwong/dum/pull/6)
**PR Title**: `refactor enums in types package to gen/ subpackage`

- [x] **1.1** Create `pkg/types/gen/` directory
- [x] **1.2** Move `task_type.go` to `pkg/types/gen/`
- [x] **1.3** Move `task_type_enum.go` to `pkg/types/gen/`
- [x] **1.4** Move `jetbrains_type.go` to `pkg/types/gen/`
- [x] **1.5** Move `jetbrains_type_enum.go` to `pkg/types/gen/`
- [x] **1.6** Update go-enum directive in `task_type.go` to output to `gen/`
- [x] **1.7** Update go-enum directive in `jetbrains_type.go` to output to `gen/`
- [x] **1.8** Update imports in `tasks_factory.go` to use `types/gen`
- [x] **1.9** Run `go generate` to verify enum generation works
- [x] **1.10** Run tests to verify everything compiles

---

### Phase 2: Plays Package

**PR**: [#7](https://github.com/alanktwong/dum/pull/7)
**PR Title**: `refactor plays package: move PlayBookInfo to gen/ subpackage`

- [x] **2.1** Create `pkg/plays/gen/` directory
- [x] **2.2** Move `playbook_info.go` to `pkg/plays/gen/`
- [x] **2.3** Update imports in `pkg/plays/` that reference `PlayBookInfo`
- [x] **2.4** Update mockery config to generate mock to `gen/`
- [x] **2.5** Run mockery to generate `playbook_info_mocks.go` in `gen/`
- [x] **2.6** Run tests to verify everything compiles

---

### Phase 3: Logging Package

**PR**: [#8](https://github.com/alanktwong/dum/pull/8)
**PR Title**: `refactor logging package: move Logger interface to gen/ subpackage`

- [x] **3.1** Create `pkg/logging/gen/` directory
- [x] **3.2** Create `logger.go` interface in `pkg/logging/gen/`
- [x] **3.3** Update `logging.go` to re-export from `gen/`
- [x] **3.4** Update mockery config to generate mock to `gen/`
- [x] **3.5** Run mockery to generate `logging_mocks.go` in `gen/`
- [x] **3.6** Run tests to verify everything compiles

---

### Phase 4: External Package (mocks only)

**PR**: TBD
**PR Title**: `refactor external package: move its interfaces to gen/`

- [ ] **4.1** Create `pkg/external/gen/` directory
- [ ] **4.2** Identify which interfaces to move (Brew, Code, Git, etc.)
- [ ] **4.3** Move interface definitions to `pkg/external/gen/`
- [ ] **4.4** Update imports throughout codebase
- [ ] **4.5** Update mockery configs to generate to `gen/`
- [ ] **4.6** Run tests to verify everything compiles

---

### Phase 5: Factory Package (mocks only)

**PR**: TBD
**PR Title**: `refactor factory package: move its interfaces to gen/`

- [ ] **5.1** Create `pkg/factory/gen/` directory
- [ ] **5.2** Move interfaces used for mocking to `gen/`
- [ ] **5.3** Update imports throughout codebase
- [ ] **5.4** Update mockery configs to generate to `gen/`
- [ ] **5.5** Run tests to verify everything compiles

---

### Phase 6: Tasks Package (mocks only)

**PR**: TBD
**PR Title**: `refactor tasks package: move interfaces to gen/`

- [ ] **6.1** Create `pkg/tasks/gen/` and `pkg/tasks/installer/gen/` directories
- [ ] **6.2** Move interfaces used for mocking to `gen/`
- [ ] **6.3** Update imports throughout codebase
- [ ] **6.4** Update mockery configs to generate to `gen/`
- [ ] **6.5** Run tests to verify everything compiles

---

### Phase 7: Verification

**PR**: TBD
**PR Title**: `verify coverage exclusion works`

- [ ] **7.1** Run `make coverage` to verify generated code is excluded
- [ ] **7.2** Run full test suite to ensure nothing is broken
- [ ] **7.3** Update coverage config if needed

---

## Notes

- This plan executes phases sequentially, testing after each phase
- If at any point tests fail, fix before proceeding
- Coverage should improve after this change since generated code won't dilute coverage metrics