# 🔑 BYOK (Bring Your Own Key) Guide

## Overview

Grompt supports **BYOK (Bring Your Own Key)** functionality, allowing users to provide their AI provider API keys per request via HTTP headers. This eliminates the need for server-side API key configuration and enables:

- ✅ **Zero server configuration** - No need to store API keys on the server
- ✅ **Per-request authentication** - Each request can use different keys
- ✅ **Client-side control** - Users maintain full control over their API keys
- ✅ **Multi-tenant support** - Perfect for SaaS applications
- ✅ **Browser compatibility** - Works seamlessly with web applications

---

## How It Works

### Header Priority

1. **External Key (via headers)** - Takes highest priority
2. **Server Configuration** - Fallback if no external key provided
3. **Demo Mode** - If neither external key nor server config exists

### Supported Headers

| Header | Provider | Example Value |
|--------|----------|---------------|
| `X-API-Key` | Generic (all providers) | `sk-...` or `AIza...` |
| `X-OPENAI-Key` | OpenAI | `sk-proj-...` |
| `X-CLAUDE-Key` | Anthropic Claude | `sk-ant-...` |
| `X-GEMINI-Key` | Google Gemini | `AIza...` |
| `X-DEEPSEEK-Key` | DeepSeek | `...` |
| `X-CHATGPT-Key` | ChatGPT | `sk-...` |

---

## Usage Examples

### 1. cURL Examples

#### OpenAI with Generic Header

```bash
curl -X POST http://localhost:8080/api/v1/unified \
  -H "Content-Type: application/json" \
  -H "X-API-Key: sk-proj-YOUR_OPENAI_KEY_HERE" \
  -d '{
    "provider": "openai",
    "model": "gpt-4o-mini",
    "prompt": "Explain quantum computing in simple terms",
    "max_tokens": 500
  }'
```

#### Gemini with Provider-Specific Header

```bash
curl -X POST http://localhost:8080/api/v1/unified \
  -H "Content-Type: application/json" \
  -H "X-GEMINI-Key: AIzaYOUR_GEMINI_KEY_HERE" \
  -d '{
    "provider": "gemini",
    "model": "gemini-2.0-flash-exp",
    "ideas": ["quantum computing", "beginner-friendly", "visual analogies"],
    "purpose": "Educational Content",
    "max_tokens": 1000
  }'
```

#### Claude (Anthropic) Example

```bash
curl -X POST http://localhost:8080/api/v1/unified \
  -H "Content-Type: application/json" \
  -H "X-CLAUDE-Key: sk-ant-YOUR_CLAUDE_KEY" \
  -d '{
    "provider": "claude",
    "model": "claude-3-5-sonnet-20241022",
    "prompt": "Write a technical documentation example",
    "max_tokens": 2000
  }'
```

---

### 2. JavaScript/TypeScript Examples

#### Fetch API

```typescript
async function generateWithBYOK(apiKey: string, prompt: string) {
  const response = await fetch('http://localhost:8080/api/v1/unified', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-API-Key': apiKey  // External key
    },
    body: JSON.stringify({
      provider: 'openai',
      model: 'gpt-4o-mini',
      prompt: prompt,
      max_tokens: 1000
    })
  });

  if (!response.ok) {
    throw new Error(`HTTP ${response.status}: ${await response.text()}`);
  }

  return await response.json();
}

// Usage
const result = await generateWithBYOK('sk-proj-...', 'Explain AI safety');
console.log(result.response);
```

#### React Hook Example

```typescript
import { useState } from 'react';

function useGromptBYOK() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const generate = async (apiKey: string, provider: string, prompt: string) => {
    setLoading(true);
    setError(null);

    try {
      const response = await fetch('/api/v1/unified', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-API-Key': apiKey
        },
        body: JSON.stringify({ provider, prompt, max_tokens: 1000 })
      });

      if (!response.ok) {
        throw new Error(`API Error: ${response.status}`);
      }

      return await response.json();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
      return null;
    } finally {
      setLoading(false);
    }
  };

  return { generate, loading, error };
}

// Component usage
function PromptGenerator() {
  const [apiKey, setApiKey] = useState('');
  const { generate, loading } = useGromptBYOK();

  const handleGenerate = async () => {
    const result = await generate(apiKey, 'openai', 'Hello world');
    console.log(result);
  };

  return (
    <div>
      <input
        type="password"
        placeholder="Your API Key"
        value={apiKey}
        onChange={(e) => setApiKey(e.target.value)}
      />
      <button onClick={handleGenerate} disabled={loading}>
        Generate
      </button>
    </div>
  );
}
```

---

### 3. Python Examples

#### Using Requests Library

```python
import requests

def generate_prompt_byok(api_key: str, provider: str, prompt: str):
    """Generate prompt using external API key"""
    url = "http://localhost:8080/api/v1/unified"
    headers = {
        "Content-Type": "application/json",
        "X-API-Key": api_key
    }
    payload = {
        "provider": provider,
        "prompt": prompt,
        "max_tokens": 1000
    }

    response = requests.post(url, json=payload, headers=headers)
    response.raise_for_status()
    return response.json()

# Usage
result = generate_prompt_byok(
    api_key="sk-proj-YOUR_KEY",
    provider="openai",
    prompt="Explain machine learning"
)
print(result['response'])
```

#### Async with HTTPX

```python
import httpx
import asyncio

async def generate_async(api_key: str, provider: str, prompt: str):
    async with httpx.AsyncClient() as client:
        response = await client.post(
            "http://localhost:8080/api/v1/unified",
            json={"provider": provider, "prompt": prompt, "max_tokens": 1000},
            headers={"X-API-Key": api_key}
        )
        response.raise_for_status()
        return response.json()

# Usage
result = asyncio.run(generate_async("sk-...", "openai", "Hello AI"))
```

---

## Security Best Practices

### ✅ DO

- **Use HTTPS** in production to encrypt API keys in transit
- **Validate API keys** on the client before sending requests
- **Store keys securely** using browser's `sessionStorage` or secure vaults
- **Clear keys** from memory after use
- **Use environment variables** for development

### ❌ DON'T

- Don't commit API keys to version control
- Don't store keys in `localStorage` (vulnerable to XSS)
- Don't log API keys in browser console
- Don't share keys between untrusted users
- Don't expose keys in client-side code

---

## Web UI Integration

The Grompt web interface includes built-in BYOK support:

1. Navigate to **Prompt Crafter** section
2. Click **"🔑 Use Your Own API Key (BYOK)"** button
3. Enter your API key (masked input)
4. Select provider and generate prompts
5. Key is sent only for that request and never stored

### Security Notice
>
> 💡 Your API key is only used for the current request and is **never stored on our servers**. The key is transmitted securely and discarded after use.

---

## Error Handling

### Common Errors

#### Missing API Key

```json
{
  "error": "OpenAI API Key not configured and not provided via X-API-Key header"
}
```

**Solution:** Provide a valid API key via `X-API-Key` header

#### Invalid API Key

```json
{
  "error": "Error in openai API: invalid_api_key"
}
```

**Solution:** Verify your API key is correct and active

#### Unsupported Provider

```json
{
  "error": "Unsupported provider: unknown_provider"
}
```

**Solution:** Use one of: `openai`, `claude`, `gemini`, `deepseek`, `chatgpt`, `ollama`

---

## Testing BYOK

### Quick Test Script

```bash
#!/bin/bash

# Test BYOK with OpenAI
API_KEY="sk-proj-YOUR_KEY_HERE"
PROMPT="Write a haiku about programming"

curl -X POST http://localhost:8080/api/v1/unified \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $API_KEY" \
  -d "{
    \"provider\": \"openai\",
    \"prompt\": \"$PROMPT\",
    \"max_tokens\": 100
  }" | jq .

# Test multiple providers
for provider in openai claude gemini; do
  echo "Testing $provider..."
  curl -s -X POST http://localhost:8080/api/v1/unified \
    -H "Content-Type: application/json" \
    -H "X-API-Key: $API_KEY" \
    -d "{
      \"provider\": \"$provider\",
      \"prompt\": \"Hello from $provider\",
      \"max_tokens\": 50
    }" | jq -r '.provider + ": " + .response'
done
```

---

## FAQ

### Q: Is my API key stored on the server?

**A:** No. External API keys are only used for the current request and are never persisted.

### Q: Can I use different keys for different requests?

**A:** Yes! Each request can use a different API key via headers.

### Q: What happens if I provide both server config and external key?

**A:** External key takes priority and will be used instead of server configuration.

### Q: Does BYOK work with all providers?

**A:** Yes! BYOK supports OpenAI, Claude, Gemini, DeepSeek, and ChatGPT.

### Q: Can I use BYOK with the web interface?

**A:** Yes! The web UI includes a built-in BYOK input field in the Prompt Crafter.

---

## Troubleshooting

### CORS Issues

If you encounter CORS errors when using BYOK from a web app:

1. Ensure you're using the correct headers
2. Check that Grompt server allows your origin
3. Verify CORS is properly configured in `handlers.go`

The server should include these headers:

```txt
Access-Control-Allow-Headers: Content-Type, Authorization, X-API-Key, X-OPENAI-Key, X-CLAUDE-Key, X-GEMINI-Key, X-DEEPSEEK-Key, X-CHATGPT-Key
```

### Network Issues

```bash
# Test if server is running
curl http://localhost:8080/api/health

# Test with verbose output
curl -v -X POST http://localhost:8080/api/v1/unified \
  -H "X-API-Key: test" \
  -d '{"provider":"openai","prompt":"test"}'
```

---

## Implementation Details

### Backend (Go)

The BYOK implementation in `handlers.go`:

1. Extracts API key from request headers
2. Checks generic `X-API-Key` first
3. Falls back to provider-specific headers (`X-OPENAI-Key`, etc.)
4. Creates temporary API client with external key
5. Uses external key for request, discards after response

### Frontend (TypeScript)

The frontend implementation in `unifiedAIService.ts`:

1. Accepts optional `apiKey` parameter
2. Adds `X-API-Key` header if key provided
3. Sends request with external authentication
4. Handles response and errors normally

---

## Support

For issues or questions:

- GitHub Issues: <https://github.com/kubex-ecosystem/grompt/issues>
- Documentation: <https://github.com/kubex-ecosystem/grompt/tree/main/docs>

---

**Version:** 1.1.0
**Last Updated:** 2025-01-20
**License:** MIT
