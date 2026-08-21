#!/usr/bin/env bash

# Go build, quality, and release recipes for the repository dispatcher.

go_tidy() {
    go mod tidy
}

go_vendor() {
    go mod vendor
}

go_generate() {
    go generate ./...
    # The logging config must run first: it generates internal/logging/gen,
    # which later configs need to compile when they scan internal/logging
    # (mock_exports.go references the generated types).
    mockery --config cfg/mockery.logging.yml || return $?
    local config
    for config in cfg/mockery*.yml; do
        [[ "$config" == "cfg/mockery.logging.yml" ]] && continue
        mockery --config "$config" || return $?
    done
}

go_build() {
    go build -v \
        -ldflags "-X main.version=$DUM_VERSION -X main.commit=$GIT_COMMIT" \
        -o "$DUM_EXECUTABLE-$DUM_VERSION" "$DUM_SOURCE"
}

go_linters_build() {
    mkdir -p "$OUT_DIR"
    (cd cmd/linters && go build -o ../dist/linters .)
}

go_release() {
    goreleaser build --clean --snapshot --config cfg/goreleaser.local.yaml
}

go_release_build() {
    goreleaser build --clean --snapshot --config cfg/goreleaser.snapshot.yaml
}

go_lint() {
    golangci-lint run cmd/...
    golangci-lint run internal/...
}

go_vet() {
    go vet ./...
}

go_fmt() {
    golangci-lint fmt
}

go_fix() {
    go fix ./...
}
