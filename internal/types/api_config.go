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
	"github.com/kubex-ecosystem/grompt/utils"
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
	apiKey     string
	baseURL    string
	version    string
	httpClient *http.Client
	demoMode   bool
}

type Config struct {
	Logger    l.Logger
	Server    *kbx.InitArgs                        `yaml:"server"`
	Defaults  *kbx.InitArgs                        `yaml:"defaults"`
	Providers map[string]interfaces.Provider `yaml:"providers"`

	BindAddr       string `json:"bind_addr,omitempty" gorm:"default:'localhost'"`
	Port           string `json:"port" gorm:"default:8080"`
	OpenAIAPIKey   string `json:"openai_api_key,omitempty" gorm:"default:''"`
	DeepSeekAPIKey string `json:"deepseek_api_key,omitempty" gorm:"default:''"`
	OllamaEndpoint string `json:"ollama_endpoint,omitempty" gorm:"default:'http://localhost:11434'"`
	ClaudeAPIKey   string `json:"claude_api_key,omitempty" gorm:"default:''"`
	GeminiAPIKey   string `json:"gemini_api_key,omitempty" gorm:"default:''"`
	ChatGPTAPIKey  string `json:"chatgpt_api_key,omitempty" gorm:"default:''"`
	Debug          bool   `json:"debug" gorm:"default:false"`
}

func NewConfig(bindAddr, port, openAIKey, deepSeekKey, ollamaEndpoint, claudeKey, geminiKey, chatGPTKey string, logger l.Logger) *Config {
	if logger == nil {
		logger = l.LoggerG
	}
	if bindAddr == "" {
		bindAddr = kbx.DefaultServerHost
	}
	return &Config{
		Logger:         logger,
		BindAddr:       bindAddr,
		Port:           port,
		OpenAIAPIKey:   openAIKey,
		DeepSeekAPIKey: deepSeekKey,
		OllamaEndpoint: ollamaEndpoint,
		ClaudeAPIKey:   claudeKey,
		GeminiAPIKey:   geminiKey,
	}
}

func (c *Config) GetAPIConfig(provider string) interfaces.IAPIConfig {
	if c == nil {
		glgr.Log("error", "Config is nil")
		return nil
	}
	switch provider {
	case "openai":
		return NewOpenAIAPI(c.GetAPIKey("openai"))
	// case "deepseek":
	// 	return NewDeepSeekAPI(c.GetAPIKey("deepseek"))
	case "ollama":
		return NewOllamaAPI(c.GetAPIEndpoint("ollama"))
	// case "claude":
	// 	return NewClaudeAPI(c.GetAPIKey("claude"))
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
		if c.OpenAIAPIKey == "" {
			c.OpenAIAPIKey = utils.GetEnvOr("OPENAI_API_KEY", c.OpenAIAPIKey) // pragma: allowlist secret
		}
		return c.OpenAIAPIKey
	case "deepseek":
		if c.DeepSeekAPIKey == "" {
			c.DeepSeekAPIKey = utils.GetEnvOr("DEEPSEEK_API_KEY", c.DeepSeekAPIKey) // pragma: allowlist secret
		}
		return c.DeepSeekAPIKey
	case "ollama":
		if c.OllamaEndpoint == "" {
			c.OllamaEndpoint = utils.GetEnvOr("OLLAMA_ENDPOINT", c.OllamaEndpoint) // pragma: allowlist secret
		}
		return c.OllamaEndpoint
	case "claude":
		if c.ClaudeAPIKey == "" {
			c.ClaudeAPIKey = utils.GetEnvOr("CLAUDE_API_KEY", c.ClaudeAPIKey) // pragma: allowlist secret
		}
		return c.ClaudeAPIKey
	case "gemini":
		if c.GeminiAPIKey == "" {
			c.GeminiAPIKey = utils.GetEnvOr("GEMINI_API_KEY", c.GeminiAPIKey) // pragma: allowlist secret
		}
		return c.GeminiAPIKey
	case "chatgpt":
		if c.ChatGPTAPIKey == "" {
			c.ChatGPTAPIKey = utils.GetEnvOr("CHATGPT_API_KEY", c.ChatGPTAPIKey) // pragma: allowlist secret
		}
		return c.ChatGPTAPIKey
	default:
		return ""
	}
}

func (c *Config) SetAPIKey(provider string, key string) error {
	switch provider {
	case "openai":
		c.OpenAIAPIKey = key // pragma: allowlist secret
	case "deepseek":
		c.DeepSeekAPIKey = key // pragma: allowlist secret
	case "ollama":
		c.OllamaEndpoint = key // pragma: allowlist secret
	case "claude":
		c.ClaudeAPIKey = key // pragma: allowlist secret
	case "gemini":
		c.GeminiAPIKey = key // pragma: allowlist secret
	case "chatgpt":
		c.ChatGPTAPIKey = key // pragma: allowlist secret
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

func (c *Config) checkOllamaConnection() bool {
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
		pwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %v", err)
		}
		c.Defaults = &kbx.InitArgs{
			Address: "localhost",
			Port:    "8080",
			Name:    vs.NewVersionService().GetName(),
			Pwd:     pwd,
			TempDir: os.TempDir(),
			EnvFile: filepath.Join(pwd, ".env"),
			LogFile: filepath.Join(pwd, "config", "logs", "grompt.log"),
			ConfigFile: filepath.Join(pwd, "config", "development.yml"),
			ConfigType: "yaml",
			Debug: 	false,
			IsConfidential: false,
			ReleaseMode: false,
			Bind: "0.0.0.0",
			NotificationTimeoutSeconds: 30,
			NotificationProvider: "email",
			ConfigDBType: "sqlite",
			ConfigDBFile: filepath.Join(pwd, "config", "data", "config.db"),
			PubKeyPath: filepath.Join(pwd, "config", "keys", "pubkey.pem"),
			PrivKeyPath: filepath.Join(pwd, "config", "keys", "privkey.pem"),
			PubCertKeyPath: filepath.Join(pwd, "config", "keys", "pubcert.pem"),
		}
	}
	if c.Server == nil {
		c.Server = &kbx.InitArgs{}
	}
	if c.Port == "" {
		c.Port = c.Defaults.Port
	}
	if c.BindAddr == "" {
		c.Server.Address = net.JoinHostPort(c.Defaults.Address, c.GetPort())
	} else {
		c.Server.Address = net.JoinHostPort(c.BindAddr, c.GetPort())
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
	if c.Server.LogFile != "" {
		c.Server.LogFile = c.Defaults.LogFile
	}
	if c.Server.LogFile != "" {
		tmpDir = filepath.Dir(c.Server.LogFile)
	} else {
		tmpDir = os.TempDir()
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
	if c.Server.Pwd == "" {
		var err error
		c.Server.Pwd, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %v", err)
		}
	}
	var configFilePath string
	if c.Server.ConfigFile == "" {
		configFilePath = filepath.Join(c.Server.Pwd, "config", "development.yml")
	} else {
		configFilePath = c.Server.ConfigFile
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		if _, err := os.Create(configFilePath); err != nil {
			return fmt.Errorf("failed to create config file: %v", err)
		}
		glgr.Log("info", fmt.Sprintf("Config file not found. Created new config file at: %s", configFilePath))
		// Write default config to the newly created file
		mapper := NewMapper(&c, configFilePath)
		mapper.SerializeToFile("yaml")
		if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
			return fmt.Errorf("failed to create default config file: %v", err)
		}
		glgr.Log("info", fmt.Sprintf("Default config file created at: %s", configFilePath))
	}

	glgr.Log("info", fmt.Sprintf("Server name set to: %s", c.Server.Name))
	glgr.Log("info", fmt.Sprintf("Using config file: %v", configFilePath))
	glgr.Log("info", fmt.Sprintf("Bind address set to: %s", c.BindAddr))
	glgr.Log("info", fmt.Sprintf("Environment set to: %s", c.Server.EnvFile))
	glgr.Log("info", fmt.Sprintf("Log file set to: %s", c.Server.LogFile))
	glgr.Log("info", fmt.Sprintf("Temporary directory set to: %s", tmpDir))
	glgr.Log("info", fmt.Sprintf("Current working directory set to: %s", c.Server.Pwd))

	return nil
}
