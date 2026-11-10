# Design: `dum tarball` command

Date: 2026-11-03
Status: Approved (user)

## Goal

Add a `dum tarball` command with `backup` and `restore` subcommands. Each backup
produces a gzip-compressed tar archive, a SHA-256 sidecar, and a versioned YAML
inventory. Restore verifies and validates the archive before writing anything
to the requested destination.

The design builds on the list-driven behavior in
`~/projects/dotfiles/bin/tarball.sh`, while also allowing explicit backup paths.

## Decisions

- Use a plain `.tgz` plus two external sidecars; do not create a custom bundle.
- Use Go archive handling and validation, not a platform `tar` subprocess.
- Support explicit paths and a YAML-only input list for backup/restore workflows.
- Require `--input INPUT.yml` for list mode; do not use a plain-text list,
  implicit default, environment fallback, or automatic seeding.
- Use a simple top-level `paths` array; relative entries resolve from the
  directory containing the input YAML file.
- Produce one archive per source path.
- Preserve numeric UID/GID metadata when available. Restore uses the current
  process owner by default; applying archived ownership requires
  `--preserve-owner`.
- Require sidecars for restore. Archives created by the old shell script, which
  have no sidecars, are not accepted by this command.
- Keep SHA-256 as an integrity check, not an authenticity guarantee. This
  design does not add signatures or keys.
- Do not add source deletion, per-file content hashes, selective file restore,
  or restore-to-original-absolute-path behavior in v1.

## CLI and project boundaries

Register `dum tarball` in the root Cobra command. Cobra flags, arguments, and
subcommand wiring stay under `cmd/dum`; archive operations, manifest types,
checksum handling, and safe path logic live in a Cobra-free
`internal/tarball` package. The project already depends on `gopkg.in/yaml.v3`.

Proposed interface:

```text
dum tarball backup [PATH...] [--input INPUT.yml] [--target DIR] [--exclude PATTERN]
                   [--use-gitignore] [--dry-run]
dum tarball restore [ARCHIVE] [--input INPUT.yml] [--target DIR] [--dest DIR]
                    [--overwrite] [--allow-unsafe-links] [--preserve-owner]
                    [--dry-run]
```

`backup` paths replace the YAML list. With no paths, `--input INPUT.yml` is
required. Explicit relative paths resolve from the current working directory;
relative entries in the YAML file resolve from the YAML file's directory.
Both forms may select paths outside `$HOME`. Explicit paths are the complete
selection; `--input` is read only when there are no explicit paths.

The only supported list format is:

```yaml
paths:
  - Documents
  - /mnt/data/Photos
```

Accept normal YAML path strings, including spaces. Reject an empty/missing
`paths` list, empty entries, and duplicate resolved inputs. In list mode,
`--input INPUT.yml` is required; do not infer a filename or read an environment
variable. Relative `--input` filenames resolve from the current directory.

An explicit `restore ARCHIVE` uses that archive's sidecars and does not read the
input YAML. Relative archive paths resolve from the current directory; `--target`
is used only for list-based archive discovery. Bare `restore` requires
`--input INPUT.yml`; its `paths` entries drive archive discovery. One matching
archive is selected automatically, multiple matches use the available
`fzf`/`gum` picker, and multiple matches without a picker report the candidates.
Missing matches are reported and aggregated. `--input` is read only in this
list-based mode. For each path, discover `<target>/<basename>-*.tgz`; reject
case-insensitive basename collisions across input entries before selection.

Keep the current archive-target configuration names and precedence:

| Setting | Resolution |
| --- | --- |
| Archive directory | `--target` > `TGZ_TARGET` > `${XDG_DATA_HOME:-$HOME/.local/share}/tarballs` |
| Restore destination | `--dest`, otherwise `$HOME` |

At runtime, resolve the archive directory with this precedence. In Cobra help,
show the fallback symbolically as `TGZ_TARGET, else $XDG_DATA_HOME/tarballs,
else ~/.local/share/tarballs` for both backup and restore. Do not expand
environment variables or `~` to user-specific absolute paths in help. The help
display is separate from the runtime `--target` value, which continues to use
`defaultTarballTarget()`.

List mode fails clearly if `--input` is omitted or the YAML file is missing or
invalid. Explicit paths and explicit archive restore do not require an input
file. Do not support the old `tarballs.txt` format or implicitly seed a list.

## Archive and sidecar layout

For each selected source, create a timestamped archive using one timestamp per
backup invocation. `<name>` is the source's final path component; reject
case-insensitive basename collisions before writing any output:

```text
<name>-<YYYY-MM-DD-HHMMSS>.tgz
<name>-<YYYY-MM-DD-HHMMSS>.tgz.sha256
<name>-<YYYY-MM-DD-HHMMSS>.tgz.manifest.yaml
```

The `.sha256` file contains one lowercase 64-character hexadecimal digest and
a newline. The manifest is versioned YAML. Its inventory records each
archive-relative path and its type (`file`, `directory`, or `symlink`).
Symlink entries also record their literal target and, when unsafe, a warning
and unsafe marker. Entries may include numeric `uid` and `gid` when source
ownership metadata is available; both fields must be present together and match
the tar headers. If ownership metadata is unavailable or cleared, omit both
fields and write neutral owner IDs (`uid: 0`, `gid: 0`) to the tar header. Do
not include per-file hashes or the original absolute source path.

The manifest also records backup configuration at the top level: `exemptions`
lists the glob patterns used to exclude files from the archive (omitted when
none were configured), and `use_gitignore` is a boolean indicating whether a
source-root `.gitignore` was consulted. These fields describe why certain
entries are absent from the inventory.

Tar headers carry permissions and modification times. Restore preserves
ordinary permission bits and modification times, but strips setuid and setgid
bits. Do not restore ACLs or extended attributes.

The archive member root preserves the current layout for sources inside `$HOME`
(the path relative to `$HOME`). For a source outside `$HOME`, use only its
basename as the member root. No original absolute path is used as an extraction
destination. Detect path/name collisions before writing; never overwrite any
existing archive or sidecar.

Example manifest shape:

```yaml
schema_version: 1
entries:
  - path: Documents
    type: directory
  - path: Documents/project/readme.txt
    type: file
    uid: 501
    gid: 20
  - path: Documents/config-link
    type: symlink
    target: /old/location
    unsafe: true
    warning: absolute symlink target; restore requires --allow-unsafe-links
    uid: 501
    gid: 20
```

Paths in `entries` must be normalized, relative archive paths. The manifest
describes the complete archive inventory, including directory entries. Do not
silently omit unsupported filesystem entries such as devices or sockets; fail
the backup with a clear error.

## Backup behavior

1. Resolve inputs from explicit arguments or the required YAML `paths` list.
   Explicit relative paths use the working directory; relative YAML entries use
   the YAML file's directory. Validate all selected paths and output names
   before creating output.
2. Reject a backup target equal to or inside any source, preventing self-inclusion
   and inclusion of earlier backups.
3. Walk each source without following symlinks. Preserve ordinary relative links
   only when their targets remain inside the archived root.
4. For absolute or root-escaping symlinks, log a warning, preserve the symlink
   entry without dereferencing it, and record its literal target and unsafe
   status in the manifest.
5. Build each archive and its sidecars in temporary files. Calculate SHA-256
   from the completed archive bytes.
6. Publish sidecars first and the archive last. The archive is the visible
   completion marker; restore requires the full trio. Clean up temporary files
   on failure and leave any pre-existing output untouched.
7. Do not remove source paths. There is no `--delete-source` in v1.

## Exemptions

Exclude files and directories from a backup before they enter the archive.
Exemptions apply per source, match against the entry's path relative to the
source root, and drop matched directories without descending into them.

1. `--exclude PATTERN` (repeatable): one glob pattern per flag. Patterns use
   `path.Match` syntax (`*`, `?`, `[...]`). A pattern without `/` matches any
   path segment by basename at any depth (`node_modules`, `dist`, `*.log`).
   A pattern with `/` anchors from the source root (`build/output`). A trailing
   slash matches directories only (`dist/`). Invalid globs are rejected before
   walking sources.
2. `--use-gitignore`: load `.gitignore` from each source root and merge its
   patterns into the same exemption set. Patterns in a nested `.gitignore`
   apply relative to that file's own directory.
3. Exempted entries are skipped in the walk and never appear in the archive,
   manifest entries, checksum, or sidecars. Directories matched by an exemption
   are pruned with `SkipDir`, so their contents are never visited.
4. The manifest records the exemption configuration used: the list of glob
   patterns under `exemptions` and whether `.gitignore` was consulted under
   `use_gitignore`. When no exemptions were configured, `exemptions` is omitted
   and `use_gitignore` is `false`.
5. Exemptions are configured only on `backup`; restore has no exemption concept
   and restores exactly what the archive contains.

`--dry-run` validates inputs and prints the planned outputs without creating or
changing files. Exemption pruning is part of the walk, so `--dry-run` plans
assume exemptions are not materialized until an actual backup runs.

## Restore and safety behavior

1. Resolve the archive and require both sidecars. Do not restore legacy archives
   that lack either sidecar.
2. Copy the archive into a private temporary snapshot while computing SHA-256.
   Compare the digest before decompressing or parsing the archive. Parse only the
   verified snapshot to avoid a hash-then-reopen race.
3. Parse a supported manifest version and validate every tar entry. Require an
   exact inventory/type match. Optional UID/GID fields must appear as a pair and
   match the tar header; if omitted, the header must contain neutral owner IDs.
   For symlinks, require the manifest's literal target to match the tar header
   and derive unsafe status independently from that header; reject mismatched
   or falsely-safe manifest metadata. Reject absolute member paths, parent
   traversal, duplicate paths, unsupported hard links/special entries, and any
   member path that traverses a symlink.
4. Validate symlink targets. Safe relative links that remain inside the archive
   root may be restored. An absolute or escaping link causes restore to fail
   before destination writes unless `--allow-unsafe-links` is set. Even with the
   flag, create the link itself; never follow it during extraction or file
   placement. The flag does not permit archive path traversal.
5. Extract and validate the full archive in a private staging directory. Do not
   write into the destination during checksum, manifest, or archive validation.
6. Preflight destination conflicts before applying. Create missing structural
   parent directories with mode `0755` subject to umask and current-user
   ownership; reuse existing real directories even when they have no explicit
   tar entry. Reject file or symlink ancestors. Existing files or leaf symlinks
   cause failure unless `--overwrite` is set. With `--overwrite`, replace only
   conflicting files or leaf symlinks; never delete unrelated paths or replace
   directories. Any file/directory type conflict fails even with `--overwrite`.
   Refuse to traverse pre-existing symlink components below the chosen
   destination root.
7. Apply staged entries under `--dest`, defaulting to `$HOME`. In-home source
   layout is preserved relative to `$HOME`; external sources land under
   `<dest>/<basename>`. Preserve ordinary file/directory modes and modification
   times while clearing setuid/setgid bits. By default, new entries are owned by
   the current process UID/GID. With `--preserve-owner`, attempt to apply the
   manifest UID/GID; if values are absent or ownership changes are denied, warn
   and leave new entries owned by the current process user. Apply metadata from
   explicit directory entries; implicit structural parents keep their safe
   default mode and current-user ownership. Reused directories retain existing
   ownership when a change is denied. This flag can assign files to other local
   accounts when run as root, so it is explicit opt-in.

`--dry-run` performs checksum and archive/manifest validation and reports the
planned restore without changing the destination. Validation and conflict
failures happen before destination writes. A filesystem error during final
application may leave a partial restore; report the failing path and error.

SHA-256 protects against accidental or uncoordinated byte changes only when the
expected digest is trusted. It does not authenticate the creator: someone able
to replace both archive and sidecars can generate a matching digest.

## Failure handling

- Missing `--input` in list mode, missing input YAML, source, archive, checksum,
  or manifest: fail with a specific error; do not silently seed or fall back.
- Invalid digest, malformed/unsupported manifest, or archive/manifest mismatch:
  fail before destination writes.
- Unsafe archive paths, duplicate entries, unsupported types, or unsafe
  symlinks without explicit opt-in: fail before destination writes.
- Existing output artifacts: fail without overwriting them.
- Restore conflicts without `--overwrite`: report conflicts and write nothing.
- Resolve each `paths` entry to its basename for archive matching. Reject
  case-insensitive basename collisions in the input before list restore.
- Process list restore entries independently: fully verify, stage, and preflight
  each archive before applying that archive; continue to later entries after an
  item fails and return nonzero if any item failed. Earlier successful items
  remain restored if a later item fails.

## Testing

Use Go tests with `t.TempDir` and Go-generated tar/gzip fixtures; tests must not
depend on the host `tar` or checksum commands.

Cover:

- CLI registration, flags, explicit-path precedence, symbolic `--target` help
  for backup and restore without changing runtime environment precedence, YAML
  parsing and path resolution, required `--input`, and invalid/missing input
  behavior.
- Round trips for files, directories, symlinks, in-home layout, and external
  basename layout.
- Complete artifact creation, naming collisions, target/source overlap, and
  cleanup after failed backup publication.
- Digest mismatch, missing sidecars, malformed/unsupported manifest, and
  archive/manifest inventory mismatch.
- Traversal, absolute archive members, duplicate paths, unsupported entry
  types, archive symlink-parent attacks, and pre-existing destination symlink
  ancestors cannot escape staging or destination.
- Absolute/escaping symlink warnings and manifest records; default restore
  rejection and explicit `--allow-unsafe-links` behavior.
- Optional UID/GID manifest fields match tar headers; default current-user
  ownership and explicit `--preserve-owner` success/fallback behavior.
- Destination conflict behavior with and without `--overwrite`; unrelated files
  remain untouched; directory type conflicts fail even with overwrite.
- File/directory mode and modification-time preservation, setuid/setgid-bit
  stripping, and ACL/xattr non-restoration.
- `--dry-run` validates and reports without changing output or destination.
- Nested in-home file restore creates/reuses safe implicit parent directories;
  missing parents are not synthesized as tar members or manifest entries.
- Nested in-home file restore is tested with both missing and existing parent
  directories.
- List restore zero/one/multiple archive selection, relative archive paths,
  basename-collision rejection, and picker-unavailable diagnostics.
- Backup exemptions: `--exclude` repeated patterns skip matching files and
  directories at any depth; `--use-gitignore` honors a source-root `.gitignore`;
  nested `.gitignore` patterns apply relative to their own directory; exempted
  directories are pruned and never descended into; manifest, checksum, and
  sidecars contain no exempted entries.

Run the focused package tests, then the repository's normal `mise run test`,
`mise run lint`, and `mise run vet` checks during implementation.
