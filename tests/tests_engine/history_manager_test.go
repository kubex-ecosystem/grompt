package engine_test

import (
	"testing"

	eng "github.com/rafa-mori/grompt/internal/core/engine"
	"github.com/rafa-mori/grompt/internal/types"
)

func TestHistoryManager_AddAndTrim(t *testing.T) {
	h := eng.NewHistoryManager(2)
	if h == nil {
		t.Fatalf("NewHistoryManager returned nil")
	}

	// Add 3 results, max 2 -> oldest should be trimmed
	h.Add(types.Result{ID: "1", Prompt: "p1", Response: "r1"})
	h.Add(types.Result{ID: "2", Prompt: "p2", Response: "r2"})
	h.Add(types.Result{ID: "3", Prompt: "p3", Response: "r3"})

	got := h.GetHistory()
	if len(got) != 2 {
		t.Fatalf("history size = %d, want 2", len(got))
	}
	if got[0].ID != "2" || got[1].ID != "3" {
		t.Fatalf("expected IDs [2,3], got [%s,%s]", got[0].ID, got[1].ID)
	}
}
