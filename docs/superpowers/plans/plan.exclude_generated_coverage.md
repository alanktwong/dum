# Plan: Exclude Generated Code from Test Coverage

**Date**: 2026-03-25
**Design**: `../designs/design.exclude_generated_coverage.md`
**Approach**: B - Separate Packages

## Goal

Move generated code (`*_enum.go`, `*_mocks.go`) and their source type definitions into `gen/` subpackages. This allows excluding entire packages from test coverage at collection time.

---

## TODOs

### Phase 1: Types Package

- [ ] **1.1** Create `pkg/types/gen/` directory
- [ ] **1.2** Move `task_type.go` to `pkg/types/gen/`
- [ ] **1.3** Move `task_type_enum.go` to `pkg/types/gen/`
- [ ] **1.4** Move `jetbrains_type.go` to `pkg/types/gen/`
- [ ] **1.5** Move `jetbrains_type_enum.go` to `pkg/types/gen/`
- [ ] **1.6** Update go-enum directive in `task_type.go` to output to `gen/`
- [ ] **1.7** Update go-enum directive in `jetbrains_type.go` to output to `gen/`
- [ ] **1.8** Update imports in `tasks_factory.go` to use `types/gen`
- [ ] **1.9** Run `go generate` to verify enum generation works
- [ ] **1.10** Run tests to verify everything compiles

### Phase 2: Plays Package

- [ ] **2.1** Create `pkg/plays/gen/` directory
- [ ] **2.2** Move `playbook_info.go` to `pkg/plays/gen/`
- [ ] **2.3** Update imports in `pkg/plays/` that reference `PlayBookInfo`
- [ ] **2.4** Update mockery config to generate mock to `gen/`
- [ ] **2.5** Run mockery to generate `playbook_info_mocks.go` in `gen/`
- [ ] **2.6** Run tests to verify everything compiles

### Phase 3: Logging Package

- [ ] **3.1** Create `pkg/logging/gen/` directory
- [ ] **3.2** Move `logging.go` to `pkg/logging/gen/`
- [ ] **3.3** Update imports in `pkg/logging/` that reference `Logger`
- [ ] **3.4** Update mockery config to generate mock to `gen/`
- [ ] **3.5** Run mockery to generate `logging_mocks.go` in `gen/`
- [ ] **3.6** Run tests to verify everything compiles

### Phase 4: External Package (mocks only)

- [ ] **4.1** Create `pkg/external/gen/` directory
- [ ] **4.2** Identify which interfaces to move (Brew, Code, Git, etc.)
- [ ] **4.3** Move interface definitions to `pkg/external/gen/`
- [ ] **4.4** Update imports throughout codebase
- [ ] **4.5** Update mockery configs to generate to `gen/`
- [ ] **4.6** Run tests to verify everything compiles

### Phase 5: Factory Package (mocks only)

- [ ] **5.1** Create `pkg/factory/gen/` directory
- [ ] **5.2** Move interfaces used for mocking to `gen/`
- [ ] **5.3** Update imports throughout codebase
- [ ] **5.4** Update mockery configs to generate to `gen/`
- [ ] **5.5** Run tests to verify everything compiles

### Phase 6: Tasks Package (mocks only)

- [ ] **6.1** Create `pkg/tasks/gen/` and `pkg/tasks/installer/gen/` directories
- [ ] **6.2** Move interfaces used for mocking to `gen/`
- [ ] **6.3** Update imports throughout codebase
- [ ] **6.4** Update mockery configs to generate to `gen/`
- [ ] **6.5** Run tests to verify everything compiles

### Phase 7: Verification

- [ ] **7.1** Run `make coverage` to verify generated code is excluded
- [ ] **7.2** Run full test suite to ensure nothing is broken
- [ ] **7.3** Update coverage config if needed

---

## Notes

- This plan executes phases sequentially, testing after each phase
- If at any point tests fail, fix before proceeding
- Coverage should improve after this change since generated code won't dilute coverage metrics