# BashTask Design

## Overview

Add a new `bash` task type that executes shell scripts or inline bash commands, similar to GitHub Actions `run` step.

## YAML Schema

```yaml
tasks:
  - id: install-my-tool
    type: bash
    description: Install my custom tool
    enabled: true
    sudo: false
    script: ./scripts/install.sh    # file reference (xor with 'run')
    # OR
    run: |
      echo "hello"
      curl -s https://example.com | bash
```

**Constraints:**

- Exactly one of `script` OR `run` must be provided
- `script` is relative to the playbook file location
- `run` is multiline bash commands written to a temp file and executed

## Config File Location (Already Implemented)

The playbook file location is handled by the CLI layer (`pkg/cmd/cmd.go`):

| Priority | Source | Example |
|----------|--------|---------|
| 1 | `-f/--file` flag | `dum install -f ./my-playbook.yml` |
| 2 | `INSTALLER_CONFIG` env var | `INSTALLER_CONFIG=./my.yml dum install` |
| 3 | XDG default | `~/.config/installer.yml` |

No changes needed to support this feature.

## Files to Create/Modify


| File | Action |
|------|--------|
| `pkg/enums/task_type.go` | Add `bash` to enum |
| `pkg/playbook/bash_task.go` | New - main task implementation |
| `pkg/playbook/bash_task_test.go` | New - unit tests |
| `pkg/playbook/factory.go` | Add `case e.TaskTypeBash:` in switch |
| `pkg/enums/task_type_enum.go` | Regenerate with `make generate` |

## Struct Design

```go
type BashTask struct {
    Attributes
    Script    string  // file path OR
    Run       string  // inline commands (mutually exclusive)
    Utils     external.Ext
    Log       logging.Logger
}
```

## Execution Flow

1. `Install()` validates exactly one of `Script`/`Run` is set
2. If `Script` provided:
   - Resolve path relative to playbook file
   - Execute: `bash <resolved-path>`
3. If `Run` provided:
   - Write to temp file: `os.CreateTemp("", "dum-bash-*.sh")`
   - Execute: `bash <temp-file>`
   - Cleanup temp file after execution
4. Respect `input.DryRun` - log command instead of executing

## Error Handling

- Missing both `script` and `run` → error
- Both `script` and `run` provided → error
- Script file not found → error with resolved path
- Command fails → wrap error with exit code

## Testing Strategy

- Happy path: script execution
- Happy path: inline command execution
- Error: neither script/run provided
- Error: both provided
- Error: script file not found
- Dry-run mode verification