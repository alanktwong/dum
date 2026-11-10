# Tarball Command Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use subagent-driven-development
> (recommended) or executing-plans to implement this plan task-by-task. Steps use
> checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `dum tarball backup` and `dum tarball restore` with YAML inputs,
SHA-256 verification, a content manifest, and safe staged extraction.

**Architecture:** Keep Cobra wiring in `cmd/dum` and archive/input/manifest
logic in a new Cobra-free `internal/tarball` package. Use Go's tar, gzip, hash,
filesystem, and existing YAML libraries; stage and validate before destination
writes.

**Tech Stack:** Go 1.26, Cobra, `archive/tar`, `compress/gzip`, `crypto/sha256`,
`gopkg.in/yaml.v3`, testify.

---

## Delivery

Deliver as one cohesive PR on the current feature branch; no stacked PRs. Keep
the implementation split into independently testable package and CLI tasks.
Use test-first changes and small Conventional Commit commits.

## File Map

- Modify `cmd/dum/root.go` to register the new `tarball` command.
- Modify `cmd/dum/root_test.go` to assert registration and subcommands.
- Modify `.golangci.yml` and `AGENTS.md` to register the `tb` import alias for
  `internal/tarball`.
- Create `cmd/dum/tarball.go` for Cobra command construction, flags, and CLI
  orchestration.
- Create `cmd/dum/tarball_picker.go` for `fzf`/`gum` selection; keep subprocess
  and terminal behavior out of `internal/`.
- Create `cmd/dum/tarball_test.go` for command arguments, flag behavior, and
  list/explicit-archive flows.
- Create `internal/tarball/service.go` for service/options types and injected
  logging/ownership seams.
- Create `internal/tarball/input.go` for strict `input.yml` parsing and path
  resolution; `input_test.go` covers it.
- Create `internal/tarball/manifest.go` for versioned inventory serialization
  and validation; `manifest_test.go` covers it.
- Create `internal/tarball/paths.go` for path normalization, source-root
  mapping, symlink classification, and overlap checks; `paths_test.go` covers
  it.
- Create `internal/tarball/backup.go` for archive, checksum, manifest, and
  atomic publication; `backup_test.go` covers it.
- Create `internal/tarball/restore.go` for sidecar loading, checksum snapshot,
  manifest/archive comparison, and top-level restore workflow.
- Create `internal/tarball/extract.go` for safe tar validation and private
  staging; `restore_test.go` covers it with Go-generated fixtures.
- Create `internal/tarball/apply.go` for destination preflight, conflict policy,
  metadata, and application of staged entries; `apply_test.go` covers it.
- Create `internal/tarball/discovery.go` for list-based archive matching;
  `discovery_test.go` covers it.
- Create `internal/tarball/roundtrip_test.go` for end-to-end backup-to-restore
  compatibility tests.
- Create `internal/tarball/test_helpers_test.go` for deterministic tar/gzip and
  sidecar fixtures. Do not invoke host `tar`.
- Modify `README.md` to document the command, YAML input, sidecars, and restore
  safety flags.

No `go.mod` change is needed: Cobra, testify, and `yaml.v3` are already present.
Do not regenerate the installer JSON schema; `input.yml` is separate from
`dum.yml`.

## Task Structure

### Task 1: Parse the YAML input file

**Files:**

- Create: `internal/tarball/input.go`
- Test: `internal/tarball/input_test.go`

- [x] **Step 1: Write input parser tests**

  Cover a valid top-level `paths` list, absolute entries, relative entries
  resolved from the YAML file's directory, paths with spaces, missing/empty
  `paths`, empty items, malformed YAML, unknown fields, and duplicate resolved
  paths. Also cover a relative `--input` filename resolving from the current
  working directory. A missing `--input` is a CLI concern and is covered in
  Task 6.

- [x] **Step 2: Run the focused tests and confirm they fail**

  Run: `go test ./internal/tarball -run TestLoadInput -v`
  Expected: FAIL because the package/parser is not implemented yet.

- [x] **Step 3: Implement strict YAML loading**

  Add a typed input struct with `Paths []string`. Implement
  `LoadInput(ctx context.Context, inputPath string)` using `yaml.v3` with
  unknown-field rejection. Resolve a relative `inputPath` from the process
  working directory and relative entries from the input file's directory.
  Normalize entries, reject empty/duplicate results, and wrap file/parse errors.

- [x] **Step 4: Re-run the focused tests**

  Run: `go test ./internal/tarball -run TestLoadInput -v`
  Expected: PASS for valid, relative, absolute, whitespace, and invalid cases.

- [x] **Step 5: Commit the input loader**

  ```bash
  git add internal/tarball/input.go internal/tarball/input_test.go
  git commit -m "feat(tarball): parse input YAML"
  ```

### Task 2: Define manifest and path invariants

**Files:**

- Create: `internal/tarball/service.go`
- Create: `internal/tarball/manifest.go`
- Create: `internal/tarball/paths.go`
- Test: `internal/tarball/manifest_test.go`
- Test: `internal/tarball/paths_test.go`

- [x] **Step 1: Write manifest and path tests**

  Cover schema version 1 serialization, sorted relative member paths, entry
  types, optional `uid`/`gid` appearing only as a pair, symlink target/warning
  fields, source-root mapping inside/outside `$HOME`, unsafe absolute/escaping
  links, duplicate member paths, and case-insensitive source-basename collisions.

- [x] **Step 2: Run the focused tests and confirm they fail**

  Run: `go test ./internal/tarball -run 'Test(Manifest|Validate|Parse|ArchivePath|Classify)' -v`
  Expected: FAIL because the model and helpers are not implemented yet.

- [x] **Step 3: Implement service types and manifest validation**

  Define `Service`, backup/restore option structs, entry type constants, and
  versioned manifest structs. Validate normalized relative paths, required
  fields, supported entry types, paired UID/GID values, and symlink metadata.
  Keep exported service methods context-first and use the existing logging
  abstraction for warnings.

- [x] **Step 4: Implement archive path helpers**

  Map sources inside `$HOME` to HOME-relative members and outside sources to a
  basename root. Use slash-separated archive paths independent of host path
  separators. Add checks for traversal, duplicate/case-fold collisions, unsafe
  link targets, and backup-target/source overlap. Never follow symlinks while
  classifying a source tree.

- [x] **Step 5: Re-run the focused tests**

  Run: `go test ./internal/tarball -run 'Test(Manifest|Validate|Parse|ArchivePath|Classify)' -v`
  Expected: PASS for manifest round trips and path/symlink invariants.

- [x] **Step 6: Commit the manifest and path model**

  ```bash
  git add internal/tarball/service.go internal/tarball/manifest.go \
    internal/tarball/paths.go internal/tarball/manifest_test.go \
    internal/tarball/paths_test.go
  git commit -m "feat(tarball): define archive manifest"
  ```

### Task 3: Create backup archives and sidecars

**Files:**

- Modify: `internal/tarball/service.go`
- Create: `internal/tarball/backup.go`
- Test: `internal/tarball/backup_test.go`
- Test helper: `internal/tarball/test_helpers_test.go`

- [x] **Step 1: Write backup round-trip artifact tests**

  Create temporary file/directory sources. Assert one timestamped `.tgz`,
  `.tgz.sha256`, and `.tgz.manifest.yaml` per source; verify the checksum with
  Go's `crypto/sha256`; parse tar members and YAML without a shell command.
  Assert in-home member layout and external basename-root layout.

- [x] **Step 2: Run the backup tests and confirm they fail**

  Run: `go test ./internal/tarball -run TestBackup -v`
  Expected: FAIL because backup writing is not implemented.

- [x] **Step 3: Implement archive writing and inventory generation**

  Before creating temporary outputs, preflight every selected source, all
  basename collisions, every final archive/sidecar path, and every
  target/source overlap. A preflight failure for any source must create no
  artifacts for any source in that invocation.

  Walk entries without following links. Write regular files, directories, and
  symlinks only; fail on unsupported special entries. Populate tar modes and
  modification times, clear setuid/setgid bits, and record numeric UID/GID in
  headers and manifest when available; otherwise write neutral `0/0` IDs and
  omit both manifest fields. Log absolute/root-escaping symlink warnings and
  record their literal targets and unsafe state. Sort manifest entries
  deterministically.

- [x] **Step 4: Implement checksums, staging, and no-overwrite publication**

  Build all artifacts in temporary files in the target directory. Hash the
  completed archive bytes and write a strict lowercase SHA-256 sidecar. Publish
  both sidecars before the archive completion marker, without replacing any
  existing path. If publication fails after a sidecar has been published,
  remove only sidecars proven to be the staged files from this invocation; never
  remove a pre-existing or competing file. Clean up temporary outputs on errors.
  Use one timestamp per backup invocation.

- [x] **Step 5: Add backup safety and dry-run tests**

  Cover explicit paths replacing YAML inputs, output-name collisions, target
  equal-to/inside-source rejection, existing-artifact conflicts, unsafe symlink
  warnings/manifest notes, unsupported types, and dry-run creating no files.
  Include a valid first source followed by an invalid/colliding source and assert
  preflight creates no artifacts. Inject a publication failure after the first
  sidecar and assert the published sidecar and temporary files are cleaned up,
  while competing/pre-existing files remain unchanged. Use the logger
  abstraction or a test logger to assert warning records. Assert backup dry-run
  exposes all planned artifact paths while creating no files.

- [x] **Step 6: Re-run backup tests**

  Run: `go test ./internal/tarball -run TestBackup -v`
  Expected: PASS for artifact, path, symlink, collision, and dry-run behavior.

- [x] **Step 7: Commit backup support**

  ```bash
  git add internal/tarball/service.go internal/tarball/backup.go \
    internal/tarball/backup_test.go \
    internal/tarball/test_helpers_test.go
  git commit -m "feat(tarball): create checksummed backups"
  ```

### Task 4: Verify and stage restores safely

**Files:**

- Create: `internal/tarball/restore.go`
- Create: `internal/tarball/extract.go`
- Test: `internal/tarball/restore_test.go`
- Modify: `internal/tarball/test_helpers_test.go`
- Reuse: `internal/tarball/manifest.go`, `paths.go`, and test fixture helpers

- [x] **Step 1: Write restore validation tests**

  Build archives and sidecars using Go fixtures. Cover valid archives, missing
  sidecars, malformed digest, checksum mismatch, unsupported manifest version,
  archive/manifest path/type/UID/GID mismatch, and symlink target or unsafe-flag
  mismatch. Assert every invalid case leaves destination untouched.

- [x] **Step 2: Run the focused tests and confirm they fail**

  Run: `go test ./internal/tarball -run TestRestoreValidation -v`
  Expected: FAIL because restore verification is not implemented.

- [x] **Step 3: Implement verified snapshot and metadata validation**

  Copy the archive into a private temporary snapshot while hashing it; compare
  the `.sha256` digest before gzip/tar parsing. Parse only that snapshot. Require
  a supported manifest and exact entry match for the original full-inventory
  format. Task 9 refines schema version 1 to also support directory/link-only
  manifests. Verify symlink target/safety from tar headers rather than trusting
  YAML; require optional UID/GID pairs to match headers and neutral `0/0` headers
  when omitted.

- [x] **Step 4: Add malicious archive fixtures and tests**

  Cover absolute member paths, `..` traversal, duplicate paths, hardlinks,
  devices/sockets, archive paths through symlinks, truncated archives, invalid
  gzip streams, and malformed tar headers. Assert no write outside the private
  staging directory. Unsafe links fail by default and pass validation only with
  the explicit option; archive path traversal never does.

- [x] **Step 5: Implement private staging extraction**

  Extract only validated entries into a private staging directory using
  descriptor-rooted/no-follow filesystem operations. Create implicit parent
  directories as needed; do not synthesize tar or manifest entries. Keep unsafe
  link creation separate from file placement so no member is ever written
  through a symlink.

- [x] **Step 6: Re-run restore validation tests**

  Run: `go test ./internal/tarball -run 'TestRestoreValidation|TestSafeExtract' -v`
  Expected: PASS; bad digest, metadata, or archive paths produce no destination
  writes or writes outside staging.

- [x] **Step 7: Commit restore validation and staging**

  ```bash
  git add internal/tarball/restore.go internal/tarball/extract.go \
    internal/tarball/restore_test.go internal/tarball/test_helpers_test.go
  git commit -m "feat(tarball): validate restore archives"
  ```

### Task 5: Apply restores with conflict and metadata policy

**Files:**

- Modify: `internal/tarball/service.go`
- Create: `internal/tarball/apply.go`
- Test: `internal/tarball/apply_test.go`
- Test: `internal/tarball/roundtrip_test.go`
- Modify: `internal/tarball/restore.go`

- [x] **Step 1: Write destination conflict tests**

  Cover no-conflict application, existing-file/link failure by default,
  `--overwrite` replacing only archive leaf files/links, unrelated destination
  preservation, directory type conflicts failing even with overwrite, implicit
  parent creation/reuse, and pre-existing symlink ancestors being rejected.

- [x] **Step 2: Run the focused tests and confirm they fail**

  Run: `go test ./internal/tarball -run TestApplyRestore -v`
  Expected: FAIL because destination application is not implemented.

- [x] **Step 3: Implement full preflight before destination writes**

  Check all archive entry destinations for conflicts before applying. Reuse
  existing real directories, create missing structural parents with `0755`
  subject to umask, reject file/symlink parents, and never replace directories.
  Canonicalize the chosen destination root and refuse to traverse descendant
  symlink components.

- [x] **Step 4: Implement apply and metadata handling**

  Apply staged files/links without deleting unrelated paths. Restore ordinary
  modes and modification times, stripping setuid/setgid. New entries use the
  current process UID/GID by default. `--preserve-owner` attempts archived
  ownership; missing IDs or denied changes warn and leave new entries owned by
  the current process. Apply metadata from explicit directory entries; implicit
  structural parents keep their safe default mode and current-user ownership.
  Use a narrow injectable ownership setter to test permission-denied fallback
  without requiring root.

- [x] **Step 5: Add overwrite, ownership, and dry-run tests**

  Verify `--dry-run` still hashes/validates but changes no destination files;
  check mode/mtime preservation, special-bit stripping, default ownership,
  `--preserve-owner` success/fallback, and explicit unsafe-link opt-in. Assert
  owner restoration is never attempted without `--preserve-owner`, and denied
  owner restoration on a reused explicit directory leaves its existing owner
  unchanged. Verify ACL/xattr metadata is not applied; make filesystem xattr
  assertions conditional when the test platform supports them. Add an end-to-end
  backup-to-restore test using Dum-generated sidecars for files, directories,
  safe symlinks, an in-home source, and an external source. Assert restore
  dry-run reports validated entries and destination paths without mutations.

- [x] **Step 6: Re-run restore application tests**

  Run: `go test ./internal/tarball -run 'TestApplyRestore|TestRestoreDryRun' -v`
  Expected: PASS for conflict, metadata, ownership, and dry-run behavior.

- [x] **Step 7: Commit restore application**

  ```bash
  git add internal/tarball/service.go internal/tarball/apply.go \
    internal/tarball/restore.go internal/tarball/apply_test.go \
    internal/tarball/roundtrip_test.go
  git commit -m "feat(tarball): apply safe restores"
  ```

### Task 6: Add list discovery and CLI commands

**Files:**

- Create: `internal/tarball/discovery.go`
- Test: `internal/tarball/discovery_test.go`
- Create: `cmd/dum/tarball.go`
- Create: `cmd/dum/tarball_picker.go`
- Test: `cmd/dum/tarball_test.go`
- Modify: `cmd/dum/root.go`
- Modify: `cmd/dum/root_test.go`
- Modify: `.golangci.yml`
- Modify: `AGENTS.md`

- [x] **Step 1: Write archive discovery tests**

  Cover zero/one/multiple `<basename>-*.tgz` matches, missing sidecars,
  case-insensitive basename collisions, and independent per-entry failure
  aggregation.

- [x] **Step 2: Run discovery tests and confirm they fail**

  Run: `go test ./internal/tarball -run TestDiscoverArchives -v`
  Expected: FAIL because discovery is not implemented.

- [x] **Step 3: Implement list-based discovery**

  Resolve each source basename, discover candidate archives in the selected
  target directory, and return typed candidates/errors without invoking
  interactive tools. Return per-entry results so one missing candidate does
  not stop discovery for later inputs.

- [x] **Step 4: Write Cobra command tests**

  Assert `dum tarball` registers `backup` and `restore`; explicit backup paths
  bypass `--input`, list mode requires `--input`, explicit archive restore
  bypasses input parsing, and relative archive paths use the current directory.
  Cover target/destination defaults, precedence (`--target` over
  `TGZ_TARGET` over XDG), relative `--input` resolution, flags, usage errors,
  dry-run wiring, and a list restore where an invalid first archive does not
  prevent a later valid archive from being restored; return nonzero for the
  failed item. Capture CLI output for backup/restore `--dry-run` and assert it
  names the planned artifacts or restore destinations.

- [x] **Step 5: Implement command and picker wiring**

  Add the parent/subcommands in `cmd/dum/tarball.go`, register the parent from
  `root.go`, and use local flag names so `--dry-run` is spelled as designed
  rather than the existing `--dryrun` helper. Resolve `TGZ_TARGET` and XDG
  defaults in the CLI layer. Add `tb` as the required import alias for
  `internal/tarball` in `.golangci.yml` and `AGENTS.md`. Keep YAML/archive
  business logic in `internal/tarball`.

  For multiple list matches, prefer `fzf`, then `gum`; if neither is installed,
  return all candidates in the error. Keep external-command execution in
  `cmd/dum/tarball_picker.go`, use context-aware subprocesses, and test with a
  controlled fake executable or an injected picker. Process list entries
  independently: continue after selection, checksum, manifest, conflict
  preflight, or apply errors; aggregate failures and return nonzero. Fully
  validate and preflight each archive before applying it.

- [x] **Step 6: Update root and command tests**

  Run: `go test ./cmd/dum -run 'TestNewDum_Tarball|TestTarballCommand' -v`
  Expected: PASS for registration, explicit/list modes, selection, and flags.

- [x] **Step 7: Commit CLI and discovery**

  ```bash
  git add internal/tarball/discovery.go internal/tarball/discovery_test.go \
    cmd/dum/tarball.go cmd/dum/tarball_picker.go cmd/dum/tarball_test.go \
    cmd/dum/root.go cmd/dum/root_test.go .golangci.yml AGENTS.md
  git commit -m "feat(tarball): add backup and restore CLI"
  ```

### Task 7: Document the command and run repository checks

**Files:**

- Modify: `README.md`

- [x] **Step 1: Document command usage and YAML input**

  Add `tarball` to the subcommand table. Include `input.yml` example, explicit
  backup path examples, archive restore and list restore examples, target/dest
  options, sidecar names, overwrite and unsafe-link/owner opt-ins, legacy archive
  rejection, and the distinction between SHA-256 integrity and authenticity.

- [x] **Step 2: Run focused and full tests**

  Run: `go test ./internal/tarball ./cmd/dum`
  Expected: PASS for all new package and CLI tests.

  Run: `mise run test`
  Expected: repository tests pass and coverage meets configured thresholds.

- [x] **Step 3: Run build and quality checks**

  Run: `mise run lint`
  Expected: no lint errors.

  Run: `mise run vet`
  Expected: no vet errors.

  Run: `mise run build`
  Expected: build succeeds for the configured local platform.

- [x] **Step 4: Check coverage and final diff**

  Run: `mise run check-coverage`
  Expected: repository coverage meets the 80% default threshold.

  Run: `git diff --check` and inspect `git status --short`.
  Expected: no whitespace errors; only intended feature files are modified.

- [x] **Step 5: Commit documentation**

  ```bash
  git add README.md
  git commit -m "docs: document tarball backup and restore"
  ```

  If coverage or tests fail, return to the owning task, fix and re-verify before
  committing the documentation.

### Task 8: Add backup exemptions

**Files:**

- Create: `internal/tarball/exemptions.go`
- Test: `internal/tarball/exemptions_test.go`
- Modify: `internal/tarball/service.go` — add `Exemptions []string` and
  `UseGitignore bool` to `BackupOptions`.
- Modify: `internal/tarball/manifest.go` — add `Exemptions []string` and
  `UseGitignore bool` top-level manifest fields with yaml tags.
- Modify: `internal/tarball/backup.go` — skip exempted entries in `writeArchive`;
  validate patterns in `planBackup`; record exemption config in manifest.
- Modify: `cmd/dum/tarball.go` — add `--exclude PATTERN` (repeatable) and
  `--use-gitignore` flags; build `Exemptions` in `runBackup`.
- Modify: `README.md` — document `--exclude` and `--use-gitignore` with
  `node_modules` / `dist` / `.gitignore` examples.

No new package dependencies: `path.Match` and `bufio`/`os` suffice. Do not add
`**` recursion in v1.

- [x] **Step 1: Write exemption matcher tests**

  Cover `path.Match` semantics: basename-at-any-depth for pattern without `/`,
  source-root anchoring for pattern with `/`, directory-only trailing slash,
  negation not supported (invalid leading `!` rejected), invalid glob rejection,
  `.gitignore` parsing (comments, blank lines, negation prefix, nested
  directory-relative resolution), per-source isolation, recursive loading,
  and rejection of symlinked `.gitignore` files.

- [x] **Step 2: Run the focused tests and confirm they fail**

  Run: `go test ./internal/tarball -run 'TestBackup_UsesNestedGitignoreFilesInScope|TestBackup_GitignoreRulesArePerSource|TestLoadGitignorePatternsDoesNotFollowSymlinks' -v`
  Expected: FAIL because nested rules are not loaded/scoped and source rules are
  currently merged across backups.

- [x] **Step 3: Implement exemption matcher**

  Implement `internal/tarball/exemptions.go`: compile each pattern with
  `path.Match`; scope patterns from each `.gitignore` to that directory; treat a
  trailing slash as directory-only and patterns without `/` as basename matches
  below their scope. Load a directory's `.gitignore` as the archive walk enters
  it, after ancestor exemptions have had a chance to prune it. Reject negation
  and symlinked `.gitignore` files.

- [x] **Step 4: Re-run the focused tests**

  Run: `go test ./internal/tarball -run 'TestBackup_UsesNestedGitignoreFilesInScope|TestBackup_GitignoreRulesArePerSource|TestLoadGitignorePatternsDoesNotFollowSymlinks|TestExemptionMatcher' -v`
  Expected: PASS for scoped matching, recursive `.gitignore`, and source isolation.

- [x] **Step 5: Wire exemptions into backup options and walk**

  Add `Exemptions []string` and `UseGitignore bool` to `BackupOptions`. In
  `planBackup`, validate every explicit pattern before walking sources. In
  `writeArchive`'s `filepath.WalkDir`, compute each entry's path relative to its
  source root; match CLI patterns and `.gitignore` rules scoped to the current
  directory. Load each reachable directory's `.gitignore` before descending; on
  a match, return `nil` for files or `filepath.SkipDir` for directories. Each
  source gets its own matcher and manifest rules. Record exemption config
  (patterns and `use_gitignore`) in the manifest. Exempted entries never reach
  the archive or manifest entries.

- [x] **Step 6: Add CLI flags**

  Add `--exclude PATTERN` (repeatable, no default) and `--use-gitignore` (bool,
  default false) to `dum tarball backup`. Read them in `runBackup` and pass
  explicit patterns plus the flag in `BackupOptions`; load source-specific
  root/nested `.gitignore` rules during each source archive walk.

- [x] **Step 7: Add backup exemption tests**

  Assert excluded `node_modules/`, `dist/`, and `*.log` files are absent from
  the archive; assert the manifest top-level `exemptions` list and explicit
  `use_gitignore: true` or `use_gitignore: false` are recorded correctly; assert
  an exempted directory is pruned (its nested files never visited); assert
  `--use-gitignore` honors root and nested `.gitignore` files, with nested rules
  scoped to their directory and no cross-source leakage; assert symlinked
  `.gitignore` files are rejected; assert invalid `--exclude` globs fail before
  artifacts are created; assert dry-run creates no files. Manifest inventory
  assertions are covered separately by Task 9: only directories and symlinks
  are listed.

- [x] **Step 8: Re-run backup tests**

  Run: `go test ./internal/tarball -run TestBackup -v`
  Expected: PASS including exemption pruning.

- [x] **Step 9: Run repository checks**

  Run: `go test ./internal/tarball ./cmd/dum`
  Run: `mise run test`
  Run: `mise run lint`
  Run: `mise run vet`
  Run: `mise run build`
  Run: `mise run check-coverage`

- [ ] **Step 10: Commit exemptions**

  ```bash
  git add internal/tarball/exemptions.go internal/tarball/exemptions_test.go \
    internal/tarball/service.go internal/tarball/backup.go \
    cmd/dum/tarball.go README.md
  git commit -m "feat(tarball): exempt files from backup"
  ```

### Task 9: List directories and symlinks in manifests

**Files:**

- Modify: `internal/tarball/backup.go` — omit regular-file entries from newly
  written manifests and mark unavailable regular-file owner metadata in the tar
  header.
- Modify: `internal/tarball/extract.go` — validate unlisted regular files while
  requiring manifest directory/link entries to match. Continue accepting
  pre-existing development manifests with complete file entries.
- Modify: `internal/tarball/apply.go` — restore regular-file ownership from tar
  headers when the manifest omits file entries; honor the unavailable-owner
  marker.
- Test: `internal/tarball/backup_test.go`, `internal/tarball/restore_test.go`,
  and `internal/tarball/roundtrip_test.go`.
- Modify: `README.md` and this design spec to document the reduced manifest.

Keep `schema_version: 1`; this format is not released. New backups list every
directory (including empty directories) and symlink, but no regular files.
Restore remains compatible with development archives whose v1 manifests list
every file.

- [x] **Step 1: Add a failing backup manifest test**

  Create a source containing a regular file, nested/empty directories, and a
  symlink. Assert the archive still contains all non-exempted entries while the
  manifest lists only directories and symlinks. Assert directory/link owner
  metadata remains paired, file owner IDs remain in tar headers, and
  `schema_version` remains 1. Inject unavailable file ownership and assert the
  tar header uses `0:0` plus the exact PAX marker
  `dum.owner-metadata=unavailable`; separately assert known `0:0` is not marked
  unavailable.

- [x] **Step 2: Run the focused backup test and confirm it fails**

  Run: `go test ./internal/tarball -run TestBackup_ManifestListsDirectoriesAndSymlinks -v`
  Expected: FAIL because backup currently lists regular files in the manifest.

- [x] **Step 3: Omit regular files from new manifests**

  In `writeArchive`, continue writing regular-file tar headers and content, but
  append `ManifestEntry` only for directories and symlinks. Preserve file
  UID/GID in tar headers. If `fileOwner` reports unavailable metadata, write
  neutral `0:0` header IDs and PAX record
  `dum.owner-metadata=unavailable` so restore can distinguish unavailable
  metadata from a root-owned file. Rejecting invalid markers is covered in the
  restore tests. Use the existing owner metadata seam to test this branch.

- [x] **Step 4: Re-run the focused backup test**

  Run: `go test ./internal/tarball -run TestBackup_ManifestListsDirectoriesAndSymlinks -v`
  Expected: PASS; tar members remain complete while manifest entries omit files.

- [x] **Step 5: Add failing restore tests for reduced manifests**

  Build Go-generated fixtures where a regular file exists in the tar but has no
  manifest entry, while all directories and symlinks do. Assert restore accepts
  the archive, still rejects unsafe paths/duplicates/type mismatches, and
  `--preserve-owner` uses tar-header IDs for unlisted files. Assert the owner
  unavailable marker triggers the current-user fallback; reject the marker on
  non-files, unknown marker values, and marked files whose IDs are not `0:0`.
  Reject unlisted directories and symlinks. Keep tests proving complete
  pre-release v1 file inventories still restore and mismatches still fail.

- [x] **Step 6: Run focused restore tests and confirm reduced-manifest cases fail**

  Run: `go test ./internal/tarball -run 'ReducedManifest|TarHeaderOwner|UnmarkedZeroOwner|UnavailableFileOwner|InvalidOwnerMetadataMarkers' -v`
  Expected: FAIL because every tar member is currently required in the manifest
  and file ownership currently comes from manifest entries.

- [x] **Step 7: Support directory/link inventory validation and header ownership**

  In archive inspection, require every manifest entry to match a tar member.
  If the manifest includes any regular-file entries, retain exact full-inventory
  matching for that pre-release format. If it has no regular-file entries,
  accept structurally valid tar files without manifest records while continuing
  all path, duplicate, type, symlink-parent, and directory/link checks. During
  apply, use manifest ownership for listed entries and tar-header ownership for
  unlisted regular files unless the exact PAX marker
  `dum.owner-metadata=unavailable` says metadata is unavailable. Require this
  marker's IDs to be neutral `0:0` and reject invalid marker placement/values.
  Populate `RestoreResult.Paths` from all staged archive members so dry-run
  output still reports regular files omitted from the manifest.

- [x] **Step 8: Re-run focused backup, restore, and round-trip tests**

  Run: `go test ./internal/tarball -run 'Backup_Manifest|RegularFileOwnerMetadataMarker|ReducedManifest|TarHeaderOwner|UnmarkedZeroOwner|UnavailableFileOwner|InvalidOwnerMetadataMarkers|Restore_RoundTripsBackupArchives' -v`
  Expected: PASS for compact manifests, full archive validation, owner restore,
  and pre-release compatibility.

- [x] **Step 9: Update documentation and run repository checks**

  Document that the manifest lists directories and links, not regular files;
  clarify that restore validates regular-file paths from the tar and SHA-256
  covers archive bytes. Document `dum.owner-metadata=unavailable` behavior.
  Run:

  ```bash
  go test ./internal/tarball ./cmd/dum
  mise run test
  mise run lint
  mise run vet
  mise run build
  mise run check-coverage
  ```

- [ ] **Step 10: Commit the reduced manifest change**

  ```bash
  git add internal/tarball/backup.go internal/tarball/extract.go \
    internal/tarball/apply.go internal/tarball/backup_test.go \
    internal/tarball/restore_test.go internal/tarball/roundtrip_test.go \
    README.md docs/superpowers/specs/2026-11-03-tarball-command-design.md \
    docs/superpowers/plans/2026-11-03-tarball-command.md
  git commit -m "feat(tarball): omit files from manifests"
  ```
