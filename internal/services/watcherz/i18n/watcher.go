package wi18nast

// Package wi18nast implements internationalization (i18n) support.

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type OnAnalyze func(file string, usages []Usage)

type Watcher struct {
	projectPath string
	parser      *Parser
	cb          OnAnalyze

	watcher *fsnotify.Watcher
	mu      sync.Mutex
	pending map[string]time.Time // debounce por arquivo
	stopCh  chan struct{}
}

func NewWatcher(projectPath string, cb OnAnalyze) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	p := NewParser()

	wt := &Watcher{
		projectPath: projectPath,
		parser:      p,
		cb:          cb,
		watcher:     w,
		pending:     make(map[string]time.Time),
		stopCh:      make(chan struct{}),
	}

	// adiciona diretórios recursivamente
	if err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if skipDir(path) {
				return filepath.SkipDir
			}
			return w.Add(path)
		}
		return nil
	}); err != nil {
		_ = w.Close()
		return nil, err
	}

	return wt, nil

}

func (wt *Watcher) Start() {
	go wt.loop()
	go wt.debouncer()
}

func (wt *Watcher) Stop() {
	close(wt.stopCh)
	_ = wt.watcher.Close()
}

func (wt *Watcher) loop() {
	for {
		select {
		case ev, ok := <-wt.watcher.Events:
			if !ok {
				return
			}
			if !isSourceFile(filepath.Ext(ev.Name)) {
				continue
			}
			if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
				wt.mu.Lock()
				wt.pending[ev.Name] = time.Now()
				wt.mu.Unlock()
			}
			if ev.Op&fsnotify.Remove != 0 {
				wt.cb(ev.Name, nil) // arquivo removido → sem usages
			}
		case err, ok := <-wt.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("[watcher] erro: %v", err)
		case <-wt.stopCh:
			return
		}
	}
}

func (wt *Watcher) debouncer() {
	t := time.NewTicker(150 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			now := time.Now()
			var batch []string
			wt.mu.Lock()
			for f, ts := range wt.pending {
				if now.Sub(ts) > 120*time.Millisecond {
					batch = append(batch, f)
					delete(wt.pending, f)
				}
			}
			wt.mu.Unlock()

			for _, f := range batch {
				usages, err := wt.parser.ParseFile(f)
				if err != nil {
					log.Printf("[parser] %s: %v", f, err)
					continue
				}
				wt.cb(f, usages)
			}
		case <-wt.stopCh:
			return
		}
	}

}

// ---------- Parser (regex-based) ----------

type Parser struct {
	// Regex patterns for i18n usage detection
	tCallPattern          *regexp.Regexp
	useTranslationPattern *regexp.Regexp
	transPattern          *regexp.Regexp
}

func NewParser() *Parser {
	return &Parser{
		tCallPattern:          regexp.MustCompile(`\bt\s*\(\s*['"](.*?)['"]`),
		useTranslationPattern: regexp.MustCompile(`useTranslation\s*\(`),
		transPattern:          regexp.MustCompile(`<Trans\b|<Translation\b`),
	}
}

func (p *Parser) ParseFile(filePath string) ([]Usage, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", filePath, err)
	}

	src := string(b)
	lines := strings.Split(src, "\n")

	var usages []Usage
	p.parseWithRegex(src, lines, filePath, &usages)
	return usages, nil
}

func (p *Parser) parseWithRegex(src string, lines []string, file string, out *[]Usage) {
	// Find t() calls
	matches := p.tCallPattern.FindAllStringSubmatchIndex(src, -1)
	for _, match := range matches {
		if len(match) >= 4 {
			start := match[0]
			key := src[match[2]:match[3]]
			line, col := p.getLineColumn(src, start)

			usage := Usage{
				FilePath:  file,
				Line:      line,
				Column:    col,
				CallType:  "t()",
				Key:       key,
				Component: p.findComponentFromLines(lines, line),
				Nearby:    snippet(lines, line, 2),
				At:        time.Now(),
			}
			*out = append(*out, usage)
		}
	}

	// Find useTranslation hooks
	matches = p.useTranslationPattern.FindAllStringIndex(src, -1)
	for _, match := range matches {
		start := match[0]
		line, col := p.getLineColumn(src, start)

		usage := Usage{
			FilePath:  file,
			Line:      line,
			Column:    col,
			CallType:  "useTranslation",
			Key:       "hook_usage",
			Component: p.findComponentFromLines(lines, line),
			Nearby:    snippet(lines, line, 2),
			At:        time.Now(),
		}
		*out = append(*out, usage)
	}

	// Find Trans components
	matches = p.transPattern.FindAllStringIndex(src, -1)
	for _, match := range matches {
		start := match[0]
		line, col := p.getLineColumn(src, start)

		usage := Usage{
			FilePath:  file,
			Line:      line,
			Column:    col,
			CallType:  "Trans",
			Key:       "component_usage",
			Component: p.findComponentFromLines(lines, line),
			Nearby:    snippet(lines, line, 2),
			At:        time.Now(),
		}
		*out = append(*out, usage)
	}
}

func (p *Parser) getLineColumn(src string, pos int) (int, int) {
	line := 1
	col := 1
	for i := 0; i < pos && i < len(src); i++ {
		if src[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return line, col
}

func (p *Parser) findComponentFromLines(lines []string, targetLine int) string {
	// Look backwards for function/component declaration
	for i := targetLine - 1; i >= 0 && i >= targetLine-20; i-- {
		if i < len(lines) {
			line := strings.TrimSpace(lines[i])
			if strings.Contains(line, "function ") || strings.Contains(line, "const ") && strings.Contains(line, "=>") {
				// Extract component name with simple regex
				funcPattern := regexp.MustCompile(`(?:function|const)\s+(\w+)`)
				if matches := funcPattern.FindStringSubmatch(line); len(matches) > 1 {
					return matches[1]
				}
			}
		}
	}
	return "unknown"
}

// ---------- helpers ----------

func snippet(lines []string, center, radius int) []string {
	var out []string
	start := max(1, center-radius)
	end := min(len(lines), center+radius)
	for i := start; i <= end; i++ {
		out = append(out, fmt.Sprintf("%4d  %s", i, lines[i-1]))
	}
	return out
}

func isSourceFile(ext string) bool {
	switch strings.ToLower(ext) {
	case ".ts", ".tsx", ".js", ".jsx":
		return true
	}
	return false
}

func skipDir(path string) bool {
	return strings.Contains(path, "node_modules") || strings.Contains(path, string(filepath.Separator)+".git")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
