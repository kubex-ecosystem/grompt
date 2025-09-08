# Changelog

All notable changes to Grompt will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2025-09-08

### 🎉 Major Release - Welcome to Kubex.world

This is the official public launch of Grompt v2.0 on the new Kubex.world domain!

### ✨ Added

#### Core Features

- **Demo Mode**: App now works perfectly without API key using intelligent simulated responses
- **User API Key Input**: Users can securely add their own Gemini API key via UI
- **Multi-language Support**: Complete translations for EN, ES, PT-BR, ZH
- **Privacy-First Analytics**: Simple, GDPR-compliant tracking system
- **PWA Support**: Enhanced Progressive Web App capabilities

#### Infrastructure & Deployment

- **Production Domain**: Official deployment to kubex.world
- **CDN Independence**: Removed dependency on aistudiocdn.com, using esm.sh
- **Enhanced SEO**: Comprehensive meta tags, Open Graph, Twitter Cards
- **Security Headers**: Added security.txt and proper security configuration
- **Sitemap & Robots**: SEO optimization for search engines

#### Developer Experience

- **Environment Variables**: Extensive env var support for multiple AI providers
- **Better Error Handling**: Graceful fallbacks and user-friendly error messages
- **Code Documentation**: Comprehensive documentation and inline comments
- **Build Optimization**: Improved Vite configuration for production

#### UI/UX Improvements

- **Version Badge**: v2.0 indicator in header
- **Kubex Branding**: Integration with Kubex ecosystem identity
- **Enhanced Footer**: Links to ecosystem, GitHub, security, humans.txt
- **Loading States**: Better feedback during API calls
- **Responsive Design**: Improved mobile experience

### 🔄 Changed

#### Breaking Changes

- **API Key Handling**: Changed from required to optional with demo fallback
- **Build Configuration**: Updated Vite config for production deployment
- **CDN Sources**: Migrated from aistudiocdn to standard ESM CDN

#### Improvements

- **Performance**: Better bundle splitting and chunk optimization
- **Security**: Enhanced API key validation and storage
- **Accessibility**: Improved ARIA labels and keyboard navigation
- **Error Messages**: More helpful and actionable error descriptions

### 🐛 Fixed

- **White Screen Issue**: Resolved deployment issues when API key is missing
- **CDN Dependencies**: Removed reliance on non-standard CDN providers
- **Mobile Compatibility**: Fixed responsive design issues
- **Theme Persistence**: Improved dark/light mode persistence

### 🏗️ Technical Details

#### Stack Updates

- React 19.1.1
- Vite 6.2.0
- TypeScript 5.8.2
- Tailwind CSS (latest)
- Google Gemini AI integration

#### Deployment

- Vercel hosting on kubex.world
- Environment variable configuration
- Static asset optimization
- Progressive Web App manifest

#### Security

- Client-side API key storage only
- No server-side API key requirements
- HTTPS-only deployment
- Content Security Policy headers

### 🎯 Kubex Principles Applied

This release fully embodies the three core Kubex principles:

1. **Radical Simplicity**:
   - One-click deployment without configuration
   - Demo mode works immediately
   - Clear, jargon-free interface

2. **Modularity**:
   - Pluggable AI providers
   - Component-based architecture
   - Configurable through environment variables

3. **No Cages**:
   - No vendor lock-in (works with any Gemini-compatible API)
   - Open source codebase
   - Standard web technologies only
   - User controls their own API keys

### 🚀 What's Next

#### Planned for v2.1

- [ ] Additional AI provider support (OpenAI, Anthropic, Ollama)
- [ ] Advanced prompt templates
- [ ] Export functionality (PDF, Markdown)
- [ ] Collaborative features

#### Long-term Roadmap

- [ ] Plugin system for custom prompt types
- [ ] Integration with popular development tools
- [ ] Advanced analytics and usage insights
- [ ] Enterprise deployment options

---

## [1.0.8] - 2025-09-07

### Added

- Initial stable release
- Basic prompt generation with Gemini AI
- Theme switching (dark/light)
- Local storage for prompts and settings

### Fixed

- Various bug fixes and stability improvements

---

**Legend:**

- 🎉 Major release
- ✨ New features
- 🔄 Changes
- 🐛 Bug fixes
- 🏗️ Technical details
- 🎯 Philosophy & principles
- 🚀 Future plans
