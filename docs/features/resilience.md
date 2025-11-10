# Resilience Modes

**Never fails, always works**: Grompt implements a hierarchical fallback system ensuring the binary **never fails** due to missing API keys.

---

## The Problem

Traditional AI tools fail when:

- API keys are missing or invalid
- Rate limits are exceeded
- Provider services are down
- Network connectivity issues occur

This creates **friction** and **negative user experiences**.

---

## Grompt's Solution: 3-Tier Fallback

Grompt implements a **hierarchical resilience system** with three modes:

### 1. 🔑 BYOK Mode (Priority 1)

**Bring Your Own Key** - Use your external API key per request.

- **How**: Provide key via HTTP header or UI toggle
- **Storage**: Zero - keys never stored
- **Use Case**: Personal projects, cost control, maximum security
- **Indicator**: Blue badge "🔑 Using Your API Key (BYOK)"

```bash
curl -X POST http://localhost:8080/api/v1/unified \
  -H "X-API-Key: sk-your-key" \
  -d '{"prompt":"Hello","provider":"openai"}'
```

### 2. 🔧 Server Mode (Priority 2)

**Server Configuration** - Use environment variables set on the server.

- **How**: Set `OPENAI_API_KEY`, `CLAUDE_API_KEY`, etc.
- **Storage**: Server memory (process environment)
- **Use Case**: Shared environments, team projects
- **Indicator**: Green badge "🔧 Using Server Config (openai)"

```bash
export OPENAI_API_KEY=sk-...
./grompt start
```

### 3. 🎭 Demo Mode (Priority 3 - Final Fallback)

**Template-Based Responses** - Always available, no API key needed.

- **How**: Automatic fallback when no keys available
- **Storage**: N/A - uses built-in templates
- **Use Case**: Testing, learning, demonstrations
- **Indicator**: Yellow badge "🎭 Demo Mode" with warning

---

## Hierarchical Priority Flow

```plaintext
┌─────────────────────────────────────┐
│  Request Received                   │
└──────────┬──────────────────────────┘
           │
           ▼
    ┌──────────────┐
    │ BYOK Key?    │
    │ (X-API-Key)  │
    └──┬───────┬───┘
       │ YES   │ NO
       ▼       ▼
    ┌─────┐   ┌──────────────┐
    │BYOK │   │ Server Key?  │
    │Mode │   │ (ENV var)    │
    └─────┘   └──┬───────┬───┘
                 │ YES   │ NO
                 ▼       ▼
              ┌─────┐  ┌──────┐
              │Srv  │  │ Demo │
              │Mode │  │ Mode │
              └─────┘  └──────┘
```

**Key principle**: The binary **never fails** - there's always a fallback!

---

## Mode Detection and Transparency

### Backend Response

Every API response includes the `mode` field:

```json
{
  "response": "Generated content...",
  "provider": "openai",
  "model": "gpt-4",
  "mode": "byok",
  "usage": {
    "total_tokens": 150
  }
}
```

### Frontend Indicators

The UI displays color-coded mode badges:

- **Blue** 🔑: "Using Your API Key (BYOK)"
- **Green** 🔧: "Using Server Config (gemini)"
- **Yellow** 🎭: "Demo Mode - Connect your API key for AI-powered prompts"

This **prevents confusion** and ensures users understand what's happening.

---

## Demo Mode Details

### What is Demo Mode?

When no API keys are available, Grompt generates **template-based responses** following Kubex principles:

- ✅ Modularity
- ✅ Anti-jargon language
- ✅ Practical, implementable solutions

### Example Demo Response

**Input**:

- Ideas: "REST API", "authentication", "PostgreSQL"
- Purpose: "Code Generation"

**Output**:

```markdown
# Code Generation Expert Assistant

## Primary Objective
Transform the provided ideas into actionable code generation solutions
following Kubex principles of modularity.

## User Requirements
- REST API
- authentication
- PostgreSQL

## Task Instructions
You are an expert code generation specialist. Based on the requirements above,
provide a comprehensive solution that:

### Key Requirements:
- Follow KUBEX principles: Modularity, No Vendor Lock-in, Practicality
- Use clear, anti-jargon language
- Provide modular, reusable components
- Ensure outputs are platform-agnostic

### Expected Output Format:
- Use Markdown for clear structure
- Include code examples when applicable
- Provide step-by-step instructions
- Add relevant comments and documentation

### Constraints:
- Avoid vendor lock-in solutions
- Keep complexity minimal
- Focus on practical, implementable solutions
- Use open standards and formats

## Context
This prompt was generated using Grompt, part of the Kubex Ecosystem,
following principles of modularity and clarity.

---
*Generated in demo mode - Connect your AI provider API key for enhanced AI-powered prompts*
```

### Demo Mode Benefits

- ✅ **Zero friction**: Works out-of-the-box
- ✅ **Educational**: Follows Kubex principles
- ✅ **Functional**: Provides usable structured prompts
- ✅ **Clear**: Visual indicator shows it's demo mode

---

## Use Cases by Mode

### BYOK Mode

**Best for:**

- Personal projects
- Sensitive environments
- Cost tracking per user
- Testing different providers

**Example**:

```bash
# Developer testing Claude vs OpenAI
curl ... -H "X-API-Key: sk-ant-..." -d '{"provider":"claude",...}'
curl ... -H "X-API-Key: sk-..." -d '{"provider":"openai",...}'
```

### Server Mode

**Best for:**

- Team environments
- Shared deployments
- Consistent provider selection
- Centralized billing

**Example**:

```bash
# Team server with OpenAI configured
export OPENAI_API_KEY=sk-team-key
./grompt start
# All team members use the same key automatically
```

### Demo Mode

**Best for:**

- First-time users
- Learning and exploration
- Offline/air-gapped environments
- Demonstrations and presentations

**Example**:

```bash
# No configuration needed!
./grompt start
# Open browser → Works immediately
```

---

## Error Handling

### BYOK Mode Errors

If you provide a BYOK key and it fails (invalid, rate limit, etc.):

- ❌ **Does NOT fall back** to server config or demo mode
- ✅ **Returns explicit error** message
- 📊 **HTTP Status**: 400/401/429 (depending on error)

**Rationale**: If you explicitly provide a key, you want to know if it fails.

### Server Mode Errors

If server config key fails:

- ✅ **Falls back to demo mode**
- 📊 **HTTP Status**: 200 (success)
- 🎭 **Mode**: "demo"
- ⚠️ **Warning indicator** in UI

### Demo Mode (Offline)

- ✅ **Never fails** - always returns template-based response
- 📊 **HTTP Status**: 200 (always success)
- 🎭 **Mode**: "demo"

---

## Configuration Examples

### Pure BYOK (No Server Config)

```bash
# No environment variables
./grompt start
```

- BYOK: ✅ Available
- Server: ❌ Not configured
- Demo: ✅ Available

### Server + BYOK

```bash
export OPENAI_API_KEY=sk-server-key
./grompt start
```

- BYOK: ✅ Available (overrides server)
- Server: ✅ Available (fallback)
- Demo: ✅ Available (final fallback)

### Demo Only (Testing/Learning)

```bash
# No config, just run
./grompt start
```

- BYOK: ✅ Available (user can provide)
- Server: ❌ Not configured
- Demo: ✅ Available (default)

---

## Monitoring Mode Usage

### Check Current Configuration

```bash
curl http://localhost:8080/api/v1/config
```

Response:

```json
{
  "providers": {
    "openai": {
      "available": true,
      "models": ["gpt-4", "gpt-3.5-turbo"]
    },
    "claude": {
      "available": false,
      "models": []
    }
  },
  "default_provider": "openai",
  "demo_mode": false
}
```

### API Response Includes Mode

Every response shows which mode was used:

```json
{
  "response": "...",
  "provider": "openai",
  "model": "gpt-4",
  "mode": "server"  ← Shows mode used
}
```

---

## Best Practices

### Development

Use **Demo Mode** for testing UX/flows without API costs:

```bash
./grompt start  # No config needed
```

### Staging

Use **Server Mode** with test API keys:

```bash
export OPENAI_API_KEY=sk-test-key
./grompt start
```

### Production

Offer **both BYOK and Server** for maximum flexibility:

```bash
# Server has fallback keys
export OPENAI_API_KEY=sk-prod-key
export CLAUDE_API_KEY=sk-ant-prod-key

# Users can override with BYOK
./grompt gateway start
```

---

## Benefits of This Approach

### 1. Zero Friction

Users can start using Grompt **immediately** without any configuration.

### 2. Maximum Flexibility

Choose the mode that fits your needs:

- Personal → BYOK
- Team → Server
- Testing → Demo

### 3. Prevents "Non-Issues"

Clear visual indicators prevent users from reporting "bugs" that are actually expected behavior (e.g., demo mode when no keys configured).

### 4. Cost Control

BYOK allows per-user cost tracking and control.

### 5. Never Fails

The binary **always works**, even with zero configuration.

---

## FAQ

**Q: What happens if all modes fail?**

A: **Demo mode never fails**. It's the final fallback and always returns a template-based response.

**Q: Can I disable demo mode?**

A: Not currently. Demo mode is the foundation of Grompt's resilience guarantee.

**Q: How do I know which mode is active?**

A: Check the visual indicator in the UI or the `mode` field in API responses.

**Q: Does BYOK override server config?**

A: **Yes.** BYOK has the highest priority.

**Q: Can I use different modes for different features?**

A: **Yes.** Each request is independent - use BYOK for one feature and server config for another.

---

## Next Steps

- [BYOK Feature](byok.md) - Deep dive into Bring Your Own Key
- [Quick Start](../getting-started/quick-start.md) - Try all three modes
- [Installation](../getting-started/installation.md) - Configure server mode

---

**Resilience**: Never fails | **Transparency**: Visual indicators | **Flexibility**: Choose your mode
