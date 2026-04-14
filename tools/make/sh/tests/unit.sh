#!/bin/sh

set -eu

tmp=$(mktemp)
tempLog=$(mktemp)

cleanup() {
  rm -f "$tmp" "$tempLog"
}
trap cleanup EXIT

ginkgo -r ./internal/... ./cmd/... >"$tempLog" 2>&1 || ginkgo_status=$?
ginkgo_status=${ginkgo_status:-0}
sed 's|/src/||g' "$tempLog"
if [ "$ginkgo_status" -ne 0 ]; then
  exit "$ginkgo_status"
fi

go test ./internal/... ./cmd/... -coverprofile="$tmp" >"$tempLog" 2>&1 || go_status=$?
go_status=${go_status:-0}
if [ "$go_status" -ne 0 ]; then
  sed 's|/src/||g' "$tempLog"
  exit "$go_status"
fi

go tool cover -func="$tmp" | grep '^total:'
