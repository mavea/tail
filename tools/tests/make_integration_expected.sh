#!/bin/sh
set -eu

INTEGRATION_DIR="${1}"
if [ ! -d "$INTEGRATION_DIR" ]; then
  echo "integration dir not found: $INTEGRATION_DIR" >&2
  exit 1
fi

for dir in "$INTEGRATION_DIR"/*; do
  [ -d "$dir" ] || continue
  [ -f "$dir/main.go" ] || continue

  name="$(basename "$dir")"
  echo "build integration expected '$name' from $dir"
  go build -ldflags='-w -s' -o "/app/$name" "$dir"
  "/app/$name" >"$dir/expected/result.log" 2>>"$dir/expected/result.log" || true

  "/app/$name" 2>&1 | /app/app -a test -n 3 -o direct -t none -r roll > "$dir/expected/pipe_none_3_result.log" || true
  "/app/$name" 2>&1 | /app/app -a test -n 5 -o direct -t none -r roll > "$dir/expected/pipe_none_5_result.log" || true
  "/app/$name" 2>&1 | /app/app -a test -n 3 -o direct -t minimal -r roll > "$dir/expected/pipe_minimal_3_result.log" || true
  "/app/$name" 2>&1 | /app/app -a test -n 5 -o direct -t minimal -r roll > "$dir/expected/pipe_minimal_5_result.log" || true

  /app/app -a test -n 3 -o direct -t none -r roll -c /app/$name > "$dir/expected/command_none_3_result.log" || true
  /app/app -a test -n 5 -o direct -t none -r roll -c /app/$name > "$dir/expected/command_none_5_result.log" || true
  /app/app -a test -n 3 -o direct -t minimal -r roll -c /app/$name > "$dir/expected/command_minimal_3_result.log" || true
  /app/app -a test -n 5 -o direct -t minimal -r roll -c /app/$name > "$dir/expected/command_minimal_5_result.log" || true

  /app/app -a test -n 3 -o direct -t none -r roll $dir/expected/result.log > "$dir/expected/file_none_3_result.log" || true
  /app/app -a test -n 5 -o direct -t none -r roll $dir/expected/result.log > "$dir/expected/file_none_5_result.log" || true
  /app/app -a test -n 3 -o direct -t minimal -r roll $dir/expected/result.log > "$dir/expected/file_minimal_3_result.log" || true
  /app/app -a test -n 5 -o direct -t minimal -r roll $dir/expected/result.log > "$dir/expected/file_minimal_5_result.log" || true
done