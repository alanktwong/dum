# Symbolic tarball target default in help Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use subagent-driven-development (recommended) or
> executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for
> tracking.

**Goal:** Show Cobra's symbolic tarball target fallback in backup and restore help without changing
the runtime target value or environment precedence.

**Architecture:** Register each `--target` flag with the current runtime default from
`defaultTarballTarget()`, then set only pflag's `DefValue` help metadata to the static fallback
description. Add a CLI test proving both commands display the symbolic chain while their runtime
flag values still honor environment overrides.

**Tech Stack:** Go, Cobra/pflag, testify/assert.

---

## PR Strategy

Single cohesive change; no stacked PRs.

## File Structure

- Modify `cmd/dum/tarball.go`: define the symbolic help string, register both `--target` flags through
  one helper, and leave `defaultTarballTarget()` unchanged.
- Modify `cmd/dum/tarball_test.go`: cover symbolic help for backup and restore and verify it does not
  replace the runtime flag value.

No changes to `internal/tarball`, the README, or generated schema are needed.

## Task Structure

### Task 1: Add failing help regression test

**Files:**

- Modify: `cmd/dum/tarball_test.go`

- [x] **Step 1: Add test for both command help outputs**

Add this test near `TestTarballCommand_TargetDefaultPrecedence`:

```go
func TestTarballCommand_TargetHelpUsesSymbolicDefault(t *testing.T) {
	const expectedHelpDefault =
		"TGZ_TARGET, else $XDG_DATA_HOME/tarballs, else ~/.local/share/tarballs"

	home := t.TempDir()
	dataHome := filepath.Join(t.TempDir(), "xdg-data")
	envTarget := filepath.Join(t.TempDir(), "env-target")
	t.Setenv("HOME", home)
	t.Setenv("XDG_DATA_HOME", dataHome)
	t.Setenv("TGZ_TARGET", envTarget)

	for _, name := range []string{"backup", "restore"} {
		t.Run(name, func(t *testing.T) {
			command := newTestTarballCommand()
			sub := findCommand(command.Commands(), name)
			if !assert.NotNil(t, sub, "missing %s command", name) {
				return
			}

			var output bytes.Buffer
			sub.SetOut(&output)
			sub.SetErr(&output)
			if !assert.NoError(t, sub.Help()) {
				return
			}

			help := output.String()
			normalizedHelp := strings.Join(strings.Fields(help), " ")
			assert.Contains(t, normalizedHelp, expectedHelpDefault)
			assert.NotContains(t, help, home)
			assert.NotContains(t, help, dataHome)
			assert.NotContains(t, help, envTarget)

			target, err := sub.Flags().GetString(tarballTargetFlag)
			if assert.NoError(t, err) {
				assert.Equal(t, envTarget, target)
			}
		})
	}
}
```

Use the existing test imports and helpers; do not add dependencies. Follow `@test-driven-development`.

- [x] **Step 2: Run the new test and confirm it fails for the current behavior**

Run: `go test ./cmd/dum -run '^TestTarballCommand_TargetHelpUsesSymbolicDefault$' -count=1`

Expected: FAIL because help contains the resolved `envTarget` and does not show both symbolic
fallback paths.

### Task 2: Separate help default from runtime default

**Files:**

- Modify: `cmd/dum/tarball.go`

- [x] **Step 1: Add symbolic display metadata and shared flag-registration helper**

Add a constant with the exact help text:

```go
tarballTargetHelpDefault = "TGZ_TARGET, else $XDG_DATA_HOME/tarballs, else ~/.local/share/tarballs"
```

Add a helper that registers the runtime value and then changes only the displayed default:

```go
func addTarballTargetFlag(cmd *ca.Command, usage string) {
	cmd.Flags().String(tarballTargetFlag, defaultTarballTarget(), usage)
	cmd.Flags().Lookup(tarballTargetFlag).DefValue = tarballTargetHelpDefault
}
```

Replace the two direct `cmd.Flags().String(tarballTargetFlag, defaultTarballTarget(), ...)` calls in
`newBackupCommand` and `newRestoreCommand` with `addTarballTargetFlag`, preserving each existing
usage description. Do not change `defaultTarballTarget()` or runtime precedence.

- [x] **Step 2: Run help and runtime regression tests**

Run: `go test ./cmd/dum -run '^(TestTarballCommand_TargetHelpUsesSymbolicDefault|TestTarballCommand_TargetDefaultPrecedence)$' -count=1`

Expected: PASS. Help contains the symbolic chain for both commands; target values still follow
`TGZ_TARGET`, `XDG_DATA_HOME`, and home fallback as existing tests assert.

### Task 3: Verify CLI package and commit

**Files:**

- Verify: `cmd/dum/tarball.go`
- Verify: `cmd/dum/tarball_test.go`

- [x] **Step 1: Run the CLI package tests**

Run: `go test ./cmd/dum/... -count=1`

Expected: PASS.

- [x] **Step 2: Run static checks**

Run: `go vet ./cmd/dum/... && golangci-lint run ./cmd/dum/...`

Expected: both commands exit successfully with no lint issues.

- [x] **Step 3: Review and commit the focused diff**

Run: `git diff --check && git diff -- cmd/dum/tarball.go cmd/dum/tarball_test.go`

Confirm only the target help metadata and regression test changed, then commit:

```bash
git add cmd/dum/tarball.go cmd/dum/tarball_test.go
git commit -m "fix(tarball): keep target help default symbolic"
```

Use `@verification-before-completion` before claiming the change is complete.
