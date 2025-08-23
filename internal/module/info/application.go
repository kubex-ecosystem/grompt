// Package info provides functionality to read and parse the application manifest.
package info

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed manifest.json
var manifestJSONData []byte

var (
	application Manifest
)

type Reference struct {
	Name            string `json:"name"`
	ApplicationName string `json:"application"`
	Bin             string `json:"bin"`
	Version         string `json:"version"`
}

type manifest struct {
	Manifest // Interface Constraint

	Reference   `mapstructure:",squash"` // Identifiers that will be merged in the root of the struct
	Repository  string                   `json:"repository"`
	Aliases     []string                 `json:"aliases,omitempty"`
	Homepage    string                   `json:"homepage,omitempty"`
	Description string                   `json:"description,omitempty"`
	Main        string                   `json:"main,omitempty"`
	Author      string                   `json:"author,omitempty"`
	License     string                   `json:"license,omitempty"`
	Keywords    []string                 `json:"keywords,omitempty"`
	Platforms   []string                 `json:"platforms,omitempty"`
	LogLevel    string                   `json:"log_level,omitempty"`
	Debug       bool                     `json:"debug,omitempty"`
	ShowTrace   bool                     `json:"show_trace,omitempty"`
	Private     bool                     `json:"private,omitempty"`
}
type Manifest interface {
	GetName() string
	GetVersion() string
	GetAliases() []string
	GetRepository() string
	GetHomepage() string
	GetDescription() string
	GetMain() string
	GetBin() string
	GetAuthor() string
	GetLicense() string
	GetKeywords() []string
	GetPlatforms() []string
	IsPrivate() bool
}

func (m *manifest) GetName() string        { return m.Name }
func (m *manifest) GetVersion() string     { return m.Version }
func (m *manifest) GetAliases() []string   { return m.Aliases }
func (m *manifest) GetRepository() string  { return m.Repository }
func (m *manifest) GetHomepage() string    { return m.Homepage }
func (m *manifest) GetDescription() string { return m.Description }
func (m *manifest) GetMain() string        { return m.Main }
func (m *manifest) GetBin() string         { return m.Bin }
func (m *manifest) GetAuthor() string      { return m.Author }
func (m *manifest) GetLicense() string     { return m.License }
func (m *manifest) GetKeywords() []string  { return m.Keywords }
func (m *manifest) GetPlatforms() []string { return m.Platforms }
func (m *manifest) IsPrivate() bool        { return m.Private }

// lazy cache
var (
	cachedManifest *Manifest
	cachedControl  *Control
)

// GetManifest lazy, sem init() com side-effects
func GetManifest() (Manifest, error) {
	if cachedManifest != nil {
		return *cachedManifest, nil
	}
	var m Manifest
	if len(manifestJSONData) == 0 {
		return nil, fmt.Errorf("manifest.json: embed is empty")
	}
	if err := json.Unmarshal(manifestJSONData, &m); err != nil {
		return nil, fmt.Errorf("manifest.json: %w", err)
	}
	cachedManifest = &m
	return m, nil
}

// FS secOrder quiser permitir override por FS externo:
type FS interface {
	ReadFile(name string) ([]byte, error)
}

func LoadFromFS(fs FS) (Manifest, Control, error) {
	var m Manifest
	var c Control
	if b, err := fs.ReadFile("manifest.json"); err == nil {
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, Control{}, fmt.Errorf("manifest.json: %w", err)
		}
	} else {
		return nil, Control{}, fmt.Errorf("manifest.json: %w", err)
	}
	if b, err := fs.ReadFile("control.json"); err == nil {
		if err := json.Unmarshal(b, &c); err != nil {
			return nil, Control{}, fmt.Errorf("control.json: %w", err)
		}
	} else {
		return nil, Control{}, fmt.Errorf("control.json: %w", err)
	}
	return m, c, nil
}

// func GetControl() (*Control, error) {
// 	if cachedControl != nil {
// 		return cachedControl, nil
// 	}
// 	var c Control
// 	if len(controlJSONData) == 0 {
// 		return nil, fmt.Errorf("control.json: embed is empty")
// 	}
// 	if err := json.Unmarshal(controlJSONData, &c); err != nil {
// 		return nil, fmt.Errorf("control.json: %w", err)
// 	}
// 	cachedControl = &c
// 	return &c, nil
// }
