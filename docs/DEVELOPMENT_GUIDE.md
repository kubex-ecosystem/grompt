---
title: "Development Guide - Grompt V1"
version: 0.2.1
owner: kubex
audience: dev
languages: [en, pt-BR]
sources: ["Makefile", "frontend/package.json", "scripts/test_grompt_v1.sh"]
assumptions: ["Go 1.25+", "Node.js 18+", "Modern development environment"]
---
<!-- markdownlint-disable-next-line -->
# Development Guide - Grompt V1

## TL;DR

Complete development guide for Grompt V1, covering backend Go development, React frontend, MultiProvider integration, testing strategies, and deployment procedures. Includes commands, workflows, and best practices for contributing to the project.

## Development Environment Setup

### Prerequisites

**System Requirements**:

- Go 1.25.1+ (for backend development)
- Node.js 18+ with npm (for frontend development)
- Git for version control
- Modern IDE (VS Code, GoLand, or similar)

**Optional Tools**:

- Docker for containerized development
- Make utility (usually pre-installed on Unix systems)
- jq for JSON processing in shell scripts

### Initial Setup

```bash
# 1. Clone the repository
git clone https://github.com/kubex-ecosystem/grompt.git
cd grompt

# 2. Install Go dependencies (automatic via go.mod)
go mod download

# 3. Install frontend dependencies
cd frontend && npm install && cd ..

# 4. Set up environment variables
cp .env.example .env.local
# Edit .env.local with your API keys

# 5. Verify setup
make test
cd frontend && npm run typecheck && cd ..
```

### Environment Configuration

**Backend Environment Variables**:

```bash
# API Keys for providers
export OPENAI_API_KEY="sk-your-openai-key"
export ANTHROPIC_API_KEY="sk-ant-your-anthropic-key"
export GEMINI_API_KEY="your-gemini-key"

# Optional: GoBE integration
export GOBE_BASE_URL="https://gobe.example.com"

# Development settings
export DEBUG=true
export LOG_LEVEL=debug
```

**Frontend Environment Variables** (`.env.local`):

```env
# Development server
VITE_APP_TITLE="Grompt Dev"
VITE_API_BASE_URL="http://localhost:3000"

# Optional: Development API keys
VITE_OPENAI_API_KEY="sk-your-development-key"
VITE_ANTHROPIC_API_KEY="sk-ant-your-development-key"
VITE_GEMINI_API_KEY="your-gemini-development-key"
```

## Development Workflows

### Backend Development

#### Starting Development Server

```bash
# Option 1: Using make (recommended)
make run

# Option 2: Direct Go execution
go run cmd/main.go gateway serve -p 3000

# Option 3: Build and run
make build-dev
./dist/grompt gateway serve -p 3000
```

#### Backend Development Commands

```bash
# Core development
make build-dev          # Build development binary
make run                # Run application
make test              # Run Go tests
make clean             # Clean build artifacts

# Code quality
go vet ./...           # Static analysis
go fmt ./...           # Format code
golangci-lint run      # Comprehensive linting

# Specific module testing
go test ./internal/gateway/...    # Test gateway modules
go test ./internal/providers/...  # Test provider integrations
```

#### Backend Project Structure

```plaintext
cmd/
├── main.go                    # Application entry point
└── cli/                      # CLI command implementations

internal/
├── gateway/                  # Core gateway functionality
│   ├── transport/           # HTTP handlers and routing
│   │   ├── grompt_v1.go    # V1 API implementation
│   │   ├── sse_coalescer.go # Streaming optimization
│   │   └── http.go         # Core HTTP setup
│   ├── registry/           # Provider registry
│   └── middleware/         # Production middleware
├── providers/              # AI provider implementations
├── services/              # Business logic services
├── metrics/               # Observability and metrics
└── module/                # CLI framework and utilities

config/                    # Configuration files
scripts/                  # Development and deployment scripts
```

### Frontend Development

#### Starting Frontend Development

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies (first time)
npm install

# Start development server
npm run dev

# Alternative: Hot reload with backend
# Terminal 1: Start backend
make run

# Terminal 2: Start frontend dev server
cd frontend && npm run dev
```

#### Frontend Development Commands

```bash
# Core development
npm run dev            # Start Vite dev server
npm run build          # Build for production
npm run preview        # Preview production build

# Code quality
npm run typecheck      # TypeScript type checking
npm run lint           # ESLint linting
npm run format         # Prettier formatting

# Testing
npm test              # Run unit tests
npm run test:watch    # Watch mode testing
npm run test:coverage # Coverage report

# PWA and service worker
npm run build:sw      # Build service worker
npm run dev:https     # HTTPS dev server for PWA testing
```

#### Frontend Project Structure

```plaintext
frontend/src/
├── components/              # Reusable UI components
│   ├── layout/             # Layout components (Header, etc.)
│   ├── providers/          # MultiProvider configuration UI
│   ├── pwa/               # PWA-specific components
│   └── settings/          # Settings panels
├── core/                  # Core business logic
│   └── llm/               # LLM provider system
│       ├── providers/     # Individual AI providers
│       └── wrapper/       # MultiAIWrapper abstraction
├── hooks/                 # Custom React hooks
├── screens/               # Main application screens
├── services/              # API and service integrations
│   ├── api.ts            # Core API client
│   ├── enhancedAPI.ts    # Enhanced offline-capable API
│   └── multiProviderService.ts # MultiProvider integration
├── types/                 # TypeScript type definitions
└── utils/                 # Utility functions

frontend/public/           # Static assets
├── icons/                # PWA icons
├── manifest.json         # Web app manifest
└── sw.js                 # Service worker
```

## Testing Strategies

### Backend Testing

#### Unit Testing

```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/gateway/transport/
go test ./internal/providers/

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

#### Integration Testing

```bash
# Test V1 API endpoints
./scripts/test_grompt_v1.sh

# Manual endpoint testing
curl http://localhost:3000/v1/health
curl http://localhost:3000/v1/providers
curl -X POST http://localhost:3000/v1/generate \
  -H "Content-Type: application/json" \
  -d '{"provider":"openai","ideas":["test"],"purpose":"general"}'
```

#### Provider Testing

```bash
# Test individual providers
export OPENAI_API_KEY="your-key"
go test ./internal/providers/openai/ -v

export ANTHROPIC_API_KEY="your-key"
go test ./internal/providers/anthropic/ -v
```

### Frontend Testing

#### Unit and Component Testing

```bash
cd frontend

# Run all tests
npm test

# Watch mode for development
npm run test:watch

# Coverage report
npm run test:coverage

# Test specific files
npm test -- --testPathPattern=MultiProvider
npm test -- --testPathPattern=enhancedAPI
```

#### E2E Testing

```bash
# Start servers for E2E testing
# Terminal 1: Backend
make run

# Terminal 2: Frontend
cd frontend && npm run dev

# Terminal 3: E2E tests
cd frontend && npm run test:e2e
```

#### Provider Integration Testing

```bash
# Test MultiProvider functionality
cd frontend
npm run test:providers

# Test specific provider integrations
npm test -- --testPathPattern=openai
npm test -- --testPathPattern=anthropic
npm test -- --testPathPattern=gemini
```

## Code Quality and Standards

### Go Code Standards

#### Code Style

```go
// Package comments
// Package transport implements HTTP handlers for Grompt V1 API
package transport

// Struct documentation
// GromptV1Handlers handles all V1 API requests with production middleware
type GromptV1Handlers struct {
    registry             *registry.Registry
    productionMiddleware *middleware.ProductionMiddleware
}

// Function documentation with clear purpose
// NewGromptV1Handlers creates handlers with dependency injection
func NewGromptV1Handlers(reg *registry.Registry) *GromptV1Handlers {
    return &GromptV1Handlers{
        registry: reg,
    }
}
```

#### Error Handling

```go
// Standard error handling pattern
func (h *GromptV1Handlers) validateRequest(req *GenerateRequest) error {
    if req.Provider == "" {
        return fmt.Errorf("provider is required")
    }

    if len(req.Ideas) == 0 {
        return fmt.Errorf("at least one idea is required")
    }

    return nil
}

// Error response pattern
func (h *GromptV1Handlers) handleError(w http.ResponseWriter, err error, statusCode int) {
    log.Printf("Request error: %v", err)
    http.Error(w, err.Error(), statusCode)
}
```

### TypeScript Code Standards

#### Component Structure

```typescript
// Interface-first design
interface MultiProviderConfigProps {
  isOpen: boolean
  onClose: () => void
  onConfigUpdate: (config: MultiProviderConfig) => void
}

// Functional component with proper typing
export const MultiProviderConfig: React.FC<MultiProviderConfigProps> = ({
  isOpen,
  onClose,
  onConfigUpdate
}) => {
  // State with proper typing
  const [config, setConfig] = useState<MultiProviderConfig>({
    providers: {},
    fallbackToBackend: true,
    cacheResponses: true
  })

  // Event handlers
  const handleSave = useCallback(async () => {
    try {
      await multiProviderService.saveConfig(config)
      onConfigUpdate(config)
      onClose()
    } catch (error) {
      console.error('Failed to save config:', error)
    }
  }, [config, onConfigUpdate, onClose])

  return (
    // JSX implementation
  )
}
```

#### Service Implementation

```typescript
// Service class with proper error handling
export class MultiProviderService {
  private wrapper: MultiAIWrapper | null = null

  async generateContent(request: GenerateRequest): Promise<GenerateResponse> {
    try {
      // Try local providers first
      if (this.wrapper) {
        return await this.wrapper.generateContent(request)
      }
    } catch (error) {
      console.warn('Local provider failed, falling back to backend:', error)
    }

    // Fallback to backend
    return await this.backendGenerate(request)
  }
}
```

### Documentation Standards

#### Code Comments

```go
// Good: Explains why, not what
// Coalescing improves UX by reducing micro-chunks that create jittery text display
coalescer := NewSSECoalescer(flushFunc)

// Good: Documents non-obvious behavior
// hasNaturalBreak checks for punctuation that indicates good stopping points
// This prevents breaking words or phrases mid-sentence
func (c *SSECoalescer) hasNaturalBreak(content string) bool {
    // Implementation...
}
```

#### API Documentation

```typescript
/**
 * Generates AI content with automatic provider selection and fallback
 *
 * @param request - Generation request with provider preference
 * @returns Promise resolving to generated content with metadata
 * @throws {ProviderError} When all providers fail
 * @throws {ValidationError} When request parameters are invalid
 *
 * @example
 * ```typescript
 * const response = await service.generateContent({
 *   provider: 'openai',
 *   ideas: ['React', 'TypeScript'],
 *   purpose: 'code'
 * })
 * console.log(response.content)
 * ```
 */
async generateContent(request: GenerateRequest): Promise<GenerateResponse>
```

## Build and Deployment

### Development Builds

```bash
# Backend development build
make build-dev

# Frontend development build
cd frontend && npm run build

# Full development build (backend + frontend)
make build-dev && cd frontend && npm run build
```

### Production Builds

```bash
# Production build with optimizations
make build

# This automatically:
# 1. Builds optimized frontend
# 2. Embeds frontend in Go binary
# 3. Creates compressed binary
# 4. Runs basic validation
```

### Docker Development

```bash
# Build development Docker image
docker build -f Dockerfile.dev -t grompt:dev .

# Run with environment variables
docker run -p 3000:3000 \
  -e OPENAI_API_KEY="your-key" \
  -e GEMINI_API_KEY="your-key" \
  grompt:dev

# Docker Compose for full stack
docker-compose -f docker-compose.dev.yml up
```

### Cross-Platform Builds

```bash
# Build for multiple platforms
make build GOOS=linux GOARCH=amd64
make build GOOS=darwin GOARCH=arm64
make build GOOS=windows GOARCH=amd64

# Automated cross-platform builds
make build-all
```

## Debugging and Troubleshooting

### Backend Debugging

#### Development Debugging

```bash
# Enable debug logging
export DEBUG=true
export LOG_LEVEL=debug
make run

# Use Go debugging tools
go run -race cmd/main.go  # Race condition detection
dlv debug cmd/main.go     # Delve debugger
```

#### Common Issues

**Port Already in Use**:

```bash
# Find process using port 3000
lsof -i :3000

# Kill process
kill -9 $(lsof -t -i:3000)
```

**Module Issues**:

```bash
# Clean module cache
go clean -modcache
go mod download
```

### Frontend Debugging

#### Development Debugging (Frontend)

```bash
cd frontend

# Clear Vite cache
npm run dev -- --force

# Debug build issues
npm run build -- --debug

# Debug TypeScript issues
npm run typecheck
```

#### Browser Debugging

```typescript
// Enable debug mode in console
localStorage.setItem('debug', 'grompt:*')

// Check provider health
console.log(await multiProviderService.getProviderHealth())

// Inspect service worker
navigator.serviceWorker.getRegistrations().then(registrations => {
  console.log('SW registrations:', registrations)
})
```

### Performance Debugging

#### Backend Performance

```bash
# CPU profiling
go tool pprof http://localhost:3000/debug/pprof/profile

# Memory profiling
go tool pprof http://localhost:3000/debug/pprof/heap

# Goroutine debugging
go tool pprof http://localhost:3000/debug/pprof/goroutine
```

#### Frontend Performance

```typescript
// Bundle analysis
npm run build:analyze

// Performance monitoring
console.time('generation')
await service.generateContent(request)
console.timeEnd('generation')

// Memory usage
console.log('Memory:', performance.memory)
```

## Contributing Guidelines

### Development Workflow

1. **Feature Branch Creation**

   ```bash
   git checkout main
   git pull origin main
   git checkout -b feature/description
   ```

2. **Development Process**
   - Write tests first (TDD approach)
   - Implement feature with proper error handling
   - Ensure all tests pass
   - Update documentation

3. **Code Quality Checks**

   ```bash
   # Backend
   make test
   go vet ./...
   golangci-lint run

   # Frontend
   cd frontend
   npm run typecheck
   npm run lint
   npm test
   ```

4. **Commit Standards**

   ```bash
   # Use conventional commits
   git commit -m "feat: add streaming coalescing for better UX"
   git commit -m "fix: resolve provider selection edge case"
   git commit -m "docs: update MultiProvider architecture guide"
   ```

5. **Pull Request Process**
   - Create PR with descriptive title and body
   - Include testing instructions
   - Link related issues
   - Request appropriate reviewers

### Code Review Guidelines

#### Backend Review Checklist

- [ ] Error handling is comprehensive
- [ ] Resource cleanup (defer statements)
- [ ] Thread safety for concurrent operations
- [ ] Proper logging with structured data
- [ ] Performance considerations
- [ ] Security implications reviewed

#### Frontend Review Checklist

- [ ] TypeScript types are comprehensive
- [ ] React hooks follow best practices
- [ ] Error boundaries handle failures gracefully
- [ ] Accessibility considerations (a11y)
- [ ] Performance optimizations applied
- [ ] PWA functionality tested

### Release Process

#### Version Management

```bash
# Update version in relevant files
# - go.mod
# - frontend/package.json
# - internal/module/info/manifest.json

# Create release tag
git tag -a v1.2.0 -m "Release version 1.2.0"
git push origin v1.2.0
```

#### Release Checklist

- [ ] All tests pass in CI/CD
- [ ] Documentation updated
- [ ] Migration guides written (if needed)
- [ ] Performance benchmarks run
- [ ] Security review completed
- [ ] Release notes prepared

## How to Run / Repro

### Quick Development Setup

```bash
# 1-minute setup for new developers
git clone https://github.com/kubex-ecosystem/grompt.git
cd grompt
make build-dev && cd frontend && npm install && npm run build && cd ..
export GEMINI_API_KEY="your-key"
./dist/grompt gateway serve -p 3000

# Test in browser: http://localhost:3000
```

### Full Development Environment

```bash
# Complete setup with all tools
git clone https://github.com/kubex-ecosystem/grompt.git
cd grompt

# Backend setup
go mod download
make test

# Frontend setup
cd frontend
npm install
npm run typecheck
npm run lint
npm test
cd ..

# Environment configuration
cp .env.example .env.local
# Edit .env.local with your API keys

# Start development servers
# Terminal 1: Backend
make run

# Terminal 2: Frontend dev server
cd frontend && npm run dev

# Terminal 3: Run tests
./scripts/test_grompt_v1.sh
```

## Risks & Mitigations

### Development Risks

**Risk: API Key Exposure**
**Mitigation**: Use environment variables, never commit keys to repository

**Risk: Build Inconsistencies**
**Mitigation**: Use consistent Go and Node versions, containerized builds

**Risk: Provider Dependencies**
**Mitigation**: Implement comprehensive fallback mechanisms, mock providers for testing

**Risk: Frontend/Backend Version Mismatches**
**Mitigation**: Automated integration testing, version compatibility checks

## Next Steps

1. **Advanced Tooling**: Enhanced debugging tools and development utilities
2. **Automated Testing**: Comprehensive CI/CD pipeline with automated testing
3. **Performance Monitoring**: Real-time performance tracking in development
4. **Developer Experience**: Hot-reload improvements and faster feedback loops
5. **Documentation**: Interactive API documentation and tutorials

## Changelog

- **v0.2.1**: Comprehensive development guide with V1 architecture
- **v0.2.0**: Updated development workflows for MultiProvider system
- **v0.1.0**: Initial development guide
