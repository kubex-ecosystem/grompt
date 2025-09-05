// Package engine provides the core functionality for the factory engine.
package engine

import (
	"github.com/rafa-mori/grompt/internal/core/engine"
	"github.com/rafa-mori/grompt/internal/core/provider"
)

type Engine = engine.IEngine

func NewEngine(config provider.IConfig) engine.IEngine { return engine.NewEngine(config) }
