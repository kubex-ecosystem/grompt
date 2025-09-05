// Package config provides configuration management for the factory.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	providersFct "github.com/rafa-mori/grompt/factory/providers"
	providersPkg "github.com/rafa-mori/grompt/internal/core/provider"
	"github.com/rafa-mori/grompt/internal/types"
	"gopkg.in/yaml.v3"

	l "github.com/rafa-mori/logz"
)

type Config = providersPkg.IConfig

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
) providersPkg.IConfig {
	return providersPkg.NewConfig(
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

func NewConfigFromFile(filePath string) providersPkg.IConfig {
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

func NewProvider(name, apiKey, version string) any /* providersPkg.Provider */ {
	cfg := providersPkg.NewConfig("", "", "", "", "", "", "", "", nil)
	return providersFct.NewProvider(name, apiKey, version, cfg)
}
