# Quick Start

Get started with **Grompt** in less than 60 seconds!

---

## Web Interface (Recommended)

### 1. Start the Server

```bash
./grompt start
```

Default address: **<http://localhost:8080>**

### 2. Open in Browser

Navigate to `http://localhost:8080` and you'll see the welcome screen with:

- **Prompt Crafter**: Transform ideas into structured prompts
- **Chat Interface**: Conversational AI
- **Code Generator**: Generate code scaffolds
- **Content Summarizer**: Create executive summaries
- **Image Prompt**: Generate image generation briefs

### 3. Choose Your Mode

Grompt offers **three resilience modes**:

1. **🔑 BYOK Mode** (Bring Your Own Key)
   - Click "🔑 Usar Sua Própria API Key (BYOK)" in any feature
   - Enter your API key (e.g., `sk-...` for OpenAI, `AIza...` for Gemini)
   - Your key is **never stored** - used only for that request

2. **🔧 Server Mode**
   - Set environment variables before starting (see [Installation](installation.md#configuration-optional))
   - Keys stored in server config

3. **🎭 Demo Mode**
   - Always available as fallback
   - No API key needed
   - Perfect for testing and learning

### 4. Try Prompt Crafter

1. Add ideas: "REST API", "authentication", "PostgreSQL"
2. Select purpose: "Code Generation"
3. Choose provider (or use Demo Mode)
4. Click **Generate Prompt**
5. Copy and use the structured prompt in any AI model!

---

## CLI Usage

### Ask a Question

```bash
grompt ask "What is quantum computing?" --provider gemini --model gemini-2.0-flash
```

### Generate Prompt from Ideas

```bash
grompt generate \
  --idea "REST API" \
  --idea "authentication" \
  --idea "rate limiting" \
  --purpose code
```

### Get AI Squad Recommendations

```bash
grompt squad "Build a payment microservice with Stripe integration"
```

Output:

```txt
🤖 Recommended AI Squad:

1. Backend Specialist (Claude 3.5 Sonnet)
   - Architecture design
   - Payment integration
   - Security implementation

2. Code Reviewer (GPT-4)
   - Code quality
   - Best practices
   - Documentation

3. Testing Engineer (Gemini 2.0 Flash)
   - Test coverage
   - Edge cases
   - Performance testing
```

---

## Custom Port

```bash
grompt start -p 5000
```

Now access at **<http://localhost:5000>**

---

## Gateway Mode (Enterprise)

For production deployments with rate limiting, metrics, and multi-tenant support:

```bash
grompt gateway start
```

Features:

- Rate limiting per API key
- Request metrics and logging
- CORS support
- Multi-provider routing

---

## API Usage

Grompt exposes REST endpoints when running in server mode:

### Unified Endpoint (All Providers)

```bash
curl -X POST http://localhost:8080/api/v1/unified \
  -H "Content-Type: application/json" \
  -H "X-API-Key: sk-your-key-here" \
  -d '{
    "prompt": "Explain quantum entanglement",
    "provider": "openai",
    "model": "gpt-4",
    "max_tokens": 500
  }'
```

### Provider-Specific Endpoints

```bash
# OpenAI
POST /api/openai

# Claude (Anthropic)
POST /api/claude

# Gemini
POST /api/gemini

# DeepSeek
POST /api/deepseek

# Ollama (local)
POST /api/ollama
```

### Check Configuration

```bash
curl http://localhost:8080/api/v1/config
```

Returns available providers and default settings.

---

## Next Steps

- **[BYOK Feature](../features/byok.md)** - Maximum security with per-request API keys
- **[Resilience Modes](../features/resilience.md)** - Understand the fallback hierarchy
- **[GitHub Repository](https://github.com/kubex-ecosystem/grompt)** - Source code and contributions

---

**Zero configuration required** | **Three resilience modes** | **Multi-provider support**
