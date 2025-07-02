#!/usr/bin/env bash

set -euo pipefail
set -o errtrace
set -o functrace
set -o posix
IFS=$'\n\t'

build_binary() {
  local _PLATFORM_ARG="${1:-${_PLATFORM:-}}"
  local _ARCH_ARG="${2:-${_ARCH:-}}"

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
      OUTPUT_NAME=$(printf '%s_%s_%s' "${_BINARY}" "$platform_pos" "$arch_pos")
      if [[ "$platform_pos" == "windows" ]]; then
        OUTPUT_NAME=$(printf '%s.exe' "$OUTPUT_NAME")
      fi

      local build_env=("GOOS=${platform_pos}" "GOARCH=${arch_pos}")
      local build_args=(
        "-ldflags '-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)'"
        "-trimpath -o \"$OUTPUT_NAME\" \"${_CMD_PATH}\""
      )
      local build_cmd=""
      build_cmd=$(printf '%s %s %s' "${build_env[@]}" "go build " "${build_args[@]}")

      log info "Compilando para ${platform_pos}/${arch_pos}"

      if ! bash -c "${build_cmd}"; then
        log error "Falha ao compilar para ${platform_pos} ${arch_pos}"
        return 1
      else
        if [[ "$platform_pos" != "windows" ]]; then
            install_upx || return 1
            upx "$OUTPUT_NAME" --force-overwrite --lzma --no-progress --no-color -qqq || true
            log success "Binário empacotado: ${OUTPUT_NAME}"
        fi
        if [[ ! -f "$OUTPUT_NAME" ]]; then
          log error "Binário não encontrado: ${OUTPUT_NAME}"
          return 1
        else
          compress_binary "$platform_pos" "$arch_pos" || return 1
          log success "Binário criado com sucesso: ${OUTPUT_NAME}"
        fi
      fi
    done
  done
  log success "Todos os builds foram concluídos com sucesso!"
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
        _BINARY_PATH="$(dirname "${BINARY_NAME:-\.\/}")"
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
        log error "Falha ao comprimir para ${platform_pos} ${arch_pos}"
        return 1
      else
        log success "Binário comprimido: ${OUTPUT_NAME}"
      fi
    done
  done
}

export -f build_binary
export -f compress_binary
