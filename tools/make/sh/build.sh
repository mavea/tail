#!/bin/sh
set -eu

VERSION_FILE=/src/version
[ -f "$VERSION_FILE" ] || { printf 'error: version file not found: %s\n' "$VERSION_FILE" >&2; exit 1; }
APP_VERSION=$(sed -n 's/^VERSION=//p' "$VERSION_FILE" | head -n 1 | tr -d '\r')
APP_RELEASE_DATE=$(sed -n 's/^RELEASE_DATE=//p' "$VERSION_FILE" | head -n 1 | tr -d '\r')
[ -n "$APP_VERSION" ] || { printf 'error: VERSION not set in %s\n' "$VERSION_FILE" >&2; exit 1; }
[ -n "$APP_RELEASE_DATE" ] || { printf 'error: RELEASE_DATE not set in %s\n' "$VERSION_FILE" >&2; exit 1; }
LDFLAGS="-w -s -X tail/internal/bootstrap.Version=$APP_VERSION -X tail/internal/bootstrap.BuildTime=$APP_RELEASE_DATE"

tmp=$(mktemp)
set +e
go build -trimpath -ldflags="$LDFLAGS" -o "$1" "$2" >"$tmp" 2>&1
rc=$?
set -e
[ -s "$tmp" ] && sed 's|/src/||g' "$tmp" >&2
rm -f "$tmp"
[ "$rc" -eq 0 ] || exit "$rc"
chmod +x "$1"

