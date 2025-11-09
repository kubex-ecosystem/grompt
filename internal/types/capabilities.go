package types

import "github.com/kubex-ecosystem/grompt/internal/interfaces"

// LegacyCapabilities represents the capabilities of a legacy provider.
type LegacyCapabilities struct {
	MaxTokens         int                 `json:"max_tokens"`
	SupportsBatch     bool                `json:"supports_batch"`
	SupportsStreaming bool                `json:"supports_streaming"`
	SupportsImages    bool                `json:"supports_images"`
	SupportsAudio     bool                `json:"supports_audio"`
	Models            []string            `json:"models"`
	Pricing           *interfaces.Pricing `json:"pricing,omitempty"`
}
