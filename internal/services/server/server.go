// Package server implements the HTTP server for the Prompt Crafter application.
package server

import (
	"fmt"
	"io/fs"
	"log"
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
	gl "github.com/kubex-ecosystem/logz"
)

var reactApp = grompt.NewGUIGrompt()

type Server struct {
	*gateway.ServerImpl
	apiRouter *gin.RouterGroup
	config    *t.Config
	handlers  *Handlers
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
		ServerImpl: server,
		config:     baseCfg,
		handlers:   NewHandlers(cfg),
	}
}

func (s *Server) Start() error {
	// Configurar roteamento
	s.setupRoutes()

	gl.Infof("Gateway started: http://localhost:%s", s.config.Port)
	gl.Info("📁 Serving embedded React application")
	gl.Info("🔧 Available APIs:")
	gl.Info("   • /api/v1/config - Configuration")
	gl.Info("   • /api/v1/models - Available Models")
	gl.Info("   • /api/v1/test - API Test")
	gl.Info("   • /api/v1/unified - Unified API")
	gl.Info("   • /api/v1/openai - OpenAI API")
	gl.Info("   • /api/v1/deepseek - DeepSeek API")
	gl.Info("   • /api/v1/claude - Claude API")
	gl.Info("   • /api/v1/gemini - Gemini API")
	gl.Info("   • /api/v1/chatgpt - ChatGPT API")
	gl.Info("   • /api/v1/ollama - Ollama Local")
	gl.Info("   • /api/v1/health - Server Status")
	gl.Info("💡 Press Ctrl+C to stop")

	// Detecta se há a página aberta em algum lugar

	// Abrir navegador após delay
	go func(url string) {
		time.Sleep(1 * time.Second)
		openBrowser(url)
	}("http://localhost:" + s.config.Port)

	return http.ListenAndServe(
		net.JoinHostPort(s.config.BindAddr, s.config.Port),
		s.ServerImpl.Router(),
	)
}

func (s *Server) setupRoutes() {
	// 1) Rotas da API
	s.setupAPIRoutes()

	// 2) Rotas do Gateway (Providers dinâmicos)
	s.setupGatewayRoutes()

	// 3) Rotas estáticas (React App)
	s.setupStaticRoutes()

	// 4) Montar roteador da API (separado)
	s.ServerImpl.RegisterRoutes()
}

func (s *Server) setupStaticRoutes() {
	buildFS, err := fs.Sub(reactApp.GetWebFS(), "embedded/guiweb")
	if err != nil {
		gl.Log("warn", fmt.Sprintf("⚠️ build embed não encontrado: %v", err))
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
	apiRouter := gin.IRouter(s.ServerImpl.GetRouter().Group("/api/v1"))

	// Mapeamento de rotas para registro individual (se necessário)
	routeMap := make(map[string]struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	})

	// 1) Rotas básicas
	routeMap["api-v1-health-Get"] = struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	}{"/api/v1/health", http.MethodGet, s.handlers.HandleHealth}
	routeMap["api-v1-config-Get"] = struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	}{"/api/v1/config", http.MethodGet, s.handlers.HandleConfig}
	routeMap["api-v1-config-Post"] = struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	}{"/api/v1/config", http.MethodPost, s.handlers.HandleConfig}
	routeMap["api-v1-test-Get"] = struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	}{"/api/v1/test", http.MethodGet, s.handlers.HandleTest}
	routeMap["api-v1-models-Get"] = struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	}{"/api/v1/models", http.MethodGet, s.handlers.HandleModels}

	// 2) Provedores (diretos)
	routeMap["api-v1-openai-Post"] = struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	}{"/api/v1/openai", http.MethodPost, s.handlers.HandleOpenAI}
	routeMap["api-v1-claude-Post"] = struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	}{"/api/v1/claude", http.MethodPost, s.handlers.HandleClaude}
	routeMap["api-v1-gemini-Post"] = struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	}{"/api/v1/gemini", http.MethodPost, s.handlers.HandleGemini}
	routeMap["api-v1-deepseek-Post"] = struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	}{"/api/v1/deepseek", http.MethodPost, s.handlers.HandleDeepSeek}
	routeMap["api-v1-chatgpt-Post"] = struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	}{"/api/v1/chatgpt", http.MethodPost, s.handlers.HandleChatGPT}
	routeMap["api-v1-ollama-Post"] = struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	}{"/api/v1/ollama", http.MethodPost, s.handlers.HandleOllama}

	// 3) Geração Unificada e Atalhos
	routeMap["api-v1-unified-Post"] = struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	}{"/api/v1/unified", http.MethodPost, s.handlers.HandleUnified}
	routeMap["api-v1-ask-Post"] = struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	}{"/api/v1/ask", http.MethodPost, s.handlers.HandleAsk}
	routeMap["api-v1-squad-Post"] = struct {
		Path    string
		Method  string
		Handler gin.HandlerFunc
	}{"/api/v1/squad", http.MethodPost, s.handlers.HandleSquad}

	// 4) Agentes / Squad
	// s.GET("/api/v1/agents", getGinHandlerFunc(s.handlers.HandleAgents))
	// s.POST("/api/v1/agents", getGinHandlerFunc(s.handlers.HandleAgents))
	// s.POST("/api/v1/agents/generate", getGinHandlerFunc(s.handlers.HandleAgentsGenerate))
	// s.GET("/api/v1/agents/", getGinHandlerFunc(s.handlers.HandleAgent))
	// s.PUT("/api/v1/agents/", getGinHandlerFunc(s.handlers.HandleAgent))
	// s.DELETE("/api/v1/agents/", getGinHandlerFunc(s.handlers.HandleAgent))
	// s.GET("/api/v1/agents.md", getGinHandlerFunc(s.handlers.HandleAgentsMarkdown))

	// 5) Registrar rotas no roteador da API
	for _, route := range routeMap {
		switch route.Method {
		case http.MethodGet:
			apiRouter.GET(route.Path, route.Handler)
		case http.MethodPost:
			apiRouter.POST(route.Path, route.Handler)
		case http.MethodPut:
			apiRouter.PUT(route.Path, route.Handler)
		case http.MethodDelete:
			apiRouter.DELETE(route.Path, route.Handler)
		}
	}

	// 6) Rotas de teste / WASM
	s.ServerImpl.Router().GET("/wasm-test", gin.WrapH(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	const result = parse("LookAtni is awesome!");
	console.log('WASM parse result:', result);
});
</script>
</body>
</html>`))
	})))

	registry, err := registry.FromRuntimeConfig(s.config)
	if err != nil {
		gl.Log("warn", "⚠️  API runtime desabilitado: %v", err)
		return
	}
	gatewayRouter := routes.NewGatewayRoutes(registry, middleware.NewProductionMiddleware(middleware.DefaultProductionConfig()))

	gatewayRouter.Register(apiRouter)

	s.apiRouter = apiRouter.(*gin.RouterGroup)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", "", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}

	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}

func (s *Server) setupGatewayRoutes() {
	registry, err := registry.FromRuntimeConfig(s.config)
	if err != nil {
		gl.Log("warn", "⚠️  Gateway runtime desabilitado: %v", err)
		return
	}
	gatewayRouter := routes.NewGatewayRoutes(registry, middleware.NewProductionMiddleware(middleware.DefaultProductionConfig()))
	gatewayRouter.Register(s.ServerImpl.GetRouter())
}

func (s *Server) setupFallbackRoutes() error {

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

	// Config route (fallback path)
	// This route is an alternative path to access the server's configuration.
	// It serves the same purpose as the /api/v1/config route.
	s.apiRouter.GET("/api/config", s.handlers.HandleConfig)

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
