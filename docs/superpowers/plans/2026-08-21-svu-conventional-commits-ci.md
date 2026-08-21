# svu + Conventional Commits + CI/Release Pipelines Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use subagent-driven-development (recommended) or
> executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for
> tracking.

**Goal:** Ship release management: svu-driven semver tagging, local conventional-commit
enforcement, automated CI builds (PR/push/weekly), and GoReleaser releases with checksums on tag
push.

**Architecture:** All changes are configuration, workflow, and documentation files — no Go code.
svu computes versions from conventional commits; a Makefile target tags locally; pushing a `v*`
tag triggers a GitHub Actions workflow that runs GoReleaser to cross-compile four targets,
package archives, generate SHA256 checksums, and publish.

**Tech Stack:** svu, commitlint (`@commitlint/config-conventional`), pre-commit, nvm/node 22,
GitHub Actions, GoReleaser v2.

**Spec:** `docs/superpowers/specs/2026-08-21-svu-conventional-commits-ci-design.md`
**Issue:** https://github.com/alanktwong/dum/issues/32

**Branch:** continue on existing `feature/svu`.

**Commit convention:** repo gitlint requires a body (rule B6). Every commit message below has one.
Use them verbatim.

---

## Task Structure

### Task 1: GoReleaser build matrix + checksums

**Files:**

- Modify: `cfg/goreleaser.yaml:14-17` (goos/goarch)
- Modify: `cfg/goreleaser.yaml` (insert `checksum:` section after `archives:`)
- Modify: `Makefile:82` (stale comment)

- [ ] **Step 1: Expand goos/goarch**

In `cfg/goreleaser.yaml`, replace:

```yaml
    goos:
      - darwin
    goarch:
      - arm64
```

with:

```yaml
    goos:
      - darwin
      - linux
    goarch:
      - arm64
      - amd64
```

- [ ] **Step 2: Add explicit checksum section**

Insert between the `archives:` block and the `changelog:` block:

```yaml
checksum:
  name_template: "{{ .ProjectName }}_{{ .Version }}_checksums.txt"
  algorithm: sha256
```

- [ ] **Step 3: Fix stale Makefile comment**

In `Makefile`, replace:

```make
release: ## compile release binaries for all platforms (darwin/linux, arm64/amd64/386)
```

with:

```make
release: ## compile release binaries for darwin/linux on arm64/amd64
```

- [ ] **Step 4: Validate config**

Run: `goreleaser check --config cfg/goreleaser.yaml`
Expected: `configuration is valid` (or equivalent success line).

- [ ] **Step 5: Snapshot-build all four targets locally**

Run: `goreleaser release --snapshot --clean --config cfg/goreleaser.yaml --skip publish,sign`
Expected: builds succeed for darwin/arm64, darwin/amd64, linux/arm64, linux/amd64;
`dist/` contains four `.tar.gz` archives plus `dum_0.0.0-SNAPSHOT-..._checksums.txt`.
Verify checksum file lists all four archives:
`grep -c ".tar.gz" dist/*checksums*.txt` → `4`.
Note: mid-plan the tree may be dirty, so artifact names can carry a `-dirty` suffix — that is
expected and not a failure.

- [ ] **Step 6: Commit**

```bash
git add cfg/goreleaser.yaml Makefile
git commit -m "build: expand release matrix and add checksums" -m "Full goos x goarch cross product yields darwin/linux on
arm64/amd64. Explicit sha256 checksums file ships with every
release so downloads are verifiable."
```

### Task 2: svu tooling in setup.sh

**Files:**

- Modify: `tools/lib/setup.sh:13-25` (brew list)
- Modify: `tools/lib/setup.sh:67-82` (required_commands)

- [ ] **Step 1: Add svu to brew install list**

In `setup_install_tools`, insert `svu` alphabetically between `shellcheck` and `trufflehog`:

```bash
    brew install \
        direnv \
        git-lfs \
        go \
        golangci-lint \
        goreleaser \
        jq \
        mockery \
        pre-commit \
        shellcheck \
        svu \
        trufflehog \
        yamlfmt \
        yq
```

- [ ] **Step 2: Add svu to required_commands**

In `setup_check_tools`, add `"svu"` between `"shellcheck"` and `"trufflehog"`:

```bash
    local -a required_commands=(
        "direnv"
        "git-lfs"
        "go"
        "go-enum"
        "goimports"
        "golangci-lint"
        "goreleaser"
        "jq"
        "mockery"
        "pre-commit"
        "shellcheck"
        "svu"
        "trufflehog"
        "yamlfmt"
        "yq"
    )
```

- [ ] **Step 3: Verify with shellcheck**

Run: `shellcheck tools/lib/setup.sh`
Expected: no output, exit 0.

- [ ] **Step 4: Commit**

```bash
git add tools/lib/setup.sh
git commit -m "build: add svu to dev tooling" -m "svu computes semver from conventional commits; installed
via brew and required by check_tools like the rest of the
toolchain."
```

### Task 3: .nvmrc + nvm-aware node check

**Files:**

- Create: `.nvmrc`
- Modify: `tools/lib/setup.sh:54+` (`setup_check_tools`)

- [ ] **Step 1: Create .nvmrc**

Write exactly:

```
22
```

(trailing newline included — end-of-file-fixer enforces it)

- [ ] **Step 2: Add nvm bootstrap to setup_check_tools**

Insert immediately after the Homebrew availability check and before the `$GOPATH/bin` PATH check:

```bash
    # Ensure node is available via nvm before checking commands.
    if ! command -v node &>/dev/null; then
        export NVM_DIR="${NVM_DIR:-$HOME/.nvm}"
        if [[ -s "$NVM_DIR/nvm.sh" ]]; then
            # shellcheck disable=SC1091
            . "$NVM_DIR/nvm.sh"
            nvm install
        else
            echo >&2 "node is required but nvm was not found at $NVM_DIR."
            echo >&2 "Install nvm (https://github.com/nvm-sh/nvm), run 'nvm install', then retry."
            return 1
        fi
    fi
```

Note: `nvm install` with no argument reads `.nvmrc` from the working directory; build.sh runs
from the repository root.

- [ ] **Step 3: Add node/npm to required_commands**

Add `"node"` and `"npm"` between `"mockery"` and `"pre-commit"` in the array from Task 2.

- [ ] **Step 4: Verify with shellcheck**

Run: `shellcheck tools/lib/setup.sh`
Expected: no output, exit 0.

- [ ] **Step 5: Run check_tools**

Run: `./tools/build.sh check_tools`
Expected: `All good!` (node resolves via nvm or already on PATH).

- [ ] **Step 6: Commit**

```bash
git add .nvmrc tools/lib/setup.sh
git commit -m "build: manage node via nvm" -m "commitlint's pre-commit hook needs node at hook runtime.
Pin major version in .nvmrc and teach check_tools to source
nvm and install it when missing instead of adding a brew dep."
```

### Task 4: Makefile tag target

**Files:**

- Modify: `Makefile:85-87` (after `release-build`)

- [ ] **Step 1: Add tag target**

Insert after the `release-build` block:

```make
.PHONY: tag
tag: ## compute next semver from conventional commits and create git tag
	@svu next --tag
```

Use a tab character before `@svu` (Makefile recipes require tabs).

- [ ] **Step 2: Ensure svu is installed**

Task 2 only registers svu in setup.sh. Install it now if missing:
Run: `command -v svu || brew install svu`
Expected: prints svu's path, or installs it.

- [ ] **Step 3: Dry-run verification**

Run: `svu next`
Expected: prints next version computed from commits since last tag (e.g., `v1.0.1` given
`always: true`). No tag created.

- [ ] **Step 4: Verify make wiring without tagging**

Run: `make -n tag`
Expected: prints `svu next --tag`. Do NOT run `make tag` until the PR merges to main — tagging
on the feature branch would compute from the wrong base.

- [ ] **Step 5: Commit**

```bash
git add Makefile
git commit -m "build: add make tag target" -m "Delegates to svu next --tag so maintainers create release
tags from merged conventional commits with one command."
```

### Task 5: commitlint pre-commit hook

**Files:**

- Modify: `.pre-commit-config.yaml:42-46` (after gitlint block)

- [ ] **Step 1: Confirm pinned rev exists**

Run: `git ls-remote --tags https://github.com/alessandrojcm/commitlint-pre-commit-hook | grep -F 'refs/tags/v9' | tail -5`
If `v9.18.0` is absent, substitute the newest `v9.x.y` tag in Step 2.

- [ ] **Step 2: Add hook entry**

Insert after the gitlint repo block (keeping both commit-msg hooks):

```yaml
  - repo: https://github.com/alessandrojcm/commitlint-pre-commit-hook
    rev: v9.18.0
    hooks:
      - id: commitlint
        stages: [commit-msg]
        additional_dependencies: ['@commitlint/config-conventional']
```

- [ ] **Step 3: Install commit-msg hook type**

Run: `pre-commit install --hook-type commit-msg`
Expected: `pre-commit installed at .git/hooks/commit-msg` (idempotent if already present).

- [ ] **Step 4: Verify rejection of bad message**

```bash
printf 'bad message\n' > /tmp/opencode/bad-msg.txt
pre-commit run commitlint --commit-msg-file /tmp/opencode/bad-msg.txt
```
Expected: FAIL with commitlint output about subject format.

- [ ] **Step 5: Verify acceptance of good message**

```bash
printf 'feat: valid conventional subject\n' > /tmp/opencode/good-msg.txt
pre-commit run commitlint --commit-msg-file /tmp/opencode/good-msg.txt
```
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add .pre-commit-config.yaml
git commit -m "build: enforce conventional commits via commitlint" -m "Runs at commit-msg stage alongside gitlint, which keeps
general hygiene rules. Config lives in .commitlintrc.yml."
```

### Task 6: ci.yml triggers + weekly schedule

**Files:**

- Modify: `.github/workflows/ci.yml:3-4`

- [ ] **Step 1: Replace triggers**

Replace:

```yaml
on:
  workflow_dispatch:
```

with:

```yaml
on:
  workflow_dispatch:
  pull_request:
  push:
    branches: [main]
  schedule:
    - cron: '0 6 * * 1' # Mondays 06:00 UTC
```

- [ ] **Step 2: Validate YAML**

Run: `yq eval '.on' .github/workflows/ci.yml`
Expected: prints all four trigger keys (`workflow_dispatch`, `pull_request`, `push`,
`schedule`) without parse errors.

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/ci.yml
git commit -m "ci: run tests on PRs, pushes to main, weekly" -m "Weekly Monday run revalidates the build against main to
catch dependency or runner drift; scheduled runs never
publish."
```

### Task 7: release.yml

**Files:**

- Create: `.github/workflows/release.yml`

- [ ] **Step 1: Create workflow file**

Write exactly:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  setup:
    uses: ./.github/workflows/common.yml
    with:
      go-version: '1.26'

  release:
    needs: setup
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - uses: goreleaser/goreleaser-action@v6
        with:
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

Notes: `fetch-depth: 0` gives goreleaser full history and tags; reusing `common.yml` installs
the generated-tool binaries that goreleaser's `before` hooks need.

- [ ] **Step 2: Validate YAML**

Run: `yq eval '.jobs | keys' .github/workflows/release.yml`
Expected: `setup`, `release`.

- [ ] **Step 3: Commit**

```bash
git add .github/workflows/release.yml
git commit -m "ci: add tag-push release pipeline" -m "GoReleaser publishes archives and checksums when a v*
tag lands; GITHUB_TOKEN suffices for public-repo releases."
```

### Task 8: AGENTS.md updates

**Files:**

- Modify: `AGENTS.md:26-36` (commands block)
- Modify: `AGENTS.md:45-53` (after Code Generation)
- Modify: `AGENTS.md:80` (Development Setup comment)
- Modify: `AGENTS.md:195-201` (Important Files)

- [ ] **Step 1: Add make tag to commands block**

After the `make generate` line in the primary commands block, add:

```bash
make tag             # Compute next semver from conventional commits and create git tag
```

- [ ] **Step 2: Add Versioning and Releases section**

Insert after the Code Generation section (before Coverage Threshold):

```markdown
### Versioning and Releases

Versions follow [Conventional Commits](https://www.conventionalcommits.org), enforced locally by
commitlint via pre-commit. [svu](https://github.com/caarlos0/svu) computes the next version from
merged commits; configuration lives in `.svu.yml` (`always: true`, `v0: true`).

Release flow: squash-merge PRs with conventional titles, run `make tag`, then
`git push origin <tag>`. Pushing a tag triggers GoReleaser via `.github/workflows/release.yml`,
which publishes per-platform archives and a SHA256 checksums file to the GitHub Release.
```

- [ ] **Step 3: Update Development Setup comment**

Replace:

```bash
make install    # Install tools (go-enum, goimports, gocov, golangci-lint, mockery)
```

with:

```bash
make install    # Install tools (brew deps incl. svu/goreleaser; go-enum, goimports, gocov, golangci-lint, mockery). Node comes from nvm per .nvmrc
```

- [ ] **Step 4: Extend Important Files**

Add these lines to the Important Files list:

```markdown
- `.svu.yml` - svu versioning config
- `.commitlintrc.yml` - commitlint conventional-commits config
- `.nvmrc` - node version pin for the commitlint pre-commit hook
- `.pre-commit-config.yaml` - local hook enforcement (incl. commitlint at commit-msg)
- `.github/workflows/release.yml` - tag-push release pipeline
```

- [ ] **Step 5: Commit**

```bash
git add AGENTS.md
git commit -m "docs: document versioning and release flow" -m "Agents need the make tag flow, tooling notes, and new
config files surfaced in Important Files."
```

### Task 9: README.md badges + versioning docs

**Files:**

- Modify: `README.md:1-3` (badges under title)
- Modify: `README.md:164-177` (Build Commands block)
- Modify: `README.md` (new subsection after Build Commands)

- [ ] **Step 1: Add badges**

Directly after `# dum` and before the intro paragraph, insert:

```markdown
[![Release](https://img.shields.io/github/v/release/alanktwong/dum)](https://github.com/alanktwong/dum/releases)
[![License](https://img.shields.io/github/license/alanktwong/dum)](LICENSE)
[![Build status](https://github.com/alanktwong/dum/actions/workflows/ci.yml/badge.svg)](https://github.com/alanktwong/dum/actions/workflows/ci.yml)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://www.conventionalcommits.org)
[![GoDoc](https://pkg.go.dev/badge/github.com/alanktwong/dum.svg)](https://pkg.go.dev/github.com/alanktwong/dum)
```

- [ ] **Step 2: Add make tag to Build Commands**

After the `make clean` line in the Build Commands block, add:

```shell
make tag             # Compute next semver from conventional commits and create git tag
```

- [ ] **Step 3: Add Versioning and Releases subsection**

Insert after the Build Commands block:

```markdown
### Versioning and Releases

Contributions follow [Conventional Commits](https://www.conventionalcommits.org); commitlint
validates local commit messages via pre-commit (node required — managed by nvm, see `.nvmrc`).

To cut a release: run `make tag` ([svu](https://github.com/caarlos0/svu) computes the next semver
from merged conventional commits), then `git push origin <tag>`. Pushing the tag triggers
GoReleaser, which publishes per-platform archives and a SHA256 checksums file to the GitHub
Release.
```

- [ ] **Step 4: Verify badge URLs resolve**

Run: `curl -sI https://img.shields.io/github/v/release/alanktwong/dum | head -1` → `HTTP/2 200`
(same spot-check for the license badge; pkg.go.dev badge may 404 until first index fetch — acceptable).

- [ ] **Step 5: Commit**

```bash
git add README.md
git commit -m "docs: add badges and release instructions to README" -m "Release, license, build status, conventional commits, and
GoDoc badges styled after svu's README, plus the maintainer
release flow."
```

### Task 10: Full verification + PR

**Files:** none (verification only)

- [ ] **Step 1: Run full quality suite**

Run: `make check`
Expected: build, fmt, lint, vet, test, coverage threshold all pass.

- [ ] **Step 2: Confirm tree state**

Run: `git status --short && git log --oneline origin/main..HEAD`
Expected: clean tree; this branch's commits listed.

- [ ] **Step 3: Push branch**

Run: `git push -u origin feature/svu`

- [ ] **Step 4: Open PR**

Use the `stacked-prs` skill conventions for a single-PR stack (base `main`, stack overview in
description):

```bash
gh pr create --repo alanktwong/dum --base main --head feature/svu \
  --title "feat: svu release management with CI and tag-push releases" \
  --body-file - <<'EOF'
## Summary

Implements #32: svu-driven semantic versioning, local conventional-commit enforcement, automated
CI builds, and GoReleaser releases with checksums.

Design: docs/superpowers/specs/2026-08-21-svu-conventional-commits-ci-design.md

## Changes

- goreleaser: full darwin/linux x arm64/amd64 matrix + explicit sha256 checksums
- svu tooling (brew + check_tools) and `make tag` target
- node via nvm (`.nvmrc`, nvm-aware check_tools) for the commitlint hook
- commitlint at commit-msg stage (alongside gitlint)
- ci.yml: PR/push-to-main/weekly-cron triggers
- release.yml: tag-push pipeline publishing archives + checksums
- AGENTS.md / README.md documentation and badges

Closes #32
EOF
```

- [ ] **Step 5: Post-merge smoke test (after squash-merge)**

On `main`: run `svu next` (expect patch bump), then `make tag && git push origin <tag>`, and
watch the Release workflow produce four archives + `dum_<version>_checksums.txt`.
