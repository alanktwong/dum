# Design: Replace pre-commit with hk

Date: 2026-11-02
Status: Approved

## Goal

Replace pre-commit as this repository's local Git hook manager with
[hk](https://hk.jdx.dev), managed through mise. hk will own the repository's
pre-commit, commit-msg, and pre-push checks. Git LFS keeps its own pre-push hook.

The design uses repo-local hk installation, requires Git 2.54 or newer, activates
the configured TruffleHog pre-push scan, and leaves CI unchanged.

The mise-managed hk setup supports Linux x64/arm64 and macOS arm64. Intel macOS
is excluded from the local mise toolchain because hk v2.11.0 has no macOS x64
prebuilt; GoReleaser's macOS x64 product build remains unchanged.

## Non-goals and intentional policy changes

- Do not add hook checks to GitHub Actions in this migration.
- Do not install hk globally; keep setup in `mise run install`.
- Do not add a Rust source build just for Intel Macs; its extra toolchain and
  compile time are not worth retaining Intel macOS developer setup.
- Do not remove Node/npm in this migration. `.opencode/package.json` still
  declares a plugin dependency; verify its runtime requirements separately.
- Simplify commit-message validation: replace commitlint and Gitlint with hk's
  Conventional Commit check. The Gitlint non-empty-body rule is intentionally
  dropped.
- Strip all trailing whitespace, including Markdown hard-break spaces. Run a
  one-time whole-repository cleanup and review the resulting diff.
- Do not broaden that cleanup to unrelated final-newline or smart-quote fixes.
  `hk check --all` currently reports such pre-existing findings in
  `.env-sample`, `.github/CODEOWNERS`, `cfg/goreleaser.local.yaml`,
  `docs/superpowers/plans/plan.exclude_generated_coverage.md`, and
  `docs/superpowers/plans/2026-08-13-viper-config.md`; leave them as existing
  debt.

## Current state

`.pre-commit-config.yaml` configures these checks:

- Trailing whitespace, preserving Markdown hard breaks, plus a 5,000 KB limit
  for added files.
- Case conflicts, merge-conflict markers, mixed line endings, final newlines,
  executable shebangs, and private-key detection.
- ShellCheck and smart-quote fixes.
- A local `check-lfs-hook` command at pre-commit.
- Gitlint and commitlint at commit-msg.
- TruffleHog at pre-push.

`mise run install` currently installs only pre-commit and commit-msg hook
shims, then runs `git lfs install --force`. It does not install pre-commit's
pre-push hook. Therefore the configured TruffleHog scan is not activated by
normal setup, while Git LFS's pre-push hook is.

## Hook configuration

Add a version-pinned `hk.pkl` that imports the hk schema and builtins from the
same release. Use top-level `steps` for file checks so hk provides its shared
`check`, `fix`, and `pre-commit` workflows. The pre-commit workflow must fix
staged files, stage fixes, and stash/restore unstaged work so partial commits
remain safe. Set `fail_fast = false` so hk reports all failing steps, matching
pre-commit's aggregate failure reporting.

| Existing check | hk configuration | Notes |
| --- | --- | --- |
| `trailing-whitespace` | `Builtins.trailing_whitespace` | Remove the Markdown exception; all trailing spaces are stripped. |
| `check-added-large-files` | Custom hk step that skips `filter=lfs` paths, then calls `hk util check-added-large-files --maxkb 5000` | Preserve the threshold and pre-commit's Git LFS exemption. |
| `check-case-conflict` | `Builtins.check_case_conflict` | |
| `check-merge-conflict` | `Builtins.check_merge_conflict` | |
| `mixed-line-ending` | `Builtins.mixed_line_ending` | Retain automatic normalization behavior. |
| `end-of-file-fixer` | `Builtins.newlines` | |
| `check-executables-have-shebangs` | `Builtins.check_executables_have_shebangs` | |
| `detect-private-key` | `Builtins.detect_private_key` | |
| `shellcheck` | `Builtins.shellcheck` | Use hk's shell-file selection. |
| `fix-smartquotes` | `Builtins.fix_smart_quotes` | |

Add `check-lfs-hook` explicitly to the `pre-commit` hook, not to the shared
top-level steps. This retains its current purpose without making `hk check` or
`hk fix` depend on the presence of the Git LFS hook.

For `commit-msg`, use `Builtins.check_conventional_commit`. Remove the
commitlint and Gitlint hooks and the now-unused `.commitlintrc.yml`. This check
does not retain Gitlint's body requirement.

For `pre-push`, add a custom hk step that runs the existing TruffleHog Git scan
with its current base-ref fallback and flags (`--only-verified`, `--no-update`,
and `--fail`). This deliberately activates a check that setup did not previously
install. Do not substitute hk's filesystem TruffleHog builtin; it would change
the current Git-history scan semantics.

Git LFS remains independently installed with `git lfs install --force`. Keep
`check-lfs-hook` as a pre-commit sanity check. With Git 2.54+, hk's repo-local
config-based hooks coexist with Git LFS's `.git/hooks/pre-push` hook.

Implement the large-file check in `tools/hooks/check-added-large-files.sh` so
the per-file Git attribute check can exclude Git LFS paths before delegating the
remaining files to hk's utility. Do not place this script under `tools/lib/`,
which `tools/build.sh` sources as a library.

## Installation and toolchain

- Add hk as a version-pinned mise tool and update `mise.lock` for Linux x64/arm64
  and macOS arm64. The hk release has no macOS x64 artifact, so
  `mise install --locked` is not supported on Intel Macs.
- Delete `.pre-commit-config.yaml` after using it for the migration audit; no
  compatibility runner or active pre-commit configuration remains.
- Remove the pre-commit tool and Python runtime pin used only by the pre-commit
  runner. Remove those commands from `setup_check_tools`; retain Node/npm for
  now as noted under non-goals.
- Change `setup_install_tools` to check Git 2.54+ and run `hk validate` before
  changing hook state. Detect global `hook.hk-*` entries and fail with guidance
  to choose one installation scope; do not mutate global Git config or use
  `--force-local`. Then install Git LFS and run repo-local `hk install --mise`.
- Fail clearly on older Git and do not use `hk install --force`, which could
  replace Git LFS's pre-push hook on legacy Git.
- `hk install --mise` requires `mise` on Git's runtime `PATH`; retain the
  existing README guidance to activate mise in the shell.
- For existing clones, require pre-upgrade cleanup while the old mise config and
  pre-commit executable are still available: run `pre-commit uninstall` and
  `pre-commit uninstall --hook-type commit-msg` before switching to the new
  toolchain. Call this out in migration/release instructions. New clones need no
  cleanup. Preserve the Git LFS pre-push hook.

## Local and CI workflows

Developers install the mise-managed toolchain, then run `mise run install` to
set up Git LFS and hk. `hk check`, `hk fix`, and `hk check --all` remain available
for manual local use. Git commit-msg and pre-push hooks run through hk
automatically. CI workflows remain unchanged and do not install hk.

Update active setup references in `README.md`, `AGENTS.md`, and
`docs/installer.md`. Historical specs and plans remain point-in-time records and
are not swept for old pre-commit references. `README.md` and `AGENTS.md` must
state that `mise install --locked` is expected to fail on Intel macOS because
hk has no macOS x64 artifact; GoReleaser's macOS x64 product builds remain
available.

## Migration and validation

Use `hk migrate pre-commit` only as an audit aid, writing its output outside the
final configuration. Its compatibility fallbacks may continue to invoke
pre-commit/prek, so no such fallback may remain in the adopted `hk.pkl`.

Before installing hooks:

1. Validate `hk.pkl` with `hk validate`.
2. Inspect `hk check --all --plan` and the plans for pre-commit, commit-msg, and
   pre-push to confirm step selection.
3. Check that the Git-version guard rejects versions below 2.54 before hook
   configuration changes.
4. Confirm the lock contains hk artifacts for Linux x64/arm64 and macOS arm64,
   and documents macOS x64 as unsupported for the mise-managed toolchain.

After setup:

1. Remove old pre-commit and commit-msg shims on existing clones; confirm Git
   LFS's pre-push hook remains installed.
2. Exercise staged/unstaged partial commits and confirm fixes are staged while
   unstaged edits are restored.
3. Test the 5,000 KB large-file boundary, verify Git LFS-attributed files are
   skipped as before, and test valid/invalid Conventional Commit subjects. A
   subject without a body should now pass, documenting the intentional policy
   change.
4. Verify TruffleHog runs on pre-push and that Git LFS and hk both run without
   replacing one another.
5. Run `hk fix --all --step trailing-whitespace --no-stage` for the one-time
   whitespace cleanup, review the full diff (including Markdown rendering
   changes), then run `hk check --all --step trailing-whitespace`. A full
   `hk check --all` remains nonzero only for the existing findings listed above;
   do not fix those in this migration.

Hook failures block the associated commit or push. Invalid hk configuration and
unsupported Git versions or conflicting global hk installation must fail
before hook installation. No CI changes are part of this migration.
