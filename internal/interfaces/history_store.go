package interfaces

import "sync"

// ---------- History store ----------

type historyStore struct {
	mu      sync.RWMutex
	entries []Result
	limit   int
}

func NewHistoryStore(limit int) IHistoryManager {
	if limit <= 0 {
		limit = 100
	}
	return &historyStore{entries: make([]Result, 0, limit), limit: limit}
}

func (h *historyStore) Add(result Result) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.entries = append(h.entries, result)
	if len(h.entries) > h.limit {
		h.entries = h.entries[1:]
	}
}

func (h *historyStore) Snapshot() []Result {
	h.mu.RLock()
	defer h.mu.RUnlock()

	out := make([]Result, len(h.entries))
	copy(out, h.entries)
	return out
}

type IHistoryManager interface {
	// Add adds a new result to the history
	Add(result Result)

	// Snapshot returns a copy of the current history entries
	Snapshot() []Result
}
