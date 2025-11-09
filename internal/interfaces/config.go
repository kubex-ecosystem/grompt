package interfaces


type ServerConfig interface {
	GetAPIConfig(provider string) IAPIConfig
	GetPort() string
	GetAPIKey(provider string) string
	SetAPIKey(provider, key string) error
	GetAPIEndpoint(provider string) string
	GetBaseGenerationPrompt(ideas []string, purpose, purposeType, lang string, maxLength int) string
}

type IConfig interface {
	GetAPIConfig(provider string) IAPIConfig
	GetPort() string
	GetAPIKey(provider string) string
	SetAPIKey(provider string, key string) error
	GetAPIEndpoint(provider string) string
	GetBaseGenerationPrompt(ideas []string, purpose, purposeType, lang string, maxLength int) string
}
