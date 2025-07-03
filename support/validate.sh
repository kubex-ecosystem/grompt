#!/usr/bin/env bash
# lib/validate.sh – Validação da versão do Go e dependências

validate_versions() {
    local REQUIRED_GO_VERSION="${_VERSION_GO:-1.20.0}"
    local GO_VERSION
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if [[ "$(printf '%s\n' "$REQUIRED_GO_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_GO_VERSION" ]]; then
        log error "A versão do Go deve ser >= $REQUIRED_GO_VERSION. Detectado: $GO_VERSION"
        exit 1
    fi
    log success "Versão do Go válida: $GO_VERSION"
    go mod tidy || return 1
}

check_dependencies() {
    for dep in "$@"; do
        if ! command -v "$dep" > /dev/null; then
            log error "$dep não está instalado."
            exit 1
        else
            log success "$dep está instalado."
        fi
    done
}

export -f validate_versions
export -f check_dependencies
