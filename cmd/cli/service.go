package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	i "github.com/kubex-ecosystem/grompt/internal/interfaces"
	s "github.com/kubex-ecosystem/grompt/internal/services/server"
	t "github.com/kubex-ecosystem/grompt/internal/types"

	"github.com/kubex-ecosystem/grompt/internal/module/kbx"
	"github.com/spf13/cobra"

	gl "github.com/kubex-ecosystem/logz"
)

func init() {
	if initArgs == nil {
		initArgs = &kbx.InitArgs{}
	}
}

func ServerCmdList() []*cobra.Command {
	return []*cobra.Command{
		startServer(),
	}
}

func startServer() *cobra.Command {
	var startCmd = &cobra.Command{
		Use: "start",
		Annotations: GetDescriptions([]string{
			"This command starts the Grompt server.",
			"This command initializes the Grompt server and starts waiting for help to build prompts.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			gl := gl.GetLoggerZ("Grompt")
			if initArgs.Debug {
				gl.SetDebugMode(true)
			}

			if initArgs.Background {
				execPath, err := os.Executable()
				if err != nil {
					gl.Log("fatal", fmt.Sprintf("❌ Erro ao obter caminho do executável: %v", err))
				}

				cmdServer := exec.Command(execPath, "service", "start")

				err = cmdServer.Start()
				if err != nil {
					gl.Log("fatal", fmt.Sprintf("❌ Erro ao iniciar processo em segundo plano: %v", err))
				}
				if cmdServer.Process == nil {
					gl.Log("fatal", "❌ Processo em segundo plano não iniciado corretamente")
				}
				pid := cmdServer.Process.Pid
				if err := cmdServer.Process.Release(); err != nil {
					gl.Log("fatal", fmt.Sprintf("❌ Erro ao liberar processo em segundo plano: %v", err))
				}
				gl.Log("info", fmt.Sprintf("✅ Servidor iniciado em segundo plano com PID %d\n", pid))
				return
			}

			var cfg i.IConfig
			var err error

			if initArgs.ConfigFile != "" {
				cfg, err = loadConfigFile(initArgs.ConfigFile)
				if err != nil {
					gl.Log("fatal", fmt.Sprintf("❌ Erro ao carregar arquivo de configuração: %v", err))
				}
				gl.Log("info", "Arquivo de configuração carregado com sucesso")
			} else {
				cfg = getDefaultConfig(initArgs)
			}

			if cfg == nil {
				gl.Log("fatal", "Configuração não carregada")
			} else {
				gl.Log("info", "Configuração carregada com sucesso")
				// Validar configuração
				gl.Log("info", "Validando configuração...")
				if err := cfg.Validate(); err != nil {
					gl.Log("fatal", fmt.Sprintf("❌ Configuração inválida: %v", err))
				}
				gl.Log("info", "Configuração validada com sucesso")
			}

			// Inicializar servidor
			server := s.NewServer(cfg)

			// Graceful shutdown
			go func() {
				c := make(chan os.Signal, 1)
				signal.Notify(c, os.Interrupt, syscall.SIGTERM)
				<-c
				fmt.Println("\n🛑 Encerrando servidor...")
				server.Shutdown()
				os.Exit(0)
			}()

			// Iniciar servidor
			if err := server.Start(); err != nil {
				gl.Log("fatal", fmt.Sprintf("❌ Erro ao iniciar servidor: %v", err))
			}

			gl.Log("success", "Grompt server started successfully")
		},
	}

	startCmd.Flags().BoolVarP(&initArgs.Debug, "debug", "D", false, "Enable debug mode")
	startCmd.Flags().BoolVarP(&initArgs.ReleaseMode, "release", "R", false, "Enable release mode")
	startCmd.Flags().StringVarP(&initArgs.Bind, "bind", "b", "localhost", "Address to bind the server to")
	startCmd.Flags().StringVarP(&initArgs.Port, "port", "p", "8080", "Port to run the server on")
	startCmd.Flags().StringVarP(&initArgs.ConfigFile, "config", "f", "", "Path to the config file")
	startCmd.Flags().StringVarP(&initArgs.EnvFile, "env", "E", "", "Path to the env file")
	startCmd.Flags().StringVarP(&initArgs.OpenAIKey, "openai-key", "o", "", "OpenAI API key")
	startCmd.Flags().StringVarP(&initArgs.DeepSeekKey, "deepseek-key", "d", "", "DeepSeek API key")
	startCmd.Flags().StringVarP(&initArgs.OllamaEndpoint, "ollama-endpoint", "e", "http://localhost:11434", "Ollama API endpoint")
	startCmd.Flags().StringVarP(&initArgs.GeminiKey, "gemini-key", "g", "", "Gemini API key")
	startCmd.Flags().StringVarP(&initArgs.ChatGPTKey, "chatgpt-key", "c", "", "ChatGPT API key")
	startCmd.Flags().StringVarP(&initArgs.ClaudeKey, "claude-key", "C", "", "Claude API key")
	startCmd.Flags().StringVarP(&initArgs.Cwd, "working-dir", "w", "", "Working directory")
	startCmd.Flags().BoolVarP(&initArgs.Background, "background", "B", false, "Run server in background")

	return startCmd
}

// func startServerService() *cobra.Command {
// 	var startCmd = &cobra.Command{
// 		Use: "start",
// 		Annotations: GetDescriptions([]string{
// 			"This command starts the Grompt server.",
// 			"This command initializes the Grompt server and starts waiting for help to build prompts.",
// 		}, false),
// 		Run: func(cmd *cobra.Command, args []string) {
// 			gl := l.LoggerG

// 			if initArgs.Debug {
// 				gl.SetDebug(true)
// 			}

// 			var cfg i.IConfig
// 			var err error

// 			if initArgs.ConfigFile != "" {
// 				cfg, err = loadConfigFile(initArgs.ConfigFile)
// 				if err != nil {
// 					gl.Log("fatal", fmt.Sprintf("❌ Erro ao carregar arquivo de configuração: %v", err))
// 				}
// 				gl.Log("info", "Arquivo de configuração carregado com sucesso")
// 			} else {
// 				cfg = getDefaultConfig(initArgs)
// 			}

// 			if cfg == nil {
// 				gl.Log("fatal", "Configuração não carregada")
// 			}

// 			// Inicializar servidor
// 			server := s.NewServer(cfg)

// 			// Graceful shutdown
// 			go func() {
// 				c := make(chan os.Signal, 1)
// 				signal.Notify(c, os.Interrupt, syscall.SIGTERM)
// 				<-c
// 				fmt.Println("\n🛑 Encerrando servidor...")
// 				server.Shutdown()
// 				os.Exit(0)
// 			}()

// 			// Iniciar servidor
// 			if err := server.Start(); err != nil {
// 				gl.Log("fatal", fmt.Sprintf("❌ Erro ao iniciar servidor: %v", err))
// 			}

// 			gl.Log("success", "Grompt server started successfully")
// 		},
// 	}

// 	startCmd.Flags().BoolVarP(&initArgs.Debug, "debug", "D", false, "Enable debug mode")
// 	startCmd.Flags().BoolVarP(&initArgs.ReleaseMode, "release", "R", false, "Enable release mode")
// 	startCmd.Flags().StringVarP(&initArgs.Bind, "bind", "b", "localhost", "Address to bind the server to")
// 	startCmd.Flags().StringVarP(&initArgs.Port, "port", "p", "8080", "Port to run the server on")
// 	startCmd.Flags().StringVarP(&initArgs.ConfigFile, "config", "f", "", "Path to the config file")
// 	startCmd.Flags().StringVarP(&initArgs.EnvFile, "env", "E", "", "Path to the env file")
// 	startCmd.Flags().StringVarP(&initArgs.OpenAIKey, "openai-key", "o", "", "OpenAI API key")
// 	startCmd.Flags().StringVarP(&initArgs.DeepSeekKey, "deepseek-key", "d", "", "DeepSeek API key")
// 	startCmd.Flags().StringVarP(&initArgs.OllamaEndpoint, "ollama-endpoint", "e", "http://localhost:11434", "Ollama API endpoint")
// 	startCmd.Flags().StringVarP(&initArgs.GeminiKey, "gemini-key", "g", "", "Gemini API key")
// 	startCmd.Flags().StringVarP(&initArgs.ChatGPTKey, "chatgpt-key", "c", "", "ChatGPT API key")
// 	startCmd.Flags().StringVarP(&initArgs.ClaudeKey, "claude-key", "C", "", "Claude API key")
// 	startCmd.Flags().StringVarP(&initArgs.Cwd, "working-dir", "w", "", "Working directory")
// 	startCmd.Flags().BoolVarP(&initArgs.Background, "background", "B", false, "Run server in background")

// 	return startCmd
// }

func loadConfigFile(f string) (*t.Config, error) {
	var cfg = &t.Config{}

	mapper := t.NewMapper(cfg, f)
	format := strings.TrimPrefix(filepath.Ext(f), ".")
	cfgData, err := mapper.DeserializeFromFile(format)
	if err != nil {
		return nil, err
	}

	// Quando DeserializeFromFile retorna nil, tentamos fazer o parse direto dos bytes
	if cfgData == nil {
		fileData, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		cfgData, err = mapper.Deserialize(fileData, format)
		if err != nil {
			return nil, err
		}
	}

	if cfgData == nil {
		return nil, fmt.Errorf("não foi possível carregar a configuração de %s", f)
	}

	// Garante que o caminho do arquivo fique registrado para futuros carregamentos
	if cfgData.Server == nil {
		cfgData.Server = t.NewConfig(
			initArgs.Name,
			initArgs.Debug,
			gl.GetLoggerZ("Grompt"),
			initArgs.Bind,
			initArgs.Port,
			initArgs.TempDir,
			initArgs.LogFile,
			initArgs.EnvFile,
			f,
			initArgs.Cwd,
			initArgs.OpenAIKey,
			initArgs.ClaudeKey,
			initArgs.GeminiKey,
			initArgs.DeepSeekKey,
			initArgs.ChatGPTKey,
			initArgs.OllamaEndpoint,
			make(map[string]string),
			make(map[string]string),
			make(map[string]string),
			make(map[string]string),
			initArgs.DefaultProvider,
			initArgs.DefaultTemperature,
			initArgs.HistorySize,
			initArgs.Timeout,
			initArgs.ProviderConfigPath,
		).GetServerConfig().(*t.ServerConfigImpl)
	} else if cfgData.Server.ConfigFile == "" {
		cfgData.Server.ConfigFile = f
	}

	return cfgData, nil
}
