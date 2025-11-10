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
	"strconv"
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

type ReactApp struct {
	*gin.Engine

	FS          *grompt.GUIGrompt
	Wasms       *[]fs.File
	ReactRoutes map[string]string
	WasmRoutes  map[string]string
}

type Server struct {
	*gateway.ServerImpl
	config   *t.Config
	handlers *Handlers
	reactApp *ReactApp
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

	reactApp := &ReactApp{
		FS:    	grompt.NewGUIGrompt(),
		Engine: gin.New(),
	}
	reactApp.Engine.Use(gin.Recovery())

	return &Server{
		ServerImpl:   server,
		config:    baseCfg,
		reactApp:  reactApp,
		handlers:  NewHandlers(cfg),
	}
}

func (s *Server) Start() error {
	// Configurar roteamento
	s.setupRoutes()
	var url string

	// Listen and server frontend routes in a goroutine
	go func(cs *Server) {
		cs.setupStaticRoutes()
		intPort, err := strconv.Atoi(cs.config.Port)
		if err != nil {
			gl.Log("error", "❌ Invalid port number: %v", err)
			return
		}
		cs.reactApp.Engine.MaxMultipartMemory = 8 << 20 // 8 MiB
		intPort++
		url = fmt.Sprintf("http://localhost:%s", strconv.Itoa(intPort))
		if err := cs.reactApp.Engine.Run(
			net.JoinHostPort(cs.config.BindAddr, strconv.Itoa(intPort)),
		); err != nil && err != http.ErrServerClosed {
			gl.Log("error", "❌ Failed to start React frontend server: %v", err)
		}
	}(s)

	// Abrir navegador após delay
	go func() {
		<-time.After(1 * time.Second)

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

		openBrowser(url)
	}()

	return http.ListenAndServe(
		net.JoinHostPort(s.config.BindAddr, s.config.Port),
		s.ServerImpl.Router(),
	)
}

func (s *Server) setupRoutes() {
	s.setupAPIRoutes()
	s.setupGatewayRoutes()
	s.mountAPIRouter()
}

func (s *Server) setupStaticRoutes() {
	// Rota para servir arquivos estáticos do React
	s.reactApp.Engine.GET(
		"/*filepath",
		func(c *gin.Context) {
			path := filepath.Clean(c.Request.URL.Path)

			if strings.HasPrefix(path, "/api/v1/") {
				c.Next()
				return
			}
			if path == "/" {
				path = "/index.html"
			}
			fullPath := strings.TrimSuffix(path, "/")
			file, err := s.reactApp.FS.GetEmbeddedFS().F.Open(fullPath)
			if err != nil {
				// Arquivo não encontrado, retornar 404
				c.Next()
				return
			}
			defer file.Close()

			SetStaticHeaders(c.Writer, path)
			content, err := s.reactApp.FS.GetWebFile(fullPath)
			if err != nil {
				c.Next()
				return
			}
			c.Writer.Write(content)
		},
	)


}

func (s *Server) setupAPIRoutes() {

	// 1) Rotas básicas
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
	s.ServerImpl.Router().Group("/api/v1", s.handlers.setCORSHeaders)
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

	routes.NewGatewayRoutes(reg, prodMiddleware).Register(s.ServerImpl.Router().Group(""))
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
