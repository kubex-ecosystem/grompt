package interfaces

import (
	"context"
)

// LegacyProvider is the legacy provider interface exposed to consumers.
type LegacyProvider interface {
	Name() string
	Version() string
	Execute(ctx context.Context, prompt string) (string, error)
	IsAvailable() bool
	KeyEnv() string
	Type() string
	GetCapabilities(ctx context.Context) *Capabilities
	Chat(ctx context.Context, req ChatRequest) (<-chan ChatChunk, error)
	Notify(ctx context.Context, event NotificationEvent) error
}
