# Task Refactor - Updated Implementation Plan

## Current State

### Completed
- ✅ Created `pkg/types/attributes.go` with shared types
- ✅ Created `pkg/tasks/` with 17 task files moved from playbook
- ✅ Updated task imports from `playbook` to `types`
- ⚠️ `pkg/types/attributes.go` has **duplicate constants** (lines 8-12 and 27-31)

### Remaining

#### 1. Fix duplicate constants (5 min)
- [ ] Remove lines 27-31 from `pkg/types/attributes.go`

#### 2. Verify tasks package builds (10 min)
- [ ] Run `go build ./pkg/tasks/...`
- [ ] Fix any remaining import issues

#### 3. Update playbook package (30 min)
- [ ] `playbook.go` - import from `types`
- [ ] `play.go` - import from `types` (also consider moving to `pkg/play/`)
- [ ] `factory.go` - update imports, construct from `tasks` package
- [ ] `executor.go` - import from `types`
- [ ] `result.go` - check if still needed or can be removed
- [ ] `play_result.go` - check if still needed or can be removed

#### 4. Delete old type files from playbook (5 min)
- [ ] Verify these are already deleted: `attributes.go`, `input.go`, `task_result.go`, `task.go`
- [ ] Remove any remaining duplicates

#### 5. Move play orchestration to pkg/play/ (optional, 30 min)
- [ ] Create `pkg/play/` directory
- [ ] Move `play.go` → `pkg/play/play.go`
- [ ] Move `play_result.go` → `pkg/play/result.go`
- [ ] Update imports in `playbook/`

#### 6. Verify full build and tests (15 min)
- [ ] `go build ./...`
- [ ] `go test ./...`

## Files Reference

### pkg/types/attributes.go (needs fix)
- TaskEllipsis, PlayEllipsis, PlayBookEllipsis (DUPLICATE - remove lines 27-31)
- Lister, Installer, Task interfaces
- Attributes struct + methods
- TaskResult struct + NewTaskResult, CreateTaskResult
- PlayBook struct
- Input struct + NewInput

### pkg/tasks/ (18 files - should build after fix)
- bash_task.go, bash_installer.go
- brew_task.go, brew_cask_task.go, brew_cellar_task.go, brew_tap.go
- dir_task.go, function_task.go, git_task.go
- jetbrains_plugin_task.go, link_task.go, mas_task.go
- sdkman_installer.go, starship_installer.go, vim_installer.go
- vscode_plugin_task.go
- mount_installer.go, test_installer.go

### pkg/playbook/ (needs updates)
- playbook.go, play.go, factory.go, executor.go
- result.go, play_result.go
- Tests still in playbook/ (may need moving too)

## Dependency Goal
```
cmd → playbook → tasks + play
                  ↑        ↑
           (pkg/types)
```

## Next Step
Start with fixing the duplicate constants, then verify tasks builds.
