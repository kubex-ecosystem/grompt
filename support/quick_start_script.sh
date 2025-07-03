#!/usr/bin/env bash

# ╔══════════════════════════════════════════════════════════════╗
# ║                 🚀 Grompt Quick Start               ║
# ║                                                              ║
# ║          Script de instalação e execução automática         ║
# ╚══════════════════════════════════════════════════════════════╝

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

# Project info
PROJECT_NAME="grompt"
SETUP_SCRIPT_URL="https://raw.githubusercontent.com/SEU_USUARIO/grompt-setup/main/setup.sh"

print_banner() {
    echo -e "${PURPLE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                 🚀 Grompt Quick Start               ║"
    echo "║                                                              ║"
    echo "║          Instalação e execução automática                   ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo
}

print_step() {
    echo -e "${CYAN}🔧 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

check_dependencies() {
    print_step "Verificando dependências do sistema..."
    
    # Check Go
    if ! command -v go &> /dev/null; then
        print_error "Go não encontrado! Instale Go 1.21+ antes de continuar."
    fi
    
    GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | sed 's/go//')
    if [[ "$(printf '%s\n' "1.21" "$GO_VERSION" | sort -V | head -n1)" != "1.21" ]]; then
        print_error "Go 1.21+ é necessário. Versão atual: $GO_VERSION"
    fi
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        echo -e "${YELLOW}⚠️  Node.js não encontrado. Instalando via nvm...${NC}"
        install_nodejs
    fi
    
    NODE_VERSION=$(node --version | sed 's/v//')
    if [[ "$(printf '%s\n' "16.0.0" "$NODE_VERSION" | sort -V | head -n1)" != "16.0.0" ]]; then
        echo -e "${YELLOW}⚠️  Node.js 16+ recomendado. Versão atual: $NODE_VERSION${NC}"
    fi
    
    # Check npm
    if ! command -v npm &> /dev/null; then
        print_error "npm não encontrado! Instale npm junto com Node.js."
    fi
    
    print_success "Dependências verificadas"
}

install_nodejs() {
    print_step "Instalando Node.js via nvm..."
    
    # Install nvm if not present
    if ! command -v nvm &> /dev/null; then
        curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
        export NVM_DIR="$HOME/.nvm"
        [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
    fi
    
    # Install latest LTS Node.js
    nvm install --lts
    nvm use --lts
    
    print_success "Node.js instalado"
}

run_setup() {
    print_step "Executando script de configuração..."
    
    # Baixar e executar o script de setup se não existir localmente
    if [ ! -f "setup.sh" ]; then
        echo -e "${YELLOW}📥 Script de setup não encontrado localmente.${NC}"
        echo -e "${YELLOW}💡 Execute o setup.sh que você já tem ou baixe de um repositório.${NC}"
        echo -e "${YELLOW}🔧 Por enquanto, vou criar a estrutura localmente...${NC}"
        echo
        
        # Se não tiver o setup.sh, criar estrutura básica
        create_basic_structure
    else
        chmod +x setup.sh
        ./setup.sh
    fi
}

create_basic_structure() {
    print_step "Criando estrutura básica do projeto..."
    
    if [ -d "$PROJECT_NAME" ]; then
        print_error "Diretório '$PROJECT_NAME' já existe!"
    fi
    
    mkdir -p "$PROJECT_NAME"
    cd "$PROJECT_NAME"
    
    print_step "Para continuar, você precisa do arquivo setup.sh completo."
    print_step "Execute o script setup.sh que foi gerado anteriormente."
    
    cd ..
    print_success "Estrutura básica criada. Execute o setup.sh para continuar."
}

build_and_run() {
    print_step "Fazendo build e executando aplicação..."
    
    cd "$PROJECT_NAME"
    
    # Install dependencies
    if [ -f "Makefile" ]; then
        make install-deps
        make build-all
        
        if [ -f "grompt" ] || [ -f "grompt.exe" ]; then
            print_success "Build concluído com sucesso!"
            
            echo -e "${CYAN}🚀 Iniciando Grompt...${NC}"
            echo -e "${YELLOW}💡 A aplicação abrirá automaticamente no seu navegador${NC}"
            echo -e "${YELLOW}🛑 Pressione Ctrl+C para parar o servidor${NC}"
            echo
            
            # Execute the application
            if [ -f "grompt" ]; then
                ./grompt
            else
                ./grompt.exe
            fi
        else
            print_error "Build falhou! Verifique as dependências."
        fi
    else
        print_error "Makefile não encontrado! Execute o setup.sh primeiro."
    fi
}

show_manual_instructions() {
    echo
    echo -e "${CYAN}📚 Instruções Manuais:${NC}"
    echo
    echo -e "${YELLOW}1. Gerar estrutura:${NC}"
    echo -e "   ${BLUE}curl -O [URL_DO_SETUP_SCRIPT]${NC}"
    echo -e "   ${BLUE}chmod +x setup.sh${NC}"
    echo -e "   ${BLUE}./setup.sh${NC}"
    echo
    echo -e "${YELLOW}2. Build e execução:${NC}"
    echo -e "   ${BLUE}cd grompt${NC}"
    echo -e "   ${BLUE}make install-deps${NC}"
    echo -e "   ${BLUE}make build-all${NC}"
    echo -e "   ${BLUE}./grompt${NC}"
    echo
    echo -e "${YELLOW}3. Configurar APIs (opcional):${NC}"
    echo -e "   ${BLUE}export CLAUDE_API_KEY=your_claude_api_key${NC}"
    echo -e "   ${BLUE}export OLLAMA_ENDPOINT=http://localhost:11434${NC}"
    echo
}

main() {
    print_banner
    
    echo -e "${CYAN}🎯 O que você gostaria de fazer?${NC}"
    echo -e "${YELLOW}1)${NC} Instalação completa (verificar deps + setup + build + run)"
    echo -e "${YELLOW}2)${NC} Apenas verificar dependências"
    echo -e "${YELLOW}3)${NC} Apenas executar setup"
    echo -e "${YELLOW}4)${NC} Apenas build e executar (se já configurado)"
    echo -e "${YELLOW}5)${NC} Mostrar instruções manuais"
    echo
    
    read -p "Escolha uma opção (1-5): " choice
    echo
    
    case $choice in
        1)
            check_dependencies
            run_setup
            build_and_run
            ;;
        2)
            check_dependencies
            ;;
        3)
            run_setup
            ;;
        4)
            build_and_run
            ;;
        5)
            show_manual_instructions
            ;;
        *)
            echo -e "${RED}Opção inválida!${NC}"
            exit 1
            ;;
    esac
}

# Execute if script is run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi