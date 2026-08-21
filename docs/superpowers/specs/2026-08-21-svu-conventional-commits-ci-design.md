# Design: svu Semantic Versioning, Conventional Commit Enforcement, and CI Build/Release Pipelines

Date: 2026-08-21
Branch: `feature/svu`
Status: Approved

## Goal

Adopt [svu](https://github.com/caarlos0/svu) for semantic versioning driven by
[Conventional Commits](https://www.conventionalcommits.org), enforce commit
format locally via pre-commit, run the build automatically on PRs and pushes
to `main`, and publish releases via GoReleaser when a version tag is pushed.

Non-goals:

- No auto-tagging on merge to `main`; tagging stays a deliberate local action.
- No changes to how the Makefile derives `DUM_VERSION` for local builds.

## Current State

- Tags: `v0.2.2`, `v1.0.0`. History uses squash merges with conventional-style
  PR titles (e.g., `refactor: rename module ... (#28)`).
- `.svu.yml` already staged on this branch:
  ```yaml
  always: true # increment patch if no commits trigger version change
  v0: true # prevent major version increments if current version is still v0
  ```
- `.commitlintrc.yml` exists, extending `@commitlint/config-conventional`.
- pre-commit runs `gitlint` at `commit-msg` stage (general hygiene, not
  conventional format). Both hooks will coexist.
- `ci.yml` triggers only on `workflow_dispatch`; no release workflow exists.
- Node is available locally via nvm (v22); repo tooling does not manage it.

## Section 1 — svu Versioning

- Add `svu` to the Homebrew install list in `setup_install_tools`
  (`tools/lib/setup.sh`) and to `required_commands` in `setup_check_tools`.
- New Makefile target:
  ```make
  .PHONY: tag
  tag: ## compute next semver from conventional commits and create git tag
  	@svu next --tag
  ```
- `make tag` computes the next version from conventional commits since the
  last tag and creates an annotated git tag locally. The developer pushes it
  manually (`git push origin <tag>`), which triggers the release workflow.
- `.svu.yml` semantics: `always: true` guarantees at least a patch bump on
  every tag; `v0: true` caps the major version at 0 until the project opts
  into 1.x+ bumps.

## Section 2 — Conventional Commit Enforcement

### Local (pre-commit)

```yaml
- repo: https://github.com/alessandrojcm/commitlint-pre-commit-hook
  rev: v9.18.0
  hooks:
    - id: commitlint
      stages: [commit-msg]
      additional_dependencies: ['@commitlint/config-conventional']
```

- Requires node/npm at hook runtime; provided by nvm in the developer shell.
- `gitlint` remains alongside; concerns do not overlap (hygiene vs format).

### Node via nvm (not brew)

- New `.nvmrc` pinning node major version `22`.
- `setup_check_tools` gains an nvm-aware check: if `node` is missing, source
  `$NVM_DIR/nvm.sh` and run `nvm install` per `.nvmrc`; otherwise fail with an
  instruction to install nvm. Brew list unchanged.
- Node is a local-only dependency (pre-commit hook runtime); CI workflows do
  not need it.

## Section 3 — Automated Build

`ci.yml` triggers become:

```yaml
on:
  workflow_dispatch:
  pull_request:
  push:
    branches: [main]
  schedule:
    - cron: '0 6 * * 1' # Mondays 06:00 UTC
```

The existing test/coverage job is unchanged; no commitlint job is added to
CI. PR titles are not enforced — with `always: true`, an off-format squash
title only risks an imprecise bump level, never a broken release.

The weekly schedule is a build-rot sanity check: it re-runs the same test job
against `main` so dependency or runner drift is caught even when nothing was
merged. It never publishes — releases fire only on `v*` tag pushes. Note:
GitHub disables scheduled workflows after 60 days of repo inactivity; any
push or manual dispatch re-enables them.

## Section 4 — Release Pipeline

New `.github/workflows/release.yml`:

```yaml
on:
  push:
    tags: ['v*']
permissions:
  contents: write
jobs:
  setup:
    uses: ./.github/workflows/common.yml # reuses ci-setup (go-enum etc.)
    with:
      go-version: '1.26'
  release:
    needs: setup
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0 # goreleaser needs full history + tags
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - uses: goreleaser/goreleaser-action@v6
        with:
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

- Reusing `common.yml` matters: goreleaser's `before` hooks run
  `go generate ./...`, which requires generated-tool binaries installed by
  `ci-setup`.
- `GITHUB_TOKEN` suffices for releases on a public repository.
- Existing `cfg/goreleaser.yaml` reads the version from the pushed tag.

### Build matrix expansion

`cfg/goreleaser.yaml` gains linux/amd64 support via the full cross product:

```yaml
goos:
  - darwin
  - linux
goarch:
  - arm64
  - amd64
```

Four release targets: darwin/arm64, darwin/amd64, linux/arm64, linux/amd64 —
matching the README's macOS/Linux claim. The archives template's existing
linux/amd64 branches become live; windows branches remain dead (windows is
still unsupported). The stale Makefile `release` target comment ("all
platforms ... arm64/amd64/386") is corrected to list the actual four targets.

### Packaging and checksums

Each target binary is packaged into a tar.gz archive together with
`cfg/installer.schema.json`. GoReleaser generates a SHA256 checksums file
covering every archive and uploads it to the GitHub Release alongside them,
so downloads can be verified. Made explicit in `cfg/goreleaser.yaml`:

```yaml
checksum:
  name_template: "{{ .ProjectName }}_{{ .Version }}_checksums.txt"
  algorithm: sha256
```

Resulting release assets per tag: four `dum_<os>_<arch>_<version>-<commit>.tar.gz`
archives plus `dum_<version>_checksums.txt`.

## Section 5 — Documentation Updates

### AGENTS.md

1. **Build/Lint/Test Commands**: add `make tag` to the primary commands block:
   ```bash
   make tag            # Compute next semver from conventional commits and create git tag
   ```
2. **New "Versioning and Releases" section** after Code Generation:
   - Versions follow [Conventional Commits](https://www.conventionalcommits.org);
     enforced locally by commitlint via pre-commit.
   - [svu](https://github.com/caarlos0/svu) computes the next version; config
     lives in `.svu.yml` (`always: true`, `v0: true`).
   - Release flow: squash-merge PRs with conventional titles → `make tag` →
     `git push origin <tag>` → GoReleaser publishes via
     `.github/workflows/release.yml`.
3. **Development Setup**: note node is managed by nvm per `.nvmrc`
   (`nvm install` picks it up), and that `make install` installs `svu` via brew.
4. **Important Files**: add `.svu.yml`, `.commitlintrc.yml`, `.nvmrc`,
   `.pre-commit-config.yaml`, and `.github/workflows/release.yml`.

### README.md

1. **Badges** directly under the `# dum` title (style follows
   [svu's README](https://github.com/caarlos0/svu/blob/main/README.md)):
   - **Release**: `img.shields.io/github/v/release/alanktwong/dum`, linking to
     the releases page.
   - **License**: `img.shields.io/github/license/alanktwong/dum` (dynamic,
     reads the existing `LICENSE`), linking to the license file.
   - **Build status**: GitHub Actions native badge for `ci.yml`
     (`github.com/alanktwong/dum/actions/workflows/ci.yml/badge.svg`),
     linking to the Actions page.
   - **Conventional Commits**: static
     `img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg`, linking
     to https://www.conventionalcommits.org.
   - **GoDoc**: official `https://pkg.go.dev/badge/github.com/alanktwong/dum.svg`,
     linking to `https://pkg.go.dev/github.com/alanktwong/dum`. Docs publish
     automatically via pkg.go.dev on first fetch; no suffix needed while the
     module stays v1.x (add `/v2+` path suffix if major bumps happen).
2. **Development > Build Commands**: add `make tag` line matching AGENTS.md.
3. **New "Versioning and Releases" subsection** under Development:
   - Conventional Commits required for all contributions; commitlint validates
     local commit messages via pre-commit.
   - `make tag` computes and tags the next semver from merged conventional
     commits; pushing the tag triggers an automated GoReleaser publish.
   - Contributors need node (nvm-managed, see `.nvmrc`) for the pre-commit
     commitlint hook.

## End-to-End Flow

1. Contributor opens PR → tests + build run automatically.
2. Squash-merge with conventional title lands on `main` → full CI re-runs.
3. Maintainer runs `make tag` → svu computes next semver, creates local tag.
4. `git push origin <tag>` → `release.yml` builds and publishes via goreleaser.

## Error Handling

| Failure | Detection | Result |
| --- | --- | --- |
| node missing locally | `check_tools` nvm-aware check | Blocked early with install instruction |
| Non-conventional local commit | pre-commit `commit-msg` hook | Commit rejected before creation |
| No conventional commits since tag | svu with `always: true` | Patch bump instead of failure |
| Major bump while on v0.x | svu with `v0: true` | Version stays within 0.x |

## Verification Plan

- `shellcheck` (already hooked) covers `tools/lib/setup.sh` changes.
- Open a draft PR from the feature branch and watch the workflow run.
- Run `svu next` (dry-run, no `--tag`) before creating the first real tag.
- Confirm a deliberately malformed commit message is rejected by the local
  hook.
