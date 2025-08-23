package module

import (
	"os"
	"strings"
)

func RegX() *Grompt {
	var printBannerV = os.Getenv("GROMPT_PRINT_BANNER")
	if printBannerV == "" {
		printBannerV = "true"
	}

	return &Grompt{
		PrintBanner: strings.ToLower(printBannerV) == "true",
	}
}
