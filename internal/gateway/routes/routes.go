// Package routes defines and registers HTTP routes for the Grompt gateway.
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kubex-ecosystem/grompt/internal/gateway/middleware"
	"github.com/kubex-ecosystem/grompt/internal/gateway/registry"
	"github.com/kubex-ecosystem/grompt/internal/gateway/transport"
)

// GatewayRoutes centraliza o registro das rotas HTTP do gateway.
// A ideia é manter um ponto único para evoluir, versionar ou adicionar
// novos grupos de rotas sem precisar tocar diretamente no servidor.
type GatewayRoutes struct {
	registry   *registry.Registry
	middleware *middleware.ProductionMiddleware
}

// NewGatewayRoutes cria um registrador de rotas para o gateway.
func NewGatewayRoutes(reg *registry.Registry, mw *middleware.ProductionMiddleware) *GatewayRoutes {
	return &GatewayRoutes{registry: reg, middleware: mw}
}

// Register injeta todas as rotas conhecidas no router informado.
func (gr *GatewayRoutes) Register(router gin.IRouter) {
	transport.WireHTTPSSE(router, gr.registry, gr.middleware)
}
