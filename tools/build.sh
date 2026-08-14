#!/usr/bin/env bash
set -euo pipefail

# Values such as OUT_DIR are initialized by the sourced common library.
# shellcheck disable=SC2153

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd -P)"

if [[ "$(pwd -P)" != "$ROOT_DIR" ]]; then
    echo "tools/build.sh must be run from the repository root: $ROOT_DIR" >&2
    exit 1
fi

for library in "$SCRIPT_DIR"/lib/*.sh; do
    # shellcheck source=/dev/null
    source "$library"
done

case "${1:-}" in
print-go-version)
    echo "Go Version: $(go version)"
    ;;
make_out)
    # shellcheck disable=SC2153
    mkdir -p "$OUT_DIR"
    ;;
rhard)
    announce "git resetting hard ..."
    repository_rhard
    ;;
install)
    announce "Installing tools ..."
    setup_install_tools
    ;;
check_tools)
    announce "Checking tools ..."
    setup_check_tools
    ;;
ci-setup)
    announce "Installing CI tools ..."
    setup_ci_tools
    ;;
clean)
    announce "Cleaning up ..."
    repository_clean
    ;;
tidy)
    announce "Tidying ..."
    go_tidy
    ;;
vendor)
    announce "Vendoring ..."
    go_vendor
    ;;
generate)
    announce "Generating ..."
    go_generate
    ;;
build)
    announce "Building ..."
    go_build
    ;;
linters-build)
    announce "Building linters ..."
    go_linters_build
    ;;
release)
    announce "Compiling releases for every OS and Platform ..."
    go_release
    ;;
release-build)
    announce "Compiling builds for every OS and Platform ..."
    go_release_build
    ;;
check)
    announce "Checking code quality ..."
    ;;
lint)
    announce "Linting ..."
    go_lint
    ;;
vet)
    announce "Vetting ..."
    go_vet
    ;;
fmt)
    announce "Formatting ..."
    go_fmt
    ;;
fix)
    announce "Fixing ..."
    go_fix
    ;;
test)
    announce "Testing ..."
    quality_test
    ;;
coverage)
    announce "Checking coverage of tests ..."
    quality_coverage
    ;;
check-coverage)
    announce "Checking coverage threshold ..."
    quality_check_coverage
    ;;
profile)
    announce "Profiling ..."
    quality_profile
    ;;
bench)
    announce "Benchmarking ..."
    quality_bench
    ;;
all)
    announce "running all ..."
    ;;
check-lfs-hook)
    repository_check_lfs_hook
    ;;
coverage-check)
    quality_coverage_check
    ;;
cpd)
    quality_cpd
    ;;
help)
    print_dispatcher_usage
    ;;
"")
    print_dispatcher_usage >&2
    exit 2
    ;;
*)
    echo "Unknown build subcommand: $1" >&2
    print_dispatcher_usage >&2
    exit 2
    ;;
esac