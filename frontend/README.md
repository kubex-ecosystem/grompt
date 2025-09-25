# Grompt Frontend - React 19 + MultiProvider Architecture

Modern React frontend for the Grompt prompt engineering tool, featuring MultiProvider AI integration, PWA capabilities, and comprehensive offline support.

## Overview

The Grompt frontend is a cutting-edge React 19 application built with TypeScript, Vite, and TailwindCSS. It provides a seamless interface for prompt generation with support for multiple AI providers through both local SDK integration and backend gateway fallback.

### Key Features

- **🤖 MultiProvider AI Integration**: OpenAI, Anthropic, and Gemini support
- **📱 Progressive Web App**: Installable with offline capabilities
- **⚡ React 19**: Latest React features with concurrent rendering
- **🎨 Modern UI**: TailwindCSS with responsive design
- **🔄 Real-time Streaming**: SSE-based prompt generation with coalescing
- **💾 Offline-First**: IndexedDB caching and service worker support
- **🔐 Secure Configuration**: Local API key management with encryption

## Architecture

### Technology Stack

- **Framework**: React 19 + Vite + TypeScript
- **Styling**: TailwindCSS + PostCSS
- **State Management**: React hooks + Context API
- **HTTP Client**: Enhanced Fetch API with offline support
- **Database**: IndexedDB via `idb` library
- **PWA**: Service Worker + Web App Manifest

### Core Components

```
frontend/src/
├── components/          # Reusable UI components
│   ├── layout/         # Header, navigation, layout components
│   ├── providers/      # MultiProvider configuration UI
│   ├── pwa/           # PWA-specific components
│   └── settings/      # Configuration panels
├── core/              # Core business logic
│   └── llm/           # LLM provider implementations
│       ├── providers/ # Individual provider classes
│       └── wrapper/   # MultiAIWrapper abstraction
├── hooks/             # Custom React hooks
├── screens/           # Main application screens
├── services/          # API and service integrations
├── types/            # TypeScript type definitions
└── utils/            # Utility functions
```

## MultiProvider System

### Provider Support

The frontend supports multiple AI providers with seamless switching:

**OpenAI**:
- Models: GPT-4 Turbo, GPT-3.5 Turbo
- Features: Function calling, JSON mode, streaming
- SDK: Official OpenAI JavaScript SDK

**Anthropic**:
- Models: Claude 3.5 Sonnet, Claude 3 Haiku
- Features: Large context, system messages, streaming
- SDK: Official Anthropic SDK

**Gemini**:
- Models: Gemini 1.5 Pro, Gemini 1.5 Flash
- Features: Multimodal, safety settings, structured output
- SDK: Google Generative AI SDK

### Configuration Management

```typescript
// Provider configuration interface
interface MultiProviderConfig {
  providers: {
    [AIProvider.OPENAI]?: {
      apiKey: string
      defaultModel: string
      baseURL?: string
    }
    // ... other providers
  }
  fallbackToBackend?: boolean
  cacheResponses?: boolean
}
```

## PWA Features

### Offline Capabilities

- **Service Worker**: Caches static assets and API responses
- **IndexedDB**: Stores prompt history and provider configurations
- **Background Sync**: Queues requests when offline
- **Push Notifications**: Updates and alerts (when enabled)

### Installation

The app can be installed on desktop and mobile devices:

```javascript
// Installation prompt handling
const installPWA = async () => {
  if (deferredPrompt) {
    deferredPrompt.prompt()
    const { outcome } = await deferredPrompt.userChoice
    console.log(`User ${outcome} the install prompt`)
  }
}
```

## Development Setup

### Prerequisites

- Node.js 18+ (recommended: 20+)
- npm or yarn package manager
- Modern browser supporting ES2022

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Type checking
npm run typecheck

# Linting
npm run lint
```

### Environment Configuration

Create `.env.local` for development:

```env
# Development settings
VITE_APP_TITLE="Grompt Dev"
VITE_API_BASE_URL="http://localhost:3000"

# Optional: Default API keys for testing
VITE_OPENAI_API_KEY="sk-your-development-key"
VITE_ANTHROPIC_API_KEY="sk-ant-your-development-key"
VITE_GEMINI_API_KEY="your-gemini-development-key"
```

## Key Components Documentation

### MultiProviderConfig Component

```typescript
// Component for managing provider configurations
<MultiProviderConfig
  isOpen={showConfig}
  onClose={() => setShowConfig(false)}
  onConfigUpdate={(config) => {
    // Handle configuration updates
    console.log('Provider config updated:', config)
  }}
/>
```

### Enhanced API Service

```typescript
// Offline-first API with automatic fallback
const enhancedAPI = new EnhancedAPI({
  baseURL: 'http://localhost:3000',
  cacheEnabled: true,
  offlineMode: 'graceful'
})

// Generate prompt with automatic provider selection
const response = await enhancedAPI.generatePrompt({
  provider: 'openai',
  ideas: ['React', 'TypeScript', 'PWA'],
  purpose: 'code'
})
```

### Streaming Integration

```typescript
// Real-time streaming with SSE
await enhancedAPI.generatePromptStream(
  request,
  (chunk) => {
    // Handle streaming chunks
    setContent(prev => prev + chunk)
  },
  (usage) => {
    // Handle completion
    console.log('Tokens used:', usage.tokens)
  },
  (error) => {
    // Handle errors
    console.error('Stream error:', error)
  }
)
```

## Testing

### Unit Testing

```bash
# Run unit tests
npm test

# Run tests with coverage
npm run test:coverage

# Run tests in watch mode
npm run test:watch
```

### Integration Testing

```bash
# Test provider integrations
npm run test:providers

# Test PWA functionality
npm run test:pwa

# Test offline scenarios
npm run test:offline
```

### E2E Testing

```bash
# Run end-to-end tests
npm run test:e2e

# Test in multiple browsers
npm run test:cross-browser
```

## Build and Deployment

### Production Build

```bash
# Build for production
npm run build

# Preview production build
npm run preview

# Analyze bundle size
npm run analyze
```

### Build Outputs

```
dist/
├── assets/           # Compiled CSS, JS, and static assets
├── icons/           # PWA icons and favicons
├── manifest.json    # Web app manifest
├── sw.js           # Service worker
└── index.html      # Main HTML file
```

### Integration with Go Backend

The frontend is automatically embedded in the Go binary during build:

```go
//go:embed frontend/dist
var frontendAssets embed.FS

func serveFrontend() {
    // Serve embedded frontend assets
    handler := http.FileServer(http.FS(frontendAssets))
    http.Handle("/", handler)
}
```

## Performance Optimizations

### Bundle Optimization

- **Code Splitting**: Automatic route-based splitting
- **Tree Shaking**: Dead code elimination
- **Asset Optimization**: Image compression and optimization
- **Gzip Compression**: Build-time compression

### Runtime Performance

- **Virtual Scrolling**: For large lists (prompt history)
- **Memoization**: React.memo and useMemo for expensive operations
- **Debouncing**: Input debouncing for API calls
- **Lazy Loading**: Component and route lazy loading

### Caching Strategy

```typescript
// Multi-level caching
const cacheStrategy = {
  memory: new Map(),        // In-memory cache
  localStorage: window.localStorage,  // Browser storage
  indexedDB: await openDB('grompt'), // Persistent storage
  serviceWorker: await caches.open('api-cache') // HTTP cache
}
```

## Troubleshooting

### Common Issues

**Build Errors**:
```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install

# Clear Vite cache
npm run dev -- --force
```

**TypeScript Errors**:
```bash
# Restart TypeScript server
npm run typecheck
```

**Provider Connection Issues**:
- Check API keys in browser localStorage
- Verify network connectivity
- Check provider status endpoints

### Debug Mode

```typescript
// Enable debug mode
localStorage.setItem('debug', 'grompt:*')

// Check provider health
console.log(await multiProviderService.getProviderHealth())

// Inspect cache state
console.log(await enhancedAPI.getCacheStats())
```

## Contributing

### Code Style

- **Prettier**: Automatic code formatting
- **ESLint**: Code quality and consistency
- **TypeScript**: Strict type checking
- **Conventional Commits**: Standardized commit messages

### Development Workflow

1. **Feature Branch**: Create from `main`
2. **Development**: Make changes with tests
3. **Testing**: Run full test suite
4. **PR Review**: Submit for code review
5. **Integration**: Merge after approval

### File Naming Conventions

- **Components**: PascalCase (`MultiProviderConfig.tsx`)
- **Hooks**: camelCase with `use` prefix (`useMultiProvider.ts`)
- **Services**: camelCase (`multiProviderService.ts`)
- **Types**: PascalCase (`types.ts` with exported interfaces)

## Related Documentation

- [MultiProvider Architecture](../docs/MULTIPROVIDER_ARCHITECTURE.md)
- [Grompt V1 API Documentation](../GROMPT_V1_API.md)
- [PWA Integration Summary](PWA_INTEGRATION_SUMMARY.md)
- [Quick Start Guide](../QUICKSTART_V1.md)
