# Comandos CLI

O Grompt oferece uma interface de linha de comando poderosa para interagir com modelos de IA e gerenciar prompts. Este guia detalha todos os comandos disponíveis.

## 📋 Sintaxe Geral

```bash
grompt [comando] [argumentos] [opções]
```

## 🚀 Comandos Principais

### `grompt` - Iniciar Servidor Web

Inicia o servidor web do Grompt com interface React.

```bash
# Iniciar na porta padrão (8080)
grompt

# Especificar porta
grompt --port 3000

# Modo debug
grompt --debug

# Combinando opções
grompt --port 8081 --debug
```

**Opções:**

| Opção | Descrição | Padrão |
|-------|-----------|--------|
| `--port, -p` | Porta do servidor | `8080` |
| `--debug, -d` | Habilitar logs detalhados | `false` |
| `--host` | Interface de rede | `localhost` |

### `grompt ask` - Pergunta Direta

Envia uma pergunta direta para um provedor de IA.

```bash
# Pergunta simples
grompt ask "Qual é a capital do Brasil?"

# Com provedor específico
grompt ask "Como implementar JWT?" --provider openai --model gpt-4

# Com configurações avançadas
grompt ask "Explique microservices" \
  --provider claude \
  --model claude-3-sonnet \
  --max-tokens 1000 \
  --temperature 0.7
```

**Opções:**

| Opção | Descrição | Exemplo |
|-------|-----------|---------|
| `--provider, -p` | Provedor de IA | `openai`, `claude`, `gemini`, `demo` |
| `--model, -m` | Modelo específico | `gpt-4`, `claude-3-sonnet` |
| `--max-tokens` | Limite de tokens na resposta | `1000` |
| `--temperature` | Criatividade (0.0-1.0) | `0.7` |
| `--apikey` | Chave de API (sobrescreve env var) | `sk-...` |
| `--dry-run` | Simular sem executar | - |

### `grompt generate` - Geração de Prompts

Cria prompts estruturados a partir de ideias brutas usando engenharia de prompts.

```bash
# Geração básica
grompt generate \
  --idea "criar uma API REST" \
  --idea "com autenticação JWT" \
  --idea "usando Node.js"

# Com propósito específico
grompt generate \
  --idea "analisar dados de vendas" \
  --idea "criar gráficos" \
  --idea "usar Python pandas" \
  --purpose analysis \
  --provider openai

# Com configurações avançadas
grompt generate \
  --idea "sistema de chat" \
  --idea "tempo real com WebSockets" \
  --idea "React frontend" \
  --purpose code \
  --provider claude \
  --model claude-3-sonnet \
  --max-tokens 2000
```

**Opções:**

| Opção | Descrição | Valores |
|-------|-----------|---------|
| `--idea, -i` | Ideias brutas (múltiplas) | Texto livre |
| `--purpose` | Propósito do prompt | `code`, `analysis`, `creative`, `general` |
| `--provider, -p` | Provedor de IA | `openai`, `claude`, `gemini`, `demo` |
| `--model, -m` | Modelo específico | Varia por provedor |
| `--max-tokens` | Limite de tokens | Número |
| `--temperature` | Criatividade | `0.0` - `1.0` |

### `grompt squad` - Geração de Squad

Gera uma lista de agentes de IA especializados baseada nos requisitos do projeto.

```bash
# Squad básico
grompt squad "Preciso de um e-commerce com React, Node.js e PostgreSQL"

# Com provedor específico
grompt squad "App mobile para delivery com pagamento integrado" \
  --provider openai \
  --model gpt-4

# Salvar em arquivo específico
grompt squad "Sistema de gestão hospitalar" \
  --output TEAM.md \
  --provider claude
```

**Opções:**

| Opção | Descrição | Padrão |
|-------|-----------|--------|
| `--output, -o` | Arquivo de saída | `AGENTS.md` |
| `--provider, -p` | Provedor de IA | `openai` |
| `--model, -m` | Modelo específico | Padrão do provedor |

## 🔧 Opções Globais

### Configuração de Provedores

```bash
# Variáveis de ambiente (recomendado)
export OPENAI_API_KEY="sk-..."
export CLAUDE_API_KEY="sk-ant-..."
export GEMINI_API_KEY="..."
export DEEPSEEK_API_KEY="..."

# Via parâmetro (menos seguro)
grompt ask "pergunta" --provider openai --apikey "sk-..."
```

### Debug e Logs

```bash
# Logs detalhados
DEBUG=true grompt ask "teste"

# Ou via parâmetro
grompt --debug ask "teste"

# Dry run (simular sem executar)
grompt ask "teste" --dry-run
```

## 📝 Exemplos Práticos

### Desenvolvimento Web

```bash
# Frontend React
grompt generate \
  --idea "componente de formulário" \
  --idea "validação em tempo real" \
  --idea "TypeScript" \
  --idea "styled-components" \
  --purpose code

# Backend API
grompt generate \
  --idea "API REST para blog" \
  --idea "CRUD de posts" \
  --idea "autenticação JWT" \
  --idea "Node.js Express" \
  --purpose code

# Banco de dados
grompt generate \
  --idea "esquema PostgreSQL" \
  --idea "relacionamento usuários-posts" \
  --idea "índices otimizados" \
  --purpose code
```

### Análise de Dados

```bash
# Análise exploratória
grompt generate \
  --idea "dataset de vendas" \
  --idea "identificar padrões" \
  --idea "estatísticas descritivas" \
  --idea "Python pandas" \
  --purpose analysis

# Machine Learning
grompt generate \
  --idea "modelo preditivo" \
  --idea "classificação de clientes" \
  --idea "scikit-learn" \
  --idea "validação cruzada" \
  --purpose analysis

# Visualização
grompt generate \
  --idea "dashboard interativo" \
  --idea "gráficos de vendas" \
  --idea "filtros dinâmicos" \
  --idea "Plotly Dash" \
  --purpose analysis
```

### DevOps e Infraestrutura

```bash
# Docker
grompt generate \
  --idea "containerizar aplicação Node.js" \
  --idea "multi-stage build" \
  --idea "otimizar tamanho da imagem" \
  --purpose code

# CI/CD
grompt generate \
  --idea "pipeline GitHub Actions" \
  --idea "testes automatizados" \
  --idea "deploy automático" \
  --idea "AWS EC2" \
  --purpose code

# Monitoramento
grompt generate \
  --idea "alertas de performance" \
  --idea "métricas customizadas" \
  --idea "Prometheus Grafana" \
  --purpose code
```

### Criação de Conteúdo

```bash
# Blog técnico
grompt generate \
  --idea "artigo sobre GraphQL" \
  --idea "comparação com REST" \
  --idea "exemplos práticos" \
  --idea "para desenvolvedores iniciantes" \
  --purpose creative

# Documentação
grompt generate \
  --idea "documentar API REST" \
  --idea "endpoints CRUD" \
  --idea "exemplos de requisições" \
  --idea "códigos de erro" \
  --purpose creative

# Marketing técnico
grompt generate \
  --idea "copy para landing page" \
  --idea "SaaS para desenvolvedores" \
  --idea "foco em conversão" \
  --idea "benefícios técnicos" \
  --purpose creative
```

## 🔍 Comandos de Utilidade

### `grompt version`

```bash
grompt version
# ou
grompt --version
```

### `grompt help`

```bash
# Ajuda geral
grompt --help

# Ajuda de comando específico
grompt ask --help
grompt generate --help
grompt squad --help
```

### Verificação de Configuração

```bash
# Testar conectividade com provedores
grompt ask "teste" --provider openai --dry-run
grompt ask "teste" --provider claude --dry-run

# Listar provedores disponíveis
grompt ask --help | grep -A5 "provider"
```

## ⚙️ Configurações Avançadas

### Arquivo de Configuração

Crie um arquivo `.gromptrc` no diretório home:

```bash
# ~/.gromptrc
OPENAI_API_KEY=sk-...
CLAUDE_API_KEY=sk-ant-...
DEFAULT_PROVIDER=openai
DEFAULT_MODEL=gpt-4
DEFAULT_MAX_TOKENS=1000
PORT=8080
DEBUG=false
```

### Scripts e Automação

```bash
#!/bin/bash
# script-grompt.sh

# Função para gerar prompts de código
generate_code_prompt() {
    local description="$1"
    grompt generate \
        --idea "$description" \
        --purpose code \
        --provider openai \
        --model gpt-4 \
        --max-tokens 2000
}

# Uso
generate_code_prompt "API REST para e-commerce com Node.js"
```

### Integração com Editores

```bash
# Vim - adicionar ao .vimrc
command! -nargs=1 GromptAsk !grompt ask "<args>"

# VS Code - via task.json
{
  "label": "Grompt Generate",
  "type": "shell",
  "command": "grompt",
  "args": ["generate", "--idea", "${input:idea}", "--purpose", "code"]
}
```

## 🛠️ Solução de Problemas

### Erros Comuns

```bash
# Erro: comando não encontrado
which grompt
export PATH=$PATH:/usr/local/bin

# Erro: porta em uso
lsof -i :8080
grompt --port 8081

# Erro: API key inválida
grompt ask "teste" --provider openai --dry-run

# Erro: timeout de conexão
curl -I https://api.openai.com
```

### Debug Avançado

```bash
# Logs detalhados
DEBUG=true grompt ask "teste" 2>&1 | tee debug.log

# Verificar variáveis de ambiente
env | grep -E "(OPENAI|CLAUDE|GEMINI|DEEPSEEK)"

# Testar conectividade
grompt ask "hello" --provider demo --debug
```

---

## 📚 Próximos Passos

- **[Configuração](configuration.md)** - Configure provedores de IA e personalize o ambiente
- **[Exemplos de Uso](usage-examples.md)** - Veja exemplos práticos e avançados
- **[API Reference](api-reference.md)** - Integre o Grompt programaticamente
