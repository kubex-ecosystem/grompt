// Package engine provides the core functionality for the factory engine.
package engine

import (
	"github.com/kubex-ecosystem/gemx/grompt/internal/engine"
	"github.com/kubex-ecosystem/gemx/grompt/internal/types"
)

type Engine = engine.IEngine

func NewEngine(config types.IConfig) engine.IEngine { return engine.NewEngine(config) }
