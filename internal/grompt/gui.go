// Package grompt provides functionality for the Grompt application.
package grompt

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kubex-ecosystem/grompt/internal/grompt/embedded"
)


var guiWebFS *embedded.EmbeddedFS

func init(){
	guiWebFS = embedded.GetEmbeddedFS()
}

// GUIGrompt analyzes GUI-related metrics and provides insights
type GUIGrompt struct {
	guiWebFS *embedded.EmbeddedFS
}

// NewGUIGrompt creates a new instance of GUIGrompt
func NewGUIGrompt() *GUIGrompt {
	return &GUIGrompt{
		guiWebFS: guiWebFS,
	}
}

// GetWebFS returns the embedded filesystem for GUI web assets
func (g *GUIGrompt) GetWebFS() fs.FS {
	if g == nil {
		return nil
	}
	if g.guiWebFS == nil {
		g.guiWebFS = guiWebFS
	}
	return g.guiWebFS.F
}

func (g *GUIGrompt) GetEmbeddedFS() *embedded.EmbeddedFS {
	if g == nil {
		return nil
	}
	if g.guiWebFS == nil {
		g.guiWebFS = guiWebFS
	}
	return g.guiWebFS
}

// GetWebRoot returns the root directory for GUI web assets
func (g *GUIGrompt) GetWebRoot(path string) *os.DirEntry {
	if g == nil {
		return nil
	}
	for _, entry := range g.guiWebFS.S {
		if entry.Name() == path {
			return &entry
		}
	}
	return nil
}

// GetWebFile retrieves a specific file from the embedded GUI web assets
func (g *GUIGrompt) GetWebFile(path string) ([]byte, error) {
	if g == nil {
		return nil, os.ErrNotExist
	}
	file, err := g.guiWebFS.F.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := fs.ReadFile(g.guiWebFS.F, filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	return data, nil
}
