// Package types defines common types used throughout the application.
package types

import "time"

// Result represents the outcome of a processed prompt.
type Result struct {
	ID        string         `json:"id"`
	Prompt    string         `json:"prompt"`
	Response  string         `json:"response"`
	Provider  string         `json:"provider"`
	Variables map[string]any `json:"variables,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}
