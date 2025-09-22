# 🚀 Grompt V1 API - Documentação Técnica

## 📋 Visão Geral

O **Grompt V1 API** é uma reforma completa do backend seguindo a arquitetura do **Analyzer**, implementando um gateway modular e escalável para geração de prompts com MultiAI providers.

### ✨ Principais Características

- 🔄 **Arquitetura baseada no Analyzer** - Mesma topologia robusta e escalável
- 🌐 **MultiAI Provider Support** - OpenAI, Anthropic, Gemini integrados via SDKs oficiais
- 🔧 **GoBE Proxy Integration** - Delegação de auth, storage e billing
- 📊 **Observabilidade completa** - Health checks, métricas e logging estruturado
- ⚡ **SSE Streaming** - Geração de prompts em tempo real
- 🛡️ **Production-ready** - Middleware, circuit breakers, rate limiting

---

## 🛠️ Arquitetura

```plaintext
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Frontend      │───▶│  Grompt Gateway  │───▶│  MultiAI Wrapper│
│   (React)       │    │   (V1 Routes)    │    │  (Official SDKs)│
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │   GoBE Proxy     │
                       │ (Auth/Storage/   │
                       │   Billing)       │
                       └──────────────────┘
```

---

## 🎯 Endpoints Implementados

### **Core Generation Routes**

#### `POST /v1/generate`

Geração síncrona de prompts com resposta completa.

**Request:**

```json
{
  "provider": "gemini",
  "model": "gemini-2.0-flash",
  "ideas": [
    "Make a REST API",
    "User authentication",
    "Database with PostgreSQL"
  ],
  "purpose": "code",
  "temperature": 0.7,
  "context": {
    "framework": "go",
    "database": "postgresql"
  }
}
```

**Response:**

```json
{
  "id": "gen_1703123456",
  "object": "prompt.generation",
  "created_at": 1703123456,
  "provider": "gemini",
  "model": "gemini-2.0-flash",
  "prompt": "Create a RESTful API server with the following specifications...",
  "ideas": ["Make a REST API", "User authentication", "Database with PostgreSQL"],
  "purpose": "code",
  "usage": {
    "tokens": 1500,
    "latency_ms": 2340,
    "cost_usd": 0.004500,
    "provider": "gemini",
    "model": "gemini-2.0-flash"
  }
}
```

#### `GET /v1/generate/stream`

Geração com SSE streaming para visualização em tempo real.

**Query Parameters:**

- `provider` - Nome do provider (obrigatório)
- `ideas[]` - Lista de ideias (obrigatório, múltiplos)
- `purpose` - Tipo de prompt: "code", "creative", "analysis"
- `model` - Modelo específico (opcional)

**SSE Events:**

```javascript
data: {"event": "generation.started", "provider": "gemini", "model": "gemini-2.0-flash"}
data: {"event": "generation.chunk", "content": "Create a RESTful API server..."}
data: {"event": "generation.chunk", "content": " with user authentication..."}
data: {"event": "generation.complete", "usage": {"tokens": 1500, "latency_ms": 2340}}
```

### **Provider Management**

#### `GET /v1/providers`

Lista providers disponíveis com status de saúde.

**Response:**

```json
{
  "object": "list",
  "data": [
    {
      "name": "gemini",
      "type": "gemini",
      "default_model": "gemini-2.0-flash",
      "available": true
    },
    {
      "name": "openai",
      "type": "openai",
      "default_model": "gpt-4o",
      "available": true
    },
    {
      "name": "anthropic",
      "type": "anthropic",
      "default_model": "claude-3-5-sonnet-20241022",
      "available": false,
      "error": "API key not configured"
    }
  ]
}
```

### **Health Monitoring**

#### `GET /v1/health`

Health check inteligente com verificação de dependências.

**Response:**

```json
{
  "status": "healthy",
  "service": "grompt-v1",
  "timestamp": 1703123456,
  "version": "1.0.0",
  "dependencies": {
    "providers": {
      "gemini": {"status": "healthy"},
      "openai": {"status": "healthy"},
      "anthropic": {"status": "unhealthy", "error": "connection timeout"}
    },
    "gobe_proxy": {
      "status": "healthy",
      "response_time": "< 5s"
    }
  }
}
```

### **GoBE Proxy Delegation**

#### `POST /v1/proxy/*`

Proxy transparente para o GoBE (auth, storage, billing).

**Headers Preservados:**

- `Authorization` - Token de autenticação
- `X-Request-Id` - ID de rastreamento
- `X-Tenant-Id` - ID do tenant
- `X-User-Id` - ID do usuário

**Exemplos:**

- `POST /v1/proxy/auth/login` → `POST {GOBE_BASE_URL}/auth/login`
- `GET /v1/proxy/billing/usage` → `GET {GOBE_BASE_URL}/billing/usage`
- `POST /v1/proxy/storage/save` → `POST {GOBE_BASE_URL}/storage/save`

---

## ⚙️ Configuração

### **Variáveis de Ambiente**

```bash
# Provider API Keys (pelo menos uma obrigatória)
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
GEMINI_API_KEY=...

# GoBE Integration
GOBE_BASE_URL=https://gobe.example.com

# Server Configuration
PORT=3000
DEBUG=true

# Default Configurations
DEFAULT_PROVIDER=gemini
DEFAULT_MODEL=gemini-2.0-flash
REQUEST_TIMEOUT_MS=30000
RATE_LIMIT_RPS=3
```

### **Configuração YAML**

```yaml
# config/config.yml
server:
  addr: ":3000"
  debug: true
  cors:
    allow_origins: ["*"]

providers:
  gemini:
    type: "gemini"
    key_env: "GEMINI_API_KEY"
    base_url: "https://generativelanguage.googleapis.com"
    default_model: "gemini-2.0-flash"

  openai:
    type: "openai"
    key_env: "OPENAI_API_KEY"
    base_url: "https://api.openai.com"
    default_model: "gpt-4o"

  anthropic:
    type: "anthropic"
    key_env: "ANTHROPIC_API_KEY"
    base_url: "https://api.anthropic.com"
    default_model: "claude-3-5-sonnet-20241022"
```

---

## 🚀 Setup e Execução

### **1. Configurar Environment**

```bash
# Clone e configure
git clone <repo>
cd grompt

# Configure as chaves de API (pelo menos uma)
export GEMINI_API_KEY="your-key-here"
export OPENAI_API_KEY="sk-your-key-here"
export ANTHROPIC_API_KEY="sk-ant-your-key-here"

# Opcional: Configure GoBE proxy
export GOBE_BASE_URL="https://your-gobe-instance.com"
```

### **2. Executar o Servidor**

```bash
# Build e execução
make build
./dist/grompt gateway serve -p 3000 -f config/config.yml

# Ou modo desenvolvimento
make run
```

### **3. Verificar Status**

```bash
# Health check
curl http://localhost:3000/v1/health

# Providers disponíveis
curl http://localhost:3000/v1/providers
```

### **4. Teste de Geração**

```bash
# Geração síncrona
curl -X POST http://localhost:3000/v1/generate \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "gemini",
    "ideas": ["Create a todo app", "With React", "And TypeScript"],
    "purpose": "code"
  }'

# Geração streaming
curl -N "http://localhost:3000/v1/generate/stream?provider=gemini&ideas[]=Todo%20app&ideas[]=React&purpose=code"
```

---

## 📊 Observabilidade

### **Métricas Coletadas**

- ✅ **Requests totais** por endpoint/método/status
- ✅ **Tokens gerados** por provider/modelo
- ✅ **Latência de geração** por provider/modelo
- ✅ **Erros de geração** por tipo
- ✅ **Custo estimado** por provider
- ✅ **Requests de proxy** para GoBE
- ✅ **Status de saúde** dos providers

### **Logging Estruturado**

```plaintext
[METRICS] Request: endpoint=/v1/generate method=POST status=200 duration=2.34s
[METRICS] Tokens generated: provider=gemini model=gemini-2.0-flash tokens=1500
[METRICS] Generation latency: provider=gemini duration=2.34s
[METRICS] Estimated cost: provider=gemini cost_usd=0.004500
```

---

## ✅ Critérios de Aceite - Status

| Critério | Status | Detalhes |
|----------|--------|----------|
| **GET /v1/health** | ✅ | Health check com verificação de dependências |
| **GET /v1/providers** | ✅ | Lista providers com status de disponibilidade |
| **POST /v1/generate** | ✅ | Geração síncrona funcionando |
| **GET /v1/generate/stream** | ✅ | SSE streaming implementado |
| **POST /v1/proxy/\*** | ✅ | Proxy para GoBE com headers preservados |
| **Observabilidade** | ✅ | Métricas e logging estruturado |
| **MultiAI Integration** | ✅ | OpenAI, Anthropic, Gemini via SDKs oficiais |

---

## 🔧 Próximos Passos

1. **Adicionar Prometheus** - Substituir métricas simples por Prometheus
2. **Rate Limiting avançado** - Por usuario/tenant via middleware
3. **Caching** - Cache de prompts gerados
4. **Webhooks** - Notificações de eventos de geração
5. **Batch Processing** - Geração de múltiplos prompts em lote

---

## 🐛 Troubleshooting

### **Provider Indisponível**

```bash
# Verificar chaves de API
echo $GEMINI_API_KEY
echo $OPENAI_API_KEY

# Verificar conectividade
curl -v https://generativelanguage.googleapis.com
```

### **GoBE Proxy Errors**

```bash
# Verificar URL do GoBE
echo $GOBE_BASE_URL

# Testar conectividade
curl $GOBE_BASE_URL/health
```

### **Debug Mode**

```bash
# Ativar logs detalhados
export DEBUG=true
./dist/grompt gateway serve -f config/config.yml
```

---

**🎉 Reforma do Backend Concluída!**

O Grompt V1 API implementa todas as funcionalidades solicitadas seguindo a arquitetura robusta do Analyzer, com integração MultiAI, proxy GoBE e observabilidade completa.
