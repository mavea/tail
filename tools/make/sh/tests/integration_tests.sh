#!/bin/sh
set -eu

if [ "$#" -lt 2 ]; then
  echo "usage: $0 <integration_dir> <app_dir> [--color=always]" >&2
  exit 2
fi

INTEGRATION_DIR="$1"
APP_DIR="$2"
shift 2
GLOBAL_DIR="$(cd "$INTEGRATION_DIR" && pwd -P)"
LOCAL_DIR="$(cd "./" && pwd -P)"

FORCE_COLOR=0
while [ "$#" -gt 0 ]; do
  case "$1" in
    --color=always)
      FORCE_COLOR=1
      ;;
    *)
      echo "unknown option: $1" >&2
      echo "usage: $0 <integration_dir> <app_dir> [--color=always]" >&2
      exit 2
      ;;
  esac
  shift
done

if [ ! -d "$INTEGRATION_DIR" ]; then
  echo "integration dir not found: $INTEGRATION_DIR" >&2
  exit 1
fi


total=0
failed=0

# ANSI colors (enabled for TTY or forced by --color=always)
if [ -t 1 ] || [ "$FORCE_COLOR" -eq 1 ]; then
  C_RESET=$(printf '\033[0m')
  C_GREEN=$(printf '\033[32m')
  C_RED=$(printf '\033[31m')
  C_YELLOW=$(printf '\033[33m')
  C_CYAN=$(printf '\033[36m')
else
  C_RESET=""
  C_GREEN=""
  C_RED=""
  C_YELLOW=""
  C_CYAN=""
fi

# GNU diff color support (best effort)
if [ -n "$C_RESET" ] && diff --color=always -u /dev/null /dev/null >/dev/null 2>&1; then
  DIFF_COLOR_OPT="--color=always"
else
  DIFF_COLOR_OPT=""
fi

print_diff_header() {
   diff="$1"
   expected="$2"
   tmp="$3"
  first_diff=$(cmp -l "$expected" "$tmp" 2>/dev/null | head -1)

  if [ -n "$first_diff" ]; then
    offset=$(printf '%s' "$first_diff" | awk '{print $1}')
    line_num=$(head -c "$offset" "$expected" 2>/dev/null | grep -c '' || echo 1)
    char_pos=$((offset - $(head -n $((line_num - 1)) "$expected" 2>/dev/null | wc -c)))
    printf "%s:%d:%d\n" ${expected#"$LOCAL_DIR"/} "$line_num" "$char_pos"
  else
    printf "%s:1\n" ${expected#"$LOCAL_DIR"/}
  fi
}

print_diff() {
   diff="$1"
   expected="$2"
   tmp="$3"
  printf "%s--------------------------DIFF----------------------------%s\n" "$C_CYAN" "$C_RESET"
  if [ -s "$diff" ]; then
    if [ -n "$DIFF_COLOR_OPT" ]; then
      diff $DIFF_COLOR_OPT -u "$expected" "$tmp" || true
    else
      cat -v "$diff"
    fi
  else
    printf "%s(diff output is empty, showing first differing bytes)%s\n" "$C_YELLOW" "$C_RESET"
    cmp -l "$expected" "$tmp" | head -n 20 || true
  fi
  printf "%s----------------------------------------------------------%s\n" "$C_CYAN" "$C_RESET"
}

# check_cmd <label> <expected_stdout> <expected_stderr_or_empty> <expected_status> <command...>
check_cmd() {
  label="$1"
  expected_stdout="$2"
  expected_stderr="$3"
  expected_status="$4"
  shift 4

  tmp_stdout=$(mktemp)
  tmp_stderr=$(mktemp)
  tmp_diff=$(mktemp)

  # set -e safe capture of exit status
  if "$@" >"$tmp_stdout" 2>"$tmp_stderr"; then
    actual_status=0
  else
    actual_status=$?
  fi

  case_failed=0

  if [ "$expected_status" != "-" ] && [ "$actual_status" -ne "$expected_status" ]; then
    printf "%sstatus diff for '%s': expected=%s actual=%s%s\n" "$C_RED" "$label" "$expected_status" "$actual_status" "$C_RESET"
    case_failed=1
  fi

  if ! diff -u "$expected_stdout" "$tmp_stdout" >"$tmp_diff" 2>&1; then
    printf "%sstdout diff for '%s':%s " "$C_RED" "$label" "$C_RESET"
    print_diff_header "$tmp_diff" "$expected_stdout" "$tmp_stdout"
    print_diff "$tmp_diff" "$expected_stdout" "$tmp_stdout"
    case_failed=1
  fi

  if [ "$expected_stderr" = "-" ]; then
    :
  elif [ -n "$expected_stderr" ]; then
    if ! diff -u "$expected_stderr" "$tmp_stderr" >"$tmp_diff" 2>&1; then
      printf "%sstderr diff for '%s':%s " "$C_RED" "$label" "$C_RESET"
      print_diff_header "$tmp_diff" "$expected_stderr" "$tmp_stderr"
      print_diff "$tmp_diff" "$expected_stderr" "$tmp_stderr"
      case_failed=1
    fi
  elif [ -s "$tmp_stderr" ]; then
    printf "%sstderr diff for '%s': expected empty, got output%s\n" "$C_RED" "$label" "$C_RESET"
    sed -n 'l' "$tmp_stderr"
    case_failed=1
  fi

  rm -f "$tmp_stdout" "$tmp_stderr" "$tmp_diff"

  total=$((total + 1))
  if [ "$case_failed" -eq 0 ]; then
    printf "    %sOK%s   %s\n" "$C_GREEN" "$C_RESET" "$label"
  else
    printf "    %sFAIL%s %s\n" "$C_RED" "$C_RESET" "$label"
    failed=$((failed + 1))
  fi
}

# check_sh <label> <expected_stdout> <expected_stderr_or_empty> <expected_status> <script>
check_sh() {
  label="$1"
  expected_stdout="$2"
  expected_stderr="$3"
  expected_status="$4"
  script="$5"
  check_cmd "$label" "$expected_stdout" "$expected_stderr" "$expected_status" sh -c "$script"
}

for dir in "$GLOBAL_DIR"/tail/*; do
  [ -d "$dir" ] || continue
  [ -f "$dir/main.go" ] || continue

  # Ensure per-iteration path is absolute.
  case "$dir" in
    /*) ;;
    *) dir="$(cd "$dir" && pwd -P)" ;;
  esac

  name="$(basename "$dir")"
  echo ${dir#"$LOCAL_DIR"/}": test integration '$name'"

  go build -ldflags='-w -s' -o "$APP_DIR/$name" "$dir"

  expected_status=0
  if [ -f "$dir/expected/status.log" ]; then
    expected_status="$(tr -d '[:space:]' <"$dir/expected/status.log")"
  fi

  echo "   check test cmd"
  check_cmd "cmd" "$dir/expected/result.out.log" "$dir/expected/result.err.log" "$expected_status" "$APP_DIR/$name"

  echo "   check pipe"
  check_sh "pipe_none_3" "$dir/expected/pipe_none_3_result.out.log" "-" 0 "\"$APP_DIR/$name\" 2>&1 | \"$APP_DIR/tail\" -a test -n 3 -o direct -t none -r roller"
  check_sh "pipe_none_5" "$dir/expected/pipe_none_5_result.out.log" "-" 0 "\"$APP_DIR/$name\" 2>&1 | \"$APP_DIR/tail\" -a test -n 5 -o direct -t none -r roller"
  check_sh "pipe_minimal_3" "$dir/expected/pipe_minimal_3_result.out.log" "-" 0 "\"$APP_DIR/$name\" 2>&1 | \"$APP_DIR/tail\" -a test -n 3 -o direct -t minimal -r roller"
  check_sh "pipe_minimal_5" "$dir/expected/pipe_minimal_5_result.out.log" "-" 0 "\"$APP_DIR/$name\" 2>&1 | \"$APP_DIR/tail\" -a test -n 5 -o direct -t minimal -r roller"
  check_sh "pipe_full_3" "$dir/expected/pipe_full_3_result.out.log" "-" 0 "\"$APP_DIR/$name\" 2>&1 | \"$APP_DIR/tail\" -a test -n 3 -o direct -t full -r roller"
  check_sh "pipe_full_5" "$dir/expected/pipe_full_5_result.out.log" "-" 0 "\"$APP_DIR/$name\" 2>&1 | \"$APP_DIR/tail\" -a test -n 5 -o direct -t full -r roller"

  echo "   check file"
  check_cmd "file_none_3" "$dir/expected/file_none_3_result.out.log" "-" 0 "$APP_DIR/tail" -a test -n 3 -o direct -t none -r roller "$dir/expected/result.out.log"
  check_cmd "file_none_5" "$dir/expected/file_none_5_result.out.log" "-" 0 "$APP_DIR/tail" -a test -n 5 -o direct -t none -r roller "$dir/expected/result.out.log"
  check_cmd "file_minimal_3" "$dir/expected/file_minimal_3_result.out.log" "-" 0 "$APP_DIR/tail" -a test -n 3 -o direct -t minimal -r roller "$dir/expected/result.out.log"
  check_cmd "file_minimal_5" "$dir/expected/file_minimal_5_result.out.log" "-" 0 "$APP_DIR/tail" -a test -n 5 -o direct -t minimal -r roller "$dir/expected/result.out.log"
  check_cmd "file_full_3" "$dir/expected/file_full_3_result.out.log" "-" 0 "$APP_DIR/tail" -a test -n 3 -o direct -t full -r roller "$dir/expected/result.out.log"
  check_cmd "file_full_5" "$dir/expected/file_full_5_result.out.log" "-" 0 "$APP_DIR/tail" -a test -n 5 -o direct -t full -r roller "$dir/expected/result.out.log"

  echo "   check exec"
  check_cmd "command_none_3" "$dir/expected/command_none_3_result.out.log" "$dir/expected/command_none_3_result.err.log" "-" "$APP_DIR/tail" -a test -n 3 -o direct -t none -r roller -c "$APP_DIR/$name"
  check_cmd "command_none_5" "$dir/expected/command_none_5_result.out.log" "$dir/expected/command_none_5_result.err.log" "-" "$APP_DIR/tail" -a test -n 5 -o direct -t none -r roller -c "$APP_DIR/$name"
  check_cmd "command_minimal_3" "$dir/expected/command_minimal_3_result.out.log" "$dir/expected/command_minimal_3_result.err.log" "-" "$APP_DIR/tail" -a test -n 3 -o direct -t minimal -r roller -c "$APP_DIR/$name"
  check_cmd "command_minimal_5" "$dir/expected/command_minimal_5_result.out.log" "$dir/expected/command_minimal_5_result.err.log" "-" "$APP_DIR/tail" -a test -n 5 -o direct -t minimal -r roller -c "$APP_DIR/$name"
  check_cmd "command_full_3" "$dir/expected/command_full_3_result.out.log" "$dir/expected/command_full_3_result.err.log" "-" "$APP_DIR/tail" -a test -n 3 -o direct -t full -r roller -c "$APP_DIR/$name"
  check_cmd "command_full_5" "$dir/expected/command_full_5_result.out.log" "$dir/expected/command_full_5_result.err.log" "-" "$APP_DIR/tail" -a test -n 5 -o direct -t full -r roller -c "$APP_DIR/$name"
done

if [ -d "$GLOBAL_DIR/help/expected" ]; then
  echo "test integration 'help'"
  check_cmd "help_en_h" "$GLOBAL_DIR/help/expected/en_h.log" "" 0 env LANG=en "$APP_DIR/tail" -h
  check_cmd "help_en_help" "$GLOBAL_DIR/help/expected/en_help.log" "" 0 env LANG=en "$APP_DIR/tail" --help
  check_cmd "help_ru_h" "$GLOBAL_DIR/help/expected/ru_h.log" "" 0 env LANG=ru "$APP_DIR/tail" -h
  check_cmd "help_ru_help" "$GLOBAL_DIR/help/expected/ru_help.log" "" 0 env LANG=ru "$APP_DIR/tail" --help
fi

if [ -d "$GLOBAL_DIR/version/expected" ]; then
  echo "test integration 'version'"
  check_cmd "version_en_v" "$GLOBAL_DIR/version/expected/en_v.log" "" 0 env LANG=en "$APP_DIR/tail" -v
  check_cmd "version_en_version" "$GLOBAL_DIR/version/expected/en_version.log" "" 0 env LANG=en "$APP_DIR/tail" --version
  check_cmd "version_ru_v" "$GLOBAL_DIR/version/expected/ru_v.log" "" 0 env LANG=ru "$APP_DIR/tail" -v
  check_cmd "version_ru_version" "$GLOBAL_DIR/version/expected/ru_version.log" "" 0 env LANG=ru "$APP_DIR/tail" --version
fi

echo ""
echo "Integration checks: $((total - failed))/$total passed"
if [ "$failed" -gt 0 ]; then
  printf "%sFAILED: %s case(s)%s\n" "$C_RED" "$failed" "$C_RESET"
  exit 1
fi

echo "All integration checks passed"
