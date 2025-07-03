#!/usr/bin/env bash

set -euo pipefail
set -o errtrace
set -o functrace
set -o posix
IFS=$'\n\t'

# Define o diretório raiz (assumindo que este script está em lib/ no root)
_ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
_APP_NAME="${APP_NAME:-$(basename "${_ROOT_DIR}")}"
_PROJECT_NAME="$_APP_NAME"
_OWNER="${OWNER:-rafa-mori}"
# Tenta ler a versão, ou define um fallback
_VERSION=$(cat "$_ROOT_DIR/version/CLI_VERSION" 2>/dev/null || echo "v0.0.0")
# Extrai a versão do Go do go.mod (certifique-se de que este arquivo exista na raiz)
_VERSION_GO=$(grep '^go ' "$_ROOT_DIR/go.mod" | awk '{print $2}')

_LICENSE="MIT"

_ABOUT="################################################################################
  Este script instala o projeto ${_PROJECT_NAME}, versão ${_VERSION}.
  OS suportados: Linux, MacOS, Windows
  Arquiteturas suportadas: amd64, arm64, 386
  Fonte: https://github.com/${_OWNER}/${_PROJECT_NAME}
  Binary Release: https://github.com/${_OWNER}/${_PROJECT_NAME}/releases/latest
  License: ${_LICENSE}
  Notas:
    - [version] é opcional; se omitido, a última versão será utilizada.
    - Se executado localmente, o script tentará resolver a versão pelos tags do repositório.
    - Instala em ~/.local/bin para usuário não-root ou em /usr/local/bin para root.
    - Adiciona o diretório de instalação à variável PATH.
    - Instala o UPX se necessário, ou compila o binário (build) conforme o comando.
    - Faz download do binário via URL de release ou efetua limpeza de artefatos.
    - Verifica dependências e versão do Go.
################################################################################"

_BANNER="################################################################################

               ██   ██ ██     ██ ██████   ████████ ██     ██
              ░██  ██ ░██    ░██░█░░░░██ ░██░░░░░ ░░██   ██
              ░██ ██  ░██    ░██░█   ░██ ░██       ░░██ ██
              ░████   ░██    ░██░██████  ░███████   ░░███
              ░██░██  ░██    ░██░█░░░░ ██░██░░░░     ██░██
              ░██░░██ ░██    ░██░█    ░██░██        ██ ░░██
              ░██ ░░██░░███████ ░███████ ░████████ ██   ░░██
              ░░   ░░  ░░░░░░░  ░░░░░░░  ░░░░░░░░ ░░     ░░   "

# Caminhos para a compilação
_CMD_PATH="$_ROOT_DIR/cmd"
_BUILD_PATH="$(dirname "$_CMD_PATH")"
_BINARY="$_BUILD_PATH/$_APP_NAME"

# Diretórios de instalação
_LOCAL_BIN="${HOME:-"~"}/.local/bin"
_GLOBAL_BIN="/usr/local/bin"

# Caso queira, defina o OWNER (use no get_release_url)
_OWNER="rafa-mori"
