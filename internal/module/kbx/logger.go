package kbx

import (
	gl "github.com/kubex-ecosystem/grompt/internal/module/logger"
	l "github.com/kubex-ecosystem/logz"
)

func SetDebugMode(debug bool) {
	gl.SetDebug(debug)
}

func SetLogLevel(level string) {
	gl.Logger.SetLogLevel(level)
}

func SetLogTrace(enable bool) {
	gl.Logger.SetShowTrace(enable)
}

func GetLogger(name string) Logger {
	return gl.Logger
}

func SetLogger(logger l.Logger) {
	// gl.SetLogger(logger)
	// TODO: Implement this function properly
}
