// Package interfaces defines the contracts for various components in the application.
package interfaces

// LegacyAPIConfig defines the contract for legacy API configurations
type LegacyAPIConfig interface {
	IsAvailable() bool
	IsDemoMode() bool
	Version() string
	ListModels() ([]string, error)
	GetCommonModels() []string
	Complete(prompt string, maxTokens int, model string) (string, error)
}

type IAPIConfig interface {
	IsAvailable() bool
	IsDemoMode() bool
	Version() string
	ListModels() (map[string]any, error)
	GetCommonModels() []string
	Complete(prompt string, maxTokens int, model string) (string, error)
}
