// Package engine provides the core functionality for the factory engine.
package engine

import (
	"github.com/kubex-ecosystem/grompt/internal/gateway/registry"
)

type Engine = registry.Registry

//func NewEngine(config types.IConfig) engine.IEngine { return engine.NewEngine(config) }
