// Package templates defines interfaces for template management.
package templates

import "fmt"

// Manager handles prompt templates
type Manager interface {
	// Process processes a template with variables
	Process(template string, vars map[string]interface{}) (string, error)

	// LoadTemplate loads a template by name
	LoadTemplate(name string) (string, error)

	// SaveTemplate saves a template with a name
	SaveTemplate(name, content string) error

	// ListTemplates returns all available template names
	ListTemplates() []string

	// DeleteTemplate removes a template
	DeleteTemplate(name string) error
}

// Template represents a prompt template
type Template struct {
	Name        string                 `json:"name"`
	Content     string                 `json:"content"`
	Variables   []string               `json:"variables"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// NewManager creates a new template manager
func NewManager(templatesPath string) Manager {
	// This will be implemented in internal/templates with concrete implementation
	// For now, return a simple implementation to avoid compilation errors
	return &simpleManager{path: templatesPath}
}

// simpleManager is a basic implementation to prevent compilation errors
type simpleManager struct {
	path string
}

func (sm *simpleManager) Process(template string, vars map[string]interface{}) (string, error) {
	// Simple template processing - replace {{.key}} with values
	result := template
	for key, value := range vars {
		placeholder := fmt.Sprintf("{{.%s}}", key)
		valueStr := fmt.Sprintf("%v", value)
		// Simple string replacement for now
		result = template + " [" + placeholder + "=" + valueStr + "]"
	}
	return result, nil
}

func (sm *simpleManager) LoadTemplate(name string) (string, error) {
	return "", fmt.Errorf("not implemented yet")
}

func (sm *simpleManager) SaveTemplate(name, content string) error {
	return fmt.Errorf("not implemented yet")
}

func (sm *simpleManager) ListTemplates() []string {
	return []string{}
}

func (sm *simpleManager) DeleteTemplate(name string) error {
	return fmt.Errorf("not implemented yet")
}
