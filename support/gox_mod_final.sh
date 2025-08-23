#!/usr/bin/env bash

# gox_mod.sh — modular, robust and "publishable" Go build
# License: MIT (add LICENSE to repo if desired)
# Requirements: bash 4+, Go 1.18+ (recommended), git

set -o errexit
set -o errtrace
set -o functrace
set -o nounset
set -o pipefail
shopt -s inherit_errexit

IFS=$'\n\t'

# ====== Defaults (override via ENV or flags) ======
: "${DIST_DIR:=dist}"
: "${MAKE_TARGET:=build-dev}"
: "${ENABLE_UPX:=0}"                 # 1 to compress with upx
: "${DEFAULT_OS:=}"                  # e.g.: linux
: "${DEFAULT_ARCH:=}"                # e.g.: amd64
: "${BUILD_TAGS:=}"                  # e.g.: netgo,osusergo
: "${ENABLE_RACE:=0}"                # -race on host
: "${VERBOSE:=0}"                    # -v, --verbose
: "${DEBUG:=0}"                      # -x
: "${LD_EXTRA:=}"                    # user extra ldflags

# If BUILDINFO_PATH='main' (or 'mod/internal/buildinfo'),
# injects -X <path>.Version/.Commit/.BuildDate with git info.
: "${BUILDINFO_PATH:=}"

# ====== Utils ======
# Color codes for logs
_SUCCESS="\033[0;32m"
_WARN="\033[0;33m"
_ERROR="\033[0;31m"
_INFO="\033[0;36m"
_NOTICE="\033[0;35m"
_NC="\033[0m"

clear_screen() {
  printf "\033[H\033[2J"
}
now_ms() {
  if command -v gdate >/dev/null 2>&1; then
    # coreutils (brew install coreutils)
    gdate +%s%3N
  elif command -v python3 >/dev/null 2>&1; then
    python3 - <<'PY'
import time
print(int(time.time()*1000))
PY
  elif command -v perl >/dev/null 2>&1; then
    perl -MTime::HiRes=time -e 'printf("%d\n", time()*1000)'
  elif command -v ruby >/dev/null 2>&1; then
    ruby -e 'puts (Time.now.to_f*1000).to_i'
  elif command -v node >/dev/null 2>&1; then
    node -e 'console.log(Date.now())'
  else
    # fallback aproximado (segundos -> ms)
    echo "$(( $(date +%s) * 1000 ))"
  fi
}
log() {
  local type="${1:-}"
  local message="${2:-}"
  local verbose="${3:-${VERBOSE:-0}}"
  verbose="${verbose:-${DEBUG:-0}}"

  case $type in
    question|_QUESTION|-q|-Q)
      if (( verbose )) || (( DEBUG )); then
        printf '%b[QUESTION]%b - %s - ❓  %s: ' "$_NOTICE" "$_NC" "$(date +%H:%M:%S)" "$message" >&2
      fi
      ;;
    notice|_NOTICE|-n|-N)
      if (( verbose )) || (( DEBUG )); then
        printf '%b[NOTICE]%b - %s - 📝  %s\n' "$_NOTICE" "$_NC" "$(date +%H:%M:%S)" "$message" >&2
      fi
      ;;
    info|_INFO|-i|-I)
      if (( verbose )) || (( DEBUG )); then
        printf '%b[INFO]%b - %s - ℹ️  %s\n' "$_INFO" "$_NC" "$(date +%H:%M:%S)" "$message" >&2
      fi
      ;;
    warn|_WARN|-w|-W)
      if (( verbose )) || (( DEBUG )); then
        printf '%b[WARN]%b - %s - ⚠️  %s\n' "$_WARN" "$_NC" "$(date +%H:%M:%S)" "$message" >&2
      fi
      ;;
    error|_ERROR|-e|-E)
      printf '%b[ERROR]%b - %s - ❌  %s\n' "$_ERROR" "$_NC" "$(date +%H:%M:%S)" "$message" >&2
      ;;
    success|_SUCCESS|-s|-S)
      printf '%b[SUCCESS]%b - %s - ✅  %s\n' "$_SUCCESS" "$_NC" "$(date +%H:%M:%S)" "$message" >&2
      ;;
    fatal|_FATAL|-f|-F)
      printf '%b[FATAL]%b 💀 - %s - %s\n' "$_FATAL" "$_NC" "$(date +%H:%M:%S)" "$message" >&2
      if (( verbose )) || (( DEBUG )); then
        printf '%b[FATAL]%b 💀  %s\n' "$_FATAL" "$_NC" "Exiting due to fatal error." >&2
      fi
      # clear_build_artifacts
      exit 1
      ;;
    *)
      if (( verbose )) || (( DEBUG )); then
        log "info" "$message" "$verbose"
      fi
      ;;
  esac
  return 0
}
die()  { log fatal "$*"; }
have() { command -v "${1:-}" >/dev/null 2>&1; }

usage() {

  cat <<'EOF'
Usage: ./gox_mod.sh [options] [PATH]

Options:
--all                   Build for default matrix (linux,darwin,windows x amd64,arm64)
--os OS1,OS2            List of GOOS (e.g.: linux,darwin)
--arch A1,A2            List of GOARCH (e.g.: amd64,arm64)
--tags TAGS             Build tags (e.g.: 'netgo,osusergo')
--race                  Enables -race on native build
--upx                   Compress with UPX (if installed)
-d, --debug             Debug mode (set -x)
-v, --verbose           Verbose (-v)
-h, --help              Help

Environment:
DIST_DIR, MAKE_TARGET, ENABLE_UPX, DEFAULT_OS, DEFAULT_ARCH,
BUILD_TAGS, ENABLE_RACE, VERBOSE, LD_EXTRA, BUILDINFO_PATH

Examples:
./gox_mod.sh
./gox_mod.sh --all
./gox_mod.sh --os linux --arch arm64 ./cmd/myapp
BUILDINFO_PATH=main LD_EXTRA="-s -w" ./gox_mod.sh
EOF

}

_START_TIME="$(now_ms)"

# ====== Args ======
ALL=0
declare -A GOPLT_MAP=()
declare -a ARG_GOOS_LIST=()
declare -a ARG_GOARCH_LIST=()
TARGET_PATH=""

# shellcheck disable=SC2015
parse_args() {
  while (("$#")); do
    case "${1:-}" in
      --all) ALL=1 ;;
      --os) shift; IFS=',' read -r -a ARG_GOOS_LIST <<< "${1:-}";;
      --arch) shift; IFS=',' read -r -a ARG_GOARCH_LIST <<< "${1:-}";;
      --tags) shift; BUILD_TAGS="${1:-}";;
      --race) ENABLE_RACE=1 ;;
      --upx) ENABLE_UPX=1 ;;
      -d|--debug) DEBUG=1;;
      -v|--verbose) VERBOSE=1 ;;
      -h|--help) usage; exit 0 ;;
      *) TARGET_PATH="${1:-}";;
    esac
    shift || true
  done
  (( DEBUG )) && set -x || true
  (( VERBOSE )) && set -v || true
}

# ====== Project / module / git ======
PROJECT_ROOT=""
MOD_NAME=""
GIT_TAG=""
GIT_COMMIT=""
GIT_DATE=""

detect_project_root() {
  PROJECT_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
  cd "${PROJECT_ROOT}"
}

read_module_name() {
  if [[ -f go.mod ]]; then
    MOD_NAME="$(awk '/^module /{print $2}' go.mod | awk -F'/' '{print $NF}')"
  fi
  [[ -n "${MOD_NAME}" ]] || MOD_NAME="$(basename "${PROJECT_ROOT}")"
}

git_info() {
  if have git && git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    GIT_TAG="$(git describe --tags --dirty --always 2>/dev/null || true)"
    GIT_COMMIT="$(git rev-parse --short HEAD 2>/dev/null || true)"
    GIT_DATE="$(git show -s --format=%cd --date=format:%Y-%m-%d 2>/dev/null || date +%Y-%m-%d)"
    SOURCE_DATE_EPOCH="$(git log -1 --pretty=%ct 2>/dev/null || date +%s)"; export SOURCE_DATE_EPOCH
  else
    GIT_TAG="dev"; GIT_COMMIT="none"; GIT_DATE="$(date +%Y-%m-%d)"
    SOURCE_DATE_EPOCH="$(date +%s)"; export SOURCE_DATE_EPOCH
  fi
}

# ====== Discover 'main' packages ======
discover_mains() {
  local path="./..."
  if [[ -n "${TARGET_PATH}" ]]; then
    local rp; rp="$(realpath "${TARGET_PATH}")"
    if [[ -d "${rp}" ]]; then
      path="${rp}/..."
    else
      path="$(dirname "${rp}")/..."
    fi
  fi
  mapfile -t MAIN_DIRS < <(go list -f '{{if eq .Name "main"}}{{.Dir}}{{end}}' "${path}" | awk 'NF')
  (( ${#MAIN_DIRS[@]} )) || die "No 'main' package found."
}

# ====== GOOS/GOARCH Matrix ======
compute_matrix() {
  if (( ALL )); then
    GOPLT_MAP[linux]="(amd64)"
    GOPLT_MAP[darwin]="(amd64 arm64)"
    GOPLT_MAP[windows]="(amd64 386)"
    ARG_GOOS_LIST=(linux darwin windows)
    ARG_GOARCH_LIST=(amd64 arm64 386)
    return 0
  fi
  if [[ ${#ARG_GOOS_LIST[@]} -eq 0 ]]; then
    # shellcheck disable=SC2207
    if [[ -n "${DEFAULT_OS}" ]]; then
      ARG_GOOS_LIST=("${DEFAULT_OS}");
      for _GOENVOS in "${ARG_GOOS_LIST[@]}"; do
        GOPLT_MAP["${_GOENVOS}"]="(amd64)";
      done
    else
      ARG_GOOS_LIST=( $(go env GOOS) );
      for os in "${!ARG_GOOS_LIST[@]}"; do
        GOPLT_MAP[${ARG_GOOS_LIST[$os]}]="( $(go env GOARCH || echo amd64) )";
      done
    fi
  fi
  if [[ ${#ARG_GOARCH_LIST[@]} -eq 0 ]]; then
    # shellcheck disable=SC2207
    if [[ -n "${DEFAULT_ARCH}" ]]; then
      ARG_GOARCH_LIST=("${DEFAULT_ARCH}");
      for _GOENVARCH in "${ARG_GOARCH_LIST[@]}"; do
        for os in "${!ARG_GOOS_LIST[@]}"; do
          if [[ -z "${GOPLT_MAP[${ARG_GOOS_LIST[${os:-}]}]}" ]]; then
            GOPLT_MAP["${ARG_GOOS_LIST[${os:-}]}"]="(${_GOENVARCH})";
          else
            if ! [[ " ${GOPLT_MAP[${ARG_GOOS_LIST[${os:-}]}]} " == *" ${_GOENVARCH} "* ]]; then
              GOPLT_MAP["${ARG_GOOS_LIST[${os:-}]}"]+=" ${_GOENVARCH}";
            fi
          fi
        done
      done
    else
      ARG_GOARCH_LIST=( $(go env GOARCH) );
      for _GOENVARCH in "${ARG_GOARCH_LIST[@]}"; do
        for os in "${!ARG_GOOS_LIST[@]}"; do
          if [[ -z "${GOPLT_MAP[${ARG_GOOS_LIST[${os:-}]}]}" ]]; then
            GOPLT_MAP["${ARG_GOOS_LIST[${os:-}]}"]="(${_GOENVARCH})";
          else
            if ! [[ " ${GOPLT_MAP[${ARG_GOOS_LIST[${os:-}]}]} " == *" ${_GOENVARCH} "* ]]; then
              GOPLT_MAP["${ARG_GOOS_LIST[${os:-}]}"]+=" ${_GOENVARCH}";
            fi
          fi
        done
      done
    fi
  fi
}

# ====== Detect build flag support ======
supports_build_flag() {
  # e.g.: supports_build_flag buildvcs
  go help build 2>/dev/null | grep -q -- "-$1\b"
}

# ====== Build argument assembly ======
declare -a BUILD_ARGS=()

build_flags() {
  local os="${1:-}"
  local arch="${2:-}"
  BUILD_ARGS=()  # reset

  if supports_build_flag trimpath; then
    BUILD_ARGS+=(-trimpath)
  fi
  if supports_build_flag buildvcs; then
    BUILD_ARGS+=(-buildvcs)
  fi

  if (( ENABLE_RACE )); then
    local host_os;   host_os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    local host_arch; host_arch="$(uname -m)"
    if [[ "$os" == "$host_os" && "$arch" == "$host_arch" ]]; then
      BUILD_ARGS+=(-race)
    fi
  fi

  if [[ -n "${BUILD_TAGS}" ]]; then
    BUILD_ARGS+=(-tags "${BUILD_TAGS}")
  fi

  local ld_str="-s -w"
  [[ -n "${LD_EXTRA}" ]] && ld_str="${ld_str} ${LD_EXTRA}"
  if [[ -n "${BUILDINFO_PATH}" ]]; then
    ld_str="${ld_str} -X ${BUILDINFO_PATH}.Version=${GIT_TAG}"
    ld_str="${ld_str} -X ${BUILDINFO_PATH}.Commit=${GIT_COMMIT}"
    ld_str="${ld_str} -X ${BUILDINFO_PATH}.BuildDate=${GIT_DATE}"
  fi
  BUILD_ARGS+=(-ldflags "${ld_str}")
}

# ====== UPX ======
maybe_upx() {
  local bin="${1:-}"
  (( ENABLE_UPX )) || return 0
  have upx || { log warn "UPX not found; continuing without compression."; return 0; }
  log notice "UPX: compressing '${bin}'"
  upx "$bin" --force-overwrite --lzma --no-progress --no-color -qqq || log info "UPX failed (ignoring)."
}

# ====== Build ======
build_one() {
  local _dir="${1:-.}"
  local _os="${2:-$(uname -s | tr '[:upper:]' '[:lower:]')}"
  local _arch="${3:-$(uname -m)}"
  local _ext="";
  [[ "${_os}" == "windows" ]] && _ext=".exe"

  local _pkg_name;
  _pkg_name="$(basename "${_dir}")"

  local _bin_name="${MOD_NAME}_${_os}_${_arch}${_ext:-}"
  local _out_dir="${DIST_DIR}" # /${_pkg_name}"

  mkdir -p "${_out_dir}"

  local _out_bin="${_out_dir}/${_bin_name}"

  log info "Building ${_pkg_name} -> ${_out_bin}"
  build_flags "${_os}" "${_arch}"

  GOOS="${_os}" GOARCH="${_arch}" \
  go mod tidy > /dev/null || {
    log warn "go mod tidy failed (ignoring)."
    return 1
  }
  GOOS="${_os}" GOARCH="${_arch}" \
  go build \
  "${BUILD_ARGS[@]}" \
  -o "${_out_bin}" "${_dir}"

  maybe_upx "${_out_bin}"
  log success "OK: ${_out_bin}"
}

build_all() {
  mkdir -p "${DIST_DIR}"

  if [[ -f Makefile ]]; then
    log info "Makefile detected. Running target '${MAKE_TARGET}'..."
    # shellcheck disable=SC1007
    if MAKEFLAGS= FORCE=y make "${MAKE_TARGET}"; then
      log success "Build via Makefile completed."
      return 0
    else
      log error "Target '${MAKE_TARGET}' failed. Fallback to 'go build'."
    fi
  fi

  for __os in "${!GOPLT_MAP[@]}"; do
    local arch_list
    GOPLT_MAP[$__os]=${GOPLT_MAP[${__os}]//\(/}
    GOPLT_MAP[${__os}]=${GOPLT_MAP[${__os}]//\)/}
    IFS=' ' read -r -a arch_list <<< "${GOPLT_MAP[${__os}]}"
    for __arch in "${arch_list[@]}"; do
      for __dir in "${MAIN_DIRS[@]}"; do
        build_one "${__dir}" "${__os}" "${__arch}"
      done
    done
  done
}

# ====== Performance measurement ======
measure_performance() {
  local _EXIT_CODE="$?"
  _EXIT_CODE="${_EXIT_CODE:-${1:-}}"
  _EXIT_CODE="${_EXIT_CODE:-0}"

  # Measure total duration in milliseconds
  _END_TIME="$(now_ms)"
  _DURATION=$(( _END_TIME - _START_TIME ))

  if (( _DURATION < 1000 )); then
    log notice "Total duration: ${_DURATION} ms"
  else
    log notice "Total duration: $(( _DURATION / 60000 )) min $(( (_DURATION / 1000) % 60 )) sec $(( _DURATION % 1000 )) ms"
  fi

  if (( _EXIT_CODE == 0 )); then
    log success "All builds succeeded."
  else
    log error "Builds completed with errors (exit code: ${_EXIT_CODE})."
  fi
}

# ====== Main ======
main() {
  parse_args "$@" || die "Error parsing arguments."
  have go || die "Go not found in PATH."
  detect_project_root || die "Could not detect project directory."
  read_module_name || die "Could not read module name."
  git_info || true
  discover_mains || die "Could not discover 'main' packages."
  compute_matrix || die "Error computing GOOS/GOARCH matrix."
  build_all || die "Build failed."
  measure_performance $? || true
  return 0
}

main "$@"
