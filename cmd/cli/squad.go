package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/rafa-mori/grompt/internal/services/squad"
)

func SquadCmdList() []*cobra.Command {
	return []*cobra.Command{
		generateSquad(),
	}
}

func generateSquad() *cobra.Command {
	var outPath string
	cmd := &cobra.Command{
		Use:   "squad [requirements]",
		Short: "Generate AGENTS.md from free text requirements",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			req := args[0]
			if outPath == "" {
				outPath = "AGENTS.md"
			}
			md, err := squad.BuildAndSave(req, outPath)
			if err != nil {
				return err
			}
			fmt.Println(md)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outPath, "output", "o", "AGENTS.md", "Output file path")

	return cmd
}
