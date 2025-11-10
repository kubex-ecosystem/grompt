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
	"github.com/kubex-ecosystem/grompt/internal/module/kbx"
	"github.com/kubex-ecosystem/logz"
)

type LegacyFileConfigImpl struct {
	Port               string            `json:"port,omitempty" yaml:"port,omitempty"`
	DefaultProvider    string            `json:"default_provider,omitempty" yaml:"default_provider,omitempty"`
	HistoryLimit       int               `json:"history_limit,omitempty" yaml:"history_limit,omitempty"`
	TimeoutSec         int               `json:"timeout_sec,omitempty" yaml:"timeout_sec,omitempty"`
	APIKeys            map[string]string `json:"api_keys,omitempty" yaml:"api_keys,omitempty"`
	Endpoints          map[string]string `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	Models             map[string]string `json:"models,omitempty" yaml:"models,omitempty"`
	ProviderConfigPath string            `json:"provider_config_path,omitempty" yaml:"provider_config_path,omitempty"`
}



// DefaultConfig rebuilds a legacy-compatible configuration.
func DefaultConfig(configFilePath string) interfaces.IConfig {
	cfg := newServerConfig()
	if configFilePath != "" {
		_ = cfg.LoadFromFile(configFilePath)
	}
	cfg.LoadFromEnv()
	return cfg
}

// NewConfig constructs a configuration using explicit parameters.

func NewServerConfig(
	name               string,
	debug              bool,
	logger             logz.Logger,
	bindAddr           string,
	port               string,
	tempDir            string,
	logFile           string,
	envFile           string,
	configFile        string,
	cwd               string,
	openAIKey         string,
	claudeKey        string,
	geminiKey        string,
	deepSeekKey      string,
	chatGPTKey      string,
	ollamaEndpoint   string,
	apiKeys            map[string]string,
	endpoints          map[string]string,
	defaultModels      map[string]string,
	providerTypes      map[string]string,
	defaultProvider    string,
	defaultTemperature float32,
	historyLimit       int,
	timeout            time.Duration,
	providerConfigPath string,
) interfaces.IConfig {
	cfg := newServerConfig()
	cfg.BindAddr = bindAddr
	if port != "" {
		cfg.Port = port
	}
	cfg.Logger = logger
	if openAIKey != "" {
		cfg.SetAPIKey("openai", openAIKey)
	}
	if claudeKey != "" {
		cfg.SetAPIKey("claude", claudeKey)
	}
	if geminiKey != "" {
		cfg.SetAPIKey("gemini", geminiKey)
	}
	if deepSeekKey != "" {
		cfg.SetAPIKey("deepseek", deepSeekKey)
	}
	if chatGPTKey != "" {
		cfg.SetAPIKey("chatgpt", chatGPTKey)
	}
	if ollamaEndpoint != "" {
		cfg.Endpoints["ollama"] = ollamaEndpoint
	}
	cfg.Name = name
	cfg.Debug = debug
	cfg.TempDir = tempDir
	cfg.LogFile = logFile
	cfg.EnvFile = envFile
	cfg.ConfigFile = configFile
	cfg.Cwd = cwd // pragma: allowlist secret
	for k, v := range apiKeys {
		if strings.ToLower(k) == "ollama" && v != "" {
			cfg.Endpoints["ollama"] = v
		} else {
			cfg.SetAPIKey(k, v)
		}
	}
	for k, v := range endpoints {
		cfg.Endpoints[k] = v
	}
	for k, v := range defaultModels {
		cfg.DefaultModels[k] = v
	}
	for k, v := range providerTypes {
		cfg.ProviderTypes[k] = v
	}
	if defaultProvider != "" {
		cfg.DefaultProvider = defaultProvider
	}
	if defaultTemperature >= 0 {
		cfg.DefaultTemperature = defaultTemperature
	}
	if historyLimit > 0 {
		cfg.HistoryLimit = historyLimit
	}
	if timeout > 0 {
		cfg.Timeout = timeout
	}
	cfg.ProviderConfigPath = providerConfigPath

	return cfg
}

// ---------- Config implementation ----------

var legacyProviders = []string{"openai", "claude", "gemini", "deepseek", "ollama", "chatgpt", "groq"}

type ServerConfigImpl struct {
	Name               string `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	Debug              bool   `json:"debug,omitempty" yaml:"debug,omitempty" mapstructure:"debug,omitempty"`
	Logger             logz.Logger `json:"-" yaml:"-" mapstructure:"-"`
	BindAddr           string `json:"bind_addr,omitempty" yaml:"bind_addr,omitempty" mapstructure:"bind_addr,omitempty"`
	Port               string `json:"port,omitempty" yaml:"port,omitempty" mapstructure:"port,omitempty"`
	TempDir            string `json:"temp_dir,omitempty" yaml:"temp_dir,omitempty" mapstructure:"temp_dir,omitempty"`
	LogFile            string `json:"log_file,omitempty" yaml:"log_file,omitempty" mapstructure:"log_file,omitempty"`
	EnvFile            string `json:"env_file,omitempty" yaml:"env_file,omitempty" mapstructure:"env_file,omitempty"`
	ConfigFile         string `json:"config_file,omitempty" yaml:"config_file,omitempty" mapstructure:"config_file,omitempty"`
	Cwd                string `json:"cwd,omitempty" yaml:"cwd,omitempty" mapstructure:"cwd,omitempty"`
	APIKeys            map[string]string `json:"api_keys,omitempty" yaml:"api_keys,omitempty" mapstructure:"api_keys,omitempty"`
	Endpoints          map[string]string `json:"endpoints,omitempty" yaml:"endpoints,omitempty" mapstructure:"endpoints,omitempty"`
	DefaultModels      map[string]string `json:"default_models,omitempty" yaml:"default_models,omitempty" mapstructure:"default_models,omitempty"`
	ProviderTypes      map[string]string `json:"provider_types,omitempty" yaml:"provider_types,omitempty" mapstructure:"provider_types,omitempty"`
	DefaultProvider    string `json:"default_provider,omitempty" yaml:"default_provider,omitempty" mapstructure:"default_provider,omitempty"`
	DefaultTemperature float32 `json:"default_temperature,omitempty" yaml:"default_temperature,omitempty" mapstructure:"default_temperature,omitempty"`
	HistoryLimit       int    `json:"history_limit,omitempty" yaml:"history_limit,omitempty" mapstructure:"history_limit,omitempty"`
	Timeout            time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty" mapstructure:"timeout,omitempty"`
	ProviderConfigPath string `json:"provider_config_path,omitempty" yaml:"provider_config_path,omitempty" mapstructure:"provider_config_path,omitempty"`

	Engine interfaces.IEngine `json:"-" yaml:"-" mapstructure:"-"`
	Mu     sync.RWMutex 		 `json:"-" yaml:"-" mapstructure:"-"`
}

func newServerConfig() *ServerConfigImpl {
	cfg := &ServerConfigImpl{
		Port:               "8080",
		APIKeys:            map[string]string{},
		Endpoints:          map[string]string{},
		DefaultModels:      map[string]string{},
		ProviderTypes:      map[string]string{},
		DefaultProvider:    "openai",
		DefaultTemperature: 0.7,
		HistoryLimit:       100,
		Timeout:            60 * time.Second,
	}

	cfg.ProviderTypes["openai"] = "openai"
	cfg.ProviderTypes["claude"] = "anthropic"
	cfg.ProviderTypes["gemini"] = "gemini"
	cfg.ProviderTypes["groq"] = "groq"

	cfg.DefaultModels["openai"] = "gpt-4o-mini"
	cfg.DefaultModels["claude"] = "claude-3-5-sonnet-20241022"
	cfg.DefaultModels["gemini"] = "gemini-1.5-pro"
	cfg.DefaultModels["groq"] = "llama-3.1-70b-versatile"

	cfg.Endpoints["openai"] = "https://api.openai.com"
	cfg.Endpoints["claude"] = "https://api.anthropic.com"
	cfg.Endpoints["gemini"] = "https://generativelanguage.googleapis.com"
	cfg.Endpoints["groq"] = "https://api.groq.com"

	return cfg
}

// func (c *ServerConfigImpl) attachEngine(engine interfaces.IEngine) {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()
// 	c.engine = engine
// }

func (c *ServerConfigImpl) GetAPIConfig(provider string) interfaces.IAPIConfig {
	return c.RegistryConfig().GetAPIConfig(provider)
}

func (c *ServerConfigImpl) GetLegacyAPIConfig(provider string) interfaces.LegacyAPIConfig {
	return &apiConfig{Provider: provider, APIServerConfig: c}
}

func (c *ServerConfigImpl) GetPort() string {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.Port
}

func (c *ServerConfigImpl) GetAPIKey(provider string) string {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return kbx.GetEnvOrDefault(c.APIKeys[strings.ToLower(provider)], "") // pragma: allowlist secret
}

func (c *ServerConfigImpl) SetAPIKey(provider, key string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	if key == "" {
		delete(c.APIKeys, strings.ToLower(provider))
		return nil
	}
	envKey := strings.ToUpper(provider) + "_API_KEY"
	if key != "" {
		if kbx.GetEnvOrDefault(key, "") != "" {
			key = kbx.GetEnvOrDefault(key, "")
		}

		if err := os.Setenv(envKey, key); err != nil {
			return fmt.Errorf("failed to set environment variable %s: %w", envKey, err)
		}
	}
	if kbx.GetEnvOrDefault(envKey, "") != "" {
		c.APIKeys[strings.ToLower(provider)] = envKey
	} else {
		c.APIKeys[strings.ToLower(provider)] = key
	}
	return nil
}

func (c *ServerConfigImpl) GetAPIEndpoint(provider string) string {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	return c.Endpoints[strings.ToLower(provider)]
}

func (c *ServerConfigImpl) GetBaseGenerationPrompt(ideas []string, purpose, purposeType, lang string, maxLength int) string {
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

func (c *ServerConfigImpl) LoadFromEnv() {
	for _, provider := range legacyProviders {
		if envKey := defaultEnvKey(provider); envKey != "" {
			if value := strings.TrimSpace(os.Getenv(envKey)); value != "" {
				c.APIKeys[provider] = value
			}
		}
	}

	if v := strings.TrimSpace(os.Getenv("GROMPT_DEFAULT_PROVIDER")); v != "" {
		c.DefaultProvider = strings.ToLower(v)
	}
	if v := strings.TrimSpace(os.Getenv("GROMPT_HISTORY_LIMIT")); v != "" {
		if parsed, err := parsePositiveInt(v); err == nil {
			c.HistoryLimit = parsed
		}
	}
	if v := strings.TrimSpace(os.Getenv("GROMPT_REQUEST_TIMEOUT")); v != "" {
		if dur, err := time.ParseDuration(v); err == nil {
			c.Timeout = dur
		}
	}
}

func (c *ServerConfigImpl) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var fileCfg LegacyFileConfigImpl
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

	c.Mu.Lock()
	defer c.Mu.Unlock()

	if fileCfg.Port != "" {
		c.Port = fileCfg.Port
	}
	if fileCfg.DefaultProvider != "" {
		c.DefaultProvider = strings.ToLower(fileCfg.DefaultProvider)
	}
	if fileCfg.HistoryLimit > 0 {
		c.HistoryLimit = fileCfg.HistoryLimit
	}
	if fileCfg.TimeoutSec > 0 {
		c.Timeout = time.Duration(fileCfg.TimeoutSec) * time.Second
	}

	mergeMap := func(dst map[string]string, src map[string]string) {
		for k, v := range src {
			if strings.TrimSpace(v) != "" {
				dst[strings.ToLower(k)] = v
			}
		}
	}

	mergeMap(c.APIKeys, fileCfg.APIKeys)
	mergeMap(c.Endpoints, fileCfg.Endpoints)

	for provider, model := range fileCfg.Models {
		if strings.TrimSpace(model) != "" {
			c.DefaultModels[strings.ToLower(provider)] = model
		}
	}

	if fileCfg.ProviderConfigPath != "" {
		c.ProviderConfigPath = fileCfg.ProviderConfigPath
	}

	return nil
}

func (c *ServerConfigImpl) RegistryConfig() interfaces.IConfig {
	cfg := NewServerConfig(
		c.BindAddr,
		c.Debug,
		c.Logger,
		c.BindAddr,
		c.Port,
		c.TempDir,
		c.LogFile,
		c.EnvFile,
		c.ConfigFile,
		c.Cwd,
		c.APIKeys["openai"],
		c.APIKeys["claude"],
		c.APIKeys["gemini"],
		c.APIKeys["deepseek"],
		c.APIKeys["chatgpt"],
		c.Endpoints["ollama"],
		c.APIKeys,
		c.Endpoints,
		c.DefaultModels,
		c.ProviderTypes,
		c.DefaultProvider,
		c.DefaultTemperature,
		c.HistoryLimit,
		c.Timeout,
		c.ProviderConfigPath,
	)

	return cfg
}

func (c *ServerConfigImpl) Validate() error {
	// Example validation: Ensure port is a valid number
	if _, err := net.LookupPort("tcp", c.GetPort()); err != nil {
		return fmt.Errorf("invalid port: %s", c.GetPort())
	}
	return c.RegistryConfig().Validate()
}
