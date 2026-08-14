#!/usr/bin/env bash

# Repository cleanup and hook validation for the repository dispatcher.

repository_rhard() {
    git reset --hard HEAD
    for directory in vim/vim.d/plugged/*/; do
        if [[ -d "$directory/.git" ]]; then
            echo "Resetting $directory ..."
            git -C "$directory" reset --hard HEAD
            git -C "$directory" clean -fd
        fi
    done
}

repository_clean() {
    go clean -cache
    rm -rf "$OUT_DIR"
    rm -f ./internal/**/*_mocks.go
    rm -f ./internal/**/*_enum.go
}

repository_check_lfs_hook() {
    local hooks_dir
    hooks_dir="$(git -C "$REPO_ROOT" rev-parse --git-path hooks)"
    if ! grep -q '^git lfs pre-push' "$hooks_dir/pre-push"; then
        echo "git-lfs is not installed properly. Please run 'git lfs install' and try again." >&2
        return 1
    fi
}
