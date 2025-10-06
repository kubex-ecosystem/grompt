package wi18nast

import "time"

// ===== TIPOS DE DADOS =====

type I18nUsage struct {
	Key             string         `json:"key"`
	FilePath        string         `json:"filePath"`
	Line            int            `json:"line"`
	Column          int            `json:"column"`
	Component       string         `json:"component"`
	FunctionContext string         `json:"functionContext"`
	JSXContext      string         `json:"jsxContext"`
	Props           []string       `json:"props"`
	Imports         []string       `json:"imports"`
	NearbyCode      []string       `json:"nearbyCode"`
	CallType        string         `json:"callType"`
	Timestamp       time.Time      `json:"timestamp"`
	AIContext       *AIContextData `json:"aiContext,omitempty"`
}

type AIContextData struct {
	ComponentPurpose string   `json:"componentPurpose"`
	UIElementType    string   `json:"uiElementType"`
	UserInteraction  bool     `json:"userInteraction"`
	BusinessContext  string   `json:"businessContext"`
	SuggestedKeys    []string `json:"suggestedKeys"`
	QualityScore     int      `json:"qualityScore"`
}

type ProjectStats struct {
	TotalUsages       int               `json:"totalUsages"`
	CoveragePercent   float64           `json:"coveragePercent"`
	QualityScore      float64           `json:"qualityScore"`
	UsagesByType      map[string]int    `json:"usagesByType"`
	UsagesByComponent map[string]int    `json:"usagesByComponent"`
	MissingKeys       []string          `json:"missingKeys"`
	HardcodedStrings  []HardcodedString `json:"hardcodedStrings"`
	LastUpdate        time.Time         `json:"lastUpdate"`
}

type HardcodedString struct {
	Text     string `json:"text"`
	FilePath string `json:"filePath"`
	Line     int    `json:"line"`
	Context  string `json:"context"`
}

type ChangeEvent struct {
	Type      string      `json:"type"` // "added", "removed", "modified", "stats_updated"
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type I18nReport struct {
	Usages      []I18nUsage   `json:"usages"`
	Stats       ProjectStats  `json:"stats"`
	ChangeLog   []ChangeEvent `json:"changeLog"`
	GeneratedAt time.Time     `json:"generatedAt"`
}

type Usage struct {
	FilePath  string    `json:"filePath"`
	Line      int       `json:"line"`
	Column    int       `json:"column"`
	CallType  string    `json:"callType"`
	Component string    `json:"component,omitempty"`
	Key       string    `json:"key,omitempty"`
	JSXCtx    string    `json:"jsxContext,omitempty"`
	Nearby    []string  `json:"nearby,omitempty"`
	At        time.Time `json:"at"`
}

type Status string

const (
	StatusDraft    Status = "draft"
	StatusProposed Status = "proposed"
	StatusApproved Status = "approved"
)

type VaultItem struct {
	Key       string            `json:"key"`
	Text      string            `json:"text"`
	File      string            `json:"file"`
	Line      int               `json:"line"`
	Component string            `json:"component,omitempty"`
	Element   string            `json:"element,omitempty"`
	Contexts  []string          `json:"contexts,omitempty"`
	Status    Status            `json:"status"`
	Quality   int               `json:"quality,omitempty"`
	Notes     string            `json:"notes,omitempty"`
	Meta      map[string]string `json:"meta,omitempty"`
	FirstSeen time.Time         `json:"firstSeen"`
	LastSeen  time.Time         `json:"lastSeen"`
}
