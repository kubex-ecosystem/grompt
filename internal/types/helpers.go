package types

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/kubex-ecosystem/grompt/internal/interfaces"
)

// ---------- Helpers ----------

func defaultEnvKey(provider string) string {
	switch provider {
	case "openai":
		return "OPENAI_API_KEY"
	case "claude":
		return "ANTHROPIC_API_KEY"
	case "gemini":
		return "GEMINI_API_KEY"
	case "deepseek":
		return "DEEPSEEK_API_KEY"
	case "chatgpt":
		return "CHATGPT_API_KEY"
	case "groq":
		return "GROQ_API_KEY"
	default:
		return ""
	}
}

func parsePositiveInt(value string) (int, error) {
	var parsed int
	_, err := fmt.Sscanf(value, "%d", &parsed)
	if err != nil {
		return 0, err
	}
	if parsed <= 0 {
		return 0, fmt.Errorf("value must be positive")
	}
	return parsed, nil
}

func ExecuteTemplate(tmpl string, vars map[string]interface{}) (string, error) {
	tpl, err := template.New("grompt").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	if err := tpl.Execute(&buff, vars); err != nil {
		return "", err
	}

	return buff.String(), nil
}

func defaultCapabilities(provider, model string) *interfaces.Capabilities {
	caps := &interfaces.Capabilities{
		SupportsBatch:     true,
		SupportsStreaming: true,
		Models:            map[string]any{},
	}

	if model != "" {
		caps.Models[model] = struct{}{}
	}

	switch provider {
	case "openai":
		caps.MaxTokens = 128000
		caps.Pricing = &interfaces.Pricing{InputCostPer1K: 0.0005, OutputCostPer1K: 0.0015, Currency: "USD"}
	case "claude":
		caps.MaxTokens = 200000
		caps.Pricing = &interfaces.Pricing{InputCostPer1K: 0.003, OutputCostPer1K: 0.015, Currency: "USD"}
	case "gemini":
		caps.MaxTokens = 1000000
		caps.Pricing = &interfaces.Pricing{InputCostPer1K: 0.000125, OutputCostPer1K: 0.000375, Currency: "USD"}
	case "groq":
		caps.MaxTokens = 128000
		caps.Pricing = &interfaces.Pricing{InputCostPer1K: 0.0002, OutputCostPer1K: 0.0002, Currency: "USD"}
	default:
		caps.MaxTokens = 32000
	}

	return caps
}
