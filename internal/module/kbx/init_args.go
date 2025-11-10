// Package kbx provides utilities for working with initialization arguments.
package kbx

import (
	"net"
	"os"
	"path/filepath"
	"reflect"
)

var (
	validKindMap = map[string]reflect.Kind{
		reflect.Struct.String():    reflect.Struct,
		reflect.Map.String():       reflect.Map,
		reflect.Slice.String():     reflect.Slice,
		reflect.Array.String():     reflect.Array,
		reflect.Chan.String():      reflect.Chan,
		reflect.Interface.String(): reflect.Interface,
		reflect.Ptr.String():       reflect.Ptr,
		reflect.String.String():    reflect.String,
		reflect.Int.String():       reflect.Int,
		reflect.Float32.String():   reflect.Float32,
		reflect.Float64.String():   reflect.Float64,
		reflect.Bool.String():      reflect.Bool,
		reflect.Uint.String():      reflect.Uint,
		reflect.Uint8.String():     reflect.Uint8,
		reflect.Uint16.String():    reflect.Uint16,
		reflect.Uint32.String():    reflect.Uint32,
		reflect.Uint64.String():    reflect.Uint64,
	}
)

type InitArgs struct {
	ConfigFile     string `yaml:"config_file,omitempty" json:"config_file,omitempty" mapstructure:"config_file,omitempty"`
	ConfigType     string `yaml:"config_type,omitempty" json:"config_type,omitempty" mapstructure:"config_type,omitempty"`
	ConfigDBFile   string `yaml:"config_db_file,omitempty" json:"config_db_file,omitempty" mapstructure:"config_db_file,omitempty"`
	ConfigDBType   string `yaml:"config_db_type,omitempty" json:"config_db_type,omitempty" mapstructure:"config_db_type,omitempty"`
	EnvFile        string `yaml:"env_file,omitempty" json:"env_file,omitempty" mapstructure:"env_file,omitempty"`
	LogFile        string `yaml:"log_file,omitempty" json:"log_file,omitempty" mapstructure:"log_file,omitempty"`
	Name           string `yaml:"name,omitempty" json:"name,omitempty" mapstructure:"name,omitempty"`
	Debug          bool   `yaml:"debug,omitempty" json:"debug,omitempty" mapstructure:"debug,omitempty"`
	ReleaseMode    bool   `yaml:"release_mode,omitempty" json:"release_mode,omitempty" mapstructure:"release_mode,omitempty"`
	IsConfidential bool   `yaml:"is_confidential,omitempty" json:"is_confidential,omitempty" mapstructure:"is_confidential,omitempty"`
	Port           string `yaml:"port,omitempty" json:"port,omitempty" mapstructure:"port,omitempty"`
	Bind           string `yaml:"bind,omitempty" json:"bind,omitempty" mapstructure:"bind,omitempty"`
	Address        string `yaml:"address,omitempty" json:"address,omitempty" mapstructure:"address,omitempty"`
	PubCertKeyPath string `yaml:"pub_cert_key_path,omitempty" json:"pub_cert_key_path,omitempty" mapstructure:"pub_cert_key_path,omitempty"`
	PubKeyPath     string `yaml:"pub_key_path,omitempty" json:"pub_key_path,omitempty" mapstructure:"pub_key_path,omitempty"`
	PrivKeyPath    string `yaml:"priv_key_path,omitempty" json:"priv_key_path,omitempty" mapstructure:"priv_key_path,omitempty"`
	Cwd            string `yaml:"cwd,omitempty" json:"cwd,omitempty" mapstructure:"cwd,omitempty"`
	TempDir       string `yaml:"temp_dir,omitempty" json:"temp_dir,omitempty" mapstructure:"temp_dir,omitempty"`
	OpenAIKey 		string `yaml:"openai_key,omitempty" json:"openai_key,omitempty" mapstructure:"openai_key,omitempty"`
	DeepSeekKey 	string `yaml:"deepseek_key,omitempty" json:"deepseek_key,omitempty" mapstructure:"deepseek_key,omitempty"`
	OllamaEndpoint string `yaml:"ollama_endpoint,omitempty" json:"ollama_endpoint,omitempty" mapstructure:"ollama_endpoint,omitempty"`
	GeminiKey 	string `yaml:"gemini_key,omitempty" json:"gemini_key,omitempty" mapstructure:"gemini_key,omitempty"`
	ChatGPTKey 	string `yaml:"chatgpt_key,omitempty" json:"chatgpt_key,omitempty" mapstructure:"chatgpt_key,omitempty"`
	ClaudeKey 	string `yaml:"claude_key,omitempty" json:"claude_key,omitempty" mapstructure:"claude_key,omitempty"`
	Background 	bool   `yaml:"background,omitempty" json:"background,omitempty" mapstructure:"background,omitempty"`

	// Extended configurations
	NotificationTimeoutSeconds int `yaml:"notification_timeout_seconds,omitempty" json:"notification_timeout_seconds,omitempty" mapstructure:"notification_timeout_seconds,omitempty"`
	NotificationProvider		any `yaml:"notification_provider,omitempty" json:"notification_provider,omitempty" mapstructure:"notification_provider,omitempty"`
	OpenAIConfig				map[string]any `yaml:"openai_config,omitempty" json:"openai_config,omitempty" mapstructure:"openai_config,omitempty"`
	ChatGPTConfig				map[string]any `yaml:"chatgpt_config,omitempty" json:"chatgpt_config,omitempty" mapstructure:"chatgpt_config,omitempty"`
	ClaudeConfig				map[string]any `yaml:"claude_config,omitempty" json:"claude_config,omitempty" mapstructure:"claude_config,omitempty"`
	AnthropicConfig			map[string]any `yaml:"anthropic_config,omitempty" json:"anthropic_config,omitempty" mapstructure:"anthropic_config,omitempty"`
	DeepSeekConfig			map[string]any `yaml:"deep_seek_config,omitempty" json:"deep_seek_config,omitempty" mapstructure:"deep_seek_config,omitempty"`
	GeminiConfig			map[string]any `yaml:"gemini_config,omitempty" json:"gemini_config,omitempty" mapstructure:"gemini_config,omitempty"`
	OllamaConfig			map[string]any `yaml:"ollama_config,omitempty" json:"ollama_config,omitempty" mapstructure:"ollama_config,omitempty"`
}

func NewInitArgs(
	configFile string,
	configType string,
	configDBFile string,
	configDBType string,
	envFile string,
	logFile string,
	name string,
	debug bool,
	releaseMode bool,
	isConfidential bool,
	port string,
	bind string,
	address string,
	pubCertKeyPath string,
	pubKeyPath string,
	cwd string,
	openAIKey string,
	deepSeekKey string,
	ollamaEndpoint string,
	geminiKey string,
	chatGPTKey string,
	claudeKey string,
	background bool,

	openAIConfig map[string]any,
	chatGPTConfig map[string]any,
	claudeConfig map[string]any,
	anthropicConfig map[string]any,
	deepSeekConfig map[string]any,
	geminiConfig map[string]any,
	ollamaConfig map[string]any,
) *InitArgs {
	configFile = GetValueOrDefaultSimple(configFile, os.ExpandEnv(DefaultGoBEConfigPath))
	configDBFile = GetValueOrDefaultSimple(configDBFile, "dbconfig.json")
	envFile = GetValueOrDefaultSimple(envFile, os.ExpandEnv(filepath.Join("$PWD", ".env")))
	logFile = GetValueOrDefaultSimple(
		logFile,
		filepath.Join(filepath.Dir(filepath.Dir(os.ExpandEnv(os.ExpandEnv(DefaultGoBEConfigPath)))), "logs", "gobe.log"),
	)
	port = GetValueOrDefaultSimple(port, "8088")
	bind = GetValueOrDefaultSimple(bind, "0.0.0.0")
	openAIKey = GetValueOrDefaultSimple(openAIKey, "")
	deepSeekKey = GetValueOrDefaultSimple(deepSeekKey, "")
	ollamaEndpoint = GetValueOrDefaultSimple(ollamaEndpoint, "http://localhost:11434")
	geminiKey = GetValueOrDefaultSimple(geminiKey, "")
	chatGPTKey = GetValueOrDefaultSimple(chatGPTKey, "")
	claudeKey = GetValueOrDefaultSimple(claudeKey, "")

	openAIConfig = GetValueOrDefaultSimple(openAIConfig, map[string]any{})
	chatGPTConfig = GetValueOrDefaultSimple(chatGPTConfig, map[string]any{})
	claudeConfig = GetValueOrDefaultSimple(claudeConfig, map[string]any{})
	anthropicConfig = GetValueOrDefaultSimple(anthropicConfig, map[string]any{})
	deepSeekConfig = GetValueOrDefaultSimple(deepSeekConfig, map[string]any{})
	geminiConfig = GetValueOrDefaultSimple(geminiConfig, map[string]any{})
	ollamaConfig = GetValueOrDefaultSimple(ollamaConfig, map[string]any{})

	return &InitArgs{
		ConfigFile:     configFile,
		ConfigType:     filepath.Ext(configFile)[1:],
		ConfigDBFile:   configDBFile,
		ConfigDBType:   filepath.Ext(configDBFile)[1:],
		EnvFile:        envFile,
		LogFile:        logFile,
		Name:           GetValueOrDefaultSimple(name, "GoBE"),
		Debug:          GetValueOrDefaultSimple(debug, false),
		ReleaseMode:    GetValueOrDefaultSimple(releaseMode, false),
		IsConfidential: GetValueOrDefaultSimple(isConfidential, false),
		Port:           port,
		Bind:           bind,
		Address:        net.JoinHostPort(bind, port),
		PubCertKeyPath: GetValueOrDefaultSimple(pubCertKeyPath, os.ExpandEnv(DefaultGoBEKeyPath)),
		PubKeyPath:     GetValueOrDefaultSimple(pubKeyPath, os.ExpandEnv(DefaultGoBECertPath)),
		Cwd:            GetValueOrDefaultSimple(cwd, ""),
		TempDir:       os.TempDir(),
		OpenAIKey: 		openAIKey,
		DeepSeekKey: 	deepSeekKey,
		OllamaEndpoint: ollamaEndpoint,
		GeminiKey: 	    geminiKey,
		ChatGPTKey: 	chatGPTKey,
		ClaudeKey: 	    claudeKey,
		Background: 	background,
		OpenAIConfig:   openAIConfig,
		ChatGPTConfig:  chatGPTConfig,
		ClaudeConfig:   claudeConfig,
		AnthropicConfig: anthropicConfig,
		DeepSeekConfig: deepSeekConfig,
		GeminiConfig:   geminiConfig,
		OllamaConfig:   ollamaConfig,
	}
}

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetValueOrDefault[T any](value T, defaultValue T) (T, reflect.Type) {
	if !IsObjValid(value) {
		return defaultValue, reflect.TypeFor[T]()
	}
	return value, reflect.TypeFor[T]()
}

func GetValueOrDefaultSimple[T any](value T, defaultValue T) T {
	if !IsObjValid(value) {
		return defaultValue
	}
	return value
}

func IsObjValid(obj any) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return false
		}
		if v.Kind() == reflect.Ptr {
			if v.Elem().Kind() == reflect.Ptr && v.Elem().IsNil() {
				return false
			}
			v = v.Elem()
		}
	}
	if _, ok := validKindMap[v.Kind().String()]; !ok {
		return false
	}
	if !v.IsValid() {
		return false
	}
	if v.IsZero() {
		return false
	}
	if v.Kind() == reflect.String && v.Len() == 0 {
		return false
	}
	if (v.Kind() == reflect.Slice || v.Kind() == reflect.Map || v.Kind() == reflect.Array) && v.Len() == 0 {
		return false
	}
	if v.Kind() == reflect.Bool {
		return true
	}
	return true
}
