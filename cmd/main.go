package main

import (
	gl "github.com/kubex-ecosystem/logz/logger"

	"github.com/kubex-ecosystem/grompt/internal/module"
)

func main() {
	if err := module.RegX().Command().Execute(); err != nil {
		gl.Log("fatal", err.Error())
	}
}
