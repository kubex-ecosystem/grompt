# Grompt: AI Prompt Crafter

![Grompt Logo](https://github.com/kubex-ecosystem/grompt/raw/main/docs/assets/logo.png)

Grompt is an intelligent AI prompt crafting tool built with Kubex principles: Radical Simplicity, Modularity, and No Cages. Transform your raw ideas into professional, structured prompts for AI models.

## ✨ Features

- 🧠 **Smart Prompt Engineering**: Transform raw ideas into structured, professional prompts
- 🎯 **Multiple Purposes**: Code generation, creative writing, data analysis, technical documentation, and more
- 🌙 **Dark/Light Theme**: Beautiful UI that adapts to your preference
- 🌍 **Multi-language**: Support for English, Spanish, Chinese, and Portuguese
- 💾 **Local Storage**: Save your prompts and ideas locally
- 🔗 **Shareable Links**: Share your prompts via URL
- 🔑 **Flexible API Key**: Works in demo mode or with your own Gemini API key
- 📱 **Responsive Design**: Works perfectly on desktop and mobile

## 🚀 Quick Start

### Option 1: Demo Mode (No API Key Required)

1. Clone the repository
2. Install dependencies: `npm install`
3. Start the development server: `npm run dev`
4. Open <http://localhost:5173>

The app will work in **demo mode** with simulated AI responses, perfect for testing and demonstrating the interface.

### Option 2: With Gemini API Key

1. Get a free API key from [Google AI Studio](https://ai.google.dev/)
2. Either:
   - **Environment Variable**: Create `.env.local` and add `GEMINI_API_KEY=your_key_here`
   - **User Input**: Enter your API key directly in the app interface (stored locally)
3. Run the app as above

## 🔧 Configuration

### Environment Variables

```bash
# Optional - app works in demo mode without this
GEMINI_API_KEY=your_gemini_api_key_here
```

### User API Key

If no environment variable is set, users can input their own API key directly in the interface. The key is:

- Stored only in localStorage (never sent to external servers)
- Validated for format (starts with "AIza", proper length)
- Used dynamically for API calls
- Can be cleared at any time

## 🏗️ Build and Deploy

```bash
# Build for production
npm run build

# Preview production build
npm run preview

# Build static files
npm run build:static
```

The app is designed to be deployed anywhere - Vercel, Netlify, GitHub Pages, or any static hosting service.

## 🎯 Kubex Principles

- **Radical Simplicity**: One command = one result. Direct, pragmatic, anti-jargon.
- **Modularity**: Well-structured, reusable components and outputs.
- **No Cages**: Platform-agnostic, open formats, no vendor lock-in.

## 🛠️ Technology Stack

- **React 19** with TypeScript
- **Vite** for blazing fast development
- **Tailwind CSS** for styling
- **Google Gemini AI** for prompt generation
- **Lucide React** for icons
- **No vendor-specific CDNs** - uses standard ESM.sh

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please read our [contributing guidelines](docs/CONTRIBUTING.md) first.

---

Built with ❤️ following Kubex principles: **CODE FAST. OWN EVERYTHING.**
