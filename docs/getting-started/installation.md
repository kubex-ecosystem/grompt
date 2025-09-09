# Instalação

Esta página fornece instruções detalhadas para instalar e configurar o Grompt em diferentes plataformas.

## 📦 Opções de Instalação

### Opção 1: Download Binário (Recomendado)

A forma mais simples de instalar o Grompt é baixando o binário pré-compilado para sua plataforma:

#### Windows

```powershell
# PowerShell
Invoke-WebRequest -Uri "https://github.com/kubex-ecosystem/grompt/releases/latest/download/grompt_windows_amd64.exe" -OutFile "grompt.exe"
.\grompt.exe --help
```

#### Linux

```bash
# Download e instalação
curl -L https://github.com/kubex-ecosystem/grompt/releases/latest/download/grompt_linux_amd64 -o grompt
chmod +x grompt
sudo mv grompt /usr/local/bin/
```

#### macOS

```bash
# macOS Intel
curl -L https://github.com/kubex-ecosystem/grompt/releases/latest/download/grompt_darwin_amd64 -o grompt

# macOS Apple Silicon (M1/M2)
curl -L https://github.com/kubex-ecosystem/grompt/releases/latest/download/grompt_darwin_arm64 -o grompt

# Tornar executável e mover para PATH
chmod +x grompt
sudo mv grompt /usr/local/bin/
```

### Opção 2: Instalar via Make

```bash
git clone https://github.com/kubex-ecosystem/grompt
cd grompt
make install
```

Este comando irá:

1. Compilar o binário para sua plataforma
2. Instalar em `/usr/local/bin/grompt`
3. Tornar disponível globalmente

### Opção 3: Compilar do Código Fonte

#### Pré-requisitos

- **Go 1.25+** - [Instalar Go](https://golang.org/doc/install)
- **Node.js 18+** - [Instalar Node.js](https://nodejs.org/)
- **Make** - Disponível na maioria dos sistemas Unix

#### Passos de Compilação

```bash
# 1. Clonar o repositório
git clone https://github.com/kubex-ecosystem/grompt.git
cd grompt

# 2. Compilar
make build

# 3. Executar
./dist/grompt --help
```

#### Compilação para Outras Plataformas

```bash
# Compilar para Windows
make build-windows

# Compilar para Linux
make build-linux

# Compilar para macOS
make build-darwin

# Compilar para todas as plataformas
make build-all
```

## ⚙️ Configuração Inicial

### 1. Verificar Instalação

```bash
grompt --version
```

### 2. Configurar Variáveis de Ambiente (Opcional)

O Grompt funciona em modo demo sem configuração, mas para usar provedores de IA externos, configure as chaves de API:

```bash
# Adicione ao seu ~/.bashrc, ~/.zshrc, ou ~/.profile

# OpenAI
export OPENAI_API_KEY="sk-..."

# Claude (Anthropic)
export CLAUDE_API_KEY="sk-ant-..."

# DeepSeek
export DEEPSEEK_API_KEY="..."

# Gemini
export GEMINI_API_KEY="..."

# Ollama (local)
export OLLAMA_ENDPOINT="http://localhost:11434"

# Configurações do servidor (opcional)
export PORT=8080
export DEBUG=false
```

### 3. Primeiro Teste

```bash
# Testar em modo demo (sem API keys)
grompt

# Testar CLI
grompt ask "Olá mundo!" --provider demo
```

## 🔧 Configuração Avançada

### Configuração do Servidor

Por padrão, o Grompt roda na porta 8080. Para alterar:

```bash
export PORT=3000
grompt
```

Ou diretamente:

```bash
grompt --port 3000
```

### Configuração de Debug

Para habilitação de logs detalhados:

```bash
export DEBUG=true
grompt
```

### Configuração para Ollama Local

Se você tem o Ollama instalado localmente:

```bash
# Instalar Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Baixar um modelo
ollama pull llama2

# Configurar endpoint
export OLLAMA_ENDPOINT="http://localhost:11434"
export OLLAMA_MODEL="llama2"
```

## 🛠️ Solução de Problemas

### Problemas Comuns

#### "Permission denied" no Linux/macOS

```bash
chmod +x grompt
```

#### "grompt: command not found"

Certifique-se que o binário está no PATH:

```bash
echo $PATH
which grompt
```

#### Porta já em uso

```bash
# Verificar qual processo usa a porta
lsof -i :8080

# Usar porta diferente
grompt --port 8081
```

#### Problemas de Firewall

```bash
# Linux: permitir porta no firewall
sudo ufw allow 8080

# macOS: permitir no firewall do sistema
# Vá em System Preferences > Security & Privacy > Firewall
```

### Logs de Debug

```bash
DEBUG=true grompt
```

### Testar Conectividade

```bash
# Testar se o servidor está rodando
curl http://localhost:8080/api/health

# Testar provedores de IA
grompt ask "teste" --provider openai --dry-run
```

## 📋 Requisitos do Sistema

| Sistema | Requisitos Mínimos |
|---------|-------------------|
| **Memória RAM** | 100 MB |
| **Espaço em Disco** | 50 MB |
| **Processador** | x86_64 ou ARM64 |
| **Sistema Operacional** | Linux, macOS, Windows |
| **Rede** | Conectividade com internet (para provedores de IA externos) |

## 🔄 Atualizações

### Verificar Versão Atual

```bash
grompt --version
```

### Atualizar para Nova Versão

```bash
# Download manual
curl -L https://github.com/kubex-ecosystem/grompt/releases/latest/download/grompt_linux_amd64 -o grompt-new
chmod +x grompt-new
sudo mv grompt-new /usr/local/bin/grompt

# Ou recompilar do código
cd grompt
git pull
make build
sudo cp dist/grompt /usr/local/bin/
```

---

## 📚 Próximos Passos

- **[Início Rápido](quickstart.md)** - Primeiros passos com o Grompt
- **[Comandos CLI](../user-guide/cli-commands.md)** - Referência completa dos comandos
- **[Configuração](../user-guide/configuration.md)** - Configuração detalhada dos provedores de IA
