// Package cli provides common functionality for command-line interface applications.
package cli

import (
	"github.com/kubex-ecosystem/grompt/internal/module/info"
	"github.com/kubex-ecosystem/grompt/internal/module/kbx"
)

func GetDescriptions(descriptionArg []string, hideBanner bool) map[string]string {
	return info.GetDescriptions(descriptionArg, hideBanner)
}


var initArgs *kbx.InitArgs

func init() {
	if initArgs == nil {
		initArgs = &kbx.InitArgs{}
	}
}
