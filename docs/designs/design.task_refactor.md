# Task Refactor

## Problem Statement

The `playbook` package contains ~50 files that mix concerns:
- Task types (brew_task.go, bash_task.go, etc.)
- Play/Playbook orchestration (play.go, playbook.go)
- Factory construction logic (factory.go)
- Execution logic (executor.go)
- Shared types (attributes.go, input.go, result.go)

This makes the codebase harder to navigate and reason about. A developer working on tasks shouldn't need to understand play orchestration, and vice versa.

## Current State

Analysis shows no circular dependencies in the current codebase:
- `playbook` → `external` + `logging`
- `cmd` → `playbook`

The refactor must preserve this property.

## Goals

1. **Improve code organization** - separate tasks from play orchestration
2. **Maintain zero circular dependencies** - new packages should not create cycles
3. **Preserve all existing functionality** - no behavior changes

## Target Package Structure

```markdown
pkg/
├── playbook/           (stays as-is, becomes orchestration layer)
│   ├── playbook.go
│   ├── play.go         (moves to pkg/play/)
│   ├── play_result.go  (moves to pkg/play/)
│   ├── executor.go
│   ├── factory.go      (refactored to use new packages)
│   └── ...
├── tasks/              (NEW - all task types)
│   ├── task.go         (interface - STAYS, referenced by play)
│   ├── bash_task.go
│   ├── bash_installer.go
│   ├── brew_task.go
│   ├── brew_cask_task.go
│   ├── brew_cellar_task.go
│   ├── brew_tap.go
│   ├── dir_task.go
│   ├── function_task.go
│   ├── git_task.go
│   ├── jetbrains_plugin_task.go
│   ├── link_task.go
│   ├── mas_task.go
│   ├── sdkman_installer.go
│   ├── starship_installer.go
│   ├── vim_installer.go
│   └── vscode_plugin_task.go
├── play/               (NEW - play orchestration)
│   ├── play.go
│   └── play_result.go
└── cmd/                (unchanged)
```

## Shared Types Location

The following types are used by multiple packages and must remain accessible:
- `Attributes` - currently in `playbook/attributes.go`, consider moving to `pkg/types/`
- `Task` interface - currently in `playbook/task.go`, must be in shared location
- `TaskResult` - currently in `playbook/task_result.go`
- `Input` - currently in `playbook/input.go`

**Recommendation:** Create `pkg/types/` for shared types to avoid circular imports between `tasks` and `play`.

## Factory Refactoring

`factory.go` currently constructs both tasks and plays. After refactor:
- Stays in `playbook` package (orchestration layer)
- Imports `tasks` and `play` packages to construct instances
- Single factory is fine - no need to split into task_factory/play_factory

## Dependency Rules

```markdown
cmd → playbook → tasks + play
                ↑        ↑
         (shared types in pkg/types)
```

No package should import another that imports it back.

## Success Criteria

1. `go build ./...` passes
2. `go test ./...` passes  
3. No new circular dependencies introduced
4. All task types work identically to before
5. Import paths update correctly in `cmd/` package

## Files to Move

### To pkg/tasks/
- bash_task.go
- bash_installer.go
- brew_task.go
- brew_cask_task.go
- brew_cellar_task.go
- brew_tap.go
- dir_task.go
- function_task.go
- git_task.go
- jetbrains_plugin_task.go
- link_task.go
- mas_task.go
- sdkman_installer.go
- starship_installer.go
- vim_installer.go
- vscode_plugin_task.go

### To pkg/play/
- play.go
- play_result.go

### New pkg/types/ (recommended)
- attributes.go (or contents merged into new file)
- task.go
- task_result.go

## Implementation Order

1. Create `pkg/types/` with shared types
2. Create `pkg/tasks/`, move task files, update imports
3. Create `pkg/play/`, move play files, update imports
4. Update `factory.go` to use new package paths
5. Run tests and fix any import issues

## Out of Scope

- Changes to task behavior
- Adding new task types
- Modifying the installer.yml schema
