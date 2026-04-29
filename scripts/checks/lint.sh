#!/bin/sh
set -e

if ! command -v golangci-lint >/dev/null 2>&1; then
    echo "ERROR: golangci-lint is not installed."
    echo ""
    echo "Install via:"
    echo "  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest"
    echo "  or: brew install golangci-lint"
    exit 1
fi

golangci-lint run --fix ./...

# Re-stage only previously staged files that were auto-fixed by golangci-lint
FIXED=$(git diff --name-only)
if [ -n "$FIXED" ]; then
    STAGED=$(git diff --cached --name-only)
    # Intersect: only re-stage files that were already staged
    RESTAGE=$(echo "$FIXED" | grep -xF "$STAGED" 2>/dev/null || true)
    if [ -n "$RESTAGE" ]; then
        echo "$RESTAGE" | xargs git add
    fi
fi
