// Package types defines the configuration and versioning for the Grompt application.
package types

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"time"

	"github.com/kubex-ecosystem/grompt/internal/interfaces"
	"github.com/kubex-ecosystem/grompt/internal/module/kbx"
	vs "github.com/kubex-ecosystem/grompt/internal/module/version"
	l "github.com/kubex-ecosystem/logz/logger"
)

var (
	CurrentVersion string    = vs.GetVersion()
	LatestVersion  string    = vs.GetLatestVersionFromGit()
	LastCheckTime  time.Time = time.Now()
)

func init() {
	// Initialize the CurrentVersion and LatestVersion
	// This will run in a goroutine to avoid blocking the main execution
	// and will check for the latest version from Git if needed.
	// If the CurrentVersion is not set, it will use the version from the version package
	// and will update the LatestVersion if the last check was more than 24 hours ago
	// or if it is the first run.
	ctx := context.Background()
	cancel, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	// Use a goroutine to avoid blocking the main execution
	go func(cancel context.Context) {
		select {
		case <-cancel.Done():
			return // Exit if the context is cancelled
		default:
			if CurrentVersion == "" {
				CurrentVersion = vs.GetVersion()
			}
			if LastCheckTime.IsZero() || LastCheckTime.Before(time.Now().Add(-24*time.Hour)) {
				// Check for the latest version from Git
				LatestVersion = vs.GetLatestVersionFromGit()
				if LatestVersion == "" {
					LatestVersion = CurrentVersion // Fallback to current version if check fails
				}
				LastCheckTime = time.Now()
			}
		}
	}(cancel)

	// // Ensure that the CurrentVersion is always set
	// if CurrentVersion == "" {
	// 	CurrentVersion = vs.GetVersion()
	// }
	// if LatestVersion == "" {
	// 	LatestVersion = vs.GetLatestVersionFromGit()
	// }
}

var (
	vSrv = vs.NewVersionService()
	glgr   = l.LoggerG
)

var (
	AppVersion = vSrv.GetVersion()
	AppName    = vSrv.GetName()
)

type APIConfig struct {
	APIKey     string `json:"api_key,omitempty" yaml:"api_key,omitempty" mapstructure:"api_key,omitempty"`
	BaseURL    string `json:"base_url,omitempty" yaml:"base_url,omitempty" mapstructure:"base_url,omitempty"`
	APIVersion    string `json:"api_version,omitempty" yaml:"api_version,omitempty" mapstructure:"api_version,omitempty"`
	HTTPClient *http.Client `json:"-" yaml:"-" mapstructure:"-"`
	DemoMode   bool  `json:"demo_mode,omitempty" yaml:"demo_mode,omitempty" mapstructure:"demo_mode,omitempty"`
}

type Config struct {
	Logger    l.Logger 										`yaml:"-" json:"-" mapstructure:"-"`
	Server    *ServerConfigImpl              `yaml:"server,omitempty" json:"server,omitempty"`
	Defaults  *kbx.InitArgs                  `yaml:"defaults,omitempty" json:"defaults,omitempty"`
	Providers map[string]interfaces.Provider `yaml:"providers,omitempty" json:"providers,omitempty"`

	BindAddr       string `json:"bind_addr,omitempty" gorm:"default:'localhost'"`
	Port           string `json:"port" gorm:"default:8080"`
	OpenAIAPIKey   string `json:"openai_api_key,omitempty" gorm:"default:''"`
	DeepSeekAPIKey string `json:"deepseek_api_key,omitempty" gorm:"default:''"`
	OllamaEndpoint string `json:"ollama_endpoint,omitempty" gorm:"default:'http://localhost:11434'"`
	ClaudeAPIKey   string `json:"claude_api_key,omitempty" gorm:"default:''"`
	GeminiAPIKey   string `json:"gemini_api_key,omitempty" gorm:"default:''"`
	ChatGPTAPIKey  string `json:"chatgpt_api_key,omitempty" gorm:"default:''"`
	Debug          bool   `json:"debug" gorm:"default:false"`
	EnableCORS		 bool   `json:"enable_cors" gorm:"default:false"`
}

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
) *Config {
	if logger != nil {
		glgr = logger
	} else{
		glgr = l.LoggerG.GetLogger()
	}
	if bindAddr == "" {
		bindAddr = kbx.DefaultServerHost
	}
	server := NewServerConfig(
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
	).(*ServerConfigImpl)
	cfg := &Config{
		Server:         server,
		Logger:         logger,
		BindAddr:       bindAddr,
		Port:           port,
	}
	cfg.SetAPIKey("openai", openAIKey)
	cfg.SetAPIKey("deepseek", deepSeekKey)
	cfg.SetAPIKey("ollama", ollamaEndpoint)
	cfg.SetAPIKey("claude", claudeKey)
	cfg.SetAPIKey("gemini", geminiKey)
	cfg.SetAPIKey("chatgpt", chatGPTKey)

	return cfg
}

func (c *Config) GetAPIConfig(provider string) interfaces.IAPIConfig {
	if c == nil {
		glgr.Log("error", "Config is nil")
		return nil
	}
	switch provider {
	case "openai":
		return NewOpenAIAPI(c.GetAPIKey("openai"))
	case "deepseek":
		return NewDeepSeekAPI(c.GetAPIKey("deepseek"))
	case "ollama":
		return NewOllamaAPI(c.GetAPIEndpoint("ollama"))
	case "claude":
		return NewClaudeAPI(c.GetAPIKey("claude"))
	case "gemini":
		return NewGeminiAPI(c.GetAPIKey("gemini"))
	case "chatgpt":
		return NewChatGPTAPI(c.GetAPIKey("chatgpt"))
	default:
		return nil
	}
}

func (c *Config) GetPort() string {
	if c.Port == "" {
		return kbx.DefaultServerPort
	}
	return c.Port
}

func (c *Config) GetAPIKey(provider string) string {
	switch provider {
	case "openai":
		if kbx.GetEnvOrDefault("OPENAI_API_KEY", c.OpenAIAPIKey) != "" {
			return "OPENAI_API_KEY"
		}
	case "deepseek":
		if kbx.GetEnvOrDefault("DEEPSEEK_API_KEY", c.DeepSeekAPIKey) != "" {
			return "DEEPSEEK_API_KEY"
		}
	case "claude":
		if kbx.GetEnvOrDefault("CLAUDE_API_KEY", c.ClaudeAPIKey) != "" {
			return "CLAUDE_API_KEY"
		}
	case "gemini":
		if kbx.GetEnvOrDefault("GEMINI_API_KEY", c.GeminiAPIKey) != "" {
			return "GEMINI_API_KEY"
		}
	case "chatgpt":
		if kbx.GetEnvOrDefault("CHATGPT_API_KEY", c.ChatGPTAPIKey) != "" {
			return "CHATGPT_API_KEY"
		}
	case "ollama":
		if kbx.GetEnvOrDefault("OLLAMA_ENDPOINT", c.OllamaEndpoint) != "" {
			return "OLLAMA_ENDPOINT"
		}
	}
	return ""
}

func (c *Config) SetAPIKey(provider string, key string) error {
	switch provider {
	case "openai":
		if err := os.Setenv("OPENAI_API_KEY", key); err == nil && key != "" { // pragma: allowlist secret
			c.OpenAIAPIKey = "OPENAI_API_KEY" // pragma: allowlist secret
		}else {
			c.OpenAIAPIKey = ""
		}
	case "deepseek":
		if err := os.Setenv("DEEPSEEK_API_KEY", key); err == nil && key != "" { // pragma: allowlist secret
			c.DeepSeekAPIKey = "DEEPSEEK_API_KEY" // pragma: allowlist secret
		} else {
			c.DeepSeekAPIKey = ""
		}
	case "ollama":
		if err := os.Setenv("OLLAMA_ENDPOINT", key); err == nil && key != "" { // pragma: allowlist secret
			c.OllamaEndpoint = "OLLAMA_ENDPOINT" // pragma: allowlist secret
		} else {
			c.OllamaEndpoint = ""
		}
	case "claude":
		if err := os.Setenv("CLAUDE_API_KEY", key); err == nil && key != "" { // pragma: allowlist secret
			c.ClaudeAPIKey = "CLAUDE_API_KEY" // pragma: allowlist secret
		} else {
			c.ClaudeAPIKey = ""
		}
	case "gemini":
		if err := os.Setenv("GEMINI_API_KEY", key); err == nil && key != "" { // pragma: allowlist secret
			c.GeminiAPIKey = "GEMINI_API_KEY" // pragma: allowlist secret
		} else {
			c.GeminiAPIKey = ""
		}
	case "chatgpt":
		if err := os.Setenv("CHATGPT_API_KEY", key); err == nil && key != "" { // pragma: allowlist secret
			c.ChatGPTAPIKey = "CHATGPT_API_KEY" // pragma: allowlist secret
		} else {
			c.ChatGPTAPIKey = ""
		}
	default:
		return fmt.Errorf("unknown provider: %s", provider)
	}
	return nil
}

func (c *Config) GetAPIEndpoint(provider string) string {
	if provider == "ollama" {
		return c.OllamaEndpoint
	}
	return ""
}

func (c *Config) CheckOllamaConnection() bool {
	ip, err := netip.ParseAddr(c.OllamaEndpoint)
	if err != nil {
		return false
	}
	conn, err := net.DialTimeout("tcp", ip.String()+":11434", 2*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// GetBaseGenerationPrompt transforms raw ideas into a structured prompt engineering template.
// This method applies professional prompt engineering techniques to convert unorganized ideas
// into well-structured, effective prompts for AI interactions.
func (c *Config) GetBaseGenerationPrompt(ideas []string, purpose, purposeType, lang string, maxLength int) string {
	// Set default values
	if ideas == nil {
		ideas = []string{}
	}
	if lang == "" {
		lang = "english"
	}
	if maxLength <= 0 {
		maxLength = 2000
	}
	if purposeType == "" {
		purposeType = "Code"
	}

	// Build ideas list
	ideasListStr := ""
	for i, idea := range ideas {
		ideasListStr += fmt.Sprintf("%d. \"%s\"\n", i+1, idea)
	}

	// Build purpose text
	purposeText := purposeType
	if purpose != "" {
		purposeText += ` (` + purpose + `)`
	}

	engineeringPrompt := `
Você é um especialista em engenharia de prompts com conhecimento profundo em técnicas de prompt engineering. Sua tarefa é transformar ideias brutas e desorganizadas em um prompt estruturado, profissional e eficaz.

CONTEXTO: O usuário inseriu as seguintes notas/ideias brutas:
` + ideasListStr + `

PROPÓSITO DO PROMPT: ` + purposeText + `
IDIOMA DE RESPOSTA: ` + lang + `
TAMANHO MÁXIMO: ` + fmt.Sprintf("%d", maxLength) + ` caracteres

INSTRUÇÕES PARA ESTRUTURAÇÃO:
1. **Análise**: Identifique o objetivo principal e temas centrais das ideias
2. **Organização**: Estruture logicamente as informações em hierarquia clara
3. **Técnicas de Prompt Engineering**:
   - Definição clara de contexto e papel (persona)
   - Instruções específicas, mensuráveis e testáveis
   - Exemplos concretos quando apropriado
   - Formato de saída bem definido e estruturado
   - Chain-of-thought para raciocínios complexos
   - Few-shot examples se necessário
4. **Formatação**: Use markdown para estruturação visual clara
5. **Tom**: Seja preciso, objetivo, profissional e acionável
6. **Escopo**: Mantenha dentro do limite de caracteres especificado

CRITÉRIOS DE QUALIDADE:
- Clareza: Instruções inequívocas e fáceis de seguir
- Completude: Cubra todos os aspectos relevantes das ideias originais
- Eficácia: Optimize para obter os melhores resultados da IA
- Usabilidade: Pronto para uso imediato sem ajustes

IMPORTANTE: Responda APENAS com o prompt estruturado em markdown, sem explicações adicionais, metadados ou texto introdutório. O prompt deve ser completo, autossuficiente e pronto para uso.
`

	return engineeringPrompt
}

func (c *Config) Validate() error {
	if glgr == nil {
		if c.Logger != nil {
			glgr = c.Logger
		} else {
			glgr = l.LoggerG.GetLogger()
			c.Logger = glgr
		}
	}
	// Example validation: Ensure port is a valid number
	if _, err := net.LookupPort("tcp", c.GetPort()); err != nil {
		return fmt.Errorf("invalid port: %s", c.GetPort())
	}

	// Additional validations can be added here
	if c.GetAPIKey("openai") == "" {
		glgr.Log("warn", "OpenAI API key is not set")
	}
	if c.GetAPIKey("deepseek") == "" {
		glgr.Log("warn", "DeepSeek API key is not set")
	}
	if c.GetAPIKey("ollama") == "" {
		glgr.Log("warn", "Ollama API key is not set")
	}
	if c.GetAPIKey("claude") == "" {
		glgr.Log("warn", "Claude API key is not set")
	}
	if c.GetAPIKey("gemini") == "" {
		glgr.Log("warn", "Gemini API key is not set")
	}
	if c.GetAPIKey("chatgpt") == "" {
		glgr.Log("warn", "ChatGPT API key is not set")
	}
	if c.Defaults == nil {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %v", err)
		}
		c.Defaults = &kbx.InitArgs{
			Address: "localhost",
			Port:    "8080",
			Name:    vs.NewVersionService().GetName(),
			Cwd:     cwd,
			TempDir: os.TempDir(),
			EnvFile: filepath.Join(cwd, ".env"),
			LogFile: filepath.Join(cwd, "config", "logs", "grompt.log"),
			ConfigFile: filepath.Join(cwd, "config", "development.yml"),
			ConfigType: "yaml",
			Debug: 	false,
			IsConfidential: false,
			ReleaseMode: false,
			Bind: "0.0.0.0",
			NotificationTimeoutSeconds: 30,
			NotificationProvider: "email",
			ConfigDBType: "sqlite",
			ConfigDBFile: filepath.Join(cwd, "config", "data", "config.db"),
			PubKeyPath: filepath.Join(cwd, "config", "keys", "pubkey.pem"),
			PrivKeyPath: filepath.Join(cwd, "config", "keys", "privkey.pem"),
			PubCertKeyPath: filepath.Join(cwd, "config", "keys", "pubcert.pem"),
		}
	}
	if c.Server == nil {
		c.Server = newServerConfig()
	}
	if c.Port == "" {
		c.Port = c.Defaults.Port
	}
	if c.BindAddr == "" {
		c.Server.BindAddr = c.Defaults.Bind
	} else {
		c.Server.BindAddr = c.BindAddr
	}
	if c.Server.Name == "" {
		c.Server.Name = c.Defaults.Name
	}
	if c.Server.EnvFile == "" {
		c.Server.EnvFile = c.Defaults.EnvFile
	}
	if c.Server.Debug {
		glgr.SetDebug(true)
		glgr.Log("info", "Debug mode is enabled")
	}else {
		glgr.Log("info", "Debug mode is disabled")
	}
	var tmpDir string
	if c.Server.LogFile == "" {
		c.Server.LogFile = c.Defaults.LogFile
	}
	if c.Server.TempDir == "" {
		tmpDir = c.Defaults.TempDir
	} else {
		tmpDir = c.Server.TempDir
	}

	if  _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		tmpDir, err = os.MkdirTemp("", "grompt-logs")
		if err != nil {
			return fmt.Errorf("failed to create temp dir for logs: %v", err)
		}
	}
	if c.Server.LogFile == "" {
		if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
			return fmt.Errorf("log directory does not exist: %s", tmpDir)
		}
		c.Server.LogFile = tmpDir + "/grompt.log"
	}
	if _, err := os.Stat(tmpDir + "/grompt.log"); os.IsNotExist(err) {
		file, err := os.Create(tmpDir + "/grompt.log")
		if err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
		}
		file.Close()
	}
	if c.Server.Cwd == "" {
		var err error
		c.Server.Cwd, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %v", err)
		}
	}
	var configFilePath string
	if c.Server.ConfigFile == "" {
		configFilePath = filepath.Join(c.Server.Cwd, "config", "development.yml")
	} else {
		configFilePath = c.Server.ConfigFile
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		ConfigObj := c
		mapper := NewMapperW(ConfigObj, configFilePath)
		mapper.SerializeToFile(filepath.Ext(configFilePath)[1:])
		// if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// 	return fmt.Errorf("failed to create default config file: %v", err)
		// }
		glgr.Log("info", fmt.Sprintf("Default config file created at: %s", configFilePath))
	}

	glgr.Log("info", fmt.Sprintf("Server name set to: %s", c.Server.Name))
	glgr.Log("info", fmt.Sprintf("Using config file: %v", configFilePath))
	glgr.Log("info", fmt.Sprintf("Bind address set to: %s", c.BindAddr))
	glgr.Log("info", fmt.Sprintf("Environment set to: %s", c.Server.EnvFile))
	glgr.Log("info", fmt.Sprintf("Log file set to: %s", c.Server.LogFile))
	glgr.Log("info", fmt.Sprintf("Temporary directory set to: %s", tmpDir))
	glgr.Log("info", fmt.Sprintf("Current working directory set to: %s", c.Server.Cwd))

	return nil
}

func (c *Config) GetServerConfig() interfaces.IConfig {
	return c.Server
}

func (c *Config) GetProviders() map[string]interfaces.Provider {
	return c.Providers
}

func (c *Config) GetConfigFilePath() string {
	return c.Server.ConfigFile
}

func (c *Config) IsDebugMode() bool {
	return c.Server.Debug
}

func (c *Config) GetConfigArgs() kbx.InitArgs {
	return *c.Defaults
}

func (c *Config) IsCORSEnabled() bool {
	return c.Server.EnableCORS
}
