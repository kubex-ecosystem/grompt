// Package gateway provides the gateway server functionality for the grompt.
package gateway

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/kubex-ecosystem/grompt/internal/gateway/middleware"
	"github.com/kubex-ecosystem/grompt/internal/gateway/registry"
	"github.com/kubex-ecosystem/grompt/internal/gateway/routes"
	"github.com/kubex-ecosystem/grompt/internal/interfaces"
)

// ServerConfig holds configuration for the gateway server
type ServerConfig struct {
	Addr            string
	ProvidersConfig string
	Debug           bool
	EnableCORS      bool
}

// Server represents the gateway server
type Server struct {
	config     interfaces.IConfig
	registry   *registry.Registry
	middleware *middleware.ProductionMiddleware
	router     *gin.Engine
	routes     *routes.GatewayRoutes
}

// NewServer creates a new gateway server instance
func NewServer(config interfaces.IConfig) (*Server, error) {
	// Load providers registry
	reg, err := registry.Load(config.GetConfigFilePath())
	if err != nil {
		return nil, err
	}

	if !config.IsDebugMode() {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Initialize production middleware
	prodConfig := middleware.DefaultProductionConfig()
	prodMiddleware := middleware.NewProductionMiddleware(prodConfig)

	// Register all providers with production middleware
	for _, providerName := range reg.ListProviders() {
		prodMiddleware.RegisterProvider(providerName)
	}

	return &Server{
		config:     config,
		registry:   reg,
		middleware: prodMiddleware,
		router:     router,
		routes:     routes.NewGatewayRoutes(reg, prodMiddleware),
	}, nil
}

// Start starts the gateway server
func (s *Server) Start() error {
	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("🛑 Shutting down gracefully...")
		s.middleware.Stop()
		os.Exit(0)
	}()

	if s.config.IsCORSEnabled() {
		s.router.Use(corsMiddleware())
	}

	s.routes.Register(s.router)

	log.Printf("🚀 grompt-gw listening on %s with ENTERPRISE features!",
		net.JoinHostPort(
			s.config.GetConfigArgs().Bind,
			s.config.GetConfigArgs().Port,
		),
	)
	return s.router.Run(
		net.JoinHostPort(
			s.config.GetConfigArgs().Bind,
			s.config.GetConfigArgs().Port,
		),
	)
}

// corsMiddleware reproduz a lógica anterior usando o pipeline do gin.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "content-type, authorization, x-external-api-key, x-tenant-id, x-user-id")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
