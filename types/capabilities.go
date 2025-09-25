package types

// Capabilities represents the capabilities of a provider.
type Capabilities struct {
	MaxTokens         int      `json:"max_tokens"`
	SupportsBatch     bool     `json:"supports_batch"`
	SupportsStreaming bool     `json:"supports_streaming"`
	SupportsImages    bool     `json:"supports_images"`
	SupportsAudio     bool     `json:"supports_audio"`
	Models            []string `json:"models"`
	Pricing           *Pricing `json:"pricing,omitempty"`
}
