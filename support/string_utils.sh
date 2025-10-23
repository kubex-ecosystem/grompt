#!/usr/bin/env bash
# shellcheck disable=SC2015,SC1091

# Script Metadata
__secure_logic_version="1.0.0"
__secure_logic_date="$( date +%Y-%m-%d )"
__secure_logic_author="Rafael Mori"
__secure_logic_use_type="lib"
__secure_logic_init_timestamp="$(date +%s)"
__secure_logic_elapsed_time=0

# Check if verbose mode is enabled
if [[ "${MYNAME_VERBOSE:-false}" == "true" ]]; then
  set -x  # Enable debugging
fi

IFS=$'\n\t'

_ROOT_DIR="$(git rev-parse --show-toplevel 2>/dev/null)"
_ROOT_DIR="${_ROOT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"

declare -a _main_args=( "$@" )

__get_output_tty() {
  if [[ -t 1 ]]; then
    echo '/dev/tty'
  else
    echo '/dev/stderr'
  fi
}

__secure_logic_sourced_name() {
  # prefer BASH_SOURCE[0], fallback to $0 for shells where BASH_SOURCE may be absent
  local _self="${BASH_SOURCE[0]:-$0}"
  _self="${_self//${_ROOT_DIR:-}/}"
  _self="${_self//.sh/}"
  _self="${_self//-/_}"
  _self="${_self//\//_}"
  _self="${_self//./_}"
  _self="${_self// /_}"
  echo "_was_sourced_${_self//__/_}"
  return 0
}

__first(){
  if [ "$EUID" -eq 0 ] || [ "$UID" -eq 0 ]; then
    printf '%s\n' "Please do not run as root." >"$(__get_output_tty)"
    exit 1
  elif [ -n "${SUDO_USER:-}" ]; then
    printf '%s\n' "Please do not run as root, but with sudo privileges." >"$(__get_output_tty)"
    exit 1
  else
    # shellcheck disable=SC2155
    local _ws_name="$(__secure_logic_sourced_name)"

    # detect if script was sourced: compare BASH_SOURCE[0] (or fallback) with $0
    if [ "${BASH_SOURCE[0]:-$0}" != "$0" ]; then
      if test ${__secure_logic_use_type:-} != "lib"; then
        printf '%s\n' "This script is not intended to be sourced." >"$(__get_output_tty)"
        printf '%s\n' "Please run it directly." >"$(__get_output_tty)"
        exit 1
      fi
      # If the script is sourced, we set the variable to true
      # and export it to the environment without changing
      # the shell options.
      export "${_ws_name:-}"="true"
    else
      if test ${__secure_logic_use_type:-} != "exec"; then
        printf '%s\n' "This script is not intended to be executed directly." >"$(__get_output_tty)"
        printf '%s\n' "Please source it instead." >"$(__get_output_tty)"
        exit 1
      fi
      # If the script is executed directly, we set the variable to false
      # and export it to the environment. We also set the shell options
      # to ensure a safe execution.
      export "${_ws_name:-}"="false"
      set -o errexit # Exit immediately if a command exits with a non-zero status
      set -o nounset # Treat unset variables as an error when substituting
      set -o pipefail # Return the exit status of the last command in the pipeline that failed
      set -o errtrace # If a command fails, the shell will exit immediately
      set -o functrace # If a function fails, the shell will exit immediately
      # shopt -s inherit_errexit # Inherit the errexit option in functions

      if [[ "${_DEBUG:-}" == "true" ]]; then
        set -x
      fi
    fi
  fi
}

_QUIET=${_QUIET:-${QUIET:-false}}
_DEBUG=${_DEBUG:-${DEBUG:-false}}
_HIDE_ABOUT=${_HIDE_ABOUT:-${HIDE_ABOUT:-false}}

__first "${_main_args[@]}" >&2 || {
  echo "Error: This script must be run directly, not sourced." >&2
  exit 1
}

_SCRIPT_DIR="${_ROOT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}/support"
source "${_SCRIPT_DIR}/utils.sh" || {
  echo "Error: Could not source ${_SCRIPT_DIR}/utils.sh. Please ensure it exists." >&2
  exit 1
}

__main_functions() {
  if [[ $# -gt 0 ]]; then
    local func_name="${1:-}"
    local _full_args=( "${_main_args[@]}" )
    if declare -F "$func_name" >'/dev/null' 2>&1; then
      "$func_name" "${@:2}"
    else
      if test -x "${1:-}"; then
        _default "${_full_args[@]:1}"
      else
        log error "Function '${func_name}' not found."
        return 1
      fi
    fi
    return $?
  fi
}

snake_to_camel() {
  echo "$1" | sed -e 's/_\([a-z]\)/\U\1/g' -e 's/^./\L&/'
}
camel_to_snake() {
  echo "$1" | sed -r 's/([A-Z])/_\L\1/g'
}
snake_to_kebab() {
  echo "$1" | tr '_' '-'
}
camel_to_kebab() {
  echo "$1" | sed -r 's/([A-Z])/-\L\1/g' | sed 's/^-//'
}
convert_case() {
  local text="$1"
  local target_case="$2"

  case "$target_case" in
    snake)
      echo "$text" | sed -r 's/([A-Z])/_\L\1/g' | tr '-' '_'
      ;;
    kebab)
      echo "$text" | sed -r 's/([A-Z])/-\L\1/g' | tr '_' '-'
      ;;
    camel)
      echo "$text" | sed -r 's/(^|-|_)([a-z])/\U\2/g' | sed -r 's/[-_]//g'
      ;;
    *)
      echo "Invalid case specified. Use 'snake', 'kebab', or 'camel'."
      return 1
      ;;
  esac
}
toLowerCase() {
  echo "$1" | tr '[:upper:]' '[:lower:]'
}
toUpperCase() {
  echo "$1" | tr '[:lower:]' '[:upper:]'
}
type_text() {
  local text="$1"
  local delay="${2:-0.1}"  # Default delay is 0.1 seconds
  local alignment="${3:-left}"  # Default alignment is left
  local term_width
  term_width=$(tput cols)
  local text_length=${#text}
  local padding

  case "$alignment" in
    center)
      padding=$(( (term_width - text_length) / 2 ))
      ;;
    right)
      padding=$(( term_width - text_length ))
      ;;
    *)
      padding=0
      ;;
  esac

  for ((i=0; i<padding; i++)); do
    echo -n " "
  done

  for ((i=0; i<${#text}; i++)); do
    echo -n "${text:$i:1}"
    sleep "$delay"
  done
  echo
}
replace_first() {
  local text="$1"
  local search="$2"
  local replace="$3"

  echo "${text/$search/$replace}"
}
replace_all() {
  local text="$1"
  local search="$2"
  local replace="$3"

  echo "${text//"$search"/"$replace"}"
}
replace_last() {
  local text="$1"
  local search="$2"
  local replace="$3"

  # shellcheck disable=SC2295
  echo "${text%$search*}$replace${text##*$search}"
}
replace_nth() {
  local text="$1"
  local search="$2"
  local replace="$3"
  local n="$4"

  # deterministic awk: initialize counter inside BEGIN and use original line and offsets
  echo "$text" | awk -v search="$search" -v replace="$replace" -v n="$n" 'BEGIN { i=0 }
  {
    line = $0
    out = ""
    pos = 1
    while (match(line, search)) {
      # absolute position of match in original string
      R = RSTART
      L = RLENGTH
      # if this is the nth match, rebuild and print
      if (++i == n) {
        prefix = substr($0, 1, pos + R - 2)
        suffix = substr($0, pos + R - 1 + L)
        print prefix replace suffix
        next
      }
      # advance pos and cut the processed prefix from line
      pos += R + L - 1
      line = substr(line, R + L)
    }
    print $0
  }'
}
replace_case_insensitive() {
  local text="$1"
  local search="$2"
  local replace="$3"

  # shellcheck disable=SC2001
  echo "$text" | sed "s/$search/$replace/Ig"
}
replace_case_sensitive() {
  local text="$1"
  local search="$2"
  local replace="$3"

  # shellcheck disable=SC2001
  echo "$text" | sed "s/$search/$replace/g"
}
replace_case_sensitive_first() {
  local text="$1"
  local search="$2"
  local replace="$3"

  echo "$text" | sed "0,/$search/s/$search/$replace/"
}
replace_case_sensitive_last() {
  local text="$1"
  local search="$2"
  local replace="$3"

  # prefer rev (available on Linux); fallback to perl for replacing the last occurrence
  if command -v rev >/dev/null 2>&1; then
    printf '%s' "$text" | rev | sed "s/$(printf '%s' "$search" | rev)/$(printf '%s' "$replace" | rev)/" | rev
  else
    # perl fallback: capture everything before the last occurrence and append replacement
    printf '%s' "$text" | perl -0777 -pe 's/(.*)\Q'"$search"'\E/$1'"$replace"'/s'
  fi
}
__secure_logic_main() {
  local _ws_name
  _ws_name="$(__secure_logic_sourced_name)"
  local _ws_name_val
  _ws_name_val=$(eval "echo \${${_ws_name:-}}")
  if test "${_ws_name_val:-}" != "true"; then
    __main_functions "${_main_args[@]}"
    return $?
  else
    # shellcheck disable=SC2207
    local _functions_to_export=( $(compgen -A function | grep -v -e '^__') )
    for func in "${_functions_to_export[@]}"; do
      # shellcheck disable=SC2163
      export -f "$func" || {
        echo "Error: Could not export function '$func'." >&2
        return 1
      }
    done
    return 0
  fi
}

__secure_logic_main "${_main_args[@]}"

