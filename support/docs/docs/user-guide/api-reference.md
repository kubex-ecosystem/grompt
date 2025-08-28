# API Reference

O Grompt oferece uma API REST robusta para integração programática com aplicações externas. Esta página documenta todos os endpoints disponíveis.

## 🌐 Base URL

```plaintext
http://localhost:8080/api
```

## 🔐 Autenticação

A API do Grompt atualmente não requer autenticação para uso local. Para implantações em produção, considere implementar autenticação via proxy reverso ou middleware.

## 📊 Endpoints Principais

### Health Check

Verifica se o servidor está operacional.

```http
GET /api/health
```

**Response:**

```json
{
  "status": "ok",
  "version": "1.0.0",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

### Configuração

Obtém informações sobre provedores disponíveis e configuração atual.

```http
GET /api/config
```

**Response:**

```json
{
  "providers": {
    "openai": {
      "available": true,
      "models": ["gpt-4", "gpt-3.5-turbo"]
    },
    "claude": {
      "available": false,
      "reason": "API key not configured"
    },
    "demo": {
      "available": true,
      "models": ["demo"]
    }
  },
  "defaultProvider": "demo",
  "maxTokens": 2000
}
```

### Modelos Disponíveis

Lista todos os modelos disponíveis por provedor.

```http
GET /api/models
```

**Response:**

```json
{
  "openai": {
    "gpt-4": {
      "name": "GPT-4",
      "description": "Most capable model",
      "maxTokens": 8192
    },
    "gpt-3.5-turbo": {
      "name": "GPT-3.5 Turbo",
      "description": "Fast and efficient",
      "maxTokens": 4096
    }
  },
  "claude": {
    "claude-3-sonnet": {
      "name": "Claude 3 Sonnet",
      "description": "Balanced performance",
      "maxTokens": 200000
    }
  }
}
```

## 🤖 Geração de Prompts

### Endpoint Unificado

Gera prompts estruturados a partir de ideias brutas.

```http
POST /api/unified
```

**Request Body:**

```json
{
  "ideas": [
    "criar uma API REST",
    "usar Node.js e Express",
    "autenticação JWT",
    "banco PostgreSQL"
  ],
  "purpose": "code",
  "provider": "openai",
  "model": "gpt-4",
  "maxTokens": 2000,
  "temperature": 0.7
}
```

**Parameters:**

| Campo | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| `ideas` | `string[]` | ✅ | Lista de ideias brutas |
| `purpose` | `string` | ❌ | Propósito: `code`, `analysis`, `creative`, `general` |
| `provider` | `string` | ❌ | Provedor de IA (padrão: `demo`) |
| `model` | `string` | ❌ | Modelo específico |
| `maxTokens` | `number` | ❌ | Limite de tokens (padrão: 1000) |
| `temperature` | `number` | ❌ | Criatividade 0.0-1.0 (padrão: 0.7) |

**Response:**

```json
{
  "success": true,
  "data": {
    "prompt": "Crie uma API REST com as seguintes especificações:\n\n**Objetivo Principal:** Desenvolver uma API robusta usando Node.js...",
    "metadata": {
      "provider": "openai",
      "model": "gpt-4",
      "tokensUsed": 156,
      "processingTime": "1.2s"
    }
  }
}
```

### Pergunta Direta

Envia uma pergunta direta para um provedor de IA.

```http
POST /api/ask
```

**Request Body:**

```json
{
  "question": "Como implementar cache Redis em Node.js?",
  "provider": "claude",
  "model": "claude-3-sonnet",
  "maxTokens": 1500
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "answer": "Para implementar cache Redis em Node.js, você pode seguir estes passos...",
    "metadata": {
      "provider": "claude",
      "model": "claude-3-sonnet",
      "tokensUsed": 892,
      "processingTime": "2.1s"
    }
  }
}
```

### Geração de Squad

Gera uma lista de agentes de IA especializados para um projeto.

```http
POST /api/squad
```

**Request Body:**

```json
{
  "description": "Preciso desenvolver um e-commerce completo com React, Node.js e PostgreSQL",
  "provider": "openai",
  "model": "gpt-4"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "agents": [
      {
        "name": "Frontend Developer Bot",
        "role": "Desenvolver interface React",
        "skills": ["React", "TypeScript", "CSS"],
        "responsibilities": ["Componentes UI", "Estado global", "Responsividade"]
      },
      {
        "name": "Backend Developer Bot",
        "role": "Desenvolver API REST",
        "skills": ["Node.js", "Express", "PostgreSQL"],
        "responsibilities": ["Endpoints", "Autenticação", "Banco de dados"]
      }
    ],
    "summary": "Squad completo para desenvolvimento de e-commerce",
    "metadata": {
      "provider": "openai",
      "model": "gpt-4",
      "tokensUsed": 345
    }
  }
}
```

## 🔧 Endpoints Específicos por Provedor

### OpenAI

```http
POST /api/openai
```

### Claude

```http
POST /api/claude
```

### Gemini

```http
POST /api/gemini
```

### DeepSeek

```http
POST /api/deepseek
```

### Ollama

```http
POST /api/ollama
```

Todos seguem o mesmo formato do endpoint unificado, mas com configurações específicas do provedor.

## 📝 Exemplos de Integração

### JavaScript/Node.js

```javascript
const axios = require('axios');

class GromptClient {
  constructor(baseURL = 'http://localhost:8080/api') {
    this.baseURL = baseURL;
  }

  async generatePrompt(ideas, options = {}) {
    try {
      const response = await axios.post(`${this.baseURL}/unified`, {
        ideas,
        purpose: options.purpose || 'general',
        provider: options.provider || 'demo',
        model: options.model,
        maxTokens: options.maxTokens || 1000,
        temperature: options.temperature || 0.7
      });

      return response.data;
    } catch (error) {
      throw new Error(`Failed to generate prompt: ${error.message}`);
    }
  }

  async ask(question, options = {}) {
    try {
      const response = await axios.post(`${this.baseURL}/ask`, {
        question,
        provider: options.provider || 'demo',
        model: options.model,
        maxTokens: options.maxTokens || 1000
      });

      return response.data;
    } catch (error) {
      throw new Error(`Failed to ask question: ${error.message}`);
    }
  }

  async generateSquad(description, options = {}) {
    try {
      const response = await axios.post(`${this.baseURL}/squad`, {
        description,
        provider: options.provider || 'openai',
        model: options.model || 'gpt-4'
      });

      return response.data;
    } catch (error) {
      throw new Error(`Failed to generate squad: ${error.message}`);
    }
  }
}

// Uso
const client = new GromptClient();

async function example() {
  // Gerar prompt
  const prompt = await client.generatePrompt([
    'criar dashboard analytics',
    'React com TypeScript',
    'gráficos interativos'
  ], {
    purpose: 'code',
    provider: 'openai'
  });

  console.log(prompt.data.prompt);

  // Fazer pergunta
  const answer = await client.ask(
    'Como otimizar performance em React?',
    { provider: 'claude' }
  );

  console.log(answer.data.answer);

  // Gerar squad
  const squad = await client.generateSquad(
    'App mobile para delivery de comida'
  );

  console.log(squad.data.agents);
}
```

### Python

```python
import requests
import json

class GromptClient:
    def __init__(self, base_url="http://localhost:8080/api"):
        self.base_url = base_url

    def generate_prompt(self, ideas, **options):
        payload = {
            "ideas": ideas,
            "purpose": options.get("purpose", "general"),
            "provider": options.get("provider", "demo"),
            "model": options.get("model"),
            "maxTokens": options.get("max_tokens", 1000),
            "temperature": options.get("temperature", 0.7)
        }

        response = requests.post(f"{self.base_url}/unified", json=payload)
        response.raise_for_status()
        return response.json()

    def ask(self, question, **options):
        payload = {
            "question": question,
            "provider": options.get("provider", "demo"),
            "model": options.get("model"),
            "maxTokens": options.get("max_tokens", 1000)
        }

        response = requests.post(f"{self.base_url}/ask", json=payload)
        response.raise_for_status()
        return response.json()

    def generate_squad(self, description, **options):
        payload = {
            "description": description,
            "provider": options.get("provider", "openai"),
            "model": options.get("model", "gpt-4")
        }

        response = requests.post(f"{self.base_url}/squad", json=payload)
        response.raise_for_status()
        return response.json()

# Uso
client = GromptClient()

# Gerar prompt
prompt = client.generate_prompt([
    "análise de dados de vendas",
    "Python pandas",
    "visualizações matplotlib"
], purpose="analysis", provider="claude")

print(prompt["data"]["prompt"])

# Fazer pergunta
answer = client.ask(
    "Como fazer deploy de modelo ML?",
    provider="openai"
)

print(answer["data"]["answer"])

# Gerar squad
squad = client.generate_squad(
    "Sistema de monitoramento de infraestrutura"
)

for agent in squad["data"]["agents"]:
    print(f"{agent['name']}: {agent['role']}")
```

### cURL

```bash
# Gerar prompt
curl -X POST http://localhost:8080/api/unified \
  -H "Content-Type: application/json" \
  -d '{
    "ideas": [
      "sistema de autenticação",
      "JWT tokens",
      "Node.js Express"
    ],
    "purpose": "code",
    "provider": "openai",
    "model": "gpt-4"
  }'

# Fazer pergunta
curl -X POST http://localhost:8080/api/ask \
  -H "Content-Type: application/json" \
  -d '{
    "question": "Como configurar CORS em Express?",
    "provider": "claude"
  }'

# Gerar squad
curl -X POST http://localhost:8080/api/squad \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Plataforma de e-learning",
    "provider": "openai"
  }'
```

## ⚠️ Tratamento de Erros

### Códigos de Status HTTP

| Código | Descrição |
|--------|-----------|
| `200` | Sucesso |
| `400` | Request inválido |
| `401` | Não autorizado |
| `404` | Endpoint não encontrado |
| `429` | Rate limit excedido |
| `500` | Erro interno do servidor |
| `503` | Serviço indisponível |

### Formato de Erro

```json
{
  "success": false,
  "error": {
    "code": "INVALID_PROVIDER",
    "message": "Provider 'invalid' not supported",
    "details": {
      "supportedProviders": ["openai", "claude", "gemini", "demo"]
    }
  }
}
```

### Tipos de Erro Comuns

| Código | Descrição |
|--------|-----------|
| `INVALID_PROVIDER` | Provedor não suportado |
| `MISSING_API_KEY` | Chave de API não configurada |
| `RATE_LIMIT_EXCEEDED` | Limite de requisições excedido |
| `INVALID_MODEL` | Modelo não disponível |
| `PROMPT_TOO_LONG` | Prompt excede limite de tokens |
| `PROVIDER_ERROR` | Erro do provedor de IA |

## 📊 Rate Limiting

- **Limite padrão:** 100 requests por minuto por IP
- **Headers de resposta:**
  - `X-RateLimit-Limit`: Limite total
  - `X-RateLimit-Remaining`: Requests restantes
  - `X-RateLimit-Reset`: Timestamp do reset

## 🔒 Considerações de Segurança

- **CORS:** Configure adequadamente em produção
- **Rate limiting:** Implemente limites apropriados
- **Validação:** Sempre valide inputs
- **Logs:** Monitore uso e erros
- **Proxy:** Use proxy reverso em produção

---

## 📚 Próximos Passos

- **[Configuração](configuration.md)** - Configure provedores e personalização
- **[Exemplos de Uso](usage-examples.md)** - Veja casos práticos de integração
- **[Contribuindo](../development/contributing.md)** - Ajude a melhorar a API
