#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "==> gofmt"
test -z "$(gofmt -l . | tee /dev/stderr)"

echo "==> go test ./..."
go test ./...

echo "==> build examples"
for dir in examples/*/*; do
  if [ -d "$dir" ]; then
    go build "./$dir"
  fi
done

echo "All checks passed."
