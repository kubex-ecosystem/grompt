package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/rafa-mori/grompt/internal/services/squad"
)

// HandleAgents manages listing and creating agents.
func (h *Handlers) HandleAgents(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case http.MethodGet:
		agents := h.agentStore.All()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(agents)
	case http.MethodPost:
		var a squad.Agent
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		sa, err := h.agentStore.Add(a)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sa)
	default:
		h.HandleMethodNotAllowed(w, r)
	}
}

// HandleAgent manages operations on a single agent by ID.
func (h *Handlers) HandleAgent(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/api/agents/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		for _, a := range h.agentStore.All() {
			if a.ID == id {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(a)
				return
			}
		}
		http.NotFound(w, r)
	case http.MethodPut:
		var ag squad.Agent
		if err := json.NewDecoder(r.Body).Decode(&ag); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		sa, err := h.agentStore.Update(id, ag)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sa)
	case http.MethodDelete:
		if err := h.agentStore.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		h.HandleMethodNotAllowed(w, r)
	}
}

// HandleAgentsMarkdown returns current agents as AGENTS.md content.
func (h *Handlers) HandleAgentsMarkdown(w http.ResponseWriter, r *http.Request) {
	h.setCORSHeaders(w)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != http.MethodGet {
		h.HandleMethodNotAllowed(w, r)
		return
	}
	md := h.agentStore.ToMarkdown()
	w.Header().Set("Content-Type", "text/markdown")
	w.Write([]byte(md))
}
