// Package config provides legacy constructors that bridge to the new grompt configuration layer.
package config

import (
    "github.com/kubex-ecosystem/grompt"
    logz "github.com/kubex-ecosystem/logz"
)

type Config = grompt.Config

type APIConfig = grompt.APIConfig

// NewConfig reconstructs a legacy configuration using the updated engine internals.
func NewConfig(
    bindAddr string,
    port string,
    openAIKey string,
    deepSeekKey string,
    ollamaEndpoint string,
    claudeKey string,
    geminiKey string,
    chatGPTKey string,
    logger logz.Logger,
) Config {
    return grompt.NewConfig(
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
    return grompt.DefaultConfig(path)
}
