# Installation

**Grompt** is distributed as a single binary with zero dependencies. No Docker, no Node.js, no Python required.

---

## System Requirements

- **OS**: Linux, macOS, or Windows
- **Architecture**: amd64, arm64, or 386
- **Disk Space**: ~20MB
- **Memory**: ~20MB idle, ~50MB under load

---

## Quick Install

### Option 1: Download Pre-built Binary

Visit the [GitHub Releases](https://github.com/kubex-ecosystem/grompt/releases) page and download the binary for your platform:

```bash
# Linux (amd64)
wget https://github.com/kubex-ecosystem/grompt/releases/latest/download/grompt-linux-amd64
chmod +x grompt-linux-amd64
mv grompt-linux-amd64 /usr/local/bin/grompt

# macOS (arm64 - M1/M2/M3)
curl -L https://github.com/kubex-ecosystem/grompt/releases/latest/download/grompt-darwin-arm64 -o grompt
chmod +x grompt
mv grompt /usr/local/bin/grompt

# macOS (amd64 - Intel)
curl -L https://github.com/kubex-ecosystem/grompt/releases/latest/download/grompt-darwin-amd64 -o grompt
chmod +x grompt
mv grompt /usr/local/bin/grompt
```

### Option 2: Build from Source

Requires **Go 1.25.1+**:

```bash
git clone https://github.com/kubex-ecosystem/grompt.git
cd grompt
make build
```

Binary will be available at `./dist/grompt`.

### Option 3: Install to System PATH

```bash
git clone https://github.com/kubex-ecosystem/grompt.git
cd grompt
make install
```

This builds and installs to `/usr/local/bin/grompt`.

---

## Verify Installation

```bash
grompt --version
```

You should see output like:

```
Grompt v1.0.0
Modern AI Prompt Engineering Platform
Part of the Kubex Ecosystem
```

---

## Configuration (Optional)

Grompt works **without any configuration** thanks to:

- **Demo Mode**: Always available as fallback
- **BYOK Support**: Use your own API keys per request
- **Server Config**: Set environment variables for convenience

### Environment Variables

To use specific AI providers, set these **optional** variables:

```bash
# OpenAI
export OPENAI_API_KEY=sk-...

# Anthropic Claude
export CLAUDE_API_KEY=sk-ant-...

# Google Gemini
export GEMINI_API_KEY=...

# DeepSeek
export DEEPSEEK_API_KEY=...

# ChatGPT
export CHATGPT_API_KEY=...

# Ollama (local)
export OLLAMA_ENDPOINT=http://localhost:11434
```

!!! note "Optional Configuration"
    All environment variables are **optional**. Grompt uses a hierarchical fallback:

    1. **BYOK** (Your API Key via UI) → Priority
    2. **Server Config** (ENV vars) → Fallback
    3. **Demo Mode** → Never fails!

---

## Next Steps

- [Quick Start Guide](quick-start.md) - Get started in 60 seconds
- [BYOK Feature](../features/byok.md) - Use your own API keys
- [Resilience Modes](../features/resilience.md) - Understand fallback system

---

**Binary size**: ~15MB | **Memory**: ~20MB idle | **Dependencies**: Zero
