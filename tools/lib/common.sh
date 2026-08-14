#!/usr/bin/env bash

# Shared paths and build metadata for the repository-root dispatcher.

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
readonly REPO_ROOT
TOOLS_DIR="$REPO_ROOT/tools"
# shellcheck disable=SC2034
readonly TOOLS_DIR

DUM_APP="${DUM_APP:-dum}"
# shellcheck disable=SC2034
DUM_SOURCE="${DUM_SOURCE:-$REPO_ROOT/cmd/$DUM_APP}"
OUT_DIR="${OUT_DIR:-./dist}"

GIT_COMMIT="${GIT_COMMIT:-$(git -C "$REPO_ROOT" rev-parse --short HEAD 2>/dev/null || true)}"

if [[ -z "${DUM_VERSION:-}" ]]; then
    GIT_TAG="$(git -C "$REPO_ROOT" describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || true)"
    GIT_DIRTY="$(git -C "$REPO_ROOT" diff --quiet HEAD 2>/dev/null || printf '%s' '-dirty')"
    DUM_VERSION="${GIT_TAG:+$GIT_TAG-}${GIT_COMMIT}${GIT_DIRTY}"
fi

# shellcheck disable=SC2034
DUM_EXECUTABLE="${DUM_EXECUTABLE:-$OUT_DIR/$DUM_APP}"

announce() {
    printf '\033[1;36m%s\033[0m\n' "$1"
}

print_dispatcher_usage() {
    cat <<'EOF'
Usage: ./tools/build.sh <subcommand>

Subcommands:
  Shared:
    print-go-version  print the Go version
    make_out          create the output directory
  Setup:
    install           install development tools
    check_tools       check required development tools
    ci-setup          install tools needed by CI
  Repository:
    rhard             reset the repository
    clean             remove generated output
    check-lfs-hook    verify the Git LFS pre-push hook
  Go:
    tidy              tidy Go dependencies
    vendor            vendor Go dependencies
    vet               run go vet
    fmt               format Go code
    fix               run go fix
    generate          generate Go code and mocks
    lint              run linters
    build             build the dum executable
    linters-build     build custom linters
    release           compile release binaries
    release-build     compile snapshot binaries
  Quality:
    check             run the code quality aggregate recipe
    test              run tests with coverage
    coverage          generate a coverage report
    check-coverage    check the coverage threshold
    profile           collect CPU and memory profiles
    bench             run benchmarks
    coverage-check    run the standalone coverage checker
  All:
    all               run the aggregate recipe
EOF
}
