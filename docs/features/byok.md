# BYOK (Bring Your Own Key)

**Maximum flexibility and security**: Use your own API keys per request without storing them on the server.

---

## What is BYOK?

**BYOK** (Bring Your Own Key) allows you to provide your own AI provider API keys **on a per-request basis**, ensuring:

- ✅ **Zero Storage**: Your API keys are **never stored** on the server
- ✅ **Maximum Security**: Keys exist only during the request lifecycle
- ✅ **Full Control**: You control costs, usage, and provider selection
- ✅ **Multi-Provider**: Works with all supported AI providers

---

## Supported Providers

BYOK works with:

- **OpenAI** (GPT-4, GPT-3.5, etc.) - `sk-...`
- **Anthropic Claude** (3.5 Sonnet, Haiku, etc.) - `sk-ant-...`
- **Google Gemini** (2.0 Flash, etc.) - `AIza...`
- **DeepSeek** - `sk-...`
- **ChatGPT** - API key format varies
- **Ollama** (local) - No key needed

---

## How to Use BYOK

### Web Interface

1. **Navigate to any feature**:
   - Prompt Crafter
   - Chat Interface
   - Code Generator
   - Content Summarizer
   - Image Prompt Generator

2. **Click the BYOK toggle**:

   ```txt
   🔑 Usar Sua Própria API Key (BYOK)
   ```

3. **Enter your API key**:
   - Input is a **password field** (hidden from view)
   - Format examples:
     - OpenAI: `sk-proj-abc123...`
     - Claude: `sk-ant-api03-xyz789...`
     - Gemini: `AIzaSyABC123...`

4. **Submit your request**:
   - Key is sent via HTTP header (`X-API-Key`)
   - Used **only for this request**
   - **Immediately discarded** after response

### API/CLI Usage

Send API keys via HTTP headers:

```bash
# Generic header (works for all providers)
curl -X POST http://localhost:8080/api/v1/unified \
  -H "Content-Type: application/json" \
  -H "X-API-Key: sk-your-key-here" \
  -d '{
    "prompt": "Explain quantum physics",
    "provider": "openai",
    "max_tokens": 500
  }'
```

**Provider-specific headers** (optional):

```bash
# OpenAI
-H "X-OPENAI-Key: sk-..."

# Claude
-H "X-CLAUDE-Key: sk-ant-..."

# Gemini
-H "X-GEMINI-Key: AIza..."

# DeepSeek
-H "X-DEEPSEEK-Key: sk-..."

# ChatGPT
-H "X-CHATGPT-Key: ..."
```

---

## Security Features

### 1. No Persistence

```go
// Backend: Keys are NEVER stored
func HandleUnified(w http.ResponseWriter, r *http.Request) {
    externalKey := r.Header.Get("X-API-Key")
    // Key used only in this function scope
    // Discarded when function returns
}
```

### 2. Password Input (Web UI)

```html
<input type="password" placeholder="sk-... ou AIza... (opcional)" />
```

Keys are hidden in the UI and never displayed in plain text.

### 3. HTTPS Recommended

For production deployments:

```bash
# Use reverse proxy with HTTPS (nginx, Caddy, etc.)
# Example with Caddy:
grompt.example.com {
    reverse_proxy localhost:8080
}
```

### 4. No Logs

API keys are **never logged** in server logs.

---

## Use Cases

### 1. Personal Projects

Use your personal API keys without sharing them with the server:

```txt
User → BYOK Key → Grompt → AI Provider
                   (no storage)
```

### 2. Team Environments

Each team member uses their own keys:

- **Developer A**: Uses personal OpenAI key
- **Developer B**: Uses personal Claude key
- **Server**: No keys configured (pure BYOK mode)

### 3. Cost Control

Track and control costs per user/project by using separate API keys.

### 4. Testing Providers

Test different providers without reconfiguring the server:

```bash
# Test OpenAI
curl ... -H "X-API-Key: sk-openai-key" -d '{"provider":"openai",...}'

# Test Claude
curl ... -H "X-API-Key: sk-ant-claude-key" -d '{"provider":"claude",...}'

# Test Gemini
curl ... -H "X-API-Key: AIza-gemini-key" -d '{"provider":"gemini",...}'
```

---

## Hierarchical Priority

Grompt uses a **hierarchical fallback** for API keys:

### Priority Order

1. **BYOK** (Highest Priority)
   - Key from HTTP header (`X-API-Key`, etc.)
   - Used if present
   - Mode: `byok`

2. **Server Config** (Fallback)
   - Environment variables (`OPENAI_API_KEY`, etc.)
   - Used if BYOK not provided
   - Mode: `server`

3. **Demo Mode** (Final Fallback)
   - No API key needed
   - Template-based responses
   - Mode: `demo`

### Example Flow

```txt
Request with X-API-Key header
  ↓
✅ BYOK Mode → Use provided key
  ↓ (if not present)
Server has OPENAI_API_KEY?
  ↓
✅ Server Mode → Use env var
  ↓ (if not present)
✅ Demo Mode → Template response
```

**Result**: The binary **never fails** due to missing API keys!

---

## Visual Indicators

The UI shows which mode is active:

- **🔑 BYOK Mode** (Blue): Using your external API key
- **🔧 Server Mode** (Green): Using server configuration
- **🎭 Demo Mode** (Yellow): Using demo/template responses

This prevents confusion and ensures transparency.

---

## Best Practices

### ✅ Do

- Use BYOK for personal projects and sensitive environments
- Rotate your API keys regularly
- Use HTTPS in production
- Monitor your API usage on provider dashboards

### ❌ Don't

- Share API keys between users
- Hard-code keys in client-side code
- Commit keys to version control
- Use the same key for development and production

---

## FAQ

**Q: Are my API keys stored on the server?**

A: **No.** BYOK keys are used only during the request and immediately discarded.

**Q: Can I use BYOK with all features?**

A: **Yes.** All features support BYOK: Prompt Crafter, Chat, Code Generator, Summarizer, and Image Prompt.

**Q: What happens if I provide an invalid key?**

A: The request will fail with an error message. Grompt will NOT fall back to server config or demo mode when you explicitly provide a BYOK key.

**Q: Can I mix BYOK and server config?**

A: **Yes.** Use BYOK when you need it, and fall back to server config otherwise. The hierarchical fallback handles this automatically.

**Q: Is BYOK required?**

A: **No.** BYOK is **optional**. You can also use:

- Server config (environment variables)
- Demo mode (no keys needed)

---

## Next Steps

- [Resilience Modes](resilience.md) - Understand the full fallback hierarchy
- [Quick Start](../getting-started/quick-start.md) - Try BYOK in action
- [Installation](../getting-started/installation.md) - Set up server config

---

**Security**: Zero storage | **Flexibility**: Per-request keys | **Transparency**: Visual mode indicators
