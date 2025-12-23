// Package main is the entry point of the application.
package main

import (
	"github.com/kubex-ecosystem/grompt/internal/module"
	gl "github.com/kubex-ecosystem/logz"
)

func main() {
	gl.GetLoggerZ("Grompt")
	if err := module.RegX().Command().Execute(); err != nil {
		gl.Log("fatalc", err.Error())
	}
}
