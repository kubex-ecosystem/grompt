package server

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"strings"

// 	// "github.com/kubex-ecosystem/grompt/internal/services/squad"
// )

// // HandleAgents manages listing and creating agents.
// func (h *Handlers) HandleAgents(w http.ResponseWriter, r *http.Request) {
// 	h.setCORSHeaders(w)
// 	if r.Method == "OPTIONS" {
// 		return
// 	}

// 	switch r.Method {
// 	case http.MethodGet:
// 		agents := h.agentStore.All()
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(agents)
// 	case http.MethodPost:
// 		var a squad.Agent
// 		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
// 			http.Error(w, "invalid json", http.StatusBadRequest)
// 			return
// 		}
// 		sa, err := h.agentStore.Add(a)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(sa)
// 	default:
// 		h.HandleMethodNotAllowed(w, r)
// 	}
// }

// // HandleAgent manages operations on a single agent by ID.
// func (h *Handlers) HandleAgent(w http.ResponseWriter, r *http.Request) {
// 	h.setCORSHeaders(w)
// 	if r.Method == "OPTIONS" {
// 		return
// 	}
// 	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/agents/")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "invalid id", http.StatusBadRequest)
// 		return
// 	}

// 	switch r.Method {
// 	case http.MethodGet:
// 		for _, a := range h.agentStore.All() {
// 			if a.ID == id {
// 				w.Header().Set("Content-Type", "application/json")
// 				json.NewEncoder(w).Encode(a)
// 				return
// 			}
// 		}
// 		http.NotFound(w, r)
// 	case http.MethodPut:
// 		var ag squad.Agent
// 		if err := json.NewDecoder(r.Body).Decode(&ag); err != nil {
// 			http.Error(w, "invalid json", http.StatusBadRequest)
// 			return
// 		}
// 		sa, err := h.agentStore.Update(id, ag)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusNotFound)
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(sa)
// 	case http.MethodDelete:
// 		if err := h.agentStore.Delete(id); err != nil {
// 			http.Error(w, err.Error(), http.StatusNotFound)
// 			return
// 		}
// 		w.WriteHeader(http.StatusNoContent)
// 	default:
// 		h.HandleMethodNotAllowed(w, r)
// 	}
// }

// // HandleAgentsMarkdown returns current agents as AGENTS.md content.
// func (h *Handlers) HandleAgentsMarkdown(w http.ResponseWriter, r *http.Request) {
// 	h.setCORSHeaders(w)
// 	if r.Method == "OPTIONS" {
// 		return
// 	}
// 	if r.Method != http.MethodGet {
// 		h.HandleMethodNotAllowed(w, r)
// 		return
// 	}
// 	md := h.agentStore.ToMarkdown()
// 	w.Header().Set("Content-Type", "text/markdown")
// 	w.Write([]byte(md))
// }

// // HandleAgentsGenerate generates a squad of agents based on requirements using LLM
// func (h *Handlers) HandleAgentsGenerate(w http.ResponseWriter, r *http.Request) {
// 	h.setCORSHeaders(w)
// 	if r.Method == "OPTIONS" {
// 		return
// 	}

// 	if r.Method != http.MethodPost {
// 		h.HandleMethodNotAllowed(w, r)
// 		return
// 	}

// 	var req struct {
// 		Requirements string `json:"requirements"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "invalid json", http.StatusBadRequest)
// 		return
// 	}

// 	if req.Requirements == "" {
// 		http.Error(w, "requirements field is required", http.StatusBadRequest)
// 		return
// 	}

// 	// Create LLM function based on available APIs
// 	llmFunc := func(prompt string) (string, error) {
// 		// Try different providers in order of preference
// 		if h.config.GetAPIKey("claude") != "" {
// 			return h.claudeAPI.Complete(prompt, 4000, "claude-2")
// 		}
// 		if h.config.GetAPIKey("openai") != "" {
// 			return h.openaiAPI.Complete(prompt, 4000, "gpt-4")
// 		}
// 		if h.config.GetAPIKey("deepseek") != "" {
// 			return h.deepseekAPI.Complete(prompt, 4000, "deepseek-chat")
// 		}
// 		return "", fmt.Errorf("no LLM API available")
// 	}

// 	// Generate agents using LLM
// 	agents, err := squad.ParseRequirementsWithLLM(req.Requirements, llmFunc)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Generate the markdown
// 	markdown := squad.GenerateMarkdown(agents)

// 	response := struct {
// 		Agents   []squad.Agent `json:"agents"`
// 		Markdown string        `json:"markdown"`
// 	}{
// 		Agents:   agents,
// 		Markdown: markdown,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// // HandleAgentsImport handles importing AGENTS.md files with validation
// func (h *Handlers) HandleAgentsImport(w http.ResponseWriter, r *http.Request) {
// 	h.setCORSHeaders(w)
// 	if r.Method == "OPTIONS" {
// 		return
// 	}

// 	if r.Method != http.MethodPost {
// 		h.HandleMethodNotAllowed(w, r)
// 		return
// 	}

// 	var req struct {
// 		Content  string `json:"content"`
// 		Merge    bool   `json:"merge"`    // Whether to merge with existing agents
// 		Validate bool   `json:"validate"` // Whether to validate before importing
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "invalid json", http.StatusBadRequest)
// 		return
// 	}

// 	if req.Content == "" {
// 		http.Error(w, "content field is required", http.StatusBadRequest)
// 		return
// 	}

// 	// Parse the markdown content
// 	parseResult := squad.ParseAgentsMarkdown(req.Content)

// 	// If validation is requested and there are errors, return them
// 	if req.Validate && len(parseResult.Errors) > 0 {
// 		response := struct {
// 			Success bool               `json:"success"`
// 			Message string             `json:"message"`
// 			Agents  []squad.Agent      `json:"agents"`
// 			Errors  []squad.ParseError `json:"errors"`
// 		}{
// 			Success: false,
// 			Message: "Validation errors found in the imported content",
// 			Agents:  parseResult.Agents,
// 			Errors:  parseResult.Errors,
// 		}

// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	// If not merging, clear existing agents first
// 	if !req.Merge {
// 		// Get all existing agents and delete them
// 		existingAgents := h.agentStore.All()
// 		for _, agent := range existingAgents {
// 			h.agentStore.Delete(agent.ID)
// 		}
// 	}

// 	// Import the new agents
// 	var importedAgents []interface{}
// 	for _, agent := range parseResult.Agents {
// 		storedAgent, err := h.agentStore.Add(agent)
// 		if err != nil {
// 			// Log error but continue with other agents
// 			continue
// 		}
// 		importedAgents = append(importedAgents, storedAgent)
// 	}

// 	response := struct {
// 		Success        bool               `json:"success"`
// 		Message        string             `json:"message"`
// 		ImportedAgents []interface{}      `json:"imported_agents"`
// 		ImportedCount  int                `json:"imported_count"`
// 		Errors         []squad.ParseError `json:"errors"`
// 		TotalErrors    int                `json:"total_errors"`
// 	}{
// 		Success:        true,
// 		Message:        fmt.Sprintf("Successfully imported %d agents", len(importedAgents)),
// 		ImportedAgents: importedAgents,
// 		ImportedCount:  len(importedAgents),
// 		Errors:         parseResult.Errors,
// 		TotalErrors:    len(parseResult.Errors),
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }

// // HandleAgentsExportAdvanced handles advanced export with metadata options
// func (h *Handlers) HandleAgentsExportAdvanced(w http.ResponseWriter, r *http.Request) {
// 	h.setCORSHeaders(w)
// 	if r.Method == "OPTIONS" {
// 		return
// 	}

// 	if r.Method != http.MethodPost {
// 		h.HandleMethodNotAllowed(w, r)
// 		return
// 	}

// 	var req struct {
// 		IncludeMetadata bool   `json:"include_metadata"`
// 		AgentIDs        []int  `json:"agent_ids"` // Specific agents to export (empty = all)
// 		Format          string `json:"format"`    // "markdown", "json", "yaml"
// 		Filename        string `json:"filename"`  // Suggested filename
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "invalid json", http.StatusBadRequest)
// 		return
// 	}

// 	// Default values
// 	if req.Format == "" {
// 		req.Format = "markdown"
// 	}
// 	if req.Filename == "" {
// 		req.Filename = "AGENTS"
// 	}

// 	// Get agents to export
// 	allAgents := h.agentStore.All()
// 	var agentsToExport []squad.Agent

// 	if len(req.AgentIDs) > 0 {
// 		// Export specific agents
// 		agentMap := make(map[int]squad.Agent)
// 		for _, stored := range allAgents {
// 			agentMap[stored.ID] = stored.Agent
// 		}

// 		for _, id := range req.AgentIDs {
// 			if agent, exists := agentMap[id]; exists {
// 				agentsToExport = append(agentsToExport, agent)
// 			}
// 		}
// 	} else {
// 		// Export all agents
// 		for _, stored := range allAgents {
// 			agentsToExport = append(agentsToExport, stored.Agent)
// 		}
// 	}

// 	var content string
// 	var contentType string
// 	var extension string

// 	switch req.Format {
// 	case "json":
// 		data, err := json.MarshalIndent(agentsToExport, "", "  ")
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		content = string(data)
// 		contentType = "application/json"
// 		extension = ".json"

// 	case "yaml":
// 		// Simple YAML-like format for agents
// 		var b strings.Builder
// 		b.WriteString("agents:\n")
// 		for _, agent := range agentsToExport {
// 			b.WriteString(fmt.Sprintf("  - title: %q\n", agent.Title))
// 			if agent.Role != "" {
// 				b.WriteString(fmt.Sprintf("    role: %q\n", agent.Role))
// 			}
// 			if len(agent.Skills) > 0 {
// 				b.WriteString("    skills:\n")
// 				for _, skill := range agent.Skills {
// 					b.WriteString(fmt.Sprintf("      - %q\n", skill))
// 				}
// 			}
// 			if len(agent.Restrictions) > 0 {
// 				b.WriteString("    restrictions:\n")
// 				for _, restriction := range agent.Restrictions {
// 					b.WriteString(fmt.Sprintf("      - %q\n", restriction))
// 				}
// 			}
// 			if agent.PromptExample != "" {
// 				b.WriteString(fmt.Sprintf("    prompt_example: %q\n", agent.PromptExample))
// 			}
// 			b.WriteString("\n")
// 		}
// 		content = b.String()
// 		contentType = "application/x-yaml"
// 		extension = ".yaml"

// 	default: // markdown
// 		content = squad.ExportAgentsToMarkdown(agentsToExport, req.IncludeMetadata)
// 		contentType = "text/markdown"
// 		extension = ".md"
// 	}

// 	// Set headers for file download
// 	filename := req.Filename + extension
// 	w.Header().Set("Content-Type", contentType)
// 	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
// 	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))

// 	// Write content
// 	w.Write([]byte(content))
// }

// // HandleAgentsValidate validates agents content without importing
// func (h *Handlers) HandleAgentsValidate(w http.ResponseWriter, r *http.Request) {
// 	h.setCORSHeaders(w)
// 	if r.Method == "OPTIONS" {
// 		return
// 	}

// 	if r.Method != http.MethodPost {
// 		h.HandleMethodNotAllowed(w, r)
// 		return
// 	}

// 	var req struct {
// 		Content string `json:"content"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "invalid json", http.StatusBadRequest)
// 		return
// 	}

// 	if req.Content == "" {
// 		http.Error(w, "content field is required", http.StatusBadRequest)
// 		return
// 	}

// 	// Parse and validate the content
// 	parseResult := squad.ParseAgentsMarkdown(req.Content)

// 	// Additional validation for each parsed agent
// 	var allErrors []squad.ParseError
// 	allErrors = append(allErrors, parseResult.Errors...)

// 	for i, agent := range parseResult.Agents {
// 		validationErrors := squad.ValidateAgent(agent)
// 		for _, err := range validationErrors {
// 			err.Section = fmt.Sprintf("Agent %d (%s) - %s", i+1, agent.Title, err.Section)
// 			allErrors = append(allErrors, err)
// 		}
// 	}

// 	response := struct {
// 		Valid        bool               `json:"valid"`
// 		AgentsFound  int                `json:"agents_found"`
// 		Agents       []squad.Agent      `json:"agents"`
// 		Errors       []squad.ParseError `json:"errors"`
// 		ErrorCount   int                `json:"error_count"`
// 		WarningCount int                `json:"warning_count"`
// 	}{
// 		Valid:        len(allErrors) == 0,
// 		AgentsFound:  len(parseResult.Agents),
// 		Agents:       parseResult.Agents,
// 		Errors:       allErrors,
// 		ErrorCount:   len(allErrors),
// 		WarningCount: 0, // Could implement warnings separately
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }
