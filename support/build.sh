#!/usr/bin/env bash

set -euo pipefail
set -o errtrace
set -o functrace
set -o posix

IFS=$'\n\t'

# ====== Configuration from manifest.json ======
load_manifest_config() {
  local _root_dir="${_ROOT_DIR:-${ROOT_DIR:-$(git rev-parse --show-toplevel)}}"
  local _manifest_path="${_root_dir}/internal/module/info/manifest.json"

  if [[ -f "${_manifest_path}" ]]; then
    export _APP_NAME="${_APP_NAME:-$(jq -r '.application // .bin // .name' "${_manifest_path}" 2>/dev/null || echo "")}"
    export _VERSION="${_VERSION:-$(jq -r '.version' "${_manifest_path}" 2>/dev/null || echo "")}"
    export _ORGANIZATION="${_ORGANIZATION:-$(jq -r '.organization' "${_manifest_path}" 2>/dev/null || echo "")}"

    # Load supported platforms from manifest
    local _platforms_json
    _platforms_json=$(jq -r '.platforms[]?' "${_manifest_path}" 2>/dev/null || echo "")
    if [[ -n "${_platforms_json}" ]]; then
      export _SUPPORTED_PLATFORMS="${_platforms_json}"
    fi

    log info "Loaded configuration from manifest: ${_APP_NAME} v${_VERSION}"
  else
    log warn "Manifest not found at ${_manifest_path}, using defaults"
  fi
}

# ====== Git information ======
load_git_info() {
  if command -v git >/dev/null 2>&1 && git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    export _GIT_TAG="${_GIT_TAG:-$(git describe --tags --dirty --always 2>/dev/null || echo "dev")}"
    export _GIT_COMMIT="${_GIT_COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")}"
    export _GIT_DATE="${_GIT_DATE:-$(git show -s --format=%cd --date=format:%Y-%m-%d 2>/dev/null || date +%Y-%m-%d)}"
    export SOURCE_DATE_EPOCH="${SOURCE_DATE_EPOCH:-$(git log -1 --pretty=%ct 2>/dev/null || date +%s)}"

    log info "Git info: ${_GIT_TAG} (${_GIT_COMMIT}) ${_GIT_DATE}"
  else
    export _GIT_TAG="dev"
    export _GIT_COMMIT="unknown"
    export _GIT_DATE=
    _GIT_DATE="$(date +%Y-%m-%d)"
    export SOURCE_DATE_EPOCH=
    SOURCE_DATE_EPOCH="$(date +%s)"

    log warn "Git not available, using default build info"
  fi
}

# ====== Go Version Validation ======
# Source go version management functions
# shellcheck disable=SC1091
source "$(dirname "${BASH_SOURCE[0]}")/go_version.sh"

validate_go_version() {
  # Use modular go version checking
  if ! check_go_version_compatibility; then
    if [[ "${GO_VERSION_CHECK:-true}" == "true" ]]; then
      log error "Set GO_VERSION_CHECK=false to skip this check"
      return 1
    fi
  fi
  return 0
}

# ====== Load platforms from manifest.json ======
load_platforms_from_manifest() {
  local _root_dir="${_ROOT_DIR:-${ROOT_DIR:-$(git rev-parse --show-toplevel)}}"
  local _manifest_path="${_root_dir}/internal/module/info/manifest.json"

  if [[ -f "${_manifest_path}" ]] && command -v jq >/dev/null 2>&1; then
    local platforms=()
    while IFS= read -r platform; do
      if [[ -n "${platform}" && "${platform}" != "null" ]]; then
        platforms+=("${platform}")
      fi
    done < <(jq -r '.platforms[]?' "${_manifest_path}" 2>/dev/null)

    if [[ ${#platforms[@]} -gt 0 ]]; then
      printf '%s\n' "${platforms[@]}"
      return 0
    fi
  fi

  # Fallback to default platforms
  echo "linux/amd64"
  echo "linux/arm64"
  echo "linux/armv6l"
  echo "linux/386"
  echo "darwin/amd64"
  echo "darwin/arm64"
  echo "windows/amd64"
  echo "windows/386"

  return 1
}

# ====== Discover main packages ======
discover_main_packages() {
  local _root_dir="${_ROOT_DIR:-${ROOT_DIR:-$(git rev-parse --show-toplevel)}}"
  local _target_path="${1:-./cmd/...}"

  # Use relative path from project root
  cd "${_root_dir}" || return 1

  local _main_dirs_arr=()
  mapfile -t _main_dirs_arr < <(go list -f '{{if eq .Name "main"}}{{.Dir}}{{end}}' "${_target_path}" 2>/dev/null | awk 'NF')

  if [[ ${#_main_dirs_arr[@]} -eq 0 ]]; then
    log error "No 'main' packages found in ${_target_path}"
    return 1
  fi

  log info "Found ${#_main_dirs_arr[@]} main package(s): ${_main_dirs_arr[*]}"
  printf '%s\n' "${_main_dirs_arr[@]}"
}

# ====== Platform/Architecture Matrix ======
compute_build_matrix() {
  local _platform_arg="${1:-$(uname -s | tr '[:upper:]' '[:lower:]')}"
  local _arch_arg="${2:-$(uname -m | tr '[:upper:]' '[:lower:]')}"

  # Normalize architecture
  case "${_arch_arg}" in
    x86_64|amd64) _arch_arg="amd64" ;;
    armv8|aarch64|arm64) _arch_arg="arm64" ;;
    armv6l|armv7l) _arch_arg="armv6l" ;;
    i386|I386) _arch_arg="386" ;;
  esac

  # Declare associative array for platform-architecture mapping
  declare -A GOPLT_MAP
  declare -a PLATFORMS=()

  case "${_platform_arg}" in
    all|ALL)
      PLATFORMS=("linux" "darwin" "windows")
      GOPLT_MAP[linux]="amd64 arm64 armv6l 386"
      GOPLT_MAP[darwin]="amd64 arm64"
      GOPLT_MAP[windows]="amd64 386"
      ;;
    linux|LINUX)
      PLATFORMS=("linux")
      case "${_arch_arg}" in
        all|ALL) GOPLT_MAP[linux]="amd64 arm64 armv6l 386" ;;
        *) GOPLT_MAP[linux]="${_arch_arg}" ;;
      esac
      ;;
    darwin|DARWIN|macOS|MACOS)
      PLATFORMS=("darwin")
      case "${_arch_arg}" in
        all|ALL) GOPLT_MAP[darwin]="amd64 arm64" ;;
        *) GOPLT_MAP[darwin]="${_arch_arg}" ;;
      esac
      ;;
    windows|WINDOWS)
      PLATFORMS=("windows")
      case "${_arch_arg}" in
        all|ALL) GOPLT_MAP[windows]="amd64 386" ;;
        *) GOPLT_MAP[windows]="${_arch_arg}" ;;
      esac
      ;;
    *)
      # Single platform specified
      PLATFORMS=("${_platform_arg}")
      GOPLT_MAP[${_platform_arg}]="${_arch_arg}"
      ;;
  esac

  # Export for use in other functions
  for platform in "${PLATFORMS[@]}"; do
    echo "${platform}:${GOPLT_MAP[${platform}]}"
  done
}

# ====== Build flags and ldflags ======
prepare_build_flags() {
  local _os="${1:-}"
  local _arch="${2:-}"
  local _build_mode="${3:-production}" # production or development

  local _build_args=()
  local _ldflags="-s -w"

  # Add trimpath if supported
  if go help build 2>/dev/null | grep -q -- "-trimpath"; then
    _build_args+=("-trimpath")
  fi

  # Add buildvcs if supported (Go 1.18+)
  if go help build 2>/dev/null | grep -q -- "-buildvcs"; then
    _build_args+=("-buildvcs=false")
  fi

  # Add build tags if specified
  if [[ -n "${BUILD_TAGS:-}" ]]; then
    _build_args+=("-tags" "${BUILD_TAGS}")
  fi

  # Race detection for same platform builds
  if [[ "${ENABLE_RACE:-0}" == "1" ]]; then
    local host_os host_arch
    host_os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    host_arch="$(uname -m)"
    [[ "${host_arch}" == "x86_64" ]] && host_arch="amd64"
    [[ "${host_arch}" == "aarch64" ]] && host_arch="arm64"
    [[ "${host_arch}" == "armv6l" ]] && host_arch="armv6l"
    [[ "${host_arch}" == "i386" ]] && host_arch="386"

    if [[ "${_os}" == "${host_os}" && "${_arch}" == "${host_arch}" ]]; then
      _build_args+=("-race")
      log info "Enabling race detection for native build"
    fi
  fi

  # Build info injection
  if [[ -n "${_GIT_TAG:-}" && -n "${_GIT_COMMIT:-}" && -n "${_GIT_DATE:-}" ]]; then
    _ldflags="${_ldflags} -X main.version=${_GIT_TAG}"
    _ldflags="${_ldflags} -X main.commit=${_GIT_COMMIT}"
    _ldflags="${_ldflags} -X main.date=${_GIT_DATE}"
  fi

  # Extra ldflags
  if [[ -n "${LD_EXTRA:-}" ]]; then
    _ldflags="${_ldflags} ${LD_EXTRA}"
  fi

  # Development mode adjustments
  if [[ "${_build_mode}" == "development" ]]; then
    # Keep debug info in development
    _ldflags="${_ldflags/-s /}"
    _ldflags="${_ldflags/-w /}"
  fi

  _build_args+=("-ldflags" "${_ldflags}")

  printf '%s\n' "${_build_args[@]}"
}

get_output_name() {
  local _platform_pos="${1:-${_PLATFORM:-$(uname -s | tr '[:upper:]' '[:lower:]')}}"
  local _arch_pos="${2:-${_ARCH:-$(uname -m | tr '[:upper:]' '[:lower:]')}}"
  local _root_dir="${_ROOT_DIR:-${ROOT_DIR:-$(git rev-parse --show-toplevel)}}"
  local _binary_dir="${_root_dir:-}/dist"
  local _app_name="${_APP_NAME:-$(basename "${_root_dir}")}"

  [[ ! -d "${_binary_dir:-}" ]] && mkdir -p "${_binary_dir:-}"

  local _output_name="${_binary_dir}/${_app_name}_${_platform_pos}_${_arch_pos}"
  [[ "${_platform_pos}" == "windows" ]] && _output_name="${_output_name}.exe"

  echo "${_output_name}"
}

# ====== UPX compression ======
upx_packaging() {
  local _output_name="${1:-}"
  local _platform_pos="${2:-}"
  local _will_upx_pack="${3:-${_WILL_UPX_PACK_BINARY:-true}}"

  [[ "${_will_upx_pack}" != "true" ]] && return 0
  [[ ! -f "${_output_name}" ]] && return 1

  if ! command -v upx >/dev/null 2>&1; then
    log warn "UPX not found, skipping compression"
    return 0
  fi

  log info "Compressing binary with UPX: $(basename "${_output_name}")"

  local _upx_args=("--force-overwrite" "--lzma" "--no-progress" "--no-color" "-qqq")

  # macOS specific UPX handling
  if [[ "${_platform_pos}" == "darwin" ]]; then
    _upx_args+=("--force-macos")
    log warn "UPX on macOS may have limited compatibility"
  fi

  if upx "${_upx_args[@]}" "${_output_name}" 2>/dev/null; then
    log success "UPX compression successful: $(basename "${_output_name}")"
  else
    log warn "UPX compression failed for $(basename "${_output_name}"), continuing..."
  fi

  return 0
}

# ====== Single binary compilation ======
compile_binary() {
  local _platform_pos="${1:-}"
  local _arch_pos="${2:-}"
  local _main_dir="${3:-}"
  local _build_mode="${4:-production}"

  [[ -z "${_platform_pos}" || -z "${_arch_pos}" || -z "${_main_dir}" ]] && {
    log error "Missing required parameters for compile_binary"
    return 1
  }

  local _output_name
  _output_name=$(get_output_name "${_platform_pos}" "${_arch_pos}")

  log info "Building for ${_platform_pos}/${_arch_pos}: $(basename "${_main_dir}") -> $(basename "${_output_name}")"

  # Prepare build arguments
  local _build_args=()
  mapfile -t _build_args < <(prepare_build_flags "${_platform_pos}" "${_arch_pos}" "${_build_mode}")

  # Set build environment
  local _build_env=(
    "GOOS=${_platform_pos}"
    "GOARCH=${_arch_pos}"
    "CGO_ENABLED=0"
  )

  # Execute build
  if env "${_build_env[@]}" go build "${_build_args[@]}" -o "${_output_name}" "./cmd"; then
    log success "Build successful: $(basename "${_output_name}")"

    # Apply UPX compression if enabled
    if [[ "${_build_mode}" == "production" ]]; then
      upx_packaging "${_output_name}" "${_platform_pos}" "${_WILL_UPX_PACK_BINARY:-true}"
    fi

    return 0
  else
    log error "Build failed for ${_platform_pos}/${_arch_pos}"
    return 1
  fi
}

# ====== Build for single architecture ======
build_for_arch() {
  local _platform_pos="${1:-}"
  local _arch_pos="${2:-}"
  local _main_dirs="${3:-}"
  local _build_mode="${4:-production}"
  local _force="${5:-false}"

  # Normalize architecture
  case "${_arch_pos}" in
    x86_64|X86_64) _arch_pos="amd64" ;;
    armv8|aarch64|AARCH64) _arch_pos="arm64" ;;
    armv6l|ARMV6L) _arch_pos="armv6l" ;;
    i386|I386) _arch_pos="386" ;;
  esac

  # Validate platform/arch combination
  if ! is_valid_platform_arch "${_platform_pos}" "${_arch_pos}"; then
    log warn "Skipping unsupported combination: ${_platform_pos}/${_arch_pos}"
    return 0
  fi

  # Build each main package
  local main_dir
  local build_count=0
  while IFS= read -r main_dir; do
    [[ -z "${main_dir}" ]] && continue
    build_count=$((build_count + 1))
    log info "BUILD ITERATION #${build_count} for ${_platform_pos}/${_arch_pos}: ${main_dir}"

    local _output_name
    _output_name=$(get_output_name "${_platform_pos}" "${_arch_pos}")

    # Check if we should overwrite existing binary
    if ! check_overwrite_binary "${_output_name}" "${_force}" "true"; then
      continue
    fi

    # Compile the binary
    if compile_binary "${_platform_pos}" "${_arch_pos}" "${main_dir}" "${_build_mode}"; then
      # Create compressed archive if in production mode
      if [[ "${_build_mode}" == "production" ]]; then
        create_archive "${_platform_pos}" "${_arch_pos}" "${_output_name}"
      fi
    fi
  done <<< "${_main_dirs}"

  return 0
}

# ====== Validate platform/architecture combination ======
is_valid_platform_arch() {
  local _platform="${1:-}"
  local _arch="${2:-}"

  case "${_platform}" in
    linux)
      case "${_arch}" in
        amd64|arm64|armv6l|386) return 0 ;;
        *) return 1 ;;
      esac
      ;;
    darwin)
      case "${_arch}" in
        amd64|arm64) return 0 ;;
        *) return 1 ;;
      esac
      ;;
    windows)
      case "${_arch}" in
        amd64|386) return 0 ;;
        *) return 1 ;;
      esac
      ;;
    *)
      return 1
      ;;
  esac
}

# ====== Create compressed archive ======
create_archive() {
  local _platform="${1:-}"
  local _arch="${2:-}"
  local _binary_path="${3:-}"

  [[ ! -f "${_binary_path}" ]] && return 1

  local _binary_name
  _binary_name=$(basename "${_binary_path}")
  local _binary_dir
  _binary_dir=$(dirname "${_binary_path}")

  cd "${_binary_dir}" || return 1

  if [[ "${_platform}" == "windows" ]]; then
    local _archive_name="${_binary_name%.exe}.zip"
    if command -v zip >/dev/null 2>&1; then
      if zip -9 "${_archive_name}" "${_binary_name}" >/dev/null 2>&1; then
        log success "Created archive: ${_archive_name}"
      else
        log warn "Failed to create ZIP archive for ${_binary_name}"
      fi
    fi
  else
    local _archive_name="${_binary_name}.tar.gz"
    if tar -czf "${_archive_name}" "${_binary_name}" 2>/dev/null; then
      log success "Created archive: ${_archive_name}"
    else
      log warn "Failed to create tar.gz archive for ${_binary_name}"
    fi
  fi

  cd - >/dev/null || true
}

# ====== Build for all architectures of a platform ======
build_for_platform() {
  local _platform_pos="${1:-}"
  local _arch_list="${2:-}"
  local _main_dirs="${3:-}"
  local _build_mode="${4:-production}"
  local _force="${5:-false}"

  log info "Building for platform: ${_platform_pos}"

  local arch
  for arch in ${_arch_list}; do
    build_for_arch "${_platform_pos}" "${arch}" "${_main_dirs}" "${_build_mode}" "${_force}"
  done
}

# ====== Main build function ======
build_binary() {
  local _platform_args="${1:-}"
  local _arch_args="${2:-}"
  local _force="${3:-false}"
  local _build_mode="${4:-production}"

  # Initialize configuration
  load_manifest_config
  load_git_info
  validate_go_version || return 1

  # Debug logging
  log info "Build called with: platform='${_platform_args}', arch='${_arch_args}', force='${_force}', mode='${_build_mode}'"

  # Check for cross-compilation mode
  # Only enter cross-compilation if no specific platform is provided
  if [[ "${_platform_args:-all}" == "all" ]] || [[ "${_platform_args:-__CROSS_COMPILE__}" == "__CROSS_COMPILE__" ]] || [[ -z "${_platform_args}" && -z "${_arch_args}" ]]; then
    log info "Cross-compilation mode: building for all platforms in manifest.json"

    # Discover main packages once
    local _main_dirs
    _main_dirs=$(discover_main_packages "./cmd/...")
    [[ -z "${_main_dirs}" ]] && {
      log error "No main packages found"
      return 1
    }

    # log info "Starting cross-platform build process..."
    # log notice "App: ${_APP_NAME:-unknown} v${_VERSION:-unknown}"
    # log notice "Git: ${_GIT_TAG:-unknown} (${_GIT_COMMIT:-unknown})"
    # log notice "Mode: ${_build_mode}"

    # Load platforms from manifest and build directly
    local built_any=false
    while IFS='/' read -r os arch; do
      if [[ -n "${os}" && -n "${arch}" ]]; then
        log info "Building for ${os}/${arch}..."
        if build_for_arch "${os}" "${arch}" "${_main_dirs}" "${_build_mode}" "${_force}"; then
          built_any=true
        fi
      fi
    done < <(load_platforms_from_manifest)

    if [[ "${built_any}" == "true" ]]; then
      log success "Cross-platform build completed"
      return 0
    else
      log error "No platforms were successfully built"
      return 1
    fi
  else
    # Single-platform build mode
    log info "Single-platform build mode: ${_platform_args}/${_arch_args}"

    # Normalize arch for single builds
    case "${_arch_args}" in
      x86_64|X86_64) _arch_args="amd64" ;;
      armv8|aarch64|AARCH64) _arch_args="arm64" ;;
      armv6l|ARMV6L) _arch_args="armv6l" ;;
      i386|I386) _arch_args="386" ;;
    esac

    # Discover main packages
    local _main_dirs
    _main_dirs=$(discover_main_packages "./cmd/...")
    [[ -z "${_main_dirs}" ]] && {
      log error "No main packages found"
      return 1
    }

    log info "Starting single-platform build process..."
    log notice "App: ${_APP_NAME:-unknown} v${_VERSION:-unknown}"
    log notice "Git: ${_GIT_TAG:-unknown} (${_GIT_COMMIT:-unknown})"
    log notice "Mode: ${_build_mode}"
    log notice "Target: ${_platform_args}/${_arch_args}"

    # Build directly for the specified platform/arch
    if build_for_arch "${_platform_args}" "${_arch_args}" "${_main_dirs}" "${_build_mode}" "${_force}"; then
      log success "Single-platform build completed: ${_platform_args}/${_arch_args}"
      return 0
    else
      log error "Single-platform build failed: ${_platform_args}/${_arch_args}"
      return 1
    fi
  fi
}

check_overwrite_binary() {
  local _platform_pos="${1:-${_platform_pos:-${_PLATFORM:-$(uname -s | tr '[:upper:]' '[:lower:]')}}}"
  local _arch_pos="${2:-${_arch_pos:-${_ARCH:-$(uname -m | tr '[:upper:]' '[:lower:]')}}}"
  local _output_name="${3:-${_output_name:-}}"
  local _force="${4:-${_force:-n}}"
  local _will_upx_pack_binary="${5:-${_will_upx_pack_binary:-${WILL_UPX_PACK_BINARY:-true}}}"
  local _is_interactive="${6:-${_is_interactive:-${IS_INTERACTIVE:-}}}"

  if [[ -f "${_output_name:-}" ]]; then
    local REPLY="y"

    if [[ "${_is_interactive:-}" != "true" || "${CI:-}" != "true" ]]; then

      if [[ ${_force:-} =~ [yY] || ${_force:-} == "true" || "${NON_INTERACTIVE:-}" == "true" ]]; then
        REPLY="y"
      elif [[ -t 0 ]]; then
        # If the script is running interactively, prompt for confirmation
        log notice "Binary already exists: ${_output_name:-}" true
        log notice "Current binary: ${_output_name:-}"
        log notice "Press 'y' to overwrite or any other key to skip." true
        log question "(y) to overwrite, any other key to skip (default: n, 10 seconds to respond)" true
        read -t 10 -p "" -n 1 -r REPLY || REPLY="n"
        echo '' # Move to a new line after the prompt
        REPLY="${REPLY,,}"  # Convert to lowercase
        REPLY="${REPLY:-n}"  # Default to 'n' if no input
      else
        log notice "Binary already exists: ${_output_name:-}" true
        log notice "Skipping confirmation in non-interactive mode." true
      fi
    fi

    if [[ ! ${REPLY:-} =~ [yY] ]]; then
      log notice "Skipping build for ${_platform_pos:-} ${_arch_pos:-}." true
      return 0
    fi

    log warn "Overwriting existing binary: ${_output_name}" true
    if [[ "${_platform_pos:-}" == "windows" ]]; then
      rm -f "${_output_name:-}.exe" || return 1
    else
      rm -f "${_output_name:-}" || return 1
    fi

    log info "Binary built successfully: ${_output_name}"

    if compile_binary "${_platform_pos:-}" "${_arch_pos:-}" "${_output_name:-}" "${_force:-}" "${_will_upx_pack_binary:-}"; then
      log success "Binary built successfully: ${_output_name}"
      return 0
    else
      log error "Failed to build binary: ${_output_name}" true
      return 1
    fi
  fi

  return 0
}

arch_iterator() {
  # Legacy compatibility wrapper
  local _platform_pos="${1:-}"
  local _arch_pos="${2:-}"
  local _build_mode="${3:-production}"
  local _force="${4:-false}"

  local _main_dirs
  _main_dirs=$(discover_main_packages "./cmd/...")

  build_for_arch "${_platform_pos}" "${_arch_pos}" "${_main_dirs}" "${_build_mode}" "${_force}"
}

platform_iterator() {
  # Legacy compatibility wrapper
  local _platform_pos="${1:-}"
  local _arch_args="${2:-}"
  local _force="${3:-false}"
  local _build_mode="${4:-production}"

  local _main_dirs
  _main_dirs=$(discover_main_packages "./cmd/...")

  build_for_platform "${_platform_pos}" "${_arch_args}" "${_main_dirs}" "${_build_mode}" "${_force}"
}

compress_binary() {
  local _platform_arg="${1:-}"
  local _arch_arg="${2:-}"
  local _output_name="${3:-}"

  # Wrapper around create_archive for legacy compatibility
  if [[ -n "${_output_name}" && -f "${_output_name}" ]]; then
    create_archive "${_platform_arg}" "${_arch_arg}" "${_output_name}"
  else
    log warn "compress_binary called but no valid binary found: ${_output_name:-<none>}"
  fi
}

# ====== Platform/Architecture utility functions ======
_get_os_arr_from_args() {
  local _platform="${1:-$(uname -s | tr '[:upper:]' '[:lower:]')}"

  case "${_platform}" in
    all|ALL|a|A|-a|-A)
      echo "windows darwin linux"
      ;;
    win|WIN|windows|WINDOWS|w|W|-w|-W)
      echo "windows"
      ;;
    linux|LINUX|l|L|-l|-L)
      echo "linux"
      ;;
    darwin|DARWIN|macOS|MACOS|m|M|-m|-M)
      echo "darwin"
      ;;
    *)
      echo "${_platform}"
      ;;
  esac
}

_get_arch_arr_from_args() {
  local _platform="${1:-$(uname -s | tr '[:upper:]' '[:lower:]')}"
  local _arch="${2:-$(uname -m | tr '[:upper:]' '[:lower:]')}"

  # Normalize architecture names
  case "${_arch}" in
    x86_64|X86_64) _arch="amd64" ;;
    armv8|aarch64|AARCH64) _arch="arm64" ;;
    armv6l|ARMV6L) _arch="armv6l" ;;
    i386|I386) _arch="386" ;;
  esac

  case "${_platform}" in
    darwin|DARWIN|macOS|MACOS)
      case "${_arch}" in
        all|ALL|a|A|-a|-A)
          echo "amd64 arm64"
          ;;
        amd64|arm64)
          echo "${_arch}"
          ;;
        *)
          log error "Invalid architecture '${_arch}' for darwin. Valid: amd64, arm64"
          return 1
          ;;
      esac
      ;;
    linux|LINUX)
      case "${_arch}" in
        all|ALL|a|A|-a|-A)
          echo "amd64 arm64 armv6l 386"
          ;;
        amd64|arm64|armv6l|386)
          echo "${_arch}"
          ;;
        *)
          log error "Invalid architecture '${_arch}' for linux. Valid: amd64, arm64, armv6l, 386"
          return 1
          ;;
      esac
      ;;
    windows|WINDOWS)
      case "${_arch}" in
        all|ALL|a|A|-a|-A)
          echo "amd64 386"
          ;;
        amd64|386)
          echo "${_arch}"
          ;;
        *)
          log error "Invalid architecture '${_arch}' for windows. Valid: amd64, 386"
          return 1
          ;;
      esac
      ;;
    *)
      echo "${_arch}"
      ;;
  esac
}

_get_os_from_args() {
  local _platform="${1:-"$(uname -s | tr '[:upper:]' '[:lower:]')"}"

  case "${_platform:-"$(uname -s | tr '[:upper:]' '[:lower:]')"}" in
    all|ALL|a|A|-a|-A)
      echo "all"
    ;;

    win|WIN|windows|WINDOWS|w|W|-w|-W)
      echo "windows"
    ;;

    linux|LINUX|l|L|-l|-L)
      echo "linux"
    ;;

    darwin|DARWIN|macOS|MACOS|m|M|-m|-M)
      echo "darwin"
    ;;

    *)
      log fatal "Invalid platform: '${_platform:-}'. Valid options: windows, linux, darwin, all."
    ;;

  esac
}

_get_arch_from_args() {
  local _platform="${1:-$(uname -s | tr '[:upper:]' '[:lower:]')}"
  local _arch="${2:-$(uname -m | tr '[:upper:]' '[:lower:]')}"

  # Normalize common arch names
  case "${_arch}" in
    x86_64|X86_64) echo "amd64" ;;
    armv8|aarch64|AARCH64) echo "arm64" ;;
    i386|I386) echo "386" ;;
    ARMV6L|aa) echo "armv6l" ;;
    all|ALL|a|A|-a|-A) echo "all" ;;
    amd64|arm64|386|armv6l) echo "${_arch}" ;;
    *) echo "${_arch}" ;;
  esac
}

# ====== Enhanced build_binary with legacy support ======
build_binary_enhanced() {
  local _platform_args="${1:-$(uname -s | tr '[:upper:]' '[:lower:]')}"
  local _arch_args="${2:-$(uname -m | tr '[:upper:]' '[:lower:]')}"
  local _force="${3:-false}"
  local _build_mode="${4:-production}"

  # Set build mode based on UPX setting for backward compatibility
  if [[ "${_WILL_UPX_PACK_BINARY:-true}" == "false" ]]; then
    _build_mode="development"
  fi

  build_binary "${_platform_args}" "${_arch_args}" "${_force}" "${_build_mode}"
}

# ====== Performance measurement ======
measure_build_performance() {
  local _start_time="${1:-}"
  local _end_time="${2:-$(date +%s%3N 2>/dev/null || echo $(($(date +%s) * 1000)))}"

  if [[ -n "${_start_time}" ]]; then
    local _duration=$(( _end_time - _start_time ))

    if (( _duration < 1000 )); then
      log notice "Build completed in ${_duration}ms"
    elif (( _duration < 60000 )); then
      log notice "Build completed in $((_duration / 1000)).$(((_duration % 1000) / 100))s"
    else
      local minutes=$(( _duration / 60000 ))
      local seconds=$(( (_duration / 1000) % 60 ))
      log notice "Build completed in ${minutes}m ${seconds}s"
    fi
  fi
}

# ====== Clean build artifacts ======
clean_build_artifacts() {
  local _root_dir="${_ROOT_DIR:-${ROOT_DIR:-$(git rev-parse --show-toplevel)}}"
  local _dist_dir="${_root_dir}/dist"

  if [[ -d "${_dist_dir}" ]]; then
    log info "Cleaning build artifacts in ${_dist_dir}"
    rm -rf "${_dist_dir:?}"/*
    log success "Build artifacts cleaned"
  else
    log info "No build artifacts to clean"
  fi
}

# ====== Install UPX if needed ======
install_upx() {
  command -v upx >/dev/null 2>&1 && return 0

  log info "UPX not found, attempting to install..."

  case "$(uname -s)" in
    Linux*)
      if command -v apt-get >/dev/null 2>&1; then
        sudo apt-get update && sudo apt-get install -y upx-ucl
      elif command -v yum >/dev/null 2>&1; then
        sudo yum install -y upx
      elif command -v pacman >/dev/null 2>&1; then
        sudo pacman -S --noconfirm upx
      else
        log warn "Could not install UPX automatically on this Linux distribution"
        return 1
      fi
      ;;
    Darwin*)
      if command -v brew >/dev/null 2>&1; then
        brew install upx
      else
        log warn "Homebrew not found, cannot install UPX automatically"
        return 1
      fi
      ;;
    *)
      log warn "Automatic UPX installation not supported on this platform"
      return 1
      ;;
  esac

  if command -v upx >/dev/null 2>&1; then
    log success "UPX installed successfully"
    return 0
  else
    log error "UPX installation failed"
    return 1
  fi
}

# ====== Main entry point with legacy compatibility ======
# This maintains compatibility with main.sh while providing enhanced functionality
build_binary_main() {
  local _platform_args="${1:-$(uname -s | tr '[:upper:]' '[:lower:]')}"
  local _arch_args="${2:-$(uname -m | tr '[:upper:]' '[:lower:]')}"
  local _force="${3:-false}"
  local _will_upx_pack_binary="${4:-${_WILL_UPX_PACK_BINARY:-true}}"

  # Determine build mode based on UPX setting
  local _build_mode="production"
  [[ "${_will_upx_pack_binary}" == "false" ]] && _build_mode="development"

  # Use the enhanced build system
  build_binary_enhanced "${_platform_args}" "${_arch_args}" "${_force}" "${_build_mode}"

  return $?
}

# ====== Set the main build_binary function ======
# For backwards compatibility, we alias to the main function
alias build_binary=build_binary_main

have() {
  command -v "$1" >/dev/null 2>&1
}
