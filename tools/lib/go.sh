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
    for config in cfg/mockery*.yml; do
        mockery --config "$config"
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
