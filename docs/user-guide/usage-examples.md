# Exemplos de Uso

Esta página apresenta exemplos práticos e avançados de como usar o Grompt em diferentes cenários do mundo real.

## 🖥️ Desenvolvimento Web

### Frontend React com TypeScript

**Cenário:** Criar um componente de formulário complexo com validação

```bash
grompt generate \
  --idea "componente React para formulário de cadastro" \
  --idea "validação em tempo real com Yup" \
  --idea "TypeScript com tipos strict" \
  --idea "styled-components para estilização" \
  --idea "tratamento de erros e loading states" \
  --purpose code \
  --provider openai \
  --model gpt-4
```

**Resultado gerado:**

```plaintext
Crie um componente React TypeScript para formulário de cadastro com as seguintes especificações:

**Objetivo Principal:** Desenvolver um componente reutilizável de formulário com validação robusta

**Requisitos Técnicos:**
- TypeScript com tipagem estrita
- Validação em tempo real usando Yup schema
- Estados de loading, erro e sucesso
- Styled-components para estilização consistente
- Acessibilidade (ARIA labels, navegação por teclado)

**Estrutura do Formulário:**
- Campos: nome, email, senha, confirmação de senha
- Validações: email válido, senha forte, confirmação matching
- Feedback visual imediato para cada campo

**Estados de Interface:**
- Idle: formulário pronto para preenchimento
- Validating: validação em tempo real ativa
- Submitting: enviando dados (spinner/loading)
- Success: cadastro realizado com sucesso
- Error: exibir erros de validação ou servidor

**Exemplo de Uso:**
```typescript
<CadastroForm
  onSubmit={handleSubmit}
  onSuccess={handleSuccess}
  onError={handleError}
  validationSchema={cadastroSchema}
/>
/```

**Deliverables:**

- Componente CadastroForm.tsx com tipos TypeScript
- Schema de validação Yup
- Styled-components para estilização
- Hooks customizados para lógica de formulário
- Testes unitários com React Testing Library

```

### Backend API Node.js

**Cenário:** API REST para sistema de e-commerce

```bash
grompt generate \
  --idea "API REST para e-commerce" \
  --idea "autenticação JWT com refresh tokens" \
  --idea "CRUD produtos, usuários, pedidos" \
  --idea "Node.js Express TypeScript" \
  --idea "PostgreSQL com Prisma ORM" \
  --idea "middleware de validação e rate limiting" \
  --purpose code \
  --provider claude \
  --model claude-3-sonnet
```

### Sistema de Autenticação

**Cenário:** Implementar autenticação completa

```bash
grompt generate \
  --idea "sistema de autenticação completo" \
  --idea "login, registro, recuperação de senha" \
  --idea "JWT com refresh tokens" \
  --idea "middleware de autorização" \
  --idea "roles e permissões granulares" \
  --idea "rate limiting e proteção contra ataques" \
  --purpose code
```

## 📊 Análise de Dados

### Análise Exploratória com Python

**Cenário:** Analisar dados de vendas de uma empresa

```bash
grompt generate \
  --idea "análise exploratória de dados de vendas" \
  --idea "dataset com 100k registros de transações" \
  --idea "identificar padrões sazonais e tendências" \
  --idea "Python pandas, matplotlib, seaborn" \
  --idea "estatísticas descritivas e visualizações" \
  --idea "relatório executivo com insights" \
  --purpose analysis \
  --provider openai
```

### Machine Learning

**Cenário:** Modelo de predição de churn de clientes

```bash
grompt generate \
  --idea "modelo de predição de churn de clientes" \
  --idea "features demográficas e comportamentais" \
  --idea "scikit-learn com random forest e XGBoost" \
  --idea "feature engineering e seleção" \
  --idea "validação cruzada e métricas de performance" \
  --idea "interpretabilidade com SHAP" \
  --purpose analysis \
  --provider claude
```

### Dashboard Interativo

**Cenário:** Dashboard de KPIs de negócio

```bash
grompt generate \
  --idea "dashboard interativo de KPIs de vendas" \
  --idea "métricas em tempo real" \
  --idea "filtros por período, região, produto" \
  --idea "Plotly Dash ou Streamlit" \
  --idea "conexão com banco PostgreSQL" \
  --idea "cache Redis para performance" \
  --purpose analysis
```

## 🔧 DevOps e Infraestrutura

### Containerização com Docker

**Cenário:** Containerizar aplicação full-stack

```bash
grompt generate \
  --idea "containerizar aplicação full-stack" \
  --idea "frontend React, backend Node.js, PostgreSQL" \
  --idea "multi-stage builds otimizados" \
  --idea "docker-compose para desenvolvimento" \
  --idea "networks e volumes configurados" \
  --idea "health checks e restart policies" \
  --purpose code
```

### CI/CD Pipeline

**Cenário:** Pipeline completo no GitHub Actions

```bash
grompt generate \
  --idea "pipeline CI/CD completo GitHub Actions" \
  --idea "testes automatizados em múltiplas versões" \
  --idea "build e push para Docker Hub" \
  --idea "deploy automático para AWS EC2" \
  --idea "notificações Slack em falhas" \
  --idea "rollback automático em caso de erro" \
  --purpose code
```

### Monitoramento e Observabilidade

**Cenário:** Stack de monitoramento completa

```bash
grompt generate \
  --idea "stack de monitoramento para microservices" \
  --idea "Prometheus para métricas" \
  --idea "Grafana dashboards customizados" \
  --idea "alertas baseados em SLIs/SLOs" \
  --idea "logs centralizados com ELK stack" \
  --idea "tracing distribuído com Jaeger" \
  --purpose code
```

## 📝 Criação de Conteúdo

### Blog Técnico

**Cenário:** Artigo técnico detalhado

```bash
grompt generate \
  --idea "artigo sobre microservices patterns" \
  --idea "comparação com arquitetura monolítica" \
  --idea "exemplos práticos com código" \
  --idea "pros, contras e quando usar" \
  --idea "para desenvolvedores sênior" \
  --idea "8000 palavras com diagramas" \
  --purpose creative
```

### Documentação de API

**Cenário:** Documentação técnica completa

```bash
grompt generate \
  --idea "documentação completa para API REST" \
  --idea "todos os endpoints com exemplos" \
  --idea "códigos de erro detalhados" \
  --idea "guia de autenticação" \
  --idea "rate limiting e best practices" \
  --idea "formato OpenAPI 3.0" \
  --purpose creative
```

### Curso Online

**Cenário:** Estrutura de curso técnico

```bash
grompt generate \
  --idea "curso online sobre React avançado" \
  --idea "12 módulos progressivos" \
  --idea "projetos práticos em cada módulo" \
  --idea "hooks customizados, context, performance" \
  --idea "testes, TypeScript, patterns avançados" \
  --idea "público alvo: devs com experiência básica" \
  --purpose creative
```

## 🎯 Casos de Uso Específicos

### Code Review Automatizado

**Cenário:** Prompts para revisão de código

```bash
# Revisão de segurança
grompt generate \
  --idea "revisar código para vulnerabilidades" \
  --idea "OWASP Top 10" \
  --idea "injection attacks, XSS, CSRF" \
  --idea "autenticação e autorização" \
  --purpose code

# Revisão de performance
grompt generate \
  --idea "revisar código para otimização" \
  --idea "algoritmos ineficientes" \
  --idea "queries N+1, memory leaks" \
  --idea "cache strategies" \
  --purpose code
```

### Debugging Assistido

**Cenário:** Ajuda para resolver bugs complexos

```bash
grompt ask \
  "Tenho um memory leak em aplicação Node.js. CPU usage cresce gradualmente. Como investigar?" \
  --provider claude \
  --model claude-3-sonnet

grompt ask \
  "Query SQL demora 30 segundos, tem 3 JOINs e WHERE complexo. Como otimizar?" \
  --provider openai \
  --model gpt-4
```

### Geração de Testes

**Cenário:** Testes automatizados completos

```bash
grompt generate \
  --idea "suite de testes para API REST" \
  --idea "testes unitários, integração, e2e" \
  --idea "Jest, Supertest, Cypress" \
  --idea "mocks e fixtures organizados" \
  --idea "cobertura 90%+ critical paths" \
  --purpose code
```

## 🚀 Workflows Avançados

### Squad de Agentes para Projeto Completo

**Cenário:** E-commerce completo do zero

```bash
grompt squad \
  "Preciso desenvolver um e-commerce completo do zero:
  - Frontend React/Next.js responsivo
  - Backend Node.js/Express com microservices
  - Pagamentos via Stripe
  - Gestão de estoque em tempo real
  - Sistema de reviews e avaliações
  - Painel administrativo completo
  - App mobile React Native
  - Deploy em AWS com CI/CD
  - Monitoramento e alertas"
```

### Pipeline de Conteúdo

**Cenário:** Automatizar criação de conteúdo

```bash
# 1. Pesquisa de tópicos
grompt ask "Quais são os 10 tópicos mais relevantes em DevOps para 2024?"

# 2. Estrutura de artigo
grompt generate \
  --idea "artigo sobre GitOps" \
  --idea "comparação com CI/CD tradicional" \
  --idea "exemplos práticos com ArgoCD" \
  --purpose creative

# 3. SEO e metadata
grompt generate \
  --idea "otimizar artigo para SEO" \
  --idea "palavras-chave: GitOps, DevOps, CI/CD" \
  --idea "meta description, title tags" \
  --purpose creative
```

### Arquitetura de Sistema

**Cenário:** Design de arquitetura complexa

```bash
grompt generate \
  --idea "arquitetura de sistema para streaming de vídeo" \
  --idea "handle 1M usuários concorrentes" \
  --idea "CDN, encoding adaptivo, live streaming" \
  --idea "AWS ou GCP" \
  --idea "microservices com Event Sourcing" \
  --idea "cache distribuído Redis Cluster" \
  --purpose code
```

## 📈 Otimização de Workflows

### Templates Personalizados

Crie aliases para comandos frequentes:

```bash
# No seu ~/.bashrc ou ~/.zshrc
alias grompt-api='grompt generate --purpose code --provider openai --model gpt-4'
alias grompt-analysis='grompt generate --purpose analysis --provider claude'
alias grompt-creative='grompt generate --purpose creative --provider openai'

# Uso
grompt-api \
  --idea "API REST para blog" \
  --idea "Node.js TypeScript" \
  --idea "autenticação JWT"
```

### Scripts de Automação

```bash
#!/bin/bash
# generate-project.sh

PROJECT_TYPE="$1"
DESCRIPTION="$2"

case $PROJECT_TYPE in
  "api")
    grompt generate \
      --idea "API REST $DESCRIPTION" \
      --idea "Node.js TypeScript Express" \
      --idea "PostgreSQL Prisma ORM" \
      --idea "autenticação JWT" \
      --idea "testes Jest" \
      --purpose code
    ;;
  "frontend")
    grompt generate \
      --idea "Frontend React $DESCRIPTION" \
      --idea "TypeScript hooks customizados" \
      --idea "styled-components" \
      --idea "React Query estado global" \
      --purpose code
    ;;
  "fullstack")
    grompt squad "Aplicação full-stack: $DESCRIPTION"
    ;;
esac
```

### Integração com IDEs

**VS Code Task:**

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Grompt: Generate Code Prompt",
      "type": "shell",
      "command": "grompt",
      "args": [
        "generate",
        "--idea", "${input:idea1}",
        "--idea", "${input:idea2}",
        "--purpose", "code"
      ],
      "group": "build",
      "presentation": {
        "echo": true,
        "reveal": "always"
      }
    }
  ],
  "inputs": [
    {
      "id": "idea1",
      "description": "Primeira ideia",
      "type": "promptString"
    },
    {
      "id": "idea2",
      "description": "Segunda ideia",
      "type": "promptString"
    }
  ]
}
```

## 🔍 Dicas e Truques

### Prompts Mais Eficazes

1. **Seja específico sobre o contexto:**

   ```bash
   --idea "aplicação React para startup fintech B2B"
   --idea "conformidade PCI-DSS obrigatória"
   ```

2. **Inclua restrições técnicas:**

   ```bash
   --idea "deploy em Kubernetes"
   --idea "máximo 512MB RAM por container"
   ```

3. **Defina critérios de sucesso:**

   ```bash
   --idea "API deve responder em <100ms"
   --idea "suportar 10k requests/segundo"
   ```

### Iteração de Prompts

```bash
# 1. Prompt inicial amplo
grompt generate --idea "sistema de chat" --purpose code

# 2. Refinar com mais detalhes
grompt generate \
  --idea "sistema de chat em tempo real" \
  --idea "WebSockets Socket.io" \
  --idea "salas privadas e grupos" \
  --purpose code

# 3. Focar em aspecto específico
grompt generate \
  --idea "autenticação para chat WebSocket" \
  --idea "JWT tokens via query string" \
  --idea "middleware Socket.io" \
  --purpose code
```

### Combinação de Provedores

```bash
# Claude para análise e estruturação
grompt generate \
  --idea "estrutura de microservices" \
  --provider claude

# OpenAI para implementação
grompt generate \
  --idea "implementar service discovery" \
  --idea "Consul ou etcd" \
  --provider openai
```

---

## 📚 Próximos Passos

- **[API Reference](api-reference.md)** - Integre o Grompt programaticamente
- **[Arquitetura](../development/architecture.md)** - Entenda como o Grompt funciona internamente
- **[Contribuindo](../development/contributing.md)** - Ajude a melhorar o projeto
