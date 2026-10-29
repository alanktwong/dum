# Mise Dev Toolchain Adoption Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use `subagent-driven-development` (recommended) or
> `executing-plans` to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for
> tracking.

**Goal:** Make mise the repository's tool, environment, and task entry point while preserving
`tools/build.sh` as the recipe dispatcher and keeping CI and contributor workflows reproducible.

**Architecture:** Add root `mise.toml` and a four-platform `mise.lock`; tasks delegate to the
existing shell dispatcher, except `tag` and `help`. Move repository hook setup out of direnv and
into `mise run install`, switch CI to `jdx/mise-action@v5`, then remove Make/nvm/direnv entry
points only after task ordering and locked installs are verified.

**Tech Stack:** mise, Go 1.26, Bash, GitHub Actions, pre-commit, GoReleaser.

**Spec:** `docs/superpowers/specs/2026-10-18-mise-adoption-design.md`
**Branch:** Continue on `feature/mise`.
**PR strategy:** One PR. Tool configuration, task migration, CI, and documentation form one
cohesive cutover; keep implementation commits separated by task for review.
**Commit convention:** Conventional Commits. Include a commit body to satisfy repository gitlint
rules; do not tag or publish during implementation.

---

## Task Structure

### Task 1: Add mise tool, environment, and lock configuration

**Files:**

- Create: `mise.toml`
- Create: `mise.lock`
- Create if generated: `.mise/locks/`

- [x] **Step 1: Resolve backend for every required tool**

Select explicit current mise backends for `git-lfs`, `golangci-lint`, `mockery`, `goreleaser`,
`svu`, `trufflehog`, `jq`, `yq`, `shellcheck`, `pre-commit`, and `yamlfmt`. Use Mike Farah's
`yq`; do not use deprecated `ubi:`. Use the `go:` backend for `go-enum`, `goimports`, and
`go-test-coverage`. Pin `go-test-coverage` to `2.19.0`: versions `2.20.0` and `2.21.0` require
Go 1.27 and trigger an automatic toolchain download beyond the project's Go 1.26 lock. Pin mise
core Python to `3.13` because Aqua's `pre-commit` artifact is a Python zipapp. Validate
`python3` through `check_tools`.

Before recording a backend, check its executable name, supported OS/architecture assets, and
whether lockfile checksums are available. Confirm that each selected backend can resolve and
install all four supported tuples.

- [x] **Step 2: Create `mise.toml` tool and environment sections**

Declare `go = "1.26"`, `node = "22"`, and `python = "3.13"`; declare every tool from Step 1 with
an appropriate `latest` request where the spec has no existing pin, except
`go-test-coverage = "2.19.0"`. Add `[settings] lockfile = true` and `[env] _.file = ".env"`. Do
not add task definitions in this step.

- [x] **Step 3: Generate lock entries for all supported platforms**

Run:

```bash
mise lock --platform linux-x64,linux-arm64,macos-x64,macos-arm64
```

Expected: `mise.lock` contains resolved versions for every configured tool and platform, with
artifact URLs/checksums when supported. Confirm the Go tool backend uses `go@1.26.10` without
downloading a separate Go 1.27 toolchain. Review and include any generated `.mise/locks/` sidecars.

- [ ] **Step 4: Verify locked installation on every platform tuple**

Run `mise install --locked` on Linux x64, Linux arm64, macOS x64, and macOS arm64. Local native
verification has passed on macOS arm64 and Linux arm64. Defer x64 verification to native PR CI
runners; do not treat the failing Docker/QEMU Linux x64 attempt as evidence. Record each runner's
OS/architecture from CI logs. If PR CI does not cover a required tuple, leave this gate pending and
request a native runner rather than claiming four-platform verification. Do not add a permanent
CI matrix without confirming it is needed.

Execution note: macOS arm64 and Linux arm64 locked installs passed. Docker/QEMU Linux x64 crashed
inside Go 1.26.10's runtime; Rosetta is unavailable for macOS x64. The remaining native x64 checks
are explicitly deferred to PR CI.

Expected: every configured tool installs from its lock entry on all four tuples. Fix backend or
lock entries before proceeding; do not add Homebrew fallbacks.

- [x] **Step 5: Review and commit tool configuration**

Confirm no `ubi:` entries; `mise.lock` and sidecars are tracked; optional `.env` does not prevent
config loading when absent.

```bash
git add mise.toml mise.lock
if [ -d .mise/locks ]; then git add .mise/locks; fi
git commit -m "build: add mise toolchain manifest" \
  -m "Declare project tools and environment in mise." \
  -m "Lock four platform tuples; pin go-test-coverage for Go 1.26."
```

If no sidecar directory is generated, omit it from `git add`.

### Task 2: Lock code generation and migrate repository setup

**Files:**

- Modify: `internal/types/gen/task_type.go`
- Modify: `internal/types/gen/jetbrains_type.go`
- Modify: `internal/external/gen/mocks.go` (tracked generated mock, if locked Mockery output differs)
- Modify: `tools/lib/setup.sh`
- Modify: `tools/build.sh`
- Modify: `tools/lib/common.sh`

- [x] **Step 1: Run the check-tools regression harness (red)**

Before changing `setup_check_tools`, run the all-present harness with `GOPATH` set to a path absent
from `PATH`:

```bash
GOPATH=/__mise_test_gopath__ MISSING_COMMAND=__none__ bash -c \
  'source tools/lib/setup.sh; command() { [[ "$2" != "${MISSING_COMMAND:-}" ]]; }; setup_check_tools'
```

Expected before the fix: non-zero exit with the obsolete `$GOPATH/bin` PATH error. This captures
the behavior the migration must remove without modifying the host environment.

- [x] **Step 2: Route go-enum directives through the mise-managed executable**

In both Go files, replace `go run github.com/abice/go-enum@latest` with the `go-enum` executable
provided by mise:

```go
//go:generate go-enum --file=$GOFILE --names
```

- [x] **Step 3: Update `setup_install_tools` and `setup_check_tools`**

Change `setup_install_tools` to configure the repository only:

```bash
pre-commit install
pre-commit install --hook-type commit-msg
git lfs install --force
```

Remove Homebrew installs, Go `go install` loops, `direnv allow`, and the direnv activation banner.
Change `setup_check_tools` to validate the commands supplied by `mise`, including
`go-test-coverage` and `python3`; remove Homebrew, nvm bootstrap, and `$GOPATH/bin` checks. A
missing command must name the command and instruct the user to run `mise install --locked`. Keep
`node`, `npm`, and `python3` in the required command set for commitlint and its Python-based
pre-commit runner. Remove `setup_ci_tools`.

- [x] **Step 4: Verify the check-tools regression harness (green)**

Use a Bash-only stub for `command -v` so the check can run without changing host `PATH`, Git
configuration, or hooks. From the repository root, run the "all present" case:

```bash
GOPATH=/__mise_test_gopath__ MISSING_COMMAND=__none__ bash -c \
  'source tools/lib/setup.sh; command() { [[ "$2" != "${MISSING_COMMAND:-}" ]]; }; setup_check_tools'
```

Expected: exit 0 and `All good!`, without requiring `brew`, `nvm`, or `$GOPATH/bin`.

Run the missing-tool case:

```bash
GOPATH=/__mise_test_gopath__ MISSING_COMMAND=go-test-coverage bash -c \
  'source tools/lib/setup.sh; command() { [[ "$2" != "${MISSING_COMMAND:-}" ]]; }; setup_check_tools'
```

Expected: non-zero exit and an error naming `go-test-coverage` and `mise install --locked`.
Repeat with `MISSING_COMMAND=python3`; expect the same remediation and an error naming `python3`.

- [x] **Step 5: Remove obsolete CI bootstrap dispatch**

In `tools/build.sh`, remove the `ci-setup` case and update the `install` announcement to describe
hook/LFS setup. In `tools/lib/common.sh`, remove `ci-setup` from dispatcher help and update the
`install` description. Leave other excluded dispatcher commands unchanged.

- [x] **Step 6: Verify shell and generation behavior**

Run:

```bash
shellcheck tools/lib/setup.sh tools/build.sh tools/lib/common.sh
mise exec -- ./tools/build.sh check_tools
mise exec -- ./tools/build.sh generate
```

Expected: ShellCheck exits 0 with no findings; check-tools prints `All good!`; dispatcher
generation runs `go generate ./...` and mockery with mise-managed tool versions. Commit any change
to the tracked generated mock only if it matches locked Mockery output. Run generation again and
confirm the tracked mock is unchanged on the second run. Confirm the two directives no longer
contain `@latest`:

```bash
! rg 'go run github.com/abice/go-enum@latest' internal/types/gen
```

- [x] **Step 7: Commit setup and generator changes**

```bash
git add internal/types/gen/task_type.go internal/types/gen/jetbrains_type.go \
  internal/external/gen/mocks.go tools/lib/setup.sh tools/build.sh tools/lib/common.sh
git commit -m "build: move repository setup to mise" \
  -m "Use the locked go-enum executable; reserve install for hooks and Git LFS." \
  -m "Remove the obsolete CI-only Go tool bootstrap."
```

### Task 3: Define mise tasks and prove task-graph behavior

**Files:**

- Modify: `mise.toml`
- Test: mise task graph and dry runs; keep `Makefile` during this task

- [x] **Step 1: Add all supported tasks from the spec**

Add tasks for `install`, `check`, `check_tools`, `clean`, `tidy`, `vendor`, `generate`, `build`,
`linters-build`, `release`, `release-build`, `tag`, `lint`, `vet`, `fmt`, `fix`, `test`,
`coverage`, `check-coverage`, `profile`, `bench`, `all`, and `help`.

Dispatcher-backed tasks run `./tools/build.sh <subcommand>`. Keep `tag` as `svu next --tag` and
`help` as `mise tasks ls`. `build` must create `dist/` before dispatching the build:

```toml
[tasks.build]
depends = ["fmt", "vendor", "generate"]
run = "./tools/build.sh make_out && ./tools/build.sh build"
```

Replicate the Makefile edges exactly: `vendor -> tidy`; `test -> build`; `coverage` and
`check-coverage -> test`; `lint` and `release-build -> build`; `check -> build, fmt, fix, lint,
vet, test, check-coverage`; `all -> build, check`.

- [x] **Step 2: Add only the required `wait_for` ordering edges**

Use `wait_for`, not `depends`, for tasks that must be ordered only when already selected in the
same run graph:

```toml
[tasks.tidy]
wait_for = ["fmt"]

[tasks.generate]
wait_for = ["vendor"]

[tasks.fix]
wait_for = ["build"]

[tasks.lint]
wait_for = ["fix"]

[tasks.vet]
wait_for = ["lint"]

[tasks.test]
wait_for = ["vet"]

[tasks.check]
wait_for = ["build"]
```

Do not make standalone `generate` depend on `vendor`, or standalone `vendor` depend on `fmt`.

- [x] **Step 3: Validate task definitions and inspect the graph**

Run:

```bash
mise tasks validate
mise tasks deps build check all
mise run --dry-run build
mise run --dry-run check
mise run --dry-run all
mise run --dry-run vendor
mise run --dry-run generate
```

Expected: build order is format → tidy/vendor → generate → build; check follows build, then fix,
lint, vet, test, and coverage; all builds before check; shared tasks appear once. Standalone
`vendor` and `generate` do not acquire extra prerequisites through `wait_for`.

- [x] **Step 4: Compare against Make before removing it**

Capture `make -n build`, `make -n check`, and `make -n all` while `Makefile` still exists. Compare
prerequisite and recipe order with the mise dry runs; investigate every behavioral difference.
Run `mise run build` with `OUT_DIR` pointed at a fresh temporary directory and confirm it creates
the output directory and binary before the dispatcher build runs. Do not remove any existing `dist/`
directory as part of this check.

- [x] **Step 5: Commit task graph**

```bash
git add mise.toml
git commit -m "build: add mise task graph" \
  -m "Expose supported targets and preserve Make prerequisites." \
  -m "Keep output creation, deduplication, and side-effect-sensitive ordering."
```

### Task 4: Migrate CI and GoReleaser hook

**Files:**

- Modify: `.github/workflows/ci.yml`
- Modify: `.github/workflows/common.yml`
- Modify: `.github/workflows/release.yml`
- Modify: `cfg/goreleaser.yaml`

- [x] **Step 1: Switch reusable common setup to mise**

In `common.yml`, remove the `go-version` input, add a boolean `cache` input with default `true`,
and pass it to `jdx/mise-action@v5`. Give the setup job `permissions: contents: read`. Keep
`go mod download`.

Use the reusable-workflow input in the action step:

```yaml
on:
  workflow_call:
    inputs:
      cache:
        type: boolean
        default: true

jobs:
  setup:
    permissions:
      contents: read
    steps:
      - uses: jdx/mise-action@v5
        with:
          cache: ${{ inputs.cache }}
```

- [x] **Step 2: Update CI and release callers**

In `ci.yml`, remove `go-version` inputs, replace `actions/setup-go@v5` with
`jdx/mise-action@v5`, drop `make ci-setup`, and replace `make generate`/`make build` with
`mise run generate`/`mise run build`. Keep the test matrix, coverage command, coverage action,
and schedule.

In `release.yml`, call common setup with `cache: false`; remove setup-go and `make ci-setup` from
the release job; add `jdx/mise-action@v5` there with `cache: false`. Preserve
`goreleaser/goreleaser-action@v6` as an explicit version-lock exception and do not trigger a tag
or publish artifacts during validation.

- [x] **Step 3: Repoint the GoReleaser generation hook**

In `cfg/goreleaser.yaml`, replace only `- make generate` with
`- ./tools/build.sh generate`. Keep other hooks and release settings unchanged.

- [x] **Step 4: Validate YAML and GoReleaser configuration**

Run:

```bash
yq eval '.' .github/workflows/ci.yml >/dev/null
yq eval '.' .github/workflows/common.yml >/dev/null
yq eval '.' .github/workflows/release.yml >/dev/null
mise exec -- goreleaser check --config cfg/goreleaser.yaml
```

Expected: all YAML parses, GoReleaser reports a valid configuration, and the release workflow
has `cache: false` on both mise-action paths (the reusable setup call and publisher job).

- [x] **Step 5: Commit CI migration**

```bash
git add .github/workflows/ci.yml .github/workflows/common.yml \
  .github/workflows/release.yml cfg/goreleaser.yaml
git commit -m "ci: provision toolchain with mise" \
  -m "Use the locked mise toolchain in CI." \
  -m "Keep release jobs off restored installed-tool caches."
```

### Task 5: Update contributor and repository documentation

**Files:**

- Modify: `README.md`
- Modify: `AGENTS.md`
- Modify: `.github/PULL_REQUEST_TEMPLATE/build-ci.md`
- Modify: `docs/installer.md`
- Modify: `cmd/gen-schema/main.go`
- Modify: `docs/superpowers/specs/2026-10-18-mise-adoption-design.md`

- [x] **Step 1: Rewrite active developer instructions**

In `README.md`, document installing mise once, `mise trust`, `mise install --locked`, shell
activation for bash/zsh/fish, `mise run install`, and the stale-direnv-hook warning. Replace
active Make commands with `mise run <target>` or `mise tasks ls`; preserve direct dispatcher
examples where the spec calls for them.

In `AGENTS.md`, update all active build, test, release, schema, and coverage commands; distinguish
`mise install --locked` (tools) from `mise run install` (Git hooks/LFS). Remove the obsolete
Makefile and `.nvmrc` entries from Important Files. In `docs/installer.md`, update active build
guidance but retain the explicitly historical host-tooling list. Update the PR template checklist
to `mise run` commands.

- [x] **Step 2: Update the schema generator comment**

In `cmd/gen-schema/main.go`, replace the `make generate` reference with `mise run generate`.

- [x] **Step 3: Search for stale active guidance**

Search all active guidance and workflow entry points:

```bash
rg -n 'make|Makefile|\.nvmrc|nvm|direnv' \
  README.md AGENTS.md docs/installer.md .github/PULL_REQUEST_TEMPLATE/build-ci.md \
  cmd/gen-schema/main.go .github/workflows cfg/goreleaser.yaml
```

Review every hit, including targets such as `tidy`, `vendor`, `release`, and `release-build`.
Replace active Make/nvm guidance; retain only the required stale-direnv-hook warning and the
historical host-tooling references in `docs/installer.md` explicitly allowed by the spec.

- [x] **Step 4: Commit documentation**

```bash
git add README.md AGENTS.md .github/PULL_REQUEST_TEMPLATE/build-ci.md \
  docs/installer.md cmd/gen-schema/main.go \
  docs/superpowers/specs/2026-10-18-mise-adoption-design.md
git commit -m "docs: document mise development workflow" \
  -m "Document locked installs, activation, hook setup, and mise tasks."
```

### Task 6: Remove legacy entry points after migration gates pass

**Files:**

- Delete: `Makefile`
- Delete: `.nvmrc`
- Delete: `.envrc`

- [x] **Step 1: Confirm pre-deletion gates**

Require the four-platform lock, successful local installs on macOS arm64 and Linux arm64, mise
generation, task graph/dry-run ordering, `dist/` creation, workflow YAML parsing, and active-doc
search. The user approved deferring native x64 checks to PR CI, so Makefile deletion may proceed;
do not call four-platform verification complete or approve merge until CI logs show the runner
architectures and successful locked installs cover all four tuples. If current PR runners do not
cover a tuple, request a native run before merge.

- [x] **Step 2: Delete Make, nvm, and direnv files**

Delete only the three files listed above. Do not remove direct dispatcher commands needed by
GoReleaser or `.pre-commit-config.yaml`.

- [x] **Step 3: Check for stale active references**

Search the repository for active Makefile, `.nvmrc`, nvm, and direnv activation references. Do not
sweep historical specs/plans under `docs/superpowers/{specs,plans}/`.

Expected: remaining matches are historical records, the required warning to replace a stale
direnv hook (not instructions to activate direnv), the intentional GoReleaser Action exception,
or product features unrelated to development tool management.

- [x] **Step 4: Commit cutover**

```bash
git add -A Makefile .nvmrc .envrc
git commit -m "build: replace Make and direnv with mise" \
  -m "Remove old entry points after mise tasks, lockfile, CI, and docs are ready."
```

### Task 7: Run final end-to-end verification and record baseline lint exception

**Files:**

- Verify: `mise.toml`, `mise.lock`, `.mise/locks/`, workflows, dispatcher, and docs
- Verify: no additional tracked files change unexpectedly

- [x] **Step 1: Install and validate local tools**

Run from the repository root:

```bash
mise install --locked
mise run check_tools
mise tasks ls
mise run help
```

Expected: locked install succeeds; check-tools prints `All good!`; task listing and help show the
supported targets.

- [x] **Step 2: Verify hook setup in isolated repository state**

Do not run hook installation in a linked worktree: its `.git/hooks` may be shared with the primary
checkout. Create a disposable standalone clone of the committed implementation, then run
`mise run install` there with `GIT_CONFIG_GLOBAL` pointed at a temporary config file and
`PRE_COMMIT_HOME` pointed at a temporary directory. Verify pre-commit and commit-msg hooks are
installed and run:

```bash
set -euo pipefail
test_root="$(mktemp -d)"
cleanup() {
  mise trust --untrust "$test_root/repo/mise.toml" >/dev/null 2>&1 || true
  rm -rf "$test_root"
}
trap cleanup EXIT
git clone --no-hardlinks . "$test_root/repo"
touch "$test_root/gitconfig"
mkdir -p "$test_root/pre-commit"
(
  set -euo pipefail
  cd "$test_root/repo"
  mise trust
  GIT_CONFIG_GLOBAL="$test_root/gitconfig" PRE_COMMIT_HOME="$test_root/pre-commit" mise run install
  ./tools/build.sh check-lfs-hook
)
```

Expected: the LFS pre-push hook check passes; Git config, pre-commit cache, and hooks remain isolated
to the disposable test. The trap removes only the directory returned by `mktemp` and clears the
temporary mise trust entry.

- [x] **Step 3: Run build, tests, coverage, vet, and release tasks**

Run `mise run all` separately and confirm it stops at the known lint finding. Then run the
independent checks:

```bash
mise run test
mise run check-coverage
mise run vet
mise run release
mise run release-build
mise exec -- goreleaser check --config cfg/goreleaser.yaml
mise exec -- goreleaser release --snapshot --clean --config cfg/goreleaser.yaml --skip publish,sign
```

Expected: `mise run all` reaches lint and reports the pre-existing G703 finding in
`internal/external/jetbrains.go:34` under golangci-lint 2.15.0. The host-installed 2.15.0 reports
the same finding against unchanged source; preserve product scope and do not alter lint config.
Independent `test`, `check-coverage`, and `vet` tasks pass. Local release tasks produce snapshot
artifacts, the production GoReleaser config validates, and its generation hook runs the dispatcher.
Snapshot command skips publish/sign. Do not publish, push tags, or run a production release.

- [x] **Step 4: Reconfirm lock and graph invariants after cutover**

Run:

```bash
mise install --dry-run --locked
mise tasks validate
mise tasks deps build check all
mise run --dry-run build
mise run --dry-run check
mise run --dry-run all
git diff --check
git status --short
```

Expected: all tools are installed for the current macOS arm64 host; graph order and shared-task
de-duplication remain correct; whitespace check passes; no generated outputs or other unexpected
changes remain. The lockfile previously resolved all 68 tool/platform targets. Native x64 installs
remain deferred to PR CI and must be confirmed before merge.

- [ ] **Step 5: Commit the tracked implementation plan**

After updating completed checkboxes and leaving the x64 PR-CI gate explicitly pending, commit the
plan:

```bash
git add docs/superpowers/plans/2026-10-29-mise-adoption.md
git commit -m "docs: add mise adoption implementation plan" \
  -m "Record migration steps and leave native x64 PR-CI checks pending."
```

## Completion Criteria

- `mise.toml` and `mise.lock` cover Linux/macOS x64/arm64; locked installation passes on macOS
  arm64 and native Linux arm64. Native x64 installs remain a PR-CI merge gate.
- All supported task names, dependency edges, `wait_for` ordering, and fresh-clone `dist/` creation
  match the spec.
- Go generation uses mise-managed `go-enum`; `check_tools` validates `go-test-coverage` and gives
  the locked-install remedy. Local `go-test-coverage` is pinned to Go-1.26-compatible `2.19.0`;
  installation must not auto-download a Go 1.27 toolchain.
- CI uses `jdx/mise-action@v5`; release setup and publishing disable installed-tool caching; the
  GoReleaser and coverage action exceptions remain explicit.
- Makefile, `.nvmrc`, and `.envrc` are removed only after task graph and CI checks pass.
- Developer docs and PR checklist use mise commands; historical specs/plans remain unchanged.
- Tests, check-coverage, vet, release tasks, config checks, and final task graph checks pass.
- `mise run all` retains only the confirmed pre-existing G703 lint failure; resolve that separately
  from this no-product-change migration.
