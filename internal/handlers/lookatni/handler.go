// Package lookatni provides HTTP handlers for lookatni integration
package lookatni

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kubex-ecosystem/grompt/internal/services/lookatni"
)

// Handler provides HTTP endpoints for lookatni functionality
type Handler struct {
	service *lookatni.LookAtniService
}

// NewHandler creates a new lookatni HTTP handler
func NewHandler(workDir string) *Handler {
	return &Handler{
		service: lookatni.NewLookAtniService(workDir),
	}
}

// HandleExtractProject handles POST /api/v1/lookatni/extract
func (h *Handler) HandleExtractProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req lookatni.ProjectExtractionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// Set defaults
	if len(req.IncludePatterns) == 0 {
		req.IncludePatterns = []string{"*.js", "*.ts", "*.jsx", "*.tsx", "*.go", "*.py", "*.java", "*.cpp", "*.c", "*.h", "*.php", "*.rb", "*.rs", "*.swift", "*.kt", "*.scala", "*.cs", "*.vb", "*.fs", "*.clj", "*.elm", "*.dart", "*.lua", "*.pl", "*.sh", "*.bat", "*.ps1", "*.sql", "*.html", "*.css", "*.scss", "*.less", "*.sass", "*.json", "*.xml", "*.yaml", "*.yml", "*.toml", "*.ini", "*.cfg", "*.conf", "*.md", "*.txt", "*.rst", "*.adoc"}
	}
	if len(req.ExcludePatterns) == 0 {
		req.ExcludePatterns = []string{"node_modules/**", ".git/**", "dist/**", "build/**", "target/**", "*.min.js", "*.bundle.js", "vendor/**", ".vscode/**", ".idea/**", "*.log", "*.tmp", "*.cache"}
	}
	if req.MaxFileSize == 0 {
		req.MaxFileSize = 1024 * 1024 // 1MB default
	}
	if req.ContextDepth == 0 {
		req.ContextDepth = 3
	}
	if req.FragmentBy == "" {
		req.FragmentBy = "function"
	}

	extracted, err := h.service.ExtractProject(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Extraction failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(extracted)
}

// HandleCreateArchive handles POST /api/v1/lookatni/archive
func (h *Handler) HandleCreateArchive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var extracted lookatni.ExtractedProject
	if err := json.NewDecoder(r.Body).Decode(&extracted); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	archivePath, err := h.service.CreateNavigableArchive(r.Context(), &extracted)
	if err != nil {
		http.Error(w, fmt.Sprintf("Archive creation failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate download URL
	fileName := filepath.Base(archivePath)
	downloadURL := fmt.Sprintf("/api/v1/lookatni/download/%s", fileName)

	response := map[string]interface{}{
		"success":      true,
		"archive_path": archivePath,
		"download_url": downloadURL,
		"file_name":    fileName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleDownloadArchive handles GET /api/v1/lookatni/download/{filename}
func (h *Handler) HandleDownloadArchive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract filename from URL path
	path := r.URL.Path
	filename := filepath.Base(path)

	// Security: validate filename
	if filename == "" || filename == "." || filename == ".." {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	// Construct full path
	downloadDir := filepath.Join(h.service.GetWorkDir(), "downloads")
	fullPath := filepath.Join(downloadDir, filename)

	// Security: ensure file is within downloads directory
	if !filepath.HasPrefix(fullPath, downloadDir) {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set headers for download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/zip")

	// Serve file
	http.ServeFile(w, r, fullPath)
}

// HandleListExtractedProjects handles GET /api/v1/lookatni/projects
func (h *Handler) HandleListExtractedProjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Query parameters
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// For now, return mock data - in a real implementation, this would query a database
	projects := []map[string]interface{}{
		{
			"id":           "proj_001",
			"name":         "grompt",
			"repo_url":     "https://github.com/kubex-ecosystem/grompt",
			"extracted_at": "2025-09-15T10:30:00Z",
			"files_count":  45,
			"fragments":    156,
			"languages":    []string{"Go", "TypeScript", "JavaScript"},
			"download_url": "/api/v1/lookatni/download/grompt_1726389000.zip",
		},
		{
			"id":           "proj_002",
			"name":         "gobe",
			"repo_url":     "https://github.com/kubex-ecosystem/gobe",
			"extracted_at": "2025-09-14T15:45:00Z",
			"files_count":  78,
			"fragments":    289,
			"languages":    []string{"Go", "TypeScript"},
			"download_url": "/api/v1/lookatni/download/gobe_1726302300.zip",
		},
	}

	// Apply pagination
	total := len(projects)
	start := offset
	end := start + limit
	if start >= total {
		projects = []map[string]interface{}{}
	} else {
		if end > total {
			end = total
		}
		projects = projects[start:end]
	}

	response := map[string]interface{}{
		"success":  true,
		"projects": projects,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
		"has_more": (offset + limit) < total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleProjectFragments handles GET /api/v1/lookatni/projects/{id}/fragments
func (h *Handler) HandleProjectFragments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract project ID from URL
	path := r.URL.Path
	projectID := filepath.Base(path)

	// Query parameters
	fragmentType := r.URL.Query().Get("type")
	language := r.URL.Query().Get("language")
	search := r.URL.Query().Get("search")

	// For now, return mock fragments - in a real implementation, this would query a database
	fragments := []lookatni.CodeFragment{
		{
			ID:           "frag_001",
			Type:         "function",
			Name:         "ExtractProject",
			FilePath:     "internal/services/lookatni/service.go",
			StartLine:    120,
			EndLine:      145,
			Content:      "func (s *LookAtniService) ExtractProject(ctx context.Context, req ProjectExtractionRequest) (*ExtractedProject, error) {\n\t// Implementation...\n}",
			Language:     "go",
			Complexity:   8,
			Dependencies: []string{"context", "ProjectExtractionRequest"},
		},
		{
			ID:           "frag_002",
			Type:         "interface",
			Name:         "ProjectExtractor",
			FilePath:     "internal/types/interfaces.go",
			StartLine:    15,
			EndLine:      25,
			Content:      "type ProjectExtractor interface {\n\tExtract(path string) (*Project, error)\n\tListFiles() []string\n}",
			Language:     "go",
			Complexity:   2,
			Dependencies: []string{"Project"},
		},
	}

	// Apply filters
	filtered := []lookatni.CodeFragment{}
	for _, frag := range fragments {
		include := true

		if fragmentType != "" && frag.Type != fragmentType {
			include = false
		}
		if language != "" && frag.Language != language {
			include = false
		}
		if search != "" && !strings.Contains(strings.ToLower(frag.Name), strings.ToLower(search)) {
			include = false
		}

		if include {
			filtered = append(filtered, frag)
		}
	}

	response := map[string]interface{}{
		"success":    true,
		"project_id": projectID,
		"fragments":  filtered,
		"total":      len(filtered),
		"filters": map[string]interface{}{
			"type":     fragmentType,
			"language": language,
			"search":   search,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
