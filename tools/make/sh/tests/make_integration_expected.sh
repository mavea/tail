#!/bin/sh
set -eu

INTEGRATION_DIR="${1}"
APP_DIR="${2}"
if [ ! -d "$INTEGRATION_DIR" ]; then
  echo "integration dir not found: $INTEGRATION_DIR" >&2
  exit 1
fi

for dir in "$INTEGRATION_DIR"/tail/*; do
  [ -d "$dir" ] || continue
  [ -f "$dir/main.go" ] || continue

  name="$(basename "$dir")"
  echo "build integration expected '$name' from $dir"
  go build -ldflags='-w -s' -o "$APP_DIR/$name" "$dir"
  if "$APP_DIR/$name" >"$dir/expected/result.out.log" 2>"$dir/expected/result.err.log"; then
    status=0
  else
    status=$?
  fi
  printf '%s\n' "$status" >"$dir/expected/status.log"

  "$APP_DIR/$name" 2>&1 | $APP_DIR/tail -a test -n 3 -o direct -t none -r roller > "$dir/expected/pipe_none_3_result.out.log" || true
  "$APP_DIR/$name" 2>&1 | $APP_DIR/tail -a test -n 5 -o direct -t none -r roller > "$dir/expected/pipe_none_5_result.out.log" || true
  "$APP_DIR/$name" 2>&1 | $APP_DIR/tail -a test -n 3 -o direct -t full -r roller > "$dir/expected/pipe_full_3_result.out.log" || true
  "$APP_DIR/$name" 2>&1 | $APP_DIR/tail -a test -n 5 -o direct -t full -r roller > "$dir/expected/pipe_full_5_result.out.log" || true
  "$APP_DIR/$name" 2>&1 | $APP_DIR/tail -a test -n 3 -o direct -t minimal -r roller > "$dir/expected/pipe_minimal_3_result.out.log" || true
  "$APP_DIR/$name" 2>&1 | $APP_DIR/tail -a test -n 5 -o direct -t minimal -r roller > "$dir/expected/pipe_minimal_5_result.out.log" || true

  $APP_DIR/tail -a test -n 3 -o direct -t none -r roller -c $APP_DIR/$name > "$dir/expected/command_none_3_result.out.log" 2>"$dir/expected/command_none_3_result.err.log" || true
  $APP_DIR/tail -a test -n 5 -o direct -t none -r roller -c $APP_DIR/$name > "$dir/expected/command_none_5_result.out.log" 2>"$dir/expected/command_none_5_result.err.log" || true
  $APP_DIR/tail -a test -n 3 -o direct -t full -r roller -c $APP_DIR/$name > "$dir/expected/command_full_3_result.out.log" 2>"$dir/expected/command_full_3_result.err.log" || true
  $APP_DIR/tail -a test -n 5 -o direct -t full -r roller -c $APP_DIR/$name > "$dir/expected/command_full_5_result.out.log" 2>"$dir/expected/command_full_5_result.err.log" || true
  $APP_DIR/tail -a test -n 3 -o direct -t minimal -r roller -c $APP_DIR/$name > "$dir/expected/command_minimal_3_result.out.log" 2>"$dir/expected/command_minimal_3_result.err.log" || true
  $APP_DIR/tail -a test -n 5 -o direct -t minimal -r roller -c $APP_DIR/$name > "$dir/expected/command_minimal_5_result.out.log" 2>"$dir/expected/command_minimal_5_result.err.log" || true

  $APP_DIR/tail -a test -n 3 -o direct -t none -r roller $dir/expected/result.out.log > "$dir/expected/file_none_3_result.out.log" || true
  $APP_DIR/tail -a test -n 5 -o direct -t none -r roller $dir/expected/result.out.log > "$dir/expected/file_none_5_result.out.log" || true
  $APP_DIR/tail -a test -n 3 -o direct -t full -r roller $dir/expected/result.out.log > "$dir/expected/file_full_3_result.out.log" || true
  $APP_DIR/tail -a test -n 5 -o direct -t full -r roller $dir/expected/result.out.log > "$dir/expected/file_full_5_result.out.log" || true
  $APP_DIR/tail -a test -n 3 -o direct -t minimal -r roller $dir/expected/result.out.log > "$dir/expected/file_minimal_3_result.out.log" || true
  $APP_DIR/tail -a test -n 5 -o direct -t minimal -r roller $dir/expected/result.out.log > "$dir/expected/file_minimal_5_result.out.log" || true
done

echo "build integration expected help from $INTEGRATION_DIR/help"
LANG=en $APP_DIR/tail -h > "$INTEGRATION_DIR/help/expected/en_h.log" || true
LANG=en $APP_DIR/tail --help > "$INTEGRATION_DIR/help/expected/en_help.log" || true
LANG=ru $APP_DIR/tail -h > "$INTEGRATION_DIR/help/expected/ru_h.log" || true
LANG=ru $APP_DIR/tail --help > "$INTEGRATION_DIR/help/expected/ru_help.log" || true

echo "build integration expected version from $INTEGRATION_DIR/version"
LANG=en $APP_DIR/tail -v > "$INTEGRATION_DIR/version/expected/en_v.log" || true
LANG=en $APP_DIR/tail --version > "$INTEGRATION_DIR/version/expected/en_version.log" || true
LANG=ru $APP_DIR/tail -v > "$INTEGRATION_DIR/version/expected/ru_v.log" || true
LANG=ru $APP_DIR/tail --version > "$INTEGRATION_DIR/version/expected/ru_version.log" || true