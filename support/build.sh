#!/usr/bin/env bash

set -euo pipefail
set -o errtrace
set -o functrace
set -o posix

IFS=$'\n\t'

get_output_name() {
  local _platform_pos="${1:-${_PLATFORM:-$(uname -s | tr '[:upper:]' '[:lower:]')}}"
  local _arch_pos="${2:-${_ARCH:-$(uname -m | tr '[:upper:]' '[:lower:]')}}"
  local _root_dir="${_ROOT_DIR:-${ROOT_DIR:-$(git rev-parse --show-toplevel)}}"
  local _binary_dir="${_root_dir:-}/dist"
  local _cmd_path="${_CMD_PATH:-${CMD_PATH:-${_root_dir}/cmd}}"
  local _binary_name="${_BINARY_NAME:-${BINARY_NAME:-$(basename "${_cmd_path}" .go)}}"
  local _will_upx_pack_binary="${_WILL_UPX_PACK_BINARY:-${WILL_UPX_PACK_BINARY:-true}}"

  if [[ ! -d "${_binary_dir:-}" ]]; then
    mkdir -p "${_binary_dir:-}" || true
  fi

  local _output_name=""
  _output_name="$(printf '%s/%s_%s_%s' "${_binary_dir:-}" "${_binary_name:-}" "${_platform_pos:-}" "${_arch_pos:-}")"
  if [[ "$_will_upx_pack_binary" == "true" ]]; then
    _output_name="$(printf '%s/%s_%s_%s' "${_binary_dir:-}" "${_binary_name:-}" "${_platform_pos:-}" "${_arch_pos:-}")"
    if [[ "$_platform_pos" == "windows" ]]; then
      _output_name="$(printf '%s/%s_%s_%s.exe' "${_binary_dir:-}" "${_binary_name:-}" "${_platform_pos:-}" "${_arch_pos:-}")"
    else
      _output_name="$(printf '%s/%s_%s_%s' "${_binary_dir:-}" "${_binary_name:-}" "${_platform_pos:-}" "${_arch_pos:-}")"
    fi
  else
    _output_name="${_binary_dir:-}/${_APP_NAME:-}"
  fi

  echo "${_output_name:-}"
}

upx_packaging() {
  local _output_name="${1:-${_OUTPUT_NAME:-${OUTPUT_NAME:-}}}"
  local _platform_pos="${2:-${_PLATFORM:-${PLATFORM:-$(uname -s | tr '[:upper:]' '[:lower:]')}}}"
  local _arch_pos="${3:-${_ARCH:-${ARCH:-$(uname -m | tr '[:upper:]' '[:lower:]')}}}"

  install_upx || return 1

  if [[ "${_platform_pos:-}" != "darwin" ]]; then
    upx "${_output_name:-}" --force-overwrite --lzma --no-progress --no-color -qqq || true
    log success "Packed binary: ${_output_name:-}"
  else
    upx "${_output_name:-}" --force-overwrite --lzma --no-progress --force-macos --no-color -qqq || true
    log warn "UPX packing on macOS is not fully supported. The binary may not be compressed properly." true
  fi

  return 0
}

compile_binary() {
  local _platform_pos="${1:-${_PLATFORM:-$(uname -s | tr '[:upper:]' '[:lower:]')}}"
  local _arch_pos="${2:-${_ARCH:-$(uname -m | tr '[:upper:]' '[:lower:]')}}"
  local _output_name="${3:-${_OUTPUT_NAME:-${OUTPUT_NAME:-}}}"
  local _force="${4:-${_FORCE:-n}}"
  local _will_upx_pack_binary="${5:-${_WILL_UPX_PACK_BINARY:-true}}"

  local _root_dir="${_root_dir:-${ROOT_DIR:-$(git rev-parse --show-toplevel)}}"
  local _cmd_path="${_cmd_path:-${CMD_PATH:-${_root_dir}/cmd}}"

  local _binary_name="${_binary_name:-${BINARY_NAME:-$(basename "${_cmd_path}" .go)}}"
  local _app_name="${_app_name:-${APP_NAME:-$(basename "${_root_dir}")}}"
  local _build_env=("GOOS=${_platform_pos}" "GOARCH=${_arch_pos}")

  local _build_cmd=""
  local _build_args=(
    "-ldflags '-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)'"
    "-trimpath -o \"${_output_name:-}\" \"${_cmd_path:-}\""
  )
  _build_cmd=$(printf '%s %s %s' "${_build_env[@]}" "go build " "${_build_args[@]}")
  log info "Building for ${_platform_pos:-} ${_arch_pos:-}..." true

  eval "${_build_cmd}" || {
    log error "Failed to build for ${_platform_pos:-} ${_arch_pos:-}" true
    return 1
  }

  log success "Binary built: ${_output_name:-}" true

  return 0
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
  local _platform_pos="${1:-${_PLATFORM:-$(uname -s | tr '[:upper:]' '[:lower:]')}}"
  local _arch_arg="${2:-${_ARCH:-$(uname -m | tr '[:upper:]' '[:lower:]')}}"
  local _will_upx_pack_binary="${3:-${_will_upx_pack_binary:-${WILL_UPX_PACK_BINARY:-true}}}"
  local _force="${4:-${_force:-${FORCE:-}}}"

  local _root_dir="${_root_dir:-${ROOT_DIR:-$(git rev-parse --show-toplevel)}}"
  local _cmd_path="${_cmd_path:-${CMD_PATH:-${_root_dir}/cmd}}"

  local _archs=( "$(_get_arch_arr_from_args "${_platform_pos:-}" "${_arch_arg:-}")" )

  [[ -z "${_platform_pos:-}" ]] && return 0
  [[ -z "${_arch_pos:-}" ]] && return 0

  if [[ "${_platform_pos:-}" != "darwin" && "${_arch_pos:-}" == "arm64" ]]; then
    log info "Unsupported build for ${_platform_pos:-} ${_arch_pos:-}. Skipping."
    return 0
  fi
  if [[ "${_platform_pos:-}" != "windows" && "${_arch_pos:-}" == "386" ]]; then
    log info "Skipping unsupported build for ${_platform_pos:-} ${_arch_pos:-}."
    return 0
  fi

  local _old_bin_name
  _old_bin_name=$(get_output_name "${_platform_pos:-}" "${_arch_pos:-}")

  if ! check_overwrite_binary "${_platform_pos:-}" "${_arch_pos:-}" "${_old_bin_name:-}" "${_force:-}" "${_will_upx_pack_binary:-}" "true"; then
    return 0
  fi

  if ! compile_binary "${_platform_pos:-}" "${_arch_pos:-}" "${_old_bin_name:-}" "${_force:-}" "${_will_upx_pack_binary:-}"; then
    log error "Failed to build for ${_platform_pos:-} ${_arch_pos:-}" true
    return 1
  else
    if [[ "${_will_upx_pack_binary:-}" == "true" ]]; then
      if upx_packaging "${_old_bin_name:-}" "${_platform_pos:-}" "${_arch_pos:-}"; then
        log success "UPX packing successful: ${_old_bin_name:-}"
      else
        log error "UPX packing failed: ${_old_bin_name:-}" true
        return 1
      fi

      if [[ ! -f "${_old_bin_name:-}" ]]; then
        log error "Binary not found after build: ${_old_bin_name:-}" true
        return 1
      else
        compress_binary "${_platform_pos:-}" "${_arch_pos:-}" || return 1
        log success "Binary created successfully: ${_old_bin_name:-}"
      fi

    else
      log warn "UPX packing disabled. The binary will not be compressed." true
      log warn "Build indicated for development use only." true
      log success "Binary created without packing: ${_old_bin_name:-}"
    fi
  fi

  return 0
}

platform_iterator() {
  local _platform_arg="${1:-${_PLATFORM:-${PLATFORM:-$(uname -s | tr '[:upper:]' '[:lower:]')}}}"
  local _arch_arg="${2:-${_ARCH:-${ARCH:-$(uname -m | tr '[:upper:]' '[:lower:]')}}}"
  local _will_upx_pack_binary="${4:-${_will_upx_pack_binary:-${WILL_UPX_PACK_BINARY:-true}}}"
  local _force="${3:-${_force:-${FORCE:-}}}"

  local _platform_pos="${1:-${_PLATFORM:-${PLATFORM:-}}}"

  _platform_pos="${_PLATFORM_ARG:-}"
  [[ -z "${_platform_pos:-}" ]] && return 0

  # Search for spaces inside, if found, its an array
  if [[ "${_platform_pos:-}" =~ \  ]]; then
    platforms=( "${_platform_pos// /_}" )
    for _platform_pos in "${platforms[@]}"; do
      local _arch_pos=""
      _arch_pos="${2:-${_ARCH:-${ARCH:-}}}"
      _arch_pos="${_ARCH_ARG:-}"
      [[ -z "${_arch_pos:-}" ]] && return 0
      if ! arch_iterator "${_platform_pos:-}" "${_arch_pos:-}" "${_force:-}" "${_will_upx_pack_binary:-}"; then
        log error "Failed to build for ${_platform_pos:-} ${_arch_pos:-}"
        return 1
      fi
    done
  else
    for _arch_pos in "${_archs[@]}"; do
      _arch_pos="${2:-${_ARCH:-}}"
      _arch_pos="${_ARCH_ARG:-}"
      [[ -z "${_arch_pos:-}" ]] && return 0
      if ! arch_iterator "${_platform_pos:-}" "${_arch_pos:-}" "${_force:-}" "${_will_upx_pack_binary:-}"; then
        log error "Failed to build for ${_platform_pos:-} ${_arch_pos:-}" true
        return 1
      fi
    done
  fi

  return 0
}

compress_binary() {
  local _platform_arg="${1:-${_PLATFORM:-$(uname -s | tr '[:upper:]' '[:lower:]')}}"
  local _arch_arg="${2:-${_ARCH:-$(uname -m | tr '[:upper:]' '[:lower:]')}}"
  local _output_name="${3:-${_OUTPUT_NAME:-}}"

  local _ROOT_DIR="${_ROOT_DIR:-${ROOT_DIR:-$(git rev-parse --show-toplevel)}}"
  local _BINARY_PATH="${_BINARY_PATH:-${BINARY_PATH:-${_ROOT_DIR:-}/dist}}"

  # Obtém arrays de plataformas e arquiteturas
  local _platforms=( "$(_get_os_arr_from_args "${_platform_arg:-}")" )
  local _archs=( "$(_get_arch_arr_from_args "${_platform_pos:-}" "${_arch_arg:-}")" )

  [[ -z "${_platform_arg:-}" ]] && _platform_arg="linux"
  [[ -z "${_arch_arg:-}" ]] && _arch_arg="amd64"

  for _platform_pos in "${_platforms[@]}"; do
    [[ -z "${_platform_pos:-}" ]] && continue
    for _arch_pos in "${_archs[@]}"; do
      [[ -z "${_arch_pos:-}" ]] && continue

      if [[ "${_platform_pos:-}" != "darwin" && "${_arch_pos:-}" == "arm64" ]]; then
        continue
      fi
      if [[ "${_platform_pos:-}" == "linux" && "${_arch_pos:-}" == "386" ]]; then
        continue
      fi

      local BINARY_NAME
      BINARY_NAME=$(printf '%s_%s_%s' "${_BINARY:-}" "${_platform_pos:-}" "${_arch_pos:-}")
      if [[ "${_platform_pos:-}" == "windows" ]]; then
        BINARY_NAME=$(printf '%s.exe' "${BINARY_NAME:-}")
      fi

      local OUTPUT_NAME="${BINARY_NAME//.exe/}"
      local compress_cmd_exec=""

      if [[ "${_platform_pos:-}" != "windows" ]]; then
        OUTPUT_NAME="${OUTPUT_NAME:-}.tar.gz"
        _CURR_PATH="$(pwd)"
        _BINARY_PATH="${_ROOT_DIR:-}/dist"

        cd "${_BINARY_PATH:-}" || true # Just to avoid tar warning about relative paths
        if tar -czf "./$(basename "${OUTPUT_NAME:-}")" "./$(basename "${BINARY_NAME:-}")"; then
          compress_cmd_exec="true"
        else
          compress_cmd_exec="false"
        fi
        cd "${_CURR_PATH:-}" || true
      else
        OUTPUT_NAME="${OUTPUT_NAME:-}.zip"
        # log info "Comprimindo para ${_platform_pos} ${_arch_pos} em ${OUTPUT_NAME}..."
        if zip -r -9 "${OUTPUT_NAME:-}" "${BINARY_NAME:-}" >/dev/null; then
          compress_cmd_exec="true"
        else
          compress_cmd_exec="false"
        fi
      fi

      if [[ "${compress_cmd_exec:-}" == "false" ]]; then
        log error "Failed to compress for ${_platform_pos:-} ${_arch_pos:-}"
        return 1
      else
        log success "Compressed binary: ${OUTPUT_NAME:-}"
      fi
    done
  done
}

# Entrypoint for building binaries
build_binary() {
  local _platform_args="${1:-${_PLATFORM:-${PLATFORM:-$(uname -s | tr '[:upper:]' '[:lower:]')}}}"
  local _arch_args="${2:-${_ARCH:-${ARCH:-$(uname -m | tr '[:upper:]' '[:lower:]')}}}"
  local _force="${3:-${_FORCE:-${FORCE:-}}}"
  local _will_upx_pack_binary="${4:-${_WILL_UPX_PACK_BINARY:-${WILL_UPX_PACK_BINARY:-true}}}"

  local _root_dir="${_ROOT_DIR:-${ROOT_DIR:-$(git rev-parse --show-toplevel)}}"
  local _cmd_path="${_CMD_PATH:-${CMD_PATH:-${_root_dir}/cmd}}"
  local _binary_name="${_BINARY_NAME:-${BINARY_NAME:-$(basename "${_cmd_path:-}" .go)}}"
  local _app_name="${_APP_NAME:-${APP_NAME:-$(basename "${_root_dir:-}")}}"
  local _version="${_VERSION:-${VERSION:-$(git describe --tags)}}"

  log notice "Binary Name: ${_binary_name:-}"
  log notice "App Name: ${_app_name:-}"
  log notice "Version: ${_version:-}"

  local _platforms=()
  _platforms=( $(_get_os_arr_from_args "${_platform_args:-}") )

  [[ -z "${_platforms[*]}" ]] && _platforms=("linux")

  for _platform_pos in "${_platforms[@]}"; do

    # [[ -z "${_platform_pos:-}" ]] && continue

    local _archs=( $(_get_arch_arr_from_args "${_platform_pos:-}" "${_arch_args:-}") )

    if ! platform_iterator "${_platform_pos:-}" "${_archs[*]}" "${_force:-false}" "${_will_upx_pack_binary:-true}" "${_binary_name:-}"; then
      log error "Failed to build for platform: ${_platform_pos:-}" true
      return 1
    fi

  done

  return 0
}

export -f build_binary
export -f compress_binary
