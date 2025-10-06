// Package config provides legacy constructors that bridge to the new grompt configuration layer.
package config

import (
	grt "github.com/kubex-ecosystem/grompt"
	tp "github.com/kubex-ecosystem/grompt/internal/types"
	logz "github.com/kubex-ecosystem/logz"
)

type Config = grt.Config

type APIConfig = grt.APIConfig

// NewConfig reconstructs a legacy configuration using the updated engine internals.
func NewConfig(
	port string,
	bindAddr string,
	openAIKey string,
	deepSeekKey string,
	ollamaEndpoint string,
	claudeKey string,
	geminiKey string,
	chatGPTKey string,
	logger logz.Logger,
) Config {
	return tp.NewConfig(
		bindAddr,
		port,
		openAIKey,
		deepSeekKey,
		ollamaEndpoint,
		claudeKey,
		geminiKey,
		chatGPTKey,
		logger,
	)
}

// NewConfigFromFile loads a configuration file and augments it with environment overrides.
func NewConfigFromFile(path string) Config {
	return grt.DefaultConfig(path)
}
