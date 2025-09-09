// Package agents provides functionality to manage and persist agents.
// It allows adding, updating, deleting, and retrieving agents,
// as well as converting them to a markdown format for documentation.
package agents

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/kubex-ecosystem/gemx/grompt/internal/services/squad"
)

// StoredAgent wraps squad.Agent with an ID for persistence
// and JSON encoding.
type StoredAgent struct {
	ID int `json:"id"`
	squad.Agent
}

// Store manages a slice of StoredAgent and persists to a JSON file.
type Store struct {
	mu     sync.Mutex
	agents []StoredAgent
	path   string
	nextID int
}

// NewStore creates a new Store and loads data from the given file if it exists.
func NewStore(path string) *Store {
	s := &Store{path: path}
	s.load()
	return s
}

// load reads agents from the store's file.
func (s *Store) load() {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return
	}
	json.Unmarshal(data, &s.agents)
	for _, a := range s.agents {
		if a.ID >= s.nextID {
			s.nextID = a.ID + 1
		}
	}
}

// save writes agents to the store's file.
func (s *Store) save() error {
	data, err := json.MarshalIndent(s.agents, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// All returns all stored agents.
func (s *Store) All() []StoredAgent {
	s.mu.Lock()
	defer s.mu.Unlock()
	res := make([]StoredAgent, len(s.agents))
	copy(res, s.agents)
	return res
}

// Add inserts a new agent into the store.
func (s *Store) Add(a squad.Agent) (StoredAgent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sa := StoredAgent{ID: s.nextID, Agent: a}
	s.nextID++
	s.agents = append(s.agents, sa)
	return sa, s.save()
}

// Update replaces an existing agent by ID.
func (s *Store) Update(id int, a squad.Agent) (StoredAgent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, ag := range s.agents {
		if ag.ID == id {
			s.agents[i].Agent = a
			err := s.save()
			return s.agents[i], err
		}
	}
	return StoredAgent{}, os.ErrNotExist
}

// Delete removes an agent by ID.
func (s *Store) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, ag := range s.agents {
		if ag.ID == id {
			s.agents = append(s.agents[:i], s.agents[i+1:]...)
			return s.save()
		}
	}
	return os.ErrNotExist
}

// ToMarkdown converts stored agents to AGENTS.md format.
func (s *Store) ToMarkdown() string {
	agents := make([]squad.Agent, len(s.agents))
	for i, ag := range s.agents {
		agents[i] = ag.Agent
	}
	return squad.GenerateMarkdown(agents)
}
