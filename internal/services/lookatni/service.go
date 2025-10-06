// Package lookatni provides integration with lookatni library for code extraction and analysis
package lookatni

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// LookAtniService handles lookatni operations
type LookAtniService struct {
	workDir     string
	nodeModules string
	timeout     time.Duration
}

// NewLookAtniService creates a new lookatni service
func NewLookAtniService(workDir string) *LookAtniService {
	return &LookAtniService{
		workDir:     workDir,
		nodeModules: filepath.Join(workDir, "node_modules"),
		timeout:     5 * time.Minute,
	}
}

// GetWorkDir returns the working directory
func (s *LookAtniService) GetWorkDir() string {
	return s.workDir
}

// ProjectExtractionRequest represents a request to extract project files
type ProjectExtractionRequest struct {
	RepoURL         string   `json:"repo_url"`
	LocalPath       string   `json:"local_path,omitempty"`
	IncludePatterns []string `json:"include_patterns,omitempty"`
	ExcludePatterns []string `json:"exclude_patterns,omitempty"`
	MaxFileSize     int64    `json:"max_file_size,omitempty"`
	IncludeHidden   bool     `json:"include_hidden,omitempty"`
	ContextDepth    int      `json:"context_depth,omitempty"`
	FragmentBy      string   `json:"fragment_by,omitempty"` // "file", "function", "class", "module"
}

// ExtractedProject represents the result of project extraction
type ExtractedProject struct {
	ProjectName string           `json:"project_name"`
	Structure   ProjectStructure `json:"structure"`
	Files       []ExtractedFile  `json:"files"`
	Fragments   []CodeFragment   `json:"fragments"`
	Metadata    ProjectMetadata  `json:"metadata"`
	DownloadURL string           `json:"download_url,omitempty"`
	ExtractedAt time.Time        `json:"extracted_at"`
}

// ProjectStructure represents the hierarchical structure of the project
type ProjectStructure struct {
	Root        string          `json:"root"`
	Directories []DirectoryNode `json:"directories"`
	TotalFiles  int             `json:"total_files"`
	TotalSize   int64           `json:"total_size"`
}

// DirectoryNode represents a directory in the project structure
type DirectoryNode struct {
	Name     string          `json:"name"`
	Path     string          `json:"path"`
	Files    []FileNode      `json:"files"`
	Children []DirectoryNode `json:"children"`
}

// FileNode represents a file in the project structure
type FileNode struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
	Language string `json:"language"`
}

// ExtractedFile represents a file extracted from the project
type ExtractedFile struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Content   string            `json:"content"`
	Language  string            `json:"language"`
	Size      int64             `json:"size"`
	LineCount int               `json:"line_count"`
	Fragments []CodeFragment    `json:"fragments"`
	Metadata  map[string]string `json:"metadata"`
}

// CodeFragment represents a logical code fragment (function, class, etc.)
type CodeFragment struct {
	ID           string            `json:"id"`
	Type         string            `json:"type"` // "function", "class", "method", "interface", "struct"
	Name         string            `json:"name"`
	FilePath     string            `json:"file_path"`
	StartLine    int               `json:"start_line"`
	EndLine      int               `json:"end_line"`
	Content      string            `json:"content"`
	Language     string            `json:"language"`
	Complexity   int               `json:"complexity,omitempty"`
	Dependencies []string          `json:"dependencies,omitempty"`
	Metadata     map[string]string `json:"metadata"`
}

// ProjectMetadata contains metadata about the extracted project
type ProjectMetadata struct {
	Languages      map[string]int `json:"languages"`
	TotalLines     int            `json:"total_lines"`
	TotalFiles     int            `json:"total_files"`
	TotalFragments int            `json:"total_fragments"`
	ExtractionTime time.Duration  `json:"extraction_time"`
	GitInfo        GitInfo        `json:"git_info,omitempty"`
}

// GitInfo contains git repository information
type GitInfo struct {
	Branch       string    `json:"branch"`
	LastCommit   string    `json:"last_commit"`
	LastCommitAt time.Time `json:"last_commit_at"`
	Contributors []string  `json:"contributors"`
	RemoteURL    string    `json:"remote_url"`
}

// ExtractProject extracts a project using lookatni
func (s *LookAtniService) ExtractProject(ctx context.Context, req ProjectExtractionRequest) (*ExtractedProject, error) {
	startTime := time.Now()

	// Prepare temporary directory for extraction
	tempDir := filepath.Join(s.workDir, "temp", fmt.Sprintf("extract_%d", time.Now().Unix()))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Clone or copy project if needed
	projectPath := req.LocalPath
	if projectPath == "" && req.RepoURL != "" {
		clonedPath, err := s.cloneRepository(ctx, req.RepoURL, tempDir)
		if err != nil {
			return nil, fmt.Errorf("failed to clone repository: %w", err)
		}
		projectPath = clonedPath
	}

	// Execute lookatni extraction
	extracted, err := s.executeLookAtni(ctx, projectPath, req)
	if err != nil {
		return nil, fmt.Errorf("lookatni extraction failed: %w", err)
	}

	// Post-process and enhance extraction
	enhanced, err := s.enhanceExtraction(extracted, req)
	if err != nil {
		return nil, fmt.Errorf("failed to enhance extraction: %w", err)
	}

	enhanced.ExtractedAt = time.Now()
	enhanced.Metadata.ExtractionTime = time.Since(startTime)

	return enhanced, nil
}

// executeLookAtni runs lookatni as Node.js library
func (s *LookAtniService) executeLookAtni(ctx context.Context, projectPath string, req ProjectExtractionRequest) (*ExtractedProject, error) {
	// Create lookatni configuration
	config := map[string]interface{}{
		"inputPath":       projectPath,
		"includePatterns": req.IncludePatterns,
		"excludePatterns": req.ExcludePatterns,
		"maxFileSize":     req.MaxFileSize,
		"includeHidden":   req.IncludeHidden,
		"contextDepth":    req.ContextDepth,
		"fragmentBy":      req.FragmentBy,
		"outputFormat":    "json",
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	// Create Node.js script to run lookatni
	script := fmt.Sprintf(`
const lookatni = require('lookatni');
const config = %s;

async function extractProject() {
  try {
    const result = await lookatni.extract(config);
    console.log(JSON.stringify(result, null, 2));
  } catch (error) {
    console.error('Extraction error:', error);
    process.exit(1);
  }
}

extractProject();
`, string(configJSON))

	scriptPath := filepath.Join(s.workDir, "temp", "extract_script.js")
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		return nil, fmt.Errorf("failed to write script: %w", err)
	}
	defer os.Remove(scriptPath)

	// Execute Node.js script
	cmd := exec.CommandContext(ctx, "node", scriptPath)
	cmd.Dir = s.workDir

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("lookatni execution failed: %w", err)
	}

	// Parse lookatni output
	var extracted ExtractedProject
	if err := json.Unmarshal(output, &extracted); err != nil {
		return nil, fmt.Errorf("failed to parse lookatni output: %w", err)
	}

	return &extracted, nil
}

// enhanceExtraction adds additional metadata and processing
func (s *LookAtniService) enhanceExtraction(extracted *ExtractedProject, req ProjectExtractionRequest) (*ExtractedProject, error) {
	// Generate unique IDs for files and fragments
	for i := range extracted.Files {
		file := &extracted.Files[i]
		file.ID = fmt.Sprintf("file_%d_%s", i, strings.ReplaceAll(file.Name, ".", "_"))

		for j := range file.Fragments {
			fragment := &file.Fragments[j]
			fragment.ID = fmt.Sprintf("frag_%d_%d_%s", i, j, fragment.Type)
			fragment.FilePath = file.Path
		}
	}

	// Generate fragments list for easy access
	var allFragments []CodeFragment
	for _, file := range extracted.Files {
		allFragments = append(allFragments, file.Fragments...)
	}
	extracted.Fragments = allFragments

	// Update metadata
	extracted.Metadata.TotalFragments = len(allFragments)

	return extracted, nil
}

// cloneRepository clones a git repository to a temporary directory
func (s *LookAtniService) cloneRepository(ctx context.Context, repoURL, tempDir string) (string, error) {
	clonePath := filepath.Join(tempDir, "repo")

	cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "1", repoURL, clonePath)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git clone failed: %w", err)
	}

	return clonePath, nil
}

// CreateNavigableArchive creates a navigable archive for download
func (s *LookAtniService) CreateNavigableArchive(ctx context.Context, extracted *ExtractedProject) (string, error) {
	// Create archive directory
	archiveDir := filepath.Join(s.workDir, "temp", fmt.Sprintf("archive_%d", time.Now().Unix()))
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create archive directory: %w", err)
	}

	// Generate HTML navigation interface
	navHTML := s.generateNavigationHTML(extracted)
	if err := os.WriteFile(filepath.Join(archiveDir, "index.html"), []byte(navHTML), 0644); err != nil {
		return "", fmt.Errorf("failed to write navigation HTML: %w", err)
	}

	// Copy extracted files maintaining structure
	for _, file := range extracted.Files {
		filePath := filepath.Join(archiveDir, "files", file.Path)
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return "", fmt.Errorf("failed to create file directory: %w", err)
		}
		if err := os.WriteFile(filePath, []byte(file.Content), 0644); err != nil {
			return "", fmt.Errorf("failed to write file: %w", err)
		}
	}

	// Generate metadata JSON
	metadataJSON, _ := json.MarshalIndent(extracted, "", "  ")
	if err := os.WriteFile(filepath.Join(archiveDir, "metadata.json"), metadataJSON, 0644); err != nil {
		return "", fmt.Errorf("failed to write metadata: %w", err)
	}

	// Create ZIP archive
	zipPath := filepath.Join(s.workDir, "downloads", fmt.Sprintf("%s_%d.zip", extracted.ProjectName, time.Now().Unix()))
	if err := os.MkdirAll(filepath.Dir(zipPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create downloads directory: %w", err)
	}

	if err := s.createZipArchive(archiveDir, zipPath); err != nil {
		return "", fmt.Errorf("failed to create ZIP archive: %w", err)
	}

	// Clean up temporary archive directory
	os.RemoveAll(archiveDir)

	return zipPath, nil
}

// generateNavigationHTML creates an HTML interface for project navigation
func (s *LookAtniService) generateNavigationHTML(extracted *ExtractedProject) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - Project Navigator</title>
    <style>
        body { font-family: monospace; background: #1a1a1a; color: #fff; margin: 0; padding: 20px; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { border-bottom: 2px solid #333; padding-bottom: 20px; margin-bottom: 20px; }
        .project-info { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-bottom: 20px; }
        .info-card { background: #2a2a2a; padding: 15px; border-radius: 8px; }
        .file-tree { background: #2a2a2a; padding: 15px; border-radius: 8px; }
        .file-item { padding: 5px; cursor: pointer; border-radius: 4px; }
        .file-item:hover { background: #3a3a3a; }
        .fragment-list { background: #2a2a2a; padding: 15px; border-radius: 8px; margin-top: 20px; }
        .fragment-item { padding: 10px; margin: 5px 0; background: #333; border-radius: 4px; cursor: pointer; }
        .fragment-item:hover { background: #444; }
        .code-view { background: #1e1e1e; padding: 20px; border-radius: 8px; margin-top: 20px; }
        pre { margin: 0; white-space: pre-wrap; }
        .hidden { display: none; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔍 %s Project Navigator</h1>
            <p>Generated by LookAtni • %d files • %d fragments • %d total lines</p>
        </div>

        <div class="project-info">
            <div class="info-card">
                <h3>📊 Statistics</h3>
                <p>Files: %d</p>
                <p>Lines: %d</p>
                <p>Fragments: %d</p>
                <p>Languages: %v</p>
            </div>
            <div class="info-card">
                <h3>📁 Structure</h3>
                <div class="file-tree" id="fileTree">
                    <!-- File tree will be generated by JavaScript -->
                </div>
            </div>
        </div>

        <div class="fragment-list">
            <h3>🧩 Code Fragments</h3>
            <div id="fragmentList">
                <!-- Fragments will be generated by JavaScript -->
            </div>
        </div>

        <div class="code-view hidden" id="codeView">
            <h3 id="codeTitle">Code Viewer</h3>
            <pre id="codeContent"></pre>
        </div>
    </div>

    <script>
        const projectData = %s;

        function showFile(filePath) {
            const file = projectData.files.find(f => f.path === filePath);
            if (file) {
                document.getElementById('codeTitle').textContent = file.name;
                document.getElementById('codeContent').textContent = file.content;
                document.getElementById('codeView').classList.remove('hidden');
            }
        }

        function showFragment(fragmentId) {
            const fragment = projectData.fragments.find(f => f.id === fragmentId);
            if (fragment) {
                document.getElementById('codeTitle').textContent =
                    fragment.name + ' (' + fragment.type + ')';
                document.getElementById('codeContent').textContent = fragment.content;
                document.getElementById('codeView').classList.remove('hidden');
            }
        }

        // Generate file tree
        const fileTree = document.getElementById('fileTree');
        projectData.files.forEach(file => {
            const item = document.createElement('div');
            item.className = 'file-item';
            item.textContent = file.path;
            item.onclick = () => showFile(file.path);
            fileTree.appendChild(item);
        });

        // Generate fragment list
        const fragmentList = document.getElementById('fragmentList');
        projectData.fragments.forEach(fragment => {
            const item = document.createElement('div');
            item.className = 'fragment-item';
            item.innerHTML = '<strong>' + fragment.name + '</strong> (' + fragment.type + ') - ' + fragment.file_path;
            item.onclick = () => showFragment(fragment.id);
            fragmentList.appendChild(item);
        });
    </script>
</body>
</html>
`, extracted.ProjectName, extracted.ProjectName, len(extracted.Files), len(extracted.Fragments),
		extracted.Metadata.TotalLines, extracted.Metadata.TotalFiles, extracted.Metadata.TotalLines,
		extracted.Metadata.TotalFragments, extracted.Metadata.Languages, toJSON(extracted))
}

// createZipArchive creates a ZIP archive from a directory
func (s *LookAtniService) createZipArchive(sourceDir, zipPath string) error {
	cmd := exec.Command("zip", "-r", zipPath, ".")
	cmd.Dir = sourceDir
	return cmd.Run()
}

// toJSON converts struct to JSON string for JavaScript embedding
func toJSON(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
