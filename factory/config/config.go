// Package config provides configuration management for the factory.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	providersPkg "github.com/rafa-mori/grompt/internal/providers"
	"github.com/rafa-mori/grompt/internal/types"
	"gopkg.in/yaml.v3"
)

type Config = types.IConfig

func NewConfig(
	port string,
	openAIKey string,
	deepSeekKey string,
	ollamaEndpoint string,
	claudeKey string,
	geminiKey string,
) types.IConfig {
	return types.NewConfig(port, openAIKey, deepSeekKey, ollamaEndpoint, claudeKey, geminiKey)
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
	cfg := types.NewConfig("", "", "", "", "", "")
	return providersPkg.NewProvider(name, apiKey, version, cfg)
}
