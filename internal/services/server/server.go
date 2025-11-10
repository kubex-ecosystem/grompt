// Package server implements the HTTP server for the Prompt Crafter application.
package server

import (
	"fmt"
	"io/fs"
	"mime"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kubex-ecosystem/grompt/internal/gateway"
	"github.com/kubex-ecosystem/grompt/internal/gateway/middleware"
	"github.com/kubex-ecosystem/grompt/internal/gateway/registry"
	"github.com/kubex-ecosystem/grompt/internal/gateway/routes"
	"github.com/kubex-ecosystem/grompt/internal/grompt"
	"github.com/kubex-ecosystem/grompt/internal/interfaces"
	t "github.com/kubex-ecosystem/grompt/internal/types"
	gl "github.com/kubex-ecosystem/logz/logger"
)

var reactApp = grompt.NewGUIGrompt()

type Server struct {
	*gateway.ServerImpl
	apiRouter *gin.RouterGroup
	config   *t.Config
	handlers *Handlers
}

type ReactApp struct {
	FS          []fs.DirEntry
	Wasms       *[]fs.File
	ReactRoutes map[string]string
	WasmRoutes  map[string]string
}

func NewServer(cfg interfaces.IConfig) *Server {
	if cfg == nil {
		gl.Log("error", "❌ Configuração inválida fornecida ao criar o servidor")
		return nil
	}

	if !cfg.IsDebugMode() {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	if cfg.IsDebugMode() {
		engine.Use(gin.Logger())
	}
	engine.Use(gin.Recovery())

	baseCfg, ok := cfg.(*t.Config)
	if !ok {
		gl.Log("error", "❌ Configuração recebida não é do tipo esperado")
		return nil
	}

	server, err := gateway.NewServer(cfg)
	if err != nil {
		gl.Log("error", "❌ Falha ao criar o servidor gateway: %v", err)
		return nil
	}

	return &Server{
		ServerImpl:   server,
		config:    baseCfg,
		handlers:  NewHandlers(cfg),
	}
}

func (s *Server) Start() error {
	// Configurar roteamento
	s.setupRoutes()

	url := fmt.Sprintf("http://localhost:%s", s.config.Port)

	gl.Log("info", "🌐 Gateway started: %s\n", url)
	gl.Log("info", "📁 Serving embedded React application\n")
	gl.Log("info", "🔧 Available APIs:\n")
	gl.Log("info", "   • /api/v1/config - Configuration\n")
	gl.Log("info", "   • /api/v1/models - Available Models\n")
	gl.Log("info", "   • /api/v1/test - API Test\n")
	gl.Log("info", "   • /api/v1/unified - Unified API\n")
	gl.Log("info", "   • /api/v1/openai - OpenAI API\n")
	gl.Log("info", "   • /api/v1/deepseek - DeepSeek API\n")
	gl.Log("info", "   • /api/v1/claude - Claude API\n")
	gl.Log("info", "   • /api/v1/gemini - Gemini API\n")
	gl.Log("info", "   • /api/v1/chatgpt - ChatGPT API\n")
	gl.Log("info", "   • /api/v1/ollama - Ollama Local\n")
	gl.Log("info", "   • /api/v1/health - Server Status\n")
	gl.Log("info", "💡 Press Ctrl+C to stop\n\n")

	// Detecta se há a página aberta em algum lugar

	// Abrir navegador após delay
	go func() {
		time.Sleep(1 * time.Second)
		openBrowser(url)
	}()

	return http.ListenAndServe(
		net.JoinHostPort(s.config.BindAddr, s.config.Port),
		s.ServerImpl.Router(),
	)
}

func getGinHandlerFunc(f http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(c.Writer, c.Request)
	}
}

func (s *Server) setupRoutes() {
	// Register API routes first (must be before NoRoute)
	s.setupAPIRoutes()
	s.setupGatewayRoutes()
	s.mountAPIRouter()

	// Static routes with NoRoute must be last
	s.setupStaticRoutes()
}

func (s *Server) setupStaticRoutes() {
	buildFS, err := fs.Sub(reactApp.GetWebFS(), "embedded/guiweb")
	if err != nil {
		gl.LoggerG.GetLogger().Log("warn", fmt.Sprintf("⚠️ build embed não encontrado: %v", err))
		_ = s.setupFallbackRoutes()
		return
	}

	// Serve all static files including assets
	// NoRoute is only for non-API routes (SPA fallback)
	fileServer := http.FileServer(http.FS(buildFS))
	s.ServerImpl.Router().NoRoute(func(c *gin.Context) {
		// Only serve static files for non-API routes
		// API routes (/api/*, /v1/*, /healthz) are already registered and will be handled before NoRoute
		fileServer.ServeHTTP(c.Writer, c.Request)
	})
}

func (s *Server) setupAPIRoutes() {
	s.ServerImpl.Router().GET("/api/v1/health", s.handlers.HandleHealth)
	s.ServerImpl.Router().GET("/api/v1/config", s.handlers.HandleConfig)
	s.ServerImpl.Router().POST("/api/v1/config", s.handlers.HandleConfig)
	s.ServerImpl.Router().GET("/api/v1/test", s.handlers.HandleTest)
	s.ServerImpl.Router().GET("/api/v1/models", s.handlers.HandleModels)

	// 2) Provedores (diretos)
	s.ServerImpl.Router().POST("/api/v1/openai", s.handlers.HandleOpenAI)
	s.ServerImpl.Router().POST("/api/v1/claude", s.handlers.HandleClaude)
	s.ServerImpl.Router().POST("/api/v1/gemini", s.handlers.HandleGemini)
	s.ServerImpl.Router().POST("/api/v1/deepseek", s.handlers.HandleDeepSeek)
	s.ServerImpl.Router().POST("/api/v1/chatgpt", s.handlers.HandleChatGPT)
	s.ServerImpl.Router().POST("/api/v1/ollama", s.handlers.HandleOllama)

	// 3) Geração Unificada e Atalhos
	s.ServerImpl.Router().POST("/api/v1/unified", s.handlers.HandleUnified)
	s.ServerImpl.Router().POST("/api/v1/ask", s.handlers.HandleAsk)
	s.ServerImpl.Router().POST("/api/v1/squad", s.handlers.HandleSquad)

	// 4) Agentes / Squad
	// s.GET("/api/v1/agents", getGinHandlerFunc(s.handlers.HandleAgents))
	// s.POST("/api/v1/agents", getGinHandlerFunc(s.handlers.HandleAgents))
	// s.POST("/api/v1/agents/generate", getGinHandlerFunc(s.handlers.HandleAgentsGenerate))
	// s.GET("/api/v1/agents/", getGinHandlerFunc(s.handlers.HandleAgent))
	// s.PUT("/api/v1/agents/", getGinHandlerFunc(s.handlers.HandleAgent))
	// s.DELETE("/api/v1/agents/", getGinHandlerFunc(s.handlers.HandleAgent))
	// s.GET("/api/v1/agents.md", getGinHandlerFunc(s.handlers.HandleAgentsMarkdown))

	// Página de teste para WASM
	s.ServerImpl.Router().GET(
		"/wasm-test.html",
		gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetStaticHeaders(w, "wasm-test.html")
		w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<meta charset="UTF-8" />
<title>WASM Test</title>
<body>
<h1>LookAtni WASM Test</h1>
<script type="module">
import init, { parse } from '/wasm/lookatni_wasm.js';
init('/wasm/lookatni_wasm_bg.wasm').then(() => {
  console.log('WASM init OK');
  console.log('parse("Hello") =>', parse("Hello"));
}).catch(e => console.error('WASM init FAIL', e));
</script>
</body>
</html>`,
		))
	})),
	)
}

func (s *Server) mountAPIRouter() {
	router := s.ServerImpl.GetRouter()
	s.apiRouter = router.Group("/api/v1", s.handlers.setCORSHeaders)
}

func (s *Server) setupGatewayRoutes() {
	reg, err := registry.FromRuntimeConfig(s.config)
	if err != nil {
		gl.Log("warn", "⚠️  Gateway runtime desabilitado: %v", err)
		return
	}

	prodCfg := middleware.DefaultProductionConfig()
	prodMiddleware := middleware.NewProductionMiddleware(prodCfg)
	for _, providerName := range reg.ListProviders() {
		prodMiddleware.RegisterProvider(providerName)
	}

	// Register gateway routes on main router (not apiRouter)
	// This allows routes like /healthz to be at root level
	mainRouter := s.ServerImpl.GetRouter()
	routes.NewGatewayRoutes(reg, prodMiddleware).Register(mainRouter.Group(""))
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
		gl.Log("info", "🌐 Open your browser at: %s\n", url)
		return
	}

	if err != nil {
		gl.Log("warn", "⚠️  Error opening browser: %v\n", err)
		gl.Log("info", "🌐 Open your browser at: %s\n", url)
	}
}

func (s *Server) setupFallbackRoutes() error {
	// Fallback route for when the React frontend is not found
	// This route serves a simple HTML page explaining that the React frontend is not available
	// It provides instructions on how to build the React app and recompile the Go server.
	s.ServerImpl.Router().GET(
		"/",
		gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/v1/") {
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
            <li><a href="/api/v1/health" style="color: #60a5fa;">/api/v1/health</a> - Status do servidor</li>
            <li><a href="/api/v1/config" style="color: #60a5fa;">/api/v1/config</a> - Configuração das APIs</li>
        </ul>

        <p><strong>💡 Dica:</strong> Este servidor Go está funcionando corretamente. Você só precisa buildar e embarca o frontend React!</p>
    </div>
</body>
</html>`
		w.Write([]byte(html))
	})),
	)

	// Fallback API routes
	// These routes handle API requests when the React frontend is not available.
	// They provide basic functionality to ensure the server can still respond to API requests.
	s.apiRouter.GET("/api/v1/models", s.handlers.HandleModels)
	s.apiRouter.GET("/api/v1/claude", s.handlers.HandleClaude)
	s.apiRouter.GET("/api/v1/ollama", s.handlers.HandleOllama)
	s.apiRouter.GET("/api/v1/openai", s.handlers.HandleOpenAI)
	s.apiRouter.GET("/api/v1/chatgpt", s.handlers.HandleChatGPT)
	s.apiRouter.GET("/api/v1/gemini", s.handlers.HandleGemini)
	s.apiRouter.GET("/api/v1/deepseek", s.handlers.HandleDeepSeek)
	s.apiRouter.GET("/api/v1/unified", s.handlers.HandleUnified)
	// s.router.HandleFunc("/api/v1/agents", s.handlers.HandleAgents)
	// s.router.HandleFunc("/api/v1/agents/generate", s.handlers.HandleAgentsGenerate)
	// s.router.HandleFunc("/api/v1/agents/", s.handlers.HandleAgent)
	// s.router.HandleFunc("/api/v1/agents.md", s.handlers.HandleAgentsMarkdown)

	// Config route
	// This route returns the server's configuration, such as API keys and endpoints.
	// It is useful for clients to know how to interact with the server's APIs.
	s.apiRouter.GET("/api/v1/config", s.handlers.HandleConfig)

	// Test route
	// This route is used to test the server's API functionality.
	// It can be used to verify that the server is running and responding correctly.
	s.apiRouter.GET("/api/v1/test", s.handlers.HandleTest)

	// Health check route
	// This route checks the health of the server and returns a simple JSON response.
	// It is useful for monitoring and ensuring the server is running correctly.
	s.apiRouter.GET("/api/v1/health", s.handlers.HandleHealth)

	// Log the fallback routes setup
	gl.Log("warn", "⚠️  Fallback routes: Unavailable React frontend, serving API endpoints only")

	return nil
}

func stopAllServices() {
	// Aqui você pode adicionar a lógica para parar outros serviços, se necessário
}

func (s *Server) Shutdown() {
	gl.Log("info", "🧹 Cleaning resources...")
	stopAllServices()
	gl.Log("info", "✅ Server shutdown complete.")
}

func SetStaticHeaders(w http.ResponseWriter, path string) {

	// registra MIME do .wasm globalmente (belt & suspenders)
	_ = mime.AddExtensionType(".wasm", "application/wasm")

	// --- TIPOS DE ARQUIVO E MIME ---
	var mimeByExt = map[string]string{
		".wasm": "application/wasm",
		".js":   "application/javascript; charset=utf-8",
		".mjs":  "application/javascript; charset=utf-8",
		".css":  "text/css; charset=utf-8",
		".json": "application/json; charset=utf-8",
		".svg":  "image/svg+xml",
		".ico":  "image/x-icon",
		".map":  "application/octet-stream",
		".txt":  "text/plain; charset=utf-8",
		".html": "text/html; charset=utf-8",
	}

	const enableCOOPCOEP = false // mude para true se Rust/WASM usar threads

	ext := strings.ToLower(filepath.Ext(path))
	if ctype, ok := mimeByExt[ext]; ok {
		w.Header().Set("Content-Type", ctype)
	}
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	if enableCOOPCOEP && ext == ".wasm" {
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
	}
}
