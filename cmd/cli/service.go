package cli

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	gl "github.com/kubex-ecosystem/gemx/grompt/internal/module/logger"
	"github.com/kubex-ecosystem/gemx/grompt/utils"
	"gopkg.in/yaml.v3"

	s "github.com/kubex-ecosystem/gemx/grompt/internal/services/server"
	t "github.com/kubex-ecosystem/gemx/grompt/internal/types"
	l "github.com/kubex-ecosystem/logz"

	"github.com/spf13/cobra"
)

func ServerCmdList() []*cobra.Command {
	return []*cobra.Command{
		startServer(),
	}
}

func startServer() *cobra.Command {
	var debug, background bool
	var configFilePath string
	var bindAddr,
		port,
		openAIKey,
		deepSeekKey,
		ollamaEndpoint,
		claudeKey,
		geminiKey,
		chatGPTKey string

	var startCmd = &cobra.Command{
		Use: "start",
		Annotations: GetDescriptions([]string{
			"This command starts the Grompt server.",
			"This command initializes the Grompt server and starts waiting for help to build prompts.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			gl.GetLogger[l.Logger](nil)
			if debug {
				gl.SetDebug(true)
			}

			if background {
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

			var cfg t.IConfig
			var err error

			if configFilePath != "" {
				cfg, err = loadConfigFile(configFilePath)
				if err != nil {
					gl.Log("fatal", fmt.Sprintf("❌ Erro ao carregar arquivo de configuração: %v", err))
				}
				gl.Log("info", "Arquivo de configuração carregado com sucesso")
			} else {
				cfg = t.NewConfig(
					utils.GetEnvOr("BIND_ADDR", bindAddr),
					utils.GetEnvOr("PORT", port),
					utils.GetEnvOr("OPENAI_API_KEY", openAIKey),
					utils.GetEnvOr("DEEPSEEK_API_KEY", deepSeekKey),
					utils.GetEnvOr("OLLAMA_ENDPOINT", ollamaEndpoint),
					utils.GetEnvOr("CLAUDE_API_KEY", claudeKey),
					utils.GetEnvOr("GEMINI_API_KEY", geminiKey),
					utils.GetEnvOr("CHATGPT_API_KEY", chatGPTKey),
					gl.GetLogger[l.Logger](nil),
				)
			}

			if cfg == nil {
				gl.Log("fatal", "Configuração não carregada")
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

	startCmd.Flags().BoolVarP(&debug, "debug", "D", false, "Enable debug mode")
	startCmd.Flags().StringVarP(&bindAddr, "bind", "b", "localhost", "Address to bind the server to")
	startCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the server on")
	startCmd.Flags().StringVarP(&configFilePath, "config", "f", "", "Path to the config file")
	startCmd.Flags().StringVarP(&openAIKey, "openai-key", "o", "", "OpenAI API key")
	startCmd.Flags().StringVarP(&deepSeekKey, "deepseek-key", "d", "", "DeepSeek API key")
	startCmd.Flags().StringVarP(&ollamaEndpoint, "ollama-endpoint", "e", "http://localhost:11434", "Ollama API endpoint")
	startCmd.Flags().StringVarP(&geminiKey, "gemini-key", "g", "", "Gemini API key")
	startCmd.Flags().StringVarP(&chatGPTKey, "chatgpt-key", "c", "", "ChatGPT API key")
	startCmd.Flags().StringVarP(&claudeKey, "claude-key", "C", "", "Claude API key")
	startCmd.Flags().BoolVarP(&background, "background", "B", false, "Run server in background")

	return startCmd
}

func startServerService() *cobra.Command {
	var debug bool
	var configFilePath string
	var bindAddr,
		port,
		openAIKey,
		deepSeekKey,
		ollamaEndpoint,
		claudeKey,
		geminiKey,
		chatGPTKey string

	var startCmd = &cobra.Command{
		Use: "start",
		Annotations: GetDescriptions([]string{
			"This command starts the Grompt server.",
			"This command initializes the Grompt server and starts waiting for help to build prompts.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			lgr := gl.GetLogger[l.Logger](nil)

			if debug {
				gl.SetDebug(true)
			}

			var cfg t.IConfig
			var err error

			if configFilePath != "" {
				cfg, err = loadConfigFile(configFilePath)
				if err != nil {
					gl.Log("fatal", fmt.Sprintf("❌ Erro ao carregar arquivo de configuração: %v", err))
				}
				gl.Log("info", "Arquivo de configuração carregado com sucesso")
			} else {
				cfg = t.NewConfig(
					utils.GetEnvOr("BIND_ADDR", bindAddr),
					utils.GetEnvOr("PORT", port),
					utils.GetEnvOr("OPENAI_API_KEY", openAIKey),
					utils.GetEnvOr("DEEPSEEK_API_KEY", deepSeekKey),
					utils.GetEnvOr("OLLAMA_ENDPOINT", ollamaEndpoint),
					utils.GetEnvOr("CLAUDE_API_KEY", claudeKey),
					utils.GetEnvOr("GEMINI_API_KEY", geminiKey),
					utils.GetEnvOr("CHATGPT_API_KEY", chatGPTKey),
					lgr,
				)
			}

			if cfg == nil {
				gl.Log("fatal", "Configuração não carregada")
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

	startCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	startCmd.Flags().StringVarP(&bindAddr, "bind", "b", "localhost", "Address to bind the server to")
	startCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the server on")
	startCmd.Flags().StringVarP(&configFilePath, "config", "c", "config.yaml", "Path to the config file")
	startCmd.Flags().StringVarP(&openAIKey, "openai-key", "o", "", "OpenAI API key")
	startCmd.Flags().StringVarP(&deepSeekKey, "deepseek-key", "d", "", "DeepSeek API key")
	startCmd.Flags().StringVarP(&ollamaEndpoint, "ollama-endpoint", "e", "http://localhost:11434", "Ollama API endpoint")
	startCmd.Flags().StringVarP(&claudeKey, "claude-key", "C", "", "Claude API key")
	startCmd.Flags().StringVarP(&geminiKey, "gemini-key", "G", "", "Gemini API key")
	startCmd.Flags().StringVarP(&chatGPTKey, "chatgpt-key", "c", "", "ChatGPT API key")

	return startCmd
}

func loadConfigFile(f string) (*t.Config, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg t.Config
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
