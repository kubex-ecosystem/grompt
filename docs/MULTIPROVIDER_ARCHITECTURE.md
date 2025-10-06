---
title: "MultiProvider Architecture - Grompt V1"
version: 0.2.1
owner: kubex
audience: dev
languages: [en, pt-BR]
sources: ["frontend/src/core/llm/wrapper/MultiAIWrapper.ts", "frontend/src/services/multiProviderService.ts"]
assumptions: ["Provider SDKs remain stable", "LocalStorage availability"]
---

# MultiProvider Architecture - Grompt V1

## TL;DR

The Grompt MultiProvider system provides a unified interface for multiple AI providers (OpenAI, Anthropic, Gemini) with local SDK execution, automatic backend fallback, and comprehensive provider management. It enables seamless switching between providers with consistent API interfaces and enhanced user experience through caching and offline capabilities.

## Architecture Overview

### Core Components

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend MultiProvider                   │
├─────────────────────────────────────────────────────────────┤
│  MultiProviderConfig.tsx  │  multiProviderService.ts      │
│  (UI Configuration)       │  (Service Integration)         │
├─────────────────────────────────────────────────────────────┤
│                MultiAIWrapper.ts                           │
│                (Core Provider Logic)                        │
├─────────────────────────────────────────────────────────────┤
│  OpenAIProvider  │  AnthropicProvider  │  GeminiProvider   │
│  (Official SDK)  │  (Official SDK)     │  (Official SDK)   │
├─────────────────────────────────────────────────────────────┤
│              Backend Fallback (Grompt V1 API)              │
│         (When local providers fail or unavailable)          │
└─────────────────────────────────────────────────────────────┘
```

### Provider Implementation Stack

#### 1. MultiAIWrapper (`frontend/src/core/llm/wrapper/MultiAIWrapper.ts`)

**Purpose**: Core abstraction layer providing unified interface for all AI providers.

**Key Features**:
- Provider initialization and management
- Unified content generation interface
- Streaming support with standardized callbacks
- Response caching and optimization
- Error handling and retry logic

**Provider Support**:
```typescript
interface MultiAIConfig {
  providers: {
    [AIProvider.OPENAI]?: ProviderConfig
    [AIProvider.ANTHROPIC]?: ProviderConfig
    [AIProvider.GEMINI]?: ProviderConfig
  }
  caching?: boolean
  fallbackEnabled?: boolean
}
```

#### 2. Individual Provider Implementations

**OpenAI Provider** (`frontend/src/core/llm/providers/openai.ts`):
- Official OpenAI SDK integration
- GPT-4, GPT-3.5-Turbo support
- Streaming and non-streaming modes
- Token usage tracking

**Anthropic Provider** (`frontend/src/core/llm/providers/anthropic.ts`):
- Official Anthropic SDK integration
- Claude 3.5 Sonnet support
- SSE streaming implementation
- Advanced conversation handling

**Gemini Provider** (`frontend/src/core/llm/providers/gemini.ts`):
- Google Generative AI SDK integration
- Gemini Pro support
- Content safety filtering
- Structured response handling

#### 3. Service Integration (`frontend/src/services/multiProviderService.ts`)

**Purpose**: Bridges MultiAIWrapper with the Enhanced API and backend services.

**Key Responsibilities**:
- Provider configuration persistence
- API key management and validation
- Backend fallback coordination
- Provider health monitoring
- Request routing and load balancing

### Configuration Management

#### Provider Configuration Interface

```typescript
interface MultiProviderConfig {
  providers: {
    [provider: string]: {
      apiKey: string
      defaultModel: string
      baseURL?: string
      enabled?: boolean
      priority?: number
    }
  }
  fallbackToBackend: boolean
  cacheResponses: boolean
  maxRetries: number
  timeoutMs: number
}
```

#### Configuration Persistence

- **Storage**: localStorage with encrypted API keys
- **Validation**: Real-time provider connectivity testing
- **Sync**: Automatic synchronization with backend provider registry

### Execution Flow

#### 1. Content Generation Request

```
User Request → MultiProviderService → MultiAIWrapper → Provider Selection
     ↓
Local Provider Execution (Primary)
     ↓
Success → Response Caching → User
     ↓
Failure → Backend Fallback → Enhanced API → Grompt V1 API
```

#### 2. Provider Selection Logic

```typescript
class ProviderSelector {
  selectProvider(request: GenerateRequest): AIProvider {
    // 1. Check user preference
    if (request.provider && isProviderAvailable(request.provider)) {
      return request.provider
    }

    // 2. Check provider health and availability
    const healthyProviders = getHealthyProviders()

    // 3. Select based on priority and load
    return selectByPriority(healthyProviders)
  }
}
```

## Provider-Specific Features

### OpenAI Integration

**Models Supported**:
- gpt-4-turbo-preview
- gpt-4-0125-preview
- gpt-3.5-turbo-0125

**Features**:
- Function calling support
- JSON mode responses
- Vision capabilities (where supported)
- Advanced streaming with delta responses

### Anthropic Integration

**Models Supported**:
- claude-3-5-sonnet-20241022
- claude-3-haiku-20240307

**Features**:
- System message optimization
- Advanced reasoning capabilities
- Large context window support
- Streaming with SSE

### Gemini Integration

**Models Supported**:
- gemini-1.5-pro-latest
- gemini-1.5-flash

**Features**:
- Multimodal input support
- Safety settings configuration
- Structured output generation
- Content filtering and safety

## Error Handling and Resilience

### Multi-Level Fallback Strategy

```
1. Primary Provider (Local SDK)
   ↓ (on failure)
2. Alternative Local Provider
   ↓ (on failure)
3. Backend Provider (Grompt V1 API)
   ↓ (on failure)
4. Cached Response (if available)
   ↓ (on failure)
5. Error Response with User Guidance
```

### Error Categories and Responses

```typescript
enum ErrorType {
  API_KEY_INVALID = 'api_key_invalid',
  RATE_LIMIT = 'rate_limit',
  NETWORK_ERROR = 'network_error',
  PROVIDER_UNAVAILABLE = 'provider_unavailable',
  CONTENT_FILTER = 'content_filter',
  UNKNOWN = 'unknown'
}
```

### Retry Logic

- **Exponential backoff**: 1s, 2s, 4s, 8s intervals
- **Circuit breaker**: Temporary provider disabling after repeated failures
- **Health checks**: Periodic provider availability testing

## Performance Optimizations

### Response Caching

```typescript
interface CacheStrategy {
  keyGeneration: (request: GenerateRequest) => string
  ttl: number // Time to live in milliseconds
  storage: 'memory' | 'localStorage' | 'indexedDB'
  compression: boolean
}
```

### Streaming Optimizations

- **Chunk coalescence**: Intelligent buffering for smooth UX
- **Progressive loading**: Immediate response start with incremental updates
- **Cancellation support**: Request abortion and cleanup

### Resource Management

- **Connection pooling**: Reusable HTTP connections
- **Memory management**: Automatic cache cleanup and garbage collection
- **Concurrent request limiting**: Maximum parallel requests per provider

## Integration Points

### With Enhanced API

```typescript
class EnhancedAPIIntegration {
  async generatePrompt(request: GenerateRequest): Promise<GenerateResponse> {
    // Try multi-provider service first
    try {
      return await multiProviderService.generateContent(request)
    } catch (error) {
      // Fallback to backend
      return await this.backendGenerate(request)
    }
  }
}
```

### With Grompt V1 API

```typescript
class BackendFallback {
  async executeBackendGeneration(request: GenerateRequest): Promise<GenerateResponse> {
    const response = await fetch('/v1/generate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request)
    })

    return response.json()
  }
}
```

## How to Run / Repro

### Setup MultiProvider Configuration

```bash
# 1. Configure provider API keys
export OPENAI_API_KEY="sk-your-openai-key"
export ANTHROPIC_API_KEY="sk-ant-your-anthropic-key"
export GEMINI_API_KEY="your-gemini-key"

# 2. Start frontend development
cd frontend && npm run dev

# 3. Access provider configuration
# Navigate to app → Settings icon → Configure providers
```

### Test Provider Functionality

```typescript
// Test individual provider
const wrapper = new MultiAIWrapper(config)
const response = await wrapper.generateContent({
  provider: AIProvider.OPENAI,
  prompt: "Explain quantum computing",
  model: "gpt-4-turbo-preview"
})

// Test fallback mechanism
const service = new MultiProviderService()
const result = await service.generateContent({
  provider: "invalid_provider", // Will trigger fallback
  ideas: ["AI", "Machine Learning"],
  purpose: "general"
})
```

### Monitor Provider Health

```bash
# Check provider status
curl http://localhost:3000/v1/providers

# Expected response:
{
  "providers": [
    {
      "name": "openai",
      "available": true,
      "models": ["gpt-4-turbo-preview", "gpt-3.5-turbo-0125"],
      "health": "healthy"
    }
  ]
}
```

## Risks & Mitigations

### Risk: API Key Exposure
**Mitigation**:
- Store keys in localStorage with basic encoding
- Never log API keys in console or network requests
- Implement key rotation capabilities

### Risk: Provider Service Interruption
**Mitigation**:
- Multi-provider fallback strategy
- Response caching for offline scenarios
- Circuit breaker pattern for failing providers

### Risk: Rate Limiting
**Mitigation**:
- Request queuing and throttling
- Provider rotation on rate limits
- Exponential backoff with jitter

### Risk: Inconsistent Provider Responses
**Mitigation**:
- Response normalization layer
- Standardized error handling
- Unified response schema

### Risk: Configuration Complexity
**Mitigation**:
- Intuitive UI for provider configuration
- Automatic provider discovery and setup
- Sensible defaults and validation

## Next Steps

1. **Enhanced Provider Discovery**: Automatic detection of available providers and models
2. **Advanced Caching**: Intelligent cache invalidation and response similarity detection
3. **Provider Analytics**: Usage statistics and performance monitoring dashboard
4. **Custom Provider Support**: Plugin system for additional AI providers
5. **Load Balancing**: Intelligent request distribution across multiple providers

## Changelog

- **v0.2.1**: Initial MultiProvider architecture implementation
- **v0.2.0**: Core provider abstractions and fallback mechanisms
- **v0.1.0**: Basic provider integration framework