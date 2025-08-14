// Package manifest provides functionality to read and parse the application manifest.
package manifest

import (
	_ "embed"
	"encoding/json"

	l "github.com/rafa-mori/logz"
)

//go:embed manifest.json
var manifestJSONData []byte
var application Manifest

type manifest struct {
	Manifest
	Name            string   `json:"name"`
	ApplicationName string   `json:"application"`
	Bin             string   `json:"bin"`
	Version         string   `json:"version"`
	Repository      string   `json:"repository"`
	Aliases         []string `json:"aliases,omitempty"`
	Homepage        string   `json:"homepage,omitempty"`
	Description     string   `json:"description,omitempty"`
	Main            string   `json:"main,omitempty"`
	Author          string   `json:"author,omitempty"`
	License         string   `json:"license,omitempty"`
	Keywords        []string `json:"keywords,omitempty"`
	Platforms       []string `json:"platforms,omitempty"`
	LogLevel        string   `json:"log_level,omitempty"`
	Debug           bool     `json:"debug,omitempty"`
	ShowTrace       bool     `json:"show_trace,omitempty"`
	Private         bool     `json:"private,omitempty"`
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

func init() {
	_, err := GetManifest()
	if err != nil {
		l.GetLogger("Kubex")
		l.Fatal("Failed to get manifest: " + err.Error())
	}
}

func GetManifest() (Manifest, error) {
	if application != nil {
		return application, nil
	}

	var m manifest
	if err := json.Unmarshal(manifestJSONData, &m); err != nil {
		return nil, err
	}

	application = &m
	return application, nil
}
