#!/usr/bin/env bash
set -euo pipefail

TOOLS=(
  "golang.org/x/tools/cmd/goimports@latest"
  "honnef.co/go/tools/cmd/staticcheck@latest"
  "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
)

echo "Installing optional development tools..."
for tool in "${TOOLS[@]}"; do
  echo "  -> $tool"
  go install "$tool"
done

echo "Done. Make sure \$GOBIN is on your PATH."
