package main

import (
	"github.com/rafa-mori/grompt/internal/module"
	gl "github.com/rafa-mori/grompt/internal/module/logger"
)

// This file is the entry point for the Grompt CLI application.
// It initializes the logger and starts the application by executing the main command.
// It allows the application to be run as a standalone CLI tool.

// main initializes the logger and creates a new GoBE instance.
func main() {
	if err := module.RegX().Command().Execute(); err != nil {
		gl.Log("fatal", err.Error())
	}
}
