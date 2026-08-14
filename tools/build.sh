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

cmd_print_go_version() {
    echo "Go Version: $(go version)"
}

cmd_make_out() {
    # shellcheck disable=SC2153
    mkdir -p "$OUT_DIR"
}

cmd_rhard() {
    announce "git resetting hard ..."
    repository_rhard
}

cmd_install() {
    announce "Installing tools ..."
    setup_install_tools
}

cmd_check_tools() {
    announce "Checking tools ..."
    setup_check_tools
}

cmd_ci_setup() {
    announce "Installing CI tools ..."
    setup_ci_tools
}

cmd_clean() {
    announce "Cleaning up ..."
    repository_clean
}

cmd_tidy() {
    announce "Tidying ..."
    go_tidy
}

cmd_vendor() {
    announce "Vendoring ..."
    go_vendor
}

cmd_generate() {
    announce "Generating ..."
    go_generate
}

cmd_build() {
    announce "Building ..."
    go_build
}

cmd_linters_build() {
    announce "Building linters ..."
    go_linters_build
}

cmd_release() {
    announce "Compiling releases for every OS and Platform ..."
    go_release
}

cmd_release_build() {
    announce "Compiling builds for every OS and Platform ..."
    go_release_build
}

cmd_check() {
    announce "Checking code quality ..."
}

cmd_lint() {
    announce "Linting ..."
    go_lint
}

cmd_vet() {
    announce "Vetting ..."
    go_vet
}

cmd_fmt() {
    announce "Formatting ..."
    go_fmt
}

cmd_fix() {
    announce "Fixing ..."
    go_fix
}

cmd_test() {
    announce "Testing ..."
    quality_test
}

cmd_coverage() {
    announce "Checking coverage of tests ..."
    quality_coverage
}

cmd_check_coverage() {
    announce "Checking coverage threshold ..."
    quality_check_coverage
}

cmd_profile() {
    announce "Profiling ..."
    quality_profile
}

cmd_bench() {
    announce "Benchmarking ..."
    quality_bench
}

cmd_all() {
    announce "running all ..."
}

cmd_check_lfs_hook() {
    repository_check_lfs_hook
}

cmd_coverage_check() {
    quality_coverage_check
}

cmd_cpd() {
    quality_cpd
}


case "${1:-}" in
    print-go-version) cmd_print_go_version ;;
    make_out) cmd_make_out ;;
    rhard) cmd_rhard ;;
    install) cmd_install ;;
    check_tools) cmd_check_tools ;;
    ci-setup) cmd_ci_setup ;;
    clean) cmd_clean ;;
    tidy) cmd_tidy ;;
    vendor) cmd_vendor ;;
    generate) cmd_generate ;;
    build) cmd_build ;;
    linters-build) cmd_linters_build ;;
    release) cmd_release ;;
    release-build) cmd_release_build ;;
    check) cmd_check ;;
    lint) cmd_lint ;;
    vet) cmd_vet ;;
    fmt) cmd_fmt ;;
    fix) cmd_fix ;;
    test) cmd_test ;;
    coverage) cmd_coverage ;;
    check-coverage) cmd_check_coverage ;;
    profile) cmd_profile ;;
    bench) cmd_bench ;;
    all) cmd_all ;;
    check-lfs-hook) cmd_check_lfs_hook ;;
    coverage-check) cmd_coverage_check ;;
    cpd) cmd_cpd ;;
    help) print_dispatcher_usage ;;
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
