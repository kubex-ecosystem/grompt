# 🚀 Grompt V1 Quick Start Guide

Guia rápido para começar a usar o **Grompt V1 API** em menos de 5 minutos!

## 📋 Pré-requisitos

- Go 1.25+ instalado
- Pelo menos uma API key de provider AI
- (Opcional) GoBE instance para proxy

## ⚡ Setup Rápido

### 1. **Clone e Build**
```bash
git clone <repo>
cd grompt
make build
```

### 2. **Configure API Keys** (pelo menos uma)
```bash
# Gemini (Recomendado - mais rápido e barato)
export GEMINI_API_KEY="your-gemini-key-here"

# OpenAI (Alternativa robusta)
export OPENAI_API_KEY="sk-your-openai-key"

# Anthropic Claude (Para análises complexas)
export ANTHROPIC_API_KEY="sk-ant-your-anthropic-key"

# Opcional: GoBE Proxy
export GOBE_BASE_URL="https://your-gobe-instance.com"
```

### 3. **Executar o Servidor**
```bash
./dist/grompt gateway serve -p 3000 -f config/config.yml

# Ou para desenvolvimento
make run
```

### 4. **Verificar Status**
```bash
curl http://localhost:3000/v1/health
```

## 🎯 Exemplos Práticos

### **Geração de Código - Síncrona**
```bash
curl -X POST http://localhost:3000/v1/generate \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "gemini",
    "ideas": [
      "Create a REST API for a todo app",
      "Use Go with Gin framework",
      "Include CRUD operations",
      "Add user authentication"
    ],
    "purpose": "code",
    "temperature": 0.7
  }'
```

**Resposta:**
```json
{
  "id": "gen_1703123456",
  "object": "prompt.generation",
  "prompt": "Create a RESTful API for a todo application using Go and the Gin framework...",
  "provider": "gemini",
  "model": "gemini-2.0-flash",
  "usage": {
    "tokens": 2340,
    "latency_ms": 1850,
    "cost_usd": 0.007020
  }
}
```

### **Escrita Criativa - Streaming**
```bash
curl -N "http://localhost:3000/v1/generate/stream?provider=gemini&ideas[]=Sci-fi%20story&ideas[]=Time%20travel&ideas[]=Mystery%20elements&purpose=creative"
```

**SSE Stream:**
```
data: {"event": "generation.started", "provider": "gemini"}
data: {"event": "generation.chunk", "content": "Write a captivating science fiction story..."}
data: {"event": "generation.chunk", "content": " that weaves together time travel paradoxes..."}
data: {"event": "generation.complete", "usage": {"tokens": 1890}}
```

### **Análise de Dados**
```bash
curl -X POST http://localhost:3000/v1/generate \
  -H "Content-Type: application/json" \
  -d '{
    "provider": "anthropic",
    "ideas": [
      "Analyze customer churn data",
      "Identify key patterns and drivers",
      "Create visualization recommendations",
      "Suggest retention strategies"
    ],
    "purpose": "analysis",
    "context": {
      "data_source": "SQL database",
      "timeframe": "last 12 months",
      "customer_segments": ["enterprise", "SMB", "startup"]
    }
  }'
```

## 🔧 Configuração Avançada

### **config/config.yml**
```yaml
server:
  addr: ":3000"
  debug: false

providers:
  gemini:
    type: "gemini"
    key_env: "GEMINI_API_KEY"
    default_model: "gemini-2.0-flash"

  openai:
    type: "openai"
    key_env: "OPENAI_API_KEY"
    default_model: "gpt-4o"
```

### **Providers Disponíveis**
```bash
curl http://localhost:3000/v1/providers
```

```json
{
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
    }
  ]
}
```

## 🛡️ GoBE Proxy (Auth/Storage/Billing)

```bash
# Configure GoBE
export GOBE_BASE_URL="https://gobe.example.com"

# Proxy requests mantendo headers
curl -X POST http://localhost:3000/v1/proxy/auth/login \
  -H "Authorization: Bearer your-token" \
  -H "X-Request-Id: req-123" \
  -d '{"username": "user", "password": "pass"}'

# Automaticamente proxy para: https://gobe.example.com/auth/login
```

## 📊 Monitoramento

### **Health Check Detalhado**
```bash
curl http://localhost:3000/v1/health | jq
```

```json
{
  "status": "healthy",
  "service": "grompt-v1",
  "dependencies": {
    "providers": {
      "gemini": {"status": "healthy"},
      "openai": {"status": "healthy"}
    },
    "gobe_proxy": {
      "status": "healthy",
      "response_time": "< 5s"
    }
  }
}
```

### **Logs de Métricas**
```
[METRICS] endpoint=/v1/generate provider=gemini model=gemini-2.0-flash duration=1.85s tokens=2340 cost=0.007020 status=200
```

## 🧪 Testes Automatizados

```bash
# Script de teste completo
./scripts/test_grompt_v1.sh

# Teste com servidor remoto
./scripts/test_grompt_v1.sh -u https://grompt.example.com

# Com timeout customizado
TIMEOUT=30 ./scripts/test_grompt_v1.sh
```

## 🚀 Deploy em Produção

### **Docker**
```bash
# Build
make build-docker

# Run
docker run -p 3000:3000 \
  -e GEMINI_API_KEY="your-key" \
  -e GOBE_BASE_URL="https://gobe.prod.com" \
  grompt:latest
```

### **Kubernetes**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: grompt-gateway
spec:
  replicas: 3
  selector:
    matchLabels:
      app: grompt-gateway
  template:
    metadata:
      labels:
        app: grompt-gateway
    spec:
      containers:
      - name: grompt
        image: grompt:latest
        ports:
        - containerPort: 3000
        env:
        - name: GEMINI_API_KEY
          valueFrom:
            secretKeyRef:
              name: grompt-secrets
              key: gemini-api-key
```

## 🔍 Troubleshooting

### **Provider Indisponível**
```bash
# 1. Verificar chave de API
echo $GEMINI_API_KEY

# 2. Testar conectividade
curl -v https://generativelanguage.googleapis.com

# 3. Verificar logs
tail -f logs/grompt.log
```

### **GoBE Proxy Errors**
```bash
# 1. Verificar URL
echo $GOBE_BASE_URL

# 2. Testar connectivity
curl $GOBE_BASE_URL/health

# 3. Debug mode
DEBUG=true ./dist/grompt gateway serve
```

### **Performance Issues**
```bash
# Configurar timeouts
export REQUEST_TIMEOUT_MS=60000

# Aumentar rate limits
export RATE_LIMIT_RPS=20

# Monitor de recursos
htop
```

## 💡 Dicas de Uso

1. **Gemini** = Melhor custo-benefício para maioria dos casos
2. **OpenAI** = Mais estável para código complexo
3. **Anthropic** = Melhor para análises profundas
4. **Temperature 0.7** = Equilíbrio ideal criatividade/precisão
5. **Purpose específico** = Prompts mais focados e efetivos

## 🎉 Próximos Passos

- ✅ **Está funcionando?** → Integre com seu frontend
- 🔧 **Precisa customizar?** → Edite `config/config.yml`
- 📊 **Quer métricas?** → Configure Prometheus
- 🚀 **Deploy produção?** → Use Docker/K8s templates

---

**🎯 Em menos de 5 minutos você tem um gateway AI production-ready!**

Para dúvidas: [Documentação Completa](GROMPT_V1_API.md) | [Configuração Avançada](config/grompt_v1.example.yml)