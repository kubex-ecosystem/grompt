#!/usr/bin/env bash

# ╔══════════════════════════════════════════════════════════════╗
# ║                    🚀 Grompt Setup                  ║
# ║                                                              ║
# ║              Script automático de configuração              ║
# ║                     Powered by Shell + Go                   ║
# ╚══════════════════════════════════════════════════════════════╝

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Project info
PROJECT_NAME="grompt"
PROJECT_VERSION="1.0.0"
AUTHOR_NAME="Grompt Team"

# Functions
print_banner() {
    echo -e "${PURPLE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    🚀 Grompt Setup                  ║"
    echo "║                                                              ║"
    echo "║              Gerando estrutura completa do projeto          ║"
    echo "║                     Powered by Shell + Go                   ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    echo
}

print_step() {
    echo -e "${CYAN}📁 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
    exit 1
}

# Check dependencies
check_dependencies() {
    print_step "Verificando dependências..."
    
    if ! command -v node &> /dev/null; then
        print_warning "Node.js não encontrado. Certifique-se de instalar Node.js 16+ para desenvolvimento."
    fi
    
    if ! command -v npm &> /dev/null; then
        print_warning "npm não encontrado. Será necessário para instalar dependências do frontend."
    fi
    
    if ! command -v go &> /dev/null; then
        print_error "Go não encontrado. Por favor, instale Go 1.21+ antes de continuar."
    fi
    
    print_success "Verificação de dependências concluída"
}

# Create project structure
create_structure() {
    print_step "Criando estrutura do projeto..."
    
    # Create main directories
    mkdir -p "$PROJECT_NAME"
    cd "$PROJECT_NAME"
    
    # Create frontend structure
    mkdir -p frontend/src
    mkdir -p frontend/public
    mkdir -p .github/workflows
    
    print_success "Estrutura de diretórios criada"
}

# Create Go files
create_go_files() {
    print_step "Gerando arquivos Go..."
    
    # main.go
    cat > main.go << 'EOF'
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	AppName     = "Grompt"
	AppVersion  = "1.0.0"
	DefaultPort = "8080"
)

func main() {
	printBanner()

	// Configuração
	cfg := &Config{
		Port:           getEnvOr("PORT", DefaultPort),
		ClaudeAPIKey:   os.Getenv("CLAUDE_API_KEY"),
		OllamaEndpoint: getEnvOr("OLLAMA_ENDPOINT", "http://localhost:11434"),
	}

	// Inicializar servidor
	server := NewServer(cfg)

	// Graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		fmt.Println("\n🛑 Encerrando servidor...")
		server.Shutdown()
		os.Exit(0)
	}()

	// Iniciar servidor
	if err := server.Start(); err != nil {
		log.Fatal("❌ Erro ao iniciar servidor:", err)
	}
}

func printBanner() {
	fmt.Printf(`
╔══════════════════════════════════════════════════════════════╗
║                    🚀 %s v%s                    ║
║                                                              ║
║              Transforme ideias em prompts estruturados      ║
║                     Powered by Go + React                   ║
╚══════════════════════════════════════════════════════════════╝

`, AppName, AppVersion)
}

func getEnvOr(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
EOF

    # config.go
    cat > config.go << 'EOF'
package main

type Config struct {
	Port           string
	ClaudeAPIKey   string
	OllamaEndpoint string
	Debug          bool
}

type APIConfig struct {
	ClaudeAvailable bool   `json:"claude_available"`
	OllamaAvailable bool   `json:"ollama_available"`
	DemoMode        bool   `json:"demo_mode"`
	Version         string `json:"version"`
}

func (c *Config) GetAPIConfig() *APIConfig {
	return &APIConfig{
		ClaudeAvailable: c.ClaudeAPIKey != "",
		OllamaAvailable: c.checkOllamaConnection(),
		DemoMode:        true,
		Version:         AppVersion,
	}
}

func (c *Config) checkOllamaConnection() bool {
	// Implementar verificação de conexão com Ollama
	// Por simplicidade, retorna false por enquanto
	return false
}
EOF

    # server.go
    cat > server.go << 'EOF'
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

//go:embed build/*
var reactApp embed.FS

type Server struct {
	config   *Config
	handlers *Handlers
}

func NewServer(cfg *Config) *Server {
	handlers := NewHandlers(cfg)
	return &Server{
		config:   cfg,
		handlers: handlers,
	}
}

func (s *Server) Start() error {
	// Configurar roteamento
	s.setupRoutes()

	url := fmt.Sprintf("http://localhost:%s", s.config.Port)
	
	fmt.Printf("🌐 Servidor iniciado em: %s\n", url)
	fmt.Printf("📁 Servindo aplicação React embarcada\n")
	fmt.Printf("🔧 APIs disponíveis:\n")
	fmt.Printf("   • /api/config - Configuração\n")
	fmt.Printf("   • /api/claude - Claude API\n")
	fmt.Printf("   • /api/ollama - Ollama Local\n")
	fmt.Printf("💡 Pressione Ctrl+C para parar\n\n")

	// Abrir navegador após delay
	go func() {
		time.Sleep(1 * time.Second)
		openBrowser(url)
	}()

	return http.ListenAndServe(":"+s.config.Port, nil)
}

func (s *Server) setupRoutes() {
	// Servir aplicação React
	buildFS, err := fs.Sub(reactApp, "build")
	if err != nil {
		log.Fatal("Erro ao criar subfilesystem:", err)
	}

	fileServer := http.FileServer(http.FS(buildFS))

	// SPA routing - sempre servir index.html para rotas não encontradas
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			if _, err := fs.Stat(buildFS, r.URL.Path[1:]); os.IsNotExist(err) {
				r.URL.Path = "/"
			}
		}
		fileServer.ServeHTTP(w, r)
	})

	// API Routes
	http.HandleFunc("/api/config", s.handlers.HandleConfig)
	http.HandleFunc("/api/claude", s.handlers.HandleClaude)
	http.HandleFunc("/api/ollama", s.handlers.HandleOllama)
	http.HandleFunc("/api/health", s.handlers.HandleHealth)
}

func (s *Server) Shutdown() {
	fmt.Println("🧹 Limpando recursos...")
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		fmt.Printf("🌐 Abra seu navegador em: %s\n", url)
		return
	}

	if err != nil {
		fmt.Printf("⚠️  Não foi possível abrir o navegador automaticamente.\n")
		fmt.Printf("🌐 Abra manualmente: %s\n", url)
	}
}
EOF

    # handlers.go
    cat > handlers.go << 'EOF'
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Handlers struct {
	config    *Config
	claudeAPI *ClaudeAPI
	ollamaAPI *OllamaAPI
}

func NewHandlers(cfg *Config) *Handlers {
	return &Handlers{
		config:    cfg,
		claudeAPI: NewClaudeAPI(cfg.ClaudeAPIKey),
		ollamaAPI: NewOllamaAPI(cfg.OllamaEndpoint),
	}
}

func (h *Handlers) HandleConfig(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)
	
	if r.Method == "OPTIONS" {
		return
	}

	config := h.config.GetAPIConfig()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func (h *Handlers) HandleClaude(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)
	
	if r.Method == "OPTIONS" {
		return
	}
	
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req ClaudeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if h.config.ClaudeAPIKey == "" {
		http.Error(w, "Claude API Key não configurada", http.StatusServiceUnavailable)
		return
	}

	response, err := h.claudeAPI.Complete(req.Prompt, req.MaxTokens)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro na API Claude: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"response": response,
		"provider": "claude",
	})
}

func (h *Handlers) HandleOllama(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)
	
	if r.Method == "OPTIONS" {
		return
	}
	
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req OllamaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	response, err := h.ollamaAPI.Complete(req.Model, req.Prompt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro na API Ollama: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"response": response,
		"provider": "ollama",
	})
}

func (h *Handlers) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"version":   AppVersion,
		"apis": map[string]bool{
			"claude": h.config.ClaudeAPIKey != "",
			"ollama": h.ollamaAPI.IsAvailable(),
		},
	})
}

func (h *Handlers) setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
EOF

    # claude.go
    cat > claude.go << 'EOF'
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ClaudeAPI struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type ClaudeRequest struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type ClaudeAPIRequest struct {
	Model     string `json:"model"`
	MaxTokens int    `json:"max_tokens"`
	Messages  []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type ClaudeAPIResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func NewClaudeAPI(apiKey string) *ClaudeAPI {
	return &ClaudeAPI{
		apiKey:  apiKey,
		baseURL: "https://api.anthropic.com/v1/messages",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *ClaudeAPI) Complete(prompt string, maxTokens int) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("API key não configurada")
	}

	requestBody := ClaudeAPIRequest{
		Model:     "claude-3-sonnet-20240229",
		MaxTokens: maxTokens,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar request: %v", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro na requisição: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API retornou status %d: %s", resp.StatusCode, string(body))
	}

	var response ClaudeAPIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	if len(response.Content) == 0 {
		return "", fmt.Errorf("resposta vazia da API")
	}

	return response.Content[0].Text, nil
}
EOF

    # ollama.go
    cat > ollama.go << 'EOF'
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OllamaAPI struct {
	baseURL    string
	httpClient *http.Client
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func NewOllamaAPI(baseURL string) *OllamaAPI {
	return &OllamaAPI{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (o *OllamaAPI) Complete(model, prompt string) (string, error) {
	endpoint := fmt.Sprintf("%s/api/generate", o.baseURL)

	requestBody := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar request: %v", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro na requisição para Ollama: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("erro ao ler resposta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama retornou status %d: %s", resp.StatusCode, string(body))
	}

	var response OllamaResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	return response.Response, nil
}

func (o *OllamaAPI) IsAvailable() bool {
	endpoint := fmt.Sprintf("%s/api/tags", o.baseURL)
	
	resp, err := o.httpClient.Get(endpoint)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK
}
EOF

    # go.mod
    cat > go.mod << EOF
module grompt

go 1.21

require ()
EOF

    print_success "Arquivos Go criados"
}

# Create React files
create_react_files() {
    print_step "Gerando arquivos React..."
    
    # package.json
    cat > frontend/package.json << EOF
{
  "name": "grompt-frontend",
  "version": "1.0.0",
  "description": "Interface React para o Grompt - Ferramenta de Engenharia de Prompts",
  "private": true,
  "homepage": "./",
  "dependencies": {
    "@testing-library/jest-dom": "^5.17.0",
    "@testing-library/react": "^13.4.0",
    "@testing-library/user-event": "^14.5.2",
    "lucide-react": "^0.263.1",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-scripts": "5.0.1",
    "web-vitals": "^2.1.4"
  },
  "scripts": {
    "start": "react-scripts start",
    "build": "react-scripts build && npm run post-build",
    "post-build": "echo '✅ Build concluído! Arquivos prontos para embed no Go.'",
    "test": "react-scripts test",
    "eject": "react-scripts eject",
    "build:go": "npm run build && echo '📦 Build otimizado para integração Go criado em ./build/'",
    "analyze": "npm run build && npx bundle-analyzer build/static/js/*.js"
  },
  "eslintConfig": {
    "extends": [
      "react-app",
      "react-app/jest"
    ]
  },
  "browserslist": {
    "production": [
      ">0.2%",
      "not dead",
      "not op_mini all"
    ],
    "development": [
      "last 1 chrome version",
      "last 1 firefox version",
      "last 1 safari version"
    ]
  },
  "devDependencies": {
    "bundle-analyzer": "^0.1.0"
  },
  "engines": {
    "node": ">=16.0.0",
    "npm": ">=8.0.0"
  },
  "keywords": [
    "prompt-engineering",
    "ai",
    "claude",
    "react",
    "golang",
    "embedded"
  ],
  "author": "$AUTHOR_NAME",
  "license": "MIT"
}
EOF

    # public/index.html
    cat > frontend/public/index.html << 'EOF'
<!DOCTYPE html>
<html lang="pt-BR">
  <head>
    <meta charset="utf-8" />
    <link rel="icon" href="%PUBLIC_URL%/favicon.ico" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="theme-color" content="#000000" />
    <meta name="description" content="Grompt - Transforme ideias em prompts estruturados" />
    <script src="https://cdn.tailwindcss.com"></script>
    <title>Grompt</title>
  </head>
  <body>
    <noscript>Você precisa habilitar JavaScript para executar esta aplicação.</noscript>
    <div id="root"></div>
  </body>
</html>
EOF

    # src/index.js
    cat > frontend/src/index.js << 'EOF'
import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
EOF

    # src/index.css
    cat > frontend/src/index.css << 'EOF'
body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

code {
  font-family: source-code-pro, Menlo, Monaco, Consolas, 'Courier New',
    monospace;
}

.slider::-webkit-slider-thumb {
  appearance: none;
  height: 20px;
  width: 20px;
  border-radius: 50%;
  background: #3b82f6;
  cursor: pointer;
}

.slider::-moz-range-thumb {
  height: 20px;
  width: 20px;
  border-radius: 50%;
  background: #3b82f6;
  cursor: pointer;
}
EOF

    print_success "Arquivos React criados"
}

# Create configuration files
create_config_files() {
    print_step "Criando arquivos de configuração..."
    
    # Makefile
    cat > Makefile << 'EOF'
.PHONY: build-frontend build-backend build-all clean dev install-deps run help

# Variáveis
APP_NAME=grompt
VERSION=1.0.0
BUILD_DIR=build
FRONTEND_DIR=frontend

# Help
help:
	@echo "🚀 Grompt - Comandos Disponíveis:"
	@echo ""
	@echo "📦 Build:"
	@echo "  make install-deps     - Instalar dependências"
	@echo "  make build-frontend   - Build do React"
	@echo "  make build-backend    - Build do Go"
	@echo "  make build-all        - Build completo"
	@echo "  make build-cross      - Build multiplataforma"
	@echo ""
	@echo "🔧 Desenvolvimento:"
	@echo "  make dev              - Modo desenvolvimento"
	@echo "  make run              - Executar aplicação"
	@echo ""
	@echo "🧹 Limpeza:"
	@echo "  make clean            - Limpar builds"

# Instalar dependências
install-deps:
	@echo "📦 Instalando dependências do frontend..."
	cd $(FRONTEND_DIR) && npm install
	@echo "📦 Baixando módulos Go..."
	go mod tidy
	@echo "✅ Dependências instaladas!"

# Build do frontend React
build-frontend:
	@echo "⚛️  Compilando React..."
	cd $(FRONTEND_DIR) && npm run build
	@echo "📁 Copiando build para diretório raiz..."
	cp -r $(FRONTEND_DIR)/build ./
	@echo "✅ Frontend compilado!"

# Build do backend Go
build-backend:
	@echo "🐹 Compilando Go..."
	go build -ldflags="-s -w -X main.AppVersion=$(VERSION)" -o $(APP_NAME) .
	@echo "✅ Backend compilado!"

# Build completo
build-all: build-frontend build-backend
	@echo "🎉 Build completo finalizado!"
	@echo "📱 Execute: ./$(APP_NAME)"

# Build multiplataforma
build-cross: build-frontend
	@echo "🌍 Compilando para múltiplas plataformas..."
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(APP_NAME)-windows-amd64.exe .
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(APP_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(APP_NAME)-macos-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(APP_NAME)-macos-arm64 .
	@echo "✅ Builds multiplataforma concluídos!"

# Desenvolvimento
dev:
	@echo "🔧 Iniciando modo desenvolvimento..."
	cd $(FRONTEND_DIR) && npm start &
	sleep 3
	go run . --dev

# Executar aplicação
run: build-all
	@echo "🚀 Iniciando $(APP_NAME)..."
	./$(APP_NAME)

# Limpeza
clean:
	@echo "🧹 Limpando builds..."
	rm -rf $(BUILD_DIR)/
	rm -f $(APP_NAME)*
	cd $(FRONTEND_DIR) && rm -rf build/
	@echo "✅ Limpeza concluída!"
EOF

    # .gitignore
    cat > .gitignore << 'EOF'
# Binários
grompt*
!grompt/

# Build outputs
build/
dist/

# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
go.work

# Node.js
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*

# React
frontend/build/
frontend/.env.local
frontend/.env.development.local
frontend/.env.test.local
frontend/.env.production.local

# IDEs
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Logs
*.log
logs/

# Temporary files
tmp/
temp/

# Environment variables
.env
.env.local
EOF

    # LICENSE
    cat > LICENSE << EOF
MIT License

Copyright (c) 2024 $AUTHOR_NAME

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
EOF

    print_success "Arquivos de configuração criados"
}

# Create App.jsx (the main React component)
create_app_jsx() {
    print_step "Criando componente React principal..."
    
    cat > frontend/src/App.jsx << 'EOF'
import React, { useState, useEffect } from 'react';
import { Trash2, Edit3, Plus, Wand2, Sun, Moon, Copy, Check, AlertCircle } from 'lucide-react';

const Krompt = () => {
  const [darkMode, setDarkMode] = useState(true);
  const [currentInput, setCurrentInput] = useState('');
  const [ideas, setIdeas] = useState([]);
  const [editingId, setEditingId] = useState(null);
  const [editingText, setEditingText] = useState('');
  const [purpose, setPurpose] = useState('Outros');
  const [customPurpose, setCustomPurpose] = useState('');
  const [maxLength, setMaxLength] = useState(5000);
  const [generatedPrompt, setGeneratedPrompt] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [copied, setCopied] = useState(false);
  const [apiProvider, setApiProvider] = useState('claude');
  const [availableAPIs, setAvailableAPIs] = useState({
    claude_available: false,
    ollama_available: false,
    demo_mode: true
  });
  const [connectionStatus, setConnectionStatus] = useState('checking');

  useEffect(() => {
    document.documentElement.className = darkMode ? 'dark' : '';
  }, [darkMode]);

  // Verificar configuração e APIs disponíveis na inicialização
  useEffect(() => {
    checkAPIAvailability();
  }, []);

  const checkAPIAvailability = async () => {
    try {
      const response = await fetch('/api/config');
      if (response.ok) {
        const config = await response.json();
        setAvailableAPIs(config);
        setConnectionStatus('connected');
        
        // Definir provider padrão baseado na disponibilidade
        if (config.claude_available) {
          setApiProvider('claude');
        } else if (config.ollama_available) {
          setApiProvider('ollama');
        } else {
          setApiProvider('demo');
        }
      } else {
        throw new Error('Servidor não respondeu');
      }
    } catch (error) {
      console.error('Erro ao verificar APIs:', error);
      setConnectionStatus('offline');
      setAvailableAPIs({ demo_mode: true });
      setApiProvider('demo');
    }
  };

  const addIdea = () => {
    if (currentInput.trim()) {
      const newIdea = {
        id: Date.now(),
        text: currentInput.trim()
      };
      setIdeas([...ideas, newIdea]);
      setCurrentInput('');
    }
  };

  const removeIdea = (id) => {
    setIdeas(ideas.filter(idea => idea.id !== id));
  };

  const startEditing = (id, text) => {
    setEditingId(id);
    setEditingText(text);
  };

  const saveEdit = () => {
    setIdeas(ideas.map(idea => 
      idea.id === editingId 
        ? { ...idea, text: editingText }
        : idea
    ));
    setEditingId(null);
    setEditingText('');
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditingText('');
  };

  const generateDemoPrompt = () => {
    const purposeText = purpose === 'Outros' && customPurpose 
      ? customPurpose 
      : purpose;

    return `# Prompt Estruturado - ${purposeText}

## 🎯 Contexto
Você é um assistente especializado em **${purposeText.toLowerCase()}** com conhecimento profundo na área.

## 📝 Ideias do Usuário Organizadas:
${ideas.map((idea, index) => `**${index + 1}.** ${idea.text}`).join('\n')}

## 🔧 Instruções Específicas
- Analise cuidadosamente todas as ideias apresentadas acima
- Identifique o objetivo principal e objetivos secundários
- Forneça uma resposta estruturada e bem organizada
- Mantenha o foco no propósito definido: **${purposeText}**
- Use exemplos práticos quando apropriado
- Seja específico e actionável

## 📋 Formato de Resposta Esperado
1. **Análise Inicial**: Resumo do que foi solicitado
2. **Desenvolvimento**: Resposta detalhada seguindo as ideias
3. **Conclusão**: Próximos passos ou considerações finais

## ⚙️ Configurações Técnicas
- Máximo de caracteres: ${maxLength.toLocaleString()}
- Propósito: ${purposeText}
- Total de ideias processadas: ${ideas.length}

---
*Prompt gerado automaticamente pelo Grompt v1.0*
*Modo: Demo (configure APIs para funcionalidade completa)*`;
  };

  const generatePrompt = async () => {
    if (ideas.length === 0) return;
    
    setIsGenerating(true);
    
    const purposeText = purpose === 'Outros' && customPurpose 
      ? customPurpose 
      : purpose;
    
    const engineeringPrompt = `
Você é um especialista em engenharia de prompts com conhecimento profundo em técnicas de prompt engineering. Sua tarefa é transformar ideias brutas e desorganizadas em um prompt estruturado, profissional e eficaz.

CONTEXTO: O usuário inseriu as seguintes notas/ideias brutas:
${ideas.map((idea, index) => `${index + 1}. "${idea.text}"`).join('\n')}

PROPÓSITO DO PROMPT: ${purposeText}
TAMANHO MÁXIMO: ${maxLength} caracteres

INSTRUÇÕES PARA ESTRUTURAÇÃO:
1. Analise todas as ideias e identifique o objetivo principal
2. Organize as informações de forma lógica e hierárquica
3. Aplique técnicas de engenharia de prompt como:
   - Definição clara de contexto e papel
   - Instruções específicas e mensuráveis
   - Exemplos quando apropriado
   - Formato de saída bem definido
   - Chain-of-thought se necessário
4. Use markdown para estruturação clara
5. Seja preciso, objetivo e profissional
6. Mantenha o escopo dentro do limite de caracteres

IMPORTANTE: Responda APENAS com o prompt estruturado em markdown, sem explicações adicionais ou texto introdutório. O prompt deve ser completo e pronto para uso.
`;

    try {
      let response;
      
      if (apiProvider === 'demo') {
        // Simular delay para parecer real
        await new Promise(resolve => setTimeout(resolve, 2000));
        response = generateDemoPrompt();
      } else if (apiProvider === 'claude') {
        const result = await fetch('/api/claude', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            prompt: engineeringPrompt,
            max_tokens: maxLength
          })
        });
        
        if (!result.ok) {
          throw new Error(`Erro HTTP: ${result.status}`);
        }
        
        const data = await result.json();
        response = data.response || data.content || 'Resposta vazia do servidor';
      } else if (apiProvider === 'ollama') {
        const result = await fetch('/api/ollama', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            model: 'llama2',
            prompt: engineeringPrompt,
            stream: false
          })
        });
        
        if (!result.ok) {
          throw new Error(`Erro HTTP: ${result.status}`);
        }
        
        const data = await result.json();
        response = data.response || 'Resposta vazia do Ollama';
      }
      
      setGeneratedPrompt(response);
    } catch (error) {
      console.error('Erro ao gerar prompt:', error);
      setGeneratedPrompt(`# Erro ao Gerar Prompt

**Erro:** ${error.message}

**Detalhes:** Não foi possível conectar com a API selecionada. Verifique:
- Se o servidor está rodando
- Se a API está configurada corretamente
- Se há conexão com a internet (para APIs externas)

**Dica:** Tente usar o modo demo ou configure uma API diferente.`);
    }
    
    setIsGenerating(false);
  };

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(generatedPrompt);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (error) {
      console.error('Erro ao copiar:', error);
      // Fallback para navegadores mais antigos
      const textArea = document.createElement('textarea');
      textArea.value = generatedPrompt;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const theme = {
    dark: {
      bg: 'bg-gray-900',
      cardBg: 'bg-gray-800',
      text: 'text-gray-100',
      textSecondary: 'text-gray-300',
      border: 'border-gray-700',
      input: 'bg-gray-700 border-gray-600 text-gray-100',
      button: 'bg-blue-600 hover:bg-blue-700 text-white',
      buttonSecondary: 'bg-gray-700 hover:bg-gray-600 text-gray-200',
      accent: 'text-blue-400'
    },
    light: {
      bg: 'bg-gray-50',
      cardBg: 'bg-white',
      text: 'text-gray-900',
      textSecondary: 'text-gray-600',
      border: 'border-gray-300',
      input: 'bg-white border-gray-300 text-gray-900',
      button: 'bg-blue-600 hover:bg-blue-700 text-white',
      buttonSecondary: 'bg-gray-200 hover:bg-gray-300 text-gray-700',
      accent: 'text-blue-600'
    }
  };

  const currentTheme = darkMode ? theme.dark : theme.light;

  const getConnectionStatusColor = () => {
    switch (connectionStatus) {
      case 'connected': return 'text-green-500';
      case 'offline': return 'text-red-500';
      default: return 'text-yellow-500';
    }
  };

  const getConnectionStatusText = () => {
    switch (connectionStatus) {
      case 'connected': return 'Conectado';
      case 'offline': return 'Offline (Modo Demo)';
      default: return 'Verificando...';
    }
  };

  return (
    <div className={`min-h-screen ${currentTheme.bg} ${currentTheme.text} p-4 transition-colors duration-300`}>
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-4xl font-bold mb-2">
              <span className={currentTheme.accent}>Prompt</span> Crafter
            </h1>
            <p className={currentTheme.textSecondary}>
              Transforme suas ideias brutas em prompts estruturados e profissionais
            </p>
          </div>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <div className={`h-2 w-2 rounded-full ${connectionStatus === 'connected' ? 'bg-green-500' : connectionStatus === 'offline' ? 'bg-red-500' : 'bg-yellow-500'}`}></div>
              <span className={`text-sm ${getConnectionStatusColor()}`}>
                {getConnectionStatusText()}
              </span>
            </div>
            <select 
              value={apiProvider}
              onChange={(e) => setApiProvider(e.target.value)}
              className={`px-3 py-2 rounded-lg ${currentTheme.input} border focus:ring-2 focus:ring-blue-500`}
            >
              {availableAPIs.claude_available && (
                <option value="claude">Claude API</option>
              )}
              {availableAPIs.ollama_available && (
                <option value="ollama">Ollama Local</option>
              )}
              <option value="demo">Modo Demo</option>
            </select>
            <button
              onClick={() => setDarkMode(!darkMode)}
              className={`p-2 rounded-lg ${currentTheme.buttonSecondary} transition-colors`}
            >
              {darkMode ? <Sun size={20} /> : <Moon size={20} />}
            </button>
          </div>
        </div>

        {/* Status Alert */}
        {connectionStatus === 'offline' && (
          <div className="mb-6 p-4 bg-yellow-900 border border-yellow-600 rounded-lg flex items-center gap-3">
            <AlertCircle className="text-yellow-400" size={20} />
            <p className="text-yellow-100">
              <strong>Modo Offline:</strong> Executando em modo demo. Configure APIs para funcionalidade completa.
            </p>
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Input Section */}
          <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`}>
            <h2 className="text-xl font-semibold mb-4">📝 Adicionar Ideias</h2>
            <div className="space-y-4">
              <textarea
                value={currentInput}
                onChange={(e) => setCurrentInput(e.target.value)}
                placeholder="Cole suas notas, ideias brutas ou pensamentos desorganizados aqui..."
                className={`w-full h-32 px-4 py-3 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-blue-500 resize-none`}
                onKeyDown={(e) => {
                  if (e.key === 'Enter' && e.ctrlKey) {
                    addIdea();
                  }
                }}
              />
              <button
                onClick={addIdea}
                disabled={!currentInput.trim()}
                className={`w-full flex items-center justify-center gap-2 px-4 py-3 rounded-lg ${currentTheme.button} disabled:opacity-50 disabled:cursor-not-allowed transition-all`}
              >
                <Plus size={20} />
                Incluir (Ctrl+Enter)
              </button>
            </div>

            {/* Configuration */}
            <div className="mt-6 space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">Propósito do Prompt</label>
                <div className="space-y-2">
                  <div className="flex gap-2">
                    {['Código', 'Imagem', 'Outros'].map((option) => (
                      <button
                        key={option}
                        onClick={() => setPurpose(option)}
                        className={`px-3 py-2 rounded-lg text-sm border transition-colors ${
                          purpose === option 
                            ? 'bg-blue-600 text-white border-blue-600' 
                            : `${currentTheme.buttonSecondary} ${currentTheme.border}`
                        }`}
                      >
                        {option}
                      </button>
                    ))}
                  </div>
                  {purpose === 'Outros' && (
                    <input
                      type="text"
                      value={customPurpose}
                      onChange={(e) => setCustomPurpose(e.target.value)}
                      placeholder="Descreva o objetivo do prompt..."
                      className={`w-full px-3 py-2 rounded-lg border ${currentTheme.input} focus:ring-2 focus:ring-blue-500`}
                    />
                  )}
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">
                  Tamanho Máximo: {maxLength.toLocaleString()} caracteres
                </label>
                <input
                  type="range"
                  min="500"
                  max="130000"
                  step="500"
                  value={maxLength}
                  onChange={(e) => setMaxLength(parseInt(e.target.value))}
                  className="w-full h-2 bg-gray-300 rounded-lg appearance-none cursor-pointer slider"
                />
              </div>
            </div>
          </div>

          {/* Ideas List */}
          <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg`}>
            <h2 className="text-xl font-semibold mb-4">💡 Suas Ideias ({ideas.length})</h2>
            <div className="space-y-3 max-h-96 overflow-y-auto">
              {ideas.length === 0 ? (
                <p className={`${currentTheme.textSecondary} text-center py-8`}>
                  Adicione suas primeiras ideias ao lado ←
                </p>
              ) : (
                ideas.map((idea) => (
                  <div key={idea.id} className={`p-3 rounded-lg border ${currentTheme.border} bg-opacity-50`}>
                    {editingId === idea.id ? (
                      <div className="space-y-2">
                        <textarea
                          value={editingText}
                          onChange={(e) => setEditingText(e.target.value)}
                          className={`w-full px-2 py-1 rounded border ${currentTheme.input} text-sm`}
                          rows="2"
                        />
                        <div className="flex gap-1">
                          <button
                            onClick={saveEdit}
                            className="px-2 py-1 bg-green-600 text-white rounded text-xs hover:bg-green-700"
                          >
                            Salvar
                          </button>
                          <button
                            onClick={cancelEdit}
                            className={`px-2 py-1 rounded text-xs ${currentTheme.buttonSecondary}`}
                          >
                            Cancelar
                          </button>
                        </div>
                      </div>
                    ) : (
                      <>
                        <p className="text-sm mb-2">{idea.text}</p>
                        <div className="flex justify-end gap-1">
                          <button
                            onClick={() => startEditing(idea.id, idea.text)}
                            className={`p-1 rounded ${currentTheme.buttonSecondary} hover:bg-opacity-80`}
                          >
                            <Edit3 size={14} />
                          </button>
                          <button
                            onClick={() => removeIdea(idea.id)}
                            className="p-1 rounded bg-red-600 text-white hover:bg-red-700"
                          >
                            <Trash2 size={14} />
                          </button>
                        </div>
                      </>
                    )}
                  </div>
                ))
              )}
            </div>
            
            {ideas.length > 0 && (
              <button
                onClick={generatePrompt}
                disabled={isGenerating}
                className={`w-full mt-4 flex items-center justify-center gap-2 px-4 py-3 rounded-lg bg-gradient-to-r from-purple-600 to-blue-600 text-white hover:from-purple-700 hover:to-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all transform hover:scale-105`}
              >
                <Wand2 size={20} className={isGenerating ? 'animate-spin' : ''} />
                {isGenerating ? 'Gerando...' : 'Me ajude, engenheiro?!'}
              </button>
            )}
          </div>

          {/* Generated Prompt */}
          <div className={`${currentTheme.cardBg} rounded-xl p-6 border ${currentTheme.border} shadow-lg ${generatedPrompt ? 'lg:col-span-1' : ''}`}>
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-semibold">🚀 Prompt Estruturado</h2>
              {generatedPrompt && (
                <button
                  onClick={copyToClipboard}
                  className={`flex items-center gap-2 px-3 py-2 rounded-lg ${currentTheme.buttonSecondary} hover:bg-opacity-80 transition-colors`}
                >
                  {copied ? <Check size={16} /> : <Copy size={16} />}
                  {copied ? 'Copiado!' : 'Copiar'}
                </button>
              )}
            </div>
            
            {generatedPrompt ? (
              <div className="space-y-4">
                <div className={`text-xs ${currentTheme.textSecondary} flex justify-between`}>
                  <span>Caracteres: {generatedPrompt.length}</span>
                  <span>Limite: {maxLength.toLocaleString()}</span>
                </div>
                <div className={`max-h-96 overflow-y-auto p-4 rounded-lg border ${currentTheme.border} bg-opacity-50`}>
                  <pre className="whitespace-pre-wrap text-sm font-mono">{generatedPrompt}</pre>
                </div>
              </div>
            ) : (
              <div className={`${currentTheme.textSecondary} text-center py-12`}>
                <Wand2 size={48} className="mx-auto mb-4 opacity-50" />
                <p>Seu prompt estruturado aparecerá aqui</p>
                <p className="text-sm mt-2">Adicione ideias e clique em "Me ajude, engenheiro?!"</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Krompt;
EOF

    print_success "Componente React principal criado"
}

# Create README
create_readme() {
    print_step "Criando documentação..."
    
    cat > README.md << EOF
# 🚀 Grompt

> Transforme suas ideias brutas em prompts estruturados e profissionais

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![React](https://img.shields.io/badge/React-18+-blue.svg)](https://reactjs.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## ✨ Características

- 🧠 **Engenharia de Prompts Real** - Aplica técnicas genuínas de prompt engineering
- ⚛️ **Interface React Moderna** - UI responsiva e intuitiva
- 🐹 **Backend Go Robusto** - Servidor HTTP eficiente e leve
- 📦 **Binário Único** - Zero dependências no deploy
- 🌍 **Multiplataforma** - Windows, Linux, macOS
- 🔌 **APIs Integradas** - Claude, Ollama e modo demo
- 🎨 **Temas Dark/Light** - Interface personalizável

## 🚀 Início Rápido

### Pré-requisitos
- Go 1.21+
- Node.js 16+
- npm ou yarn

### Instalação

\`\`\`bash
# 1. Clone o repositório
git clone https://github.com/seu-usuario/grompt
cd grompt

# 2. Instale dependências
make install-deps

# 3. Build completo
make build-all

# 4. Execute
./grompt
\`\`\`

A aplicação abrirá automaticamente em \`http://localhost:8080\`

## 🔧 Desenvolvimento

\`\`\`bash
# Modo desenvolvimento (hot reload)
make dev

# Build apenas frontend
make build-frontend

# Build apenas backend
make build-backend

# Build multiplataforma
make build-cross
\`\`\`

## ⚙️ Configuração

### Variáveis de Ambiente

\`\`\`bash
# Porta do servidor (padrão: 8080)
export PORT=3000

# Claude API Key (opcional)
export CLAUDE_API_KEY=your_claude_api_key

# Ollama Endpoint (padrão: http://localhost:11434)
export OLLAMA_ENDPOINT=http://localhost:11434
\`\`\`

### APIs Suportadas

- **Claude API** - Configure \`CLAUDE_API_KEY\`
- **Ollama Local** - Instale Ollama localmente
- **Modo Demo** - Funciona sem configuração

## 📁 Estrutura do Projeto

\`\`\`
grompt/
├── 📁 frontend/          # Aplicação React
│   ├── public/
│   ├── src/
│   └── package.json
├── 📄 main.go           # Entrada principal
├── 📄 server.go         # Servidor HTTP
├── 📄 handlers.go       # Manipuladores de rotas
├── 📄 claude.go         # Integração Claude
├── 📄 ollama.go         # Integração Ollama
├── 📄 config.go         # Configurações
├── 📄 Makefile          # Scripts de build
└── 📄 README.md
\`\`\`

## 🎯 Como Usar

1. **Adicione Ideias** - Cole suas notas brutas no primeiro campo
2. **Configure Propósito** - Escolha entre Código, Imagem ou Outros
3. **Ajuste Tamanho** - Define limite de caracteres (500-130k)
4. **Gere Prompt** - Clique em "Me ajude, engenheiro?!"
5. **Copie Resultado** - Use o prompt estruturado gerado

## 🔌 Integrações

### Claude API
\`\`\`bash
export CLAUDE_API_KEY=your_api_key
./grompt
\`\`\`

### Ollama Local
\`\`\`bash
# Instalar Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# Baixar modelo
ollama pull llama2

# Executar Grompt
./grompt
\`\`\`

## 📦 Distribuição

\`\`\`bash
# Build para produção
make build-cross

# Arquivos gerados:
# grompt-windows-amd64.exe
# grompt-linux-amd64
# grompt-macos-amd64
# grompt-macos-arm64
\`\`\`

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch (\`git checkout -b feature/nova-funcionalidade\`)
3. Commit suas mudanças (\`git commit -am 'Adiciona nova funcionalidade'\`)
4. Push para a branch (\`git push origin feature/nova-funcionalidade\`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🙏 Agradecimentos

- [Anthropic](https://anthropic.com) pela API Claude
- [Ollama](https://ollama.ai) pela plataforma de LLMs locais
- [React](https://reactjs.org) pela biblioteca de UI
- [Go](https://golang.org) pela linguagem robusta

---

<div align="center">
Feito com ❤️ em Go + React
</div>
EOF

    print_success "Documentação criada"
}

# Create GitHub Actions
create_github_actions() {
    print_step "Configurando GitHub Actions..."
    
    cat > .github/workflows/build.yml << 'EOF'
name: Build and Test

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
        
    - name: Set up Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
        
    - name: Install frontend dependencies
      run: cd frontend && npm install
      
    - name: Build frontend
      run: cd frontend && npm run build
      
    - name: Copy build to root
      run: cp -r frontend/build ./
      
    - name: Build Go application
      run: go build -v ./...
      
    - name: Test
      run: go test -v ./...
EOF

    print_success "GitHub Actions configurado"
}

# Initialize git repository
init_git() {
    print_step "Inicializando repositório Git..."
    
    git init
    git add .
    git commit -m "🚀 Initial commit: Grompt v$PROJECT_VERSION

✨ Features:
- Interface React moderna com temas dark/light
- Backend Go com servidor HTTP embarcado
- Integração Claude API e Ollama
- Engenharia de prompts real
- Binário único sem dependências
- Build multiplataforma

🔧 Tech Stack:
- Go 1.21+ (backend)
- React 18+ (frontend)
- TailwindCSS (styling)
- Embedded filesystem (go:embed)

📦 Zero dependencies deployment
🌍 Cross-platform binary"

    print_success "Repositório Git inicializado"
}

# Test build
test_build() {
    print_step "Testando build inicial..."
    
    if command -v npm &> /dev/null; then
        echo "🔧 Instalando dependências do npm..."
        cd frontend && npm install && cd ..
        
        echo "⚛️  Fazendo build do React..."
        cd frontend && npm run build && cd ..
        
        echo "📁 Copiando build..."
        cp -r frontend/build ./
    else
        print_warning "npm não encontrado. Pulando build do frontend."
    fi
    
    echo "📦 Inicializando módulos Go..."
    go mod tidy
    
    echo "🐹 Testando build Go..."
    go build -o grompt-test .
    
    if [ -f "grompt-test" ]; then
        print_success "Build teste bem-sucedido!"
        rm grompt-test
    else
        print_warning "Build teste falhou, mas estrutura foi criada."
    fi
}

# Final instructions
show_final_instructions() {
    echo
    print_step "🎉 Projeto Grompt criado com sucesso!"
    echo
    echo -e "${GREEN}📁 Estrutura criada em: ${BLUE}./$PROJECT_NAME/${NC}"
    echo
    echo -e "${CYAN}🚀 Próximos passos:${NC}"
    echo -e "   ${YELLOW}1.${NC} cd $PROJECT_NAME"
    echo -e "   ${YELLOW}2.${NC} make install-deps    # Instalar dependências"
    echo -e "   ${YELLOW}3.${NC} make build-all       # Build completo"
    echo -e "   ${YELLOW}4.${NC} ./grompt     # Executar aplicação"
    echo
    echo -e "${CYAN}🔧 Comandos úteis:${NC}"
    echo -e "   ${YELLOW}•${NC} make help            # Ver todos os comandos"
    echo -e "   ${YELLOW}•${NC} make dev             # Modo desenvolvimento"
    echo -e "   ${YELLOW}•${NC} make build-cross     # Build multiplataforma"
    echo -e "   ${YELLOW}•${NC} make clean           # Limpar builds"
    echo
    echo -e "${CYAN}📚 Configuração de APIs:${NC}"
    echo -e "   ${YELLOW}•${NC} export CLAUDE_API_KEY=your_key"
    echo -e "   ${YELLOW}•${NC} export OLLAMA_ENDPOINT=http://localhost:11434"
    echo
    echo -e "${CYAN}🐙 Para subir no GitHub:${NC}"
    echo -e "   ${YELLOW}1.${NC} Crie um novo repo no GitHub"
    echo -e "   ${YELLOW}2.${NC} git remote add origin https://github.com/SEU_USUARIO/grompt.git"
    echo -e "   ${YELLOW}3.${NC} git branch -M main"
    echo -e "   ${YELLOW}4.${NC} git push -u origin main"
    echo
    print_success "Estrutura completa gerada! Bom desenvolvimento! 🚀"
}

# Main execution
main() {
    print_banner
    
    # Check if project already exists
    if [ -d "$PROJECT_NAME" ]; then
        print_error "Diretório '$PROJECT_NAME' já existe! Remova-o ou escolha outro nome."
    fi
    
    check_dependencies
    create_structure
    create_go_files
    create_react_files
    create_config_files
    create_app_jsx
    create_readme
    create_github_actions
    init_git
    test_build
    show_final_instructions
}

# Script execution
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
