# Grompt Documentation

Comprehensive documentation for the Grompt prompt engineering tool, covering architecture, development, API usage, and deployment.

## 📚 Documentation Index

### Getting Started

- **[Quick Start Guide](../QUICKSTART_V1.md)** - Get up and running in 5 minutes
- **[Main Project README](../README.md)** - Project overview and basic setup
- **[Claude.md Instructions](../CLAUDE.md)** - Complete project reference for Claude Code

### Architecture Documentation

- **[MultiProvider Architecture](MULTIPROVIDER_ARCHITECTURE.md)** - Multi-AI provider system design
- **[SSE Streaming Architecture](SSE_STREAMING_ARCHITECTURE.md)** - Real-time streaming with coalescing
- **[Grompt V1 API Documentation](../GROMPT_V1_API.md)** - Complete V1 API technical reference

### Development

- **[Development Guide](DEVELOPMENT_GUIDE.md)** - Complete developer setup and workflows
- **[Frontend Documentation](../frontend/README.md)** - React frontend development guide
- **[PWA Integration Summary](../frontend/PWA_INTEGRATION_SUMMARY.md)** - Progressive Web App features

### Configuration

- **[Backend Configuration Example](../config/grompt_v1.example.yml)** - Production-ready configuration
- **[Frontend Environment Setup](../frontend/.env.example)** - Frontend environment variables
- **[API Testing Script](../scripts/test_grompt_v1.sh)** - Automated endpoint testing

## 🏗️ Architecture Overview

Grompt is built with a modern, scalable architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                   Frontend (React 19)                       │
│  • MultiProvider UI  • PWA Support  • Real-time Streaming  │
├─────────────────────────────────────────────────────────────┤
│                 Enhanced API Layer                          │
│     • Offline Caching  • Provider Fallback  • IndexedDB    │
├─────────────────────────────────────────────────────────────┤
│                MultiProvider System                         │
│  • Local SDK Execution  • Automatic Fallback  • Config UI  │
├─────────────────────────────────────────────────────────────┤
│                 Grompt V1 Gateway                           │
│    • Production Middleware  • SSE Coalescing  • Health      │
├─────────────────────────────────────────────────────────────┤
│              Provider Integrations                          │
│     • OpenAI SDK  • Anthropic SDK  • Gemini SDK           │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 Key Features

### Backend (Go)

- **Production-Ready Gateway**: Following Analyzer architecture patterns
- **Multi-Provider Support**: OpenAI, Anthropic, Gemini with official SDKs
- **SSE Streaming**: Intelligent chunk coalescing for smooth UX
- **Health Monitoring**: Comprehensive health checks and observability
- **Middleware Stack**: Timeout, concurrency, rate limiting, metrics

### Frontend (React)

- **Modern React 19**: Latest features with concurrent rendering
- **MultiProvider UI**: Seamless provider switching and configuration
- **Progressive Web App**: Installable with offline capabilities
- **Real-time Streaming**: SSE integration with smooth text display
- **Enhanced API**: Offline-first with intelligent caching

### AI Provider Integration

- **Local SDK Execution**: Direct provider integration for optimal performance
- **Automatic Fallback**: Graceful degradation to backend when needed
- **Configuration Management**: Secure API key storage and provider setup
- **Health Monitoring**: Real-time provider status and connectivity

## 🛠️ Development Workflow

### Quick Setup

```bash
# 1. Clone and setup
git clone https://github.com/kubex-ecosystem/grompt.git
cd grompt

# 2. Install dependencies
go mod download
cd frontend && npm install && cd ..

# 3. Configure environment
export OPENAI_API_KEY="your-key"
export GEMINI_API_KEY="your-key"

# 4. Build and run
make build
./dist/grompt gateway serve -p 3000
```

### Development Commands

**Backend Development**:
```bash
make build-dev    # Development build
make run          # Start development server
make test         # Run Go tests
```

**Frontend Development**:
```bash
cd frontend
npm run dev       # Start Vite dev server
npm run typecheck # TypeScript checking
npm run lint      # Code quality
npm test          # Run tests
```

**Testing**:
```bash
./scripts/test_grompt_v1.sh  # Test all V1 endpoints
curl http://localhost:3000/v1/health  # Health check
```

## 📋 API Reference

### Core Endpoints

- `POST /v1/generate` - Synchronous prompt generation
- `GET /v1/generate/stream` - Real-time SSE streaming
- `GET /v1/providers` - List available providers with health
- `GET /v1/health` - Comprehensive system health check
- `POST /v1/proxy/*` - Transparent proxy to GoBE backend

### Frontend API

```typescript
// MultiProvider service
const response = await multiProviderService.generateContent({
  provider: 'openai',
  ideas: ['AI', 'Streaming', 'React'],
  purpose: 'code'
})

// Enhanced API with offline support
await enhancedAPI.generatePromptStream(
  request,
  onChunk,
  onComplete,
  onError
)
```

## 🔧 Configuration

### Environment Variables

**Backend**:
```bash
export OPENAI_API_KEY="sk-your-openai-key"
export ANTHROPIC_API_KEY="sk-ant-your-anthropic-key"
export GEMINI_API_KEY="your-gemini-key"
export GOBE_BASE_URL="https://gobe.example.com"
```

**Frontend** (`.env.local`):
```env
VITE_API_BASE_URL="http://localhost:3000"
VITE_APP_TITLE="Grompt Development"
```

### Provider Configuration

The frontend provides a visual interface for configuring AI providers:

1. Click the Settings icon in the header
2. Configure API keys for each provider
3. Test connections and set default models
4. Save configuration (stored locally)

## 🧪 Testing

### Automated Testing

```bash
# Backend tests
make test

# Frontend tests
cd frontend && npm test

# Integration tests
./scripts/test_grompt_v1.sh
```

### Manual Testing

```bash
# Test streaming endpoint
curl -N -H "Accept: text/event-stream" \
  "http://localhost:3000/v1/generate/stream?provider=openai&ideas=test"

# Test provider health
curl http://localhost:3000/v1/providers

# Test frontend
open http://localhost:3000
```

## 📊 Performance

### Backend Optimizations

- **SSE Coalescing**: 75ms timeout with natural language boundaries
- **Concurrent Limiting**: 50 requests max with timeout controls
- **Provider Pooling**: Reusable connections and circuit breakers
- **Metrics Collection**: Structured logging and observability

### Frontend Optimizations

- **Bundle Splitting**: Automatic code splitting and lazy loading
- **PWA Caching**: Multi-level caching strategy (memory, localStorage, IndexedDB)
- **Streaming UI**: Smooth text display with intelligent buffering
- **Provider Selection**: Smart fallback and health-based routing

## 🛡️ Security

### API Key Management

- Local storage with basic encoding (frontend)
- Environment variables (backend)
- No API keys in logs or network requests
- Provider rotation on key failures

### Request Security

- CORS configuration for frontend access
- Request validation and sanitization
- Rate limiting and concurrent request controls
- Error boundaries preventing information leakage

## 🚀 Deployment

### Single Binary Deployment

```bash
# Production build
make build

# Deploy binary (includes embedded frontend)
./dist/grompt gateway serve -p 3000 -f config/config.yml
```

### Docker Deployment

```bash
# Build Docker image
docker build -t grompt:latest .

# Run with environment
docker run -p 3000:3000 \
  -e OPENAI_API_KEY="your-key" \
  grompt:latest
```

### Health Monitoring

- `/v1/health` - Comprehensive health checks
- Provider connectivity monitoring
- Metrics collection via structured logging
- Circuit breaker status and recovery

## 🐛 Troubleshooting

### Common Issues

**Build Errors**:
```bash
# Clean and rebuild
make clean && make build-dev
cd frontend && rm -rf node_modules && npm install
```

**Provider Connection Issues**:
- Check API keys in environment variables
- Verify provider endpoints are accessible
- Check `/v1/providers` for provider health status

**Streaming Problems**:
- Ensure browser supports EventSource API
- Check network connectivity for SSE
- Verify CORS configuration for cross-origin requests

### Debug Mode

```bash
# Enable backend debugging
export DEBUG=true
export LOG_LEVEL=debug

# Enable frontend debugging
localStorage.setItem('debug', 'grompt:*')
```

## 📈 Monitoring

### Metrics Available

- Request latency and throughput
- Provider health and availability
- Token usage and costs
- Error rates and types
- Streaming performance metrics

### Health Endpoints

- `/v1/health` - Overall system health
- `/v1/providers` - Provider-specific health
- Debug endpoints available in development mode

## 🤝 Contributing

### Development Process

1. **Setup**: Follow [Development Guide](DEVELOPMENT_GUIDE.md)
2. **Branch**: Create feature branch from `main`
3. **Develop**: Implement with tests and documentation
4. **Test**: Run full test suite
5. **PR**: Submit with clear description and testing instructions

### Code Standards

- **Go**: Follow standard Go conventions, comprehensive error handling
- **TypeScript**: Strict typing, React best practices, accessibility
- **Documentation**: Clear comments explaining "why", not "what"
- **Testing**: Unit tests for new functionality, integration tests for APIs

## 📞 Support

### Getting Help

- **Documentation**: Start with this index and linked guides
- **Issues**: GitHub Issues for bug reports and feature requests
- **Development**: [Development Guide](DEVELOPMENT_GUIDE.md) for setup help

### Reporting Issues

1. Check existing documentation and known issues
2. Include environment details (Go version, Node version, OS)
3. Provide reproduction steps and expected vs actual behavior
4. Include relevant log output and error messages

---

**Last Updated**: 2024-01-01
**Version**: 0.2.1
**Maintainers**: Kubex Ecosystem Team