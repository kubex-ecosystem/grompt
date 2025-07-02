package main

import (
	cc "github.com/rafa-mori/grompt/cmd/cli"
	gl "github.com/rafa-mori/grompt/logger"
	vs "github.com/rafa-mori/grompt/version"
	"github.com/spf13/cobra"

	"os"
	"strings"
)

type Grompt struct {
	parentCmdName string
	printBanner   bool
}

func (m *Grompt) Alias() string {
	return ""
}
func (m *Grompt) ShortDescription() string {
	return "Grompt a tool for building prompts with AI assistance."
}
func (m *Grompt) LongDescription() string {
	return `Grompt: A tool for building prompts with AI assistance using real engineering practices. Better prompts, better results.., Awesome prompts, AMAZING results !!!`
}
func (m *Grompt) Usage() string {
	return "grompt [command] [args]"
}
func (m *Grompt) Examples() []string {
	return []string{
		"grompt start",
		"grompt stop",
		"grompt status",
	}
}
func (m *Grompt) Active() bool {
	return true
}
func (m *Grompt) Module() string {
	return "grompt"
}
func (m *Grompt) Execute() error {
	return m.Command().Execute()
}
func (m *Grompt) Command() *cobra.Command {
	gl.Log("debug", "Starting Grompt CLI...")

	var rtCmd = &cobra.Command{
		Use:     m.Module(),
		Aliases: []string{m.Alias()},
		Example: m.concatenateExamples(),
		Version: vs.GetVersion(),
		Annotations: cc.GetDescriptions([]string{
			m.LongDescription(),
			m.ShortDescription(),
		}, m.printBanner),
	}

	rtCmd.AddCommand(cc.ServerCmdList()...)
	rtCmd.AddCommand(vs.CliCommand())

	// Set usage definitions for the command and its subcommands
	setUsageDefinition(rtCmd)
	for _, c := range rtCmd.Commands() {
		setUsageDefinition(c)
		if !strings.Contains(strings.Join(os.Args, " "), c.Use) {
			if c.Short == "" {
				c.Short = c.Annotations["description"]
			}
		}
	}

	return rtCmd
}
func (m *Grompt) SetParentCmdName(rtCmd string) {
	m.parentCmdName = rtCmd
}
func (m *Grompt) concatenateExamples() string {
	examples := ""
	rtCmd := m.parentCmdName
	if rtCmd != "" {
		rtCmd = rtCmd + " "
	}
	for _, example := range m.Examples() {
		examples += rtCmd + example + "\n  "
	}
	return examples
}
func RegX() *Grompt {
	var printBannerV = os.Getenv("GROMPT_PRINT_BANNER")
	if printBannerV == "" {
		printBannerV = "true"
	}

	return &Grompt{
		printBanner: strings.ToLower(printBannerV) == "true",
	}
}
