package types

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/kubex-ecosystem/grompt/internal/interfaces"
	"github.com/kubex-ecosystem/logz"
)

type legacyFileConfig struct {
	Port               string            `json:"port" yaml:"port"`
	DefaultProvider    string            `json:"default_provider" yaml:"default_provider"`
	HistoryLimit       int               `json:"history_limit" yaml:"history_limit"`
	TimeoutSec         int               `json:"timeout_sec" yaml:"timeout_sec"`
	APIKeys            map[string]string `json:"api_keys" yaml:"api_keys"`
	Endpoints          map[string]string `json:"endpoints" yaml:"endpoints"`
	Models             map[string]string `json:"models" yaml:"models"`
	ProviderConfigPath string            `json:"provider_config" yaml:"provider_config"`
}



// DefaultConfig rebuilds a legacy-compatible configuration.
func DefaultConfig(configFilePath string) interfaces.ServerConfig {
	cfg := newServerConfig()
	if configFilePath != "" {
		_ = cfg.loadFromFile(configFilePath)
	}
	cfg.loadFromEnv()
	return cfg
}

// NewConfig constructs a configuration using explicit parameters.

func NewServerConfig(
	bindAddr string,
	port string,
	openAIKey string,
	deepSeekKey string,
	ollamaEndpoint string,
	claudeKey string,
	geminiKey string,
	chatGPTKey string,
	logger logz.Logger,
) interfaces.ServerConfig {
	cfg := newServerConfig()
	cfg.bindAddr = bindAddr
	if port != "" {
		cfg.port = port
	}
	cfg.logger = logger
	if openAIKey != "" {
		cfg.apiKeys["openai"] = openAIKey
	}
	if claudeKey != "" {
		cfg.apiKeys["claude"] = claudeKey
	}
	if geminiKey != "" {
		cfg.apiKeys["gemini"] = geminiKey
	}
	if deepSeekKey != "" {
		cfg.apiKeys["deepseek"] = deepSeekKey
	}
	if chatGPTKey != "" {
		cfg.apiKeys["chatgpt"] = chatGPTKey
	}
	if ollamaEndpoint != "" {
		cfg.endpoints["ollama"] = ollamaEndpoint
	}
	return cfg
}

// ---------- Config implementation ----------

var legacyProviders = []string{"openai", "claude", "gemini", "deepseek", "ollama", "chatgpt", "groq"}

type configImpl struct {
	logger             logz.Logger
	bindAddr           string
	port               string
	apiKeys            map[string]string
	endpoints          map[string]string
	defaultModels      map[string]string
	providerTypes      map[string]string
	defaultProvider    string
	defaultTemperature float32
	historyLimit       int
	timeout            time.Duration
	providerConfigPath string

	engine interfaces.IEngine
	mu     sync.RWMutex
}

func newServerConfig() *configImpl {
	cfg := &configImpl{
		port:               "8080",
		apiKeys:            map[string]string{},
		endpoints:          map[string]string{},
		defaultModels:      map[string]string{},
		providerTypes:      map[string]string{},
		defaultProvider:    "openai",
		defaultTemperature: 0.7,
		historyLimit:       100,
		timeout:            60 * time.Second,
	}

	cfg.providerTypes["openai"] = "openai"
	cfg.providerTypes["claude"] = "anthropic"
	cfg.providerTypes["gemini"] = "gemini"
	cfg.providerTypes["groq"] = "groq"

	cfg.defaultModels["openai"] = "gpt-4o-mini"
	cfg.defaultModels["claude"] = "claude-3-5-sonnet-20241022"
	cfg.defaultModels["gemini"] = "gemini-1.5-pro"
	cfg.defaultModels["groq"] = "llama-3.1-70b-versatile"

	cfg.endpoints["openai"] = "https://api.openai.com"
	cfg.endpoints["claude"] = "https://api.anthropic.com"
	cfg.endpoints["gemini"] = "https://generativelanguage.googleapis.com"
	cfg.endpoints["groq"] = "https://api.groq.com"

	return cfg
}

// func (c *configImpl) attachEngine(engine interfaces.IEngine) {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	c.engine = engine
// }

func (c *configImpl) GetAPIConfig(provider string) interfaces.IAPIConfig {
	return c.registryConfig().GetAPIConfig(provider)
}

func (c *configImpl) GetLegacyAPIConfig(provider string) interfaces.LegacyAPIConfig {
	return &apiConfig{provider: provider, cfg: c}
}

func (c *configImpl) GetPort() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.port
}

func (c *configImpl) GetAPIKey(provider string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.apiKeys[strings.ToLower(provider)]
}

func (c *configImpl) SetAPIKey(provider, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if key == "" {
		delete(c.apiKeys, strings.ToLower(provider))
		return nil
	}
	c.apiKeys[strings.ToLower(provider)] = key
	return nil
}

func (c *configImpl) GetAPIEndpoint(provider string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.endpoints[strings.ToLower(provider)]
}

func (c *configImpl) GetBaseGenerationPrompt(ideas []string, purpose, purposeType, lang string, maxLength int) string {
	var builder strings.Builder
	builder.WriteString("You are Kubex Grompt assistant.")
	if purpose != "" {
		builder.WriteString("Purpose: " + purpose + "")
	}
	if purposeType != "" {
		builder.WriteString("Type: " + purposeType + "")
	}
	if lang != "" {
		builder.WriteString("Language: " + lang + "")
	}
	if maxLength > 0 {
		builder.WriteString(fmt.Sprintf("Limit response to %d characters.", maxLength))
	}
	if len(ideas) > 0 {
		builder.WriteString("Ideas:")
		for _, idea := range ideas {
			builder.WriteString("- " + idea + "")
		}
	}
	builder.WriteString("Respond with detailed, actionable output.")
	return builder.String()
}

func (c *configImpl) loadFromEnv() {
	for _, provider := range legacyProviders {
		if envKey := defaultEnvKey(provider); envKey != "" {
			if value := strings.TrimSpace(os.Getenv(envKey)); value != "" {
				c.apiKeys[provider] = value
			}
		}
	}

	if v := strings.TrimSpace(os.Getenv("GROMPT_DEFAULT_PROVIDER")); v != "" {
		c.defaultProvider = strings.ToLower(v)
	}
	if v := strings.TrimSpace(os.Getenv("GROMPT_HISTORY_LIMIT")); v != "" {
		if parsed, err := parsePositiveInt(v); err == nil {
			c.historyLimit = parsed
		}
	}
	if v := strings.TrimSpace(os.Getenv("GROMPT_REQUEST_TIMEOUT")); v != "" {
		if dur, err := time.ParseDuration(v); err == nil {
			c.timeout = dur
		}
	}
}

func (c *configImpl) loadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var fileCfg legacyFileConfig
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &fileCfg); err != nil {
			return err
		}
	case ".json":
		if err := json.Unmarshal(data, &fileCfg); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported config file extension: %s", filepath.Ext(path))
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if fileCfg.Port != "" {
		c.port = fileCfg.Port
	}
	if fileCfg.DefaultProvider != "" {
		c.defaultProvider = strings.ToLower(fileCfg.DefaultProvider)
	}
	if fileCfg.HistoryLimit > 0 {
		c.historyLimit = fileCfg.HistoryLimit
	}
	if fileCfg.TimeoutSec > 0 {
		c.timeout = time.Duration(fileCfg.TimeoutSec) * time.Second
	}

	mergeMap := func(dst map[string]string, src map[string]string) {
		for k, v := range src {
			if strings.TrimSpace(v) != "" {
				dst[strings.ToLower(k)] = v
			}
		}
	}

	mergeMap(c.apiKeys, fileCfg.APIKeys)
	mergeMap(c.endpoints, fileCfg.Endpoints)

	for provider, model := range fileCfg.Models {
		if strings.TrimSpace(model) != "" {
			c.defaultModels[strings.ToLower(provider)] = model
		}
	}

	if fileCfg.ProviderConfigPath != "" {
		c.providerConfigPath = fileCfg.ProviderConfigPath
	}

	return nil
}

func (c *configImpl) registryConfig() interfaces.ServerConfig {
	// cfg := Config{
	// 	Providers: map[string]ProviderImpl{},
	// }

	// c.mu.RLock()
	// defer c.mu.RUnlock()

	// for name, apiKey := range c.apiKeys { // pragma: allowlist secret
	// 	providerType := c.providerTypes[name]
	// 	if providerType == "" {
	// 		providerType = name
	// 	}


	// }
	cfg := NewServerConfig(
		c.bindAddr,
		c.port,
		c.apiKeys["openai"],
		c.apiKeys["deepseek"],
		c.endpoints["ollama"],
		c.apiKeys["claude"],
		c.apiKeys["gemini"],
		c.apiKeys["chatgpt"],
		c.logger,
	)

	return cfg
}


func (c *configImpl) Validate() error {
	// Example validation: Ensure port is a valid number
	if _, err := net.LookupPort("tcp", c.GetPort()); err != nil {
		return fmt.Errorf("invalid port: %s", c.GetPort())
	}
	return c.registryConfig().Validate()
}
