// Package main is the entry point of the application.
package main

import (
	"github.com/kubex-ecosystem/grompt/internal/module"
	gl "github.com/kubex-ecosystem/logz/logger"
	l "github.com/kubex-ecosystem/logz"
)

func main() {
	l.GetLogger("Grompt")
	if err := module.RegX().Command().Execute(); err != nil {
		gl.Log("fatalc", err.Error())
	}
}
