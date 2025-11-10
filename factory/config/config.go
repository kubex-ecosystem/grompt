// Package config provides legacy constructors that bridge to the new grompt configuration layer.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/interfaces"
	"github.com/kubex-ecosystem/grompt/internal/module/kbx"
	providersPkg "github.com/kubex-ecosystem/grompt/internal/providers"
	"github.com/kubex-ecosystem/grompt/internal/types"
	l "github.com/kubex-ecosystem/logz/logger"
	"gopkg.in/yaml.v3"
)

type Config = interfaces.IConfig

// NewConfig reconstructs a legacy configuration using the updated engine internals.
func NewConfig(
	name               string,
	debug              bool,
	logger             l.Logger,
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
	return types.NewConfig(
		name,
		debug,
		logger,
		bindAddr,
		port,
		tempDir,
		logFile,
		envFile,
		configFile,
		cwd,
		openAIKey,
		claudeKey,
		geminiKey,
		deepSeekKey,
		chatGPTKey,
		ollamaEndpoint,
		apiKeys,
		endpoints,
		defaultModels,
		providerTypes,
		defaultProvider,
		defaultTemperature,
		historyLimit,
		timeout,
		providerConfigPath,
	)
}

func NewConfigFromFile(filePath string) interfaces.IConfig {
	var cfg types.Config
	if _, statErr := os.Stat(filePath); statErr != nil {
		return &types.Config{}
	}
	switch fileExt := filepath.Ext(filePath); fileExt {
	case ".json":
		if err := readJSONFile(filePath, &cfg); err != nil {
			return &types.Config{}
		}
	case ".yaml", ".yml":
		if err := readYAMLFile(filePath, &cfg); err != nil {
			return &types.Config{}
		}
	default:
		return &types.Config{}
	}
	return &cfg
}

func readJSONFile(filePath string, cfg *types.Config) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(cfg)
}

func readYAMLFile(filePath string, cfg *types.Config) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	return decoder.Decode(cfg)
}

func NewProvider(name, apiKey, version string) providersPkg.Provider {
	defaultCfg := types.DefaultConfig("").(*types.Config)
	cfg := types.NewConfig(
		defaultCfg.Defaults.Name,
		defaultCfg.Defaults.Debug,
		l.LoggerG,
		defaultCfg.Defaults.Bind,
		defaultCfg.Defaults.Port,
		defaultCfg.Defaults.TempDir,
		defaultCfg.Defaults.LogFile,
		defaultCfg.Defaults.EnvFile,
		defaultCfg.Defaults.ConfigFile,
		defaultCfg.Defaults.Cwd,
		defaultCfg.Defaults.OpenAIKey,
		defaultCfg.Defaults.ClaudeKey,
		defaultCfg.Defaults.GeminiKey,
		defaultCfg.Defaults.DeepSeekKey,
		defaultCfg.Defaults.ChatGPTKey,
		defaultCfg.Defaults.OllamaEndpoint,
		make(map[string]string),
		make(map[string]string),
		make(map[string]string),
		make(map[string]string),
		kbx.DefaultLLMProvider,
		kbx.DefaultLLMTemperature,
		kbx.DefaultLLMHistoryLimit,
		kbx.DefaultTimeout,
		kbx.DefaultProviderConfigPath,
	)
	return providersPkg.NewProvider(name, apiKey, version, cfg)
}
