// Package server implements the HTTP server for the Prompt Crafter application.
package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	t "github.com/rafa-mori/grompt/internal/types"
)

//go:embed build/*
var reactApp embed.FS

type Server struct {
	config   *t.Config
	handlers *Handlers
}

func NewServer(cfg t.IConfig) *Server {
	handlers := NewHandlers(cfg)
	return &Server{
		config:   cfg.(*t.Config),
		handlers: handlers,
	}
}

func (s *Server) Start() error {
	// Configurar roteamento
	s.setupRoutes()

	url := fmt.Sprintf("http://localhost:%s", s.config.Port)

	fmt.Printf("🌐 Servidor iniciado em: %s\n", url)
	fmt.Printf("📁 Servindo aplicação React embarcada\n")
	fmt.Printf("🔧 APIs disponíveis:\n")
	fmt.Printf("   • /api/config - Configuração\n")
	fmt.Printf("   • /api/models - Modelos disponíveis\n")
	fmt.Printf("   • /api/test - Teste de API\n")
	fmt.Printf("   • /api/unified - Unified API\n")
	fmt.Printf("   • /api/openai - OpenAI API\n")
	fmt.Printf("   • /api/deepseek - DeepSeek API\n")
	fmt.Printf("   • /api/claude - Claude API\n")
	fmt.Printf("   • /api/ollama - Ollama Local\n")
	fmt.Printf("   • /api/health - Status do servidor\n")
	fmt.Printf("💡 Pressione Ctrl+C para parar\n\n")

	// Abrir navegador após delay
	go func() {
		time.Sleep(1 * time.Second)
		openBrowser(url)
	}()

	return http.ListenAndServe(":"+s.config.Port, nil)
}

func (s *Server) setupRoutes() {
	// EMBED REACT FRONTEND

	// Make sure the React build directory exists
	buildFS, err := fs.Sub(reactApp, "build")
	if err != nil {
		log.Printf("⚠️  Aviso: Não foi possível acessar arquivos React embarcados: %v", err)
		log.Printf("💡 Certifique-se de fazer 'npm run build' antes de compilar o Go")
		// Proceed with fallback routes
		// This will allow the server to run without the React frontend
		// and still serve the API endpoints.
		s.setupFallbackRoutes()
		return
	}

	// Handler that serve static files from the React build directory
	// This will serve files like index.html, main.js, styles.css, etc.
	staticHandler := http.FileServer(http.FS(buildFS))

	// Main route handler
	// This will handle all requests to the root path and serve the React app
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")

		// If the path starts with "api/", we return a 404 Not Found.
		// This prevents API routes from being handled by the React app
		// and ensures they are handled by the API handlers defined below.
		// This is important to avoid conflicts between API routes and React routing.
		if strings.HasPrefix(path, "api/") {
			http.NotFound(w, r)
			return
		}

		// Server static files directly if they exist
		// This allows serving files like /static/js/main.js, /static/css/styles.css, etc.
		// It checks if the path contains a dot (.) to identify files
		// and serves them directly from the build directory.
		if strings.Contains(path, ".") {
			// Check if the file exists in the embedded filesystem
			if _, err := fs.Stat(buildFS, path); err == nil {
				staticHandler.ServeHTTP(w, r)
				return
			}
		}

		// If the request is for the root path or any other path that doesn't match a file,
		// we serve the index.html file from the React build directory.
		// This allows the React app to handle routing internally.
		// It is important to serve index.html for all non-file requests
		// so that React Router can take over and handle the routing on the client side.
		indexFile, err := buildFS.Open("index.html")
		if err != nil {
			log.Printf("❌ Erro ao abrir index.html: %v", err)
			http.Error(w, "Frontend não disponível", http.StatusInternalServerError)
			return
		}
		defer indexFile.Close()

		// Set the Content-Type header to serve HTML
		// This is important to ensure the browser interprets the response as HTML
		// and renders the React app correctly.
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Check if the index.html file exists in the embedded filesystem
		// If it does, we read its content and write it to the response.
		// This is the main entry point for the React app and should be served for all
		// non-file requests to allow React Router to handle the routing.
		if _, err := fs.ReadFile(buildFS, "index.html"); err != nil {
			http.Error(w, "Erro ao ler frontend", http.StatusInternalServerError)
			return
		}

		// Read the content of index.html and write it to the response
		// This serves the React app for all non-file requests,
		// allowing React Router to handle the routing on the client side.
		content, err := fs.ReadFile(buildFS, "index.html")
		if err != nil {
			http.Error(w, "Erro ao carregar frontend", http.StatusInternalServerError)
			return
		}

		w.Write(content)
	})

	// API Routes
	// These routes handle API requests and are defined separately from the React app.
	// They are prefixed with "/api/" to distinguish them from the React app routes.
	// This allows the React app to handle client-side routing while the server handles API requests.
	// Each API route is handled by a specific handler function defined in the Handlers struct.
	http.HandleFunc("/api/claude", s.handlers.HandleClaude)
	http.HandleFunc("/api/openai", s.handlers.HandleOpenAI)
	http.HandleFunc("/api/deepseek", s.handlers.HandleDeepSeek)
	http.HandleFunc("/api/ollama", s.handlers.HandleOllama)
	http.HandleFunc("/api/unified", s.handlers.HandleUnified)
	http.HandleFunc("/api/models", s.handlers.HandleModels)
	http.HandleFunc("/api/agents", s.handlers.HandleAgents)
	http.HandleFunc("/api/agents/generate", s.handlers.HandleAgentsGenerate)
	http.HandleFunc("/api/agents/import", s.handlers.HandleAgentsImport)
	http.HandleFunc("/api/agents/export-advanced", s.handlers.HandleAgentsExportAdvanced)
	http.HandleFunc("/api/agents/validate", s.handlers.HandleAgentsValidate)
	http.HandleFunc("/api/agents/", s.handlers.HandleAgent)
	http.HandleFunc("/api/agents.md", s.handlers.HandleAgentsMarkdown)

	// This route handles the configuration API endpoint
	// It returns the server's configuration, such as API keys and endpoints.
	http.HandleFunc("/api/config", s.handlers.HandleConfig)

	// This route handles the test API endpoint
	// It is used to test the server's API functionality.
	http.HandleFunc("/api/test", s.handlers.HandleTest)

	// This route handles the health check for the server
	// It returns a simple JSON response indicating the server is healthy.
	http.HandleFunc("/api/health", s.handlers.HandleHealth)

	// Log the successful setup of routes
	// This log message indicates that the server has successfully set up the routes
	// and is ready to serve both the React app and the API endpoints.
	log.Println("✅ Rotas configuradas: Frontend React + APIs")
}

func (s *Server) setupFallbackRoutes() {

	// Fallback route for when the React frontend is not found
	// This route serves a simple HTML page explaining that the React frontend is not available
	// It provides instructions on how to build the React app and recompile the Go server.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Prompt Crafter - Setup Necessário</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background: #1a1a1a;
            color: #ffffff;
        }
        .container {
            background: #2d2d2d;
            padding: 30px;
            border-radius: 12px;
            border: 1px solid #404040;
        }
        h1 { color: #60a5fa; margin-bottom: 20px; }
        h2 { color: #34d399; margin-top: 30px; }
        pre {
            background: #1a1a1a;
            padding: 15px;
            border-radius: 8px;
            overflow-x: auto;
            border: 1px solid #404040;
        }
        code { color: #fbbf24; }
        .warning {
            background: #451a03;
            border: 1px solid #f59e0b;
            padding: 15px;
            border-radius: 8px;
            margin: 20px 0;
        }
        .step {
            background: #1e3a8a;
            border: 1px solid #3b82f6;
            padding: 15px;
            border-radius: 8px;
            margin: 15px 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Prompt Crafter</h1>

        <div class="warning">
            <strong>⚠️ Frontend React não encontrado!</strong><br>
            O servidor Go está rodando, mas o frontend React não foi embarcado no binário.
        </div>

        <h2>🔧 Como corrigir:</h2>

        <div class="step">
            <strong>Passo 1:</strong> Build do Frontend React
            <pre><code>cd frontend
npm install
npm run build
cd ..</code></pre>
        </div>

        <div class="step">
            <strong>Passo 2:</strong> Recompilar Go com Frontend Embarcado
            <pre><code>go build -o prompt-crafter .</code></pre>
        </div>

        <div class="step">
            <strong>Passo 3:</strong> Executar Novamente
            <pre><code>./prompt-crafter</code></pre>
        </div>

        <h2>📚 Ou use o Makefile:</h2>
        <pre><code>make build-all</code></pre>

        <h2>🔗 APIs Disponíveis:</h2>
        <ul>
            <li><a href="/api/health" style="color: #60a5fa;">/api/health</a> - Status do servidor</li>
            <li><a href="/api/config" style="color: #60a5fa;">/api/config</a> - Configuração das APIs</li>
        </ul>

        <p><strong>💡 Dica:</strong> Este servidor Go está funcionando corretamente. Você só precisa buildar e embarca o frontend React!</p>
    </div>
</body>
</html>`
		w.Write([]byte(html))
	})

	// Fallback API routes
	// These routes handle API requests when the React frontend is not available.
	// They provide basic functionality to ensure the server can still respond to API requests.
	http.HandleFunc("/api/models", s.handlers.HandleModels)
	http.HandleFunc("/api/claude", s.handlers.HandleClaude)
	http.HandleFunc("/api/ollama", s.handlers.HandleOllama)
	http.HandleFunc("/api/openai", s.handlers.HandleOpenAI)
	http.HandleFunc("/api/deepseek", s.handlers.HandleDeepSeek)
	http.HandleFunc("/api/unified", s.handlers.HandleUnified)
	http.HandleFunc("/api/agents", s.handlers.HandleAgents)
	http.HandleFunc("/api/agents/generate", s.handlers.HandleAgentsGenerate)
	http.HandleFunc("/api/agents/", s.handlers.HandleAgent)
	http.HandleFunc("/api/agents.md", s.handlers.HandleAgentsMarkdown)

	// Config route
	// This route returns the server's configuration, such as API keys and endpoints.
	// It is useful for clients to know how to interact with the server's APIs.
	http.HandleFunc("/api/config", s.handlers.HandleConfig)

	// Test route
	// This route is used to test the server's API functionality.
	// It can be used to verify that the server is running and responding correctly.
	http.HandleFunc("/api/test", s.handlers.HandleTest)

	// Health check route
	// This route checks the health of the server and returns a simple JSON response.
	// It is useful for monitoring and ensuring the server is running correctly.
	http.HandleFunc("/api/health", s.handlers.HandleHealth)

	// Log the fallback routes setup
	log.Println("⚠️  Fallback routes: Unavailable React frontend, serving API endpoints only")
}

func (s *Server) Shutdown() {
	fmt.Println("🧹 Cleaning resourses...")
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		fmt.Printf("🌐 Open your browser at: %s\n", url)
		return
	}

	if err != nil {
		fmt.Printf("⚠️  Error opening browser: %v\n", err)
		fmt.Printf("🌐 Open your browser at: %s\n", url)
	}
}

// Função para verificar se o build React existe
func (s *Server) checkReactBuild() bool {
	buildDir := "build"
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		return false
	}

	indexPath := filepath.Join(buildDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return false
	}

	return true
}
