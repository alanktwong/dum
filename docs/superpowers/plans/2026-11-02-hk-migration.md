# Replace pre-commit with hk Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use subagent-driven-development
> (recommended) or executing-plans to implement this plan task-by-task. Steps use
> checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace pre-commit with mise-managed hk for this repository's local Git
hooks while preserving Git LFS integration and the approved policy changes.

**Architecture:** Pin hk v2.11.0 in mise and define file checks in `hk.pkl` using
hk builtins plus the existing TruffleHog Git scan. Install hooks repo-locally
through `hk install --mise`, require Git 2.54+, keep Git LFS's own pre-push hook,
and leave CI unchanged. The mise-managed hk setup supports Linux x64/arm64 and
macOS arm64; Intel macOS developer setup is unsupported, while GoReleaser's
macOS x64 product build remains unchanged.

**Tech Stack:** mise, hk v2.11.0, Pkl configuration, Bash, Git 2.54+, Git LFS,
ShellCheck, TruffleHog.

**Spec:** `docs/superpowers/specs/2026-11-02-hk-migration-design.md`

**PR Strategy:** One cohesive PR targeting `main`; no stack needed.

---

## File Map

| Path | Change | Responsibility |
| --- | --- | --- |
| `hk.pkl` | Create | Shared file checks plus explicit commit-msg and pre-push hooks. |
| `mise.toml` | Modify | Pin hk; remove pre-commit and Python. Keep Node/npm. |
| `mise.lock` | Regenerate | Lock hk for Linux x64/arm64 and macOS arm64. |
| `tools/lib/setup.sh` | Modify | Enforce Git 2.54+, validate hk config, install hk and Git LFS. |
| `tools/hooks/check-added-large-files.sh` | Create | Skip Git LFS-attributed files, then check remaining paths with hk. |
| `README.md` | Modify | Explain hk setup, platform support, local commands, Git minimum, and upgrade cleanup. |
| `AGENTS.md` | Modify | Replace active pre-commit instructions and document supported dev platforms. |
| `docs/installer.md` | Modify | Remove pre-commit from the developer-tool reference list. |
| `.pre-commit-config.yaml` | Delete | Remove the retired hook configuration after the migration audit. |
| `.commitlintrc.yml` | Delete | Remove commitlint-only configuration. |

Do not edit historical specs or plans for this migration. No Go source, CI
workflow, or `tools/build.sh` changes are expected.

Start each task with a clean index. Before each commit, inspect `git status
--short`, `git diff --cached --stat`, and `git diff --cached --check`; stage only
the paths listed for that task. If a touched path has pre-existing user changes,
stop and separate them before continuing.

## Task Structure

### Task 1: Add hk configuration while retaining the old runner temporarily

**Files:**

- Create: `hk.pkl`
- Create: `tools/hooks/check-added-large-files.sh`
- Modify: `mise.toml`
- Modify: `mise.lock`
- Read only: `.pre-commit-config.yaml`

- [x] **Step 1: Add hk v2.11.0 to the project toolchain**

From the repository root, add the pinned hk release with mise:

```bash
mise use hk@2.11.0
mise lock --platform linux-x64,linux-arm64,macos-arm64
```

Keep pre-commit and Python in the tool list for this task so the existing
configuration remains runnable during the audit. Review mise's selected hk
backend; keep it consistent with the project's pinned binary-tool policy.

- [x] **Step 2: Install the locked tools and verify hk version**

Run:

```bash
mise install --locked
hk --version
```

Expected: mise installs the lockfile successfully on Linux x64/arm64 and macOS
arm64, and hk reports version `2.11.0`. Intel macOS is intentionally
unsupported by this locked setup.

- [x] **Step 3: Audit the legacy config without adopting compatibility fallbacks**

Write migration output outside the repository:

```bash
audit_dir="$(mktemp -d)"
hk migrate pre-commit \
  --config .pre-commit-config.yaml \
  --runner pre-commit \
  --output "$audit_dir/hk-audit.pkl"
```

Inspect the generated file for hook coverage and unsupported arguments. It may
contain steps that invoke pre-commit; use it only as an audit artifact and do
not copy those compatibility steps into `hk.pkl`.

Before removing the old config, record its LFS exemption in a disposable clone:
track and stage a file larger than 5,000 KB with Git LFS, then run
`pre-commit run check-added-large-files --files large-lfs.bin` and confirm it
passes. An ordinary newly added file above the threshold must fail.
Also verify that an oversized staged modification fails, matching the old hook's
behavior for non-LFS files.

- [x] **Step 4: Add the native hk configuration**

Create `hk.pkl` with this version-matched schema, builtins, and hook map:

```pkl
amends "package://github.com/jdx/hk/releases/download/v2.11.0/hk@2.11.0#/Config.pkl"
import "package://github.com/jdx/hk/releases/download/v2.11.0/hk@2.11.0#/Builtins.pkl"

fail_fast = false

steps {
  ["trailing-whitespace"] = Builtins.trailing_whitespace
  ["check-added-large-files"] = (Builtins.check_added_large_files) {
    check = new Command {
      argv = List("bash", "tools/hooks/check-added-large-files.sh", "5000", "{{files}}")
    }
  }
  ["check-case-conflict"] = Builtins.check_case_conflict
  ["check-merge-conflict"] = Builtins.check_merge_conflict
  ["mixed-line-ending"] = Builtins.mixed_line_ending
  ["end-of-file-fixer"] = Builtins.newlines
  ["check-executables-have-shebangs"] = Builtins.check_executables_have_shebangs
  ["detect-private-key"] = Builtins.detect_private_key
  ["shellcheck"] = Builtins.shellcheck
  ["fix-smartquotes"] = Builtins.fix_smart_quotes
}

hooks {
  ["pre-commit"] {
    steps {
      ["check-lfs-hook"] {
        check = "./tools/build.sh check-lfs-hook"
      }
    }
  }

  ["commit-msg"] {
    steps {
      ["conventional-commit"] = Builtins.check_conventional_commit
    }
  }

  ["pre-push"] {
    steps {
      ["trufflehog"] {
        check = #"bash -c 'default_branch=$(git rev-parse --abbrev-ref origin/HEAD 2>/dev/null || printf "%s" main); base_commit=$(git rev-parse --verify "$default_branch" 2>/dev/null || git rev-parse --verify HEAD^ 2>/dev/null || git rev-parse HEAD); trufflehog git file://. --since-commit "$base_commit" --branch HEAD --only-verified --no-update --fail'"#
      }
    }
  }
}
```

Create `tools/hooks/check-added-large-files.sh` with:

```bash
#!/usr/bin/env bash
set -euo pipefail

maxkb="${1:?maximum size in KB is required}"
shift
files=()
for file in "$@"; do
    if git check-attr filter -- "$file" | grep -q ': filter: lfs$'; then
        continue
    fi
    files+=("$file")
done

if ((${#files[@]} > 0)); then
    hk util check-added-large-files --maxkb "$maxkb" "${files[@]}"
fi
```

Top-level steps create hk's shared `check`, `fix`, and `pre-commit` behavior.
Keep `check-lfs-hook` only on pre-commit; keep TruffleHog only on pre-push. Do
not replace the Git-history scan with hk's filesystem TruffleHog builtin.

- [x] **Step 5: Validate config and inspect selected steps**

Run:

```bash
hk validate
hk check --all --plan
hk run pre-commit --plan
```

Create a temporary commit-message file and inspect the commit-msg and pre-push
plans:

```bash
plan_dir="$(mktemp -d)"
printf 'feat: validate hk config\n' > "$plan_dir/commit-message.txt"
hk run commit-msg --plan "$plan_dir/commit-message.txt"
hk run pre-push --plan origin https://example.invalid/repo.git
```

Expected: hk validates the Pkl config, the plans include the intended hooks,
and no final hk step invokes pre-commit or prek.

- [x] **Step 6: Commit the config and temporary hk pin**

```bash
git add hk.pkl mise.toml mise.lock
git diff --cached --stat
git diff --cached --check
git commit -m "build(hooks): add hk configuration" \
  -m "Add native hook steps before retiring the legacy runner."
```

### Task 2: Switch setup to hk and remove pre-commit dependencies

**Files:**

- Modify: `tools/lib/setup.sh`
- Modify: `mise.toml`
- Modify: `mise.lock`
- Delete: `.pre-commit-config.yaml`
- Delete: `.commitlintrc.yml`

- [x] **Step 1: Remove existing pre-commit shims before changing the tool list**

On an existing clone, while the old mise config and pre-commit executable are
still available, run:

```bash
pre-commit uninstall
pre-commit uninstall --hook-type commit-msg
```

Do not uninstall Git LFS. Verify its pre-push hook remains:

```bash
hooks_dir="$(git rev-parse --git-path hooks)"
grep -q '^git lfs pre-push' "$hooks_dir/pre-push"
```

Expected: both pre-commit shims are removed and the Git LFS check succeeds.

- [x] **Step 2: Write the failing Git-version install test**

Use a temporary Bash harness that sources `tools/lib/setup.sh`, stubs `git` and
`hk`, and puts a logging `pre-commit` stub first on `PATH`. The `git` function
must also intercept `git lfs install`; never let this red test invoke real hook
installers. The harness must assert that Git `2.53.0` makes
`setup_install_tools` fail before `hk validate`, the global-hook query,
`git lfs install`, or `hk install`. For Git `2.54.0` with no global hook, allow
the read-only preflights followed by these state-changing calls in this order:

```text
hk validate
git lfs install --force
hk install --mise
```

Run the harness against the current setup and confirm it fails because no Git
minimum is enforced yet. Keep the harness in the approved temporary directory;
do not add a new test framework for this shell-only guard. The `git` stub must
return status 1 for the global-hook query when no global hk install exists and
status 0 when simulating `hook.hk-*` entries. Also simulate status 128 for an
unreadable global config and assert fail-closed behavior.

- [x] **Step 3: Add a portable Bash Git-version guard**

In `tools/lib/setup.sh`, add a helper that reads `git --version`, parses the
major/minor components with Bash regular-expression matching, and rejects
versions below 2.54 with a clear error naming the required and detected
versions. Do not rely on GNU-only `sort -V`.

Use this implementation shape so the comparison works on macOS Bash without
GNU `sort` and rejects malformed tokens such as `2.54oops`:

```bash
setup_check_git_version() {
    local version major minor
    version="$(git --version | awk '{print $3}')"
    if [[ ! "$version" =~ ^([0-9]+)\.([0-9]+)(\.[0-9]+([.-][[:alnum:]]+)*)?$ ]]; then
        echo >&2 "Could not parse Git version '$version'; Git 2.54 or newer is required."
        return 1
    fi

    major=$((10#${BASH_REMATCH[1]}))
    minor=$((10#${BASH_REMATCH[2]}))
    if (( major < 2 || (major == 2 && minor < 54) )); then
        echo >&2 "Git 2.54 or newer is required (found $version)."
        return 1
    fi
}
```

Add a second read-only helper for global hook scope:

```bash
setup_check_no_global_hk_hooks() {
    local status
    if git config --global --get-regexp '^hook\.hk-' >/dev/null 2>&1; then
        echo >&2 "Global hk hooks are configured; this repository uses repo-local hk. Choose one installation scope. Run 'hk uninstall --global' only if you intend to switch to repo-local hooks."
        return 1
    else
        status=$?
        if (( status != 1 )); then
            echo >&2 "Could not inspect global Git hook configuration; refusing repo-local installation."
            return 1
        fi
    fi
    return 0
}
```

Call `setup_check_git_version || return` from both `setup_check_tools` and
`setup_install_tools`. In `setup_install_tools`, run `hk validate || return` and
`setup_check_no_global_hk_hooks || return` before any hook changes. These
explicit returns must stop setup even when a caller invokes the function from a
conditional and Bash `errexit` is suppressed. If a global hk hook is configured,
fail with its scope-conflict message; do not call `hk install --force-local` or
modify global Git config. If none is configured, run `git lfs install --force`
followed by `hk install --mise`.

- [x] **Step 4: Run the version-guard test green**

Re-run the temporary harness from Step 2.

Expected: Git `2.53.0` is rejected with no hook-changing calls; Git `2.54.0`
with no global hk entry passes and records the setup commands in the expected
order. A simulated global `hook.hk-*` entry, `2.54oops`, and malformed
Git-version output, and a global Git config read error must each fail before Git
LFS or hk installation, with readable errors. Only exit status 1 from the
global hook query means "no matching global hk hooks."

- [x] **Step 5: Update the required-tool list and retire old dependencies**

In `mise.toml`, remove the pre-commit and Python entries; keep the hk pin added
in Task 1. In `tools/lib/setup.sh`, require `git` and `hk`, retain Git LFS,
Node, and npm, and remove `pre-commit` and `python3`. Keep the existing missing
tool remediation message pointing to `mise install --locked`.

Regenerate lock entries for supported hk platforms and install them:

```bash
mise lock --platform linux-x64,linux-arm64,macos-arm64
mise install --locked
mise run check_tools
```

Expected: hk and Git LFS are found, no pre-commit or Python runner is required,
Node/npm remain available, and hk is locked for Linux x64/arm64 and macOS arm64.
Intel macOS cannot run `mise install --locked` for this toolchain.

- [x] **Step 6: Delete legacy hook and commitlint configuration**

Delete `.pre-commit-config.yaml` and `.commitlintrc.yml`. Do not delete
historical references under `docs/superpowers/specs/` or
`docs/superpowers/plans/`.

- [x] **Step 7: Install hooks and verify Git LFS coexistence**

On Git 2.54+, run:

```bash
mise run install
hk validate
git config --local --get-regexp '^hook\.'
```

Expected: with no global hk hooks configured, hk has repo-local config-based
hooks, and
`$(git rev-parse --git-path hooks)/pre-push` still contains the Git LFS hook.
If global `hook.hk-*` entries are detected, setup must stop before changing
hooks and explain the scope conflict. Do not pass `--force` or `--force-local`.

- [x] **Step 8: Commit the runner transition**

```bash
git add -A -- tools/lib/setup.sh mise.toml mise.lock \
  .pre-commit-config.yaml .commitlintrc.yml
git commit -m "build(hooks): replace pre-commit with hk" \
  -m "Use mise-managed hk and preserve Git LFS hook ownership."
```

### Task 3: Update active setup documentation

**Files:**

- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `docs/installer.md`

- [x] **Step 1: Document the existing-clone upgrade order**

In `README.md`, explain that existing clones must run both pre-commit
uninstall commands while the old tool is available, before switching to the new
mise tool list. Then document `mise trust && mise install --locked` and
`mise run install` for hk and Git LFS setup. State the Git 2.54 minimum and
repo-local `hk install --mise` behavior. Retain the bash/zsh/fish mise
activation instructions and explain that local hk launchers need `mise` on
Git's runtime `PATH`. State that the locked developer toolchain supports Linux
x64/arm64 and macOS arm64; on Intel macOS, `mise install --locked` is expected
to fail because hk has no x64 prebuilt. Clarify that GoReleaser's macOS x64
product builds remain available. Explain that setup stops with guidance if
global hk hooks are configured, so contributors must choose one installation
scope.

- [x] **Step 2: Document local hk commands and intentional policy changes**

In `README.md`, replace active pre-commit/commitlint setup text with `hk check`,
`hk fix`, and `hk check --all`. Say Conventional Commit subjects are validated
locally, commit bodies are optional, CI remains unchanged, and TruffleHog now
runs on push.

- [x] **Step 3: Update the agent and installer guides**

In `AGENTS.md`, update setup commands, commit-message guidance, and Important
Files to refer to `hk.pkl` and hk rather than `.pre-commit-config.yaml`. State
the same Linux x64/arm64 and macOS arm64 dev-tool support, with Intel macOS
`mise install --locked` expected to fail.
In `docs/installer.md`, remove the obsolete pre-commit tool reference and add
the hk reference in the same list.

- [x] **Step 4: Check active documentation references**

Run:

```bash
rg -n 'pre-commit|commitlint|gitlint' README.md AGENTS.md docs/installer.md
```

Expected: `pre-commit` references are limited to hk's Git hook event and the
one-time uninstall instruction for old hooks. No current instructions imply the
pre-commit framework remains installed or supported, and no active references
to commitlint/Gitlint remain. Historical specs/plans remain unchanged.

- [ ] **Step 5: Commit documentation**

```bash
git add README.md AGENTS.md docs/installer.md
git commit -m "docs(hooks): document hk setup" \
  -m "Explain the Git minimum, upgrade cleanup, and local hk workflow."
```

### Task 4: Clean repository trailing whitespace only

**Files:**

- Modify: tracked files reported by hk's `trailing-whitespace` step

- [ ] **Step 1: Preview affected files**

After the documentation commit, require a clean worktree and index before the
repository-wide operation:

```bash
test -z "$(git status --porcelain --untracked-files=all)"
git diff --cached --quiet
```

If either command fails, stop and separate existing work first.

Run:

```bash
hk check --all --step trailing-whitespace --plan
```

Confirm the selection is the repository-wide tracked/eligible file set and that
only the trailing-whitespace step is selected.

- [ ] **Step 2: Apply only the approved whitespace cleanup**

Run:

```bash
hk fix --all --step trailing-whitespace --no-stage
```

Expected: trailing whitespace is removed across the repository, including
Markdown hard-break spaces; no other formatter/fixer runs and no files are
staged automatically.

- [ ] **Step 3: Review and verify the cleanup diff**

Run:

```bash
git diff --stat
git diff --check
hk check --all --step trailing-whitespace
```

Review Markdown diffs for the approved hard-break rendering change. In
`.github/PULL_REQUEST_TEMPLATE/brief.md`, replace the stripped blank bullet with
`- Describe change` to preserve a valid list item. Expected: the
trailing-whitespace check passes and `git diff --check` reports no errors.

- [ ] **Step 4: Commit the whitespace-only diff separately**

Before staging, confirm `git status --short` contains only the hk-generated
cleanup and confirm `git diff --cached --quiet`. Stage only those reviewed paths,
then inspect the staged diff before committing:

```bash
paths_file="$(mktemp)"
git diff --name-only -z > "$paths_file"
if [[ -s "$paths_file" ]]; then
  xargs -0 git add -- < "$paths_file"
fi
rm -f "$paths_file"
git diff --cached --stat
git diff --cached --check
if ! git diff --cached --quiet; then
  git commit -m "style: strip repository trailing whitespace" \
    -m "Apply the approved whitespace rule, including Markdown files."
else
  echo "No whitespace cleanup changes; skip commit."
fi
```

### Task 5: Verify hook behavior end to end

**Files:**

- Verify: `hk.pkl`, `mise.toml`, `mise.lock`, `tools/lib/setup.sh`
- Verify: Git hooks in the current clone; no tracked file changes expected

- [ ] **Step 1: Verify hk configuration and execution plans**

Run:

```bash
hk validate
hk check --all --plan
hk run pre-commit --plan
```

Also inspect commit-msg and pre-push plans using a temporary commit-message file
and a dummy remote URL. Expected: file checks appear under check/pre-commit,
`check-lfs-hook` is pre-commit-only, Conventional Commit validation is
commit-msg-only, and TruffleHog is pre-push-only.

- [ ] **Step 2: Test commit-message policy**

Run `hk run commit-msg <file>` for a valid Conventional Commit subject with no
body and for an invalid subject.

Expected: the no-body message passes and the invalid subject fails. Do not
reintroduce Gitlint's body requirement.

- [ ] **Step 3: Test staged-file selection and the large-file limit**

Create an independent local clone so the checks have the project's `hk.pkl`,
`tools/build.sh`, and `tools/lib/` files without touching the target clone's
index or hooks:

```bash
test_root="$(mktemp -d)"
hk_bin="$(mise which hk)"
git_lfs_bin="$(mise which git-lfs)"
test_path="$(dirname "$hk_bin"):$(dirname "$git_lfs_bin"):/usr/bin:/bin:/usr/sbin:/sbin"
git clone --local --no-hardlinks . "$test_root/repo"
```

Run the cases in `$test_root/repo`, invoking `"$hk_bin"` directly and setting
`PATH="$test_path"` so the clone needs no separate mise trust entry while hk
builtins and Git LFS remain available. Set
`GIT_CONFIG_GLOBAL="$test_root/gitconfig"` for `git lfs install --force` so
Git LFS setup cannot modify the user's global Git config. Stage a file with a
trailing-space fix, add a separate unstaged edit to the same file, run
`"$hk_bin" run pre-commit`, and confirm the fix is staged while the unstaged
edit is restored. Then test a newly added file at 5,000 KB, one byte over the
limit, an oversized staged modification, and an LFS-attributed file over the
limit. Expected: the boundary passes, non-LFS additions and modifications over
the threshold fail, and the LFS-attributed file passes. The Task 1 baseline test
must also confirm the legacy pre-commit hook skipped the LFS pointer.

- [ ] **Step 4: Verify TruffleHog runs on pre-push**

Create a separate disposable clone, resolve `hk_bin` with `mise which hk`, and
put a temporary `trufflehog` executable earlier on `PATH` that records its
arguments. Invoke `"$hk_bin" run pre-push origin
https://example.invalid/repo.git` with a local parent commit available and
`PATH` containing the stub, hk's bin directory, and system Git. Confirm the
existing fallback selects the base ref and the command receives
`--only-verified --no-update --fail`. Do not scan the target repository or trust
the clone's mise config for this wiring test.

- [ ] **Step 5: Verify Git LFS and hk are both installed**

On Git 2.54+, confirm the local hk hook command exists in Git config and the
Git LFS pre-push script remains under `git rev-parse --git-path hooks`. Run the
pre-commit `check-lfs-hook` step and confirm it passes. Never use
`hk install --force` as a workaround.

- [ ] **Step 6: Run project and migration checks**

Run:

```bash
mise run check_tools
mise run check
hk validate
hk check --all
git diff --check
git status --short --untracked-files=all
```

Expected: project checks pass, hk validates and checks all file steps, and no
unintended working-tree changes remain. Do not modify CI workflows.

## Completion Criteria

- hk v2.11.0 is mise-pinned; `hk.pkl` imports matching schema/builtins.
- No active pre-commit config, runner dependency, Python requirement, or
  commitlint config remains; historical docs are untouched.
- `mise run install` rejects Git below 2.54 before mutating hooks, then installs
  Git LFS and repo-local hk without `--force`.
- A pre-existing global hk installation is detected before mutation; the setup
  does not force local hooks or alter global Git configuration.
- Partial commits remain safe; the 5,000 KB limit is verified; the approved
  commit-message simplification and Markdown whitespace change are tested.
- TruffleHog runs at pre-push and Git LFS's hook remains installed.
- CI is unchanged, and the final whitespace-only cleanup is reviewed separately.
