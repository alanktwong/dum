# `dum.yml` Default Configuration Design

Date: 2026-08-20
Status: Approved

## Goal

Make `dum.yml` the default configuration file name while keeping existing
`installer.yml` setups working through a legacy fallback.

## Resolution Order

`GetDefaultConfig()` in `cmd/dum/cli.go` resolves the default config path.
First match wins:

1. `DUM_CONFIG` environment variable — returned verbatim, no existence check
   (replaces `INSTALLER_CONFIG`).
2. `./dum.yml` — must exist (`os.Stat`).
3. `$XDG_CONFIG_HOME/dum/dum.yml` — must exist (`XDG_CONFIG_HOME` defaults to
   `~/.config`).
4. Legacy fallback: `$XDG_CONFIG_HOME/dum/installer.yml` — must exist.
5. Nothing found — return `$XDG_CONFIG_HOME/dum/dum.yml` as the canonical
   default so help text and downstream errors point at the expected location.

## Unchanged Behavior

- `--file` flag overrides every resolution step.
- `~` remains literal in returned paths; expansion happens downstream in the
  factory/loader as it does today.
- Resolution still happens once at flag construction via an isolated Viper
  instance bound to env vars.

## Documentation Updates

- `cmd/dum/install.go`: Config section paths and the "Configuring the
  installer.yml:" heading become `dum.yml`; override line references
  `DUM_CONFIG`.
- `cmd/dum/schema.go`: short/long text references `dum.yml`.
- `README.md`: configuration section paths and env var name.
- `AGENTS.md`: overview line referencing `installer.yml`.
- Internal `cfg/installer.schema.json` filename stays as-is (out of scope).

## Testing

- Rewrite `flags_test.go` `GetDefaultConfig` coverage: env precedence, cwd
  hit, XDG hit, legacy fallback, miss returns canonical default, repeated-call
  environment isolation.
- Update `install_test.go` and any root tests referencing
  `INSTALLER_CONFIG`.
- Fixture-based cases use `t.TempDir()` with `t.Setenv("XDG_CONFIG_HOME")`
  and `t.Chdir` (or equivalent) for cwd search.

## Error Handling

No new error paths. Missing files fail downstream in factory loading with
existing clear messages.
