#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<EOF
Usage: $(basename "$0") [major|minor|patch]

Calculates the next semantic version, updates pkg/version/version.go,
and refreshes documentation references. Defaults to patch bump.
EOF
}

if [[ "${1:-}" =~ ^(-h|--help)$ ]]; then
  usage
  exit 0
fi

BUMP="${1:-patch}"
case "$BUMP" in
  major|minor|patch) ;;
  *)
    echo "error: bump type must be major, minor, or patch" >&2
    usage
    exit 1
    ;;
esac

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

if ! git diff --quiet --ignore-submodules --exit-code || \
   ! git diff --quiet --ignore-submodules --cached --exit-code; then
  echo "error: please start from a clean working tree before bumping a release" >&2
  exit 1
fi

VERSION_FILE="pkg/version/version.go"
if [[ ! -f "$VERSION_FILE" ]]; then
  echo "error: $VERSION_FILE not found" >&2
  exit 1
fi

CURRENT_VERSION="$(grep -Eo 'Version = \"v[0-9]+\.[0-9]+\.[0-9]+\"' "$VERSION_FILE" | sed -E 's/.*\"(v[0-9]+\.[0-9]+\.[0-9]+)\"/\1/')"
if [[ -z "$CURRENT_VERSION" ]]; then
  echo "error: unable to parse current version from $VERSION_FILE" >&2
  exit 1
fi

BASE="${CURRENT_VERSION#v}"
IFS='.' read -r MAJOR MINOR PATCH <<<"$BASE"
case "$BUMP" in
  major)
    ((MAJOR++))
    MINOR=0
    PATCH=0
    ;;
  minor)
    ((MINOR++))
    PATCH=0
    ;;
  patch)
    ((PATCH++))
    ;;
esac
NEXT_VERSION="v${MAJOR}.${MINOR}.${PATCH}"

if [[ "$NEXT_VERSION" == "$CURRENT_VERSION" ]]; then
  echo "version unchanged ($CURRENT_VERSION)" >&2
  exit 0
fi

cat >"$VERSION_FILE" <<EOF
package version

import "strings"

// Version follows semantic versioning (vMAJOR.MINOR.PATCH).
const Version = "${NEXT_VERSION}"

// UserAgent returns the default user agent string shared by the SDK.
func UserAgent() string {
	return "port-go-sdk/" + strings.TrimPrefix(Version, "v")
}
EOF

gofmt -w "$VERSION_FILE"

NEW_VERSION="$NEXT_VERSION" python3 <<'PY'
import os
import pathlib
import re
import sys

ver = os.environ["NEW_VERSION"]
readme = pathlib.Path("README.md")
text = readme.read_text()
pattern = r'(go get github\.com/port-experimental/port-go-sdk@)v\d+\.\d+\.\d+'
new_text, count = re.subn(pattern, r'\1' + ver, text, count=1)
if count == 0:
    sys.stderr.write("warning: README snippet not updated; pattern missing?\n")
else:
    readme.write_text(new_text)
PY

echo "Updated version to ${NEXT_VERSION}"
echo
echo "Next steps:"
echo "  1. Review changes (git status)."
echo "  2. Run tests/lint (e.g., make test, golangci-lint)."
echo "  3. Commit and tag: git commit -am \"Release ${NEXT_VERSION}\" && git tag ${NEXT_VERSION}"
echo "  4. Push: git push origin HEAD && git push origin ${NEXT_VERSION}"
