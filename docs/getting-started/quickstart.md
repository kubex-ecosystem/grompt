# Início Rápido

Comece a usar o Grompt em apenas 5 minutos! Este guia mostra como transformar suas primeiras ideias em prompts profissionais.

## 🚀 Configuração passo a passo

### 1. Baixar e Executar

#### 📦 Downloads Rápidos

- **Linux (amd64)** — [grompt_linux_amd64](https://github.com/kubex-ecosystem/gemx/grompt/releases/latest/download/grompt_linux_amd64)
- **macOS Intel** — [grompt_darwin_amd64](https://github.com/kubex-ecosystem/gemx/grompt/releases/latest/download/grompt_darwin_amd64)
- **macOS Apple Silicon** — [grompt_darwin_arm64](https://github.com/kubex-ecosystem/gemx/grompt/releases/latest/download/grompt_darwin_arm64)
- **Windows (amd64)** — [grompt_windows_amd64.exe](https://github.com/kubex-ecosystem/gemx/grompt/releases/latest/download/grompt_windows_amd64.exe)

👉 Ou veja todas as opções na [página de Releases](https://github.com/kubex-ecosystem/gemx/grompt/releases).

#### Verificação opcional de integridade (Linux/macOS)

```bash
sha256sum grompt_linux_amd64
# compare com o checksum publicado no release
```

#### 📦 Download e Execução pelo Terminal

```bash
# Linux/macOS
curl -L https://github.com/kubex-ecosystem/gemx/grompt/releases/latest/download/grompt_linux_amd64 -o grompt
chmod +x grompt
./grompt
```

```powershell
# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/kubex-ecosystem/gemx/grompt/releases/latest/download/grompt_windows_amd64.exe" -OutFile "grompt.exe"
.\grompt.exe
```

#### Abrir no Navegador

O Grompt abre automaticamente em: **<http://localhost:8080>**

#### Começar a Usar

Nenhuma configuração adicional é necessária. O modo demo está pronto para uso.

---

## 🎯 Primeiro Prompt em 2 Minutos

### Via Interface Web

1. **Abra** <http://localhost:8080>
2. **Digite suas ideias** na caixa de texto:

   ```plaintext
   quero uma função javascript
   que calcule desconto
   e aplique em produtos
   ```

3. **Clique em "Gerar Prompt"**
4. **Copie o resultado estruturado!**

### Via CLI

```bash
# Gerar prompt a partir de ideias
grompt generate \
  --idea "quero uma função javascript" \
  --idea "que calcule desconto" \
  --idea "e aplique em produtos" \
  --provider demo
```

**Resultado:**

```markdown
Crie uma função JavaScript para cálculo e aplicação de descontos com as seguintes especificações:

**Objetivo Principal:** Implementar lógica de desconto flexível para produtos

**Requisitos Funcionais:**
- Aceitar parâmetros: produto (objeto), percentual de desconto (número)
- Calcular o valor do desconto baseado no preço original
- Aplicar o desconto e retornar o preço final
- Validar se o desconto está entre 0-100%

**Estrutura de Dados:**
- Produto: { id, nome, preco, categoria }
- Retorno: { precoOriginal, desconto, precoFinal, economia }

**Exemplo de Uso:**
\```javascript
const produto = { id: 1, nome: "Notebook", preco: 2000, categoria: "tech" };
const resultado = aplicarDesconto(produto, 15);
// resultado: { precoOriginal: 2000, desconto: 300, precoFinal: 1700, economia: 300 }
\```

**Considerações:**

- Tratar casos de erro (desconto inválido, produto sem preço)
- Usar precisão decimal adequada para valores monetários
- Incluir documentação JSDoc

```

---

## 🔧 Configurando um Provedor de IA

Para usar provedores reais de IA (OpenAI, Claude, etc.), configure uma chave de API:

### Configuração Rápida

```bash
# Exemplo com OpenAI
export OPENAI_API_KEY="sk-sua-chave-aqui"
grompt
```

### Teste com Provedor Real

```bash
# Via CLI
grompt ask "Como implementar autenticação JWT?" \
  --provider openai \
  --model gpt-4

# Via interface web
# 1. Vá para Configurações (ícone de engrenagem)
# 2. Adicione sua chave de API
# 3. Selecione o provedor desejado
```

## 📖 Exemplos Práticos

### Exemplo 1: Prompt de Código

**Suas ideias:**

- Sistema de login
- Com React e Node.js
- Usando JWT
- Banco PostgreSQL

**CLI:**

```bash
grompt generate \
  --idea "Sistema de login" \
  --idea "Com React e Node.js" \
  --idea "Usando JWT" \
  --idea "Banco PostgreSQL" \
  --purpose code \
  --provider openai
```

### Exemplo 2: Análise de Dados

**Suas ideias:**

- Analisar vendas do trimestre
- Identificar produtos mais vendidos
- Criar gráficos
- Python pandas

**CLI:**

```bash
grompt generate \
  --idea "Analisar vendas do trimestre" \
  --idea "Identificar produtos mais vendidos" \
  --idea "Criar gráficos" \
  --idea "Python pandas" \
  --purpose analysis
```

### Exemplo 3: Geração de Squad

```bash
# Gerar time de agentes de IA para seu projeto
grompt squad "Preciso de um app mobile para delivery de comida, com pagamento integrado, notificações push e painel administrativo web"
```

Isso cria um arquivo `AGENTS.md` com agentes especializados recomendados.

## 🎪 Casos de Uso Comuns

### Para Desenvolvedores

```bash
# Revisão de código
grompt generate --idea "revisar código React" --idea "melhorar performance" --idea "detectar bugs"

# Documentação
grompt generate --idea "documentar API REST" --idea "incluir exemplos" --idea "especificar erros"

# Testes
grompt generate --idea "testes unitários" --idea "Jest framework" --idea "cobertura 90%"
```

### Para Criação de Conteúdo

```bash
# Blog post
grompt generate --idea "artigo sobre microservices" --idea "para iniciantes" --idea "com exemplos práticos"

# Marketing
grompt generate --idea "copy para landing page" --idea "SaaS B2B" --idea "foco em conversão"
```

### Para Análise

```bash
# Business Intelligence
grompt generate --idea "dashboard de vendas" --idea "KPIs principais" --idea "Power BI"

# Pesquisa
grompt generate --idea "análise sentimento" --idea "redes sociais" --idea "Python NLTK"
```

## 🔍 Comandos Essenciais

| Comando | Uso | Exemplo |
|---------|-----|---------|
| `grompt` | Iniciar servidor web | `grompt --port 3000` |
| `grompt ask` | Pergunta direta à IA | `grompt ask "Como fazer deploy?"` |
| `grompt generate` | Gerar prompt a partir de ideias | `grompt generate --idea "..."` |
| `grompt squad` | Gerar squad de agentes | `grompt squad "projeto..."` |
| `grompt --help` | Ajuda geral | `grompt --help` |

## ⚙️ Opções Importantes

### Provedores Disponíveis

```bash
--provider demo      # Modo demonstração (padrão)
--provider openai    # OpenAI GPT models
--provider claude    # Anthropic Claude
--provider gemini    # Google Gemini
--provider deepseek  # DeepSeek AI
--provider ollama    # Modelos locais via Ollama
```

### Propósitos de Prompt

```bash
--purpose code      # Prompts para geração de código
--purpose analysis  # Prompts para análise de dados
--purpose creative  # Prompts para conteúdo criativo
--purpose general   # Prompts de uso geral (padrão)
```

### Configuração de Modelo

```bash
--model gpt-4              # Modelo específico
--max-tokens 2000          # Limite de tokens
--temperature 0.7          # Criatividade (0.0-1.0)
```

## 🛠️ Solução Rápida de Problemas

### Porta já em uso

```bash
grompt --port 8081
```

### Problemas de permissão

```bash
chmod +x grompt
```

### API key não funciona

```bash
# Testar conectividade
grompt ask "teste" --provider openai --dry-run
```

### Ver logs detalhados

```bash
DEBUG=true grompt
```

## 📚 Próximos Passos

Agora que você tem o básico funcionando:

1. **[Comandos CLI](../user-guide/cli-commands.md)** - Explore todos os comandos disponíveis
2. **[Configuração](../user-guide/configuration.md)** - Configure provedores de IA e personalize
3. **[Exemplos de Uso](../user-guide/usage-examples.md)** - Veja exemplos avançados e práticos
4. **[API Reference](../user-guide/api-reference.md)** - Integre o Grompt em seus próprios projetos

---

**Dica:** 💡 Mantenha o Grompt rodando em uma aba do navegador para acesso rápido durante o desenvolvimento!
