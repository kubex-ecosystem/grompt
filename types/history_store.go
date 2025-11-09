package types

import (
	"sync"

	"github.com/kubex-ecosystem/grompt/internal/interfaces"
)

// ---------- History store ----------

type historyStore struct {
	mu      sync.RWMutex
	entries []interfaces.Result
	limit   int
}

func newHistoryStore(limit int) *historyStore {
	if limit <= 0 {
		limit = 100
	}
	return &historyStore{entries: make([]interfaces.Result, 0, limit), limit: limit}
}

func (h *historyStore) add(result interfaces.Result) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.entries = append(h.entries, result)
	if len(h.entries) > h.limit {
		h.entries = h.entries[1:]
	}
}

func (h *historyStore) snapshot() []interfaces.Result {
	h.mu.RLock()
	defer h.mu.RUnlock()

	out := make([]interfaces.Result, len(h.entries))
	copy(out, h.entries)
	return out
}
