#!/usr/bin/env bash
set -euo pipefail

maxkb="${1:?maximum size in KB is required}"
shift

files=()
for file in "$@"; do
    if git check-attr filter -- "$file" | grep -q ': filter: lfs$'; then
        continue
    fi

    files+=("$file")
done

if ((${#files[@]} == 0)); then
    exit 0
fi

hk util check-added-large-files --maxkb "$maxkb" "${files[@]}"
