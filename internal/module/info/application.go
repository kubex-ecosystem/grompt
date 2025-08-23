package manifest

import (
	_ "embed"
	"encoding/json"

	l "github.com/rafa-mori/logz"
)

//go:embed manifest.json
var manifestJSONData []byte
var application Manifest

type mmanifest struct {
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

func (m *mmanifest) GetName() string        { return m.Name }
func (m *mmanifest) GetVersion() string     { return m.Version }
func (m *mmanifest) GetAliases() []string   { return m.Aliases }
func (m *mmanifest) GetRepository() string  { return m.Repository }
func (m *mmanifest) GetHomepage() string    { return m.Homepage }
func (m *mmanifest) GetDescription() string { return m.Description }
func (m *mmanifest) GetMain() string        { return m.Main }
func (m *mmanifest) GetBin() string         { return m.Bin }
func (m *mmanifest) GetAuthor() string      { return m.Author }
func (m *mmanifest) GetLicense() string     { return m.License }
func (m *mmanifest) GetKeywords() []string  { return m.Keywords }
func (m *mmanifest) GetPlatforms() []string { return m.Platforms }
func (m *mmanifest) IsPrivate() bool        { return m.Private }

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

	var m mmanifest
	if err := json.Unmarshal(manifestJSONData, &m); err != nil {
		return nil, err
	}

	application = &m
	return application, nil
}
