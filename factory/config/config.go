// Package config provides legacy constructors that bridge to the new grompt configuration layer.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	providersPkg "github.com/kubex-ecosystem/grompt/internal/providers"
	"github.com/kubex-ecosystem/grompt/internal/types"
	"gopkg.in/yaml.v3"

	l "github.com/kubex-ecosystem/logz"
)

type Config = types.IConfig

// NewConfig reconstructs a legacy configuration using the updated engine internals.
func NewConfig(
	bindAddr,
	port,
	openAIKey,
	deepSeekKey,
	ollamaEndpoint,
	claudeKey,
	geminiKey,
	chatGPTKey string,
	logger l.Logger,
) types.IConfig {
	return types.NewConfig(
		bindAddr,
		port,
		openAIKey,
		deepSeekKey,
		ollamaEndpoint,
		claudeKey,
		geminiKey,
		chatGPTKey,
		logger,
	)
}

func NewConfigFromFile(filePath string) types.IConfig {
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
	cfg := types.NewConfig("", "", "", "", "", "", "", "", nil)
	return providersPkg.NewProvider(name, apiKey, version, cfg)
}
