#!/usr/bin/env bash
# lib/install_funcs.sh – Functions for installing binaries and managing paths

install_upx() {
    if ! command -v upx &> /dev/null; then
        if ! sudo -v &> /dev/null; then
            log error "You do not have permission to install UPX."
            log warn "If you want binary packing, please install UPX manually."
            log warn "See: https://upx.github.io/"
            return 1
        fi
        if [[ "$(uname)" == "Darwin" ]]; then
            brew install upx >/dev/null
        elif command -v apt-get &> /dev/null; then
            sudo apt-get install -y upx >/dev/null
        elif command -v yum &> /dev/null; then
            sudo yum install -y upx >/dev/null
        elif command -v dnf &> /dev/null; then
            sudo dnf install -y upx >/dev/null
        elif command -v pacman &> /dev/null; then
            sudo pacman -S --noconfirm upx >/dev/null
        elif command -v zypper &> /dev/null; then
            sudo zypper install -y upx >/dev/null
        elif command -v apk &> /dev/null; then
            sudo apk add upx >/dev/null
        elif command -v port &> /dev/null; then
            sudo port install upx >/dev/null
        elif command -v snap &> /dev/null; then
            sudo snap install upx >/dev/null
        elif command -v flatpak &> /dev/null; then
            sudo flatpak install flathub org.uptane.upx -y >/dev/null
        else
            log warn "If you want binary packing, please install UPX manually."
            log warn "See: https://upx.github.io/"
            return 1
        fi
    fi

    return 0
}
detect_shell_rc() {
    local shell_rc_file
    local user_shell
    user_shell=$(basename "$SHELL")

    case "$user_shell" in
        bash) shell_rc_file="${HOME:-~}/.bashrc" ;;
        zsh) shell_rc_file="${HOME:-~}/.zshrc" ;;
        sh) shell_rc_file="${HOME:-~}/.profile" ;;
        fish) shell_rc_file="${HOME:-~}/.config/fish/config.fish" ;;
        *)
            log warn "Unsupported shell: $user_shell"
            log warn "Please add the path manually to your shell configuration file."
            log warn "Supported shells: bash, zsh, sh, fish"
            return 1
            ;;
    esac

    if [ ! -f "$shell_rc_file" ]; then
        log error "Configuration file not found: ${shell_rc_file}"
        return 1
    fi

    echo "$shell_rc_file"

    return 0
}
add_to_path() {
    local target_path="${1:-}"

    local shell_rc_file=""

    local path_expression=""

    path_expression="export PATH=\"${target_path}:\$PATH\""

    shell_rc_file="$(detect_shell_rc)"

    if [[ -z "$shell_rc_file" ]]; then
        log error "Could not identify the shell configuration file."
        return 1
    fi
    if grep -q "${path_expression}" "$shell_rc_file" 2>/dev/null; then
        log success "$target_path is already in PATH for $shell_rc_file."
        return 0
    fi

    if [[ -z "${target_path}" ]]; then
        log error "Target path not provided."
        return 1
    fi

    if [[ ! -d "${target_path}" ]]; then
        log error "Target path is not a valid directory: $target_path"
        return 1
    fi

    if [[ ! -f "${shell_rc_file}" ]]; then
        log error "Configuration file not found: ${shell_rc_file}"
        return 1
    fi

    # echo "export PATH=${target_path}:\$PATH" >> "$shell_rc_file"
    printf '%s\n' "${path_expression}" | tee -a "$shell_rc_file" >/dev/null || {
        log error "Failed to add $target_path to PATH in $shell_rc_file."
        return 1
    }

    log success "Added $target_path to PATH in $shell_rc_file."

    "$SHELL" -c "source ${shell_rc_file}" || {
        log warn "Failed to reload shell. Please run 'source ${shell_rc_file}' manually."
    }

    return 0
}
install_binary() {
    local SUFFIX="${_PLATFORM_WITH_ARCH}"
    local BINARY_TO_INSTALL="${_BINARY/\/${_APP_NAME}/\/${_APP_NAME}\/dist}${SUFFIX:+_${SUFFIX}}"
    log info "Installing binary: '${BINARY_TO_INSTALL}' as '$_APP_NAME'"

    if [[ "$(id -u)" -ne 0 ]]; then
        log info "Non-root user detected. Installing to ${_LOCAL_BIN}..."
        mkdir -p "$_LOCAL_BIN"
        cp "$BINARY_TO_INSTALL" "$_LOCAL_BIN/$_APP_NAME" || exit 1
        add_to_path "$_LOCAL_BIN"
    else
        log info "Root user detected. Installing to ${_GLOBAL_BIN}..."
        cp "$BINARY_TO_INSTALL" "$_GLOBAL_BIN/$_APP_NAME" || exit 1
        add_to_path "$_GLOBAL_BIN"
    fi

    if [[ -n "$shell_rc_file" ]]; then
      # shellcheck source=/dev/null
      . "${shell_rc_file:-$(detect_shell_rc)}" || {
          log warn "Failed to reload shell configuration. Please run 'source ${shell_rc_file}' manually."
      }
    fi
}
uninstall_binary() {
    log info "Uninstalling binary: '$_APP_NAME'"

    if [ "$(id -u)" -ne 0 ]; then
        log info "Non-root user detected. Uninstalling from ${_LOCAL_BIN}..."
        rm -f "$_LOCAL_BIN/$_APP_NAME"
    else
        log info "Root user detected. Uninstalling from ${_GLOBAL_BIN}..."
        rm -f "$_GLOBAL_BIN/$_APP_NAME"
    fi

    log success "Binary '$_APP_NAME' uninstalled successfully."

    local shell_rc_file=""
    shell_rc_file="$(detect_shell_rc)"

    if [[ -n "$shell_rc_file" ]]; then
      # shellcheck source=/dev/null
      . "${shell_rc_file:-$(detect_shell_rc)}" || {
          log warn "Failed to reload shell configuration. Please run 'source ${shell_rc_file}' manually."
      }
    fi

    return 0
}
download_binary() {
    if [[ "${_PRIVATE_REPOSITORY:-}" == "true" ]]; then
        log error "It is not possible to download from a private repository." true
        log error "Please clone the repository and build the binary locally." true
        return 1
    fi
    if ! what_platform; then
        log error "Failed to detect platform."
        return 1
    fi
    if [[ -z "${_PLATFORM}" ]]; then
        log error "Unsupported platform: ${_PLATFORM}"
        return 1
    fi
    local version
    if [[ -z "${_VERSION}" ]]; then
        version="${_REPOSITORY}/releases/latest"  # Use latest release if version is not specified
    else
        version="${_REPOSITORY}/releases/tag/${_VERSION}"
    fi
    if [ -z "$version" ]; then
        log error "Failed to determine latest version."
        return 1
    fi

    local release_url
    release_url=$(get_release_url)
    log info "Downloading binary ${_APP_NAME} for OS=${_PLATFORM}, ARCH=${_ARCH}, Version=${version}..."
    log info "Release URL: ${release_url}"

    local archive_path="${_TEMP_DIR}/${_APP_NAME}.tar.gz"
    if ! curl -L -o "${archive_path}" "${release_url}"; then
        log error "Failed to download binary from: ${release_url}"
        return 1
    fi
    log success "Binary downloaded successfully."

    log info "Extracting binary to: $(dirname "${_BINARY}")"
    if ! tar -xzf "${archive_path}" -C "$(dirname "${_BINARY}")"; then
        log error "Failed to extract binary from: ${archive_path}"
        rm -rf "${_TEMP_DIR}"
        return 1
    fi

    rm -rf "${_TEMP_DIR}"
    log success "Binary extracted successfully."

    if [ ! -f "$_BINARY" ]; then
        log error "Binary not found after extraction: ${_BINARY}"
        return 1
    fi
    log success "Download and extraction of ${_APP_NAME} completed successfully!"
    return 0
}
install_from_release() {
    download_binary || return 1
    install_binary || return 1
}
check_path() {
    log info "Checking if $1 is in PATH..."
    if ! echo "$PATH" | grep -q "$1"; then
        log warn "$1 is not in PATH."
        log warn "Add: export PATH=$1:\$PATH"
    else
        log success "$1 is already in PATH."
    fi
}

export -f install_upx
export -f detect_shell_rc
export -f add_to_path
export -f install_binary
export -f download_binary
export -f install_from_release
export -f check_path
