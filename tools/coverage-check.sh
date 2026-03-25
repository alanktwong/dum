#!/bin/bash

# Set default coverage threshold
COVERAGE_THRESHOLD=${COVERAGE_THRESHOLD:-80}

# Set output directory
OUT_DIR="./dist"

# Check coverage threshold
echo "Checking coverage threshold ..."
grep -v '_enum.go' "${OUT_DIR}/coverage.txt" > "${OUT_DIR}/coverage-filtered.txt"
COVERAGE=$(go tool cover -func="${OUT_DIR}/coverage-filtered.txt" | grep total | awk '{print $3}' | tr -d '%')
if [ "$(echo "$COVERAGE < $COVERAGE_THRESHOLD" | bc -l)" -eq 1 ]; then
    echo "Coverage ${COVERAGE}% is below threshold ${COVERAGE_THRESHOLD}%"
    exit 1
else
    echo "Coverage ${COVERAGE}% meets threshold ${COVERAGE_THRESHOLD}%"
fi