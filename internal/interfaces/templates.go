package interfaces

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
