#!/bin/sh
set -eu

INTEGRATION_DIR="${1}"
APP_DIR="${2}"
if [ ! -d "$INTEGRATION_DIR" ]; then
  echo "integration dir not found: $INTEGRATION_DIR" >&2
  exit 1
fi

for dir in "$INTEGRATION_DIR"/*; do
  [ -d "$dir" ] || continue
  [ -f "$dir/main.go" ] || continue

  name="$(basename "$dir")"
  echo "test integration '$name' from $dir"
  go build -ldflags='-w -s' -o "$APP_DIR/$name" "$dir"
  "$APP_DIR/$name" >"$dir/expected/result.log" 2>>"$dir/expected/result.log" || true

  "$APP_DIR/$name" 2>&1 | $APP_DIR/tail -a test -n 3 -o direct -t none -r roll | diff -U 3 - "$dir/expected/pipe_none_3_result.log" || true
  "$APP_DIR/$name" 2>&1 | $APP_DIR/tail -a test -n 5 -o direct -t none -r roll | diff -U 3 - "$dir/expected/pipe_none_5_result.log" || true
  "$APP_DIR/$name" 2>&1 | $APP_DIR/tail -a test -n 3 -o direct -t minimal -r roll | diff -U 3 - "$dir/expected/pipe_minimal_3_result.log" || true
  "$APP_DIR/$name" 2>&1 | $APP_DIR/tail -a test -n 5 -o direct -t minimal -r roll | diff -U 3 - "$dir/expected/pipe_minimal_5_result.log" || true

  $APP_DIR/tail -a test -n 3 -o direct -t none -r roll -c $APP_DIR/$name | diff -U 3 - "$dir/expected/command_none_3_result.log" || true
  $APP_DIR/tail -a test -n 5 -o direct -t none -r roll -c $APP_DIR/$name | diff -U 3 - "$dir/expected/command_none_5_result.log" || true
  $APP_DIR/tail -a test -n 3 -o direct -t minimal -r roll -c $APP_DIR/$name | diff -U 3 - "$dir/expected/command_minimal_3_result.log" || true
  $APP_DIR/tail -a test -n 5 -o direct -t minimal -r roll -c $APP_DIR/$name | diff -U 3 - "$dir/expected/command_minimal_5_result.log" || true

  $APP_DIR/tail -a test -n 3 -o direct -t none -r roll $dir/expected/result.log | diff -U 3 - "$dir/expected/file_none_3_result.log" || true
  $APP_DIR/tail -a test -n 5 -o direct -t none -r roll $dir/expected/result.log | diff -U 3 - "$dir/expected/file_none_5_result.log" || true
  $APP_DIR/tail -a test -n 3 -o direct -t minimal -r roll $dir/expected/result.log | diff -U 3 - "$dir/expected/file_minimal_3_result.log" || true
  $APP_DIR/tail -a test -n 5 -o direct -t minimal -r roll $dir/expected/result.log | diff -U 3 - "$dir/expected/file_minimal_5_result.log" || true
done