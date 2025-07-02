#!/usr/bin/env bash
# lib/install_funcs.sh – Funções para instalação e manipulação de PATH

install_upx() {
    if ! command -v upx &> /dev/null; then
        if ! sudo -v &> /dev/null; then
            log error "Você não tem permissões de superusuário para instalar o empacotador de binários."
            log warn "Se deseja o empacotamento de binários, instale o UPX manualmente."
            log warn "Veja: https://upx.github.io/"
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
            log warn "Se deseja o empacotamento de binários, instale o UPX manualmente."
            log warn "Veja: https://upx.github.io/"
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
            log warn "Shell não suportado; ajuste o PATH manualmente."
            return 1
            ;;
    esac
    
    if [ ! -f "$shell_rc_file" ]; then
        log error "Arquivo de configuração não encontrado: ${shell_rc_file}"
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


    if [ -z "$shell_rc_file" ]; then
        log error "Não foi possível identificar o arquivo de configuração do shell."
        return 1
    fi
    if grep -q "${path_expression}" "$shell_rc_file" 2>/dev/null; then
        log success "$target_path já está no PATH do $shell_rc_file."
        return 0
    fi

    if [[ -z "${target_path}" ]]; then
        log error "Caminho de destino não fornecido."
        return 1
    fi

    if [[ ! -d "${target_path}" ]]; then
        log error "Caminho de destino não é um diretório válido: $target_path"
        return 1
    fi

    if [[ ! -f "${shell_rc_file}" ]]; then
        log error "Arquivo de configuração não encontrado: ${shell_rc_file}"
        return 1
    fi

    # echo "export PATH=${target_path}:\$PATH" >> "$shell_rc_file"
    printf '%s\n' "${path_expression}" | tee -a "$shell_rc_file" >/dev/null || {
        log error "Falha ao adicionar $target_path ao PATH em $shell_rc_file."
        return 1
    }

    log success "Adicionado $target_path ao PATH em $shell_rc_file."
    
    "$SHELL" -c "source ${shell_rc_file}" || {
        log warn "Falha ao recarregar o shell. Por favor, execute 'source ${shell_rc_file}' manualmente."
    }

    return 0
}

install_binary() {
    local SUFFIX="${_PLATFORM_WITH_ARCH}"
    local BINARY_TO_INSTALL="${_BINARY}${SUFFIX:+_${SUFFIX}}"
    log info "Instalando o binário: '${BINARY_TO_INSTALL}' como '$_APP_NAME'"

    if [ "$(id -u)" -ne 0 ]; then
        log info "Usuário não-root detectado. Instalando em ${_LOCAL_BIN}..."
        mkdir -p "$_LOCAL_BIN"
        cp "$BINARY_TO_INSTALL" "$_LOCAL_BIN/$_APP_NAME" || exit 1
        add_to_path "$_LOCAL_BIN"
    else
        log info "Usuário root detectado. Instalando em ${_GLOBAL_BIN}..."
        cp "$BINARY_TO_INSTALL" "$_GLOBAL_BIN/$_APP_NAME" || exit 1
        add_to_path "$_GLOBAL_BIN"
    fi
}

download_binary() {
    if ! what_platform; then
        log error "Falha ao detectar a plataforma."
        return 1
    fi
    if [[ -z "${_PLATFORM}" ]]; then
        log error "Plataforma não suportada: ${_PLATFORM}"
        return 1
    fi
    local version
    version=$(curl -s "https://api.github.com/repos/${_OWNER}/${_PROJECT_NAME}/releases/latest" | grep "tag_name" | cut -d '"' -f 4 || echo "latest")
    if [ -z "$version" ]; then
        log error "Falha ao determinar a última versão."
        return 1
    fi

    local release_url
    release_url=$(get_release_url)
    log info "Baixando o binário ${_APP_NAME} para OS=${_PLATFORM}, ARCH=${_ARCH}, Versão=${version}..."
    log info "URL de Release: ${release_url}"

    local archive_path="${_TEMP_DIR}/${_APP_NAME}.tar.gz"
    if ! curl -L -o "${archive_path}" "${release_url}"; then
        log error "Falha ao baixar o binário de: ${release_url}"
        return 1
    fi
    log success "Binário baixado com sucesso."

    log info "Extraindo o binário para: $(dirname "${_BINARY}")"
    if ! tar -xzf "${archive_path}" -C "$(dirname "${_BINARY}")"; then
        log error "Falha ao extrair o binário de: ${archive_path}"
        rm -rf "${_TEMP_DIR}"
        exit 1
    fi

    rm -rf "${_TEMP_DIR}"
    log success "Binário extraído com sucesso."

    if [ ! -f "$_BINARY" ]; then
        log error "Binário não encontrado após extração: ${_BINARY}"
        exit 1
    fi
    log success "Download e extração de ${_APP_NAME} concluídos!"
}

install_from_release() {
    download_binary
    install_binary
}

check_path() {
    log info "Verificando se o diretório de instalação está no PATH..."
    if ! echo "$PATH" | grep -q "$1"; then
        log warn "$1 não está no PATH."
        log warn "Adicione: export PATH=$1:\$PATH"
    else
        log success "$1 já está no PATH."
    fi
}

export -f install_upx
export -f detect_shell_rc
export -f add_to_path
export -f install_binary
export -f download_binary
export -f install_from_release
export -f check_path