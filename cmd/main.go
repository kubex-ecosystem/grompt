// Package main is the entry point of the application.
package main

import (
	"github.com/kubex-ecosystem/grompt/internal/module"
	gl "github.com/kubex-ecosystem/grompt/internal/module/logger"
)

func main() {
	if err := module.RegX().Command().Execute(); err != nil {
		gl.Log("fatal", err.Error())
	}
}
