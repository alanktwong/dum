# Git Hooks

This repository uses [hk](https://github.com/jdx/hk) to run local checks before commits and pushes.
The goal is to catch formatting problems, accidental secrets, and malformed commit messages before
they reach CI, without slowing down day-to-day work.

The hook configuration lives in [`hk.pkl`](../../hk.pkl). Everything else on this page is context for
working with it.

## One-time setup

You need:

- Git 2.54 or newer
- The mise-managed developer tools installed (`mise trust && mise install --locked`)
- The repository cloned and the current directory set to its root

Then run:

```shell
mise run install
```

This task does four things:

1. Verifies the Git version is 2.54 or newer.
2. Validates `hk.pkl`.
3. Refuses to continue if hk hooks are configured globally (this repository expects repo-local
   hooks only).
4. Runs `git lfs install --force` and `hk install --mise`.

`hk install --mise` wires the hooks into your **local Git config** rather than writing scripts into
`.git/hooks/`. Each hook command is wrapped in `mise x --` so the pinned tools are available even if
mise is not activated in your shell. Git's own LFS hooks in `.git/hooks/` stay in place and run
alongside the hk hooks.

You only need to run this once per clone. Re-run `mise run install` after pulling changes to
`hk.pkl` or `tools/`.

> **Note:** Intel macOS is not supported for the mise-managed developer setup because hk v2.11.0
> has no macOS x64 binary. GoReleaser's macOS x64 product builds are unaffected.

## What runs when

| Event        | What runs                                                                 | Behavior                          |
| ------------ | ------------------------------------------------------------------------- | --------------------------------- |
| `pre-commit` | Whitespace, file, and shell-script hygiene checks, plus an LFS guard       | Fixes files, stages them, and restores your unstaged edits |
| `commit-msg` | Conventional Commits validation on the subject                            | Rejects the commit on failure     |
| `pre-push`   | TruffleHog secret scan                                                    | Rejects the push on failure       |

`fail_fast` is disabled, so all applicable steps run and report together.

### pre-commit

The shared steps in `hk.pkl` run in **fix** mode before each commit. They operate on staged files:

- `trailing-whitespace` removes trailing whitespace.
- `check-added-large-files` blocks files over 5000 KB. It is backed by
  [`tools/hooks/check-added-large-files.sh`](../../tools/hooks/check-added-large-files.sh), which skips
  files tracked by Git LFS.
- `check-case-conflict` catches filenames that collide on case-insensitive filesystems.
- `check-merge-conflict` catches leftover merge markers.
- `mixed-line-ending` normalizes line endings.
- `end-of-file-fixer` ensures files end with a newline.
- `check-executables-have-shebangs` requires executable scripts to have a shebang.
- `detect-private-key` blocks committed private keys.
- `shellcheck` lints shell scripts.
- `fix-smartquotes` replaces smart quotes with plain quotes.
- `check-lfs-hook` runs `./tools/build.sh check-lfs-hook`, which verifies that Git LFS's own
  `pre-push` hook is installed.

If a step modifies a file, hk stages the fix. Any unstaged edits you had in the same files are
restored afterward, so `git diff` keeps showing your work as you left it.

### commit-msg

Commit subjects must follow [Conventional Commits](https://www.conventionalcommits.org). Bodies are
optional.

Allowed types: `build`, `chore`, `ci`, `docs`, `feat`, `fix`, `perf`, `refactor`, `revert`, `style`,
`test`.

```text
feat: add brew task retries
fix(playbook): skip disabled tasks during dry run
docs: explain git hooks
```

Temporary rebase commits (`fixup!`, `squash!`, `amend!`) and Git-generated merge/revert subjects are
skipped automatically.

### pre-push

TruffleHog scans the commits you are about to push, starting from the origin default branch. If it
finds a verified secret, fix the commit rather than bypassing the hook. A secret that already
reached history should be rotated and removed with `git filter-repo` or an equivalent history
rewrite.

## Running checks by hand

You do not need to create a commit to see what the hooks would do.

```shell
hk validate                     # Check hk.pkl syntax
hk check --plan                 # Show steps and file selection, run nothing
hk run pre-commit --plan        # Show the pre-commit plan (fix mode)
hk run pre-push --plan          # Show the pre-push plan
hk run commit-msg --plan /tmp/msg
```

Run the checks directly:

```shell
hk check                        # Check modified files
hk check --all                  # Check every tracked and eligible untracked file
hk fix                          # Apply fixes and stage them
hk fix --no-stage               # Apply fixes without staging

hk check -S shellcheck --all    # Run one step
hk fix -S trailing-whitespace path/to/file
hk check --all --why shellcheck # Explain why a file was included or skipped
```

`hk check` and `hk fix` are read/write opposites: `check` only reports, `fix` can modify and stage.
For a read-only, CI-style gate, use `hk check --all`.

## Testing the hook wiring

The commands above test hk itself. To test what Git actually invokes:

```shell
git hook list pre-commit
git hook run pre-commit
```

Test the commit-msg hook with a temporary message file:

```shell
printf 'bad message\n' > /tmp/msg
git hook run commit-msg -- /tmp/msg    # should fail

printf 'feat: add thing\n' > /tmp/msg
git hook run commit-msg -- /tmp/msg    # should pass
```

For `pre-push`, call hk directly instead of `git hook run pre-push`, because Git supplies ref
information on stdin:

```shell
hk run pre-push
```

An end-to-end smoke test on a throwaway branch:

```shell
git switch -c test/hooks
# create a scratch file with trailing whitespace or no final newline
git add <file>
git commit -m "feat: hook smoke test"   # pre-commit and commit-msg fire
git restore --staged <file>
git restore <file>
git switch -
git branch -D test/hooks
```

### Verify the install

```shell
git config --get-regexp '^hook\.hk-'          # repo-local hk hooks
git config --global --get-regexp '^hook\.hk-' # should be empty
./tools/build.sh check-lfs-hook               # LFS guard passes
```

Re-run `mise run install` to repair the installation.

### Bypassing hooks

Use these sparingly:

```shell
HK=0 git commit -m "wip"         # Skip hk hooks, keep other hooks
git commit --no-verify -m "wip"  # Skip all hooks
```

`HK=0` is part of every installed hook command, so it works for `pre-commit`, `commit-msg`, and
`pre-push`.

## Troubleshooting

| Symptom                                  | Likely cause and fix                                                                 |
| ---------------------------------------- | ------------------------------------------------------------------------------------ |
| Hooks never run                          | Git older than 2.54, or `mise run install` was never run. Check `git --version`, then reinstall. |
| `Global hk hooks are configured`         | hk is installed at both scopes. Run `hk uninstall --global`, then `mise run install`. |
| `hk: command not found` during a hook    | mise is not available on the hook's `PATH`. Re-run `mise run install`, which installs hooks through `mise x --`. |
| `git-lfs is not installed properly`      | Run `git lfs install --force` or `mise run install`.                                  |
| A commit is rejected for its message     | Reword the subject to Conventional Commits, for example `docs: clarify setup steps`.  |
| `hk check` passes but `git commit` edits files | Expected: pre-commit runs in fix mode. Review the staged changes and commit again.    |
| TruffleHog finds a secret                | Remove the secret from the commit, rotate it, and rewrite history if it was already pushed. |

## Related files

| File                                                                                 | Purpose                                              |
| ------------------------------------------------------------------------------------ | ---------------------------------------------------- |
| [`hk.pkl`](../../hk.pkl)                                                                | Hook definitions, steps, and customizations           |
| [`tools/lib/setup.sh`](../../tools/lib/setup.sh)                                        | Installation checks and `hk install --mise`           |
| [`tools/build.sh`](../../tools/build.sh)                                                | Dispatcher, including `check-lfs-hook`                |
| [`tools/hooks/check-added-large-files.sh`](../../tools/hooks/check-added-large-files.sh) | Large-file check that excludes Git LFS files          |
| [`mise.toml`](../../mise.toml)                                                          | Pinned hk version and the `install` task              |
