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
    git reset --hard HEAD
    for directory in vim/vim.d/plugged/*/; do
        if [[ -d "$directory/.git" ]]; then
            echo "Resetting $directory ..."
            git -C "$directory" reset --hard HEAD
            git -C "$directory" clean -fd
        fi
    done
}

cmd_install() {
    announce "Installing tools ..."
    bash "$TOOLS_DIR/dev-prep.sh"
}

cmd_check_tools() {
    announce "Checking tools ..."
    bash "$TOOLS_DIR/dev-env-check.sh"
}

cmd_ci_setup() {
    announce "Installing CI tools ..."
    go install github.com/abice/go-enum@latest
    go install github.com/vektra/mockery/v3@latest
    go install golang.org/x/tools/cmd/goimports@latest
}

cmd_clean() {
    announce "Cleaning up ..."
    go clean -cache
    rm -rf "$OUT_DIR"
    rm -f ./internal/**/*_mocks.go
    rm -f ./internal/**/*_enum.go
}

cmd_tidy() {
    announce "Tidying ..."
    go mod tidy
}

cmd_vendor() {
    announce "Vendoring ..."
    go mod vendor
}

cmd_generate() {
    announce "Generating ..."
    go generate ./...
    for config in cfg/mockery*.yml; do
        mockery --config "$config"
    done
}

cmd_build() {
    announce "Building ..."
    go build -v \
        -ldflags "-X main.version=$DUM_VERSION -X main.commit=$GIT_COMMIT" \
        -o "$DUM_EXECUTABLE-$DUM_VERSION" "$DUM_SOURCE"
}

cmd_linters_build() {
    announce "Building linters ..."
    mkdir -p "$OUT_DIR"
    (cd cmd/linters && go build -o ../dist/linters .)
}

cmd_release() {
    announce "Compiling releases for every OS and Platform ..."
    goreleaser build --clean --snapshot --config cfg/goreleaser.local.yaml
}

cmd_release_build() {
    announce "Compiling builds for every OS and Platform ..."
    goreleaser build --clean --snapshot --config cfg/goreleaser.snapshot.yaml
}

cmd_check() {
    announce "Checking code quality ..."
}

cmd_lint() {
    announce "Linting ..."
    golangci-lint run cmd/...
    golangci-lint run internal/...
}

cmd_vet() {
    announce "Vetting ..."
    go vet ./...
}

cmd_fmt() {
    announce "Formatting ..."
    golangci-lint fmt
}

cmd_fix() {
    announce "Fixing ..."
    go fix ./...
}

cmd_test() {
    announce "Testing ..."
    go test -v -cover \
        -coverprofile="$OUT_DIR/cover.out" \
        -covermode=atomic \
        -coverpkg=./... \
        -race ./...
}

cmd_coverage() {
    announce "Checking coverage of tests ..."
    go-test-coverage --config=./cfg/testcoverage.yml || true
    announce "Generating HTML report ..."
    go tool cover -html="$OUT_DIR/cover.out" -o "$OUT_DIR/coverage.html"
    announce "HTML report: $OUT_DIR/coverage.html"
}

cmd_check_coverage() {
    announce "Checking coverage threshold ..."
    go-test-coverage --config=./cfg/testcoverage.yml
}

cmd_profile() {
    announce "Profiling ..."
    go test -cpuprofile cpu.prof -memprofile mem.prof -v -bench ./...
}

cmd_bench() {
    announce "Benchmarking ..."
    go test -v -benchmem -bench=. ./...
}

cmd_all() {
    announce "running all ..."
}

cmd_check_lfs_hook() {
    if ! grep -q '^git lfs pre-push' "$REPO_ROOT/.git/hooks/pre-push"; then
        echo "git-lfs is not installed properly. Please run 'git lfs install' and try again." >&2
        return 1
    fi
}

cmd_coverage_check() {
    bash "$TOOLS_DIR/coverage-check.sh"
}

cmd_cpd() {
    bash "$TOOLS_DIR/cpd.sh"
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
