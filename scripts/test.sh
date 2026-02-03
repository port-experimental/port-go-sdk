#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

echo "==> make lint"
make lint

echo "==> make test"
make test

echo "==> make examples"
make examples

echo "All checks passed."
