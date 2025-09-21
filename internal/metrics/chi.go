// Package metrics - Code Health Index (CHI) calculation engine
package metrics

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/types"
)

// CHICalculator calculates Code Health Index metrics
type CHICalculator struct {
	repoPath string
}

// NewCHICalculator creates a new CHI calculator
func NewCHICalculator(repoPath string) *CHICalculator {
	return &CHICalculator{
		repoPath: repoPath,
	}
}

// CodeFile represents a source code file analysis
type CodeFile struct {
	Path                 string
	Language             string
	Lines                int
	LinesOfCode          int
	CyclomaticComplexity int
	Functions            int
	TestFile             bool
	Duplications         []Duplication
}

// Duplication represents code duplication detection
type Duplication struct {
	StartLine int
	EndLine   int
	Hash      string
	Content   string
}

// Calculate computes Code Health Index for a repository
func (c *CHICalculator) Calculate(ctx context.Context, repo types.Repository) (*types.CHIMetrics, error) {
	if c.repoPath == "" {
		return nil, fmt.Errorf("repository path not set")
	}

	files, err := c.analyzeCodebase(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze codebase: %w", err)
	}

	duplicationPct := c.calculateDuplication(files)
	cyclomaticAvg := c.calculateCyclomaticComplexity(files)
	testCoverage := c.calculateTestCoverage(files)
	maintainabilityIndex := c.calculateMaintainabilityIndex(files)
	technicalDebt := c.calculateTechnicalDebt(files)

	// Calculate overall CHI score (0-100)
	chiScore := c.calculateCHIScore(duplicationPct, cyclomaticAvg, testCoverage, maintainabilityIndex)

	return &types.CHIMetrics{
		Score:                chiScore,
		DuplicationPercent:   duplicationPct,
		CyclomaticComplexity: cyclomaticAvg,
		TestCoverage:         testCoverage,
		MaintainabilityIndex: maintainabilityIndex,
		TechnicalDebt:        technicalDebt,
		Period:               60, // Default to 60 days
		CalculatedAt:         time.Now(),
	}, nil
}

// analyzeCodebase walks through the repository and analyzes code files
func (c *CHICalculator) analyzeCodebase(ctx context.Context) ([]CodeFile, error) {
	var files []CodeFile

	err := filepath.Walk(c.repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-code files
		if info.IsDir() || !c.isCodeFile(path) {
			return nil
		}

		// Skip vendor, node_modules, .git etc.
		if c.shouldSkipPath(path) {
			return nil
		}

		file, err := c.analyzeFile(path)
		if err != nil {
			// Log error but continue with other files
			return nil
		}

		files = append(files, *file)
		return nil
	})

	return files, err
}

// isCodeFile determines if a file is a source code file
func (c *CHICalculator) isCodeFile(path string) bool {
	ext := filepath.Ext(path)
	codeExtensions := []string{
		".go", ".js", ".ts", ".jsx", ".tsx",
		".py", ".java", ".kt", ".cs", ".cpp", ".c",
		".rb", ".php", ".swift", ".rs", ".scala",
	}

	for _, codeExt := range codeExtensions {
		if ext == codeExt {
			return true
		}
	}
	return false
}

// shouldSkipPath determines if a path should be skipped during analysis
func (c *CHICalculator) shouldSkipPath(path string) bool {
	skipDirs := []string{
		".git", "node_modules", "vendor", "dist", "build",
		"target", ".next", "coverage", "__pycache__",
	}

	for _, skipDir := range skipDirs {
		if strings.Contains(path, skipDir) {
			return true
		}
	}
	return false
}

// analyzeFile analyzes a single source code file
func (c *CHICalculator) analyzeFile(path string) (*CodeFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	linesOfCode := c.countLinesOfCode(lines)

	file := &CodeFile{
		Path:        path,
		Language:    c.detectLanguage(path),
		Lines:       len(lines),
		LinesOfCode: linesOfCode,
		TestFile:    c.isTestFile(path),
	}

	// Language-specific analysis
	if file.Language == "go" {
		c.analyzeGoFile(file, content)
	}

	// Generic duplication detection
	file.Duplications = c.detectDuplications(lines)

	return file, nil
}

// detectLanguage detects the programming language of a file
func (c *CHICalculator) detectLanguage(path string) string {
	ext := filepath.Ext(path)
	languageMap := map[string]string{
		".go":    "go",
		".js":    "javascript",
		".ts":    "typescript",
		".jsx":   "javascript",
		".tsx":   "typescript",
		".py":    "python",
		".java":  "java",
		".kt":    "kotlin",
		".cs":    "csharp",
		".cpp":   "cpp",
		".c":     "c",
		".rb":    "ruby",
		".php":   "php",
		".swift": "swift",
		".rs":    "rust",
		".scala": "scala",
	}

	if lang, ok := languageMap[ext]; ok {
		return lang
	}
	return "unknown"
}

// isTestFile determines if a file is a test file
func (c *CHICalculator) isTestFile(path string) bool {
	testPatterns := []string{"_test.go", ".test.js", ".test.ts", "test_", "_spec."}

	for _, pattern := range testPatterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

// countLinesOfCode counts actual lines of code (excluding comments and blank lines)
func (c *CHICalculator) countLinesOfCode(lines []string) int {
	loc := 0
	inBlockComment := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Handle block comments (/* */)
		if strings.Contains(trimmed, "/*") {
			inBlockComment = true
		}
		if strings.Contains(trimmed, "*/") {
			inBlockComment = false
			continue
		}
		if inBlockComment {
			continue
		}

		// Skip single line comments
		if strings.HasPrefix(trimmed, "//") ||
			strings.HasPrefix(trimmed, "#") ||
			strings.HasPrefix(trimmed, "--") {
			continue
		}

		loc++
	}

	return loc
}

// analyzeGoFile performs Go-specific analysis
func (c *CHICalculator) analyzeGoFile(file *CodeFile, content []byte) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file.Path, content, 0)
	if err != nil {
		return
	}

	// Count functions and calculate cyclomatic complexity
	ast.Inspect(node, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			file.Functions++
			file.CyclomaticComplexity += c.calculateGoComplexity(fn)
		}
		return true
	})
}

// calculateGoComplexity calculates cyclomatic complexity for a Go function
func (c *CHICalculator) calculateGoComplexity(fn *ast.FuncDecl) int {
	complexity := 1 // Base complexity

	ast.Inspect(fn, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt,
			*ast.TypeSwitchStmt, *ast.SelectStmt:
			complexity++
		case *ast.CaseClause:
			complexity++
		}
		return true
	})

	return complexity
}

// detectDuplications detects code duplications in a file
func (c *CHICalculator) detectDuplications(lines []string) []Duplication {
	var duplications []Duplication
	// Simplified duplication detection - look for identical line sequences
	// In production, use more sophisticated algorithms like suffix trees

	minDuplicationLines := 5

	for i := 0; i < len(lines)-minDuplicationLines; i++ {
		for j := i + minDuplicationLines; j < len(lines)-minDuplicationLines; j++ {
			matchLength := 0
			for k := 0; k < minDuplicationLines && i+k < len(lines) && j+k < len(lines); k++ {
				if strings.TrimSpace(lines[i+k]) == strings.TrimSpace(lines[j+k]) && strings.TrimSpace(lines[i+k]) != "" {
					matchLength++
				} else {
					break
				}
			}

			if matchLength >= minDuplicationLines {
				content := strings.Join(lines[i:i+matchLength], "\n")
				duplications = append(duplications, Duplication{
					StartLine: i + 1,
					EndLine:   i + matchLength,
					Hash:      fmt.Sprintf("%d", hash(content)),
					Content:   content,
				})
			}
		}
	}

	return duplications
}

// Simple hash function for duplication detection
func hash(s string) uint32 {
	h := uint32(0)
	for _, c := range s {
		h = h*31 + uint32(c)
	}
	return h
}

// calculateDuplication calculates the percentage of duplicated code
func (c *CHICalculator) calculateDuplication(files []CodeFile) float64 {
	totalLines := 0
	duplicatedLines := 0

	for _, file := range files {
		totalLines += file.LinesOfCode
		for _, dup := range file.Duplications {
			duplicatedLines += dup.EndLine - dup.StartLine + 1
		}
	}

	if totalLines == 0 {
		return 0
	}

	return (float64(duplicatedLines) / float64(totalLines)) * 100.0
}

// calculateCyclomaticComplexity calculates average cyclomatic complexity
func (c *CHICalculator) calculateCyclomaticComplexity(files []CodeFile) float64 {
	totalComplexity := 0
	totalFunctions := 0

	for _, file := range files {
		totalComplexity += file.CyclomaticComplexity
		totalFunctions += file.Functions
	}

	if totalFunctions == 0 {
		return 0
	}

	return float64(totalComplexity) / float64(totalFunctions)
}

// calculateTestCoverage estimates test coverage based on test files
func (c *CHICalculator) calculateTestCoverage(files []CodeFile) float64 {
	testFiles := 0
	codeFiles := 0

	for _, file := range files {
		if file.TestFile {
			testFiles++
		} else {
			codeFiles++
		}
	}

	if codeFiles == 0 {
		return 0
	}

	// Simplified coverage estimation: ratio of test files to code files
	coverage := (float64(testFiles) / float64(codeFiles)) * 100.0
	if coverage > 100 {
		coverage = 100
	}

	return coverage
}

// calculateMaintainabilityIndex calculates maintainability index
func (c *CHICalculator) calculateMaintainabilityIndex(files []CodeFile) float64 {
	if len(files) == 0 {
		return 0
	}

	totalLOC := 0
	totalComplexity := 0
	totalFunctions := 0

	for _, file := range files {
		totalLOC += file.LinesOfCode
		totalComplexity += file.CyclomaticComplexity
		totalFunctions += file.Functions
	}

	if totalFunctions == 0 {
		return 0
	}

	avgComplexity := float64(totalComplexity) / float64(totalFunctions)
	avgLOC := float64(totalLOC) / float64(len(files))

	// Simplified maintainability index calculation
	// Real formula is more complex and includes Halstead metrics
	mi := 171 - 5.2*math.Log(avgLOC) - 0.23*avgComplexity - 16.2*math.Log(avgLOC/float64(totalFunctions))

	if mi < 0 {
		mi = 0
	}
	if mi > 100 {
		mi = 100
	}

	return mi
}

// calculateTechnicalDebt estimates technical debt in hours
func (c *CHICalculator) calculateTechnicalDebt(files []CodeFile) float64 {
	debt := 0.0

	for _, file := range files {
		// Debt from duplication: 30 minutes per duplicated block
		debt += float64(len(file.Duplications)) * 0.5

		// Debt from high complexity: 1 hour per complex function
		if file.Functions > 0 {
			avgComplexity := float64(file.CyclomaticComplexity) / float64(file.Functions)
			if avgComplexity > 10 { // High complexity threshold
				debt += (avgComplexity - 10) * 1.0
			}
		}

		// Debt from large files: 2 hours per 1000 LOC over threshold
		if file.LinesOfCode > 500 {
			debt += float64(file.LinesOfCode-500) / 1000.0 * 2.0
		}
	}

	return debt
}

// calculateCHIScore calculates overall Code Health Index score (0-100)
func (c *CHICalculator) calculateCHIScore(duplication, complexity, testCoverage, maintainability float64) int {
	// Weighted scoring
	duplicationScore := math.Max(0, 100-(duplication*2))  // Weight: 2x penalty for duplication
	complexityScore := math.Max(0, 100-(complexity-5)*10) // Penalty starts at complexity > 5
	testScore := testCoverage                             // Direct test coverage percentage
	maintainabilityScore := maintainability               // Direct maintainability index

	// Weighted average
	weights := []float64{0.3, 0.25, 0.25, 0.2} // duplication, complexity, tests, maintainability
	scores := []float64{duplicationScore, complexityScore, testScore, maintainabilityScore}

	weightedSum := 0.0
	for i, score := range scores {
		weightedSum += score * weights[i]
	}

	chi := int(math.Round(weightedSum))
	if chi < 0 {
		chi = 0
	}
	if chi > 100 {
		chi = 100
	}

	return chi
}
