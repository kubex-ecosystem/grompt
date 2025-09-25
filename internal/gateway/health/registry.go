package health

import "sync"

// ProberRegistry gerencia todos os probers disponíveis
type ProberRegistry struct {
	probers map[string]Prober
	mu      sync.RWMutex
}

// NewProberRegistry cria um novo registry
func NewProberRegistry() *ProberRegistry {
	return &ProberRegistry{
		probers: make(map[string]Prober),
	}
}

// Register adiciona um prober ao registry
func (r *ProberRegistry) Register(prober Prober) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.probers[prober.Name()] = prober
}

// Get retorna um prober pelo nome
func (r *ProberRegistry) Get(name string) (Prober, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	prober, exists := r.probers[name]
	return prober, exists
}

// List retorna todos os probers registrados
func (r *ProberRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.probers))
	for name := range r.probers {
		names = append(names, name)
	}
	return names
}

// Remove remove um prober do registry
func (r *ProberRegistry) Remove(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.probers, name)
}

// DefaultProberRegistry é o registry global padrão
var DefaultProberRegistry = NewProberRegistry()

// RegisterDefaultProbers registra todos os probers padrão
func RegisterDefaultProbers() {
	DefaultProberRegistry.Register(NewGroqProber())
	DefaultProberRegistry.Register(NewGeminiProber())

	// TODO: Adicionar outros probers conforme implementados
	// DefaultProberRegistry.Register(NewOpenAIProber())
	// DefaultProberRegistry.Register(NewClaudeProber())
	// DefaultProberRegistry.Register(NewDeepSeekProber())
	// DefaultProberRegistry.Register(NewOllamaProber())
}

// GetProber é um helper para acessar probers do registry padrão
func GetProber(name string) (Prober, bool) {
	return DefaultProberRegistry.Get(name)
}

// ListProbers é um helper para listar probers do registry padrão
func ListProbers() []string {
	return DefaultProberRegistry.List()
}
