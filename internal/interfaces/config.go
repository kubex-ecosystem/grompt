package interfaces

import "github.com/kubex-ecosystem/grompt/internal/module/kbx"


type IConfig interface {
	GetAPIConfig(provider string) IAPIConfig
	GetPort() string
	GetAPIKey(provider string) string
	SetAPIKey(provider, key string) error
	GetAPIEndpoint(provider string) string
	GetBaseGenerationPrompt(ideas []string, purpose, purposeType, lang string, maxLength int) string
	GetServerConfig() IConfig
	GetProviders() map[string]Provider
	GetConfigFilePath() string
	IsCORSEnabled() bool
	IsDebugMode() bool
	GetConfigArgs() kbx.InitArgs
	Validate() error
}
