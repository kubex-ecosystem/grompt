#!/usr/bin/env bash
set -euo pipefail

_ROOT_DIR="$(git rev-parse --show-toplevel 2>/dev/null)"
_ROOT_DIR="${_ROOT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
_SCRIPT="$(realpath "$_ROOT_DIR/support/string_utils.sh")"
_TMP="$(mktemp)"
trap 'rm -f "$_TMP"' EXIT

# create a temporary exec-mode copy (do not modify original)
awk 'BEGIN{r=0} { if(!r && match($0,/__secure_logic_use_type=/)){ sub(/__secure_logic_use_type="[^"]*"/,"__secure_logic_use_type=\"exec\""); r=1 } print }' "$_SCRIPT" > "$_TMP"

print_ok() { printf '[OK]  %s\n' "${1:-}"; }
print_fail() { printf '[FAIL] %s\n  expected: %s\n  got:      %s\n' "${1:-}" "${2:-}" "${3:-}"; }

run_exec() {
  if [ ! -f "$_TMP" ]; then
    echo "Temporary script '$_TMP' not found!"
    exit 1
  fi
  if [ ! -x "$_TMP" ]; then
    chmod +x "$_TMP" || {
      echo "Failed to make '$_TMP' executable!"
      exit 1
    }
  fi
  local name="${1:-}"; shift
  local expected="${1:-}"; shift
  local out
  out="$( "$_TMP" "$@" )" || out="$out"
  if [ "$out" = "$expected" ]; then
    print_ok "exec:$name"
  else
    print_fail "exec:$name" "$expected" "$out"
  fi
}

run_source() {
  local name="${1:-}"; shift
  local expected="${1:-}"; shift
  local out
  out="$( "$@" )" || out="$out"
  if [ "$out" = "$expected" ]; then
    print_ok "src:$name"
  else
    print_fail "src:$name" "$expected" "$out"
  fi
}

echo "== Exec-mode tests (uses temporary copy) =="
run_exec "toLowerCase" "hello" toLowerCase "HeLLo"
run_exec "replace_first" "baz bar foo" replace_first "foo bar foo" "foo" "baz"
run_exec "replace_nth" "one two ONE2 two one" replace_nth "one two one two one" "one" "ONE2" "2"
run_exec "replace_case_sensitive_last" "alpha beta alpha gamma ALAST" replace_case_sensitive_last "alpha beta alpha gamma alpha" "alpha" "ALAST"

echo
echo "== Source-mode tests (source original) =="
# force lib mode in environment before sourcing to be explicit
export __secure_logic_use_type="lib"
# shellcheck source=/dev/null
. "$_SCRIPT"

# call functions now defined in this shell
run_source "toLowerCase" "hello" toLowerCase "HeLLo"
run_source "replace_first" "baz bar foo" replace_first "foo bar foo" "foo" "baz"
run_source "replace_nth" "one two ONE2 two one" replace_nth "one two one two one" "one" "ONE2" "2"
run_source "replace_case_sensitive_last" "alpha beta alpha gamma ALAST" replace_case_sensitive_last "alpha beta alpha gamma alpha" "alpha" "ALAST"

echo
echo "All tests done."
