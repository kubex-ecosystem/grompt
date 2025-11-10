package interfaces

import (
	"github.com/gin-gonic/gin"
	"github.com/kubex-ecosystem/grompt/internal/gateway/middleware"
)

type GatewayRoutes interface {
	Register(router gin.IRouter)
	GetRegistry() Registry
	GetMiddleware() *middleware.ProductionMiddleware
}

// Server represents the gateway server
type Server interface {
	// Start starts the server
	Start() error
	// Stop stops the server
	Stop() error
	// Router returns the underlying HTTP router
	Router() *gin.Engine
	// GetRouter returns the underlying HTTP router
	GetRouter() *gin.Engine
	// RegisterRoutes registers multiple gateway routes
	RegisterRoutes(routes ...GatewayRoutes) error
	// RegisterRoutesByHandlers registers routes based on provided handlers
	RegisterRoutesByHandlers(handlers Handlers) error
}
