# Grompt

**Modern AI Prompt Engineering Platform**

Transform unstructured ideas into clean, effective prompts for AI models. Built with Go (backend) and React 19 (frontend), runs as a single binary with zero dependencies.

---

## Quick Start

```bash
# Download and run
./grompt start

# Open browser at http://localhost:8080
```

---

## Key Features

### 🔑 BYOK (Bring Your Own Key)
Use your own API keys per request - maximum flexibility and security.

### 🔧 Multi-Provider Support
- OpenAI (GPT-4, GPT-3.5)
- Anthropic Claude (3.5 Sonnet, Haiku)
- Google Gemini (2.0 Flash)
- DeepSeek
- Ollama (local)

### 💪 Resilience by Design
**Hierarchical Fallback:**
1. BYOK (Your API Key) → Priority
2. Server Config (ENV vars) → Fallback  
3. Demo Mode → Never fails!

### ⚡ Zero Dependencies
Single ~15MB binary - no Docker, no Node, no Python required.

---

## Use Cases

- **Prompt Engineering**: Structure ideas into professional prompts
- **Code Generation**: Generate scaffolds aligned with Kubex principles
- **Content Summarization**: Executive summaries and action plans
- **Image Prompts**: Briefings for image generation models
- **Chat Interface**: Conversational AI with context

---

## Links

- [GitHub Repository](https://github.com/kubex-ecosystem/grompt)
- [Getting Started](getting-started/installation.md)
- [BYOK Guide](features/byok.md)

---

**MIT License** • Part of the **Kubex Ecosystem**
