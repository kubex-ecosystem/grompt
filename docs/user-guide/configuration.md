# Configuração

Este guia detalha como configurar o Grompt para trabalhar com diferentes provedores de IA, personalizar o ambiente e otimizar seu fluxo de trabalho.

## 🧭 Matriz de Provedores (Cheat Sheet)

| Provedor | Variável | Padrão do Endpoint/Chave | Modelo padrão (endpoints) | Exemplos de modelos |
|---|---|---|---|---|
| OpenAI | `OPENAI_API_KEY` | chave secreta | `gpt-4o-mini` | `gpt-4o`, `gpt-4`, `gpt-3.5-turbo` |
| Claude | `CLAUDE_API_KEY` | chave secreta | `claude-3-5-sonnet-20241022` | `claude-3-sonnet-20240229`, `claude-3-5-haiku-20241022` |
| DeepSeek | `DEEPSEEK_API_KEY` | chave secreta | `deepseek-chat` | `deepseek-coder`, `deepseek-math`, `deepseek-reasoner` |
| Gemini | `GEMINI_API_KEY` | chave secreta | `gemini-1.5-flash` | `gemini-2.0-flash`, `gemini-2.0-pro`, `gemini-1.5-pro` |
| ChatGPT | `CHATGPT_API_KEY` | chave secreta | `gpt-4o-mini` | `gpt-4o`, `gpt-4`, `gpt-3.5-turbo` |
| Ollama | `OLLAMA_ENDPOINT` | `http://localhost:11434` | `llama3.2` | `mistral`, `codellama`, `llama3.1`, `llama2` |

Observações:
- Configure as variáveis no ambiente ou em `~/.gromptrc` (permissões 600).
- Endpoints HTTP do servidor: `PORT` (padrão `8080`) e `BindAddr` (padrão `localhost`).

## 🔧 Configuração de Provedores de IA

### OpenAI (GPT)

#### Obter Chave de API

1. Acesse [platform.openai.com](https://platform.openai.com)
2. Faça login ou crie uma conta
3. Vá para "API Keys" no menu
4. Clique em "Create new secret key"
5. Copie a chave (começa com `sk-`)

#### Configurar

```bash
# Via variável de ambiente (recomendado)
export OPENAI_API_KEY="sk-proj-..."

# Ou no arquivo de configuração
echo 'OPENAI_API_KEY=sk-proj-...' >> ~/.gromptrc
```

#### Modelos Disponíveis

| Modelo | Descrição | Uso Recomendado |
|--------|-----------|-----------------|
| `gpt-4` | Mais capaz, mais lento | Tarefas complexas, código crítico |
| `gpt-4-turbo` | Balanceado | Uso geral, boa velocidade |
| `gpt-3.5-turbo` | Rápido, econômico | Tarefas simples, prototipagem |

#### Exemplo de Uso

```bash
# Pergunta simples
grompt ask "Como implementar cache em Redis?" \
  --provider openai \
  --model gpt-4o-mini

# Geração de prompt
grompt generate \
  --idea "API de pagamentos" \
  --idea "Stripe integration" \
  --provider openai \
  --model gpt-4-turbo
```

### Claude (Anthropic)

#### Obter Chave de API (sk-ant-...)

1. Acesse [console.anthropic.com](https://console.anthropic.com)
2. Crie uma conta ou faça login
3. Vá para "API Keys"
4. Clique em "Create Key"
5. Copie a chave (começa com `sk-ant-`)

#### Configurar (Antes de usar, aceite os termos em <https://www.anthropic.com/policies/api-terms-of-service>)

```bash
# Via variável de ambiente
export CLAUDE_API_KEY="sk-ant-..."

# Ou no arquivo de configuração
echo 'CLAUDE_API_KEY=sk-ant-...' >> ~/.gromptrc
```

#### Modelos Disponíveis (Claude)

| Modelo | Descrição | Uso Recomendado |
|--------|-----------|-----------------|
| `claude-3-5-sonnet-20241022` | Avançado recente | Análises complexas, raciocínio |
| `claude-3-sonnet-20240229` | Balanceado | Uso geral, desenvolvimento |
| `claude-3-5-haiku-20241022` | Rápido, eficiente | Tarefas simples, iteração rápida |

#### Exemplo de Uso (Claude)

```bash
grompt ask "Explique clean architecture" \
  --provider claude \
  --model claude-3-sonnet \
  --max-tokens 1500
```

### Gemini (Google)

#### Obter Chave de API (AIza...)

1. Acesse [ai.google.dev](https://ai.google.dev)
2. Clique em "Get API key"
3. Configure um projeto no Google Cloud
4. Ative a API do Gemini
5. Copie a chave de API

#### Configurar (Antes de usar, configure faturamento no Google Cloud)

```bash
# Via variável de ambiente
export GEMINI_API_KEY="AIza..."

# Ou no arquivo de configuração
echo 'GEMINI_API_KEY=AIza...' >> ~/.gromptrc
```

#### Modelos Disponíveis (Gemini)

| Modelo | Descrição | Uso Recomendado |
|--------|-----------|-----------------|
| `gemini-1.5-flash` | Rápido | Uso geral, prototipagem |
| `gemini-2.0-flash` | Equilíbrio | Uso geral |
| `gemini-2.0-pro` | Capaz | Tarefas complexas |

#### Exemplo de Uso (Gemini)

```bash
grompt generate \
  --idea "dashboard analytics" \
  --idea "React TypeScript" \
  --provider gemini \
  --model gemini-pro
```

### DeepSeek

#### Obter Chave de API (sk-deep-...)

1. Acesse [platform.deepseek.com](https://platform.deepseek.com)
2. Registre-se ou faça login
3. Vá para "API Keys"
4. Crie uma nova chave
5. Copie a chave de API

#### Configurar (Antes de usar, aceite os termos em <https://deepseek.com/terms>)

```bash
# Via variável de ambiente
export DEEPSEEK_API_KEY="..."

# Ou no arquivo de configuração
echo 'DEEPSEEK_API_KEY=...' >> ~/.gromptrc
```

#### Modelos Disponíveis (DeepSeek)

| Modelo | Descrição | Uso Recomendado |
|--------|-----------|-----------------|
| `deepseek-chat` | Conversação geral | Discussões, explicações |
| `deepseek-coder` | Especializado em código | Desenvolvimento, debugging |

#### Exemplo de Uso (DeepSeek)

```bash
grompt ask "Otimizar query SQL complexa" \
  --provider deepseek \
  --model deepseek-coder
```

### Ollama (Local)

#### Instalar Ollama

```bash
# Linux/macOS
curl -fsSL https://ollama.ai/install.sh | sh

# Windows
# Baixe de https://ollama.ai/download/windows
```

#### Configurar (Antes de usar, aceite os termos em <https://ollama.ai/terms>)

```bash
# Iniciar Ollama
ollama serve

# Baixar modelos
ollama pull llama2
ollama pull codellama
ollama pull mistral

# Configurar endpoint
export OLLAMA_ENDPOINT="http://localhost:11434"
```

#### Modelos Populares

| Modelo | Tamanho | Descrição |
|--------|---------|-----------|
| `llama3.2` | 1B/3B/11B | Uso geral |
| `codellama` | 7B/13B/34B | Programação |
| `mistral` | 7B | Rápido e eficiente |
| `llama3.1`, `llama2` | vários | Alternativas |

#### Exemplo de Uso (Ollama)

```bash
grompt ask "Como funciona Docker?" \
  --provider ollama \
  --model llama2
```

## ⚙️ Configuração do Servidor

### Variáveis de Ambiente

```bash
# Arquivo ~/.bashrc ou ~/.zshrc
export PORT=8080                    # Porta do servidor
export HOST=localhost               # Interface de rede
export DEBUG=false                  # Logs detalhados
export CORS_ORIGINS="*"             # Origens CORS permitidas

# Provedores de IA
export OPENAI_API_KEY="sk-..."
export CLAUDE_API_KEY="sk-ant-..."
export GEMINI_API_KEY="AIza..."
export DEEPSEEK_API_KEY="..."
export OLLAMA_ENDPOINT="http://localhost:11434"

# Configurações padrão
export DEFAULT_PROVIDER="openai"
export DEFAULT_MODEL="gpt-4"
export DEFAULT_MAX_TOKENS=1000
export DEFAULT_TEMPERATURE=0.7
```

### Arquivo de Configuração

Crie um arquivo `~/.gromptrc`:

```ini
# ~/.gromptrc

# Servidor
PORT=8080
HOST=localhost
DEBUG=false

# Provedores de IA
OPENAI_API_KEY=sk-...
CLAUDE_API_KEY=sk-ant-...
GEMINI_API_KEY=AIza...
DEEPSEEK_API_KEY=...
OLLAMA_ENDPOINT=http://localhost:11434

# Padrões
DEFAULT_PROVIDER=openai
DEFAULT_MODEL=gpt-4
DEFAULT_MAX_TOKENS=1000
DEFAULT_TEMPERATURE=0.7

# Interface
THEME=dark
LANGUAGE=pt-BR
AUTO_SAVE=true
```

### Configuração Avançada

#### Proxy e Rede

```bash
# Para ambientes corporativos
export HTTP_PROXY="http://proxy.empresa.com:8080"
export HTTPS_PROXY="http://proxy.empresa.com:8080"
export NO_PROXY="localhost,127.0.0.1"

# Timeout personalizado
export REQUEST_TIMEOUT=30
export CONNECTION_TIMEOUT=10
```

> Nota: funcionalidades como rate limiting e cache não estão expostas via variáveis de ambiente na versão atual.

## 🔒 Segurança

### Proteção de Chaves de API

```bash
# Nunca commitar chaves no git
echo ".env" >> .gitignore
echo ".gromptrc" >> .gitignore

# Usar arquivo de ambiente separado
touch .env
chmod 600 .env  # Apenas owner pode ler/escrever

# Verificar permissões
ls -la ~/.gromptrc
# Deve mostrar: -rw------- (600)
```

### Produção

> Nota: Use um reverse proxy (Nginx/Caddy) para TLS/CORS. O binário não expõe TLS nativo.

### Rotação de Chaves

```bash
#!/bin/bash
# script-rotate-keys.sh

# Backup da configuração atual
cp ~/.gromptrc ~/.gromptrc.bak.$(date +%Y%m%d)

# Atualizar chave (substitua pela nova)
sed -i 's/OPENAI_API_KEY=.*/OPENAI_API_KEY=nova-chave/' ~/.gromptrc

# Testar nova configuração
grompt ask "teste" --provider openai --dry-run
```

## 🎨 Personalização
No momento, não há personalização via variáveis de ambiente.

## 📊 Monitoramento e Logs
Sem variáveis dedicadas para logs/métricas. Use o health `GET /api/health` e observabilidade do seu proxy.

## 🔧 Configuração para Desenvolvimento

### Configuração de Teste

```bash
# Ambiente de teste
export NODE_ENV=test
export TEST_API_KEYS=true
export MOCK_PROVIDERS=true

# Arquivo .env.test
cp .env .env.test
```

### Debug
Sem flags dedicadas; valide chamadas com `make test` e providers via `/api/test?provider=...`.

## 🚀 Configuração para Produção

### Variáveis de Produção

```bash
# Arquivo .env.production
NODE_ENV=production
DEBUG=false
PORT=80
HOST=0.0.0.0

# Segurança
CORS_ORIGINS=https://meudominio.com
RATE_LIMIT_ENABLED=true
HTTPS_REDIRECT=true

# Performance
CACHE_ENABLED=true
COMPRESSION_ENABLED=true
STATIC_FILES_CACHE=3600
```

### Docker

```dockerfile
# Dockerfile
FROM node:18-alpine

# Variáveis de ambiente
ENV NODE_ENV=production
ENV PORT=8080

# Copiar arquivos
COPY . /app
WORKDIR /app

# Instalar dependências e compilar
RUN npm ci --only=production
RUN npm run build

# Usuário não-root
USER node

# Comando de início
CMD ["./grompt", "--port", "8080"]
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  grompt:
    build: .
    ports:
      - "8080:8080"
    environment:
      - NODE_ENV=production
      - PORT=8080
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - CLAUDE_API_KEY=${CLAUDE_API_KEY}
    volumes:
      - ./logs:/app/logs
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

## 📋 Checklist de Configuração

### ✅ Configuração Básica

- [ ] Baixar e instalar o binário do Grompt
- [ ] Configurar pelo menos um provedor de IA
- [ ] Testar conexão com `grompt ask "teste" --dry-run`
- [ ] Configurar variáveis de ambiente
- [ ] Testar interface web em `http://localhost:8080`

### ✅ Configuração Avançada

- [ ] Configurar arquivo `~/.gromptrc`
- [ ] Configurar múltiplos provedores de IA
- [ ] Configurar Ollama para uso local
- [ ] Configurar logs e monitoramento
- [ ] Configurar segurança (permissões de arquivo)

### ✅ Configuração de Produção

- [ ] Configurar HTTPS e certificados
- [ ] Configurar CORS para domínios específicos
- [ ] Configurar rate limiting
- [ ] Configurar backup de configurações
- [ ] Configurar rotação de chaves de API

---

## 📚 Próximos Passos

- **[Exemplos de Uso](usage-examples.md)** - Veja exemplos práticos e avançados
- **[API Reference](api-reference.md)** - Integre o Grompt programaticamente
- **[Solução de Problemas](../development/troubleshooting.md)** - Resolva problemas comuns
