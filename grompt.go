package grompt

import "github.com/spf13/cobra"

// This file/package allows the grompt module to be used as a library.
// It defines the Grompt interface which can be implemented by any module
// that wants to be part of the grompt ecosystem.

type Grompt interface {
	// Alias returns the alias for the command.
	Alias() string
	// ShortDescription returns a brief description of the command.
	ShortDescription() string
	// LongDescription returns a detailed description of the command.
	LongDescription() string
	// Usage returns the usage string for the command.
	Usage() string
	// Examples returns a list of example usages for the command.
	Examples() []string
	// Active returns true if the command is active and should be executed.
	Active() bool
	// Module returns the name of the module.
	Module() string
	// Execute runs the command and returns an error if it fails.
	Execute() error
	// Command returns the cobra.Command associated with this module.
	Command() *cobra.Command
}
