package interfaces

type Handlers struct {
	Config      IConfig
	// Gateway     Server
	ClaudeAPI   IAPIConfig
	OpenaiAPI   IAPIConfig
	DeepseekAPI IAPIConfig
	ChatGPTAPI  IAPIConfig
	GeminiAPI   IAPIConfig
	OllamaAPI   IAPIConfig
}
