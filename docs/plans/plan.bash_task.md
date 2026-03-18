# BashTask Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a new `bash` task type that executes shell scripts or inline bash commands, similar to GitHub Actions `run` step.

**Architecture:** BashTask follows the same Task interface pattern as other tasks (DirTask, BrewTask). It validates exactly one of `Script` or `Run` is set, resolves script paths relative to the playbook file, and executes commands via `bash`. Inline commands are written to a temp file and executed.

**Tech Stack:** Go 1.26+, go-enum for enum generation, testify for testing/mocking

---

## File Structure

| File | Action |
|------|--------|
| `pkg/enums/task_type.go` | Add `bash` to enum directive |
| `pkg/enums/task_type_enum.go` | Regenerate with `make generate` |
| `pkg/playbook/bash_task.go` | Create - main task implementation |
| `pkg/playbook/bash_task_test.go` | Create - unit tests |
| `pkg/playbook/factory.go` | Add `case e.TaskTypeBash:` in switch |
| `pkg/testdata/test_installer.yml` | Add bash task test data |
| `pkg/playbook/playbook_test.go` | Add `createTestBashTask` helper |

---

## Task 1: Add `bash` to TaskType Enum

**Files:**
- Modify: `pkg/enums/task_type.go:5`
- Output: `pkg/enums/task_type_enum.go` (regenerated)

- [ ] **Step 1: Add `bash` to the enum directive**

Modify `pkg/enums/task_type.go` line 5:
```go
// ENUM(cask, cellar, brew, dir, function,  git, link, jetbrains, mas, vscode, bash).
```

- [ ] **Step 2: Regenerate the enum**

Run: `make generate`
Expected: `pkg/enums/task_type_enum.go` updated with `TaskTypeBash` constant

- [ ] **Step 3: Commit**

```bash
git add pkg/enums/task_type.go pkg/enums/task_type_enum.go
git commit -m "feat: add bash to TaskType enum"
```

---

## Task 2: Create BashTask Implementation

**Files:**
- Create: `pkg/playbook/bash_task.go`

- [ ] **Step 1: Write the BashTask struct and constructor**

```go
package playbook

import (
	"awong/dotfiles/pkg/external"
	"awong/dotfiles/pkg/logging"
	"context"
	"fmt"
	"os"
	"os/exec"
)

// BashTask executes bash scripts or inline commands.
type BashTask struct {
	Attributes
	Script string
	Run    string
	Utils  external.Ext
	Log    logging.Logger
}

// NewBashTask constructs a BashTask.
func NewBashTask(attributes *Attributes, script, run string) (*BashTask, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}
	if attributes.ID == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}
	if script == "" && run == "" {
		return nil, fmt.Errorf("task %s: either 'script' or 'run' must be provided", attributes.ID)
	}
	if script != "" && run != "" {
		return nil, fmt.Errorf("task %s: cannot specify both 'script' and 'run'", attributes.ID)
	}
	return &BashTask{
		Attributes: *attributes,
		Script:     script,
		Run:        run,
		Utils:      external.NewExt(),
		Log:        logging.Log(),
	}, nil
}

// GetAttributes implements Attributable.
func (t *BashTask) GetAttributes() Attributes {
	return t.Attributes
}

// GetID implements Identifiable.
func (t *BashTask) GetID() string {
	return t.ID
}

// IsEnabled implements Enableable.
func (t *BashTask) IsEnabled() bool {
	return t.Enabled
}

// List implements Lister.
func (t *BashTask) List(_ context.Context, input *Input) (*TaskResult, error) {
	var cmd string
	if t.Script != "" {
		cmd = fmt.Sprintf("bash %s", t.Script)
	} else {
		cmd = fmt.Sprintf("bash (inline): %s", t.Run)
	}
	err := t.Log.Printlnf("%v %s", TaskEllipsis, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to list bash task: %v", err)
	}
	return t.CreateTaskResult(input, true)
}

// Install implements Installer.
func (t *BashTask) Install(ctx context.Context, input *Input) (*TaskResult, error) {
	t.Log.Debugf("%s START ... play: %s taskID: %s", TaskEllipsis, input.Play, t.ID)
	if !t.Enabled {
		return t.CreateTaskResult(input, false)
	}

	var scriptPath string
	if t.Script != "" {
		absPath, err := t.Utils.ToAbsolutePath(t.Script)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve script path %s: %w", t.Script, err)
		}
		scriptPath = absPath
		if !t.Utils.IsDir(scriptPath) && !fileExists(scriptPath) {
			return nil, fmt.Errorf("script file not found: %s", scriptPath)
		}
	} else {
		tmpFile, err := os.CreateTemp("", "dum-bash-*.sh")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp file: %w", err)
		}
		defer os.Remove(tmpFile.Name())
		if _, err := tmpFile.WriteString(t.Run); err != nil {
			return nil, fmt.Errorf("failed to write to temp file: %w", err)
		}
		if err := tmpFile.Close(); err != nil {
			return nil, fmt.Errorf("failed to close temp file: %w", err)
		}
		scriptPath = tmpFile.Name()
	}

	t.Log.Infof("%s %s: bash %s", TaskEllipsis, input.Play, scriptPath)
	if !input.DryRun {
		cmd := exec.CommandContext(ctx, "bash", scriptPath)
		if t.Sudo {
			cmd = exec.CommandContext(ctx, "sudo", "bash", scriptPath)
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("bash command failed with exit code %d: %s: %w",
				cmd.ProcessState.ExitCode(), string(output), err)
		}
	}
	return t.CreateTaskResult(input, true)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
```

- [ ] **Step 2: Commit**

```bash
git add pkg/playbook/bash_task.go
git commit -m "feat: implement BashTask for executing bash scripts and inline commands"
```

---

## Task 3: Wire BashTask into Factory

**Files:**
- Modify: `pkg/playbook/factory.go:167-198`

- [ ] **Step 1: Add `case e.TaskTypeBash:` to the switch statement**

Add after the `TaskTypeCellar` case (around line 195):

```go
	case e.TaskTypeBash:
		return NewBashTask(attributes,
			f.getString(yml, "script", ""),
			f.getString(yml, "run", ""))
```

- [ ] **Step 2: Commit**

```bash
git add pkg/playbook/factory.go
git commit -m "feat: wire BashTask into factory"
```

---

## Task 4: Add Test Helper and Test Data

**Files:**
- Modify: `pkg/playbook/playbook_test.go`
- Modify: `pkg/testdata/test_installer.yml`

- [ ] **Step 1: Add `createTestBashTask` helper to playbook_test.go**

Add after the `createTestVsCodeTask` function (around line 118):

```go
func createTestBashTask(t *testing.T, attr Attributes, script, run string) *BashTask {
	task, err := NewBashTask(&attr, script, run)
	assert.NoError(t, err)
	return task
}
```

- [ ] **Step 2: Add bash task test data to test_installer.yml**

Add to play-1 tasks (after line 46, before play-2):

```yaml
        - type: "bash"
          id: "test-bash-script"
          description: "test bash task"
          script: "./test_script.sh"
        - type: "bash"
          id: "test-bash-inline"
          description: "test bash inline"
          run: "echo 'hello from bash'"
```

- [ ] **Step 3: Commit**

```bash
git add pkg/playbook/playbook_test.go pkg/testdata/test_installer.yml
git commit -m "feat: add BashTask test helpers and test data"
```

---

## Task 5: Write BashTask Tests

**Files:**
- Create: `pkg/playbook/bash_task_test.go`

- [ ] **Step 1: Write tests for BashTask**

```go
package playbook

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findBashTask(t *testing.T, input *Input, id string) *BashTask {
	tk, err := findTask(input.PlayBook, id)
	assert.NoError(t, err)
	task, ok := tk.(*BashTask)
	assert.True(t, ok)
	return task
}

func TestBashTask_Install_Script(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBashTask(t, input, "test-bash-script")
	ctx := context.Background()

	tmpDir := t.TempDir()
	scriptPath := tmpDir + "/test_script.sh"
	err := os.WriteFile(scriptPath, []byte("echo 'hello'"), 0755)
	assert.NoError(t, err)

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().ToAbsolutePath(task.Script).Return(scriptPath, nil).Once()
	mockUtils.EXPECT().IsDir(scriptPath).Return(false).Once()
	task.Utils = mockUtils
	task.Script = scriptPath

	got, err := task.Install(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.True(t, got.Success)
}

func TestBashTask_Install_Inline(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBashTask(t, input, "test-bash-inline")
	ctx := context.Background()

	got, err := task.Install(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.True(t, got.Success)
}

func TestBashTask_DryRun(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = true
	task := findBashTask(t, input, "test-bash-inline")
	ctx := context.Background()

	got, err := task.Install(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.True(t, got.Success)
	assert.True(t, got.DryRun)
}

func TestBashTask_Disabled(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBashTask(t, input, "test-bash-inline")
	task.Enabled = false
	ctx := context.Background()

	got, err := task.Install(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.False(t, got.Success)
}

func TestBashTask_NewBashTask_Error_NeitherScriptNorRun(t *testing.T) {
	attr := &Attributes{
		ID:          "test",
		Description: "test",
		Enabled:     true,
		Sudo:        false,
	}
	_, err := NewBashTask(attr, "", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "either 'script' or 'run' must be provided")
}

func TestBashTask_NewBashTask_Error_BothScriptAndRun(t *testing.T) {
	attr := &Attributes{
		ID:          "test",
		Description: "test",
		Enabled:     true,
		Sudo:        false,
	}
	_, err := NewBashTask(attr, "script.sh", "echo hello")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot specify both 'script' and 'run'")
}

func TestBashTask_ScriptNotFound(t *testing.T) {
	input := createTestInput(t)
	input.DryRun = false
	task := findBashTask(t, input, "test-bash-script")
	ctx := context.Background()

	mockUtils := NewMockExt(t)
	mockUtils.EXPECT().ToAbsolutePath(task.Script).Return("/nonexistent/path.sh", nil).Once()
	task.Utils = mockUtils
	task.Script = "/nonexistent/path.sh"

	got, err := task.Install(ctx, input)
	assert.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "script file not found")
}

func TestBashTask_List(t *testing.T) {
	input := createTestInput(t)
	task := findBashTask(t, input, "test-bash-inline")
	ctx := context.Background()

	got, err := task.List(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.True(t, got.Success)
}
```

- [ ] **Step 2: Run tests to verify they pass**

Run: `go test -v ./pkg/playbook/... -run TestBashTask`
Expected: PASS for all tests

- [ ] **Step 3: Commit**

```bash
git add pkg/playbook/bash_task_test.go
git commit -m "test: add BashTask unit tests"
```

---

## Task 6: Run Full Test Suite and Lint

**Files:**
- None (verification only)

- [ ] **Step 1: Run full test suite**

Run: `make test`
Expected: All tests pass

- [ ] **Step 2: Run linter**

Run: `make lint`
Expected: No lint errors

- [ ] **Step 3: Commit (if any fixes were needed)**

```bash
git add -A
git commit -m "fix: address test/lint issues"
```

---

## Summary

After completing all tasks:
- `pkg/enums/task_type.go` - Added `bash` to enum
- `pkg/enums/task_type_enum.go` - Regenerated
- `pkg/playbook/bash_task.go` - New BashTask implementation
- `pkg/playbook/bash_task_test.go` - Unit tests
- `pkg/playbook/factory.go` - Wired into factory
- `pkg/playbook/playbook_test.go` - Added helper function
- `pkg/testdata/test_installer.yml` - Added test data

The feature is complete when all tests pass and lint is clean.
