// Package engine provides the core functionality for the factory engine.
package engine

import (
	"github.com/kubex-ecosystem/grompt/internal/engine"
	"github.com/kubex-ecosystem/grompt/internal/interfaces"
)

type Engine = interfaces.IEngine

func NewEngine(config interfaces.IConfig) interfaces.IEngine { return engine.NewEngine(config) }
