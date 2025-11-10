// Package main is the entry point of the application.
package main

import (
	"github.com/kubex-ecosystem/grompt/internal/module"
	gl "github.com/kubex-ecosystem/logz/logger"
)

func main() {
	if err := module.RegX().Command().Execute(); err != nil {
		gl.LoggerG.GetLogger().Log("fatalc", err.Error())
	}
}
