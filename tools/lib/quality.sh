#!/usr/bin/env bash

# Test, coverage, profiling, benchmarking, and duplication recipes.

quality_test() {
    go test -v -cover \
        -coverprofile="$OUT_DIR/cover.out" \
        -covermode=atomic \
        -coverpkg=./... \
        -race ./...
}

quality_write_coverage_config() {
    local config_file="$1"
    local profile="${OUT_DIR%/}/cover.out"

    awk -v profile="$profile" -v threshold="$COVERAGE_THRESHOLD" '
        /^profile:/ {
            print "profile: " profile
            next
        }
        /^  total:/ {
            print "  total: " threshold
            next
        }
        { print }
    ' "$REPO_ROOT/cfg/testcoverage.yml" > "$config_file"
}

quality_run_coverage_checker() (
    local config_file
    config_file="$(mktemp "${TMPDIR:-/tmp}/dum-testcoverage.XXXXXX")"
    trap 'rm -f -- "$config_file"' EXIT

    quality_write_coverage_config "$config_file"
    go-test-coverage --config="$config_file"
)

quality_coverage() {
    quality_run_coverage_checker || true
    go tool cover -html="$OUT_DIR/cover.out" -o "$OUT_DIR/coverage.html"
    announce "HTML report: $OUT_DIR/coverage.html"
}

quality_check_coverage() {
    quality_run_coverage_checker
}

quality_profile() {
    mkdir -p "$OUT_DIR"
    go test -cpuprofile="$OUT_DIR/cpu.prof" -memprofile="$OUT_DIR/mem.prof" -v -bench ./...
}

quality_bench() {
    go test -v -benchmem -bench=. ./...
}

quality_coverage_check() (
    local coverage_profile filtered_profile coverage

    if [[ -f "$OUT_DIR/coverage.txt" ]]; then
        coverage_profile="$OUT_DIR/coverage.txt"
    elif [[ -f "$OUT_DIR/cover.out" ]]; then
        coverage_profile="$OUT_DIR/cover.out"
    else
        echo "Coverage profile not found in $OUT_DIR." >&2
        return 1
    fi

    filtered_profile="$OUT_DIR/coverage-filtered.txt"
    grep -v '_enum.go' "$coverage_profile" > "$filtered_profile" || {
        local grep_status=$?
        if ((grep_status != 1)); then
            return "$grep_status"
        fi
    }

    coverage="$(go tool cover -func="$filtered_profile" | awk '/total:/ { gsub("%", "", $3); print $3 }')"
    if [[ -z "$coverage" ]]; then
        echo "Unable to determine coverage from $filtered_profile." >&2
        return 1
    fi

    if awk -v coverage="$coverage" -v threshold="$COVERAGE_THRESHOLD" \
        'BEGIN { exit !(coverage < threshold) }'; then
        echo "Coverage ${coverage}% is below threshold ${COVERAGE_THRESHOLD}%"
        return 1
    fi

    echo "Coverage ${coverage}% meets threshold ${COVERAGE_THRESHOLD}%"
)

quality_cpd() (
    local tmp found_dir
    local -a source_dirs=()

    tmp="$(mktemp "${TMPDIR:-/tmp}/cpd.XXXXXX")"
    trap 'rm -f -- "$tmp"' EXIT

    while IFS= read -r found_dir; do
        if [[ -n "$found_dir" ]]; then
            source_dirs+=("$found_dir")
        fi
    done < <(find backend core frontend -name java -type d 2>/dev/null | grep '/src/main/java' || true)

    if ((${#source_dirs[@]} > 0)); then
        find "${source_dirs[@]}" -name '*.java' > "$tmp"
    else
        : > "$tmp"
    fi

    pmd cpd --debug --file-list "$tmp" --format text --minimum-tokens 100 | tee /tmp/cpd.txt
)
