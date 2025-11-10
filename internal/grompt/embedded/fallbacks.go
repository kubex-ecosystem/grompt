// Package embedded provides embedded resources for the Grompt application.
package embedded

import (
	"embed"
	"io/fs"
)

type EmbeddedFS struct {
	F fs.FS
	S []fs.DirEntry
}

//go:embed all:guiweb/*
var fallbacksFS embed.FS

var (
	EFb *EmbeddedFS
)

// getEmbeddedFS retrieves a specific fallback file from the embedded resources
func getEmbeddedFS() fs.FS {
	if EFb == nil {
		EFb = &EmbeddedFS{}
		var err error
		EFb.F, err = fs.Sub(&fallbacksFS, "guiweb")
		if err != nil {
			return nil
		}
		EFb.S, err = fs.ReadDir(EFb.F, ".")
		if err != nil {
			return nil
		}
	}
	return EFb.F
}

func GetEmbeddedDirEntries() []fs.DirEntry {
	if EFb.S == nil {
		return nil
	}
	return EFb.S
}

func GetEmbeddedFile(name string) (fs.File, error) {
	if EFb.F == nil {
		return nil, fs.ErrNotExist
	}
	return EFb.F.Open(name)
}

func GetEmbeddedFS() *EmbeddedFS {
	if EFb == nil {
		EFb = &EmbeddedFS{
			F: getEmbeddedFS(),
			S: GetEmbeddedDirEntries(),
		}
	}
	return EFb
}

func init() {
	if EFb != nil {
		return
	}
	EFb = &EmbeddedFS{
		F: getEmbeddedFS(),
		S: GetEmbeddedDirEntries(),
	}
}
