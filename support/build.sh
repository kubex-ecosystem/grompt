#!/usr/bin/env bash

set -euo pipefail
set -o errtrace
set -o functrace
set -o posix

IFS=$'\n\t'

get_output_name() {
  local platform_pos="${1:-${_PLATFORM:-}}"
  local arch_pos="${2:-${_ARCH:-}}"

  local BINARY_DIR="${_ROOT_DIR}/dist"
  if [[ ! -d "$BINARY_DIR" ]]; then
    mkdir -p "$BINARY_DIR" || true
  fi
  local OUTPUT_NAME
  OUTPUT_NAME="$(printf '%s/%s_%s_%s' "${BINARY_DIR}" "${_BINARY_NAME:-${_BINARY}}" "$platform_pos" "$arch_pos")"
  if [[ "$_WILL_UPX_PACK_BINARY" == "true" ]]; then
    OUTPUT_NAME="$(printf '%s/%s_%s_%s' "${BINARY_DIR}" "${_BINARY_NAME:-${_BINARY}}" "$platform_pos" "$arch_pos")"
    if [[ "$platform_pos" == "windows" ]]; then
      OUTPUT_NAME="$(printf '%s/%s_%s_%s.exe' "${BINARY_DIR}" "${_BINARY_NAME:-${_BINARY}}" "$platform_pos" "$arch_pos")"
    else
      OUTPUT_NAME="$(printf '%s/%s_%s_%s' "${BINARY_DIR}" "${_BINARY_NAME:-${_BINARY}}" "$platform_pos" "$arch_pos")"
    fi
  else
    OUTPUT_NAME="${BINARY_DIR}/${_APP_NAME}"
  fi

  echo "$OUTPUT_NAME"
}

build_binary() {
  local _PLATFORM_ARG="${1:-${_PLATFORM:-}}"
  local _ARCH_ARG="${2:-${_ARCH:-}}"
  local _WILL_UPX_PACK_BINARY="${3:-true}"
  local FORCE="${4:-${_FORCE:-n}}"

  # Obtém arrays de plataformas e arquiteturas
  local platforms=( "$(_get_os_arr_from_args "$_PLATFORM_ARG")" )
  local archs=( "$(_get_arch_arr_from_args "$_ARCH_ARG")" )

  for platform_pos in "${platforms[@]}"; do
    [[ -z "$platform_pos" ]] && continue
    for arch_pos in "${archs[@]}"; do
      [[ -z "$arch_pos" ]] && continue
      if [[ "$platform_pos" != "darwin" && "$arch_pos" == "arm64" ]]; then
        continue
      fi
      if [[ "$platform_pos" != "windows" && "$arch_pos" == "386" ]]; then
        continue
      fi

      local OUTPUT_NAME
      OUTPUT_NAME=$(get_output_name "$platform_pos" "$arch_pos")
      if [[ -f "$OUTPUT_NAME" ]]; then
        local REPLY="y"
        if [[ "${IS_INTERACTIVE:-}" != "true" || "${CI:-}" != "true" ]]; then
          if [[ $FORCE =~ [yY] || $FORCE == "true" || "${NON_INTERACTIVE:-}" == "true" ]]; then
            REPLY="y"
          elif [[ -t 0 ]]; then
            # If the script is running interactively, prompt for confirmation
            log notice "Binary already exists: ${OUTPUT_NAME}" true
            log notice "Current binary: ${OUTPUT_NAME}"
            log notice "Press 'y' to overwrite or any other key to skip." true
            log question "(y) to overwrite, any other key to skip (default: n, 10 seconds to respond)" true
            read -t 10 -p "" -n 1 -r REPLY || REPLY="n"
            echo  # Move to a new line after the prompt
            REPLY="${REPLY,,}"  # Convert to lowercase
            REPLY="${REPLY:-n}"  # Default to 'n' if no input
          else
            log notice "Binary already exists: ${OUTPUT_NAME}" true
            log notice "Skipping confirmation in non-interactive mode." true
          fi
        fi
        if [[ ! $REPLY =~ [yY] ]]; then
          log notice "Skipping build for ${platform_pos} ${arch_pos}." true
          continue
        fi
        log warn "Overwriting existing binary: ${OUTPUT_NAME}" true
        if [[ "$platform_pos" == "windows" ]]; then
          rm -f "${OUTPUT_NAME}.exe" || return 1
        else
          rm -f "${OUTPUT_NAME}" || return 1
        fi
      fi
      local build_env=("GOOS=${platform_pos}" "GOARCH=${arch_pos}")
      local build_args=(
        "-ldflags '-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)'"
        "-trimpath -o \"${OUTPUT_NAME}\" \"${_CMD_PATH}\""
      )
      local build_cmd=""
      build_cmd=$(printf '%s %s %s' "${build_env[@]}" "go build " "${build_args[@]}")

      log info "Building for ${platform_pos} ${arch_pos}..."

      if ! bash -c "${build_cmd}"; then
        log error "Failed to build for ${platform_pos} ${arch_pos}"
        return 1
      else
        if [[ "$_WILL_UPX_PACK_BINARY" == "true" ]]; then
          if [[ "$platform_pos" != "windows" ]]; then
              install_upx || return 1
              upx "$OUTPUT_NAME" --force-overwrite --lzma --no-progress --no-color -qqq || true
              log success "Packed binary: ${OUTPUT_NAME}"
          fi
          if [[ ! -f "$OUTPUT_NAME" ]]; then
            log error "Binary not found: ${OUTPUT_NAME}"
            return 1
          else
            compress_binary "$platform_pos" "$arch_pos" || return 1
            log success "Binary created successfully: ${OUTPUT_NAME}"
          fi
        else
          log warn "UPX packing disabled. The binary will not be compressed." true
          log warn "Build indicated for development use only." true
          log success "Binary created without packing: ${OUTPUT_NAME}"
        fi
      fi
    done
  done
  return 0
}

compress_binary() {
  local platform_arg="${1:-${_PLATFORM:-}}"
  local arch_arg="${2:-${_ARCH:-}}"

  # Obtém arrays de plataformas e arquiteturas
  local platforms=( "$(_get_os_arr_from_args "$platform_arg")" )
  local archs=( "$(_get_arch_arr_from_args "$arch_arg")" )

  for platform_pos in "${platforms[@]}"; do
    [[ -z "$platform_pos" ]] && continue
    for arch_pos in "${archs[@]}"; do
      [[ -z "$arch_pos" ]] && continue
      if [[ "$platform_pos" != "darwin" && "$arch_pos" == "arm64" ]]; then
        continue
      fi
      if [[ "$platform_pos" == "linux" && "$arch_pos" == "386" ]]; then
        continue
      fi
      local BINARY_NAME
      BINARY_NAME=$(printf '%s_%s_%s' "${_BINARY}" "$platform_pos" "$arch_pos")
      if [[ "$platform_pos" == "windows" ]]; then
        BINARY_NAME=$(printf '%s.exe' "${BINARY_NAME}")
      fi
      local OUTPUT_NAME="${BINARY_NAME//.exe/}"
      local compress_cmd_exec=""
      if [[ "$platform_pos" != "windows" ]]; then
        OUTPUT_NAME="${OUTPUT_NAME}.tar.gz"
        _CURR_PATH="$(pwd)"
        _BINARY_PATH="${_ROOT_DIR}/dist"

        cd "${_BINARY_PATH}" || true # Just to avoid tar warning about relative paths
        if tar -czf "./$(basename "${OUTPUT_NAME}")" "./$(basename "${BINARY_NAME}")"; then
          compress_cmd_exec="true"
        else
          compress_cmd_exec="false"
        fi
        cd "${_CURR_PATH}" || true
      else
        OUTPUT_NAME="${OUTPUT_NAME}.zip"
        # log info "Comprimindo para ${platform_pos} ${arch_pos} em ${OUTPUT_NAME}..."
        if zip -r -9 "${OUTPUT_NAME}" "${BINARY_NAME}" >/dev/null; then
          compress_cmd_exec="true"
        else
          compress_cmd_exec="false"
        fi
      fi
      if [[ "$compress_cmd_exec" == "false" ]]; then
        log error "Failed to compress for ${platform_pos} ${arch_pos}"
        return 1
      else
        log success "Compressed binary: ${OUTPUT_NAME}"
      fi
    done
  done
}

export -f build_binary
export -f compress_binary
