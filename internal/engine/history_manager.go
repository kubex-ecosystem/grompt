package engine

// IHistoryManager defines the interface for history management
type IHistoryManager interface {
	// Add adds a result to history
	Add(result Result)
	// GetHistory returns the prompt history
	GetHistory() []Result
}

// HistoryManager manages prompt history and versioning
type HistoryManager struct {
	entries    []Result
	maxEntries int
}

// NewHistoryManager creates a new history manager
func newHistoryManager(maxEntries int) *HistoryManager {
	return &HistoryManager{
		entries:    make([]Result, 0),
		maxEntries: maxEntries,
	}
}

// NewHistoryManager creates a new history manager
func NewHistoryManager(maxEntries int) IHistoryManager {
	return newHistoryManager(maxEntries)
}

// Add adds a result to history
func (h *HistoryManager) Add(result Result) {
	if h == nil {
		return
	}

	h.entries = append(h.entries, result)

	// Maintain max entries limit
	if len(h.entries) > h.maxEntries {
		h.entries = h.entries[1:] // Remove oldest entry
	}
}

// GetHistory returns the prompt history
func (h *HistoryManager) GetHistory() []Result {
	if h == nil {
		return nil
	}
	return h.entries
}
